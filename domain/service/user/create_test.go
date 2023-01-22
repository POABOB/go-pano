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

func TestCreateService(test *testing.T) {
	test.Run("成功：成功從Repository插入資料。", func(test *testing.T) {
		urm := new(UserRepositoryMock)

		jsonObject := &model.User{
			UserId:   1,
			Name:     "User1",
			Account:  "user1",
			Password: "ppaass",
			Passconf: "ppaass",
			Roles:    []string{"admin"},
			Status:   1,
		}
		// Mock funcs
		urm.On("Create", &model.User{
			UserId:      1,
			Name:        "User1",
			Account:     "user1",
			Password:    "ppaass",
			Passconf:    "ppaass",
			RolesString: `["admin"]`,
			Roles:       []string{"admin"},
			Status:      1,
		}).Return(nil)

		// 將Ｍock注入真的Service
		CreateService := NewCreateService(urm)
		err := CreateService.Create(jsonObject)
		assert.NoError(test, err)
		urm.AssertExpectations(test)
	})

	test.Run("失敗：DB出現錯誤。", func(test *testing.T) {
		urm := new(UserRepositoryMock)

		errorString := "DB出現錯誤"
		jsonObject := &model.User{
			UserId:   1,
			Name:     "User1",
			Account:  "user1",
			Password: "ppaass",
			Passconf: "ppaass",
			Roles:    []string{"admin"},
			Status:   1,
		}
		// Mock funcs
		urm.On("Create", &model.User{
			UserId:      1,
			Name:        "User1",
			Account:     "user1",
			Password:    "ppaass",
			Passconf:    "ppaass",
			RolesString: `["admin"]`,
			Roles:       []string{"admin"},
			Status:      1,
		}).Return(errors.New(errorString))

		// 將Ｍock注入真的Service
		CreateService := NewCreateService(urm)
		err := CreateService.Create(jsonObject)
		assert.EqualError(test, err, errorString)
		urm.AssertExpectations(test)
	})
}
