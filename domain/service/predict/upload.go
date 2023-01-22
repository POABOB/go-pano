package service

import (
	"errors"
	"fmt"
	"time"

	"go-pano/domain/model"
	repository "go-pano/domain/repository/predict"
	"go-pano/utils"

	"github.com/gin-gonic/gin"
)

// 預設服務就是Predict
type IUploadService interface {
	Upload(*model.PredictForm, *gin.Context) (string, string, error)
}

// 實例化
func NewUploadService(predictRepository repository.IPredictRepository, fileUtils utils.IFileUtils) IUploadService {
	return &UploadService{predictRepository, fileUtils}
}

// class
type UploadService struct {
	PredictRepository repository.IPredictRepository
	FileUtils         utils.IFileUtils
}

// 查找DB，如果沒有資料就處理檔案上傳
func (s *UploadService) Upload(predictForm *model.PredictForm, ctx *gin.Context) (string, string, error) {
	// multipart 獲取檔案
	form, err := ctx.MultipartForm()
	if err != nil {
		return "", "", err
	}
	// 確保使用nhicode
	files := form.File["nhicode"]
	if len(files) == 0 {
		return "", "", errors.New("請使用'nhicode'作為上傳名稱！")
	}

	// db 查找是否重複上傳
	if p, err := s.PredictRepository.GetFirstByIDAndFileName(predictForm.ClinicId, files[0].Filename); err == nil {
		return "", p.PredictString, nil
	}

	// 基本路徑，獲取現在時間，並轉成字符串
	time_dir := fmt.Sprintf("%d", time.Now().UnixMicro())
	base_dir := "./static/img/" + time_dir

	// 判斷是否存在資料夾，不存在則建立
	// err無測試，因為PathExist已經被Mock所以測試也無用
	if err := s.FileUtils.PathExist(base_dir); err != nil {
		return "", "", err
	}

	// 協程上傳
	// err無測試，因為SaveWithGo已經被Mock所以測試也無用
	if err := s.FileUtils.SaveWithGo(ctx, base_dir, files); err != nil {
		return "", "", err
	}

	return time_dir, files[0].Filename, nil
}
