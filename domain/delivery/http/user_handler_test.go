package http

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
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

type UserServiceMock struct {
	mock.Mock
}

func (usm *UserServiceMock) Get() ([]model.User, error) {
	args := usm.Called()
	return args.Get(0).([]model.User), args.Error(1)
}

func (usm *UserServiceMock) Create(user *model.UserCreateForm) error {
	args := usm.Called(user)
	return args.Error(0)
}

func (usm *UserServiceMock) Update(user *model.UserUpdateForm) error {
	args := usm.Called(user)
	return args.Error(0)
}

func (usm *UserServiceMock) UpdatePassword(user *model.UserPasswordForm) error {
	args := usm.Called(user)
	return args.Error(0)
}

func TestUserHandler(test *testing.T) {

	// Get()
	test.Run("成功：Get()，獲取所有使用者。", func(test *testing.T) {
		usm := new(UserServiceMock)

		// 轉成返回格式
		var jsonObject = []model.User{
			{
				UserUpdateForm: model.UserUpdateForm{
					UserId:  1,
					Name:    "User1",
					Account: "user1",
					Roles:   []string{"admin"},
					Status:  1,
				},
			},
		}
		expected, _ := json.Marshal(utils.H200(jsonObject, ""))

		// Mock router
		usm.On("Get").Return(jsonObject, nil)
		router := gin.Default()
		NewUserHandler(router.Group("/api"), usm)

		// Request
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/api/user", nil)
		router.ServeHTTP(w, req)

		assert.Equal(test, 200, w.Code)
		assert.Equal(test, string(expected), w.Body.String())
		usm.AssertExpectations(test)
	})

	test.Run("成功：Get()，獲取所有使用者，但沒有使用者。", func(test *testing.T) {
		usm := new(UserServiceMock)

		// 轉成返回格式
		var jsonObject = []model.User{}
		expected, _ := json.Marshal(utils.H200(jsonObject, ""))

		// Mock router
		usm.On("Get").Return(jsonObject, nil)
		router := gin.Default()
		NewUserHandler(router.Group("/api"), usm)

		// Request
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/api/user", nil)
		router.ServeHTTP(w, req)

		assert.Equal(test, 200, w.Code)
		assert.Equal(test, string(expected), w.Body.String())
		usm.AssertExpectations(test)
	})

	test.Run("失敗：Get()，獲取所有使用者，但DB有問題。", func(test *testing.T) {
		usm := new(UserServiceMock)

		// 轉成返回格式
		jsonString := "DB沒有連線"
		expected, _ := json.Marshal(utils.H500(jsonString))

		// Mock router
		usm.On("Get").Return([]model.User{}, errors.New(jsonString))
		router := gin.Default()
		NewUserHandler(router.Group("/api"), usm)

		// Request
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/api/user", nil)
		router.ServeHTTP(w, req)

		assert.Equal(test, 400, w.Code)
		assert.Equal(test, string(expected), w.Body.String())
		usm.AssertExpectations(test)
	})

	// Create()
	test.Run("成功：Create()，新增使用者。", func(test *testing.T) {
		usm := new(UserServiceMock)

		// 轉成返回格式
		var jsonObject = &model.UserCreateForm{
			Name:     "User1",
			Account:  "user1",
			Password: "ppaass",
			Passconf: "ppaass",
			Roles:    []string{"index-1", "index-2"},
			Status:   1,
		}

		// Mock router
		usm.On("Create", jsonObject).Return(nil)
		router := gin.Default()
		NewUserHandler(router.Group("/api"), usm)

		// Mock Body
		jsonBody := []byte(`{
			"name": "User1",
			"account": "user1",
			"password": "ppaass",
			"passconf": "ppaass",
			"roles": ["index-1","index-2"],
			"status": 1
		}`)
		body := bytes.NewReader(jsonBody)

		// Request Body
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/user", body)
		router.ServeHTTP(w, req)

		assert.Equal(test, 200, w.Code)
		assert.Equal(test, string([]byte(`{"status":true,"msg":"新增成功","data":null}`)), w.Body.String())
		usm.AssertExpectations(test)
	})

	test.Run("失敗：Create()，新增使用者，密碼兩次不一致。", func(test *testing.T) {
		usm := new(UserServiceMock)

		// Mock router
		router := gin.Default()
		NewUserHandler(router.Group("/api"), usm)

		// Mock Body
		jsonBody := []byte(`{
			"name": "User1",
			"account": "user1",
			"password": "ppaass1",
			"passconf": "ppaass2",
			"roles": ["index-1","index-2"],
			"status": 1
		}`)
		body := bytes.NewReader(jsonBody)

		// Request Body
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/user", body)
		router.ServeHTTP(w, req)

		assert.Equal(test, 400, w.Code)
		assert.Equal(test, string([]byte(`{"status":false,"msg":"Key: 'UserCreateForm.Password' Error:Field validation for 'Password' failed on the 'eqfield' tag","data":null}`)), w.Body.String())
		usm.AssertExpectations(test)
	})

	test.Run("失敗：Create()，新增使用者，提交空表單。", func(test *testing.T) {
		usm := new(UserServiceMock)

		// Mock router
		router := gin.Default()
		NewUserHandler(router.Group("/api"), usm)

		// Mock Body
		jsonBody := []byte(`{
			"name": "",
			"account": "",
			"password": "",
			"passconf": "",
			"roles": [],
			"status": 0
		}`)
		body := bytes.NewReader(jsonBody)

		// Request Body
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/user", body)
		router.ServeHTTP(w, req)

		assert.Equal(test, 400, w.Code)
		assert.Equal(test, string([]byte(`{"status":false,"msg":"Key: 'UserCreateForm.Name' Error:Field validation for 'Name' failed on the 'required' tag\nKey: 'UserCreateForm.Account' Error:Field validation for 'Account' failed on the 'required' tag\nKey: 'UserCreateForm.Password' Error:Field validation for 'Password' failed on the 'required' tag\nKey: 'UserCreateForm.Passconf' Error:Field validation for 'Passconf' failed on the 'required' tag","data":null}`)), w.Body.String())
		usm.AssertExpectations(test)
	})

	test.Run("失敗：Create()，新增使用者，但DB有問題。", func(test *testing.T) {
		usm := new(UserServiceMock)

		// 轉成返回格式
		var jsonObject = &model.UserCreateForm{
			Name:     "User1",
			Account:  "user1",
			Password: "ppaass",
			Passconf: "ppaass",
			Roles:    []string{"index-1", "index-2"},
			Status:   1,
		}

		jsonString := "DB沒有連線"
		expected, _ := json.Marshal(utils.H500(jsonString))

		// Mock router
		usm.On("Create", jsonObject).Return(errors.New(jsonString))
		router := gin.Default()
		NewUserHandler(router.Group("/api"), usm)

		// Mock Body
		jsonBody := []byte(`{
			"name": "User1",
			"account": "user1",
			"password": "ppaass",
			"passconf": "ppaass",
			"roles": ["index-1","index-2"],
			"status": 1
		}`)
		body := bytes.NewReader(jsonBody)

		// Request Body
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/user", body)
		router.ServeHTTP(w, req)

		assert.Equal(test, 400, w.Code)
		assert.Equal(test, string(expected), w.Body.String())
		usm.AssertExpectations(test)
	})

	// Update()
	test.Run("成功：Update()，更新使用者。", func(test *testing.T) {
		usm := new(UserServiceMock)

		// 轉成返回格式
		var jsonObject = &model.UserUpdateForm{
			UserId:  1,
			Name:    "User1",
			Account: "user1",
			Roles:   []string{"index-1", "index-2"},
		}

		// Mock router
		usm.On("Update", jsonObject).Return(nil)
		router := gin.Default()
		NewUserHandler(router.Group("/api"), usm)

		// Mock Body
		jsonBody := []byte(`{
			"user_id":  1,
			"name": "User1",
			"account": "user1",
			"roles": ["index-1","index-2"]
		}`)
		body := bytes.NewReader(jsonBody)

		// Request Body
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("PUT", "/api/user", body)
		router.ServeHTTP(w, req)

		assert.Equal(test, 200, w.Code)
		assert.Equal(test, string([]byte(`{"status":true,"msg":"更新成功","data":null}`)), w.Body.String())
		usm.AssertExpectations(test)
	})

	test.Run("失敗：Update()，更新使用者，Roles格式錯誤。", func(test *testing.T) {
		usm := new(UserServiceMock)

		// Mock router
		router := gin.Default()
		NewUserHandler(router.Group("/api"), usm)

		// Mock Body
		jsonBody := []byte(`{
			"user_id":  1,
			"name": "User1",
			"account": "user1",
			"roles": ["index-1","index-2"
		}`)
		body := bytes.NewReader(jsonBody)

		// Request Body
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("PUT", "/api/user", body)
		router.ServeHTTP(w, req)

		assert.Equal(test, 400, w.Code)
		assert.Equal(test, string([]byte(`{"status":false,"msg":"invalid character '}' after array element","data":null}`)), w.Body.String())
	})

	test.Run("失敗：Update()，更新使用者，表單為空。", func(test *testing.T) {
		usm := new(UserServiceMock)

		// Mock router
		router := gin.Default()
		NewUserHandler(router.Group("/api"), usm)

		// Mock Body
		jsonBody := []byte(`{
			"user_id":  0,
			"name": "",
			"account": "",
			"roles": []
		}`)
		body := bytes.NewReader(jsonBody)

		// Request Body
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("PUT", "/api/user", body)
		router.ServeHTTP(w, req)

		assert.Equal(test, 400, w.Code)
		assert.Equal(test, string([]byte(`{"status":false,"msg":"Key: 'UserUpdateForm.UserId' Error:Field validation for 'UserId' failed on the 'required' tag\nKey: 'UserUpdateForm.Name' Error:Field validation for 'Name' failed on the 'required' tag\nKey: 'UserUpdateForm.Account' Error:Field validation for 'Account' failed on the 'required' tag","data":null}`)), w.Body.String())
	})

	test.Run("失敗：Update()，更新使用者，DB有問題。", func(test *testing.T) {
		usm := new(UserServiceMock)

		// 轉成返回格式
		var jsonObject = &model.UserUpdateForm{
			UserId:  1,
			Name:    "User1",
			Account: "user1",
			Roles:   []string{"index-1", "index-2"},
		}
		var errorString string = "有錯誤發生"
		expected, _ := json.Marshal(utils.H500(errorString))

		// Mock router
		usm.On("Update", jsonObject).Return(errors.New(errorString))
		router := gin.Default()
		NewUserHandler(router.Group("/api"), usm)

		// Mock Body
		jsonBody := []byte(`{
			"user_id":  1,
			"name": "User1",
			"account": "user1",
			"roles": ["index-1","index-2"]
		}`)
		body := bytes.NewReader(jsonBody)

		// Request Body
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("PUT", "/api/user", body)
		router.ServeHTTP(w, req)

		assert.Equal(test, 400, w.Code)
		assert.Equal(test, string(expected), w.Body.String())
		usm.AssertExpectations(test)
	})

	// UpdatePassword()
	test.Run("成功：UpdatePassword()，更新密碼。", func(test *testing.T) {
		usm := new(UserServiceMock)

		// 轉成返回格式
		var jsonObject = &model.UserPasswordForm{
			UserId:   1,
			Password: "ppaass",
			Passconf: "ppaass",
		}

		// Mock router
		usm.On("UpdatePassword", jsonObject).Return(nil)
		router := gin.Default()
		NewUserHandler(router.Group("/api"), usm)

		// Mock Body
		jsonBody := []byte(`{
			"user_id":  1,
			"password": "ppaass",
			"passconf": "ppaass"
		}`)
		body := bytes.NewReader(jsonBody)

		// Request Body
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("PATCH", "/api/user/password", body)
		router.ServeHTTP(w, req)

		assert.Equal(test, 200, w.Code)
		assert.Equal(test, string([]byte(`{"status":true,"msg":"更新成功","data":null}`)), w.Body.String())
		usm.AssertExpectations(test)
	})

	test.Run("失敗：UpdatePassword()，更新密碼，密碼兩次不一致。", func(test *testing.T) {
		usm := new(UserServiceMock)

		// Mock router
		router := gin.Default()
		NewUserHandler(router.Group("/api"), usm)

		// Mock Body
		jsonBody := []byte(`{
			"user_id":  1,
			"password": "ppaass1",
			"passconf": "ppaass2"
		}`)
		body := bytes.NewReader(jsonBody)

		// Request Body
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("PATCH", "/api/user/password", body)
		router.ServeHTTP(w, req)

		assert.Equal(test, 400, w.Code)
		assert.Equal(test, string([]byte(`{"status":false,"msg":"Key: 'UserPasswordForm.Password' Error:Field validation for 'Password' failed on the 'eqfield' tag","data":null}`)), w.Body.String())
	})

	test.Run("失敗：UpdatePassword()，更新密碼，表單為空。", func(test *testing.T) {
		usm := new(UserServiceMock)

		router := gin.Default()
		NewUserHandler(router.Group("/api"), usm)

		// Mock Body
		jsonBody := []byte(`{
			"user_id":  0,
			"password": "",
			"passconf": ""
		}`)
		body := bytes.NewReader(jsonBody)

		// Request Body
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("PATCH", "/api/user/password", body)
		router.ServeHTTP(w, req)

		assert.Equal(test, 400, w.Code)
		assert.Equal(test, string([]byte(`{"status":false,"msg":"Key: 'UserPasswordForm.UserId' Error:Field validation for 'UserId' failed on the 'required' tag\nKey: 'UserPasswordForm.Password' Error:Field validation for 'Password' failed on the 'required' tag\nKey: 'UserPasswordForm.Passconf' Error:Field validation for 'Passconf' failed on the 'required' tag","data":null}`)), w.Body.String())
	})

	test.Run("失敗：UpdatePassword()，更新裝態，Body格式有誤。", func(test *testing.T) {
		usm := new(UserServiceMock)

		router := gin.Default()
		NewUserHandler(router.Group("/api"), usm)

		// Mock Body
		jsonBody := []byte(`{
			"user_id":  1,
			"password": "ppaass",
			"passconf": "ppaass",
		}`)
		body := bytes.NewReader(jsonBody)

		// Request Body
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("PATCH", "/api/user/password", body)
		router.ServeHTTP(w, req)

		assert.Equal(test, 400, w.Code)
		assert.Equal(test, string([]byte(`{"status":false,"msg":"invalid character '}' looking for beginning of object key string","data":null}`)), w.Body.String())
	})

	test.Run("失敗：UpdatePassword()，更新狀態，DB有問題。", func(test *testing.T) {
		usm := new(UserServiceMock)

		// 轉成返回格式
		var jsonObject = &model.UserPasswordForm{
			UserId:   1,
			Password: "ppaass",
			Passconf: "ppaass",
		}
		var errorString string = "有錯誤發生"
		expected, _ := json.Marshal(utils.H500(errorString))

		// Mock router
		usm.On("UpdatePassword", jsonObject).Return(errors.New(errorString))
		router := gin.Default()
		NewUserHandler(router.Group("/api"), usm)

		// Mock Body
		jsonBody := []byte(`{
			"user_id":  1,
			"password": "ppaass",
			"passconf": "ppaass"
		}`)
		body := bytes.NewReader(jsonBody)

		// Request Body
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("PATCH", "/api/user/password", body)
		router.ServeHTTP(w, req)

		assert.Equal(test, 400, w.Code)
		assert.Equal(test, string(expected), w.Body.String())
		usm.AssertExpectations(test)
	})
}