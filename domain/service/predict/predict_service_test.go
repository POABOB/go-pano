package predict_service

// import (
// 	"bytes"
// 	"context"
// 	"errors"
// 	"log"
// 	"mime/multipart"
// 	"net"
// 	"net/http"
// 	"net/http/httptest"
// 	"testing"
// 	"time"

// 	"go-pano/config"
// 	"go-pano/domain/model"
// 	pb "go-pano/protos/predict"
// 	"go-pano/utils"

// 	"github.com/gin-gonic/gin"
// 	"github.com/stretchr/testify/assert"
// 	"github.com/stretchr/testify/mock"
// 	"google.golang.org/grpc"
// )

// func init() {
// 	gin.SetMode(gin.TestMode)
// 	config.LoadConfigTest()
// 	utils.Reset()
// }

// func (prm *PredictRepositoryMock) GetFirstByIDAndFileName(p *model.Predict) (*model.Predict, error) {
// 	args := prm.Called(p)
// 	return args.Get(0).(*model.Predict), args.Error(1)
// }

// func (prm *PredictRepositoryMock) Create(p *model.Predict) (*model.Predict, error) {
// 	args := prm.Called(p)
// 	return args.Error(0)
// }

// // Mock self method
// func (ps *PredictService) GrpcReq(string_ string, url string) (*model.Result, error) {
// 	args := sm.Called(string_, url)
// 	return args.Get(0).(*model.Result), args.Error(1)
// }

// type PredictServer struct {
// 	mock.Mock
// }

// func (s *PredictServer) Predict(ctx context.Context, req *pb.Request) (*pb.Response, error) {
// 	args := s.Called(ctx, req)
// 	return args.Get(0).(*pb.Response), args.Error(1)
// }

// // 以下區塊是使用sm去注入UploadService來假裝函數
// type FileUtilsMock struct {
// 	mock.Mock
// }

// // mock service 的 PathExist SaveWithGo
// func (fum *FileUtilsMock) PathExist(string_ string) error {
// 	args := fum.Called(string_)
// 	return args.Error(0)
// }

// func (fum *FileUtilsMock) SaveWithGo(ctx *gin.Context, base_dir string, files []*multipart.FileHeader) error {
// 	args := fum.Called(ctx, base_dir, files)
// 	return args.Error(0)
// }

// // 以下是測試環節
// // 正常測試func
// func TestUploadService(test *testing.T) {
// 	// 建立一個假的gin context
// 	ctx, _ := gin.CreateTestContext(httptest.NewRecorder())
// 	// Mock files
// 	buf := new(bytes.Buffer)
// 	// func NewWriter(w io.Writer) *Writer
// 	// Official: NewWriter returns a new multipart Writer with a random boundary, writing to w.
// 	// 意思是這個函數會把我們要寫入的buffer寫進去
// 	mw := multipart.NewWriter(buf)
// 	// func (w *Writer) CreateFormFile(fieldname, filename string) (io.Writer, error)
// 	// Official: CreateFormFile is a convenience wrapper around CreatePart.
// 	// 			It creates a new form-data header with the provided field name and file name.
// 	// 我們會使用mw將File欄位寫入，一個是欄位名稱，另為一個是檔案名稱
// 	w, _ := mw.CreateFormFile("nhicode", "00008026.jpg")
// 	// 將"Test"字串寫入，正常來說這裡是檔案
// 	_, _ = w.Write([]byte("Test"))
// 	mw.Close()
// 	// 請求攜帶上傳的檔案
// 	ctx.Request, _ = http.NewRequest("POST", "/api/predict/http", buf)
// 	// 設定Header為 FormData
// 	ctx.Request.Header.Set("Content-Type", mw.FormDataContentType())

// 	test.Run("成功：PredictForm和File都有正確攜帶，DB有找到clinicId和fileName。", func(test *testing.T) {
// 		rm := new(PredictRepositoryMock)
// 		um := new(FileUtilsMock)

// 		// Mock funcs
// 		// time.Time型別 Timestamp
// 		t := time.Now()
// 		rm.On("GetFirstByIDAndFileName", 1, "00008026.jpg").
// 			Return(&model.Predict{
// 				ID:        1,
// 				ClinicId:  1,
// 				Filename:  "00008026.jpg",
// 				Predict:   "{\"00008026.jpg\": \"test\"}",
// 				CreatedAt: t,
// 				UpdatedAt: t,
// 			}, nil)

