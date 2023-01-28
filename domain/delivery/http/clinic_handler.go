package http

import (
	"go-pano/domain/model"
	clinic_service "go-pano/domain/service/clinic"
	"go-pano/middleware"
	"go-pano/utils"

	"github.com/gin-gonic/gin"
)

// 接口
type IClinicHandler interface {
	Get(ctx *gin.Context)
	Create(ctx *gin.Context)
	Update(ctx *gin.Context)
	UpdateToken(ctx *gin.Context)
}

// 實例化
func NewClinicHandler(e *gin.RouterGroup, s clinic_service.IClinicService) {
	handler := ClinicHandler{ClinicService: s}

	// router
	clinic := e.Group("/clinic")
	clinic.Use(middleware.JWTAuthMiddleware())
	// TODO RBAC
	{
		clinic.GET("", handler.Get)
		clinic.POST("", handler.Create)
		clinic.PUT("", handler.Update)
		clinic.PATCH("/token", handler.UpdateToken)
	}
}

// class
type ClinicHandler struct {
	ClinicService clinic_service.IClinicService
}

// @Summary 獲取所有診所
// @Id 7
// @Tags Clinic
// @version 1.0
// @accept application/json
// @produce application/json
// @Security BearerAuth
// @Success 200 {object} utils.IH200{data=[]model.Clinic} "診所"
// @Failure 500 {object} utils.IH500
// @Router /api/clinic [get]
func (ch *ClinicHandler) Get(ctx *gin.Context) {
	result, err := ch.ClinicService.Get()
	if err != nil {
		ctx.JSON(400, utils.H500(err.Error()))
		return
	}

	ctx.JSON(200, utils.H200(result, ""))
}

// @Summary 新增診所
// @Id 8
// @Tags Clinic
// @version 1.0
// @accept application/json
// @produce application/json
// @Security BearerAuth
// @param data body model.ClinicCreateForm true "body"
// @Success 200 {object} utils.IH200
// @Failure 500 {object} utils.IH500
// @Router /api/clinic [post]
func (ch *ClinicHandler) Create(ctx *gin.Context) {
	// 表單驗證
	var clinicForm model.ClinicCreateForm
	if err := ctx.ShouldBindJSON(&clinicForm); err != nil {
		ctx.JSON(400, utils.H500(err.Error()))
		return
	}

	// 判斷Date是否格式正確
	if !(utils.DateRegex(clinicForm.StartAt) && utils.DateRegex(clinicForm.EndAt)) {
		ctx.JSON(400, utils.H500("Date格式錯誤"))
		return
	}

	if err := ch.ClinicService.Create(&clinicForm); err != nil {
		ctx.JSON(400, utils.H500(err.Error()))
		return
	}

	ctx.JSON(200, utils.H200(nil, "新增成功"))
}

// @Summary 更新診所資訊
// @Id 9
// @Tags Clinic
// @version 1.0
// @accept application/json
// @produce application/json
// @Security BearerAuth
// @param data body model.ClinicUpdateForm true "body"
// @Success 200 {object} utils.IH200
// @Failure 500 {object} utils.IH500
// @Router /api/clinic [put]
func (ch *ClinicHandler) Update(ctx *gin.Context) {
	// 表單驗證
	var clinicForm model.ClinicUpdateForm
	if err := ctx.ShouldBindJSON(&clinicForm); err != nil {
		ctx.JSON(400, utils.H500(err.Error()))
		return
	}

	// 判斷Date是否格式正確
	if !(utils.DateRegex(clinicForm.StartAt) && utils.DateRegex(clinicForm.EndAt)) {
		ctx.JSON(400, utils.H500("Date格式錯誤"))
		return
	}

	if err := ch.ClinicService.Update(&clinicForm); err != nil {
		ctx.JSON(400, utils.H500(err.Error()))
		return
	}

	ctx.JSON(200, utils.H200(nil, "更新成功"))
}

// @Summary 更新診所Token
// @Id 10
// @Tags Clinic
// @version 1.0
// @accept application/json
// @produce application/json
// @Security BearerAuth
// @param data body model.ClinicTokenForm true "body"
// @Success 200 {object} utils.IH200
// @Failure 500 {object} utils.IH500
// @Router /api/clinic/token [patch]
func (ch *ClinicHandler) UpdateToken(ctx *gin.Context) {
	// 表單驗證
	var clinicForm model.ClinicTokenForm
	if err := ctx.ShouldBindJSON(&clinicForm); err != nil {
		ctx.JSON(400, utils.H500(err.Error()))
		return
	}

	if err := ch.ClinicService.UpdateToken(&clinicForm); err != nil {
		ctx.JSON(400, utils.H500(err.Error()))
		return
	}

	ctx.JSON(200, utils.H200(nil, "更新成功"))
}
