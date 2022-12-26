package service

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"go-pano/config"
	"go-pano/domain/model"
	"go-pano/domain/repository"
	pb "go-pano/protos/predict"
	"go-pano/utils"

	"github.com/gin-gonic/gin"
	"github.com/go-resty/resty/v2"
	"google.golang.org/grpc"
)

// 預設服務就是Predict
type IPredictService interface {
	ImageUpload(*gin.Context)
}

// new Class()
func NewPredictService(predictRepository repository.IPredictRepository) IPredictService {
	return &PredictService{predictRepository}
}

// class
type PredictService struct {
	PredictRepository repository.IPredictRepository
}

func (ctrl *PredictService) ImageUpload(ctx *gin.Context) {

	// 檔案處理
	string_, fileName, clinic_id, err := ctrl.processFile(ctx)
	if err != nil {
		ctx.JSON(500, utils.H500(err.Error()))
		return
	} else if err == nil && string_ == "" {
		// DB已經找到資料了
		var p_j interface{}
		json.Unmarshal([]byte(string(fileName)), &p_j)
		ctx.JSON(200, utils.H200(p_j, ""))
		return
	}

	var result *model.Result
	var errstr string
	method := ctx.Param("method")
	if method == "http" {
		// var result pb.Response
		result, errstr = ctrl.httpReq(string_)
	} else if method == "grpc" {
		result, errstr = ctrl.grpcReq(string_)
	} else {
		ctx.JSON(500, utils.H500("參數只接受：http和grpc"))
		return
	}

	if errstr != "" {
		ctx.JSON(500, utils.H500(errstr))
		return
	}

	// JSON TO STRING
	p_s, err := json.Marshal(result.Predict)
	if err != nil {
		utils.LogInstance.Error(result.Msg + err.Error())
		ctx.JSON(500, utils.H500(result.Msg+err.Error()))
		return
	}

	// 插入DB
	// var predict repository.PredictRepository
	if err := ctrl.PredictRepository.Create(clinic_id, fileName, string(p_s)); err != nil {
		utils.LogInstance.Error(err.Error())
		ctx.JSON(500, utils.H500(err.Error()))
		return
	}

	ctx.JSON(200, utils.H200(result.Predict, ""))
}

// @Summary GRPC上傳圖片和AI辨識
// @Id 2
// @Tags predict
// @version 1.0
// @produce application/json
// @param clinic_id formData int true "請使用診所ID"
// @param nhicode formData file true "請選擇牙齒的X光圖"
// @Success 200 {object} utils.IH200
// @Failure 500 {object} utils.IH500
// @Router /predict/grpc [post]
// TEST: GRPC，與 Python Server 溝通
func (ctrl *PredictService) grpcReq(string_ string) (*model.Result, string) {
	// 建立GRPC連線資訊
	conn, err := grpc.Dial(config.PythonHost+":5001", grpc.WithInsecure())
	if err != nil {
		utils.LogInstance.Error("gRPC Dial失敗: " + err.Error())
		return nil, "gRPC Dial失敗: " + err.Error()
	}
	// 最後請關閉連線
	defer conn.Close()

	// 使用Client連接，並call函數
	client := pb.NewPredictClient(conn)
	pbResult, err := client.Predict(context.Background(), &pb.Request{Dir: string_})
	if err != nil {
		utils.LogInstance.Error("Error came from python: " + err.Error())
		return nil, "Error came from python: " + err.Error()
	}

	result := &model.Result{
		IsSuccessful: pbResult.IsSuccessful,
		Msg:          pbResult.Msg,
		Predict:      pbResult.Predict,
	}
	if !result.IsSuccessful {
		utils.LogInstance.Error("Error came from python: " + result.Msg)
		return nil, "Error came from python: " + result.Msg
	}

	return result, ""
}

// @Summary HTTP上傳圖片和AI辨識
// @Id 1
// @Tags predict
// @version 1.0
// @produce application/json
// @param clinic_id formData int true "請使用診所ID"
// @param nhicode formData file true "請選擇牙齒的X光圖"
// @Success 200 {object} utils.IH200
// @Failure 500 {object} utils.IH500
// @Router /predict/http [post]
// HTTP，與 Python Server 溝通
func (ctrl *PredictService) httpReq(string_ string) (*model.Result, string) {
	result := &model.Result{}
	client := resty.New()
	client.R().
		SetResult(&result).
		SetQueryString("Dir=" + string_).
		ForceContentType("application/json").
		Get("http://" + config.PythonHost + ":5000/")
	if !result.IsSuccessful {
		utils.LogInstance.Error("Error came from python: " + result.Msg)
		return nil, result.Msg
	}

	return result, ""
}

// 處理檔案上傳
func (ctrl *PredictService) processFile(ctx *gin.Context) (string, string, int, error) {
	// 判斷是否已有fileName
	clinic_id, err := strconv.Atoi(ctx.PostForm("clinic_id"))
	if err != nil {
		utils.LogInstance.Error(err.Error())
		return err.Error(), "", 0, err
	}

	// multipart 獲取檔案
	form, err := ctx.MultipartForm()
	if err != nil {
		utils.LogInstance.Error(err.Error())
		return err.Error(), "", 0, err
	}
	// 確保使用nhicode
	files := form.File["nhicode"]
	if len(files) == 0 {
		// 找不到檔案
		utils.LogInstance.Error("請使用'nhicode'作為上傳名稱！")
		return "請使用'nhicode'作為上傳名稱！", "", 0, err
	}

	// // db 方法
	// var db repository.PredictRepository
	if p, err := ctrl.PredictRepository.GetFirstByIDAndFileName(clinic_id, files[0].Filename); err == nil {
		return "", p.Predict, 0, nil
	}

	// 基本路徑，獲取現在時間，並轉成字符串
	time_dir := fmt.Sprintf("%d", time.Now().UnixMicro())
	base_dir := "./static/img/" + time_dir
	// 判斷是否存在資料夾，不存在則建立
	if err := utils.PathExist(base_dir); err != nil {
		utils.LogInstance.Error(err.Error())
		return err.Error(), "", 0, err
	}

	// 協程上傳
	if err := utils.SaveWithGo(ctx, base_dir, files); err != nil {
		utils.LogInstance.Error(err.Error())
		return err.Error(), "", 0, err
	}

	return time_dir, files[0].Filename, clinic_id, nil
}
