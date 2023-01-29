package middleware

import (
	"go-pano/utils"
	"strings"

	"github.com/gin-gonic/gin"
)

// JWTAuthMiddleware 基于JWT的驗證的中間件
func JWTAuthMiddleware() func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		// 客户端携带Token有三种方式 1.放在请求头 2.放在请求体 3.放在URI
		// 这里假设Token放在Header的Authorization中，并使用Bearer开头
		// 这里的具体实现方式要依据你的实际业务情况决定
		authHeader := ctx.Request.Header.Get("Authorization")
		if authHeader == "" {
			ctx.JSON(400, utils.H500("Header不存在Authorization"))
			ctx.Abort()
			return
		}
		// 按空格分割
		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			ctx.JSON(400, utils.H500("Authorization請記得使用Bearer當作開頭"))
			ctx.Abort()
			return
		}

		// parts[1]是獲取到的tokenString
		claims, err := utils.ParseToken(parts[1])
		if err != nil {
			ctx.JSON(400, utils.H500("無效的Token"))
			ctx.Abort()
			return
		}
		// 將claims資訊傳遞至context中，方便獲取資訊
		ctx.Set("jwtClaims", claims)

		ctx.Next()
	}
}
