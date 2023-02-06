package user_repository

import (
	"errors"
	"go-pano/domain/model"
	"go-pano/utils"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// 接口
type IUserRepository interface {
	GetAll() ([]model.User, error)
	Update(*model.UserUpdateForm) error
	UpdateAccount(*model.UserUpdateAccountForm) error
	Create(*model.UserCreateForm) error
	UpdatePassword(*model.UserPasswordForm, string) error
	Login(*model.UserLoginForm) (*model.User, error)
	Delete(*model.UserDeleteForm) error
}

// gorm使用參考
// https://pjchender.dev/golang/note-gorm-example/
// new Class()
func NewUserRepository(db utils.IDBInstance) IUserRepository {
	return &UserRepository{mysql: db}
}

// class
type UserRepository struct {
	mysql utils.IDBInstance
}

// 獲取全部User
func (ctrl *UserRepository) GetAll() ([]model.User, error) {
	user := []model.User{}
	if err := ctrl.mysql.DB().Limit(500).Find(&user).Error; err != nil {
		return []model.User{}, err
	}

	return user, nil
}

// 插入User
func (ctrl *UserRepository) Create(user *model.UserCreateForm) error {
	encrypted, _ := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	user.Password = string(encrypted)
	if db := ctrl.mysql.DB().Table("Users").Create(user); db.Error != nil {
		return db.Error
	}

	return nil
}

// 更新User
func (ctrl *UserRepository) Update(user *model.UserUpdateForm) error {
	if db := ctrl.mysql.DB().Table("Users").Updates(user); db.Error != nil {
		return db.Error
	}

	return nil
}

func (ctrl *UserRepository) UpdateAccount(user *model.UserUpdateAccountForm) error {
	if db := ctrl.mysql.DB().Table("Users").Updates(user); db.Error != nil {
		return db.Error
	}

	return nil
}

// 刪除/復原User
func (ctrl *UserRepository) UpdatePassword(user *model.UserPasswordForm, obj string) error {
	encrypted, _ := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)

	var db *gorm.DB
	if obj == "self" {
		u := &model.User{}
		if err := ctrl.mysql.DB().Where("user_id = ?", user.UserId).First(&u).Error; err != nil {
			return errors.New("不是自己的id")
		}

		if err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(user.OldPass)); err != nil {
			return errors.New("密碼錯誤")
		}
		// 更新自己
		db = ctrl.mysql.DB().Table("Users").Where("user_id = ?", user.UserId).Update("password", string(encrypted))
	} else {
		// 更新全部
		db = ctrl.mysql.DB().Table("Users").Where("user_id = ?", user.UserId).Update("password", string(encrypted))
	}

	if db.Error != nil {
		return db.Error
	}

	return nil
}

func (ctrl *UserRepository) Login(userForm *model.UserLoginForm) (*model.User, error) {
	user := &model.User{}
	if err := ctrl.mysql.DB().Where("account=?", userForm.Account).First(&user).Error; err != nil {
		return &model.User{}, errors.New("帳號或密碼錯誤")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(userForm.Password)); err != nil {
		return &model.User{}, errors.New("帳號或密碼錯誤")
	}

	if user.Status == -1 {
		return &model.User{}, errors.New("該帳號已被停用")
	}

	return user, nil
}

func (ctrl *UserRepository) Delete(user *model.UserDeleteForm) error {
	if db := ctrl.mysql.DB().Table("Users").Delete(user); db.Error != nil {
		return db.Error
	}

	return nil
}
