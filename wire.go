//go:build wireinject
// +build wireinject

package main

// TODO 依賴注入


// import (
// 	"github.com/google/wire"
// 	clinic_action "go-pano/domain/action/clinic"
// 	clinic_repository "go-pano/domain/repository/clinic"
// 	clinic_service "go-pano/domain/service/clinic"
// 	"go-pano/utils"
// 	// user_action "go-pano/domain/action/user"
// 	// user_repository "go-pano/domain/repository/user"
// 	// user_service "go-pano/domain/service/user"
// 	// predict_action "go-pano/domain/action/predict"
// 	// predict_repository "go-pano/domain/repository/predict"
// 	// predict_service "go-pano/domain/service/predict"
// )

// func ClinicGetService() *clinic_action.GetAction {
// 	panic(wire.Build(
// 		clinic_action.NewGetAction,
// 		clinic_service.NewGetService,
// 		clinic_repository.NewClinicRepository,
// 		utils.NewDBInstance,
// 	))
// 	// return &clinic_action.GetAction{}
// }
