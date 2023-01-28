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

type ClinicServiceMock struct {
	mock.Mock
}

func (csm *ClinicServiceMock) Get() ([]model.Clinic, error) {
	args := csm.Called()
	return args.Get(0).([]model.Clinic), args.Error(1)
}

func (csm *ClinicServiceMock) Create(c *model.ClinicCreateForm) error {
	args := csm.Called(c)
	return args.Error(0)
}

func (csm *ClinicServiceMock) Update(c *model.ClinicUpdateForm) error {
	args := csm.Called(c)
	return args.Error(0)
}

func (csm *ClinicServiceMock) UpdateToken(c *model.ClinicTokenForm) error {
	args := csm.Called(c)
	return args.Error(0)
}

func TestClinicHandler(test *testing.T) {
	// 100年
	jwtToken := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoxLCJuYW1lIjoiQURNSU4iLCJyb2xlcyI6WyJhZG1pbiJdLCJleHAiOjQ4Mjg0ODI4MDAsImlzcyI6Imdpbi1nby1zZXJ2ZXIifQ.O7rDIkQ8eo7VOevkzKgXmfLoMKUYuRwVRg5t12JyImg"

	// Get
	test.Run("成功：Get()，獲取所有診所。", func(test *testing.T) {
		csm := new(ClinicServiceMock)

		// 轉成返回格式
		var jsonObject = []model.Clinic{
			{
				ClinicTokenForm: model.ClinicTokenForm{
					ClinicId: 1,
				},
				ClinicCreateForm: model.ClinicCreateForm{
					Name:          "診所",
					StartAt:       "2022-10-31",
					EndAt:         "9999-12-31",
					QuotaPerMonth: 100,
				},
				Token: "token",
			},
		}
		expected, _ := json.Marshal(utils.H200(jsonObject, ""))

		// Mock router
		csm.On("Get").Return(jsonObject, nil)
		router := gin.Default()
		NewClinicHandler(router.Group("/api"), csm)

		// Request
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/api/clinic", nil)
		req.Header.Set("Authorization", "Bearer "+jwtToken)
		router.ServeHTTP(w, req)

		assert.Equal(test, 200, w.Code)
		assert.Equal(test, string(expected), w.Body.String())
		csm.AssertExpectations(test)
	})

	test.Run("成功：Get()，獲取所有診所，但沒有診所。", func(test *testing.T) {
		csm := new(ClinicServiceMock)

		// 轉成返回格式
		var jsonObject = []model.Clinic{}
		expected, _ := json.Marshal(utils.H200(jsonObject, ""))

		// Mock router
		csm.On("Get").Return(jsonObject, nil)
		router := gin.Default()
		NewClinicHandler(router.Group("/api"), csm)

		// Request
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/api/clinic", nil)
		req.Header.Set("Authorization", "Bearer "+jwtToken)
		router.ServeHTTP(w, req)

		assert.Equal(test, 200, w.Code)
		assert.Equal(test, string(expected), w.Body.String())
		csm.AssertExpectations(test)
	})

	test.Run("失敗：Get()，獲取所有診所，但DB有問題。", func(test *testing.T) {
		csm := new(ClinicServiceMock)

		// 轉成返回格式
		jsonString := "DB沒有連線"
		expected, _ := json.Marshal(utils.H500(jsonString))

		// Mock router
		csm.On("Get").Return([]model.Clinic{}, errors.New(jsonString))
		router := gin.Default()
		NewClinicHandler(router.Group("/api"), csm)

		// Request
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/api/clinic", nil)
		req.Header.Set("Authorization", "Bearer "+jwtToken)
		router.ServeHTTP(w, req)

		assert.Equal(test, 400, w.Code)
		assert.Equal(test, string(expected), w.Body.String())
		csm.AssertExpectations(test)
	})

	// Create
	test.Run("成功：Create()，新增診所。", func(test *testing.T) {
		csm := new(ClinicServiceMock)

		// 轉成返回格式
		var jsonObject = &model.ClinicCreateForm{
			Name:          "診所",
			StartAt:       "2022-10-31",
			EndAt:         "9999-12-31",
			QuotaPerMonth: 100,
		}

		// Mock router
		csm.On("Create", jsonObject).Return(nil)
		router := gin.Default()
		NewClinicHandler(router.Group("/api"), csm)

		// Mock Body
		jsonBody := []byte(`{
			"name":         "診所",
			"start_at":       "2022-10-31",
			"end_at":         "9999-12-31",
			"quota_per_month": 100
		}`)
		body := bytes.NewReader(jsonBody)

		// Request Body
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/clinic", body)
		req.Header.Set("Authorization", "Bearer "+jwtToken)
		router.ServeHTTP(w, req)

		assert.Equal(test, 200, w.Code)
		assert.Equal(test, string([]byte(`{"status":true,"msg":"新增成功","data":null}`)), w.Body.String())
		csm.AssertExpectations(test)
	})

	test.Run("失敗：Create()，新增診所，Date格式錯誤。", func(test *testing.T) {
		csm := new(ClinicServiceMock)

		// Mock router
		router := gin.Default()
		NewClinicHandler(router.Group("/api"), csm)

		// Mock Body
		jsonBody := []byte(`{
			"name":         "診所",
			"start_at":       "2022-13-31",
			"end_at":         "2022-13-31",
			"quota_per_month": 100
		}`)
		body := bytes.NewReader(jsonBody)

		// Request Body
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/clinic", body)
		req.Header.Set("Authorization", "Bearer "+jwtToken)
		router.ServeHTTP(w, req)

		assert.Equal(test, 400, w.Code)
		assert.Equal(test, string([]byte(`{"status":false,"msg":"Date格式錯誤","data":null}`)), w.Body.String())
		csm.AssertExpectations(test)
	})

	test.Run("失敗：Create()，新增診所，提交空表單。", func(test *testing.T) {
		csm := new(ClinicServiceMock)

		// Mock router
		router := gin.Default()
		NewClinicHandler(router.Group("/api"), csm)

		// Mock Body
		jsonBody := []byte(`{
			"name":         "",
			"start_at":       "",
			"end_at":         "",
			"quota_per_month": 0
		}`)
		body := bytes.NewReader(jsonBody)

		// Request Body
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/clinic", body)
		req.Header.Set("Authorization", "Bearer "+jwtToken)
		router.ServeHTTP(w, req)

		assert.Equal(test, 400, w.Code)
		assert.Equal(test, string([]byte(`{"status":false,"msg":"Key: 'ClinicCreateForm.Name' Error:Field validation for 'Name' failed on the 'required' tag\nKey: 'ClinicCreateForm.StartAt' Error:Field validation for 'StartAt' failed on the 'required' tag\nKey: 'ClinicCreateForm.EndAt' Error:Field validation for 'EndAt' failed on the 'required' tag\nKey: 'ClinicCreateForm.QuotaPerMonth' Error:Field validation for 'QuotaPerMonth' failed on the 'required' tag","data":null}`)), w.Body.String())
		csm.AssertExpectations(test)
	})

	test.Run("失敗：Create()，新增診所，但DB有問題。", func(test *testing.T) {
		csm := new(ClinicServiceMock)

		// 轉成返回格式
		var jsonObject = &model.ClinicCreateForm{
			Name:          "診所",
			StartAt:       "2022-10-31",
			EndAt:         "9999-12-31",
			QuotaPerMonth: 100,
		}

		jsonString := "DB沒有連線"
		expected, _ := json.Marshal(utils.H500(jsonString))

		// Mock router
		csm.On("Create", jsonObject).Return(errors.New(jsonString))
		router := gin.Default()
		NewClinicHandler(router.Group("/api"), csm)

		// Mock Body
		jsonBody := []byte(`{
			"name":         "診所",
			"start_at":       "2022-10-31",
			"end_at":         "9999-12-31",
			"quota_per_month": 100
		}`)
		body := bytes.NewReader(jsonBody)

		// Request Body
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/clinic", body)
		req.Header.Set("Authorization", "Bearer "+jwtToken)
		router.ServeHTTP(w, req)

		assert.Equal(test, 400, w.Code)
		assert.Equal(test, string(expected), w.Body.String())
		csm.AssertExpectations(test)
	})

	// Update
	test.Run("成功：Update()，更新診所。", func(test *testing.T) {
		csm := new(ClinicServiceMock)

		// 轉成返回格式
		var jsonObject = &model.ClinicUpdateForm{
			ClinicTokenForm: model.ClinicTokenForm{
				ClinicId: 1,
			},
			ClinicCreateForm: model.ClinicCreateForm{
				Name:          "診所",
				StartAt:       "2022-10-31",
				EndAt:         "9999-12-31",
				QuotaPerMonth: 100,
			},
		}

		// Mock router
		csm.On("Update", jsonObject).Return(nil)
		router := gin.Default()
		NewClinicHandler(router.Group("/api"), csm)

		// Mock Body
		jsonBody := []byte(`{
			"clinic_id":  1,
			"name":         "診所",
			"start_at":       "2022-10-31",
			"end_at":         "9999-12-31",
			"quota_per_month": 100
		}`)
		body := bytes.NewReader(jsonBody)

		// Request Body
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("PUT", "/api/clinic", body)
		req.Header.Set("Authorization", "Bearer "+jwtToken)
		router.ServeHTTP(w, req)

		assert.Equal(test, 200, w.Code)
		assert.Equal(test, string([]byte(`{"status":true,"msg":"更新成功","data":null}`)), w.Body.String())
		csm.AssertExpectations(test)
	})

	test.Run("失敗：Update()，更新診所，Body格式錯誤。", func(test *testing.T) {
		csm := new(ClinicServiceMock)

		// Mock router
		router := gin.Default()
		NewClinicHandler(router.Group("/api"), csm)

		// Mock Body
		jsonBody := []byte(`{
			"clinic_id":  1,
			"name":         "診所",
			"start_at":       "2022-10-31",
			"end_at":         "9999-12-31",
			"quota_per_month": 100,
		}`)
		body := bytes.NewReader(jsonBody)

		// Request Body
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("PUT", "/api/clinic", body)
		req.Header.Set("Authorization", "Bearer "+jwtToken)
		router.ServeHTTP(w, req)

		assert.Equal(test, 400, w.Code)
		assert.Equal(test, string([]byte(`{"status":false,"msg":"invalid character '}' looking for beginning of object key string","data":null}`)), w.Body.String())
	})

	test.Run("失敗：Update()，更新診所，Date格式錯誤。", func(test *testing.T) {
		csm := new(ClinicServiceMock)

		// Mock router
		router := gin.Default()
		NewClinicHandler(router.Group("/api"), csm)

		// Mock Body
		jsonBody := []byte(`{
			"clinic_id":  1,
			"name":         "診所",
			"start_at":       "2022-13-31",
			"end_at":         "2022-13-31",
			"quota_per_month": 100
		}`)
		body := bytes.NewReader(jsonBody)

		// Request Body
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("PUT", "/api/clinic", body)
		req.Header.Set("Authorization", "Bearer "+jwtToken)
		router.ServeHTTP(w, req)

		assert.Equal(test, 400, w.Code)
		assert.Equal(test, string([]byte(`{"status":false,"msg":"Date格式錯誤","data":null}`)), w.Body.String())
		csm.AssertExpectations(test)
	})

	test.Run("失敗：Update()，更新診所，表單為空。", func(test *testing.T) {
		csm := new(ClinicServiceMock)

		// Mock router
		router := gin.Default()
		NewClinicHandler(router.Group("/api"), csm)

		// Mock Body
		jsonBody := []byte(`{
			"clinic_id":  0,
			"name":         "",
			"start_at":       "",
			"end_at":         "",
			"quota_per_month": 0
		}`)
		body := bytes.NewReader(jsonBody)

		// Request Body
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("PUT", "/api/clinic", body)
		req.Header.Set("Authorization", "Bearer "+jwtToken)
		router.ServeHTTP(w, req)

		assert.Equal(test, 400, w.Code)
		assert.Equal(test, string([]byte(`{"status":false,"msg":"Key: 'ClinicUpdateForm.ClinicTokenForm.ClinicId' Error:Field validation for 'ClinicId' failed on the 'required' tag\nKey: 'ClinicUpdateForm.ClinicCreateForm.Name' Error:Field validation for 'Name' failed on the 'required' tag\nKey: 'ClinicUpdateForm.ClinicCreateForm.StartAt' Error:Field validation for 'StartAt' failed on the 'required' tag\nKey: 'ClinicUpdateForm.ClinicCreateForm.EndAt' Error:Field validation for 'EndAt' failed on the 'required' tag\nKey: 'ClinicUpdateForm.ClinicCreateForm.QuotaPerMonth' Error:Field validation for 'QuotaPerMonth' failed on the 'required' tag","data":null}`)), w.Body.String())
	})

	test.Run("失敗：Update()，更新診所，DB有問題。", func(test *testing.T) {
		csm := new(ClinicServiceMock)

		// 轉成返回格式
		var jsonObject = &model.ClinicUpdateForm{
			ClinicTokenForm: model.ClinicTokenForm{
				ClinicId: 1,
			},
			ClinicCreateForm: model.ClinicCreateForm{
				Name:          "診所",
				StartAt:       "2022-10-31",
				EndAt:         "9999-12-31",
				QuotaPerMonth: 100,
			},
		}
		var errorString string = "有錯誤發生"
		expected, _ := json.Marshal(utils.H500(errorString))

		// Mock router
		csm.On("Update", jsonObject).Return(errors.New(errorString))
		router := gin.Default()
		NewClinicHandler(router.Group("/api"), csm)

		// Mock Body
		jsonBody := []byte(`{
			"clinic_id":  1,
			"name":         "診所",
			"start_at":       "2022-10-31",
			"end_at":         "9999-12-31",
			"quota_per_month": 100
		}`)
		body := bytes.NewReader(jsonBody)

		// Request Body
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("PUT", "/api/clinic", body)
		req.Header.Set("Authorization", "Bearer "+jwtToken)
		router.ServeHTTP(w, req)

		assert.Equal(test, 400, w.Code)
		assert.Equal(test, string(expected), w.Body.String())
		csm.AssertExpectations(test)
	})

	// UpdateToken
	test.Run("成功：UpdateToken()，更新Token。", func(test *testing.T) {
		csm := new(ClinicServiceMock)

		// 轉成返回格式
		var jsonObject = &model.ClinicTokenForm{
			ClinicId: 1,
		}

		// Mock router
		csm.On("UpdateToken", jsonObject).Return(nil)
		router := gin.Default()
		NewClinicHandler(router.Group("/api"), csm)

		// Mock Body
		jsonBody := []byte(`{"clinic_id":  1}`)
		body := bytes.NewReader(jsonBody)

		// Request Body
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("PATCH", "/api/clinic/token", body)
		req.Header.Set("Authorization", "Bearer "+jwtToken)
		router.ServeHTTP(w, req)

		assert.Equal(test, 200, w.Code)
		assert.Equal(test, string([]byte(`{"status":true,"msg":"更新成功","data":null}`)), w.Body.String())
		csm.AssertExpectations(test)
	})

	test.Run("失敗：UpdateToken()，更新Token，Body格式錯誤。", func(test *testing.T) {
		csm := new(ClinicServiceMock)

		// Mock router
		router := gin.Default()
		NewClinicHandler(router.Group("/api"), csm)

		// Mock Body
		jsonBody := []byte(`{"clinic_id":  1,}`)
		body := bytes.NewReader(jsonBody)

		// Request Body
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("PATCH", "/api/clinic/token", body)
		req.Header.Set("Authorization", "Bearer "+jwtToken)
		router.ServeHTTP(w, req)

		assert.Equal(test, 400, w.Code)
		assert.Equal(test, string([]byte(`{"status":false,"msg":"invalid character '}' looking for beginning of object key string","data":null}`)), w.Body.String())
	})

	test.Run("失敗：UpdateToken()，更新Token，表單為空。", func(test *testing.T) {
		csm := new(ClinicServiceMock)

		router := gin.Default()
		NewClinicHandler(router.Group("/api"), csm)

		// Mock Body
		jsonBody := []byte(`{"clinic_id":  0}`)
		body := bytes.NewReader(jsonBody)

		// Request Body
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("PATCH", "/api/clinic/token", body)
		req.Header.Set("Authorization", "Bearer "+jwtToken)
		router.ServeHTTP(w, req)

		assert.Equal(test, 400, w.Code)
		assert.Equal(test, string([]byte(`{"status":false,"msg":"Key: 'ClinicTokenForm.ClinicId' Error:Field validation for 'ClinicId' failed on the 'required' tag","data":null}`)), w.Body.String())
	})

	test.Run("失敗：UpdateToken()，更新Token，DB有問題。", func(test *testing.T) {
		csm := new(ClinicServiceMock)

		// 轉成返回格式
		var jsonObject = &model.ClinicTokenForm{
			ClinicId: 1,
		}
		var errorString string = "有錯誤發生"
		expected, _ := json.Marshal(utils.H500(errorString))

		// Mock router
		csm.On("UpdateToken", jsonObject).Return(errors.New(errorString))
		router := gin.Default()
		NewClinicHandler(router.Group("/api"), csm)

		// Mock Body
		jsonBody := []byte(`{"clinic_id":  1}`)
		body := bytes.NewReader(jsonBody)

		// Request Body
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("PATCH", "/api/clinic/token", body)
		req.Header.Set("Authorization", "Bearer "+jwtToken)
		router.ServeHTTP(w, req)

		assert.Equal(test, 400, w.Code)
		assert.Equal(test, string(expected), w.Body.String())
		csm.AssertExpectations(test)
	})
}
