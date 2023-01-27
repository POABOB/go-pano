package clinic_service

import (
	"go-pano/domain/model"
	repository "go-pano/domain/repository/clinic"
	"go-pano/utils"
)

// interface
type IClinicService interface {
	Get() ([]model.Clinic, error)
	Create(*model.ClinicCreateForm) error
	Update(*model.ClinicUpdateForm) error
	UpdateToken(*model.ClinicTokenForm) error
}

// 實例化
func NewClinicService(r repository.IClinicRepository) IClinicService {
	return &ClinicService{ClinicRepository: r}
}

// class
type ClinicService struct {
	ClinicRepository repository.IClinicRepository
}

func (cs *ClinicService) Get() ([]model.Clinic, error) {
	result, err := cs.ClinicRepository.GetAll()
	if err != nil {
		return []model.Clinic{}, err
	}

	return result, nil
}

func (cs *ClinicService) Create(clinic *model.ClinicCreateForm) error {
	// 插入
	if err := cs.ClinicRepository.Create(clinic, utils.RandStringRunes(50)); err != nil {
		return err
	}

	return nil
}

func (cs *ClinicService) Update(clinic *model.ClinicUpdateForm) error {
	if err := cs.ClinicRepository.Update(clinic); err != nil {
		return err
	}

	return nil
}

func (cs *ClinicService) UpdateToken(clinic *model.ClinicTokenForm) error {
	// 自動產生TOKEN
	if err := cs.ClinicRepository.UpdateToken(clinic, utils.RandStringRunes(50)); err != nil {
		return err
	}

	return nil
}
