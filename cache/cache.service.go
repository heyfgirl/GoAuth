package cache

import (
	redis "GoAuth/models/redis"
	"encoding/json"
)

//SetDataToCache ...
func SetDataToCache(key string, data interface{}) error {
	jsonBytes, err := json.Marshal(data)
	if err != nil {
		return err
	}
	_, err = redis.SetJSON(key, jsonBytes)
	if err != nil {
		return err
	}
	return nil
}

//GetDataInCache ...
func GetDataInCache(key string, data interface{}) error {
	redisInfo, err := redis.GetJSON(key)
	if err != nil {
		return err
	}
	if redisInfo == nil {
		return nil
	}
	json.Unmarshal([]byte(redisInfo), data)
	return nil
}

//DelDataInCache 设置缓存
func DelDataInCache(key string) error {
	_, err := (*redis.Client).Do("del", key)
	if err != nil {
		return err
	}
	return nil
}

//GetAndDelDataInCache 设置缓存
func GetAndDelDataInCache(key string) (map[string]interface{}, error) {
	var data map[string]interface{}
	redisInfo, err := redis.GetJSON(key)
	if err != nil {
		return nil, err
	}
	if redisInfo == nil {
		return nil, nil
	}
	err = json.Unmarshal([]byte(redisInfo), &data)
	if err != nil {
		return nil, err
	}
	_, err = (*redis.Client).Do("del", key)
	if err != nil {
		return nil, err
	}
	return data, nil
}
