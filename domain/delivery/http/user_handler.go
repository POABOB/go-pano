package http

import (
	"go-pano/domain/model"
	user_service "go-pano/domain/service/user"
	"go-pano/middleware"
	"go-pano/utils"

	"github.com/gin-gonic/gin"
)

// 接口
type IUserHandler interface {
	Get(ctx *gin.Context)
	Create(ctx *gin.Context)
	Update(ctx *gin.Context)
	UpdateAccount(ctx *gin.Context)
	UpdatePassword(ctx *gin.Context)
	Login(ctx *gin.Context)
	GetUserInfo(ctx *gin.Context)
	Delete(ctx *gin.Context)
}

// 實例化
func NewUserHandler(e *gin.RouterGroup, s user_service.IUserService) {
	handler := UserHandler{UserService: s}

	// router

	user := e.Group("/user")
	user.POST("/login", handler.Login)
	user.Use(middleware.JWTAuthMiddleware())
	// TODO RBAC
	{
		user.GET("/info", handler.GetUserInfo)
		user.GET("", handler.Get)
		user.POST("", handler.Create)
		user.PATCH("", handler.Update)
		user.PUT("", handler.UpdateAccount)
		user.PATCH("/password", handler.UpdatePassword)
		user.PATCH("/password/self", handler.UpdateSelfPassword)
		user.DELETE("", handler.Delete)
	}
}

// class
type UserHandler struct {
	UserService user_service.IUserService
}

// @Summary 獲取所有使用者
// @Id 2
// @Tags User
// @version 1.0
// @accept application/json
// @produce application/json
// @Security BearerAuth
// @Success 200 {object} utils.IH200{data=[]model.User} "使用者"
// @Failure 500 {object} utils.IH500
// @Router /api/user [get]
func (uh *UserHandler) Get(ctx *gin.Context) {
	result, err := uh.UserService.Get()
	if err != nil {
		ctx.JSON(400, utils.H500(err.Error()))
		return
	}

	ctx.JSON(200, utils.H200(result, ""))
}

// @Summary 新增使用者
// @Id 3
// @Tags User
// @version 1.0
// @accept application/json
// @produce application/json
// @Security BearerAuth
// @param data body model.UserCreateForm true "body"
// @Success 200 {object} utils.IH200
// @Failure 500 {object} utils.IH500
// @Router /api/user [post]
func (uh *UserHandler) Create(ctx *gin.Context) {
	// 表單驗證
	var userForm model.UserCreateForm
	// form-data不能使用ctx.ShouldBind(&userForm)
	if err := ctx.ShouldBindJSON(&userForm); err != nil {
		ctx.JSON(400, utils.H500(err.Error()))
		return
	}

	if err := uh.UserService.Create(&userForm); err != nil {
		ctx.JSON(400, utils.H500(err.Error()))
		return
	}

	ctx.JSON(200, utils.H200(nil, "新增成功"))
}

// @Summary 更新使用者資訊
// @Id 4
// @Tags User
// @version 1.0
// @accept application/json
// @produce application/json
// @Security BearerAuth
// @param data body model.UserUpdateForm true "body"
// @Success 200 {object} utils.IH200
// @Failure 500 {object} utils.IH500
// @Router /api/user [patch]
func (uh *UserHandler) Update(ctx *gin.Context) {
	// 表單驗證
	var userForm model.UserUpdateForm
	// form-data不能使用ctx.ShouldBind(&userForm)
	if err := ctx.ShouldBindJSON(&userForm); err != nil {
		ctx.JSON(400, utils.H500(err.Error()))
		return
	}

	if err := uh.UserService.Update(&userForm); err != nil {
		ctx.JSON(400, utils.H500(err.Error()))
		return
	}

	ctx.JSON(200, utils.H200(nil, "更新成功"))
}

// @Summary 更新使用者帳號資訊
// @Id 15
// @Tags User
// @version 1.0
// @accept application/json
// @produce application/json
// @Security BearerAuth
// @param data body model.UserUpdateAccountForm true "body"
// @Success 200 {object} utils.IH200
// @Failure 500 {object} utils.IH500
// @Router /api/user [put]
func (uh *UserHandler) UpdateAccount(ctx *gin.Context) {
	// 表單驗證
	var userForm model.UserUpdateAccountForm
	if err := ctx.ShouldBindJSON(&userForm); err != nil {
		ctx.JSON(400, utils.H500(err.Error()))
		return
	}

	err := uh.UserService.UpdateAccount(&userForm)
	if err != nil {
		ctx.JSON(400, utils.H500(err.Error()))
		return
	}

	ctx.JSON(200, utils.H200(nil, "更新成功"))
}

