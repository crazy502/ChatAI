package ai

import (
	"context"
	"fmt"
	"io"
	"strings"

	"server/infra/config"

	"github.com/cloudwego/eino-ext/components/model/openai"
	einomodel "github.com/cloudwego/eino/components/model"
	"github.com/cloudwego/eino/schema"
)

type QwenProvider struct {
	llm einomodel.ToolCallingChatModel
}

func NewQwenProvider(ctx context.Context) (*QwenProvider, error) {
	cfg := config.GetConfig()
	modelName := cfg.QwenConfig.ModelName
	if modelName == "" {
		modelName = "qwen-plus"
	}
	baseURL := cfg.QwenConfig.BaseURL
	if baseURL == "" {
		baseURL = "https://dashscope.aliyuncs.com/compatible-mode/v1"
	}

	llm, err := openai.NewChatModel(ctx, &openai.ChatModelConfig{
		BaseURL: baseURL,
		Model:   modelName,
		APIKey:  cfg.QwenConfig.APIKey,
	})
	if err != nil {
		return nil, fmt.Errorf("create qwen model failed: %v", err)
	}

	return &QwenProvider{llm: llm}, nil
}

func (p *QwenProvider) GenerateResponse(ctx context.Context, messages []*schema.Message) (*schema.Message, error) {
	resp, err := p.llm.Generate(ctx, messages)
	if err != nil {
		return nil, fmt.Errorf("qwen generate failed: %v", err)
	}
	return resp, nil
}

// StreamResponse 流式响应
// ctx: 上下文
// messages: 消息队列
// cb: 流式回调函数
// return: 流式响应
// err: 错误
func (p *QwenProvider) StreamResponse(ctx context.Context, messages []*schema.Message, cb StreamCallback) (string, error) {
	//1. 创建流式响应
	stream, err := p.llm.Stream(ctx, messages)
	if err != nil {
		return "", fmt.Errorf("qwen stream failed: %v", err)
	}
	defer stream.Close()

	//2. 处理流式响应
	var fullResp strings.Builder

	for {
		msg, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return "", fmt.Errorf("qwen stream recv failed: %v", err)
		}
		if len(msg.Content) == 0 {
			continue
		}

		//3. 累加流式响应
		fullResp.WriteString(msg.Content)
		//4. 调用流式回调函数
		cb(msg.Content)
	}

	return fullResp.String(), nil
}

func (p *QwenProvider) Name() string {
	return "qwen"
}
