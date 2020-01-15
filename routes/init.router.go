package router

import (
	"github.com/gin-gonic/gin"
)

//InitRouter 初始化路由
func InitRouter(router *gin.Engine) {
	//加载账号路由
	InitAccountRouter(router)
	InitUserRouter(router)
}
