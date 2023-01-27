package clinic_service

import (
	"errors"
	"testing"

	"go-pano/config"
	"go-pano/domain/model"
	"go-pano/utils"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func init() {
	gin.SetMode(gin.TestMode)
	config.LoadConfigTest()
	utils.Reset()
}

type ClinicRepositoryMock struct {
	mock.Mock
}

func (rm *ClinicRepositoryMock) GetAll() ([]model.Clinic, error) {
	args := rm.Called()
	return args.Get(0).([]model.Clinic), args.Error(1)
}

func (rm *ClinicRepositoryMock) Create(clinic *model.ClinicCreateForm, token string) error {
	args := rm.Called(clinic, token)
	return args.Error(0)
}

func (rm *ClinicRepositoryMock) Update(clinic *model.ClinicUpdateForm) error {
	args := rm.Called(clinic)
	return args.Error(0)
}

func (rm *ClinicRepositoryMock) UpdateToken(clinic *model.ClinicTokenForm, token string) error {
	args := rm.Called(clinic, token)
	return args.Error(0)
}

func TestClinicService(test *testing.T) {
	// GetAll
	test.Run("成功：GetAll()，成功從Repository獲取資料。", func(test *testing.T) {
		urm := new(ClinicRepositoryMock)

		jsonObject := []model.Clinic{
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
		// Mock funcs
		urm.On("GetAll").Return(jsonObject, nil)

		// 將Ｍock注入真的Service
		cs := NewClinicService(urm)
		expected, err := cs.Get()
		assert.NoError(test, err)
		assert.Equal(test, expected, jsonObject)
		urm.AssertExpectations(test)
	})

	test.Run("成功：GetAll()，成功從Repository獲取資料，但是沒有Clinic。", func(test *testing.T) {
		urm := new(ClinicRepositoryMock)

		jsonObject := []model.Clinic{}
		// Mock funcs
		urm.On("GetAll").Return(jsonObject, nil)

		// 將Ｍock注入真的Service
		cs := NewClinicService(urm)
		expected, err := cs.Get()
		assert.NoError(test, err)
		assert.Equal(test, expected, []model.Clinic{})
		urm.AssertExpectations(test)
	})

	test.Run("失敗：GetAll()，DB出現錯誤。", func(test *testing.T) {
		urm := new(ClinicRepositoryMock)

		jsonString := "DB無法連線"
		// Mock funcs
		urm.On("GetAll").Return([]model.Clinic{}, errors.New(jsonString))

		// 將Ｍock注入真的Service
		cs := NewClinicService(urm)
		expected, err := cs.Get()
		assert.EqualError(test, err, jsonString)
		assert.Equal(test, expected, []model.Clinic{})
		// 驗證.On()是否真的有被 call 到
		urm.AssertExpectations(test)
	})

	//  Create
	test.Run("成功：Create()，成功從Repository插入資料。", func(test *testing.T) {
		crm := new(ClinicRepositoryMock)

		jsonObject := &model.ClinicCreateForm{
			Name:          "診所",
			StartAt:       "2022-10-31",
			EndAt:         "9999-12-31",
			QuotaPerMonth: 100,
		}
		// Mock funcs
		crm.On("Create", jsonObject, mock.Anything).Return(nil)

		// 將Ｍock注入真的Service
		cs := NewClinicService(crm)
		err := cs.Create(jsonObject)
		assert.NoError(test, err)
		crm.AssertExpectations(test)
	})

	test.Run("失敗：Create()，DB出現錯誤。", func(test *testing.T) {
		crm := new(ClinicRepositoryMock)

		errorString := "DB出現錯誤"
		jsonObject := &model.ClinicCreateForm{
			Name:          "診所",
			StartAt:       "2022-10-31",
			EndAt:         "9999-12-31",
			QuotaPerMonth: 100,
		}
		// Mock funcs
		crm.On("Create", jsonObject, mock.Anything).Return(errors.New(errorString))

		// 將Ｍock注入真的Service
		cs := NewClinicService(crm)
		err := cs.Create(jsonObject)
		assert.EqualError(test, err, errorString)
		crm.AssertExpectations(test)
	})

	// Update
	test.Run("成功：Update()，成功從Repository更新資料。", func(test *testing.T) {
		urm := new(ClinicRepositoryMock)

		jsonObject := &model.ClinicUpdateForm{
			ClinicTokenForm: model.ClinicTokenForm{
				ClinicId: 1,
			},
			ClinicCreateForm: model.ClinicCreateForm{
				Name:          "診所",
				StartAt:       "2022-10-31",
				EndAt:         "9999-12-31",
				QuotaPerMonth: 100,
			},
		}
		// Mock funcs
		urm.On("Update", jsonObject).Return(nil)

		// 將Ｍock注入真的Service
		cs := NewClinicService(urm)
		err := cs.Update(jsonObject)
		assert.NoError(test, err)
		urm.AssertExpectations(test)
	})

	test.Run("失敗：Update()，DB出現錯誤。", func(test *testing.T) {
		urm := new(ClinicRepositoryMock)

		errorString := "DB出現錯誤"
		jsonObject := &model.ClinicUpdateForm{
			ClinicTokenForm: model.ClinicTokenForm{
				ClinicId: 1,
			},
			ClinicCreateForm: model.ClinicCreateForm{
				Name:          "診所",
				StartAt:       "2022-10-31",
				EndAt:         "9999-12-31",
				QuotaPerMonth: 100,
			},
		}
		// Mock funcs
		urm.On("Update", jsonObject).Return(errors.New(errorString))

		// 將Ｍock注入真的Service
		cs := NewClinicService(urm)
		err := cs.Update(jsonObject)
		assert.EqualError(test, err, errorString)
		urm.AssertExpectations(test)
	})

	// UpdateToken
	test.Run("成功：UpdateToken()，成功從Repository更新Token。", func(test *testing.T) {
		crm := new(ClinicRepositoryMock)

		jsonObject := &model.ClinicTokenForm{
			ClinicId: 1,
		}
		// Mock funcs
		crm.On("UpdateToken", jsonObject, mock.Anything).Return(nil)

		// 將Ｍock注入真的Service
		cs := NewClinicService(crm)
		err := cs.UpdateToken(jsonObject)
		assert.NoError(test, err)
		crm.AssertExpectations(test)
	})

	test.Run("失敗：UpdateToken()，DB出現錯誤。", func(test *testing.T) {
		crm := new(ClinicRepositoryMock)

		errorString := "DB出現錯誤"
		jsonObject := &model.ClinicTokenForm{
			ClinicId: 1,
		}
		// Mock funcs
		crm.On("UpdateToken", jsonObject, mock.Anything).Return(errors.New(errorString))

		// 將Ｍock注入真的Service
		cs := NewClinicService(crm)
		err := cs.UpdateToken(jsonObject)
		assert.EqualError(test, err, errorString)
		crm.AssertExpectations(test)
	})
}
