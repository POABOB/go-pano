package user_repository

import (
	"errors"
	"go-pano/domain/model"
	"go-pano/utils"

	"golang.org/x/crypto/bcrypt"
)

// 接口
type IUserRepository interface {
	GetAll() ([]model.User, error)
	Update(*model.UserUpdateForm) error
	Create(*model.UserCreateForm) error
	UpdatePassword(*model.UserPasswordForm) error
	Login(*model.UserLoginForm) (*model.User, error)
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
	} else if db.RowsAffected == 0 {
		return utils.ErrFailed
	}

	return nil
}

// 刪除/復原User
func (ctrl *UserRepository) UpdatePassword(user *model.UserPasswordForm) error {
	encrypted, _ := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if db := ctrl.mysql.DB().Table("Users").Where("user_id = ?", user.UserId).Update("password", string(encrypted)); db.Error != nil {
		return db.Error
	} else if db.RowsAffected == 0 {
		return utils.ErrFailed
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

	if user.Status == 0 {
		return &model.User{}, errors.New("該帳號已被停用")
	}

	return user, nil
}
