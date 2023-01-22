package service

import (
	"errors"
	"testing"

	"go-pano/config"
	"go-pano/domain/model"
	"go-pano/utils"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func init() {
	gin.SetMode(gin.TestMode)
	config.LoadConfigTest()
	utils.Reset()
}

func TestUpdateStatusService(test *testing.T) {
	test.Run("成功：成功從Repository更新狀態。", func(test *testing.T) {
		urm := new(UserRepositoryMock)

		jsonObject := &model.UserStatusForm{
			UserId: 1,
			Status: 1,
		}
		// Mock funcs
		urm.On("UpdateStatus", &model.UserStatusForm{
			UserId: 1,
			Status: 1,
		}).Return(nil)

		// 將Ｍock注入真的Service
		UpdateStatusService := NewStatusService(urm)
		err := UpdateStatusService.UpdateStatus(jsonObject)
		assert.NoError(test, err)
		urm.AssertExpectations(test)
	})

	test.Run("失敗：DB出現錯誤。", func(test *testing.T) {
		urm := new(UserRepositoryMock)

		errorString := "DB出現錯誤"
		jsonObject := &model.UserStatusForm{
			UserId: 1,
			Status: 1,
		}
		// Mock funcs
		urm.On("UpdateStatus", &model.UserStatusForm{
			UserId: 1,
			Status: 1,
		}).Return(errors.New(errorString))

		// 將Ｍock注入真的Service
		UpdateStatusService := NewStatusService(urm)
		err := UpdateStatusService.UpdateStatus(jsonObject)
		assert.EqualError(test, err, errorString)
		urm.AssertExpectations(test)
	})
}
