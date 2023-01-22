package action

import (
	"go-pano/domain/model"
	service "go-pano/domain/service/user"
	"go-pano/utils"

	"github.com/gin-gonic/gin"
)

// 接口
type ICreateAction interface {
	Create(ctx *gin.Context)
}

// 實例化
func NewCreateAction(CreateService service.ICreateService) ICreateAction {
	return &CreateAction{CreateService}
}

// class
type CreateAction struct {
	CreateService service.ICreateService
}

// @Summary 新增使用者
// @Id 3
// @Tags User
// @version 1.0
// @produce application/json
// @param name formData string true "名稱"
// @param account formData string true "帳號"
// @param password formData string true "密碼"
// @param passconf formData string true "密碼確認"
// @param roles formData []string true "權限"
// @param status formData int true "狀態"
// @Success 200 {object} utils.IH200
// @Failure 500 {object} utils.IH500
// @Router /api/user [post]
func (a *CreateAction) Create(ctx *gin.Context) {
	// 表單驗證
	var userForm model.User
	// form-data不能使用ctx.ShouldBind(&userForm)
	if err := ctx.ShouldBindJSON(&userForm); err != nil {
		ctx.JSON(400, utils.H500(err.Error()))
		return
	}

	if err := a.CreateService.Create(&userForm); err != nil {
		ctx.JSON(400, utils.H500(err.Error()))
		return
	}

	ctx.JSON(200, utils.H200(nil, "新增成功"))
}
