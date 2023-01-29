package model

// 表單驗證 https://pkg.go.dev/github.com/go-playground/validator/v10#section-readme
// Clinic model
type Clinic struct {
	ClinicTokenForm  `gorm:"embedded"`
	ClinicCreateForm `gorm:"embedded"`
	Predict          []Predict
	Token            string `gorm:"<-:create;column:token;size:100;not null;index:idx_clinic_token;comment:Token" json:"token" binding:"required,max=1024" example:"token1231324568913"`
}

type ClinicTokenForm struct {
	ClinicId int `gorm:"<-:update;column:clinic_id;primaryKey;not null;autoIncrement;comment:診所ID" json:"clinic_id" binding:"required,max=11" example:"1"`
}

type ClinicCreateForm struct {
	Name          string `gorm:"column:name;size:128;not null;comment:診所名稱" json:"name" binding:"required,max=128" example:"診所"`
	StartAt       string `gorm:"column:start_at;size:10;not null;index:idx_clinic_start_at;comment:開始日期" json:"start_at" binding:"required,len=10" example:"2022-10-10"`
	EndAt         string `gorm:"column:end_at;size:10;not null;index:idx_clinic_end_at;comment:結束日期" json:"end_at" binding:"required,len=10" example:"2099-12-31"`
	QuotaPerMonth int    `gorm:"column:quota_per_month;type:int(11);not null;comment:額度" json:"quota_per_month" binding:"required,max=100000" example:"200"`
}

type ClinicUpdateForm struct {
	ClinicTokenForm
	ClinicCreateForm
}

// TableName for gorm
func (Clinic) TableName() string {
	return "Clinic"
}