// 		// 將Ｍock注入真的Service
// 		uploadServiceWithMockData := NewUploadService(rm, um)
// 		_, _, err := uploadServiceWithMockData.Upload(
// 			&model.PredictForm{
// 				ClinicId: 1,
// 				Method:   "http",
// 			},
// 			ctx,
// 		)
// 		assert.NoError(test, err)
// 		// 驗證.On()是否真的有被 call 到
// 		rm.AssertExpectations(test)
// 	})

// 	test.Run("成功：PredictForm和File都有正確攜帶，DB沒有找到資料。", func(test *testing.T) {
// 		rm := new(PredictRepositoryMock)
// 		um := new(FileUtilsMock)

// 		// Mock funcs
// 		rm.On("GetFirstByIDAndFileName", 1, "00008026.jpg").
// 			Return(&model.Predict{}, errors.New("找不到該筆資料！"))

// 		um.On("PathExist", mock.AnythingOfType("string")).Return(nil)
// 		um.On("SaveWithGo",
// 			mock.AnythingOfType("*gin.Context"),
// 			mock.AnythingOfType("string"),
// 			mock.AnythingOfType("[]*multipart.FileHeader"),
// 		).Return(nil)

// 		// 將Ｍock注入真的Service
// 		uploadServiceWithMockData := NewUploadService(rm, um)
// 		_, _, err := uploadServiceWithMockData.Upload(
// 			&model.PredictForm{
// 				ClinicId: 1,
// 				Method:   "http",
// 			},
// 			ctx,
// 		)
// 		assert.NoError(test, err)

// 		rm.AssertExpectations(test)
// 		um.AssertExpectations(test)
// 	})

// 	test.Run("失敗：PredictForm有正確攜帶，但是File沒有FormData傳入。", func(test *testing.T) {
// 		rm := new(PredictRepositoryMock)
// 		um := new(FileUtilsMock)

// 		// 將Ｍock注入真的Service
// 		ctxWithoutFormData, _ := gin.CreateTestContext(httptest.NewRecorder())
// 		ctxWithoutFormData.Request, _ = http.NewRequest("POST", "/api/predict/http", nil)
// 		uploadServiceWithMockData := NewUploadService(rm, um)
// 		_, _, err := uploadServiceWithMockData.Upload(
// 			&model.PredictForm{
// 				ClinicId: 1,
// 				Method:   "http",
// 			},
// 			ctxWithoutFormData,
// 		)

// 		assert.EqualError(test, err, "request Content-Type isn't multipart/form-data")
// 	})

// 	test.Run("失敗：PredictForm有正確攜帶，但是File的欄位不是'nhicode'。", func(test *testing.T) {
// 		rm := new(PredictRepositoryMock)
// 		um := new(FileUtilsMock)

// 		// 將Ｍock注入真的Service
// 		ctxWithWrongField, _ := gin.CreateTestContext(httptest.NewRecorder())
// 		buf := new(bytes.Buffer)
// 		mw := multipart.NewWriter(buf)
// 		w, _ := mw.CreateFormFile("nhicode2", "00008026.jpg")
// 		_, _ = w.Write([]byte("Test"))
// 		mw.Close()
// 		ctxWithWrongField.Request, _ = http.NewRequest("POST", "/api/predict/http", buf)
// 		ctxWithWrongField.Request.Header.Set("Content-Type", mw.FormDataContentType())
// 		uploadServiceWithMockData := NewUploadService(rm, um)
// 		_, _, err := uploadServiceWithMockData.Upload(
// 			&model.PredictForm{
// 				ClinicId: 1,
// 				Method:   "http",
// 			},
// 			ctxWithWrongField,
// 		)

// 		assert.EqualError(test, err, "請使用'nhicode'作為上傳名稱！")
// 	})

// 	jsonString := `{"00008026.jpg":"test"}`
// 	test.Run("成功：使用HTTP請求，並且沒有返回錯誤。", func(test *testing.T) {
// 		rm := new(PredictRepositoryMock)
// 		sm := new(CreateServiceMock)

