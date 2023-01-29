package predict_service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"time"

	"go-pano/config"
	"go-pano/domain/model"
	predict_repository "go-pano/domain/repository/predict"
	pb "go-pano/protos/predict"
	"go-pano/utils"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"gorm.io/gorm"
)

// 預設服務就是Predict
type IPredictService interface {
	Upload(int, *gin.Context) (*model.Predict, error)
}

// 實例化
func NewPredictService(p predict_repository.IPredictRepository, f utils.IFileUtils, r IRequestPredictService) IPredictService {
	return &PredictService{PredictRepository: p, FileUtils: f, RequestPredictService: r}
}

// class
type PredictService struct {
	PredictRepository     predict_repository.IPredictRepository
	FileUtils             utils.IFileUtils
	RequestPredictService IRequestPredictService
}

type IRequestPredictService interface {
	GrpcReq(string, string) (*model.Result, error)
}

func NewRequestPredictService() IRequestPredictService {
	return &RequestPredictService{}
}

type RequestPredictService struct {
}

// TODO 找不到診所 && 驗證TOKEN

// 查找DB，如果沒有資料就處理檔案上傳
func (ps *PredictService) Upload(clinicId int, ctx *gin.Context) (*model.Predict, error) {
	// multipart 獲取檔案
	form, err := ctx.MultipartForm()
	if err != nil {
		return &model.Predict{}, err
	}
	// 確保使用nhicode
	files := form.File["nhicode"]
	if len(files) == 0 {
		return &model.Predict{}, errors.New("請使用'nhicode'作為上傳名稱！")
	}

	predict := &model.Predict{
		ClinicId: clinicId,
		Filename: files[0].Filename,
	}
	// db 查找是否重複上傳
	if p, err := ps.PredictRepository.GetFirstByIDAndFileName(predict); err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return &model.Predict{}, err // 如果不是找不到的錯誤
	} else if err == nil {
		fmt.Println(p.PredictString)
		if err := json.Unmarshal([]byte(p.PredictString), &p.Predict); err != nil {
			return &model.Predict{}, err
		}
		return p, nil // 找到
	}
	// 基本路徑，獲取現在時間，並轉成字符串
	timestamp := time.Now()
	time_dir := fmt.Sprintf("%d", timestamp.UnixMicro())
	base_dir := "./static/img/" + time_dir

	// 判斷是否存在資料夾，不存在則建立
	// err無測試，因為PathExist已經被Mock所以測試也無用
	if err := ps.FileUtils.PathExist(base_dir); err != nil {
		return &model.Predict{}, err
	}

	// 協程上傳
	// err無測試，因為SaveWithGo已經被Mock所以測試也無用
	if err := ps.FileUtils.SaveWithGo(ctx, base_dir, files); err != nil {
		return &model.Predict{}, err
	}

	// Python gRPC
	result, err := ps.RequestPredictService.GrpcReq(time_dir, config.PythonHost+":5001")
	if err != nil {
		os.RemoveAll(base_dir)
		return &model.Predict{}, err
	}

	timeString := timestamp.Format("2006-01-02 15:04:05")
	predict.PredictString = result.Predict
	predict.CreatedAt = timeString
	predict.UpdatedAt = timeString
	predict.Dir = time_dir

	// 插入DB
	predictWithId, err := ps.PredictRepository.Create(predict)
	if err != nil {
		os.RemoveAll(base_dir)
		return &model.Predict{}, err
	}

	if err := json.Unmarshal([]byte(result.Predict), &predictWithId.Predict); err != nil {
		return &model.Predict{}, err
	}

	return predictWithId, nil
}

// GRPC，與 Python Server 溝通
func (rps *RequestPredictService) GrpcReq(dir string, url string) (*model.Result, error) {

	// 建立GRPC連線資訊
	conn, err := grpc.Dial(url, grpc.WithInsecure())
	if err != nil {
		return &model.Result{}, errors.New("gRPC Dial失敗: " + err.Error())
	}
	// 最後請關閉連線
	defer conn.Close()

	// 使用Client連接，並call函數
	client := pb.NewPredictClient(conn)
	pbResult, err := client.Predict(context.Background(), &pb.Request{Dir: dir})
	fmt.Println(dir)
	if err != nil {
		return &model.Result{}, errors.New("Error came from python: " + err.Error())
	}

	if !pbResult.IsSuccessful {
		return &model.Result{}, errors.New("Error came from python: " + pbResult.Msg)
	}

	return &model.Result{
		IsSuccessful: pbResult.IsSuccessful,
		Msg:          pbResult.Msg,
		Predict:      pbResult.Predict,
	}, nil
}
