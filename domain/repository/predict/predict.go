package repository

import (
	"errors"
	"go-pano/domain/model"
	"go-pano/utils"

	"gorm.io/gorm"
)

// 接口
type IPredictRepository interface {
	GetFirstByIDAndFileName(int, string) (*model.Predict, error)
	Create(int, string, string) error
}

// gorm使用參考
// https://pjchender.dev/golang/note-gorm-example/

// new Class()
func NewPredictRepository(db utils.IDBInstance) IPredictRepository {
	return &PredictRepository{mysql: db}
}

// class
type PredictRepository struct {
	mysql utils.IDBInstance
}

// 有(p *PredictRepository)都是PredictRepository的public函數
// GetFirstByID gets the user by his ID
func (ctrl *PredictRepository) GetFirstByIDAndFileName(clinicId int, fileName string) (*model.Predict, error) {
	predictModel := &model.Predict{}
	err := ctrl.mysql.DB().Where("clinic_id=?", clinicId).Where("filename=?", fileName).First(predictModel).Error

	// First找不到會報錯
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return predictModel, nil
	}

	// 抱錯
	if err != nil {
		return nil, err
	}

	return predictModel, nil
}

// Insert a Predict
func (ctrl *PredictRepository) Create(clinic_id int, fileName string, predict string) error {
	predictModel := &model.Predict{
		ClinicId:      clinic_id,
		Filename:      fileName,
		PredictString: predict,
	}

	db := ctrl.mysql.DB().Create(predictModel)

	// 處理不同ERR種類
	if db.Error != nil {
		return db.Error
	} else if db.RowsAffected == 0 {
		return utils.ErrKeyConflict
	}

	return nil
}
