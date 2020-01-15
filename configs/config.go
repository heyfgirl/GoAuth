package config

import (
	"flag"
	"io/ioutil"

	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

//Conf 全局使用配置
var Conf *Config

//Config 配置文件结构
type Config struct {
	Env        string `yaml:"Env"`        // 环境：prod、dev
	BaseURL    string `yaml:"BaseUrl"`    // base url
	Port       string `yaml:"Port"`       // 端口
	LogFile    string `yaml:"LogFile"`    // 日志文件
	ShowSQL    bool   `yaml:"ShowSql"`    // 是否显示日志
	StaticPath string `yaml:"StaticPath"` // 静态文件目录

	Pg struct {
		UserName string `yaml:"UserName"` // 数据库连接用户
		PassWord string `yaml:"PassWord"` // 数据库连接密码
		Host     string `yaml:"Host"`     // 数据库连接Host
		Port     int    `yaml:"Port"`     // 数据库连接端口
		Db       string `yaml:"Db"`       // 数据库连接数据库名称
	} `yaml:"Pg"`
	DefaultUser struct {
		NickNmaePrefix string `yaml:"NickNmaePrefix"` // 默认用户昵称前缀
		Avatar         string `yaml:"Avatar"`         // 默认用户头像
	} `yaml:"DefaultUser"`
	Redis struct {
		MaxClients int    `yaml:"MaxClients"`
		MinClients int    `yaml:"MinClients"`
		Port       int    `yaml:"Port"`
		Host       string `yaml:"Host"`
		Password   string `yaml:"Password"`
		Db         int    `yaml:"Db"`
	} `yaml:"Redis"`
	Wechat struct {
		AppID     string `yaml:"AppID"`
		AppSecret string `yaml:"AppSecret"`

		WebID     string `yaml:"WebID"`
		WebSecret string `yaml:"WebSecret"`
	} `yaml:"Wechat"`
	RedisKey struct {
		WechatAccount string `yaml:"WechatAccount"`
	} `yaml:"RedisKey"`
}

//Init 初始化
func init() {
	var configFile = flag.String("config", "./config.yaml", "配置文件路径")
	yamlFile, err := ioutil.ReadFile(*configFile)
	if err != nil {
		logrus.Error(err)
		return
	}

	Conf = &Config{}
	err = yaml.Unmarshal(yamlFile, Conf)
	if err != nil {
		logrus.Error(err)
	}
}
