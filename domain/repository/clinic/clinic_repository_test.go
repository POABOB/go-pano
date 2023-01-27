package clinic_repository

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

type ClinicRepositorySuite struct {
	DB         *gorm.DB
	mock       sqlmock.Sqlmock
	repository IClinicRepository
}

// https://github.com/go-gorm/gorm/issues/3565

func TestClinicRepository(test *testing.T) {
	s := &ClinicRepositorySuite{}
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
	s.repository = NewClinicRepository(mockInstance)

	var getClinic []model.Clinic = []model.Clinic{
		{
			ClinicTokenForm: model.ClinicTokenForm{
				ClinicId: 1,
			},
			ClinicCreateForm: model.ClinicCreateForm{
				Name:          "診所",
				StartAt:       "2022-10-31",
				EndAt:         "9999-12-31",
				QuotaPerMonth: 100,
			},
			Token: "token",
		},
	}

	// GetAll()
	test.Run("成功：GetAll，有找到資料。", func(test *testing.T) {
		s.mock.ExpectQuery(regexp.QuoteMeta(
			"SELECT `clinic_id`,SUBSTR(start_at,1,10) AS start_at,SUBSTR(end_at,1,10) AS end_at,`name`,`quota_per_month`,`token` FROM `Clinic` LIMIT 500")).
			WithArgs().
			WillReturnRows(sqlmock.NewRows([]string{"clinic_id", "name", "start_at", "end_at", "quota_per_month", "token"}).
				AddRow(getClinic[0].ClinicId, getClinic[0].Name, getClinic[0].StartAt, getClinic[0].EndAt, getClinic[0].QuotaPerMonth, getClinic[0].Token))

		res, err := s.repository.GetAll()

		assert.NoError(test, err)
		assert.True(test, reflect.DeepEqual(getClinic, res))
	})

	test.Run("成功：GetAll，沒有找到資料。", func(test *testing.T) {
		s.mock.ExpectQuery(regexp.QuoteMeta(
			"SELECT `clinic_id`,SUBSTR(start_at,1,10) AS start_at,SUBSTR(end_at,1,10) AS end_at,`name`,`quota_per_month`,`token` FROM `Clinic` LIMIT 500")).
			WithArgs().
			WillReturnRows(sqlmock.NewRows([]string{"clinic_id", "name", "start_at", "end_at", "quota_per_month", "token"}))

		res, err := s.repository.GetAll()

		assert.NoError(test, err)
		assert.True(test, reflect.DeepEqual([]model.Clinic{}, res))

	})

	test.Run("失敗：GetAll，有錯誤發生。", func(test *testing.T) {
		s.mock.ExpectQuery(regexp.QuoteMeta(
			"SELECT `clinic_id`,SUBSTR(start_at,1,10) AS start_at,SUBSTR(end_at,1,10) AS end_at,`name`,`quota_per_month`,`token` FROM `Clinic` LIMIT 500")).
			WillReturnError(errors.New("有錯誤發生"))

		res, err := s.repository.GetAll()

		assert.EqualError(test, err, "有錯誤發生")
		assert.True(test, reflect.DeepEqual([]model.Clinic{}, res))

	})

	// Create()
	test.Run("成功：Create成功插入。", func(test *testing.T) {
		s.mock.ExpectBegin()
		s.mock.ExpectExec(
			regexp.QuoteMeta("INSERT INTO `Clinic` (`name`,`start_at`,`end_at`,`quota_per_month`,`token`) VALUES (?,?,?,?,?)")).
			WithArgs(getClinic[0].Name, getClinic[0].StartAt, getClinic[0].EndAt, getClinic[0].QuotaPerMonth, sqlmock.AnyArg()).
			WillReturnResult(sqlmock.NewResult(1, 1))
		s.mock.ExpectCommit()

		err := s.repository.Create(
			&model.ClinicCreateForm{
				Name:          getClinic[0].Name,
				StartAt:       getClinic[0].StartAt,
				EndAt:         getClinic[0].EndAt,
				QuotaPerMonth: getClinic[0].QuotaPerMonth,
			}, "token")

		assert.NoError(test, err)
	})

	test.Run("失敗：Create，有錯誤發生。", func(test *testing.T) {
		s.mock.ExpectBegin()
		s.mock.ExpectExec(
			regexp.QuoteMeta("INSERT INTO `Clinic` (`name`,`start_at`,`end_at`,`quota_per_month`,`token`) VALUES (?,?,?,?,?)")).
			WithArgs(getClinic[0].Name, getClinic[0].StartAt, getClinic[0].EndAt, getClinic[0].QuotaPerMonth, sqlmock.AnyArg()).
			WillReturnError(errors.New("有錯誤發生"))
		s.mock.ExpectRollback()
		err := s.repository.Create(&model.ClinicCreateForm{
			Name:          getClinic[0].Name,
			StartAt:       getClinic[0].StartAt,
			EndAt:         getClinic[0].EndAt,
			QuotaPerMonth: getClinic[0].QuotaPerMonth,
		}, "token")

		assert.EqualError(test, err, "有錯誤發生")
	})

	// Update()
	test.Run("成功：Update成功更新。", func(test *testing.T) {
		s.mock.ExpectBegin()
		s.mock.ExpectExec(
			regexp.QuoteMeta("UPDATE `Clinic` SET `name`=?,`start_at`=?,`end_at`=?,`quota_per_month`=? WHERE clinic_id = ?")).
			WithArgs(getClinic[0].Name, getClinic[0].StartAt, getClinic[0].EndAt, getClinic[0].QuotaPerMonth, getClinic[0].ClinicId).
			WillReturnResult(sqlmock.NewResult(0, 1))
		s.mock.ExpectCommit()

		err := s.repository.Update(
			&model.ClinicUpdateForm{
				ClinicTokenForm: model.ClinicTokenForm{
					ClinicId: 1,
				},
				ClinicCreateForm: model.ClinicCreateForm{
					Name:          getClinic[0].Name,
					StartAt:       getClinic[0].StartAt,
					EndAt:         getClinic[0].EndAt,
					QuotaPerMonth: getClinic[0].QuotaPerMonth,
				},
			})

		assert.NoError(test, err)
	})

	test.Run("失敗：Update，沒有更新到資料。", func(test *testing.T) {
		s.mock.ExpectBegin()
		s.mock.ExpectExec(
			regexp.QuoteMeta("UPDATE `Clinic` SET `name`=?,`start_at`=?,`end_at`=?,`quota_per_month`=? WHERE clinic_id = ?")).
			WithArgs(getClinic[0].Name, getClinic[0].StartAt, getClinic[0].EndAt, getClinic[0].QuotaPerMonth, getClinic[0].ClinicId).
			WillReturnResult(sqlmock.NewResult(0, 0))
		s.mock.ExpectCommit()

		err := s.repository.Update(
			&model.ClinicUpdateForm{
				ClinicTokenForm: model.ClinicTokenForm{
					ClinicId: 1,
				},
				ClinicCreateForm: model.ClinicCreateForm{
					Name:          getClinic[0].Name,
					StartAt:       getClinic[0].StartAt,
					EndAt:         getClinic[0].EndAt,
					QuotaPerMonth: getClinic[0].QuotaPerMonth,
				},
			})

		assert.EqualError(test, err, "Failed")
	})

	test.Run("失敗：Update，有錯誤發生。", func(test *testing.T) {
		s.mock.ExpectBegin()
		s.mock.ExpectExec(
			regexp.QuoteMeta("UPDATE `Clinic` SET `name`=?,`start_at`=?,`end_at`=?,`quota_per_month`=? WHERE clinic_id = ?")).
			WithArgs(getClinic[0].Name, getClinic[0].StartAt, getClinic[0].EndAt, getClinic[0].QuotaPerMonth, getClinic[0].ClinicId).
			WillReturnError(errors.New("有錯誤發生"))
		s.mock.ExpectRollback()

		err := s.repository.Update(
			&model.ClinicUpdateForm{
				ClinicTokenForm: model.ClinicTokenForm{
					ClinicId: 1,
				},
				ClinicCreateForm: model.ClinicCreateForm{
					Name:          getClinic[0].Name,
					StartAt:       getClinic[0].StartAt,
					EndAt:         getClinic[0].EndAt,
					QuotaPerMonth: getClinic[0].QuotaPerMonth,
				},
			})

		assert.EqualError(test, err, "有錯誤發生")
	})

	// UpdateToken()
	test.Run("成功：UpdateToken成功更新。", func(test *testing.T) {
		s.mock.ExpectBegin()
		s.mock.ExpectExec(
			regexp.QuoteMeta("UPDATE `Clinic` SET `token`=? WHERE clinic_id = ?")).
			WithArgs(sqlmock.AnyArg(), getClinic[0].ClinicId).
			WillReturnResult(sqlmock.NewResult(0, 1))
		s.mock.ExpectCommit()

		err := s.repository.UpdateToken(
			&model.ClinicTokenForm{
				ClinicId: 1,
			}, "token")

		assert.NoError(test, err)
	})

	test.Run("失敗：UpdateToken，沒有更新到資料。", func(test *testing.T) {
		s.mock.ExpectBegin()
		s.mock.ExpectExec(
			regexp.QuoteMeta("UPDATE `Clinic` SET `token`=? WHERE clinic_id = ?")).
			WithArgs(sqlmock.AnyArg(), getClinic[0].ClinicId).
			WillReturnResult(sqlmock.NewResult(0, 0))
		s.mock.ExpectCommit()

		err := s.repository.UpdateToken(
			&model.ClinicTokenForm{
				ClinicId: 1,
			}, "token")

		assert.EqualError(test, err, "Failed")
	})

	test.Run("失敗：UpdateToken，有錯誤發生。", func(test *testing.T) {
		s.mock.ExpectBegin()
		s.mock.ExpectExec(
			regexp.QuoteMeta("UPDATE `Clinic` SET `token`=? WHERE clinic_id = ?")).
			WithArgs(sqlmock.AnyArg(), getClinic[0].ClinicId).
			WillReturnError(errors.New("有錯誤發生"))
		s.mock.ExpectRollback()

		err := s.repository.UpdateToken(
			&model.ClinicTokenForm{
				ClinicId: 1,
			}, "token")

		assert.EqualError(test, err, "有錯誤發生")
	})
}
