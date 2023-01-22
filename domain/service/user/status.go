package service

import (
	"go-pano/domain/model"
	repository "go-pano/domain/repository/user"
)

// interface
type IStatusService interface {
	UpdateStatus(*model.UserStatusForm) error
}

// 實例化
func NewStatusService(userRepository repository.IUserRepository) IStatusService {
	return &StatusService{userRepository}
}

// class
type StatusService struct {
	userRepository repository.IUserRepository
}

// 這個Service就是為了call Ai辨識，並將結果存入DB
func (s *StatusService) UpdateStatus(user *model.UserStatusForm) error {
	if err := s.userRepository.UpdateStatus(user); err != nil {
		return err
	}

	return nil
}
