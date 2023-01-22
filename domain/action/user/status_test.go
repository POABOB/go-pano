package action

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

// 以下區塊是使用sm去注入UpdateStatusService來假裝函數
type UpdateStatusServiceMock struct {
	mock.Mock
}

func (sm *UpdateStatusServiceMock) UpdateStatus(user *model.UserStatusForm) error {
	args := sm.Called(user)
	return args.Error(0)
}

func TestUpdateStatusAction(test *testing.T) {

	test.Run("成功：刪除使用者。", func(test *testing.T) {
		ussm := new(UpdateStatusServiceMock)

		// 轉成返回格式
		var jsonObject = &model.UserStatusForm{
			UserId: 1,
			Status: 0,
		}

		// Mock router
		ussm.On("UpdateStatus", jsonObject).Return(nil)
		router := gin.Default()
		action := NewStatusAction(ussm)
		router.PATCH("/api/user/status", action.UpdateStatus)

		// Mock Body
		jsonBody := []byte(`{
			"user_id":  1,
			"status": 0
		}`)
		body := bytes.NewReader(jsonBody)

		// Request Body
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("PATCH", "/api/user/status", body)
		router.ServeHTTP(w, req)

		assert.Equal(test, 200, w.Code)
		assert.Equal(test, string([]byte(`{"status":true,"msg":"刪除成功","data":null}`)), w.Body.String())
		ussm.AssertExpectations(test)
	})

	test.Run("成功：復原使用者。", func(test *testing.T) {
		ussm := new(UpdateStatusServiceMock)

		// 轉成返回格式
		var jsonObject = &model.UserStatusForm{
			UserId: 1,
			Status: 1,
		}

		// Mock router
		ussm.On("UpdateStatus", jsonObject).Return(nil)
		router := gin.Default()
		action := NewStatusAction(ussm)
		router.PATCH("/api/user/status", action.UpdateStatus)

		// Mock Body
		jsonBody := []byte(`{
			"user_id":  1,
			"status": 1
		}`)
		body := bytes.NewReader(jsonBody)

		// Request Body
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("PATCH", "/api/user/status", body)
		router.ServeHTTP(w, req)

		assert.Equal(test, 200, w.Code)
		assert.Equal(test, string([]byte(`{"status":true,"msg":"復原成功","data":null}`)), w.Body.String())
		ussm.AssertExpectations(test)
	})

	test.Run("失敗：更新裝態，表單為空。", func(test *testing.T) {
		ussm := new(UpdateStatusServiceMock)

		router := gin.Default()
		action := NewStatusAction(ussm)
		router.PATCH("/api/user/status", action.UpdateStatus)

		// Mock Body
		jsonBody := []byte(`{
			"user_id":  0,
			"status": 1
		}`)
		body := bytes.NewReader(jsonBody)

		// Request Body
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("PATCH", "/api/user/status", body)
		router.ServeHTTP(w, req)

		assert.Equal(test, 400, w.Code)
		assert.Equal(test, string([]byte(`{"status":false,"msg":"Key: 'UserStatusForm.UserId' Error:Field validation for 'UserId' failed on the 'required' tag","data":null}`)), w.Body.String())
	})

	test.Run("失敗：更新裝態，Body格式有誤。", func(test *testing.T) {
		ussm := new(UpdateStatusServiceMock)

		router := gin.Default()
		action := NewStatusAction(ussm)
		router.PATCH("/api/user/status", action.UpdateStatus)

		// Mock Body
		jsonBody := []byte(`{
			"user_id":  1,
			"status": 1,
		}`)
		body := bytes.NewReader(jsonBody)

		// Request Body
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("PATCH", "/api/user/status", body)
		router.ServeHTTP(w, req)

		assert.Equal(test, 400, w.Code)
		assert.Equal(test, string([]byte(`{"status":false,"msg":"invalid character '}' looking for beginning of object key string","data":null}`)), w.Body.String())
	})

	test.Run("失敗：更新狀態，DB有問題。", func(test *testing.T) {
		ussm := new(UpdateStatusServiceMock)

		// 轉成返回格式
		var jsonObject = &model.UserStatusForm{
			UserId: 1,
			Status: 1,
		}
		var errorString string = "有錯誤發生"
		expected, _ := json.Marshal(utils.H500(errorString))

		// Mock router
		ussm.On("UpdateStatus", jsonObject).Return(errors.New(errorString))
		router := gin.Default()
		action := NewStatusAction(ussm)
		router.PATCH("/api/user/status", action.UpdateStatus)

		// Mock Body
		jsonBody := []byte(`{
			"user_id":  1,
			"status": 1
		}`)
		body := bytes.NewReader(jsonBody)

		// Request Body
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("PATCH", "/api/user/status", body)
		router.ServeHTTP(w, req)

		assert.Equal(test, 400, w.Code)
		assert.Equal(test, string(expected), w.Body.String())
		ussm.AssertExpectations(test)
	})
}
