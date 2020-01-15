package account

import (
	"GoAuth/cache"
	common "GoAuth/commons"
	config "GoAuth/configs"
	model "GoAuth/models"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"

	"github.com/jinzhu/gorm"
)

//GetWechatInfo 登录
func GetWechatInfo(code string, vsf string) (*map[string]interface{}, error) {
	var AppID string
	var AppSecret string
	if vsf == "app" {
		AppID = config.Conf.Wechat.AppID
		AppSecret = config.Conf.Wechat.AppSecret
	} else {
		AppID = config.Conf.Wechat.WebID
		AppSecret = config.Conf.Wechat.WebSecret
	}
	response, err := http.Get(`https://api.weixin.qq.com/sns/oauth2/access_token?appid=` + AppID + `&secret=` + AppSecret + `&code=` + code + `&grant_type=authorization_code`)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	body, _ := ioutil.ReadAll(response.Body)
	var JSONBody map[string]interface{}
	//注意：反序列化map,不需要make，因为make操作被封装到Unmarshal函数
	err = json.Unmarshal([]byte(body), &JSONBody)
	if JSONBody["errcode"] != nil {
		return nil, errors.New(JSONBody["errmsg"].(string))
	}
	var accessToken string
	var openID string
	if JSONBody["access_token"] != nil {
		accessToken = JSONBody["access_token"].(string)
		openID = JSONBody["openid"].(string)
	} else {
		return nil, errors.New("未获取到数据")
	}
	response, err = http.Get(`https://api.weixin.qq.com/sns/userinfo?access_token=` + accessToken + `&openid=` + openID)
	if err != nil {
		return nil, err
	}
	body, _ = ioutil.ReadAll(response.Body)
	//注意：反序列化map,不需要make，因为make操作被封装到Unmarshal函数
	err = json.Unmarshal([]byte(body), &JSONBody)
	if JSONBody["errcode"] != nil {
		return nil, errors.New(JSONBody["errmsg"].(string))
	}
	if JSONBody["sex"] != nil {
		JSONBody["sex"] = int(JSONBody["sex"].(float64))
	} else {
		var sex int = 0
		JSONBody["sex"] = sex
	}
	// fmt.Println(JSONBody)
	if err != nil {
		return nil, err
	}
	return &JSONBody, nil
}

//BindWechat 绑定微信
func BindWechat(userID int, unionid string) error {
	accountInfo := &model.Account{}
	//获取过微信信息的账号必然存在  但其userid为0
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
	return nil
}

//GetCacheWechatAccountUnionid 获取缓存中的Unionid
func GetCacheWechatAccountUnionid(token string) string {
	//返回为nil 说明是未获取到微信缓存信息
	wechatTokenInfo, _ := cache.GetAndDelDataInCache(config.Conf.RedisKey.WechatAccount + token)
	unionid := ""
	if wechatTokenInfo != nil {
		unionid = wechatTokenInfo["token"].(string)
	} else {
		return ""
	}
	accountInfo := &(model.Account{})
	if err := model.DB.Where("accounts.token = ? and accounts.type = ?", unionid, "wechat").
		First(accountInfo).Error; err != nil {
		//获取账号失败
		common.ThrowError(500, errors.New("错误"))
	}
	if accountInfo.UserID != 0 {
		//微信已被绑定则不再操作
		return ""
	}
	return unionid
}

//SetTokenCache 登陆成功设置缓存等操作
func SetTokenCache(token *string, hash string, vsf string) {
	//设置用户缓存
	SetUserCache(hash)
	//将token信息存储到token中
	tokenInfo := cache.TokenCache{
		F:    common.GetRandomStr(16),
		Hash: hash,
		Vsf:  "vsf",
	}
	if err := cache.SetDataToCache("token_cache:"+*token, tokenInfo); err != nil {
		common.ThrowError(500, errors.New("设置用token缓存失败"))
	}
	// // 用户的各个平台的统一登陆 记录缓存
	TokenDecodedInfo := cache.TokenDecoded{}
	if err := cache.GetDataInCache("token_cache:"+*token, TokenDecodedInfo); err != nil {
		common.ThrowError(500, errors.New("获取oken缓存失败"))
	}
	if len(TokenDecodedInfo) == 0 {
		TokenDecodedInfo = map[string]interface{}{
			vsf: tokenInfo.F,
		}
	} else {
		TokenDecodedInfo[vsf] = tokenInfo.F
	}
	if err := cache.SetDataToCache("decoded_cache:"+hash, TokenDecodedInfo); err != nil {
		common.ThrowError(500, errors.New("获取oken缓存失败"))
	}
}

//SetUserCache 设置用户缓存
func SetUserCache(hash string) {
	// 登录成功 进行token 缓存 等操作
	userInfo := &(model.User{})
	err := model.DB.Where("users.hash = ?", hash).
		Joins(`left join accounts on accounts.user_id = users.id`).
		Select([]string{"users.id", "users.hash", "users.created_at", "users.nick_name", "users.real_name", "users.sex", "users.avatar", "users.state", "users.state"}).
		Preload("Account").
		First(userInfo).Error
	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
		}
		common.ThrowError(500, err)
	}
	//将用户信息存储在redis中
	cacheUserInfo := cache.UserCacheInfo.RenderUserCache(userInfo)
	if err = cache.SetDataToCache("user_cache:"+userInfo.Hash, cacheUserInfo); err != nil {
		common.ThrowError(500, errors.New("设置用户缓存失败"))
	}
	if err != nil {
		common.ThrowError(500, errors.New("设置用户缓存失败"))
	}
}
