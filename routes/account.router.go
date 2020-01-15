package router

import (
	controller "GoAuth/controllers"
	middleware "GoAuth/middlewares"

	"github.com/gin-gonic/gin"
)

//InitAccountRouter 初始化账号相关路由
func InitAccountRouter(_router *gin.Engine) {
	router := _router.Group("account")
	{
		router.GET("/login_code", controller.SendLoginCode)      //发送手机登录验证码
		router.POST("/login/phone", controller.LoginPhone)       //手机登录
		router.POST("/login/wechat", controller.LoginWechat)     //微信登录
		router.POST("/login/username", controller.LoginUserName) //账号密码
		router.GET("/logout", controller.Logout)                 //退出登录

		router.POST("/wechat/bind/phone", controller.WechatBindPhone)                     //绑定手机
		router.POST("/bind/phone", middleware.UserAuth(true), controller.BindPhone)       //绑定手机
		router.POST("/bind/wechat", middleware.UserAuth(true), controller.BindWechat)     //绑定微信
		router.POST("/bind/username", middleware.UserAuth(true), controller.BindUserName) //绑定账号密码
		router.GET("/unbind/:type", middleware.UserAuth(true), controller.UnbindAccount)  //解绑账号
	}
}
