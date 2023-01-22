package service

import (
	"encoding/json"
	"fmt"
	"go-pano/domain/model"
	repository "go-pano/domain/repository/user"
)

// interface
type IGetService interface {
	Get() ([]model.User, error)
}

// 實例化
func NewGetService(userRepository repository.IUserRepository) IGetService {
	return &GetService{userRepository}
}

// class
type GetService struct {
	userRepository repository.IUserRepository
}

// 這個Service就是為了call Ai辨識，並將結果存入DB
func (s *GetService) Get() ([]model.User, error) {
	result, err := s.userRepository.GetAll()
	if err != nil {
		return []model.User{}, err
	}

	// 將string轉[]string
	for key, val := range result {
		var p []string
		if err := json.Unmarshal([]byte(val.RolesString), &p); err != nil {
			fmt.Println(err.Error())
			return []model.User{}, err
		}
		result[key].Roles = p
		result[key].RolesString = ""
	}

	return result, nil
}
