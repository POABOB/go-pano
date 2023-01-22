package service

import (
	"context"
	"encoding/json"
	"errors"

	"go-pano/config"
	"go-pano/domain/model"
	repository "go-pano/domain/repository/predict"
	pb "go-pano/protos/predict"
	"go-pano/utils"

	"github.com/go-resty/resty/v2"
	"google.golang.org/grpc"
)

// interface
type ICreateService interface {
	Create(*model.PredictForm, string, string) (interface{}, error)
	GrpcReq(string, string) (*model.Result, error)
	HttpReq(string, string) (*model.Result, error)
}

// 實例化
func NewCreateService(predictRepository repository.IPredictRepository, self ICreateService) ICreateService {
	return &CreateService{predictRepository, self}
}

// class
type CreateService struct {
	PredictRepository repository.IPredictRepository
	self              ICreateService
}

// 這個Service就是為了call Ai辨識，並將結果存入DB
func (s *CreateService) Create(predictForm *model.PredictForm, dir string, fileName string) (interface{}, error) {
	// 判斷使用什麼方式
	var result *model.Result
	var _err error
	if predictForm.Method == "http" {
		result, _err = s.self.HttpReq(dir, "http://"+config.PythonHost+":5000/")
	} else if predictForm.Method == "grpc" {
		result, _err = s.self.GrpcReq(dir, config.PythonHost+":5001")
	} else {
		return nil, errors.New("參數只接受：http和grpc")
	}

	if _err != nil {
		return nil, _err
	}

	// 字串轉成JSON
	var p interface{}
	if err := json.Unmarshal([]byte(result.Predict), &p); err != nil {
		return nil, err
	}

	// 插入DB
	if err := s.PredictRepository.Create(predictForm.ClinicId, fileName, result.Predict); err != nil {
		return nil, err
	}

	return p, nil
}

// TEST: GRPC，與 Python Server 溝通
func (s *CreateService) GrpcReq(dir string, url string) (*model.Result, error) {

	// 建立GRPC連線資訊
	conn, err := grpc.Dial(url, grpc.WithInsecure())
	if err != nil {
		return &model.Result{}, errors.New("gRPC Dial失敗: " + err.Error())
	}
	// 最後請關閉連線
	defer conn.Close()

	// 使用Client連接，並call函數
	client := pb.NewPredictClient(conn)
	pbResult, err := client.Predict(context.Background(), &pb.Request{Dir: dir})
	if err != nil {
		return &model.Result{}, errors.New("Error came from python: " + err.Error())
	}

	if !pbResult.IsSuccessful {
		return &model.Result{}, errors.New("Error came from python: " + pbResult.Msg)
	}

	return &model.Result{
		IsSuccessful: pbResult.IsSuccessful,
		Msg:          pbResult.Msg,
		Predict:      pbResult.Predict,
	}, nil
}

func (s *CreateService) HttpReq(dir string, url string) (*model.Result, error) {
	result := &model.Result{}
	client := resty.New()
	client.R().
		SetResult(&result).
		SetQueryString("Dir=" + dir).
		ForceContentType("application/json").
		Get(url)
	if !result.IsSuccessful {
		utils.LogInstance.Error("Error came from python: " + result.Msg)
		return &model.Result{}, errors.New("Error came from python: " + result.Msg)
	}

	return result, nil
}
