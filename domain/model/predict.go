package model

import "time"

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

// python的Response格式
type Result struct {
	IsSuccessful bool        `json:"isSuccessful"`
	Msg          string      `json:"msg"`
	Predict      interface{} `json:"predict"`
}
