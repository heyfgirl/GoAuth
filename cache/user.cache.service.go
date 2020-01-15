package cache

//底层共用包
import (
	common "GoAuth/commons"
	model "GoAuth/models"
	"time"

	"github.com/gin-gonic/gin"
)

//UserCache 用户缓存数据
type UserCache struct {
	ID        int               `json:"id" form:"id"`
	Hash      string            `json:"hash" form:"hash"`
	NickName  string            `json:"nickname" form:"nickname"`
	RealName  string            `json:"realname" form:"realname"`
	Sex       int               `json:"sex" form:"sex"`
	Avatar    string            `json:"avatar" form:"avatar"`
	State     int               `json:"state" form:"state"`
	Role      string            `json:"role" form:"role"`
	Account   map[string]string `json:"account" form:"account"`
	CreatedAt time.Time         `json:"created_at" form:"created_at"`
}

//TokenCache Token缓存数据
type TokenCache struct {
	F    string `json:"f" form:"f"`
	Hash string `json:"hash" form:"hash"` //用户hash
	Vsf  string `json:"vsf" form:"vsf"`   //平台
}

//TokenDecoded 不同平台Token唯一性炎症
type TokenDecoded map[string]interface{}

//UserCacheInfo 初始化缓存信息
var UserCacheInfo = newUserCacheF()

func newUserCacheF() *UserCache {
	return &UserCache{}
}

//RenderUserCache 加载用户信息到缓存
func (uc2 *UserCache) RenderUserCache(user *model.User) *UserCache {
	uc := &UserCache{}
	//将用户信息结构体同步到缓存用户模块
	uc.ID = user.ID
	uc.Hash = user.Hash
	uc.NickName = user.NickName
	uc.Sex = user.Sex
	uc.Avatar = user.Avatar
	uc.State = user.State
	uc.CreatedAt = user.CreatedAt

	uc.Account = map[string]string{}
	for _, account := range user.Account {
		uc.Account[account.Type] = account.Token
	}
	*uc2 = *uc
	return uc2
}

//CheckoutToken 检验token
func CheckoutToken(token string) string {
	if len(token) != 32 {
		token = common.GetRandomStr(32)
	}
	return token
}

//SetUserToCtx ...设置用户信息到上下文
func (uc *UserCache) SetUserToCtx(ctx *gin.Context) {
	ctx.Set("userID", uc.ID)
	ctx.Set("userInfo", uc)
}

//GetUserIDToCtx ...获取用户id
func (uc *UserCache) GetUserIDToCtx(ctx *gin.Context) int {
	ucA, exist := ctx.Get("userID")
	if !exist {
		return 0
	}
	userID := ucA.(int)
	return userID
}

//GetUserInCtx ...在上下文中获取用户信息
func (uc *UserCache) GetUserInCtx(ctx *gin.Context) *UserCache {
	ucA, exist := ctx.Get("userInfo")
	if !exist {
		return nil
	}
	uc = ucA.(*UserCache)
	return uc
}
