package action

import (
	service "go-pano/domain/service/user"
	"go-pano/utils"

	"github.com/gin-gonic/gin"
)

// 接口
type IGetAction interface {
	Get(ctx *gin.Context)
}

// 實例化
func NewGetAction(GetService service.IGetService) IGetAction {
	return &GetAction{GetService}
}

// class
type GetAction struct {
	GetService service.IGetService
}

// @Summary 獲取所有使用者
// @Id 2
// @Tags User
// @version 1.0
// @produce application/json
// @Success 200 {object} utils.IH200{data=[]model.User} "使用者"
// @Failure 500 {object} utils.IH500
// @Router /api/user [get]
func (a *GetAction) Get(ctx *gin.Context) {
	result, err := a.GetService.Get()
	if err != nil {
		ctx.JSON(400, utils.H500(err.Error()))
		return
	}

	ctx.JSON(200, utils.H200(result, ""))
}
