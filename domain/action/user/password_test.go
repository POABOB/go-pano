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

// 以下區塊是使用sm去注入UpdatePasswordService來假裝函數
type UpdatePasswordServiceMock struct {
	mock.Mock
}

func (sm *UpdatePasswordServiceMock) UpdatePassword(user *model.UserPasswordForm) error {
	args := sm.Called(user)
	return args.Error(0)
}

func TestPasswordAction(test *testing.T) {

	test.Run("成功：更新密碼。", func(test *testing.T) {
		upsm := new(UpdatePasswordServiceMock)

		// 轉成返回格式
		var jsonObject = &model.UserPasswordForm{
			UserId:   1,
			Password: "ppaass",
			Passconf: "ppaass",
		}

		// Mock router
		upsm.On("UpdatePassword", jsonObject).Return(nil)
		router := gin.Default()
		action := NewPasswordAction(upsm)
		router.PATCH("/api/user/password", action.UpdatePassword)

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
		upsm.AssertExpectations(test)
	})

	test.Run("失敗：更新密碼，密碼兩次不一致。", func(test *testing.T) {
		upsm := new(UpdatePasswordServiceMock)

		// Mock router
		router := gin.Default()
		action := NewPasswordAction(upsm)
		router.PATCH("/api/user/password", action.UpdatePassword)

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

	test.Run("失敗：更新密碼，表單為空。", func(test *testing.T) {
		upsm := new(UpdatePasswordServiceMock)

		router := gin.Default()
		action := NewPasswordAction(upsm)
		router.PATCH("/api/user/password", action.UpdatePassword)

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

	test.Run("失敗：更新裝態，Body格式有誤。", func(test *testing.T) {
		upsm := new(UpdatePasswordServiceMock)

		router := gin.Default()
		action := NewPasswordAction(upsm)
		router.PATCH("/api/user/password", action.UpdatePassword)

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

	test.Run("失敗：更新狀態，DB有問題。", func(test *testing.T) {
		upsm := new(UpdatePasswordServiceMock)

		// 轉成返回格式
		var jsonObject = &model.UserPasswordForm{
			UserId:   1,
			Password: "ppaass",
			Passconf: "ppaass",
		}
		var errorString string = "有錯誤發生"
		expected, _ := json.Marshal(utils.H500(errorString))

		// Mock router
		upsm.On("UpdatePassword", jsonObject).Return(errors.New(errorString))
		router := gin.Default()
		action := NewPasswordAction(upsm)
		router.PATCH("/api/user/password", action.UpdatePassword)

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
		upsm.AssertExpectations(test)
	})
}
