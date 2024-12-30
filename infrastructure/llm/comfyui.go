package llm

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/leor-w/kid/config"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

type IComfyUi interface {
}

type ComfyUi struct {
	options *Options
	cli     *http.Client
}

type PromptReq struct {
	ClientId    string `json:"client_id"`     //客户端id
	CkptName    string `json:"ckpt_name"`     //模型名称
	PictureNum  string `json:"picture_num"`   //图片数量
	Prompt      string `json:"prompt"`        //prompt
	ParaFileUrl string `json:"para_file_url"` //参数文件的路径
}

type PromptResp struct {
	PromptId   string `json:"prompt_id"`
	Number     int64  `json:"number"`
	NodeErrors any    `json:"node_errors"`
}

// 定义用于解析 JSON 的结构体
type OutputImage struct {
	Filename  string `json:"filename"`
	Subfolder string `json:"subfolder"`
	Type      string `json:"type"`
}

type Output struct {
	Images []OutputImage `json:"images"`
}

type HistoryDetail struct {
	//Prompt any `json:"prompt"`
	Outputs map[string]Output `json:"outputs"`
	Status  StatusObj         `json:"status"`
	//Meta    any               `json:"meta"`
}

type StatusObj struct {
	StatusStr string `json:"status_str"`
	Completed bool   `json:"completed"`
	Messages  any    `json:"messages"`
}

type ViewReq struct {
	FileName  string `json:"file_name"`
	Type      string `json:"type"`
	Subfolder string `json:"subfolder"`
}

func (comfyUi *ComfyUi) Provide(ctx context.Context) any {
	return comfyUi.New(
		WithBaseURL(config.GetString("llm.comfyUi.baseUrl")),
		WithApiKey(config.GetString("llm.comfyUi.apiKey")),
	)
}

func (comfyUi *ComfyUi) New(opts ...Option) *ComfyUi {
	option := &Options{}
	for _, opt := range opts {
		opt(option)
	}
	comfyUi.options = option
	comfyUi.cli = &http.Client{}
	return comfyUi
}

func (comfyUi *ComfyUi) NewComfyUi() *ComfyUi {
	baseUrl := config.GetString("llm.comfyUi.baseUrl")
	if comfyUi.options == nil {
		comfyUi.options = &Options{
			BaseURL: baseUrl,
		}
	}
	comfyUi.cli = &http.Client{}
	return comfyUi
}

func (comfyUi *ComfyUi) Prompt(para PromptReq) (*PromptResp, error) {
	dataByte, err := os.ReadFile(para.ParaFileUrl)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("读取参数文件失败fileUrl:%s,，err:%s", para.ParaFileUrl, err.Error()))
	}
	data := string(dataByte)
	//替换
	data = strings.Replace(data, "{{client_id}}", para.ClientId, -1)
	data = strings.Replace(data, "{{ckpt_name}}", para.CkptName, -1)
	data = strings.Replace(data, "{{picture_num}}", para.PictureNum, -1)
	data = strings.Replace(data, "{{prompt_str}}", para.Prompt, -1)
	data = strings.Replace(data, "{{seed}}", strconv.FormatInt(time.Now().UnixMicro(), 10), -1)
	req, err := NewRequest("POST", comfyUi.options.BaseURL+"/prompt", []byte(data), false)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("请求失败，fileUrl:%s,para:%+v,err:%s", para.ParaFileUrl, para, err.Error()))
	}
	// 发送请求并获取响应
	resp, err := comfyUi.cli.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	// 读取响应体
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	res := PromptResp{}
	err = json.Unmarshal(body, &res)
	if err != nil {
		return nil, err
	}
	return &res, nil
}

//// GetHistoryList
//// @Description: 获取历史列表（还没写好的）
//// @receiver comfyUi
//func (comfyUi *ComfyUi) GetHistoryList() {
//	req, err := NewRequest("GET", comfyUi.options.BaseURL+"/history", nil, false)
//	if err != nil {
//		fmt.Println(err.Error())
//	}
//	// 创建一个HTTP客户端
//	client := &http.Client{}
//	// 发送请求并获取响应
//	resp, err := client.Do(req)
//	if err != nil {
//		log.Fatal(err)
//	}
//	defer resp.Body.Close()
//	// 读取响应体
//	body, err := io.ReadAll(resp.Body)
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Println("Response Body:", string(body))
//}

func (comfyUi *ComfyUi) GetHistoryDetail(promptId string) (map[string]HistoryDetail, error) {
	req, err := NewRequest("GET", comfyUi.options.BaseURL+"/history/"+promptId, nil, false)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("获取prompt详情失败，promptId:%s,err:%s", promptId, err.Error()))
	}
	// 发送请求并获取响应
	resp, err := comfyUi.cli.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	// 读取响应体
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var result map[string]HistoryDetail
	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (comfyUi *ComfyUi) View(para ViewReq) (any, error) {
	url := "/view?filename={{filename}}&subfolder{{subfolder}}=&type={{type}}"
	url = strings.Replace(url, "{{filename}}", para.FileName, 1)
	url = strings.Replace(url, "{{subfolder}}", para.Subfolder, 1)
	url = strings.Replace(url, "{{type}}", para.Type, 1)

	req, err := NewRequest("GET", comfyUi.options.BaseURL+url, nil, false)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("获取图片失败，para:%+v,err:%s", para, err.Error()))
	}
	// 发送请求并获取响应
	resp, err := comfyUi.cli.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	// 读取响应体
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}
