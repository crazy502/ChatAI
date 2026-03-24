package admin

import "server/infra/metrics"

type Service struct{}

func NewService() *Service {
	return &Service{}
}

func (s *Service) AllMetricsSnapshot() metrics.AllMetricsSnapshot {
	return metrics.GetCollector().AllMetricsSnapshot()
}