// @Summary 更新其他使用者密碼
// @Id 5
// @Tags User
// @version 1.0
// @accept application/json
// @produce application/json
// @Security BearerAuth
// @param data body model.UserPasswordForm true "body"
// @Success 200 {object} utils.IH200
// @Failure 500 {object} utils.IH500
// @Router /api/user/password [patch]
func (uh *UserHandler) UpdatePassword(ctx *gin.Context) {
	// 表單驗證
	var userForm model.UserPasswordForm
	if err := ctx.ShouldBindJSON(&userForm); err != nil {
		ctx.JSON(400, utils.H500(err.Error()))
		return
	}

	if err := uh.UserService.UpdatePassword(&userForm, "all"); err != nil {
		ctx.JSON(400, utils.H500(err.Error()))
		return
	}

	ctx.JSON(200, utils.H200(nil, "更新成功"))
}

// @Summary 更新自己的密碼
// @Id 12
// @Tags User
// @version 1.0
// @accept application/json
// @produce application/json
// @Security BearerAuth
// @param data body model.UserPasswordForm true "body"
// @Success 200 {object} utils.IH200
// @Failure 500 {object} utils.IH500
// @Router /api/user/password/self [patch]
func (uh *UserHandler) UpdateSelfPassword(ctx *gin.Context) {
	// 獲取JWT Claims
	claims, exists := ctx.Get("jwtClaims")
	if !exists {
		ctx.JSON(400, utils.H500("invalid token"))
		return
	}

	// 表單驗證
	var userForm model.UserPasswordForm
	if err := ctx.ShouldBindJSON(&userForm); err != nil {
		ctx.JSON(400, utils.H500(err.Error()))
		return
	}

	// 判斷是不是自己
	userId := claims.(*utils.MyClaims).UserId
	if userId != userForm.UserId {
		ctx.JSON(400, utils.H500("只能更改自己的帳號密碼!"))
		return
	}

	// any轉utils.MyClaims
	userForm.UserId = userId

	if err := uh.UserService.UpdatePassword(&userForm, "self"); err != nil {
		ctx.JSON(400, utils.H500(err.Error()))
		return
	}

	ctx.JSON(200, utils.H200(nil, "更新成功"))
}

// @Summary 使用者登入
// @Id 11
// @Tags User
// @version 1.0
// @accept application/json
// @produce application/json
// @param data body model.UserLoginForm true "body"
// @Success 200 {object} utils.IH200{data=model.UserToken}
// @Failure 500 {object} utils.IH500
// @Router /api/user/login [post]
func (uh *UserHandler) Login(ctx *gin.Context) {
	// 表單驗證
	var userForm model.UserLoginForm
	if err := ctx.ShouldBindJSON(&userForm); err != nil {
		ctx.JSON(400, utils.H500(err.Error()))
		return
	}

	token, err := uh.UserService.Login(&userForm)
	if err != nil {
		ctx.JSON(400, utils.H500(err.Error()))
		return
	}

	ctx.JSON(200, utils.H200(model.UserToken{Token: token}, "登入成功"))
}

// @Summary 使用者JWT資訊
// @Id 13
// @Tags User
// @version 1.0
// @accept application/json
// @produce application/json
// @Security BearerAuth
// @Success 200 {object} utils.IH200{data=utils.MyClaims}
// @Failure 500 {object} utils.IH500
// @Router /api/user/info [get]
func (uh *UserHandler) GetUserInfo(ctx *gin.Context) {
	claims, exists := ctx.Get("jwtClaims")
	if !exists {
		ctx.JSON(400, utils.H500("invalid token"))
		return
	}

	ctx.JSON(200, utils.H200(claims, ""))
}

// @Summary 刪除使用者
// @Id 14
// @Tags User
// @version 1.0
// @accept application/json
// @produce application/json
// @Security BearerAuth
// @Success 200 {object} utils.IH200
// @Failure 500 {object} utils.IH500
// @Router /api/user [delete]
func (uh *UserHandler) Delete(ctx *gin.Context) {
	// 表單驗證
	var userForm model.UserDeleteForm
	if err := ctx.ShouldBindJSON(&userForm); err != nil {
		ctx.JSON(400, utils.H500(err.Error()))
		return
	}

	err := uh.UserService.Delete(&userForm)
	if err != nil {
		ctx.JSON(400, utils.H500(err.Error()))
		return
	}

	ctx.JSON(200, utils.H200(nil, "刪除成功"))
}
