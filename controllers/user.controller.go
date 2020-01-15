package controller

import (
	cache "GoAuth/cache"
	common "GoAuth/commons"
	userService "GoAuth/services/user"

	"errors"
	"regexp"

	"github.com/gin-gonic/gin"
)

//PhoneParams 手机登录请求参数
type PhoneParams struct {
	Phone string `form:"phone" json:"phone" xml:"phone" binding:"required"`
	Code  int    `form:"code" json:"code" xml:"code" binding:"required"`
}

//WechatParams 手机登录请求参数
type WechatParams struct {
	Code string `form:"code" json:"code" xml:"code" binding:"required"`
}

type bindWechatParams struct {
	Code string `form:"code" json:"code" xml:"code" binding:"required"`
}

// LoginPhone  用户账号密码登录
func LoginPhone(ctx *gin.Context) {
	var reqInfo PhoneParams
	if err := ctx.ShouldBind(&reqInfo); err != nil {
		common.ThrowError(500, err)
	}
	reg := `^1([38][0-9]|14[57]|5[^4])\d{8}$`
	rgx := regexp.MustCompile(reg)

	if !rgx.MatchString(reqInfo.Phone) {
		common.ThrowError(500, errors.New("手机号错误"))
	}
	var vsf = ctx.Request.Header.Get("vsf")
	var Xtoken = cache.CheckoutToken(ctx.Request.Header.Get("x-token"))
	var token *string = &Xtoken

	//使用手机账号获取用户信息
	_, err := userService.PhoneSignIn(reqInfo.Phone, reqInfo.Code, token, vsf)
	if err != nil {
		common.ThrowError(500, err)
	}
	data := gin.H{
		"token": &token,
	}
	common.JSONData(ctx, data)
}

//LoginWechat 使用微信登录
func LoginWechat(ctx *gin.Context) {
	var reqInfo WechatParams
	if err := ctx.ShouldBind(&reqInfo); err != nil {
		common.ThrowError(500, err)
	}
	vsf := ctx.Request.Header.Get("vsf")
	//使用手机账号获取用户信息
	var Xtoken = cache.CheckoutToken(ctx.Request.Header.Get("x-token"))
	var token *string = &Xtoken
	_, err := userService.WechatSignIn(reqInfo.Code, token, vsf)
	if err != nil {
		common.ThrowError(500, err)
	}
	//查看当前用户状态
	data := gin.H{
		"token": token,
	}
	common.JSONData(ctx, data)
}

//LoginUserName 使用username登录
func LoginUserName(ctx *gin.Context) {

}

// GetSelfInfo  用户账号密码登录
func GetSelfInfo(ctx *gin.Context) {
	UserInfo := cache.UserCacheInfo.GetUserInCtx(ctx)
	if UserInfo == nil {
		common.ThrowError(500, errors.New("获取空值，未获取到用户信息"))
	}
	data := gin.H{
		"user": UserInfo,
	}
	common.JSONData(ctx, data)
}

//UpdateSelf 更新自己信息
func UpdateSelf(ctx *gin.Context) {

}

//CancelSelf 注销自己
func CancelSelf(ctx *gin.Context) {

}

//Logout 退出登录
func Logout(ctx *gin.Context) {

}

//SendLoginCode 发送手机验证码
func SendLoginCode(ctx *gin.Context) {

}

//BindPhone 绑定手机号
func BindPhone(ctx *gin.Context) {
	var reqInfo PhoneParams
	var err error
	if err = ctx.ShouldBind(&reqInfo); err != nil {
		common.ThrowError(500, err)
	}
	reg := `^1([38][0-9]|14[57]|5[^4])\d{8}$`
	rgx := regexp.MustCompile(reg)

	if !rgx.MatchString(reqInfo.Phone) {
		common.ThrowError(500, errors.New("手机号错误"))
	}

	data := gin.H{
		"success": true,
	}
	common.JSONData(ctx, data)
}

//BindWechat 绑定微信
func BindWechat(ctx *gin.Context) {
	var reqInfo WechatParams
	if err := ctx.ShouldBind(&reqInfo); err != nil {
		common.ThrowError(500, err)
	}
	selfUID := cache.UserCacheInfo.GetUserIDToCtx(ctx)
	if selfUID == 0 {
		common.ThrowError(500, errors.New("无用户信息"))
	}
	vsf := ctx.Request.Header.Get("vsf")
	token := ctx.Request.Header.Get("x-token")
	err := userService.BindWechat(selfUID, reqInfo.Code, &token, vsf)
	if err != nil {
		common.ThrowError(500, err)
	}
	data := gin.H{
		"success": true,
	}
	common.JSONData(ctx, data)
}

//BindUserName 绑定Username
func BindUserName(ctx *gin.Context) {

}

//UnbindAccount 解绑账号
func UnbindAccount(ctx *gin.Context) {

}

//WechatBindPhone 微信缓存绑定手机
func WechatBindPhone(ctx *gin.Context) {
	var reqInfo PhoneParams
	var err error
	if err = ctx.ShouldBind(&reqInfo); err != nil {
		common.ThrowError(500, err)
	}
	reg := `^1([38][0-9]|14[57]|5[^4])\d{8}$`
	rgx := regexp.MustCompile(reg)

	if !rgx.MatchString(reqInfo.Phone) {
		common.ThrowError(500, errors.New("手机号错误"))
	}
	var vsf = ctx.Request.Header.Get("vsf")
	token := ctx.Request.Header.Get("x-token")
	if err = userService.WechatBindPhone(reqInfo.Phone, reqInfo.Code, &token, vsf); err != nil {
		common.ThrowError(500, errors.New("获取数据失败"))
	}
	data := gin.H{
		"success": true,
	}
	common.JSONData(ctx, data)
}
