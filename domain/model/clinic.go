package model

// 表單驗證 https://pkg.go.dev/github.com/go-playground/validator/v10#section-readme
// Clinic model
type Clinic struct {
	ClinicTokenForm
	ClinicCreateForm
	Token string `json:"token" gorm:"<-:create" binding:"required,max=1024" example:"token1231324568913"`
}

type ClinicTokenForm struct {
	ClinicId int `gorm:"primaryKey;<-:update" json:"clinic_id" binding:"required,max=11" example:"1"`
}

type ClinicCreateForm struct {
	Name          string `json:"name" binding:"required,max=128" example:"診所"`
	StartAt       string `json:"start_at" binding:"required,len=10" example:"2022-10-10"`
	EndAt         string `json:"end_at" binding:"required,len=10" example:"2099-12-31"`
	QuotaPerMonth int    `json:"quota_per_month" binding:"required,max=100000" example:"200"`
}

type ClinicUpdateForm struct {
	ClinicTokenForm
	ClinicCreateForm
}

// TableName for gorm
func (Clinic) TableName() string {
	return "Clinic"
}
