package action

import (
	"encoding/json"
	"fmt"

	"go-pano/domain/model"
	service "go-pano/domain/service/predict"
	"go-pano/utils"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

// 接口
type IUploadAction interface {
	Upload(ctx *gin.Context)
}

// 實例化
func NewUploadAction(UploadService service.IUploadService, CreateService service.ICreateService) IUploadAction {
	return &UploadAction{UploadService, CreateService}
}

// class
type UploadAction struct {
	UploadService service.IUploadService
	CreateService service.ICreateService
}

// @Summary HTTP上傳圖片和AI辨識
// @Id 1
// @Tags predict
// @version 1.0
// @produce application/json
// @param clinic_id formData int true "請使用診所ID"
// @param nhicode formData file true "請選擇牙齒的X光圖"
// @Success 200 {object} utils.IH200{data=model.Predict} "辨識資料"
// @Failure 500 {object} utils.IH500
// @Router /api/predict [post]
// HTTP，與 Python Server 溝通
// 上傳檔案Action，負責表單驗證，返回訊息
func (a *UploadAction) Upload(ctx *gin.Context) {
	// 表單驗證
	var predictForm model.PredictForm
	// form-data不能使用ctx.ShouldBind(&predictForm)
	if err := ctx.ShouldBindWith(&predictForm, binding.FormMultipart); err != nil {
		ctx.JSON(400, utils.H500(err.Error()))
		return
	}

	if ctx.Param("method") == "" {
		predictForm.Method = ctx.DefaultPostForm("method", "http")
	} else {
		predictForm.Method = ctx.Param("method")
	}

	// 先查DB確認有無該照片資料，沒有就上傳照片
	dir, fileName, err := a.UploadService.Upload(&predictForm, ctx)
	if err != nil {
		utils.LogInstance.Error(err.Error())
		fmt.Println(err.Error())
		ctx.JSON(400, utils.H500(err.Error()))
		return
	} else if err == nil && fileName != "" {
		// DB已經找到資料了
		var f interface{}
		json.Unmarshal([]byte(string(fileName)), &f)
		ctx.JSON(200, utils.H200(f, ""))
		return
	}

	// 如果圖片資料庫找不到，call CreateService
	result, err := a.CreateService.Create(&predictForm, dir, fileName)
	if err != nil {
		utils.LogInstance.Error(err.Error())
		ctx.JSON(400, utils.H500(err.Error()))
		return
	}

	ctx.JSON(200, utils.H200(result, ""))
}
