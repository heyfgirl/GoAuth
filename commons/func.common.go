package common

import (
	"crypto/md5"
	"encoding/hex"
	"math/rand"

	// "fmt"
	// // "ginDemo/config"
	// "net/url"

	"time"
)

//GetTimeStr ...
func GetTimeStr() string {
	return time.Now().Format("2006-01-02 15:04:05")
}

//GetTimeUnix 获取当前时间戳
func GetTimeUnix() int64 {
	return time.Now().Unix()
}

//GetTimeUnixNano 获取当前时间戳毫秒
func GetTimeUnixNano() int64 {
	return time.Now().UnixNano() / 1e6
}

// MD5 方法
func MD5(str string) string {
	s := md5.New()
	s.Write([]byte(str))
	return hex.EncodeToString(s.Sum(nil))
}

//GetRandomStr 生成随机字符穿
func GetRandomStr(length int) string {
	str := "0123456789abcdefghijklmnopqrstuvwxyz_"
	bytes := []byte(str)
	result := []byte{}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < length; i++ {
		result = append(result, bytes[r.Intn(len(bytes))])
	}
	return string(result)
}

//CreateSign 生成签名
// func CreateSign(params url.Values) string {
// 	var key []string
// 	var str = ""
// 	for k := range params {
// 		if k != "sn" && k != "ts" && k != "debug" {
// 			key = append(key, k)
// 		}
// 	}
// 	sort.Strings(key)
// 	for i := 0; i < len(key); i++ {
// 		if i == 0 {
// 			str = fmt.Sprintf("%v=%v", key[i], params.Get(key[i]))
// 		} else {
// 			str = str + fmt.Sprintf("&%v=%v", key[i], params.Get(key[i]))
// 		}
// 	}

// 	// 自定义签名算法
// 	// sign := MD5(MD5(str) + MD5(config.APP_NAME + config.APP_SECRET))
// 	return sign
// }
