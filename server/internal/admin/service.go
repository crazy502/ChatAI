package admin

import (
	"context"

	"server/infra/metrics"
)

type Service struct{}

func NewService() *Service {
	return &Service{}
}

func (s *Service) AllMetricsSnapshot(ctx context.Context) metrics.AllMetricsSnapshot {
	return metrics.GetCollector().AllMetricsSnapshot()
}
