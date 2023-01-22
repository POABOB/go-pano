package service

import (
	"go-pano/domain/model"
	repository "go-pano/domain/repository/user"
)

// interface
type IPasswordService interface {
	UpdatePassword(*model.UserPasswordForm) error
}

// 實例化
func NewPasswordService(userRepository repository.IUserRepository) IPasswordService {
	return &PasswordService{userRepository}
}

// class
type PasswordService struct {
	userRepository repository.IUserRepository
}

// 這個Service就是為了call Ai辨識，並將結果存入DB
func (s *PasswordService) UpdatePassword(user *model.UserPasswordForm) error {
	if err := s.userRepository.UpdatePassword(user); err != nil {
		return err
	}

	return nil
}
