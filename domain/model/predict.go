package model

// Predict model
type Predict struct {
	PredictId int    `gorm:"column:predict_id;primaryKey;not null;autoIncrement;comment:辨識ID" json:"predict_id"`
	ClinicId  int    `gorm:"column:clinic_id;not null;index:idx_predict_clinic_id;comment:診所ID" json:"clinic_id"`
	Dir       string `gorm:"column:dir;size:64;not null;comment:檔案目錄" json:"dir"`
	Filename  string `gorm:"column:filename;size:128;not null;index:idx_predict_filename;comment:檔案名稱" json:"filename"`
	// JSON不要顯示，但是要存入DB
	PredictString string `gorm:"column:predict_string;size:1024;not null;comment:辨識結果" json:"-"`
	// JSON要顯示，但是不要存入DB
	Predict   []PredictDetail `gorm:"-:all" json:"predict"`
	CreatedAt string          `gorm:"column:created_at;size:20;not null;index:idx_predict_created_at;comment:建立日期;autoCreateTime" json:"created_at"`
	UpdatedAt string          `gorm:"column:updated_at;size:20;not null;index:idx_predict_updated_at;comment:更新日期;autoUpdateTime" json:"updated_at"`
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
	ClinicId int `form:"clinic_id" json:"clinic_id" binding:"required,max=11"`
}
