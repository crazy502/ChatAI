package admin

import (
	"server/infra/metrics"
	"server/pkg/response"
)

type AllMetricsResponse struct {
	response.Response
	Snapshot metrics.AllMetricsSnapshot `json:"snapshot"`
}
