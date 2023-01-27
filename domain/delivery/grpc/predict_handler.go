package grpc

import (
	"fmt"

	"go-pano/domain/model"
	predict_service "go-pano/domain/service/predict"
	"go-pano/utils"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

// 接口
type IPredictHandler interface {
	Upload(ctx *gin.Context)
}

// 實例化
func NewPredictHandler(e *gin.RouterGroup, s predict_service.IPredictService) {
	handler := PredictHandler{PredictService: s}

	// router
	e.POST("/predict", handler.Upload)
}

// class
type PredictHandler struct {
	PredictService predict_service.IPredictService
}

// @Summary GRPC上傳圖片和AI辨識
// @Id 1
// @Tags Predict
// @version 1.0
// @produce application/json
// @param clinic_id formData int true "請使用診所ID"
// @param nhicode formData file true "請選擇牙齒的X光圖"
// @Success 200 {object} utils.IH200{data=model.Predict} "辨識資料"
// @Failure 500 {object} utils.IH500
// @Router /api/predict [post]
// gRPC，與 Python Server 溝通
func (ph *PredictHandler) Upload(ctx *gin.Context) {
	// 表單驗證
	var predictForm model.PredictForm
	// form-data不能使用ctx.ShouldBind(&predictForm)
	if err := ctx.ShouldBindWith(&predictForm, binding.FormMultipart); err != nil {
		ctx.JSON(400, utils.H500(err.Error()))
		return
	}

	// 先查DB確認有無該照片資料，沒有就上傳照片，並且Ai辨識後建立資料紀錄
	predict, err := ph.PredictService.Upload(predictForm.ClinicId, ctx)
	if err != nil {
		utils.LogInstance.Error(err.Error())
		fmt.Println(err.Error())
		ctx.JSON(400, utils.H500(err.Error()))
		return
	}

	ctx.JSON(200, utils.H200(predict, ""))
}
