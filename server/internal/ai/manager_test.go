package ai

import (
	"context"
	"errors"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/cloudwego/eino/schema"
)

type serialTestProvider struct {
	mu        sync.Mutex
	active    int
	maxActive int
}

func (p *serialTestProvider) GenerateResponse(_ context.Context, messages []*schema.Message) (*schema.Message, error) {
	p.mu.Lock()
	p.active++
	if p.active > p.maxActive {
		p.maxActive = p.active
	}
	p.mu.Unlock()

	time.Sleep(20 * time.Millisecond)
	content := messages[len(messages)-1].Content

	p.mu.Lock()
	p.active--
	p.mu.Unlock()
	return &schema.Message{Role: schema.Assistant, Content: "reply:" + content}, nil
}

func (p *serialTestProvider) StreamResponse(context.Context, []*schema.Message, StreamCallback) (string, error) {
	return "", errors.New("not implemented")
}

func (p *serialTestProvider) Name() string {
	return "test"
}

func TestHelperSerializesConversationTurns(t *testing.T) {
	provider := new(serialTestProvider)
	helper := NewHelper(provider, "session-1")
	helper.SetSaveFunc(func(*StoredMessage) error { return nil })

	start := make(chan struct{})
	errorsCh := make(chan error, 2)
	var waitGroup sync.WaitGroup
	for _, question := range []string{"first", "second"} {
		waitGroup.Add(1)
		go func(value string) {
			defer waitGroup.Done()
			<-start
			_, err := helper.GenerateResponse("user@qq.com", context.Background(), value)
			errorsCh <- err
		}(question)
	}
	close(start)
	waitGroup.Wait()
	close(errorsCh)

	for err := range errorsCh {
		if err != nil {
			t.Fatalf("generate response: %v", err)
		}
	}
	if provider.maxActive != 1 {
		t.Fatalf("conversation turns overlapped: max active calls = %d", provider.maxActive)
	}
}

func TestAddMessageDoesNotMutateContextWhenPersistenceFails(t *testing.T) {
	helper := NewHelper(new(serialTestProvider), "session-1")
	helper.SetSaveFunc(func(*StoredMessage) error { return errors.New("database unavailable") })

	if _, err := helper.AddMessage("question", "user@qq.com", true, true); err == nil {
		t.Fatal("expected persistence error")
	}
	if helper.HasMessages() {
		t.Fatal("failed message must not remain in model context")
	}
}

func TestEnsureHydratedLoadsHistoryOnce(t *testing.T) {
	helper := NewHelper(new(serialTestProvider), "session-1")
	var loads atomic.Int32
	load := func() ([]StoredMessage, error) {
		loads.Add(1)
		time.Sleep(10 * time.Millisecond)
		return []StoredMessage{{Content: "history", IsUser: true}}, nil
	}

	var waitGroup sync.WaitGroup
	for index := 0; index < 5; index++ {
		waitGroup.Add(1)
		go func() {
			defer waitGroup.Done()
			if err := helper.EnsureHydrated(load); err != nil {
				t.Errorf("hydrate helper: %v", err)
			}
		}()
	}
	waitGroup.Wait()

	if got := loads.Load(); got != 1 {
		t.Fatalf("unexpected history load count: got %d want 1", got)
	}
}
