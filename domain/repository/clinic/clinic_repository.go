package clinic_repository

import (
	"go-pano/domain/model"
	"go-pano/utils"
)

// 接口
type IClinicRepository interface {
	GetAll() ([]model.Clinic, error)
	Update(*model.ClinicUpdateForm) error
	Create(*model.ClinicCreateForm, string) error
	UpdateToken(*model.ClinicTokenForm, string) error
}

// gorm使用參考
// https://pjchender.dev/golang/note-gorm-example/
// new Class()
func NewClinicRepository(db utils.IDBInstance) IClinicRepository {
	return &ClinicRepository{mysql: db}
}

// class
type ClinicRepository struct {
	mysql utils.IDBInstance
}

// 獲取全部Clinic
func (ctrl *ClinicRepository) GetAll() ([]model.Clinic, error) {
	clinic := []model.Clinic{}
	if err := ctrl.mysql.DB().Table("Clinic").
		Select("clinic_id", "SUBSTR(start_at,1,10) AS start_at", "SUBSTR(end_at,1,10) AS end_at", "name", "quota_per_month", "token").
		Limit(500).Find(&clinic).Error; err != nil {
		return []model.Clinic{}, err
	}

	return clinic, nil
}

// 插入Clinic
func (ctrl *ClinicRepository) Create(clinic *model.ClinicCreateForm, token string) error {
	if db := ctrl.mysql.DB().Omit("clinic_id").Create(
		&model.Clinic{
			ClinicCreateForm: *clinic,
			Token:            token,
		},
	); db.Error != nil {
		return db.Error
	}

	return nil
}

// 更新Clinic
func (ctrl *ClinicRepository) Update(clinic *model.ClinicUpdateForm) error {
	if db := ctrl.mysql.DB().Table("Clinic").Omit("clinic_id").Updates(clinic); db.Error != nil {
		return db.Error
	} else if db.RowsAffected == 0 {
		return utils.ErrFailed
	}

	return nil
}

// 更新Token
func (ctrl *ClinicRepository) UpdateToken(clinic *model.ClinicTokenForm, token string) error {
	if db := ctrl.mysql.DB().Table("Clinic").Where("clinic_id = ?", clinic.ClinicId).Update("token", token); db.Error != nil {
		return db.Error
	} else if db.RowsAffected == 0 {
		return utils.ErrFailed
	}

	return nil
}
