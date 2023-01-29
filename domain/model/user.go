package model

// 表單驗證 https://pkg.go.dev/github.com/go-playground/validator/v10#section-readme
// User model
type User struct {
	UserUpdateForm `gorm:"embedded"`
	Password       string `gorm:"<-;column:password;size:1024;index:idx_user_password;comment:密碼" json:"-" swaggerignore:"true"`
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
	UserId      int      `gorm:"column:user_id;primaryKey;not null;autoIncrement;comment:使用者ID" json:"user_id" binding:"required,max=11" example:"1"`
	Name        string   `gorm:"column:name;size:64;not null;comment:名稱" json:"name" binding:"required,max=64" example:"User"`
	Account     string   `gorm:"column:account;size:64;not null;unique;comment:帳號" json:"account" binding:"required,max=64" example:"user"`
	RolesString string   `gorm:"column:roles_string;size:512;not null;default:[];comment:權限" json:"-" swaggerignore:"true"`
	Roles       []string `gorm:"-:all" json:"roles" binding:"required,max=1024" example:"admin"`
	Status      int      `gorm:"column:status;type:tinyint;not null;default:1;comment:狀態" json:"status" binding:"max=11" example:"1"`
}

// 更新Password表單
type UserPasswordForm struct {
	UserId   int    `gorm:"primaryKey" json:"user_id" binding:"required,max=11" example:"1"`
	Password string `gorm:"<-" json:"password" binding:"required,max=1024,eqfield=Passconf" example:"password"` // 只能創建和更新
	Passconf string `gorm:"-:all" json:"passconf" binding:"required" example:"password"`                        // 不對DB做任何事
}

// 登入表單
type UserLoginForm struct {
	Account  string `json:"account" binding:"required,max=64" example:"user"`
	Password string `gorm:"<-:false" json:"password" binding:"required,max=1024" example:"password"`
}

type UserToken struct {
	Token string `json:"token"`
}

// TableName for gorm
func (User) TableName() string {
	return "Users"
}
