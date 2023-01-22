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

// 以下區塊是使用sm去注入UpdateService來假裝函數
type UpdateServiceMock struct {
	mock.Mock
}

func (sm *UpdateServiceMock) Update(user *model.UserUpdateForm) error {
	args := sm.Called(user)
	return args.Error(0)
}

func TestUpdateAction(test *testing.T) {

	test.Run("成功：更新使用者。", func(test *testing.T) {
		usm := new(UpdateServiceMock)

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
		action := NewUpdateAction(usm)
		router.PUT("/api/user", action.Update)

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

	test.Run("失敗：更新使用者，Roles格式錯誤。", func(test *testing.T) {
		usm := new(UpdateServiceMock)

		// Mock router
		router := gin.Default()
		action := NewUpdateAction(usm)
		router.PUT("/api/user", action.Update)

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

	test.Run("失敗：更新使用者，表單為空。", func(test *testing.T) {
		usm := new(UpdateServiceMock)

		// Mock router
		router := gin.Default()
		action := NewUpdateAction(usm)
		router.PUT("/api/user", action.Update)

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

	test.Run("失敗：更新使用者，DB有問題。", func(test *testing.T) {
		usm := new(UpdateServiceMock)

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
		action := NewUpdateAction(usm)
		router.PUT("/api/user", action.Update)

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
}
