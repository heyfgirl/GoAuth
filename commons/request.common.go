package common

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"sync"

	jsoniter "github.com/json-iterator/go"
)

//Adapter ...
type Adapter struct {
	pool sync.Pool
}

//New ...
func New() *Adapter {
	return &Adapter{
		pool: sync.Pool{
			New: func() interface{} {
				return bytes.NewBuffer(make([]byte, 4096))
			},
		},
	}
}

//Request ..
func (api *Adapter) Request() (*http.Response, error) {
	var err error
	buffer := api.pool.Get().(*bytes.Buffer)
	buffer.Reset()
	defer func() {
		if buffer != nil {
			api.pool.Put(buffer)
			buffer = nil
		}
	}()
	// e := jsoniter.NewEncoder(buffer)
	// err = e.Encode(req)

	req, err := http.NewRequest("GET", `https://api.weixin.qq.com/sns/oauth2/access_token?appid=wx05fb5b232abf83e7&secret=614d33541efe9a589a33147a7e10d6f6&code=`+"sfsdfds"+`&grant_type=authorization_code`, nil)
	if err != nil {
	}
	// req.Header.Set("User-Agent", "xxx")
	httpResponse, err := http.DefaultClient.Do(req)
	if httpResponse != nil {
		defer func() {
			io.Copy(ioutil.Discard, httpResponse.Body)
			httpResponse.Body.Close()
		}()
	}
	if err != nil {
	}
	if httpResponse.StatusCode != 200 {
	}
	_, err = io.Copy(buffer, httpResponse.Body)
	if err != nil {
		return nil, fmt.Errorf("adapter io.copy failure error:%v", err)
	}
	respData := buffer.Bytes()
	res := &http.Response{}
	err = jsoniter.Unmarshal(respData, res)
	fmt.Println(string(respData))
	buffer = nil
	return nil, nil
}
