package service

import (
	"go-pano/domain/model"
	repository "go-pano/domain/repository/user"

	"github.com/goccy/go-json"
)

// interface
type ICreateService interface {
	Create(*model.User) error
}

// 實例化
func NewCreateService(userRepository repository.IUserRepository) ICreateService {
	return &CreateService{userRepository}
}

// class
type CreateService struct {
	userRepository repository.IUserRepository
}

// 這個Service就是為了call Ai辨識，並將結果存入DB
func (s *CreateService) Create(user *model.User) error {
	// Roles to RolesString
	jsonString, err := json.Marshal(user.Roles)
	if err != nil {
		return err
	}
	user.RolesString = string(jsonString)

	// 插入
	if err := s.userRepository.Create(user); err != nil {
		return err
	}

	return nil
}
