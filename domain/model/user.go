package model

// 表單驗證 https://pkg.go.dev/github.com/go-playground/validator/v10#section-readme
// User model
type User struct {
	UserUpdateForm
}

// 新增User表單
type UserCreateForm struct {
	Name        string   `json:"name" binding:"required,max=64" example:"User"`
	Account     string   `json:"account" binding:"required,max=64" example:"user"`
	RolesString string   `json:"-" swaggerignore:"true"`
	Roles       []string `json:"roles" gorm:"-:all" binding:"required,max=1024" example:"admin"`
	Password    string   `gorm:"<-" json:"password" binding:"required,max=1024,eqfield=Passconf" example:"password"` // 只能創建和更新
	Passconf    string   `gorm:"-:all" json:"passconf" binding:"required" example:"password"`
	Status      int      `json:"status" binding:"max=11" example:"1"`
}

// 更新User表單
type UserUpdateForm struct {
	UserId      int      `gorm:"primaryKey" json:"user_id" binding:"required,max=11" example:"1"`
	Name        string   `json:"name" binding:"required,max=64" example:"User"`
	Account     string   `json:"account" binding:"required,max=64" example:"user"`
	RolesString string   `json:"-" swaggerignore:"true"`
	Roles       []string `json:"roles" gorm:"-:all" binding:"required,max=1024" example:"admin"`
	Status      int      `json:"status" binding:"max=11" example:"1"`
}

// 更新Password表單
type UserPasswordForm struct {
	UserId   int    `gorm:"primaryKey" json:"user_id" binding:"required,max=11" example:"1"`
	Password string `gorm:"<-" json:"password" binding:"required,max=1024,eqfield=Passconf" example:"password"` // 只能創建和更新
	Passconf string `gorm:"-:all" json:"passconf" binding:"required" example:"password"`                        // 不對DB做任何事
}

// TableName for gorm
func (User) TableName() string {
	return "Users"
}
