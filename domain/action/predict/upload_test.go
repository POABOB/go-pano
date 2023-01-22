package action

import (
	"bytes"
	"encoding/json"
	"errors"
	"mime/multipart"
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
// mock service 的 grpcReq httpReq processFile
type CreateServiceMock struct {
	mock.Mock
}

func (sm *CreateServiceMock) Create(predictForm *model.PredictForm, dir string, fileName string) (interface{}, error) {
	args := sm.Called(predictForm, dir, fileName)
	return args.Get(0), args.Error(1)
}

func (sm *CreateServiceMock) GrpcReq(string_ string, url string) (*model.Result, error) {
	args := sm.Called(string_, url)
	return args.Get(0).(*model.Result), args.Error(1)
}

func (sm *CreateServiceMock) HttpReq(string_ string, url string) (*model.Result, error) {
	args := sm.Called(string_, url)
	return args.Get(0).(*model.Result), args.Error(1)
}

type UploadServiceMock struct {
	mock.Mock
}

func (sm *UploadServiceMock) Upload(predictForm *model.PredictForm, ctx *gin.Context) (string, string, error) {
	args := sm.Called(predictForm, ctx)
	return args.Get(0).(string), args.Get(1).(string), args.Error(2)
}

func TestUploadAction(test *testing.T) {

	test.Run("成功：在DB抓到資料，不用再辨識。", func(test *testing.T) {
		csm := new(CreateServiceMock)
		usm := new(UploadServiceMock)

		// 轉成返回格式
		jsonString := `{"test": "predict"}`
		var jsonObject interface{}
		json.Unmarshal([]byte(jsonString), &jsonObject)
		expected, _ := json.Marshal(utils.H200(jsonObject, ""))

		// Mock router
		usm.On("Upload",
			&model.PredictForm{
				ClinicId: 1,
				Method:   "http",
			},
			mock.AnythingOfType("*gin.Context"),
		).Return("dirPath", jsonString, nil)

		router := gin.Default()
		action := NewUploadAction(usm, csm)
		router.POST("/api/predict", action.Upload)

		// Mock Body
		body := new(bytes.Buffer)
		mw := multipart.NewWriter(body)
		// 檔案的欄位
		writer, _ := mw.CreateFormFile("nhicode", "00008026.jpg")
		_, _ = writer.Write([]byte("Test"))
		// clinic_id的欄位
		mw.WriteField("clinic_id", "1")
		mw.Close()

		// Request
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/predict", body)
		req.Header.Set("Content-Type", mw.FormDataContentType())
		router.ServeHTTP(w, req)

		assert.Equal(test, 200, w.Code)
		assert.Equal(test, string(expected), w.Body.String())
		usm.AssertExpectations(test)
	})

	test.Run("成功：在DB抓不到資料，請求Python辨識。", func(test *testing.T) {
		csm := new(CreateServiceMock)
		usm := new(UploadServiceMock)

		// 轉成返回格式
		jsonString := `{"test": "predict"}`
		var jsonObject interface{}
		json.Unmarshal([]byte(jsonString), &jsonObject)
		expected, _ := json.Marshal(utils.H200(jsonObject, ""))

		var predictForm *model.PredictForm = &model.PredictForm{
			ClinicId: 2,
			Method:   "http",
		}
		// Mock router
		usm.On("Upload", predictForm, mock.AnythingOfType("*gin.Context")).Return("dirPath", "", nil)
		csm.On("Create", predictForm, "dirPath", "").Return(jsonObject, nil)

		router := gin.Default()
		action := NewUploadAction(usm, csm)
		router.POST("/api/predict", action.Upload)

		// Mock Body
		body := new(bytes.Buffer)
		mw := multipart.NewWriter(body)
		// 檔案的欄位
		writer, _ := mw.CreateFormFile("nhicode", "00008026.jpg")
		_, _ = writer.Write([]byte("Test"))
		// clinic_id的欄位
		mw.WriteField("clinic_id", "2")
		mw.Close()

		// Request
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/predict", body)
		req.Header.Set("Content-Type", mw.FormDataContentType())
		router.ServeHTTP(w, req)

		assert.Equal(test, 200, w.Code)
		assert.Equal(test, string(expected), w.Body.String())
		usm.AssertExpectations(test)
		csm.AssertExpectations(test)
	})

	test.Run("失敗：表單驗證沒有通過(expected: int, actually: string)。", func(test *testing.T) {
		csm := new(CreateServiceMock)
		usm := new(UploadServiceMock)

		// 轉成返回格式
		expected, _ := json.Marshal(utils.H500(`strconv.ParseInt: parsing "a": invalid syntax`))

		router := gin.Default()
		action := NewUploadAction(usm, csm)
		router.POST("/api/predict", action.Upload)

		// Mock Body
		body := new(bytes.Buffer)
		mw := multipart.NewWriter(body)
		// clinic_id的欄位
		mw.WriteField("clinic_id", "a")
		mw.Close()

		// Request
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/predict", body)
		req.Header.Set("Content-Type", mw.FormDataContentType())
		router.ServeHTTP(w, req)

		assert.Equal(test, 400, w.Code)
		assert.Equal(test, string(expected), w.Body.String())
	})

	test.Run("失敗：表單驗證沒有通過(required)。", func(test *testing.T) {
		csm := new(CreateServiceMock)
		usm := new(UploadServiceMock)

		// 轉成返回格式
		expected, _ := json.Marshal(utils.H500(`Key: 'PredictForm.ClinicId' Error:Field validation for 'ClinicId' failed on the 'required' tag`))

		router := gin.Default()
		action := NewUploadAction(usm, csm)
		router.POST("/api/predict", action.Upload)

		// Mock Body
		body := new(bytes.Buffer)
		mw := multipart.NewWriter(body)
		mw.Close()

		// Request
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/predict", body)
		req.Header.Set("Content-Type", mw.FormDataContentType())
		router.ServeHTTP(w, req)

		assert.Equal(test, 400, w.Code)
		assert.Equal(test, string(expected), w.Body.String())
	})

	test.Run("失敗：表單驗證沒有通過(max)。", func(test *testing.T) {
		csm := new(CreateServiceMock)
		usm := new(UploadServiceMock)

		// 轉成返回格式
		expected, _ := json.Marshal(utils.H500(`Key: 'PredictForm.ClinicId' Error:Field validation for 'ClinicId' failed on the 'max' tag`))

		router := gin.Default()
		action := NewUploadAction(usm, csm)
		router.POST("/api/predict", action.Upload)

		// Mock Body
		body := new(bytes.Buffer)
		mw := multipart.NewWriter(body)
		mw.WriteField("clinic_id", "100000000000")
		mw.Close()

		// Request
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/predict", body)
		req.Header.Set("Content-Type", mw.FormDataContentType())
		router.ServeHTTP(w, req)

		assert.Equal(test, 400, w.Code)
		assert.Equal(test, string(expected), w.Body.String())
	})

	test.Run("失敗：沒有傳照片，Upload()返回錯誤。", func(test *testing.T) {
		csm := new(CreateServiceMock)
		usm := new(UploadServiceMock)

		// 轉成返回格式
		expected, _ := json.Marshal(utils.H500(`沒有照片傳入！`))

		usm.On("Upload",
			&model.PredictForm{
				ClinicId: 1,
				Method:   "http",
			},
			mock.AnythingOfType("*gin.Context"),
		).Return("", "", errors.New("沒有照片傳入！"))

		router := gin.Default()
		action := NewUploadAction(usm, csm)
		router.POST("/api/predict", action.Upload)

		// Mock Body
		body := new(bytes.Buffer)
		mw := multipart.NewWriter(body)
		mw.WriteField("clinic_id", "1")
		mw.Close()

		// Request
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/predict", body)
		req.Header.Set("Content-Type", mw.FormDataContentType())
		router.ServeHTTP(w, req)

		assert.Equal(test, 400, w.Code)
		assert.Equal(test, string(expected), w.Body.String())
		usm.AssertExpectations(test)
	})

	test.Run("失敗：參數傳錯，Create()返回錯誤。", func(test *testing.T) {
		csm := new(CreateServiceMock)
		usm := new(UploadServiceMock)

		// 轉成返回格式
		expected, _ := json.Marshal(utils.H500(`參數只接受：http和grpc！`))

		var predictForm *model.PredictForm = &model.PredictForm{
			ClinicId: 2,
			Method:   "h666",
		}
		// Mock router
		usm.On("Upload", predictForm, mock.AnythingOfType("*gin.Context")).Return("dirPath", "", nil)
		csm.On("Create", predictForm, "dirPath", "").Return(nil, errors.New(`參數只接受：http和grpc！`))

		router := gin.Default()
		action := NewUploadAction(usm, csm)
		router.POST("/api/predict/:method", action.Upload)

		// Mock Body
		body := new(bytes.Buffer)
		mw := multipart.NewWriter(body)
		writer, _ := mw.CreateFormFile("nhicode", "00008026.jpg")
		_, _ = writer.Write([]byte("Test"))
		mw.WriteField("clinic_id", "2")
		mw.Close()

		// Request
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/predict/h666", body)
		req.Header.Set("Content-Type", mw.FormDataContentType())
		router.ServeHTTP(w, req)

		assert.Equal(test, 400, w.Code)
		assert.Equal(test, string(expected), w.Body.String())
		usm.AssertExpectations(test)
		csm.AssertExpectations(test)
	})
}
