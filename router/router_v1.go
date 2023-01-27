package router_v1

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "go-pano/docs"
	"go-pano/utils"

	"go-pano/domain/delivery/grpc"
	"go-pano/domain/delivery/http"

	clinic_repository "go-pano/domain/repository/clinic"
	predict_repository "go-pano/domain/repository/predict"
	user_repository "go-pano/domain/repository/user"

	clinic_service "go-pano/domain/service/clinic"
	predict_service "go-pano/domain/service/predict"
	user_service "go-pano/domain/service/user"
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
		db := utils.NewDBInstance()

		// PREDICT
		pr := predict_repository.NewPredictRepository(db)
		fu := utils.NewFileUtils()
		rps := &predict_service.RequestPredictService{}
		ps := predict_service.NewPredictService(pr, fu, rps)
		grpc.NewPredictHandler(api, ps)

		// USER
		ur := user_repository.NewUserRepository(db)
		us := user_service.NewUserService(ur)
		http.NewUserHandler(api, us)

		// CLINIC
		cr := clinic_repository.NewClinicRepository(db)
		cs := clinic_service.NewClinicService(cr)
		http.NewClinicHandler(api, cs)

		// TODO STATISTIC
	}

	url := ginSwagger.URL("/swagger/doc.json") // The url pointing to API definition
	app.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))

}
