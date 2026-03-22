package aihelper

import (
	"context"
	"fmt"
	"io"
	"strings"

	"server/config"

	"github.com/cloudwego/eino-ext/components/model/openai"
	"github.com/cloudwego/eino/components/model"
	"github.com/cloudwego/eino/schema"
)

type StreamCallback func(msg string)

// AIModel 瀹氫箟AI妯″瀷鎺ュ彛
type AIModel interface {
	GenerateResponse(ctx context.Context, messages []*schema.Message) (*schema.Message, error)
	StreamResponse(ctx context.Context, messages []*schema.Message, cb StreamCallback) (string, error)
	GetModelType() string
}

// =================== QWEN 瀹炵幇 ===================
type QWENModel struct {
	llm model.ToolCallingChatModel
}

func NewQWENModel(ctx context.Context) (*QWENModel, error) {
	cfg := config.GetConfig()
	key := cfg.QwenConfig.APIKey
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
		APIKey:  key,
	})
	if err != nil {
		return nil, fmt.Errorf("create qwen model failed: %v", err)
	}
	return &QWENModel{llm: llm}, nil
}

func (q *QWENModel) GenerateResponse(ctx context.Context, messages []*schema.Message) (*schema.Message, error) {
	resp, err := q.llm.Generate(ctx, messages)
	if err != nil {
		return nil, fmt.Errorf("qwen generate failed: %v", err)
	}
	return resp, nil
}

func (q *QWENModel) StreamResponse(ctx context.Context, messages []*schema.Message, cb StreamCallback) (string, error) {
	stream, err := q.llm.Stream(ctx, messages)
	if err != nil {
		return "", fmt.Errorf("qwen stream failed: %v", err)
	}
	defer stream.Close()

	var fullResp strings.Builder

	for {
		msg, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return "", fmt.Errorf("qwen stream recv failed: %v", err)
		}
		if len(msg.Content) > 0 {
			fullResp.WriteString(msg.Content)
			cb(msg.Content)
		}
	}

	return fullResp.String(), nil
}

func (q *QWENModel) GetModelType() string { return "qwen" }

// =================== DeepSeek 瀹炵幇 ===================
type DeepSeekModel struct {
	llm model.ToolCallingChatModel
}

func NewDeepSeekModel(ctx context.Context) (*DeepSeekModel, error) {
	cfg := config.GetConfig()
	key := cfg.DeepSeekConfig.APIKey
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
		APIKey:  key,
	})
	if err != nil {
		return nil, fmt.Errorf("create deepseek model failed: %v", err)
	}
	return &DeepSeekModel{llm: llm}, nil
}

func (d *DeepSeekModel) GenerateResponse(ctx context.Context, messages []*schema.Message) (*schema.Message, error) {
	resp, err := d.llm.Generate(ctx, messages)
	if err != nil {
		return nil, fmt.Errorf("deepseek generate failed: %v", err)
	}
	return resp, nil
}

func (d *DeepSeekModel) StreamResponse(ctx context.Context, messages []*schema.Message, cb StreamCallback) (string, error) {
	stream, err := d.llm.Stream(ctx, messages)
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
		if len(msg.Content) > 0 {
			fullResp.WriteString(msg.Content)
			cb(msg.Content)
		}
	}

	return fullResp.String(), nil
}

func (d *DeepSeekModel) GetModelType() string { return "deepseek" }
