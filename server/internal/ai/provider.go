package ai

import (
	"context"
	"fmt"
	"sync"

	"github.com/cloudwego/eino/schema"
)

type Provider interface {
	GenerateResponse(ctx context.Context, messages []*schema.Message) (*schema.Message, error)
	StreamResponse(ctx context.Context, messages []*schema.Message, cb StreamCallback) (string, error)
	Name() string
}

type ProviderCreator func(ctx context.Context, config map[string]interface{}) (Provider, error)

type Factory struct {
	creators map[string]ProviderCreator
}

var (
	globalFactory *Factory
	factoryOnce   sync.Once
)

func GetGlobalFactory() *Factory {
	factoryOnce.Do(func() {
		globalFactory = &Factory{
			creators: make(map[string]ProviderCreator),
		}
		globalFactory.registerCreators()
	})
	return globalFactory
}

func (f *Factory) registerCreators() {
	f.creators["qwen"] = func(ctx context.Context, config map[string]interface{}) (Provider, error) {
		return NewQwenProvider(ctx)
	}

	f.creators["deepseek"] = func(ctx context.Context, config map[string]interface{}) (Provider, error) {
		return NewDeepSeekProvider(ctx)
	}
}

func (f *Factory) CreateProvider(ctx context.Context, modelType string, config map[string]interface{}) (Provider, error) {
	creator, ok := f.creators[modelType]
	if !ok {
		return nil, fmt.Errorf("unsupported model type: %s", modelType)
	}
	return creator(ctx, config)
}

func (f *Factory) RegisterProvider(modelType string, creator ProviderCreator) {
	f.creators[modelType] = creator
}
