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

func TestUpdateService(test *testing.T) {
	test.Run("成功：成功從Repository更新資料。", func(test *testing.T) {
		urm := new(UserRepositoryMock)

		jsonObject := &model.UserUpdateForm{
			UserId:  1,
			Name:    "User1",
			Account: "user1",
			Roles:   []string{"admin"},
		}
		// Mock funcs
		urm.On("Update", &model.UserUpdateForm{
			UserId:      1,
			Name:        "User1",
			Account:     "user1",
			RolesString: `["admin"]`,
			Roles:       []string{"admin"},
		}).Return(nil)

		// 將Ｍock注入真的Service
		UpdateService := NewUpdateService(urm)
		err := UpdateService.Update(jsonObject)
		assert.NoError(test, err)
		urm.AssertExpectations(test)
	})

	test.Run("失敗：DB出現錯誤。", func(test *testing.T) {
		urm := new(UserRepositoryMock)

		errorString := "DB出現錯誤"
		jsonObject := &model.UserUpdateForm{
			UserId:  1,
			Name:    "User1",
			Account: "user1",
			Roles:   []string{"admin"},
		}
		// Mock funcs
		urm.On("Update", &model.UserUpdateForm{
			UserId:      1,
			Name:        "User1",
			Account:     "user1",
			RolesString: `["admin"]`,
			Roles:       []string{"admin"},
		}).Return(errors.New(errorString))

		// 將Ｍock注入真的Service
		UpdateService := NewUpdateService(urm)
		err := UpdateService.Update(jsonObject)
		assert.EqualError(test, err, errorString)
		urm.AssertExpectations(test)
	})
}
