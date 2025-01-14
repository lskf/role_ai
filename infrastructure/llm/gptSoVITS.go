package llm

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/leor-w/kid/config"
	"io"
	"net/http"
)

type IGptSoVITS interface {
	Tts(para GptSovitsTtsParam) (any, error)
	SetGptWeights(url string) (any, error)
	SetSovitsWeights(url string) (any, error)
}

type GptSoVITS struct {
	options *Options
	cli     *http.Client
}

type GptSovitsTtsParam struct {
	Text              string   `json:"text"`           //required
	TextLang          string   `json:"text_lang"`      //required
	RefAudioPath      string   `json:"ref_audio_path"` //required
	AuxRefAudioPaths  []string `json:"aux_ref_audio_paths,omitempty"`
	PromptText        string   `json:"prompt_text,omitempty"`
	PromptLang        string   `json:"prompt_lang"` //required
	TopK              int64    `json:"top_k,omitempty"`
	TopP              float64  `json:"top_p,omitempty"`
	Temperature       float64  `json:"temperature,omitempty"`
	TextSplitMethod   string   `json:"text_split_method,omitempty"`
	BatchSize         int64    `json:"batch_size,omitempty"`
	BatchThreshold    float64  `json:"batch_threshold,omitempty"`
	SplitBucket       bool     `json:"split_bucket,omitempty"`
	SpeedFactor       float64  `json:"speed_factor,omitempty"`
	StreamingMode     bool     `json:"streaming_mode,omitempty"`
	Seed              int64    `json:"seed,omitempty"`
	ParallelInfer     bool     `json:"parallel_infer,omitempty"`
	RepetitionPenalty float64  `json:"repetition_penalty,omitempty"`
}

func (gptSoVITS *GptSoVITS) Provide(ctx context.Context) any {
	return gptSoVITS.New(
		WithBaseURL(config.GetString("llm.tts.baseUrl")),
		WithApiKey(config.GetString("llm.tts.apiKey")),
	)
}

func (gptSoVITS *GptSoVITS) New(opts ...Option) *GptSoVITS {
	option := &Options{}
	for _, opt := range opts {
		opt(option)
	}
	gptSoVITS.options = option
	gptSoVITS.cli = &http.Client{}
	return gptSoVITS
}

func (gptSoVITS *GptSoVITS) NewGptSoVITS() *GptSoVITS {
	baseUrl := config.GetString("llm.gptSoVITS.baseUrl")
	if gptSoVITS.options == nil {
		gptSoVITS.options = &Options{
			BaseURL: baseUrl,
		}
	}
	gptSoVITS.cli = &http.Client{}
	return gptSoVITS
}

func (gptSoVITS *GptSoVITS) Tts(para GptSovitsTtsParam) (any, error) {

	data, err := json.Marshal(para)
	if err != nil {
		return nil, err
	}

	req, err := NewRequest("POST", gptSoVITS.options.BaseURL+"/tts", data, false)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("请求失败，para:%+v,err:%s", para, err.Error()))
	}
	// 发送请求并获取响应
	resp, err := gptSoVITS.cli.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, errors.New(resp.Status)
	}
	// 读取响应体
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return &body, nil
}

func (gptSoVITS *GptSoVITS) SetGptWeights(url string) (any, error) {
	req, err := NewRequest("GET", gptSoVITS.options.BaseURL+"/set_gpt_weights?weights_path="+url, nil, false)
	return req, err
}

func (gptSoVITS *GptSoVITS) SetSovitsWeights(url string) (any, error) {
	req, err := NewRequest("GET", gptSoVITS.options.BaseURL+"set_sovits_weights?weights_path="+url, nil, false)
	return req, err
}
