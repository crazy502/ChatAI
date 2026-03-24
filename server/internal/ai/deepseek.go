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

type DeepSeekProvider struct {
	llm einomodel.ToolCallingChatModel
}

func NewDeepSeekProvider(ctx context.Context) (*DeepSeekProvider, error) {
	cfg := config.GetConfig()
	modelName := cfg.DeepSeekConfig.ModelName
	if modelName == "" {
		modelName = "deepseek-chat"
	}
	baseURL := cfg.DeepSeekConfig.BaseURL
	if baseURL == "" {
		baseURL = "https://api.deepseek.com"
	}

	llm, err := openai.NewChatModel(ctx, &openai.ChatModelConfig{
		BaseURL: baseURL,
		Model:   modelName,
		APIKey:  cfg.DeepSeekConfig.APIKey,
	})
	if err != nil {
		return nil, fmt.Errorf("create deepseek model failed: %v", err)
	}

	return &DeepSeekProvider{llm: llm}, nil
}

func (p *DeepSeekProvider) GenerateResponse(ctx context.Context, messages []*schema.Message) (*schema.Message, error) {
	resp, err := p.llm.Generate(ctx, messages)
	if err != nil {
		return nil, fmt.Errorf("deepseek generate failed: %v", err)
	}
	return resp, nil
}

func (p *DeepSeekProvider) StreamResponse(ctx context.Context, messages []*schema.Message, cb StreamCallback) (string, error) {
	stream, err := p.llm.Stream(ctx, messages)
	if err != nil {
		return "", fmt.Errorf("deepseek stream failed: %v", err)
	}
	defer stream.Close()

	var fullResp strings.Builder

	for {
		msg, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return "", fmt.Errorf("deepseek stream recv failed: %v", err)
		}
		if len(msg.Content) == 0 {
			continue
		}

		fullResp.WriteString(msg.Content)
		cb(msg.Content)
	}

	return fullResp.String(), nil
}

func (p *DeepSeekProvider) Name() string {
	return "deepseek"
}
