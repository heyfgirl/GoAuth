package model

import (
	"bytes"
	"database/sql/driver"
	"errors"
	"fmt"
	"strings"

	config "GoAuth/configs"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres" //加载postgres驱动
	"github.com/sirupsen/logrus"
)

//JSON 自定义json类型
type JSON []byte

//Value 自定义JSON格式
func (j JSON) Value() (driver.Value, error) {
	if j.IsNull() {
		return nil, nil
	}
	return string(j), nil
}

//Scan 自定义JSON格式
func (j *JSON) Scan(value interface{}) error {
	if value == nil {
		*j = nil
		return nil
	}
	s, ok := value.([]byte)
	if !ok {
		return errors.New("Invalid Scan Source")
	}
	*j = append((*j)[0:0], s...)
	return nil
}

//MarshalJSON 自定义JSON格式
func (j JSON) MarshalJSON() ([]byte, error) {
	if j == nil {
		return []byte("null"), nil
	}
	return j, nil
}

//UnmarshalJSON 自定义JSON格式
func (j *JSON) UnmarshalJSON(data []byte) error {
	if j == nil {
		return errors.New("null point exception")
	}
	*j = append((*j)[0:0], data...)
	return nil
}

//IsNull 自定义JSON格式
func (j JSON) IsNull() bool {
	return len(j) == 0 || string(j) == "null"
}

//Equals 自定义JSON格式
func (j JSON) Equals(j1 JSON) bool {
	return bytes.Equal([]byte(j), []byte(j1))
}

//Models 数据库统一
var Models = []interface{}{
	&User{}, &Account{},
}

//DB 导出数据库orm.Ormer 全局使用
var DB *gorm.DB

func init() {
	var err error
	var pgUser = config.Conf.Pg.UserName
	var PassWord = config.Conf.Pg.PassWord
	var Host = config.Conf.Pg.Host
	// var Port = config.Conf.Pg.Port
	var Db = config.Conf.Pg.Db
	var pgMaxOpenConns = 1000
	var pgMaxIdleConns = 20
	var dataBaseURL = "host=" + Host + " user=" + pgUser + " dbname=" + Db + " sslmode=disable password=" + PassWord
	// fmt.Printf("%p, %T\n", DB, DB)
	DB, err = gorm.Open("postgres", dataBaseURL)
	// fmt.Printf("%p, %T\n", DB, DB)
	if err != nil {
		// panic(err)
		logrus.Error(err)
	}
	DB.LogMode(true)                        //调试模式
	DB.DB().SetMaxIdleConns(pgMaxOpenConns) //配置连接池
	DB.DB().SetMaxOpenConns(pgMaxIdleConns)
	if err = DB.AutoMigrate(Models...).Error; nil != err {
		fmt.Printf("auto migrate tables failed: %s", err.Error())
	}
	autoTags(DB, Models...)
	// defer DB.Close()
}

func autoTags(s *gorm.DB, values ...interface{}) {
	db := s.Unscoped()
	for _, value := range values {
		scope := db.NewScope(value)
		var tags = map[string]string{}
		var uniqueConstraints = map[string][]string{}
		for _, field := range scope.GetStructFields() {
			if name, ok := field.TagSettings["UNIQUE_CONSTRAINT"]; ok {
				names := strings.Split(name, ",")
				for _, name := range names {
					if name == "UNIQUE_CONSTRAINT" || name == "" {
						name = fmt.Sprintf("uix_%v_%v", scope.TableName(), field.DBName)
					}
					uniqueConstraints[name] = append(uniqueConstraints[name], field.DBName)
				}
			}
		}
		for _, columns := range uniqueConstraints {
			tagKey := scope.TableName() + "_" + strings.Join(columns, "_")
			fmt.Println(tagKey)
			if len(fmt.Sprintf("UNIQUE (%s)", strings.Join(columns, ","))) != 0 {
				tags[tagKey] = fmt.Sprintf("UNIQUE (%s)", strings.Join(columns, ","))
			}
			// tags[tagKey] = append(tags[tagKey], fmt.Sprintf("UNIQUE (%s)", strings.Join(columns, ",")))
		}
		fmt.Println(tags)
		for key, tag := range tags {
			if len(tag) != 0 {
				sqlCreate := "alter table" + " " + scope.TableName() + " " +
					"add constraint" + " " + key + " " +
					tag
				scope.Raw(sqlCreate).Exec()
				// fmt.Println(sqlCreate)
			}
		}
	}
}
