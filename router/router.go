package router

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "go-pano/docs"
	predict_action "go-pano/domain/action/predict"
	user_action "go-pano/domain/action/user"
	// clinic_action "go-pano/domain/action/clinic"
)

// @title Gin Go Pano
// @version 1.0
// @description Swagger API.
// @host localhost
func NewRouter(app *gin.Engine) {

	// TODO JWT MIDDLEWARE
	// TODO RBAC
	api := app.Group("/api")
	{
		// 路由分組
		predict := api.Group("/predict")
		{
			UploadAction := new(predict_action.UploadAction)
			// 預設Http
			predict.POST("", UploadAction.Upload)
			// OPTION Http || gRPC
			predict.POST("/*method", UploadAction.Upload)
		}

		user := api.Group("/user")
		{
			// 預設Http
			user.GET("", new(user_action.GetAction).Get)
			user.POST("", new(user_action.CreateAction).Create)
			user.PUT("", new(user_action.UpdateAction).Update)
			user.PATCH("/status", new(user_action.StatusAction).UpdateStatus)
			user.PATCH("/password", new(user_action.PasswordAction).UpdatePassword)
		}

		// TODO CLINIC
		// clinic := api.Group("/clinic")
		// {
		// 預設Http
		// clinic.GET("", new(clinic_action.GetAction).Get)
		// clinic.POST("", new(clinic_action.InsertAction).Insert)
		// clinic.PUT("", new(clinic_action.UpdateAction).Update)
		// clinic.DELETE("", new(clinic_action.DeleteAction).Delete)
		// }

		// TODO STATISTIC
	}

	url := ginSwagger.URL("/swagger/doc.json") // The url pointing to API definition
	app.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))

}
