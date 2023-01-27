package utils

import (
	"errors"
	"go-pano/config"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
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

// JWTAuthMiddleware 基于JWT的驗證的中間件
func JWTAuthMiddleware() func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		// 客户端携带Token有三种方式 1.放在请求头 2.放在请求体 3.放在URI
		// 这里假设Token放在Header的Authorization中，并使用Bearer开头
		// 这里的具体实现方式要依据你的实际业务情况决定
		authHeader := ctx.Request.Header.Get("Authorization")
		if authHeader == "" {
			ctx.JSON(400, H500("Header不存在Authorization"))
			ctx.Abort()
			return
		}
		// 按空格分割
		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			ctx.JSON(400, H500("Authorization請記得使用Bearer當作開頭"))
			ctx.Abort()
			return
		}

		// parts[1]是獲取到的tokenString
		claims, err := ParseToken(parts[1])
		if err != nil {
			ctx.JSON(400, H500("無效的Token"))
			ctx.Abort()
			return
		}
		// 將claims資訊傳遞至context中，方便獲取資訊
		ctx.Set("jwtUserId", claims.UserId)
		ctx.Set("jwtName", claims.Name)
		ctx.Set("jwtRoles", claims.Roles)
		ctx.Next()
	}
}
