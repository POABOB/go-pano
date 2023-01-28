package user_service

import (
	"encoding/json"
	"fmt"
	"go-pano/domain/model"
	user_repository "go-pano/domain/repository/user"
	"go-pano/utils"
)

// interface
type IUserService interface {
	Get() ([]model.User, error)
	Create(*model.UserCreateForm) error
	Update(*model.UserUpdateForm) error
	UpdatePassword(*model.UserPasswordForm) error
	Login(*model.UserLoginForm) (string, error)
}

// 實例化
func NewUserService(r user_repository.IUserRepository) IUserService {
	return &UserService{UserRepository: r}
}

// class
type UserService struct {
	UserRepository user_repository.IUserRepository
}

func (us *UserService) Get() ([]model.User, error) {
	result, err := us.UserRepository.GetAll()
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

func (us *UserService) Create(user *model.UserCreateForm) error {
	// Roles to RolesString
	jsonString, err := json.Marshal(user.Roles)
	if err != nil {
		return err
	}
	user.RolesString = string(jsonString)

	// 插入
	if err := us.UserRepository.Create(user); err != nil {
		return err
	}

	return nil
}

func (us *UserService) Update(user *model.UserUpdateForm) error {
	// Roles to RolesString
	jsonString, err := json.Marshal(user.Roles)
	if err != nil {
		return err
	}
	user.RolesString = string(jsonString)

	if err := us.UserRepository.Update(user); err != nil {
		return err
	}

	return nil
}

func (us *UserService) UpdatePassword(user *model.UserPasswordForm) error {
	if err := us.UserRepository.UpdatePassword(user); err != nil {
		return err
	}

	return nil
}

func (us *UserService) Login(user *model.UserLoginForm) (string, error) {
	result, err := us.UserRepository.Login(user)
	if err != nil {
		return "", err
	}

	// 將string轉[]string
	var p []string
	if err := json.Unmarshal([]byte(result.RolesString), &p); err != nil {
		fmt.Println(err.Error())
		return "", err
	}
	result.Roles = p
	result.RolesString = ""

	// 簽發Token
	tokenString, _ := utils.GenToken(result.UserId, result.Name, result.Roles)

	return tokenString, nil
}
