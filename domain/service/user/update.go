package service

import (
	"go-pano/domain/model"
	repository "go-pano/domain/repository/user"

	"github.com/goccy/go-json"
)

// interface
type IUpdateService interface {
	Update(*model.UserUpdateForm) error
}

// 實例化
func NewUpdateService(userRepository repository.IUserRepository) IUpdateService {
	return &UpdateService{userRepository}
}

// class
type UpdateService struct {
	userRepository repository.IUserRepository
}

// 這個Service就是為了call Ai辨識，並將結果存入DB
func (s *UpdateService) Update(user *model.UserUpdateForm) error {
	// Roles to RolesString
	jsonString, err := json.Marshal(user.Roles)
	if err != nil {
		return err
	}
	user.RolesString = string(jsonString)

	if err := s.userRepository.Update(user); err != nil {
		return err
	}

	return nil
}
