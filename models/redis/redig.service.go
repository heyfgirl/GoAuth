package redis

import (
	"fmt"
	"log"
	"os"
	"strconv"

	config "GoAuth/configs"

	redis "github.com/gomodule/redigo/redis"
)

//Client redis客户端
var Client *redis.Conn

//init 初始化
func init() {
	var err error
	ClientD, err := redis.Dial("tcp", config.Conf.Redis.Host+":"+strconv.Itoa(config.Conf.Redis.Port))
	if err != nil {
		log.Fatalln(err)
		os.Exit(-1)
	}
	if _, err := ClientD.Do("AUTH", config.Conf.Redis.Password); err != nil {
		fmt.Println(err)
		ClientD.Close()
	}
	if _, err := ClientD.Do("SELECT", config.Conf.Redis.Db); err != nil {
		fmt.Println(err)
		ClientD.Close()
	}
	redis.DialDatabase(1)

	Client = &ClientD
	// defer Client.Close()
}

//SetJSON 塞入json数据
func SetJSON(key string, data []byte) (bool, error) {
	// value, err := json.Marshal(data)
	// if err != nil {
	// 	common.ThrowError(500, err)
	// }
	nv, err := (*Client).Do("SET", key, data)
	n := false
	if nv == int64(1) {
		n = true
	} else {
		n = false
	}
	return n, err
}

//GetJSON 塞入json数据
func GetJSON(key string) ([]byte, error) {
	// var imapGet map[string]interface{}
	valueGet, err := redis.Bytes((*Client).Do("GET", key))
	if err != nil {
		//获取到空值
		if err == redis.ErrNil {
			return nil, nil
		}
		return nil, err
	}
	// errShal := json.Unmarshal(valueGet, &imapGet)
	// if errShal != nil {
	// 	return nil, errShal
	// }
	// fmt.Println(imapGet["username"])
	return valueGet, nil
}
