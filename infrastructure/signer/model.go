package signer

import "fmt"

// Headers 请求头
type Headers string

const (
	HeaderAppId       Headers = "X-Tsign-App-Id"    // AppId
	HeaderTimestamp   Headers = "X-Tsign-Timestamp" // 时间戳
	HeaderContentType Headers = "Content-Type"      // 请求体类型
	HeaderContentMD5  Headers = "Content-MD5"       // 请求体 MD5
	HeaderNonce       Headers = "X-Tsign-Nonce"     // 随机数
	HeaderSign        Headers = "X-Tsign-Sign"      // 签名
)

var JSONContentType = "application/json;charset=UTF-8"

func (h Headers) Value() string {
	return string(h)
}

type SignData struct {
	Method    string
	Url       string
	Data      interface{}
	AppId     string
	AppSecret string
}

type Signed struct {
	Method      string // 请求方法
	AppId       string // AppId
	Timestamp   string // 时间戳
	ContentType string // 请求体类型
	ContentMD5  string // 请求体 MD5
	Nonce       string // 随机数
	Url         string // 请求路径

	AppSecret string // AppSecret
	Sign      string // 签名
}

func (s *Signed) Print() {
	fmt.Println("Method:", s.Method, "AppId:", s.AppId,
		"Timestamp:", s.Timestamp, "ContentType:", s.ContentType,
		"ContentMD5:", s.ContentMD5, "Url:", s.Url, "Nonce:", s.Nonce,
		"Sign:", s.Sign)
}
