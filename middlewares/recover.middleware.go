package middleware

import (
	common "GoAuth/commons"

	"github.com/gin-gonic/gin"
)

//ErrorInfo 错误详情【code和详情】
// type ErrorInfo map[string]interface{}

//Recover 统一处理中间件
func Recover() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		token := ctx.GetHeader("x-token")
		ctx.Header("x-token", token)
		ctx.Set("token", token)
		// ctx.Writer.Header()["id"] = []string{"abc"} go会自动将header中的key格式化首字母大小写，由于 header是一个map所以可以越过set add等方法
		//错误捕获以及处理
		defer common.RecoverError(ctx)
		ctx.Next()
	}
}
