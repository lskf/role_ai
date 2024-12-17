package signer

import (
	"bytes"
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"io"
	"net/http"
	"role_ai/infrastructure/utils"
	"time"
)

// Sign 签名数据
func Sign(data *SignData) (*Signed, error) {
	nonce, err := utils.GenerateNonce()
	if err != nil {
		return nil, err
	}
	signHeader := &Signed{
		Method:      data.Method,
		AppId:       data.AppId,
		Timestamp:   fmt.Sprintf("%d", time.Now().Unix()),
		ContentType: JSONContentType,
		Nonce:       nonce,
		Url:         data.Url,
		AppSecret:   data.AppSecret,
	}
	if data.Data != nil {
		// 编译为 JSON 字符串
		dataJson, err := json.Marshal(data.Data)
		if err != nil {
			return nil, err
		}
		// 计算内容 MD5
		signHeader.ContentMD5, err = doContentMD5(string(dataJson))
		if err != nil {
			return nil, err
		}
	}
	signHeader.Sign, err = doSign(signHeader)
	if err != nil {
		return nil, err
	}
	return signHeader, nil
}

func Verify(req *http.Request, appSecret string) (bool, error) {
	appId := req.Header.Get(string(HeaderAppId))
	timestamp := req.Header.Get(string(HeaderTimestamp))
	contentType := req.Header.Get(string(HeaderContentType))
	reqMd5 := req.Header.Get(string(HeaderContentMD5))
	nonce := req.Header.Get(string(HeaderNonce))
	sign := req.Header.Get(string(HeaderSign))
	url := req.URL.Path
	var contentMD5 string
	body, err := io.ReadAll(req.Body)
	if err != nil {
		return false, errors.Wrap(err, "读取请求体错误: ")
	}
	_ = req.Body.Close()
	bodyStr := string(body)
	if (utils.TrimString(bodyStr) != "" && reqMd5 == "") ||
		(utils.TrimString(bodyStr) == "" && reqMd5 != "") {
		return false, errors.New("请求体MD5校验错误")
	}
	// 重新设置请求体
	req.Body = io.NopCloser(bytes.NewBuffer(body))

	// 计算内容 MD5
	if len(body) > 0 {
		contentMD5, err = doContentMD5(string(body))
	}
	if err != nil {
		return false, errors.Wrap(err, "计算内容 MD5 错误: ")
	}
	// 验证请求体 MD5
	if contentMD5 != reqMd5 {
		return false, errors.New("请求体MD5校验错误")
	}
	// 验证签名
	signHeader := &Signed{
		Method:      req.Method,
		AppId:       appId,
		Timestamp:   timestamp,
		ContentType: contentType,
		ContentMD5:  contentMD5,
		Nonce:       nonce,
		Url:         url,
		AppSecret:   appSecret,
		Sign:        sign,
	}
	signature, err := doSign(signHeader)
	if err != nil {
		return false, errors.Wrap(err, "验证签名错误:")
	}
	if sign == signature {
		return true, nil
	}
	return false, errors.New("签名验证失败")
}

// doSign 计算签名 HMAC-SHA256
func doSign(data *Signed) (string, error) {
	var signData bytes.Buffer
	signData.WriteString(data.Method)
	signData.WriteString("\n")
	signData.WriteString(data.AppId)
	signData.WriteString("\n")
	signData.WriteString(data.Timestamp)
	signData.WriteString("\n")
	signData.WriteString(data.ContentType)
	signData.WriteString("\n")
	signData.WriteString(data.ContentMD5)
	signData.WriteString("\n")
	signData.WriteString(data.Nonce)
	signData.WriteString("\n")
	signData.WriteString(data.Url)
	h := hmac.New(sha256.New, []byte(data.AppSecret))
	if _, err := h.Write(signData.Bytes()); err != nil {
		return "", err
	}
	return hex.EncodeToString(h.Sum(nil)), nil
}

// doContentMD5 计算内容数据的 MD5 值
func doContentMD5(data string) (string, error) {
	hash := md5.New()
	_, err := hash.Write([]byte(data))
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", hash.Sum(nil)), nil
}
