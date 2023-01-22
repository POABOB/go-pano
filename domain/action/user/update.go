package action

import (
	"go-pano/domain/model"
	service "go-pano/domain/service/user"
	"go-pano/utils"

	"github.com/gin-gonic/gin"
)

// 接口
type IUpdateAction interface {
	Update(ctx *gin.Context)
}

// 實例化
func NewUpdateAction(UpdateService service.IUpdateService) IUpdateAction {
	return &UpdateAction{UpdateService}
}

// class
type UpdateAction struct {
	UpdateService service.IUpdateService
}

// @Summary 更新使用者資訊
// @Id 4
// @Tags User
// @version 1.0
// @produce application/json
// @param name formData string true "名稱"
// @param account formData string true "帳號"
// @param roles formData []string true "權限"
// @Success 200 {object} utils.IH200
// @Failure 500 {object} utils.IH500
// @Router /api/user [put]
func (a *UpdateAction) Update(ctx *gin.Context) {
	// 表單驗證
	var userForm model.UserUpdateForm
	// form-data不能使用ctx.ShouldBind(&userForm)
	if err := ctx.ShouldBindJSON(&userForm); err != nil {
		ctx.JSON(400, utils.H500(err.Error()))
		return
	}

	if err := a.UpdateService.Update(&userForm); err != nil {
		ctx.JSON(400, utils.H500(err.Error()))
		return
	}

	ctx.JSON(200, utils.H200(nil, "更新成功"))
}
