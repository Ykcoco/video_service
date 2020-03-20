package util

import(
	"log"
	"fmt"
	"encoding/json"
)

//ResMsg:http响应数据的通用结构
type RespMsg struct {
	Code int `json:"code"`
	Msg string `json:"msg"`
	Data interface{} `json:"data"`
}

//生成response对象
func NewRespMsg(code int, msg string, data interface{}) *RespMsg {
     return &RespMsg{
		 Code:code,
		 Msg : msg,
		 Data: data,
	 }
}

//对象转json格式的二进制数组
func (resp *RespMsg) JSONBytes() []byte {
	r, err := json.Marshal(resp)
	if err != nil {
       log.Println(err)
	}
	return r
}


//对象json转string
func (resp *RespMsg) JSONString() string {
	r, err := json.Marshal(resp)
	if err != nil {
		log.Println(err)
	}

	return string(r)
}


//只包含code和message的响应体
func GenSimpleRespStream(code int, msg string) []byte {
	return []byte(fmt.Sprintf(`{"code":%d,"msg":"%s"}`, code, msg))
}


// GenSimpleRespString : 只包含code和message的响应体(string)
func GenSimpleRespString(code int, msg string) string {
	return fmt.Sprintf(`{"code":%d,"msg":"%s"}`, code, msg)
}

