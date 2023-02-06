package user_service

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

func (urm *UserRepositoryMock) GetAll() ([]model.User, error) {
	args := urm.Called()
	return args.Get(0).([]model.User), args.Error(1)
}

func (urm *UserRepositoryMock) Create(user *model.UserCreateForm) error {
	args := urm.Called(user)
	return args.Error(0)
}

func (urm *UserRepositoryMock) Update(user *model.UserUpdateForm) error {
	args := urm.Called(user)
	return args.Error(0)
}

func (urm *UserRepositoryMock) UpdateAccount(user *model.UserUpdateAccountForm) error {
	args := urm.Called(user)
	return args.Error(0)
}

func (urm *UserRepositoryMock) UpdatePassword(user *model.UserPasswordForm, obj string) error {
	args := urm.Called(user, obj)
	return args.Error(0)
}

func (urm *UserRepositoryMock) Login(user *model.UserLoginForm) (*model.User, error) {
	args := urm.Called(user)
	return args.Get(0).(*model.User), args.Error(1)
}

func (urm *UserRepositoryMock) Delete(user *model.UserDeleteForm) error {
	args := urm.Called(user)
	return args.Error(0)
}

func TestUserService(test *testing.T) {
	// GetAll()
	test.Run("成功：GetAll()，成功從Repository獲取資料。", func(test *testing.T) {
		urm := new(UserRepositoryMock)

		jsonObject := []model.User{
			{
				Name: "User1",
				UserUpdateAccountForm: model.UserUpdateAccountForm{
					UserId:      1,
					Account:     "user1",
					RolesString: `["admin"]`,
					Status:      1,
				},
			},
		}
		// Mock funcs
		urm.On("GetAll").Return(jsonObject, nil)

		// 將Ｍock注入真的Service
		UserService := NewUserService(urm)
		expected, err := UserService.Get()
		assert.NoError(test, err)
		assert.Equal(test, expected, []model.User{
			{
				Name: "User1",
				UserUpdateAccountForm: model.UserUpdateAccountForm{
					UserId:  1,
					Account: "user1",
					Roles:   []string{"admin"},
					Status:  1,
				},
			},
		})
		urm.AssertExpectations(test)
	})

	test.Run("成功：GetAll()，成功從Repository獲取資料，但是沒有User。", func(test *testing.T) {
		urm := new(UserRepositoryMock)

		jsonObject := []model.User{}
		// Mock funcs
		urm.On("GetAll").Return(jsonObject, nil)

		// 將Ｍock注入真的Service
		UserService := NewUserService(urm)
		expected, err := UserService.Get()
		assert.NoError(test, err)
		assert.Equal(test, expected, []model.User{})
		urm.AssertExpectations(test)
	})

	test.Run("失敗：GetAll()，DB出現錯誤。", func(test *testing.T) {
		urm := new(UserRepositoryMock)

		jsonString := "DB無法連線"
		// Mock funcs
		urm.On("GetAll").Return([]model.User{}, errors.New(jsonString))

		// 將Ｍock注入真的Service
		UserService := NewUserService(urm)
		expected, err := UserService.Get()
		assert.EqualError(test, err, jsonString)
		assert.Equal(test, expected, []model.User{})
		// 驗證.On()是否真的有被 call 到
		urm.AssertExpectations(test)
	})

	test.Run("失敗：GetAll()，Roles轉[]string出現錯誤。", func(test *testing.T) {
		urm := new(UserRepositoryMock)

		jsonString := "unexpected end of JSON input"
		jsonObject := []model.User{
			{
				Name: "User1",
				UserUpdateAccountForm: model.UserUpdateAccountForm{
					UserId:      1,
					Account:     "user1",
					RolesString: `["admin"`,
					Status:      1,
				},
			},
		}
		// Mock funcs
		urm.On("GetAll").Return(jsonObject, nil)

		// 將Ｍock注入真的Service
		UserService := NewUserService(urm)
		expected, err := UserService.Get()
		assert.EqualError(test, err, jsonString)
		assert.Equal(test, expected, []model.User{})
		urm.AssertExpectations(test)
	})

	// Create()
	test.Run("成功：Create()，成功從Repository插入資料。", func(test *testing.T) {
		urm := new(UserRepositoryMock)

		jsonObject := &model.UserCreateForm{
			Name:     "User1",
			Account:  "user1",
			Password: "ppaass",
			Passconf: "ppaass",
			Roles:    []string{"admin"},
			Status:   1,
		}
		// Mock funcs
		urm.On("Create", &model.UserCreateForm{
			Name:        "User1",
			Account:     "user1",
			Password:    "ppaass",
			Passconf:    "ppaass",
			RolesString: `["admin"]`,
			Roles:       []string{"admin"},
			Status:      1,
		}).Return(nil)

		// 將Ｍock注入真的Service
		UserService := NewUserService(urm)
		err := UserService.Create(jsonObject)
		assert.NoError(test, err)
		urm.AssertExpectations(test)
	})

	test.Run("失敗：Create()，DB出現錯誤。", func(test *testing.T) {
		urm := new(UserRepositoryMock)

		errorString := "DB出現錯誤"
		jsonObject := &model.UserCreateForm{
			Name:     "User1",
			Account:  "user1",
			Password: "ppaass",
			Passconf: "ppaass",
			Roles:    []string{"admin"},
			Status:   1,
		}
		// Mock funcs
		urm.On("Create", &model.UserCreateForm{
			Name:        "User1",
			Account:     "user1",
			Password:    "ppaass",
			Passconf:    "ppaass",
			RolesString: `["admin"]`,
			Roles:       []string{"admin"},
			Status:      1,
		}).Return(errors.New(errorString))

		// 將Ｍock注入真的Service
		UserService := NewUserService(urm)
		err := UserService.Create(jsonObject)
		assert.EqualError(test, err, errorString)
		urm.AssertExpectations(test)
	})

	// Update()
	test.Run("成功：Update()，成功從Repository更新資料。", func(test *testing.T) {
		urm := new(UserRepositoryMock)

		jsonObject := &model.UserUpdateForm{
			UserId: 1,
			Name:   "User1",
		}
		// Mock funcs
		urm.On("Update", &model.UserUpdateForm{
			UserId: 1,
			Name:   "User1",
		}).Return(nil)

		// 將Ｍock注入真的Service
		UserService := NewUserService(urm)
		err := UserService.Update(jsonObject)
		assert.NoError(test, err)
		urm.AssertExpectations(test)
	})

	test.Run("失敗：Update()，DB出現錯誤。", func(test *testing.T) {
		urm := new(UserRepositoryMock)

		errorString := "DB出現錯誤"
		jsonObject := &model.UserUpdateForm{
			UserId: 1,
			Name:   "User1",
		}
		// Mock funcs
		urm.On("Update", &model.UserUpdateForm{
			UserId: 1,
			Name:   "User1",
		}).Return(errors.New(errorString))

		// 將Ｍock注入真的Service
		UserService := NewUserService(urm)
		err := UserService.Update(jsonObject)
		assert.EqualError(test, err, errorString)
		urm.AssertExpectations(test)
	})

	// Password()
	test.Run("成功：UpdatePassowrd()，成功從Repository更新密碼。", func(test *testing.T) {
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
		}, "all").Return(nil)

		// 將Ｍock注入真的Service
		UpdateUserService := NewUserService(urm)
		err := UpdateUserService.UpdatePassword(jsonObject, "all")
		assert.NoError(test, err)
		urm.AssertExpectations(test)
	})

	test.Run("失敗：UpdatePassowrd()，DB出現錯誤。", func(test *testing.T) {
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
		}, "all").Return(errors.New(errorString))

		// 將Ｍock注入真的Service
		UpdateUserService := NewUserService(urm)
		err := UpdateUserService.UpdatePassword(jsonObject, "all")
		assert.EqualError(test, err, errorString)
		urm.AssertExpectations(test)
	})

	// Login()
	test.Run("成功：Login()，成功從Repository獲取登入資料。", func(test *testing.T) {
		urm := new(UserRepositoryMock)

		jsonObject := &model.UserLoginForm{
			Account:  "user1",
			Password: "ppaass",
		}
		// Mock funcs
		user := &model.User{
			Name: "User1",
			UserUpdateAccountForm: model.UserUpdateAccountForm{
				UserId:      1,
				Account:     "user1",
				RolesString: `["admin"]`,
				Status:      1,
			},
			Password: "hashpassword",
		}
		urm.On("Login", jsonObject).Return(user, nil)

		// 將Ｍock注入真的Service
		UserService := NewUserService(urm)
		token, err := UserService.Login(jsonObject)
		assert.NoError(test, err)
		assert.NotEmpty(test, token)
		urm.AssertExpectations(test)
	})

	test.Run("失敗：Login()，帳號或密碼錯誤。", func(test *testing.T) {
		urm := new(UserRepositoryMock)

		jsonObject := &model.UserLoginForm{
			Account:  "user1",
			Password: "ppaass",
		}
		// Mock funcs
		user := &model.User{}
		urm.On("Login", jsonObject).Return(user, errors.New("帳號或密碼錯誤"))

		// 將Ｍock注入真的Service
		UserService := NewUserService(urm)
		token, err := UserService.Login(jsonObject)
		assert.EqualError(test, err, "帳號或密碼錯誤")
		assert.Empty(test, token)
		urm.AssertExpectations(test)
	})

	test.Run("失敗：Login()，[]string轉Roles失敗。", func(test *testing.T) {
		urm := new(UserRepositoryMock)

		jsonObject := &model.UserLoginForm{
			Account:  "user1",
			Password: "ppaass",
		}
		// Mock funcs
		user := &model.User{
			Name: "User1",
			UserUpdateAccountForm: model.UserUpdateAccountForm{
				UserId:      1,
				Account:     "user1",
				RolesString: `["admin"`,
				Status:      1,
			},
			Password: "hashpassword",
		}
		urm.On("Login", jsonObject).Return(user, nil)

		// 將Ｍock注入真的Service
		UserService := NewUserService(urm)
		token, err := UserService.Login(jsonObject)
		assert.EqualError(test, err, "unexpected end of JSON input")
		assert.Empty(test, token)
		urm.AssertExpectations(test)
	})

}
