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

func TestUpdatePasswordService(test *testing.T) {
	test.Run("成功：成功從Repository更新密碼。", func(test *testing.T) {
		urm := new(UserRepositoryMock)

		jsonObject := &model.UserPasswordForm{
			UserId:   1,
			Password: "ppaass",
			Passconf: "ppaass",
		}
		// Mock funcs
		urm.On("UpdatePassword", &model.UserPasswordForm{
			UserId:   1,
			Password: "ppaass",
			Passconf: "ppaass",
		}).Return(nil)

		// 將Ｍock注入真的Service
		UpdatePasswordService := NewPasswordService(urm)
		err := UpdatePasswordService.UpdatePassword(jsonObject)
		assert.NoError(test, err)
		urm.AssertExpectations(test)
	})

	test.Run("失敗：DB出現錯誤。", func(test *testing.T) {
		urm := new(UserRepositoryMock)

		errorString := "DB出現錯誤"
		jsonObject := &model.UserPasswordForm{
			UserId:   1,
			Password: "ppaass",
			Passconf: "ppaass",
		}
		// Mock funcs
		urm.On("UpdatePassword", &model.UserPasswordForm{
			UserId:   1,
			Password: "ppaass",
			Passconf: "ppaass",
		}).Return(errors.New(errorString))

		// 將Ｍock注入真的Service
		UpdatePasswordService := NewPasswordService(urm)
		err := UpdatePasswordService.UpdatePassword(jsonObject)
		assert.EqualError(test, err, errorString)
		urm.AssertExpectations(test)
	})
}
