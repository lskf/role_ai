package tests

import (
	"encoding/json"
	"fmt"
	"github.com/imroc/req/v3"
	"github.com/tidwall/gjson"
	"net/http"
	"role_ai/infrastructure/encrypt"
	"role_ai/infrastructure/signer"
	"testing"
)

const (
	BaseURL = "http://localhost:8086"
	//BaseURL = "https://token-otc.xxdev.top"
)

type Resp struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

func ParseResp(resp string) (*Resp, error) {
	var r Resp
	if err := json.Unmarshal([]byte(resp), &r); err != nil {
		return nil, err
	}
	return &r, nil
}

func ProcessingResp(t *testing.T, resp *req.Response) {
	if resp.Response.StatusCode != http.StatusOK {
		t.Fatalf("请求失败: %s", resp.String())
	}
	r, err := ParseResp(resp.String())
	if err != nil {
		t.Fatal("解析失败: ", err.Error())
	}
	if r.Code != 0 {
		t.Fatalf("请求失败: %s", r.Msg)
	}
	t.Logf("请求成功: %s", r.Data)
}

func doGet(signed *signer.Signed, uri string) (*req.Response, error) {
	return req.R().
		SetHeader("X-Tsign-App-Id", signed.AppId).
		SetHeader("X-Tsign-Timestamp", signed.Timestamp).
		SetHeader("Content-Type", signed.ContentType).
		SetHeader("X-Tsign-Nonce", signed.Nonce).
		SetHeader("Content-MD5", signed.ContentMD5).
		SetHeader("X-Tsign-Sign", signed.Sign).
		Get(BaseURL + uri)
}

func doPost(signed *signer.Signed, uri string, data interface{}) (*req.Response, error) {
	return req.R().SetBody(data).
		SetHeader("X-Tsign-App-Id", signed.AppId).
		SetHeader("X-Tsign-Timestamp", signed.Timestamp).
		SetHeader("Content-Type", signed.ContentType).
		SetHeader("X-Tsign-Nonce", signed.Nonce).
		SetHeader("Content-MD5", signed.ContentMD5).
		SetHeader("X-Tsign-Sign", signed.Sign).
		Post(BaseURL + uri)
}

// TestQueryTransferAmount 测试计算兑换可得对应代币数量
func TestQueryTransferAmount(t *testing.T) {
	reqUri := "/api/v1/transform/transfer/amount"
	type Req struct {
		FromSymbol   string `json:"from_symbol"`   // 转换发送代币
		TargetSymbol string `json:"target_symbol"` // 转换目标代币
		Amount       string `json:"amount"`        // 转换发送数量
	}
	usecases := []struct {
		Name string
		Data Req
	}{
		{
			Name: "xxy->xsee",
			Data: Req{
				FromSymbol:   "xxy",
				TargetSymbol: "xsee",
				Amount:       "1000",
			},
		},
	}
	for _, usecase := range usecases {
		t.Run(usecase.Name, func(t *testing.T) {
			t.Logf("FromSymbol: %s, TargetSymbol: %s, Amount: %s", usecase.Data.FromSymbol, usecase.Data.TargetSymbol, usecase.Data.Amount)
			signed, err := signer.Sign(&signer.SignData{
				Method:    http.MethodPost,
				Url:       reqUri,
				Data:      usecase.Data,
				AppId:     "243058621222",
				AppSecret: "TO-H1Exk8Z8bBNDEEnNFfkgi3otHkXMEx8vYxHNor87zwYE",
			})
			if err != nil {
				t.Fatal("签名失败: ", err.Error())
			}
			reqUri := reqUri + fmt.Sprintf("?from_symbol=%s&target_symbol=%s&amount=%s",
				usecase.Data.FromSymbol, usecase.Data.TargetSymbol, usecase.Data.Amount)
			resp, err := doPost(signed, reqUri, usecase.Data)
			ProcessingResp(t, resp)
			fmt.Println(resp.String())
		})
	}
}

