package service

import (
	"bytes"
	"errors"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

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

// PredictRepositoryMock的聲明在create_test.go
// mock service 的 GetFirstByIDAndFileName
func (rm *PredictRepositoryMock) GetFirstByIDAndFileName(clinicId int, fileName string) (*model.Predict, error) {
	args := rm.Called(clinicId, fileName)
	return args.Get(0).(*model.Predict), args.Error(1)
}

// 以下區塊是使用sm去注入UploadService來假裝函數
type FileUtilsMock struct {
	mock.Mock
}

// mock service 的 PathExist SaveWithGo
func (um *FileUtilsMock) PathExist(string_ string) error {
	args := um.Called(string_)
	return args.Error(0)
}

func (um *FileUtilsMock) SaveWithGo(ctx *gin.Context, base_dir string, files []*multipart.FileHeader) error {
	args := um.Called(ctx, base_dir, files)
	return args.Error(0)
}

// 以下是測試環節
// 正常測試func
func TestUploadService(test *testing.T) {
	// 建立一個假的gin context
	ctx, _ := gin.CreateTestContext(httptest.NewRecorder())
	// Mock files
	buf := new(bytes.Buffer)
	// func NewWriter(w io.Writer) *Writer
	// Official: NewWriter returns a new multipart Writer with a random boundary, writing to w.
	// 意思是這個函數會把我們要寫入的buffer寫進去
	mw := multipart.NewWriter(buf)
	// func (w *Writer) CreateFormFile(fieldname, filename string) (io.Writer, error)
	// Official: CreateFormFile is a convenience wrapper around CreatePart.
	// 			It creates a new form-data header with the provided field name and file name.
	// 我們會使用mw將File欄位寫入，一個是欄位名稱，另為一個是檔案名稱
	w, _ := mw.CreateFormFile("nhicode", "00008026.jpg")
	// 將"Test"字串寫入，正常來說這裡是檔案
	_, _ = w.Write([]byte("Test"))
	mw.Close()
	// 請求攜帶上傳的檔案
	ctx.Request, _ = http.NewRequest("POST", "/api/predict/http", buf)
	// 設定Header為 FormData
	ctx.Request.Header.Set("Content-Type", mw.FormDataContentType())

	test.Run("成功：PredictForm和File都有正確攜帶，DB有找到clinicId和fileName。", func(test *testing.T) {
		rm := new(PredictRepositoryMock)
		um := new(FileUtilsMock)

		// Mock funcs
		// time.Time型別 Timestamp
		t := time.Now()
		rm.On("GetFirstByIDAndFileName", 1, "00008026.jpg").
			Return(&model.Predict{
				ID:        1,
				ClinicId:  1,
				Filename:  "00008026.jpg",
				Predict:   "{\"00008026.jpg\": \"test\"}",
				CreatedAt: t,
				UpdatedAt: t,
			}, nil)

		// 將Ｍock注入真的Service
		uploadServiceWithMockData := NewUploadService(rm, um)
		_, _, err := uploadServiceWithMockData.Upload(
			&model.PredictForm{
				ClinicId: 1,
				Method:   "http",
			},
			ctx,
		)
		assert.NoError(test, err)
		// 驗證.On()是否真的有被 call 到
		rm.AssertExpectations(test)
	})

	test.Run("成功：PredictForm和File都有正確攜帶，DB沒有找到資料。", func(test *testing.T) {
		rm := new(PredictRepositoryMock)
		um := new(FileUtilsMock)

		// Mock funcs
		rm.On("GetFirstByIDAndFileName", 1, "00008026.jpg").
			Return(&model.Predict{}, errors.New("找不到該筆資料！"))

		um.On("PathExist", mock.AnythingOfType("string")).Return(nil)
		um.On("SaveWithGo",
			mock.AnythingOfType("*gin.Context"),
			mock.AnythingOfType("string"),
			mock.AnythingOfType("[]*multipart.FileHeader"),
		).Return(nil)

		// 將Ｍock注入真的Service
		uploadServiceWithMockData := NewUploadService(rm, um)
		_, _, err := uploadServiceWithMockData.Upload(
			&model.PredictForm{
				ClinicId: 1,
				Method:   "http",
			},
			ctx,
		)
		assert.NoError(test, err)

		rm.AssertExpectations(test)
		um.AssertExpectations(test)
	})

	test.Run("失敗：PredictForm有正確攜帶，但是File沒有FormData傳入。", func(test *testing.T) {
		rm := new(PredictRepositoryMock)
		um := new(FileUtilsMock)

		// 將Ｍock注入真的Service
		ctxWithoutFormData, _ := gin.CreateTestContext(httptest.NewRecorder())
		ctxWithoutFormData.Request, _ = http.NewRequest("POST", "/api/predict/http", nil)
		uploadServiceWithMockData := NewUploadService(rm, um)
		_, _, err := uploadServiceWithMockData.Upload(
			&model.PredictForm{
				ClinicId: 1,
				Method:   "http",
			},
			ctxWithoutFormData,
		)

		assert.EqualError(test, err, "request Content-Type isn't multipart/form-data")
	})

	test.Run("失敗：PredictForm有正確攜帶，但是File的欄位不是'nhicode'。", func(test *testing.T) {
		rm := new(PredictRepositoryMock)
		um := new(FileUtilsMock)

		// 將Ｍock注入真的Service
		ctxWithWrongField, _ := gin.CreateTestContext(httptest.NewRecorder())
		buf := new(bytes.Buffer)
		mw := multipart.NewWriter(buf)
		w, _ := mw.CreateFormFile("nhicode2", "00008026.jpg")
		_, _ = w.Write([]byte("Test"))
		mw.Close()
		ctxWithWrongField.Request, _ = http.NewRequest("POST", "/api/predict/http", buf)
		ctxWithWrongField.Request.Header.Set("Content-Type", mw.FormDataContentType())
		uploadServiceWithMockData := NewUploadService(rm, um)
		_, _, err := uploadServiceWithMockData.Upload(
			&model.PredictForm{
				ClinicId: 1,
				Method:   "http",
			},
			ctxWithWrongField,
		)

		assert.EqualError(test, err, "請使用'nhicode'作為上傳名稱！")
	})
}
