//go:build wireinject
// +build wireinject

package router_v1

import (
	"github.com/google/wire"

	clinic_repository "go-pano/domain/repository/clinic"
	predict_repository "go-pano/domain/repository/predict"
	user_repository "go-pano/domain/repository/user"

	clinic_service "go-pano/domain/service/clinic"
	predict_service "go-pano/domain/service/predict"
	user_service "go-pano/domain/service/user"
	"go-pano/utils"
)

func initClinicService() clinic_service.IClinicService {
	wire.Build(
		clinic_service.NewClinicService,
		clinic_repository.NewClinicRepository,
		utils.NewDBInstance,
	)
	return nil
}

func initUserService() user_service.IUserService {
	wire.Build(
		user_service.NewUserService,
		user_repository.NewUserRepository,
		utils.NewDBInstance,
	)
	return nil
}
func initPredictService() predict_service.IPredictService {
	wire.Build(
		predict_service.NewPredictService,
		predict_repository.NewPredictRepository,
		predict_service.NewRequestPredictService,
		utils.NewFileUtils,
		utils.NewDBInstance,
	)
	return nil
}
