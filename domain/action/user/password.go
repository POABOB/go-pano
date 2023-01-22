package action

import (
	"go-pano/domain/model"
	service "go-pano/domain/service/user"
	"go-pano/utils"

	"github.com/gin-gonic/gin"
)

// 接口
type IPasswordAction interface {
	UpdatePassword(ctx *gin.Context)
}

// 實例化
func NewPasswordAction(PasswordService service.IPasswordService) IPasswordAction {
	return &PasswordAction{PasswordService}
}

// class
type PasswordAction struct {
	PasswordService service.IPasswordService
}

// @Summary 更新使用者密碼
// @Id 5
// @Tags User
// @version 1.0
// @produce application/json
// @param user_id formData int true "使用者ID"
// @param password formData string true "密碼"
// @param passconf formData string true "密碼確認"
// @Success 200 {object} utils.IH200
// @Failure 500 {object} utils.IH500
// @Router /api/user/password [patch]
func (a *PasswordAction) UpdatePassword(ctx *gin.Context) {
	// 表單驗證
	var userForm model.UserPasswordForm
	if err := ctx.ShouldBindJSON(&userForm); err != nil {
		ctx.JSON(400, utils.H500(err.Error()))
		return
	}

	if err := a.PasswordService.UpdatePassword(&userForm); err != nil {
		ctx.JSON(400, utils.H500(err.Error()))
		return
	}

	ctx.JSON(200, utils.H200(nil, "更新成功"))
}
