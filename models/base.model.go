package model

import (
	"time"

	"github.com/jinzhu/gorm/dialects/postgres"
	_ "github.com/jinzhu/gorm/dialects/postgres" //加载postgres驱动
)

//DefaultModel 基本模型的定义
type DefaultModel struct {
	ID        int        `gorm:"not null;PRIMARY_KEY;AUTO_INCREMENT" json:"id" form:"id"`
	Hash      string     `gorm:"not null;size:32;unique;" json:"hash" form:"hash"`
	CreatedAt time.Time  `gorm:"not null;default:now()" json:"created_at" form:"created_at"`
	UpdatedAt time.Time  `gorm:"not null;default:now()" json:"updated_at" form:"updated_at"`
	DeletedAt *time.Time `gorm:"default:NULL::timestamp with time zone" json:"deleted_at" form:"deleted_at"`
}

//Account 账号表
type Account struct {
	DefaultModel
	UserID    int            `gorm:"not null;index" json:"user_id" form:"user_id"`                                         // 外键 (属于), tag `index`是为该列创建索引
	Type      string         `gorm:"not null;type:varchar(6);unique_constraint:token_type_unique" json:"type" form:"type"` //unique_index:account_type_token
	Token     string         `gorm:"not null;type:varchar(100);unique_constraint:token_type_unique" json:"token" form:"token"`
	Access    string         `gorm:"type:varchar(100);" json:"access" form:"access"`
	Refresh   string         `gorm:"type:varchar(6);" json:"refresh" form:"refresh"`
	Extension postgres.Jsonb `gorm:"not null;default:'{}'::jsonb" json:"extension" form:"extension"`
}

//User 用户表
type User struct {
	DefaultModel
	NickName  string         `gorm:"size:16;default:''" json:"nickname" form:"nickname"`
	RealName  string         `gorm:"size:16;default:''" json:"realname" form:"realname"`
	Sex       int            `gorm:"default:0" json:"sex" form:"sex"`
	Avatar    string         `gorm:"not null" json:"avatar" form:"avatar"`
	State     int            `gorm:"not null;default:1" json:"state" form:"state"`
	Role      string         `gorm:"size:16;not null;default:'normal'" json:"role" form:"role"`
	Extension postgres.Jsonb `gorm:"not null;default:'{}'::jsonb" json:"extension" form:"extension"`
	Account   []Account      `json:"account" form:"account"`
}
