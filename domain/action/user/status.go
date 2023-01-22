package action

import (
	"go-pano/domain/model"
	service "go-pano/domain/service/user"
	"go-pano/utils"

	"github.com/gin-gonic/gin"
)

// 接口
type IStatusAction interface {
	UpdateStatus(ctx *gin.Context)
}

// 實例化
func NewStatusAction(StatusService service.IStatusService) IStatusAction {
	return &StatusAction{StatusService}
}

// class
type StatusAction struct {
	StatusService service.IStatusService
}

// @Summary 更新使用者狀態（status:0 - 已刪除，status:1 - 使用中）
// @Id 6
// @Tags User
// @version 1.0
// @produce application/json
// @param user_id formData int true "使用者ID"
// @param status formData int true "狀態"
// @Success 200 {object} utils.IH200
// @Failure 500 {object} utils.IH500
// @Router /api/user/status [patch]
func (a *StatusAction) UpdateStatus(ctx *gin.Context) {
	// 表單驗證
	var userForm model.UserStatusForm
	if err := ctx.ShouldBindJSON(&userForm); err != nil {
		ctx.JSON(400, utils.H500(err.Error()))
		return
	}

	if err := a.StatusService.UpdateStatus(&userForm); err != nil {
		ctx.JSON(400, utils.H500(err.Error()))
		return
	}

	if userForm.Status == 1 {
		ctx.JSON(200, utils.H200(nil, "復原成功"))
		return
	}
	ctx.JSON(200, utils.H200(nil, "刪除成功"))
}
