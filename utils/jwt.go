package utils

import (
	"errors"
	"go-pano/config"
	"time"

	"github.com/dgrijalva/jwt-go"
)

// jwt.StandardClaims會有預設的欄位
// 可以自行增加欄位，來夾帶在JWT之中
type MyClaims struct {
	UserId int      `json:"user_id"`
	Name   string   `json:"name"`
	Roles  []string `json:"roles"`
	jwt.StandardClaims
}

// ７天的使用期限
const TokenExpireDuration = time.Hour * 24 * 7

// 加密的私鑰
var MySecret = []byte(config.Server.Secret)

// GenToken 生成JWT
func GenToken(userId int, name string, roles []string) (string, error) {
	// 指定簽名方法和夾帶的的Claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, MyClaims{
		userId,
		name,
		roles,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(TokenExpireDuration).Unix(), // 過期日
			Issuer:    "gin-go-server",                            // 簽發人
		},
	})
	// 加密後返回token
	return token.SignedString(MySecret)
}

// ParseToken 解析JWT
func ParseToken(tokenString string) (*MyClaims, error) {
	// 解析token
	token, err := jwt.ParseWithClaims(tokenString, &MyClaims{}, func(token *jwt.Token) (interface{}, error) {
		return MySecret, nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(*MyClaims); ok && token.Valid { // 校驗token
		return claims, nil
	}
	return nil, errors.New("invalid token")
}
