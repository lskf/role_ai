package llm

import (
	"context"
	"github.com/anthropics/anthropic-sdk-go"
	"github.com/anthropics/anthropic-sdk-go/option"
	"github.com/leor-w/kid/config"
)

type CompletionReq = anthropic.CompletionNewParams
type CompletionRes = anthropic.Completion
type MessageReq struct {
	Model    string
	MaxToken int64
	Messages []MessageObj
}
type MessageObj struct {
	User      string
	Assistant string
}
type MessageRes = anthropic.Message

type IClaude interface {
	Complete(req *CompletionReq) (res *CompletionRes, err error)
	Message(req *MessageReq) (res *MessageReq, err error)
}

type Claude struct {
	BaseUrl string
	ApiKey  string
	cli     *anthropic.Client
}

func (claude *Claude) Provide(ctx context.Context) any {
	return anthropic.NewClient(
		option.WithBaseURL(config.GetString("llm.claude.baseUrl")),
		option.WithAPIKey(config.GetString("llm.claude.apiKey")),
	)
}

func (claude *Claude) NewClient() *Claude {
	baseUrl := config.GetString("llm.claude.baseUrl")
	apiKey := config.GetString("llm.claude.apiKey")
	if claude.BaseUrl != "" {
		baseUrl = claude.BaseUrl
	}
	if claude.ApiKey != "" {
		apiKey = claude.ApiKey
	}
	claude.cli = anthropic.NewClient(
		option.WithBaseURL(baseUrl),
		option.WithAPIKey(apiKey),
	)
	return claude
}

func (claude *Claude) Complete(req CompletionReq) (res *CompletionRes, err error) {
	return nil, nil
}

func (claude *Claude) Message(req MessageReq) (res *MessageRes, err error) {
	messageParam := make([]anthropic.MessageParam, 0)
	for _, v := range req.Messages {
		if v.User != "" {
			messageParam = append(messageParam, anthropic.NewUserMessage(anthropic.NewTextBlock(v.User)))
		}
		if v.Assistant != "" {
			messageParam = append(messageParam, anthropic.NewAssistantMessage(anthropic.NewTextBlock(v.Assistant)))
		}
	}
	res, err = claude.cli.Messages.New(context.TODO(), anthropic.MessageNewParams{
		Model:     anthropic.F(req.Model),
		MaxTokens: anthropic.F(req.MaxToken),
		Messages:  anthropic.F(messageParam),
	})
	return
}
