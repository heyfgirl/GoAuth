package router

import (
	controller "GoAuth/controllers"
	middleware "GoAuth/middlewares"

	"github.com/gin-gonic/gin"
)

//InitUserRouter 初始化账号相关路由
func InitUserRouter(_router *gin.Engine) {
	router := _router.Group("user")
	{
		router.GET("/self/info", middleware.UserAuth(true), controller.GetSelfInfo)   //获取用户自己信息
		router.POST("/self/update", middleware.UserAuth(true), controller.UpdateSelf) //更新用户自己信息
		router.GET("/self/cancel", middleware.UserAuth(true), controller.CancelSelf)  //注销用户自己
	}
}
