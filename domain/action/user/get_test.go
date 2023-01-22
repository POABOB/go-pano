package action

import (
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

// 以下區塊是使用sm去注入GetService來假裝函數
type GetServiceMock struct {
	mock.Mock
}

func (sm *GetServiceMock) Get() ([]model.User, error) {
	args := sm.Called()
	return args.Get(0).([]model.User), args.Error(1)
}

func TestGetAction(test *testing.T) {

	test.Run("成功：獲取所有使用者。", func(test *testing.T) {
		gsm := new(GetServiceMock)

		// 轉成返回格式
		var jsonObject = []model.User{
			{
				UserId:  1,
				Name:    "User1",
				Account: "user1",
				Roles:   []string{"admin"},
				Status:  1,
			},
		}
		expected, _ := json.Marshal(utils.H200(jsonObject, ""))

		// Mock router
		gsm.On("Get").Return(jsonObject, nil)
		router := gin.Default()
		action := NewGetAction(gsm)
		router.GET("/api/user", action.Get)

		// Request
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/api/user", nil)
		router.ServeHTTP(w, req)

		assert.Equal(test, 200, w.Code)
		assert.Equal(test, string(expected), w.Body.String())
		gsm.AssertExpectations(test)
	})

	test.Run("成功：獲取所有使用者，但沒有使用者。", func(test *testing.T) {
		gsm := new(GetServiceMock)

		// 轉成返回格式
		var jsonObject = []model.User{}
		expected, _ := json.Marshal(utils.H200(jsonObject, ""))

		// Mock router
		gsm.On("Get").Return(jsonObject, nil)
		router := gin.Default()
		action := NewGetAction(gsm)
		router.GET("/api/user", action.Get)

		// Request
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/api/user", nil)
		router.ServeHTTP(w, req)

		assert.Equal(test, 200, w.Code)
		assert.Equal(test, string(expected), w.Body.String())
		gsm.AssertExpectations(test)
	})

	test.Run("失敗：獲取所有使用者，但DB有問題。", func(test *testing.T) {
		gsm := new(GetServiceMock)

		// 轉成返回格式
		jsonString := "DB沒有連線"
		expected, _ := json.Marshal(utils.H500(jsonString))

		// Mock router
		gsm.On("Get").Return([]model.User{}, errors.New(jsonString))
		router := gin.Default()
		action := NewGetAction(gsm)
		router.GET("/api/user", action.Get)

		// Request
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/api/user", nil)
		router.ServeHTTP(w, req)

		assert.Equal(test, 400, w.Code)
		assert.Equal(test, string(expected), w.Body.String())
		gsm.AssertExpectations(test)
	})
}
