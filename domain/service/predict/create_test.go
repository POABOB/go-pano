package service

import (
	"context"
	"errors"
	"log"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"

	"go-pano/config"
	"go-pano/domain/model"
	"go-pano/utils"

	pb "go-pano/protos/predict"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc"
)

var httpPythonUrl string
var grpcPythonUrl string

func init() {
	gin.SetMode(gin.TestMode)
	config.LoadConfigTest()
	utils.Reset()
	httpPythonUrl = "http://" + config.PythonHost + ":5000/"
	grpcPythonUrl = config.PythonHost + ":5001"
}

// 以下區塊是使用rm去注入PredictRepository來假裝函數
// mock repository 的 GetFirstByIDAndFileName Create
type PredictRepositoryMock struct {
	mock.Mock
}

func (rm *PredictRepositoryMock) Create(clinicId int, fileName string, predict string) error {
	args := rm.Called(clinicId, fileName, predict)
	return args.Error(0)
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

type PredictServer struct {
	mock.Mock
}

func (s *PredictServer) Predict(ctx context.Context, req *pb.Request) (*pb.Response, error) {
	args := s.Called(ctx, req)
	return args.Get(0).(*pb.Response), args.Error(1)
}

// 以下是測試環節
// 正常測試func
func TestCreateService(test *testing.T) {
	jsonString := `{"00008026.jpg":"test"}`
	test.Run("成功：使用HTTP請求，並且沒有返回錯誤。", func(test *testing.T) {
		rm := new(PredictRepositoryMock)
		sm := new(CreateServiceMock)

		rm.On("Create", 1, "00008026.jpg", jsonString).Return(nil)
		sm.On("HttpReq", "dirPath", httpPythonUrl).Return(&model.Result{
			IsSuccessful: true,
			Msg:          "success",
			Predict:      jsonString,
		}, nil)

		// 將Ｍock注入真的Service
		createServiceWithMockData := NewCreateService(rm, sm)
		_, err := createServiceWithMockData.Create(
			&model.PredictForm{
				ClinicId: 1,
				Method:   "http",
			},
			"dirPath",
			"00008026.jpg",
		)

		assert.NoError(test, err)
		rm.AssertExpectations(test)
		sm.AssertExpectations(test)
	})

	test.Run("成功：使用gRPC請求，並且沒有返回錯誤。", func(test *testing.T) {
		rm := new(PredictRepositoryMock)
		sm := new(CreateServiceMock)

		rm.On("Create", 1, "00008026.jpg", jsonString).Return(nil)
		sm.On("GrpcReq", "dirPath", grpcPythonUrl).Return(&model.Result{
			IsSuccessful: true,
			Msg:          "success",
			Predict:      jsonString,
		}, nil)

		// 將Ｍock注入真的Service
		createServiceWithMockData := NewCreateService(rm, sm)
		_, err := createServiceWithMockData.Create(
			&model.PredictForm{
				ClinicId: 1,
				Method:   "grpc",
			},
			"dirPath",
			"00008026.jpg",
		)

		assert.NoError(test, err)
		rm.AssertExpectations(test)
		sm.AssertExpectations(test)
	})

	test.Run("失敗：使用請求的參數不支援，並且沒有返回錯誤。", func(test *testing.T) {
		rm := new(PredictRepositoryMock)
		sm := new(CreateServiceMock)

		createServiceWithMockData := NewCreateService(rm, sm)
		_, err := createServiceWithMockData.Create(
			&model.PredictForm{
				ClinicId: 1,
				Method:   "smtp",
			},
			"dirPath",
			"00008026.jpg",
		)

		assert.EqualError(test, err, "參數只接受：http和grpc")
	})

	test.Run("失敗：使用其中一種請求，但是Python的Service報錯。", func(test *testing.T) {
		rm := new(PredictRepositoryMock)
		sm := new(CreateServiceMock)

		// 因為報錯，所以返回的Result是空物件
		sm.On("HttpReq", "dirPath", httpPythonUrl).Return(&model.Result{}, errors.New("Error came from python: "))

		createServiceWithMockData := NewCreateService(rm, sm)
		_, err := createServiceWithMockData.Create(
			&model.PredictForm{
				ClinicId: 1,
				Method:   "http",
			},
			"dirPath",
			"00008026.jpg",
		)

		assert.EqualError(test, err, "Error came from python: ")
		sm.AssertExpectations(test)
	})

	test.Run("失敗：使用其中一種請求，但是Python返回的字串格式有誤。", func(test *testing.T) {
		rm := new(PredictRepositoryMock)
		sm := new(CreateServiceMock)

		// 因為報錯，所以返回的Result是空物件
		sm.On("HttpReq", "dirPath", httpPythonUrl).Return(&model.Result{
			IsSuccessful: true,
			Msg:          "success",
			Predict:      `{"123":123,,,,,}`,
		}, nil)

		createServiceWithMockData := NewCreateService(rm, sm)
		_, err := createServiceWithMockData.Create(
			&model.PredictForm{
				ClinicId: 1,
				Method:   "http",
			},
			"dirPath",
			"00008026.jpg",
		)
		assert.EqualError(test, err, "invalid character ',' looking for beginning of object key string")
		sm.AssertExpectations(test)
	})
}

func TestHttpReq(test *testing.T) {
	jsonStringSuccess := `{"isSuccessful": true,"msg": "成功","predict": "predict"}`
	jsonStringFailed := `{"isSuccessful": false,"msg": "不知道怎麼了","predict": ""}`

	test.Run("成功：測試HttpReq，接收Python正確返回資訊。", func(test *testing.T) {
		rm := new(PredictRepositoryMock)
		sm := new(CreateServiceMock)

		s := httptest.NewServer(
			http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.Write([]byte(jsonStringSuccess))
			}),
		)
		defer s.Close()

		createServiceWithMockData := NewCreateService(rm, sm)
		result, err := createServiceWithMockData.HttpReq("dirPath", s.URL)
		assert.NoError(test, err)
		assert.Equal(test, result, &model.Result{
			IsSuccessful: true,
			Msg:          "成功",
			Predict:      "predict",
		})
	})

	test.Run("失敗：測試HttpReq，接收Python返回報錯資訊。", func(test *testing.T) {
		rm := new(PredictRepositoryMock)
		sm := new(CreateServiceMock)

		s := httptest.NewServer(
			http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.Write([]byte(jsonStringFailed))
			}),
		)
		defer s.Close()

		createServiceWithMockData := NewCreateService(rm, sm)
		result, err := createServiceWithMockData.HttpReq("dirPath", s.URL)

		assert.EqualError(test, err, "Error came from python: 不知道怎麼了")
		assert.Equal(test, result, &model.Result{
			IsSuccessful: false,
			Msg:          "",
			Predict:      "",
		})
	})
}