// 		rm.On("Create", 1, "00008026.jpg", jsonString).Return(nil)
// 		sm.On("HttpReq", "dirPath", httpPythonUrl).Return(&model.Result{
// 			IsSuccessful: true,
// 			Msg:          "success",
// 			Predict:      jsonString,
// 		}, nil)

// 		// 將Ｍock注入真的Service
// 		createServiceWithMockData := NewCreateService(rm, sm)
// 		_, err := createServiceWithMockData.Create(
// 			&model.PredictForm{
// 				ClinicId: 1,
// 				Method:   "http",
// 			},
// 			"dirPath",
// 			"00008026.jpg",
// 		)

// 		assert.NoError(test, err)
// 		rm.AssertExpectations(test)
// 		sm.AssertExpectations(test)
// 	})

// 	test.Run("成功：使用gRPC請求，並且沒有返回錯誤。", func(test *testing.T) {
// 		rm := new(PredictRepositoryMock)
// 		sm := new(CreateServiceMock)

// 		rm.On("Create", 1, "00008026.jpg", jsonString).Return(nil)
// 		sm.On("GrpcReq", "dirPath", grpcPythonUrl).Return(&model.Result{
// 			IsSuccessful: true,
// 			Msg:          "success",
// 			Predict:      jsonString,
// 		}, nil)

// 		// 將Ｍock注入真的Service
// 		createServiceWithMockData := NewCreateService(rm, sm)
// 		_, err := createServiceWithMockData.Create(
// 			&model.PredictForm{
// 				ClinicId: 1,
// 				Method:   "grpc",
// 			},
// 			"dirPath",
// 			"00008026.jpg",
// 		)

// 		assert.NoError(test, err)
// 		rm.AssertExpectations(test)
// 		sm.AssertExpectations(test)
// 	})

// 	test.Run("失敗：使用請求的參數不支援，並且沒有返回錯誤。", func(test *testing.T) {
// 		rm := new(PredictRepositoryMock)
// 		sm := new(CreateServiceMock)

// 		createServiceWithMockData := NewCreateService(rm, sm)
// 		_, err := createServiceWithMockData.Create(
// 			&model.PredictForm{
// 				ClinicId: 1,
// 				Method:   "smtp",
// 			},
// 			"dirPath",
// 			"00008026.jpg",
// 		)

// 		assert.EqualError(test, err, "參數只接受：http和grpc")
// 	})

// 	test.Run("失敗：使用其中一種請求，但是Python的Service報錯。", func(test *testing.T) {
// 		rm := new(PredictRepositoryMock)
// 		sm := new(CreateServiceMock)

// 		// 因為報錯，所以返回的Result是空物件
// 		sm.On("HttpReq", "dirPath", httpPythonUrl).Return(&model.Result{}, errors.New("Error came from python: "))

// 		createServiceWithMockData := NewCreateService(rm, sm)
// 		_, err := createServiceWithMockData.Create(
// 			&model.PredictForm{
// 				ClinicId: 1,
// 				Method:   "http",
// 			},
// 			"dirPath",
// 			"00008026.jpg",
// 		)

// 		assert.EqualError(test, err, "Error came from python: ")
// 		sm.AssertExpectations(test)
// 	})

// 	test.Run("失敗：使用其中一種請求，但是Python返回的字串格式有誤。", func(test *testing.T) {
// 		rm := new(PredictRepositoryMock)
// 		sm := new(CreateServiceMock)

// 		// 因為報錯，所以返回的Result是空物件
// 		sm.On("HttpReq", "dirPath", httpPythonUrl).Return(&model.Result{
// 			IsSuccessful: true,
// 			Msg:          "success",
// 			Predict:      `{"123":123,,,,,}`,
// 		}, nil)

// 		createServiceWithMockData := NewCreateService(rm, sm)
// 		_, err := createServiceWithMockData.Create(
// 			&model.PredictForm{
// 				ClinicId: 1,
// 				Method:   "http",
// 			},
// 			"dirPath",
// 			"00008026.jpg",
// 		)
// 		assert.EqualError(test, err, "invalid character ',' looking for beginning of object key string")
// 		sm.AssertExpectations(test)
// 	})