// TestTransfer 测试兑换
func TestTransfer(t *testing.T) {
	var uri = "/api/v1/transform/transfer"
	type reqData struct {
		Uid          string `json:"uid"`
		FromSymbol   string `json:"from_symbol"`
		TargetSymbol string `json:"target_symbol"`
		Amount       string `json:"amount"`
	}
	usecases := []struct {
		Name string
		Data *reqData
	}{
		{
			Name: "xxy->xsee",
			Data: &reqData{
				Uid:          "123456",
				FromSymbol:   "xxy",
				TargetSymbol: "xsee",
				Amount:       "1",
			},
		},
	}
	for _, usecase := range usecases {
		t.Run(usecase.Name, func(t *testing.T) {
			t.Logf("Uid: %s, FromSymbol: %s, TargetSymbol: %s, Amount: %s",
				usecase.Data.Uid, usecase.Data.FromSymbol, usecase.Data.TargetSymbol, usecase.Data.Amount)
			signed, err := signer.Sign(&signer.SignData{
				Method:    http.MethodPost,
				Url:       "/api/v1/transform/transfer",
				Data:      usecase.Data,
				AppId:     "243058621222",
				AppSecret: "TO-H1Exk8Z8bBNDEEnNFfkgi3otHkXMEx8vYxHNor87zwYE",
			})
			if err != nil {
				t.Fatal("签名失败: ", err.Error())
			}
			resp, err := doPost(signed, uri, usecase.Data)
			if err != nil {
				t.Fatal("请求失败: ", err.Error())
			}
			if resp.Response.StatusCode != http.StatusOK {
				t.Fatalf("请求失败: %s", resp.String())
			}

			ProcessingResp(t, resp)
			t.Logf("请求成功: %s", resp.String())
			rest := gjson.ParseBytes(resp.Bytes())
			decrypted, err := encrypt.Decrypt(rest.Get("data").String(), "w3iWiQgezOe38eR8VvVKAQTxrz7KuemNHLloZ2uqA+4=")
			if err != nil {
				t.Fatal("解密失败: ", err.Error())
			}
			t.Logf("解密结果: %s", decrypted)
		})
	}
}

// TestQueryTransformHistory 测试查询兑换历史记录
func TestQueryTransformHistory(t *testing.T) {
	uri := "/api/v1/transform/history"
	data := []struct {
		Uid       string
		PageNum   int
		PageSize  int
		AppId     string
		AppSecret string
	}{
		{
			Uid:       "123456",
			PageNum:   1,
			PageSize:  10,
			AppId:     "243058621222",
			AppSecret: "TO-H1Exk8Z8bBNDEEnNFfkgi3otHkXMEx8vYxHNor87zwYE",
		},
		{
			Uid:       "123456",
			AppId:     "243058621222",
			AppSecret: "TO-H1Exk8Z8bBNDEEnNFfkgi3otHkXMEx8vYxHNor87zwYE",
		},
	}
	for _, d := range data {
		t.Run("", func(t *testing.T) {
			signed, err := signer.Sign(&signer.SignData{
				Method:    http.MethodGet,
				Url:       uri,
				Data:      nil,
				AppId:     d.AppId,
				AppSecret: d.AppSecret,
			})
			if err != nil {
				t.Fatal("签名失败: ", err.Error())
			}
			uri = fmt.Sprintf("%s?uid=%s&page_num=%d&page_size=%d", uri, d.Uid, d.PageNum, d.PageSize)
			resp, err := doGet(signed, uri)
			if err != nil {
				t.Fatal("请求失败: ", err.Error())
			}
			ProcessingResp(t, resp)
			t.Logf("请求成功: %s", resp.String())

			rest := gjson.ParseBytes(resp.Bytes())
			decrypted, err := encrypt.Decrypt(rest.Get("data").String(), "w3iWiQgezOe38eR8VvVKAQTxrz7KuemNHLloZ2uqA+4=")
			if err != nil {
				t.Fatal("解密失败: ", err.Error())
			}
			t.Logf("解密结果: %s", decrypted)
		})
	}
}

func TestUseExchangeCode(t *testing.T) {
	uri := "/api/v1/transform/exchange"
	data := []struct {
		Code      string
		Uid       string
		Symbol    string
		AppId     string
		AppSecret string
		Encrypt   string
	}{
		{
			Code:      "Alst1QGH8fwCIsBi9LCTKqSTuUwMIXQH",
			Uid:       "234567",
			Symbol:    "xsee",
			AppId:     "242201119538",
			AppSecret: "TO-GMDm3yaTNgDtQ7r9Tm2Lr3ztgctpbr2eUQbKjq5ASusy",
			Encrypt:   "UfbTiPgz7ckMAjh22ZzSKxpQnM+/7yt58d7rqJBXad4=",
		},
	}
	for _, d := range data {
		t.Run("", func(t *testing.T) {
			var reqData = struct {
				Code string `json:"code"`
				Uid  string `json:"uid"`
			}{
				Code: d.Code,
				Uid:  d.Uid,
			}
			dataJson, err := json.Marshal(reqData)
			if err != nil {
				t.Fatal("序列化失败: ", err.Error())
			}
			encrypted, err := encrypt.Encrypt(string(dataJson), d.Encrypt)

			var enctrypedData = struct {
				Data string `json:"data"`
			}{
				Data: encrypted,
			}
			signed, err := signer.Sign(&signer.SignData{
				Method:    http.MethodPost,
				Url:       uri,
				Data:      enctrypedData,
				AppId:     d.AppId,
				AppSecret: d.AppSecret,
			})
			if err != nil {
				t.Fatal("签名失败: ", err.Error())
			}

			resp, err := doPost(signed, uri, enctrypedData)
			if err != nil {
				t.Fatal("请求失败: ", err.Error())
			}
			ProcessingResp(t, resp)
			t.Logf("请求成功: %s", resp.String())
		})
	}
}

