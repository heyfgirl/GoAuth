package common

import (
	"errors"
	"fmt"
	"net/http"
	"path/filepath"
	"regexp"
	"runtime"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

//Result 定义 Result 结构体
type Result struct {
	XToken  string `json:"x-token"`
	Cache   bool   `json:"cache"`
	Success bool   `json:"success"`
	// Duration int64       `json:"duration"`
	Data      interface{} `json:"data"`
	ErrorCode int         `json:"error_code"`
	Error     interface{} `json:"error"`
}

// 定义错误码
const (
	// 成功
	CodeSuccess int = 0
	//失败
	CodeError int = -1
	//自定义...
)

// type MyError struct {
// 	Code    int
// 	Message string
// }

// func (e MyError) Error() string {
// 	return fmt.Sprintf("%v: %v", e.Code, e.Message)
// }

//ThrowError 抛出错误方法
func ThrowError(code int, errorinfo error) {
	//重新组织错误
	s := errorinfo.Error()
	// fmt.Println(s)
	errorinfo = errors.New("code[" + strconv.Itoa(code) + "]::" + s)
	panic(errorinfo)
}

//RecoverError 抛出错误让统一处理中间件捕获
func RecoverError(ctx *gin.Context) {
	if r := recover(); r != nil {
		// r := r.(gin.H)
		// code := (r["code"]).(int)
		currentTime := GetTimeStr()
		// 定义 文件名、行号、方法名
		fileName, line, functionName := "?", 0, "?"
		pc, fileName, line, ok := runtime.Caller(3)
		if ok {
			functionName = runtime.FuncForPC(pc).Name()
			functionName = filepath.Ext(functionName)
			functionName = strings.TrimPrefix(functionName, ".")
		}
		message := fmt.Sprintf("%s", r)
		fmt.Println(message)

		Reg := `code\[([\d]+)\]::([\s\S]+?)$`
		Rgx, _ := regexp.Compile(Reg)
		RgxData := Rgx.FindSubmatch([]byte(message))
		//获取errorcode
		var ErrorCodeInt int
		var codeString string
		var ErrorMsgString string
		if RgxData != nil {
			codeString = string(RgxData[1])
		}
		if len(codeString) == 0 {
			ErrorCodeInt = 500
		} else {
			ErrorCodeInt2, err := strconv.Atoi(codeString)
			if err != nil {
				ErrorCodeInt = 500
			}
			ErrorCodeInt = ErrorCodeInt2
		}
		//获取errormessage
		if RgxData != nil {
			ErrorMsgString = string(RgxData[2])
		}
		if len(ErrorMsgString) == 0 {
			ErrorMsgString = message
		}

		errorJSONInfo := gin.H{
			"message":  ErrorMsgString,
			"filename": fileName,
			"line":     line,
			"function": functionName,
			"time":     currentTime,
		}

		JSONError(ctx, errorJSONInfo, ErrorCodeInt)
		return
	}
}

//JSONData 返回json数据
func JSONData(ctx *gin.Context, data interface{}) {
	token := ctx.Request.Header.Get("x-token")
	// durationStart, _ := ctx.Get("duration_start")
	// durationStart64 := durationStart.(int64)
	// durationEnd := GetTimeUnixNano()
	// fmt.Print(durationEnd)
	result := Result{
		XToken:  token,
		Cache:   false,
		Success: true,
		// Duration: durationEnd - durationStart64,
		ErrorCode: CodeSuccess,
		Data:      data,
	}
	ctx.JSON(http.StatusOK, result)
	ctx.Abort()
	return
}

//JSONError 返回json数据
func JSONError(ctx *gin.Context, data interface{}, code int) {
	token := ctx.Request.Header.Get("x-token")
	result := Result{
		XToken:    token,
		Cache:     false,
		Success:   false,
		ErrorCode: code,
		Error:     data,
	}
	ctx.JSON(http.StatusOK, result)
	ctx.Abort()
	return
}
