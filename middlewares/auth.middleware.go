package middleware

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

//UserAuth 用户验证中间件
func UserAuth(have bool) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		token := ctx.Request.Header.Get("x-token")
		fmt.Println(token)
		// UserCacheInfo, err := cache.UserCacheInfo.GetUserInCache(token)
		// if err != nil && have {
		// 	common.ThrowError(500, err)
		// }
		// //token已被占用 创建新的token
		// if UserCacheInfo == nil {
		// 	if have {
		// 		common.ThrowError(500, errors.New("异常错误，该用户未登录"))
		// 	}
		// }
		// //将当前用户信息塞入ctx上下文中
		// UserCacheInfo.SetUserToCtx(ctx)
		ctx.Next()
	}
}