func TestQueryTransformDetail(t *testing.T) {
	orderId, err := encrypt.Encrypt("1714212745622840", "w3iWiQgezOe38eR8VvVKAQTxrz7KuemNHLloZ2uqA+4=")
	if err != nil {
		t.Fatal("解密失败: ", err.Error())
	}
	data := []struct {
		OrderId string
	}{
		{
			OrderId: orderId,
		},
	}
	for _, d := range data {
		t.Run(d.OrderId, func(t *testing.T) {
			t.Logf("orderId: %s", d.OrderId)
			signed, err := signer.Sign(&signer.SignData{
				Method:    http.MethodGet,
				Url:       "/api/v1/transform/query/detail",
				Data:      d,
				AppId:     "243058621222",
				AppSecret: "TO-H1Exk8Z8bBNDEEnNFfkgi3otHkXMEx8vYxHNor87zwYE",
			})
			if err != nil {
				t.Fatal("签名失败: ", err.Error())
			}
			t.Logf("签名结果: %s", signed.Sign)
			resp, err := req.R().SetBody(d).
				SetHeader("X-Tsign-App-Id", signed.AppId).
				SetHeader("X-Tsign-Timestamp", signed.Timestamp).
				SetHeader("Content-Type", signed.ContentType).
				SetHeader("X-Tsign-Nonce", signed.Nonce).
				SetHeader("Content-MD5", signed.ContentMD5).
				SetHeader("X-Tsign-Sign", signed.Sign).
				Get(fmt.Sprintf("http://localhost:8086/api/v1/transform/query/detail?order_id=%v", d.OrderId))
			if err != nil {
				t.Fatal("请求失败: ", err.Error())
			}
			if resp.Response.StatusCode != http.StatusOK {
				t.Fatalf("请求失败: %s", resp.String())
			}
			t.Logf("请求成功: %s", resp.String())
			//var respData respData
			//if err := resp.UnmarshalJson(&respData); err != nil {
			//	t.Fatal("解析失败: ", err.Error())
			//}
			rest := gjson.ParseBytes(resp.Bytes())
			decrypted, err := encrypt.Decrypt(rest.Get("data").String(), "w3iWiQgezOe38eR8VvVKAQTxrz7KuemNHLloZ2uqA+4=")
			if err != nil {
				t.Fatal("解密失败: ", err.Error())
			}
			t.Logf("解密结果: %s", decrypted)
		})
	}
}

func TestQuery(t *testing.T) {
	uri := "/api/v1/transform/redeemable"
	data := []struct {
		FromSymbol string
		AppId      string
		AppSecret  string
	}{
		{
			FromSymbol: "xxy",
			AppId:      "243058621222",
			AppSecret:  "TO-H1Exk8Z8bBNDEEnNFfkgi3otHkXMEx8vYxHNor87zwYE",
		},
	}
	for _, d := range data {
		t.Run("", func(t *testing.T) {
			signed, err := signer.Sign(&signer.SignData{
				Method:    http.MethodGet,
				Url:       uri,
				AppId:     d.AppId,
				AppSecret: d.AppSecret,
			})
			if err != nil {
				t.Fatal("签名失败: ", err.Error())
			}
			t.Logf("签名结果: %s", signed.Sign)
			uri = fmt.Sprintf("%s?from_symbol=%s", uri, d.FromSymbol)
			resp, err := doGet(signed, uri)
			if err != nil {
				t.Fatal("请求失败: ", err.Error())
			}
			ProcessingResp(t, resp)
			t.Logf("请求成功: %s", resp.String())
		})
	}
}