func TestGrpcReq(test *testing.T) {
	test.Run("失敗：測試GrpcReq，但是Server connection refused。", func(test *testing.T) {
		rm := new(PredictRepositoryMock)
		sm := new(CreateServiceMock)
		createServiceWithMockData := NewCreateService(rm, sm)
		result, err := createServiceWithMockData.GrpcReq("dirPath", "127.0.0.1:5001")
		assert.EqualError(test, err, `Error came from python: rpc error: code = Unavailable desc = connection error: desc = "transport: Error while dialing dial tcp 127.0.0.1:5001: connect: connection refused"`)
		assert.Equal(test, result, &model.Result{
			IsSuccessful: false,
			Msg:          "",
			Predict:      "",
		})
	})

	test.Run("成功：測試HttpReq，接收Python正確返回資訊。", func(test *testing.T) {
		rm := new(PredictRepositoryMock)
		sm := new(CreateServiceMock)
		ps := new(PredictServer)

		// STEP 2-1：定義要監聽的 port 號
		lis, err := net.Listen("tcp", grpcPythonUrl)
		if err != nil {
			log.Fatalf("failed to listed: %v", err)
		}
		// STEP 2-2：使用 gRPC 的 NewServer 方法來建立 gRPC Server 的實例
		grpcServer := grpc.NewServer()

		// STEP 2-3：在 gRPC Server 中註冊 service 的實作
		// 使用 proto 提供的 RegisterRouteGuideServer 方法，並將 routeGuideServer 作為參數傳入
		pb.RegisterPredictServer(grpcServer, ps)

		// STEP 2-4：啟動 grpcServer，並阻塞在這裡直到該程序被 kill 或 stop
		go func() {
			err = grpcServer.Serve(lis)
			if err != nil {
				panic(err)
			}
		}()
		defer grpcServer.Stop()

		ps.On("Predict",
			mock.AnythingOfType("*context.valueCtx"),
			&pb.Request{Dir: "dirPath"},
		).Return(&pb.Response{
			IsSuccessful: true,
			Msg:          "成功",
			Predict:      `predict`,
		}, nil)

		createServiceWithMockData := NewCreateService(rm, sm)
		result, err := createServiceWithMockData.GrpcReq("dirPath", grpcPythonUrl)
		assert.NoError(test, err)
		assert.Equal(test, result, &model.Result{
			IsSuccessful: true,
			Msg:          "成功",
			Predict:      "predict",
		})
	})

	test.Run("成功：測試HttpReq，接收Python正確返回資訊。", func(test *testing.T) {
		rm := new(PredictRepositoryMock)
		sm := new(CreateServiceMock)
		ps := new(PredictServer)

		// STEP 2-1：定義要監聽的 port 號
		lis, err := net.Listen("tcp", grpcPythonUrl)
		if err != nil {
			log.Fatalf("failed to listed: %v", err)
		}
		// STEP 2-2：使用 gRPC 的 NewServer 方法來建立 gRPC Server 的實例
		grpcServer := grpc.NewServer()

		// STEP 2-3：在 gRPC Server 中註冊 service 的實作
		// 使用 proto 提供的 RegisterRouteGuideServer 方法，並將 routeGuideServer 作為參數傳入
		pb.RegisterPredictServer(grpcServer, ps)

		// STEP 2-4：啟動 grpcServer，並阻塞在這裡直到該程序被 kill 或 stop
		go func() {
			err = grpcServer.Serve(lis)
			if err != nil {
				panic(err)
			}
		}()
		defer grpcServer.Stop()
		ps.On("Predict",
			mock.AnythingOfType("*context.valueCtx"),
			&pb.Request{Dir: "dirPath"},
		).Return(&pb.Response{
			IsSuccessful: false,
			Msg:          "有問題",
			Predict:      "",
		}, nil)

		createServiceWithMockData := NewCreateService(rm, sm)
		result, err := createServiceWithMockData.GrpcReq("dirPath", grpcPythonUrl)
		assert.EqualError(test, err, "Error came from python: 有問題")
		assert.Equal(test, result, &model.Result{
			IsSuccessful: false,
			Msg:          "",
			Predict:      "",
		})
	})
}
