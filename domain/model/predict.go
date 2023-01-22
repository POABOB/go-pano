package model

import (
	"time"
)

// Predict model
type Predict struct {
	ID       int    `gorm:"primary_key" json:"id"`
	ClinicId int    `json:"clinic_id"`
	Filename string `json:"filename"`
	// JSON不要顯示，但是要存入DB
	PredictString string `json:"-" gorm:"predict"`
	// JSON要顯示，但是不要存入DB
	Predict   []PredictDetail `json:"predict" gorm:"-:all"`
	CreatedAt time.Time       `json:"created_at"`
	UpdatedAt time.Time       `json:"updated_at"`
}

type PredictDetail struct {
	Text     string `json:"text"`
	FileName string `json:"filename"`
	Data     struct {
		Caries   []int `json:"caries,omitempty"`
		Filling  []int `json:"filling,omitempty"`
		Endo     []int `json:"endo,omitempty"`
		Crown    []int `json:"crown,omitempty"`
		Bridge   []int `json:"bridge,omitempty"`
		Imp      []int `json:"imp,omitempty"`
		Post     []int `json:"post,omitempty"`
		Embedded []int `json:"embedded,omitempty"`
		Impacted []int `json:"impacted,omitempty"`
		R_R      []int `json:"r_r,omitempty"`
	}
}

// TableName for gorm
func (Predict) TableName() string {
	return "Predict"
}

// python的Response格式
type Result struct {
	IsSuccessful bool   `json:"isSuccessful"`
	Msg          string `json:"msg"`
	Predict      string `json:"predict"`
}

// 表單驗證 https://pkg.go.dev/github.com/go-playground/validator/v10#section-readme
type PredictForm struct {
	ClinicId int    `form:"clinic_id" json:"clinic_id" binding:"required,max=11"`
	Method   string `uri:"method" binding:"max=4"`
}
