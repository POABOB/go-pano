package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"sync"
	"testing"

	"go-pano/config"
	"go-pano/domain/model"
	"go-pano/utils"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func init() {
	gin.SetMode(gin.TestMode)
	config.LoadConfigTest()
	utils.Reset()
}

type UserRepositorySuite struct {
	DB         *gorm.DB
	mock       sqlmock.Sqlmock
	repository IUserRepository
}

// https://github.com/go-gorm/gorm/issues/3565

func TestUserRepository(test *testing.T) {
	s := &UserRepositorySuite{}
	var (
		db  *sql.DB
		err error
	)
	db, s.mock, err = sqlmock.New()
	if err != nil {
		fmt.Println(err)
	}

	dialector := mysql.New(mysql.Config{
		DSN:                       "sqlmock_db",
		DriverName:                "mysql",
		Conn:                      db,
		SkipInitializeWithVersion: true,
	})
	s.DB, err = gorm.Open(dialector, &gorm.Config{})

	if err != nil {
		fmt.Println(err)
	}

	var o sync.Once
	mockInstance := utils.NewMockInstance(s.DB, o)

	// 依賴注入
	s.repository = NewUserRepository(mockInstance)

	var getUser []model.User = []model.User{
		{
			UserId:      1,
			Name:        "User1",
			Account:     "user1",
			RolesString: `["admin"]`,
			Status:      1,
		},
	}

	// GetAll()
	test.Run("成功：GetAll，有找到資料。", func(test *testing.T) {
		s.mock.ExpectQuery(regexp.QuoteMeta(
			"SELECT * FROM `Users` LIMIT 500")).
			WithArgs().
			WillReturnRows(sqlmock.NewRows([]string{"user_id", "name", "account", "roles_string", "status"}).
				AddRow(getUser[0].UserId,
					getUser[0].Name,
					getUser[0].Account,
					getUser[0].RolesString,
					getUser[0].Status,
				))

		res, err := s.repository.GetAll()

		assert.NoError(test, err)
		assert.True(test, reflect.DeepEqual(getUser, res))
	})

	test.Run("成功：GetAll，沒有找到資料。", func(test *testing.T) {
		s.mock.ExpectQuery(regexp.QuoteMeta(
			"SELECT * FROM `Users` LIMIT 500")).
			WithArgs().
			WillReturnRows(sqlmock.NewRows([]string{"user_id", "name", "account", "roles_string", "status"}))

		res, err := s.repository.GetAll()

		assert.NoError(test, err)
		assert.True(test, reflect.DeepEqual([]model.User{}, res))

	})

	test.Run("失敗：GetAll，有錯誤發生。", func(test *testing.T) {
		s.mock.ExpectQuery(regexp.QuoteMeta(
			"SELECT * FROM `Users` LIMIT 500")).
			WillReturnError(errors.New("有錯誤發生"))

		res, err := s.repository.GetAll()

		assert.EqualError(test, err, "有錯誤發生")
		assert.True(test, reflect.DeepEqual([]model.User{}, res))

	})

	// Create()
	test.Run("成功：Create成功插入。", func(test *testing.T) {
		s.mock.ExpectBegin()
		s.mock.ExpectExec(
			regexp.QuoteMeta("INSERT INTO `Users` (`name`,`account`,`roles_string`,`password`,`status`) VALUES (?,?,?,?,?)")).
			WithArgs(getUser[0].Name, getUser[0].Account, getUser[0].RolesString, sqlmock.AnyArg(), getUser[0].Status).
			WillReturnResult(sqlmock.NewResult(1, 1))
		s.mock.ExpectCommit()

		err := s.repository.Create(
			&model.User{
				UserId:      1,
				Name:        "User1",
				Account:     "user1",
				Password:    "ppaass",
				Passconf:    "ppaass",
				Roles:       []string{"admin"},
				RolesString: `["admin"]`,
				Status:      1,
			})

		assert.NoError(test, err)
	})

	test.Run("失敗：Create，有錯誤發生。", func(test *testing.T) {
		s.mock.ExpectBegin()
		s.mock.ExpectExec(
			regexp.QuoteMeta("INSERT INTO `Users` (`name`,`account`,`roles_string`,`password`,`status`) VALUES (?,?,?,?,?)")).
			WithArgs(getUser[0].Name, getUser[0].Account, getUser[0].RolesString, sqlmock.AnyArg(), getUser[0].Status).
			WillReturnError(errors.New("有錯誤發生"))
		s.mock.ExpectRollback()
		err := s.repository.Create(&model.User{
			UserId:      1,
			Name:        "User1",
			Account:     "user1",
			Password:    "ppaass",
			Passconf:    "ppaass",
			Roles:       []string{"admin"},
			RolesString: `["admin"]`,
			Status:      1,
		})

		assert.EqualError(test, err, "有錯誤發生")
	})

	// Update()
	test.Run("成功：Update成功更新。", func(test *testing.T) {
		s.mock.ExpectBegin()
		s.mock.ExpectExec(
			regexp.QuoteMeta("UPDATE `Users` SET `name`=?,`account`=?,`roles_string`=? WHERE `user_id` = ?")).
			WithArgs(getUser[0].Name, getUser[0].Account, getUser[0].RolesString, getUser[0].UserId).
			WillReturnResult(sqlmock.NewResult(0, 1))
		s.mock.ExpectCommit()

		err := s.repository.Update(
			&model.UserUpdateForm{
				UserId:      1,
				Name:        "User1",
				Account:     "user1",
				Roles:       []string{"admin"},
				RolesString: `["admin"]`,
			})

		assert.NoError(test, err)
	})

	test.Run("失敗：Update，沒有更新到資料。", func(test *testing.T) {
		s.mock.ExpectBegin()
		s.mock.ExpectExec(
			regexp.QuoteMeta("UPDATE `Users` SET `name`=?,`account`=?,`roles_string`=? WHERE `user_id` = ?")).
			WithArgs(getUser[0].Name, getUser[0].Account, getUser[0].RolesString, getUser[0].UserId).
			WillReturnResult(sqlmock.NewResult(0, 0))
		s.mock.ExpectCommit()

		err := s.repository.Update(
			&model.UserUpdateForm{
				UserId:      1,
				Name:        "User1",
				Account:     "user1",
				Roles:       []string{"admin"},
				RolesString: `["admin"]`,
			})

		assert.EqualError(test, err, "Failed")
	})

	test.Run("失敗：Update，有錯誤發生。", func(test *testing.T) {
		s.mock.ExpectBegin()
		s.mock.ExpectExec(
			regexp.QuoteMeta("UPDATE `Users` SET `name`=?,`account`=?,`roles_string`=? WHERE `user_id` = ?")).
			WithArgs(getUser[0].Name, getUser[0].Account, getUser[0].RolesString, getUser[0].UserId).
			WillReturnError(errors.New("有錯誤發生"))
		s.mock.ExpectRollback()

		err := s.repository.Update(
			&model.UserUpdateForm{
				UserId:      1,
				Name:        "User1",
				Account:     "user1",
				Roles:       []string{"admin"},
				RolesString: `["admin"]`,
			})

		assert.EqualError(test, err, "有錯誤發生")
	})

	// UpdateStatus()
	test.Run("成功：UpdateStatus成功更新。", func(test *testing.T) {
		s.mock.ExpectBegin()
		s.mock.ExpectExec(
			regexp.QuoteMeta("UPDATE `Users` SET `status`=? WHERE user_id = ?")).
			WithArgs(getUser[0].Status, getUser[0].UserId).
			WillReturnResult(sqlmock.NewResult(0, 1))
		s.mock.ExpectCommit()

		err := s.repository.UpdateStatus(
			&model.UserStatusForm{
				UserId: 1,
				Status: 1,
			})

		assert.NoError(test, err)
	})

	test.Run("失敗：UpdateStatus，沒有更新到資料。", func(test *testing.T) {
		s.mock.ExpectBegin()
		s.mock.ExpectExec(
			regexp.QuoteMeta("UPDATE `Users` SET `status`=? WHERE user_id = ?")).
			WithArgs(getUser[0].Status, getUser[0].UserId).
			WillReturnResult(sqlmock.NewResult(0, 0))
		s.mock.ExpectCommit()

		err := s.repository.UpdateStatus(
			&model.UserStatusForm{
				UserId: 1,
				Status: 1,
			})

		assert.EqualError(test, err, "Failed")
	})

	test.Run("失敗：UpdateStatus，有錯誤發生。", func(test *testing.T) {
		s.mock.ExpectBegin()
		s.mock.ExpectExec(
			regexp.QuoteMeta("UPDATE `Users` SET `status`=? WHERE user_id = ?")).
			WithArgs(getUser[0].Status, getUser[0].UserId).
			WillReturnError(errors.New("有錯誤發生"))
		s.mock.ExpectRollback()

		err := s.repository.UpdateStatus(
			&model.UserStatusForm{
				UserId: 1,
				Status: 1,
			})

		assert.EqualError(test, err, "有錯誤發生")
	})

	// UpdatePassword()
	test.Run("成功：UpdatePassword成功更新。", func(test *testing.T) {
		s.mock.ExpectBegin()
		s.mock.ExpectExec(
			regexp.QuoteMeta("UPDATE `Users` SET `password`=? WHERE user_id = ?")).
			WithArgs(sqlmock.AnyArg(), getUser[0].UserId).
			WillReturnResult(sqlmock.NewResult(0, 1))
		s.mock.ExpectCommit()

		err := s.repository.UpdatePassword(
			&model.UserPasswordForm{
				UserId:   1,
				Password: "ppaass",
				Passconf: "ppaass",
			})

		assert.NoError(test, err)
	})

	test.Run("失敗：UpdatePassword，沒有更新到資料。", func(test *testing.T) {
		s.mock.ExpectBegin()
		s.mock.ExpectExec(
			regexp.QuoteMeta("UPDATE `Users` SET `password`=? WHERE user_id = ?")).
			WithArgs(sqlmock.AnyArg(), getUser[0].UserId).
			WillReturnResult(sqlmock.NewResult(0, 0))
		s.mock.ExpectCommit()

		err := s.repository.UpdatePassword(
			&model.UserPasswordForm{
				UserId:   1,
				Password: "ppaass",
				Passconf: "ppaass",
			})

		assert.EqualError(test, err, "Failed")
	})

	test.Run("失敗：UpdatePassword，有錯誤發生。", func(test *testing.T) {
		s.mock.ExpectBegin()
		s.mock.ExpectExec(
			regexp.QuoteMeta("UPDATE `Users` SET `password`=? WHERE user_id = ?")).
			WithArgs(sqlmock.AnyArg(), getUser[0].UserId).
			WillReturnError(errors.New("有錯誤發生"))
		s.mock.ExpectRollback()

		err := s.repository.UpdatePassword(
			&model.UserPasswordForm{
				UserId:   1,
				Password: "ppaass",
				Passconf: "ppaass",
			})

		assert.EqualError(test, err, "有錯誤發生")
	})
}
