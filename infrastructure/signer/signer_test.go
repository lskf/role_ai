package signer

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"testing"
)

func TestGenerateAppId(t *testing.T) {
	appId := GenerateAppId()
	t.Log(appId)
}

func TestGenerateSecret(t *testing.T) {
	secret, err := GenerateSecret()
	if err != nil {
		t.Error(err)
	}
	t.Log(secret)
}

//func InitSigner() *Signer {
//	return NewSigner("243447485446", "TO-BmuS11vdzHep1PUHPYjcx4uECV4icuSZtQA62AyZP67B")
//}

type ReqTest struct {
	Symbol string `json:"symbol"`
	Amount string `json:"amount"`
}

type ReqGetTransferAmount struct {
	Target string `json:"target"`
	Amount string `json:"amount"`
}

type ReqGetTransferList struct{}

type ReqTransfer struct {
	Target string `json:"target"`
	Amount string `json:"amount"`
}

type ReqUseCode struct {
	Code string `json:"code"`
}

type ReqGetTransferDetail struct {
	OrderId string `json:"order_id"`
}

func TestSigner_Sign(t *testing.T) {
	//signer := InitSigner()
	data := &ReqTest{
		Symbol: "test1",
		Amount: "1",
	}
	signData := &SignData{
		Method:    "POST",
		Url:       "/api/v1/transfer",
		Data:      data,
		AppId:     "IIStTBnP",
		AppSecret: "RV70jlOj5tLl0FsvkMJ0zZS8NvzZd1kC",
	}
	signed, err := Sign(signData)
	if err != nil {
		t.Fatal(err)
	}
	signed.Print()
	// 测试将验签与签名的值修改为不同的内容
	//data.SendAmount = "2"
	dataJson, _ := json.Marshal(data)
	// 正确的url
	url := "http://localhost:8080/api/v1/transfer"
	// 错误的url 测试验签失败
	//url := "http://localhost:8080/api/v1/transfer/failed"
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader([]byte(dataJson)))
	if err != nil {
		t.Fatal("创建测试请求错误：", err.Error())
	}
	if req.Header == nil {
		req.Header = make(http.Header)
	}
	req.Header.Set(string(HeaderAppId), signed.AppId)
	req.Header.Set(string(HeaderTimestamp), signed.Timestamp)
	req.Header.Set(string(HeaderContentType), signed.ContentType)
	req.Header.Set(string(HeaderContentMD5), signed.ContentMD5)
	req.Header.Set(string(HeaderNonce), signed.Nonce)
	req.Header.Set(string(HeaderSign), signed.Sign)
	req.Body = io.NopCloser(bytes.NewReader(dataJson))
	ok, err := Verify(req, "TO-BmuS11vdzHep1PUHPYjcx4uECV4icuSZtQA62AyZP67B")
	if err != nil {
		t.Fatal("验签报错：", err.Error())
	}
	if !ok {
		t.Fatal("验证失败")
	}
	t.Log("验证成功")
}

func TestGetTransferAmountSign(t *testing.T) {
	//查询转换所得金额
	data := &ReqGetTransferAmount{
		Target: "test2",
		Amount: "10",
	}
	signData := &SignData{
		Method:    "GET",
		Url:       "/api/v1/transform/query/amount",
		Data:      data,
		AppId:     "IIStTBnP",
		AppSecret: "RV70jlOj5tLl0FsvkMJ0zZS8NvzZd1kC",
	}
	signed, err := Sign(signData)
	if err != nil {
		t.Fatal(err)
	}
	signed.Print()
}

func TestGetTransferListSign(t *testing.T) {
	//查询转换所得金额
	data := &ReqGetTransferList{}
	signData := &SignData{
		Method:    "GET",
		Url:       "/api/v1/transform/query/list",
		Data:      data,
		AppId:     "IIStTBnP",
		AppSecret: "RV70jlOj5tLl0FsvkMJ0zZS8NvzZd1kC",
	}
	signed, err := Sign(signData)
	if err != nil {
		t.Fatal(err)
	}
	signed.Print()
}

func TestTransferSign(t *testing.T) {
	//转换金额
	data := &ReqTransfer{
		Target: "test2",
		Amount: "10",
	}
	signData := &SignData{
		Method:    "POST",
		Url:       "/api/v1/transform",
		Data:      data,
		AppId:     "IIStTBnP",
		AppSecret: "RV70jlOj5tLl0FsvkMJ0zZS8NvzZd1kC",
	}
	signed, err := Sign(signData)
	if err != nil {
		t.Fatal(err)
	}
	signed.Print()
}

func TestUseCode(t *testing.T) {
	//使用兑换码
	data := &ReqUseCode{
		Code: "ByL0Cd2m7fBDAYpTKRwNUXaSzuUvUIkR",
	}
	signData := &SignData{
		Method:    "POST",
		Url:       "/api/v1/transform/exchangeCode",
		Data:      data,
		AppId:     "xYXek0dc",
		AppSecret: "26DhgDiW7OfHLVub10ldNOt2Il1YS5Ea",
	}
	signed, err := Sign(signData)
	if err != nil {
		t.Fatal(err)
	}
	signed.Print()
}

func TestGetTr(t *testing.T) {
	//使用兑换码
	data := &ReqGetTransferDetail{
		OrderId: "1714033127682567",
	}
	signData := &SignData{
		Method:    "GET",
		Url:       "/api/v1/transform/query/detail",
		Data:      data,
		AppId:     "IIStTBnP",
		AppSecret: "RV70jlOj5tLl0FsvkMJ0zZS8NvzZd1kC",
	}
	signed, err := Sign(signData)
	if err != nil {
		t.Fatal(err)
	}
	signed.Print()
}
