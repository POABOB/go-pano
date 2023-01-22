package service

import (
	"errors"
	"testing"

	"go-pano/config"
	"go-pano/domain/model"
	"go-pano/utils"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func init() {
	gin.SetMode(gin.TestMode)
	config.LoadConfigTest()
	utils.Reset()
}

type UserRepositoryMock struct {
	mock.Mock
}

func (rm *UserRepositoryMock) GetAll() ([]model.User, error) {
	args := rm.Called()
	return args.Get(0).([]model.User), args.Error(1)
}

func (rm *UserRepositoryMock) Update(user *model.UserUpdateForm) error {
	args := rm.Called(user)
	return args.Error(0)
}

func (rm *UserRepositoryMock) UpdateStatus(user *model.UserStatusForm) error {
	args := rm.Called(user)
	return args.Error(0)
}

func (rm *UserRepositoryMock) UpdatePassword(user *model.UserPasswordForm) error {
	args := rm.Called(user)
	return args.Error(0)
}

func (rm *UserRepositoryMock) Create(user *model.User) error {
	args := rm.Called(user)
	return args.Error(0)
}

func TestGetService(test *testing.T) {
	test.Run("成功：成功從Repository獲取資料。", func(test *testing.T) {
		urm := new(UserRepositoryMock)

		jsonObject := []model.User{
			{
				UserId:      1,
				Name:        "User1",
				Account:     "user1",
				RolesString: `["admin"]`,
				Status:      1,
			},
		}
		// Mock funcs
		urm.On("GetAll").Return(jsonObject, nil)

		// 將Ｍock注入真的Service
		GetService := NewGetService(urm)
		expected, err := GetService.Get()
		assert.NoError(test, err)
		assert.Equal(test, expected, []model.User{
			{
				UserId:  1,
				Name:    "User1",
				Account: "user1",
				Roles:   []string{"admin"},
				Status:  1,
			},
		})
		urm.AssertExpectations(test)
	})

	test.Run("成功：成功從Repository獲取資料，但是沒有User。", func(test *testing.T) {
		urm := new(UserRepositoryMock)

		jsonObject := []model.User{}
		// Mock funcs
		urm.On("GetAll").Return(jsonObject, nil)

		// 將Ｍock注入真的Service
		GetService := NewGetService(urm)
		expected, err := GetService.Get()
		assert.NoError(test, err)
		assert.Equal(test, expected, []model.User{})
		urm.AssertExpectations(test)
	})

	test.Run("失敗：DB出現錯誤。", func(test *testing.T) {
		urm := new(UserRepositoryMock)

		jsonString := "DB無法連線"
		// Mock funcs
		urm.On("GetAll").Return([]model.User{}, errors.New(jsonString))

		// 將Ｍock注入真的Service
		GetService := NewGetService(urm)
		expected, err := GetService.Get()
		assert.EqualError(test, err, jsonString)
		assert.Equal(test, expected, []model.User{})
		// 驗證.On()是否真的有被 call 到
		urm.AssertExpectations(test)
	})

	test.Run("失敗：Roles轉[]string出現錯誤。", func(test *testing.T) {
		urm := new(UserRepositoryMock)

		jsonString := "unexpected end of JSON input"
		jsonObject := []model.User{
			{
				UserId:      1,
				Name:        "User1",
				Account:     "user1",
				RolesString: `["admin"`,
				Status:      1,
			},
		}
		// Mock funcs
		urm.On("GetAll").Return(jsonObject, nil)

		// 將Ｍock注入真的Service
		GetService := NewGetService(urm)
		expected, err := GetService.Get()
		assert.EqualError(test, err, jsonString)
		assert.Equal(test, expected, []model.User{})
		urm.AssertExpectations(test)
	})
}
