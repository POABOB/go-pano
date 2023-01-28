package router_v1

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "go-pano/docs"

	"go-pano/domain/delivery/grpc"
	"go-pano/domain/delivery/http"
)

// @title Gin Go Pano
// @version 1.0
// @description Swagger API.
// @host localhost
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func NewRouter(app *gin.Engine) {

	api := app.Group("/api")
	{
		// PREDICT，使用wire
		grpc.NewPredictHandler(api, initPredictService())

		// USER，使用wire
		http.NewUserHandler(api, initUserService())

		// CLINIC，使用wire
		http.NewClinicHandler(api, initClinicService())

		// TODO STATISTIC
	}

	url := ginSwagger.URL("/swagger/doc.json") // The url pointing to API definition
	app.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))

}
