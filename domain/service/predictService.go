package service

import (
	// "bytes"
	// "encoding/binary"
	"encoding/json"
	"fmt"
	"mime/multipart"

	// "os/exec"
	"strconv"
	"sync"
	"time"

	"go-pano/config"
	"go-pano/domain/repository"
	"go-pano/utils"

	"github.com/gin-gonic/gin"
	"github.com/go-resty/resty/v2"
)

// 預設服務就是PredictService
type PredictService struct{}

// python的Response格式
type Result struct {
	IsSuccessful bool        `json:"isSuccessful"`
	Msg          string      `json:"msg"`
	Predict      interface{} `json:"predict"`
}

// // 基本路徑
// var basePath string = "/app/go/static/img/"

// func (ctrl *PredictService) ImageUploadExec(ctx *gin.Context) {

// 	// 檔案處理
// 	string_, fileName, clinic_id, err := processFile(ctx)
// 	if err != nil {
// 		ctx.JSON(500, gin.H{"status": false, "msg": err.Error()})
// 		return
// 	} else if err == nil && string_ == "" {
// 		// DB已經找到資料了
// 		return
// 	}

// 	// 法一：CMD 執行python
// 	cmd := exec.Command("python3", "Start.py", "--Dir", basePath+string_)
// 	cmd.Dir = "/python"
// 	// 添加Buffer儲存輸出
// 	var out bytes.Buffer
// 	var stderr bytes.Buffer
// 	cmd.Stdout = &out
// 	cmd.Stderr = &stderr
// 	if err := cmd.Run(); err != nil {
// 		utils.LogInstance.Error(stderr.String() + err.Error())
// 		ctx.JSON(500, gin.H{"status": false, "msg": stderr.String() + err.Error()})
// 		return
// 	}

// 	// 存取問題
// 	var result Result
// 	if err := binary.Read(&out, binary.BigEndian, &result); err != nil {
// 		utils.LogInstance.Error(err.Error())
// 		ctx.JSON(500, gin.H{"status": false, "msg": err.Error()})
// 		return
// 	}

// 	// 插入DB
// 	var predict repository.Predict
// 	predict.ClinicId = clinic_id
// 	predict.Filename = fileName
// 	p_s, err := json.Marshal(result.Predict)
// 	if err != nil {
// 		utils.LogInstance.Error(result.Msg + err.Error())
// 		ctx.JSON(500, gin.H{"status": false, "msg": result.Msg + err.Error()})
// 		return
// 	}

// 	predict.Predict = string(p_s)
// 	if err := predict.Create(); err != nil {
// 		utils.LogInstance.Error(err.Error())
// 		ctx.JSON(500, gin.H{"error": err.Error()})
// 		return
// 	}

// 	ctx.JSON(200, gin.H{"status": true, "msg": "辨識完成", "predict": out})
// }

// @Summary 上傳圖片和AI辨識
// @Id 1
// @Tags predict
// @version 1.0
// @produce application/json
// @param clinic_id formData int true "請使用診所ID"
// @param nhicode formData file true "請選擇牙齒的X光圖"
// @Success 200 {object} utils.IH200
// @Failure 500 {object} utils.IH500
// @Router /predict [post]
func (ctrl *PredictService) ImageUploadHTTP(ctx *gin.Context) {

	// 檔案處理
	string_, fileName, clinic_id, err := processFile(ctx)
	if err != nil {
		ctx.JSON(500, utils.H500(err.Error()))
		return
	} else if err == nil && string_ == "" {
		// DB已經找到資料了
		return
	}

	// 法三：HTTP，與 Python Server 溝通
	result := &Result{}
	client := resty.New()
	client.R().
		SetResult(&result).
		SetQueryString("Dir=" + string_).
		ForceContentType("application/json").
		Get("http://" + config.PythonHost + ":5000/")
	if !result.IsSuccessful {
		utils.LogInstance.Error("Error came from python: " + result.Msg)
		ctx.JSON(500, utils.H500("Error came from python: "+result.Msg))
		return
	}

	// 插入DB
	var predict repository.Predict
	predict.ClinicId = clinic_id
	predict.Filename = fileName
	p_s, err := json.Marshal(result.Predict)
	if err != nil {
		utils.LogInstance.Error(result.Msg + err.Error())
		ctx.JSON(500, utils.H500(result.Msg+err.Error()))
		return
	}

	predict.Predict = string(p_s)
	if err := predict.Create(); err != nil {
		utils.LogInstance.Error(err.Error())
		ctx.JSON(500, utils.H500(err.Error()))
		return
	}

	ctx.JSON(200, utils.H200(result.Predict, ""))
}

// 使用協程上傳
func saveWithGo(ctx *gin.Context, base_dir string, files []*multipart.FileHeader) error {
	startTime := time.Now().UnixMicro()
	// 處理檔案的channel
	working := make(chan *multipart.FileHeader)
	// 處理錯誤的channel
	failures := make(chan error)

	// 定義waiting Group 避免程式直接跑完
	var wg sync.WaitGroup

	// goroutine worker * 5 分工上傳，然後把結果寫出來
	for i := 0; i < 5; i++ {
		wg.Add(1)
		// worker 處理檔案
		go worker(&wg, working, failures, ctx, base_dir)
	}

	// goroutine 將單個file傳給chan，讓協程自動執行
	for _, file := range files {
		working <- file
	}

	// 關閉 chan
	close(working)
	// 等待 goroutine 執行完畢
	wg.Wait()
	waitingTime := time.Now().UnixMicro() - startTime
	fmt.Println("Upload taking times in", float64(waitingTime)/1000.00, "ms...")

	// 檢查看看有沒有錯誤
	close(failures)
	if err := <-failures; err != nil {
		return err
	}
	return nil
}

// worker上傳檔案，並記錄錯誤
func worker(wg *sync.WaitGroup, working <-chan *multipart.FileHeader, failures chan<- error, ctx *gin.Context, base_dir string) {
	// 最後執行 Done
	defer wg.Done()
	for f := range working {
		err := ctx.SaveUploadedFile(f, base_dir+"/"+f.Filename)
		if err != nil {
			fmt.Println(err.Error())
			failures <- err
		}
	}
}

// 處理檔案上傳
func processFile(ctx *gin.Context) (string, string, int, error) {
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

	// db 方法
	var db repository.Predict
	if err := db.GetFirstByIDAndFileName(clinic_id, files[0].Filename); err == nil {
		var p_j interface{}
		json.Unmarshal([]byte(string(db.Predict)), &p_j)
		ctx.JSON(200, utils.H200(p_j, ""))
		return "", "", 0, nil
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
	if err := saveWithGo(ctx, base_dir, files); err != nil {
		utils.LogInstance.Error(err.Error())
		return err.Error(), "", 0, err
	}

	return time_dir, files[0].Filename, clinic_id, nil
}