// 	test.Run("失敗：測試GrpcReq，但是Server connection refused。", func(test *testing.T) {
// 		rm := new(PredictRepositoryMock)
// 		sm := new(CreateServiceMock)
// 		createServiceWithMockData := NewCreateService(rm, sm)
// 		result, err := createServiceWithMockData.GrpcReq("dirPath", "127.0.0.1:5001")
// 		assert.EqualError(test, err, `Error came from python: rpc error: code = Unavailable desc = connection error: desc = "transport: Error while dialing dial tcp 127.0.0.1:5001: connect: connection refused"`)
// 		assert.Equal(test, result, &model.Result{
// 			IsSuccessful: false,
// 			Msg:          "",
// 			Predict:      "",
// 		})
// 	})

// 	test.Run("成功：測試HttpReq，接收Python正確返回資訊。", func(test *testing.T) {
// 		rm := new(PredictRepositoryMock)
// 		sm := new(CreateServiceMock)
// 		ps := new(PredictServer)

// 		// STEP 2-1：定義要監聽的 port 號
// 		lis, err := net.Listen("tcp", grpcPythonUrl)
// 		if err != nil {
// 			log.Fatalf("failed to listed: %v", err)
// 		}
// 		// STEP 2-2：使用 gRPC 的 NewServer 方法來建立 gRPC Server 的實例
// 		grpcServer := grpc.NewServer()

// 		// STEP 2-3：在 gRPC Server 中註冊 service 的實作
// 		// 使用 proto 提供的 RegisterRouteGuideServer 方法，並將 routeGuideServer 作為參數傳入
// 		pb.RegisterPredictServer(grpcServer, ps)

// 		// STEP 2-4：啟動 grpcServer，並阻塞在這裡直到該程序被 kill 或 stop
// 		go func() {
// 			err = grpcServer.Serve(lis)
// 			if err != nil {
// 				panic(err)
// 			}
// 		}()
// 		defer grpcServer.Stop()

// 		ps.On("Predict",
// 			mock.AnythingOfType("*context.valueCtx"),
// 			&pb.Request{Dir: "dirPath"},
// 		).Return(&pb.Response{
// 			IsSuccessful: true,
// 			Msg:          "成功",
// 			Predict:      `predict`,
// 		}, nil)

// 		createServiceWithMockData := NewCreateService(rm, sm)
// 		result, err := createServiceWithMockData.GrpcReq("dirPath", grpcPythonUrl)
// 		assert.NoError(test, err)
// 		assert.Equal(test, result, &model.Result{
// 			IsSuccessful: true,
// 			Msg:          "成功",
// 			Predict:      "predict",
// 		})
// 	})

// 	test.Run("成功：測試HttpReq，接收Python正確返回資訊。", func(test *testing.T) {
// 		rm := new(PredictRepositoryMock)
// 		sm := new(CreateServiceMock)
// 		ps := new(PredictServer)

// 		// STEP 2-1：定義要監聽的 port 號
// 		lis, err := net.Listen("tcp", grpcPythonUrl)
// 		if err != nil {
// 			log.Fatalf("failed to listed: %v", err)
// 		}
// 		// STEP 2-2：使用 gRPC 的 NewServer 方法來建立 gRPC Server 的實例
// 		grpcServer := grpc.NewServer()

// 		// STEP 2-3：在 gRPC Server 中註冊 service 的實作
// 		// 使用 proto 提供的 RegisterRouteGuideServer 方法，並將 routeGuideServer 作為參數傳入
// 		pb.RegisterPredictServer(grpcServer, ps)

// 		// STEP 2-4：啟動 grpcServer，並阻塞在這裡直到該程序被 kill 或 stop
// 		go func() {
// 			err = grpcServer.Serve(lis)
// 			if err != nil {
// 				panic(err)
// 			}
// 		}()
// 		defer grpcServer.Stop()
// 		ps.On("Predict",
// 			mock.AnythingOfType("*context.valueCtx"),
// 			&pb.Request{Dir: "dirPath"},
// 		).Return(&pb.Response{
// 			IsSuccessful: false,
// 			Msg:          "有問題",
// 			Predict:      "",
// 		}, nil)

// 		createServiceWithMockData := NewCreateService(rm, sm)
// 		result, err := createServiceWithMockData.GrpcReq("dirPath", grpcPythonUrl)
// 		assert.EqualError(test, err, "Error came from python: 有問題")
// 		assert.Equal(test, result, &model.Result{
// 			IsSuccessful: false,
// 			Msg:          "",
// 			Predict:      "",
// 		})
// 	})
// }
