package repository

import (
	"errors"
	db "go-pano/utils"
	error_ "go-pano/utils"
	"time"

	"gorm.io/gorm"
)

// Predict model
type Predict struct {
	ID        int       `gorm:"primary_key" json:"id"`
	ClinicId  int       `json:"clinicId"`
	Filename  string    `json:"fileName"`
	Predict   string    `json:"predict"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// TableName for gorm
func (Predict) TableName() string {
	return "Predict"
}

// GetFirstByID gets the user by his ID
func (p *Predict) GetFirstByIDAndFileName(clinicId int, fileName string) error {
	err := db.DB().Where("clinic_id=?", clinicId).Where("filename=?", fileName).First(p).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return error_.ErrDataNotFound
	}

	return err
}

// Insert a Predict
func (p *Predict) Create() error {
	db := db.DB().Create(p)

	if db.Error != nil {
		return db.Error
	} else if db.RowsAffected == 0 {
		return error_.ErrKeyConflict
	}

	return nil
}
