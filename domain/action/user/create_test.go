package action

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
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

// 以下區塊是使用sm去注入CreateService來假裝函數
type CreateServiceMock struct {
	mock.Mock
}

func (sm *CreateServiceMock) Create(user *model.User) error {
	args := sm.Called(user)
	return args.Error(0)
}

func TestCreateAction(test *testing.T) {

	test.Run("成功：新增使用者。", func(test *testing.T) {
		csm := new(CreateServiceMock)

		// 轉成返回格式
		var jsonObject = &model.User{
			UserId:   0,
			Name:     "User1",
			Account:  "user1",
			Password: "ppaass",
			Passconf: "ppaass",
			Roles:    []string{"index-1", "index-2"},
			Status:   1,
		}

		// Mock router
		csm.On("Create", jsonObject).Return(nil)
		router := gin.Default()
		action := NewCreateAction(csm)
		router.POST("/api/user", action.Create)

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

		fmt.Println(1111)
		fmt.Println(w.Body.String())

		assert.Equal(test, 200, w.Code)
		assert.Equal(test, string([]byte(`{"status":true,"msg":"新增成功","data":null}`)), w.Body.String())
		csm.AssertExpectations(test)
	})

	test.Run("失敗：新增使用者，密碼兩次不一致。", func(test *testing.T) {
		csm := new(CreateServiceMock)

		// Mock router
		router := gin.Default()
		action := NewCreateAction(csm)
		router.POST("/api/user", action.Create)

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
		assert.Equal(test, string([]byte(`{"status":false,"msg":"Key: 'User.Password' Error:Field validation for 'Password' failed on the 'eqfield' tag","data":null}`)), w.Body.String())
		csm.AssertExpectations(test)
	})

	test.Run("失敗：新增使用者，提交空表單。", func(test *testing.T) {
		csm := new(CreateServiceMock)

		// Mock router
		router := gin.Default()
		action := NewCreateAction(csm)
		router.POST("/api/user", action.Create)

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
		assert.Equal(test, string([]byte(`{"status":false,"msg":"Key: 'User.Name' Error:Field validation for 'Name' failed on the 'required' tag\nKey: 'User.Account' Error:Field validation for 'Account' failed on the 'required' tag\nKey: 'User.Password' Error:Field validation for 'Password' failed on the 'required' tag\nKey: 'User.Passconf' Error:Field validation for 'Passconf' failed on the 'required' tag\nKey: 'User.Status' Error:Field validation for 'Status' failed on the 'required' tag","data":null}`)), w.Body.String())
		csm.AssertExpectations(test)
	})

	test.Run("失敗：新增使用者，但DB有問題。", func(test *testing.T) {
		csm := new(CreateServiceMock)

		// 轉成返回格式
		var jsonObject = &model.User{
			UserId:   0,
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
		csm.On("Create", jsonObject).Return(errors.New(jsonString))
		router := gin.Default()
		action := NewCreateAction(csm)
		router.POST("/api/user", action.Create)

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

		fmt.Println(1111)
		fmt.Println(w.Body.String())

		assert.Equal(test, 400, w.Code)
		assert.Equal(test, string(expected), w.Body.String())
		csm.AssertExpectations(test)
	})
}
