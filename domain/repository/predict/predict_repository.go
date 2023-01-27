package predict_repository

import (
	"fmt"
	"go-pano/domain/model"
	"go-pano/utils"
)

// 接口
type IPredictRepository interface {
	GetFirstByIDAndFileName(*model.Predict) (*model.Predict, error)
	Create(*model.Predict) (*model.Predict, error)
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

func (ctrl *PredictRepository) GetFirstByIDAndFileName(predict *model.Predict) (*model.Predict, error) {
	if err := ctrl.mysql.DB().Where("clinic_id=?", predict.ClinicId).Where("filename=?", predict.Filename).First(predict).Error; err != nil {
		return &model.Predict{}, err
	}

	return predict, nil
}

// Insert a Predict
func (ctrl *PredictRepository) Create(predict *model.Predict) (*model.Predict, error) {
	db := ctrl.mysql.DB().Create(predict)
	// 處理不同ERR種類
	if db.Error != nil {
		return &model.Predict{}, db.Error
	} else if db.RowsAffected == 0 {
		return &model.Predict{}, utils.ErrKeyConflict
	}
	fmt.Println(predict.PredictId)

	return predict, nil
}
