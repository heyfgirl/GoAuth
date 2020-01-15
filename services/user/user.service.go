package user

import (
	"GoAuth/cache"
	common "GoAuth/commons"
	config "GoAuth/configs"
	model "GoAuth/models"
	accountService "GoAuth/services/account"
	smsService "GoAuth/services/sms"
	"encoding/json"

	"errors"

	"github.com/jinzhu/gorm"
	"github.com/jinzhu/gorm/dialects/postgres"
)

//PhoneSignIn 手机登录
func PhoneSignIn(phone string, code int, token *string, vsf string) (*model.User, error) {
	//检验code
	if !smsService.CheckCode(phone, code) {
		return nil, errors.New("验证码验证失败")
	}
	//获取用户信息
	userInfo := &(model.User{})
	err := model.DB.Where("accounts.token = ? and accounts.type = ?", phone, "phone").
		Joins(`left join accounts on accounts.user_id = users.id`).
		Select([]string{"users.id", "users.hash", "users.created_at", "users.nick_name", "users.real_name", "users.sex", "users.avatar", "users.state", "users.state"}).
		Preload("Account").
		First(userInfo).Error
	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			//未获取到用户无该用户，此时应该为该手机号注册用户
			userInfo, err = registerPhoneUser(phone, nil)
			if err != nil {
				return nil, err
			}
		} else {
			//其他错误
			return nil, err
		}
	}
	accountService.SetTokenCache(token, userInfo.Hash, vsf)
	return userInfo, nil
}

//registerPhoneUser 注册创建手机用户
func registerPhoneUser(phone string, user *model.User) (*model.User, error) {
	var createUserInfo = model.User{}
	if user == nil {
		//用户名为手机用户加手机尾号后4位
		//账号数据
		accountInfo := model.Account{
			DefaultModel: model.DefaultModel{
				Hash: common.GetRandomStr(32),
			},
			Type:  "phone",
			Token: phone,
		}
		//用户数据
		var sex int = 0
		nickname := config.Conf.DefaultUser.NickNmaePrefix + phone[7:len(phone)]
		if user != nil {
			if user.Sex != 0 {
				sex = user.Sex
			}
			if user.NickName != "" {
				nickname = user.NickName
			}
		}
		createUserInfo = model.User{
			DefaultModel: model.DefaultModel{
				Hash: common.GetRandomStr(32),
			},
			NickName: nickname,
			Sex:      sex,
			Avatar:   config.Conf.DefaultUser.Avatar,
			Account:  []model.Account{accountInfo},
		}
	} else {
		createUserInfo = *user
	}
	if err := model.DB.Create(&createUserInfo).Error; err != nil {
		common.ThrowError(500, err)
	}
	return &createUserInfo, nil
}

//SendPhoneCode 发送手机登录验证码短信
func SendPhoneCode(phone string) {

}

//WechatSignIn 使用微信登录
func WechatSignIn(code string, token *string, vsf string) (*model.User, error) {
	accountInfo, err := UpsertWechatAccount(code, vsf)
	if err != nil {
		return nil, err
	}
	//获取完账号 ，获取改账号的用户信息
	if accountInfo.UserID == 0 {
		//新账号,无用户 [只创建了微信账号，未绑定用户]【给前端返回让其进行绑定操作】
		// 设置微信的账号信息缓存。
		if err = cache.SetDataToCache(config.Conf.RedisKey.WechatAccount+*token, map[string]interface{}{
			"token": accountInfo.Token,
		}); err != nil {
			common.ThrowError(500, err)
		}
		common.ThrowError(2500, errors.New("缺少必须账号手机号"))
	}
	userInfo := &model.User{}
	err = model.DB.Where("users.id = ?", accountInfo.UserID).
		Joins(`left join accounts on accounts.user_id = users.id`).
		Select([]string{"users.id", "users.hash", "users.created_at", "users.nick_name", "users.real_name", "users.sex", "users.avatar", "users.state", "users.state"}).
		Preload("Account").
		First(userInfo).Error
	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			//该账号下的用户已经获取不到，说明该用户已经处于删除状态 或异常问题
			//【特殊日志处理】
			common.ThrowError(500, errors.New("异常错误，获取不到该账号的用户信息"))
		} else {
			//其他错误
			return nil, err
		}
	}
	accountService.SetTokenCache(token, userInfo.Hash, vsf)
	return userInfo, nil
}

//BindWechat 绑定微信
func BindWechat(userID int, code string, token *string, vsf string) error {
	wechatInfo, err := accountService.GetWechatInfo(code, vsf)
	if err != nil {
		common.ThrowError(500, errors.New("异常错误"))
	}
	//获取过微信信息的账号必然存在  但其userid为0
	accountInfo := &model.Account{}
	unionid := (*wechatInfo)["unionid"].(string)
	if err := model.DB.Where("accounts.token = ? and accounts.type = ?", unionid, "wechat").
		First(accountInfo).Error; err != nil {
		common.ThrowError(500, errors.New("异常错误"))
	}
	if accountInfo.UserID != 0 {
		common.ThrowError(500, errors.New("该微信账号已经手机用户"))
	}
	if err := model.DB.Model(&accountInfo).Update("user_id", userID).
		Error; err != nil {
		common.ThrowError(500, errors.New("异常错误"))
	}
	// UpdateUserCache(token, accountInfo.UserID)
	return nil
}

//BindPhone 绑定手机账号
func BindPhone(userID int, phone string, code int) error {

	return nil
}

//UpsertWechatAccount 更新或者创建 微信账号【根据code获取微信账号，如果不存在则创建该账号】
func UpsertWechatAccount(code string, vsf string) (*model.Account, error) {
	wechatInfo, err := accountService.GetWechatInfo(code, vsf)
	if err != nil {
		return nil, err
	}
	//查看数据库是否存在该账号信息如果没有则创建该账号
	unionid := (*wechatInfo)["unionid"].(string)
	accountInfo := &model.Account{}
	err = model.DB.Where("accounts.token = ? and accounts.type = ?", unionid, "wechat").
		First(accountInfo).Error
	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			//未获取到用户无该微信账号，注册成为账号但不注册用户
			extension, err := json.Marshal(*wechatInfo)
			if err != nil {
				return nil, err
			}
			accountInfo := model.Account{
				DefaultModel: model.DefaultModel{
					Hash: common.GetRandomStr(32),
				},
				Type:  "wechat",
				Token: unionid,
				Extension: postgres.Jsonb{
					RawMessage: extension,
				},
			}
			if err := model.DB.Create(&accountInfo).Error; err != nil {
				return nil, err
			}
		}
	}
	return accountInfo, nil
}

//WechatBindPhone 微信绑定手机
func WechatBindPhone(phone string, code int, token *string, vsf string) error {
	//使用需要绑定的手机进行登录操作
	tokenF := *token
	userInfo, err := PhoneSignIn(phone, code, token, vsf)
	if tokenF != *token {
		common.ThrowError(500, errors.New("异常错误，该token已被占用不可重新用于用户设置"))
	}
	//使用手机账号获取用户信息
	unionid := accountService.GetCacheWechatAccountUnionid(*token)
	if unionid != "" {
		//获取到用户在缓存中的微信信息则进行绑定操作
		if err = accountService.BindWechat(userInfo.ID, unionid); err != nil {
			common.ThrowError(500, err)
		}
		//将用户缓存重新设置
		// UpdateUserCache(token, userInfo.ID)
	} else {
		common.ThrowError(500, errors.New("绑定失败"))
	}
	return nil
}
