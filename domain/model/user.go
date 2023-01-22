package model

// 表單驗證 https://pkg.go.dev/github.com/go-playground/validator/v10#section-readme
// User model
type User struct {
	UserId      int      `gorm:"primary_key,<-:update" json:"user_id,omitempty" binding:"max=11"`
	Name        string   `json:"name" binding:"required,max=64"`
	Account     string   `json:"account" binding:"required,max=64"`
	RolesString string   `json:"-" gorm:"roles"`
	Roles       []string `json:"roles" gorm:"-:all" binding:"required,max=1024"`
	Password    string   `gorm:"<-" json:"password,omitempty" binding:"required,max=1024,eqfield=Passconf"` // 只能創建和更新
	Passconf    string   `gorm:"-:all" json:"passconf,omitempty" binding:"required"`
	Status      int      `json:"status" binding:"required,max=11"`
}

// 更新User表單
type UserUpdateForm struct {
	UserId      int      `gorm:"primary_key" json:"user_id" binding:"required,max=11"`
	Name        string   `json:"name" binding:"required,max=64"`
	Account     string   `json:"account" binding:"required,max=64"`
	RolesString string   `json:"-" gorm:"roles"`
	Roles       []string `json:"roles" gorm:"-:all" binding:"required,max=1024"`
}

// 更新Password表單
type UserPasswordForm struct {
	UserId   int    `gorm:"primary_key" json:"user_id" binding:"required,max=11"`
	Password string `gorm:"<-" json:"password,omitempty" binding:"required,max=1024,eqfield=Passconf"` // 只能創建和更新
	Passconf string `gorm:"-:all" json:"passconf,omitempty" binding:"required"`                        // 不對DB做任何事
}

// 更新狀態表單
type UserStatusForm struct {
	UserId int `gorm:"primary_key" json:"user_id" binding:"required,max=11"`
	Status int `json:"status" binding:"max=11"`
}

// TableName for gorm
func (User) TableName() string {
	return "Users"
}
