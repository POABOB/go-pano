package router

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "go-pano/docs"
	"go-pano/domain/service"
)

// @title Gin Go Pano
// @version 1.0
// @description Swagger API.
// @host localhost
func NewRouter(app *gin.Engine) {
	// 路由分組
	predictService := new(service.PredictService)
	predict := app.Group("/predict")
	{
		predict.POST("/grpc", predictService.ImageUpload)
	}

	url := ginSwagger.URL("/swagger/doc.json") // The url pointing to API definition
	app.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))

}
