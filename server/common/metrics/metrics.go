package metrics

import (
	"sort"
	"sync"
	"time"

	"server/common/code"
)

type Overview struct {
	UptimeSeconds int64   `json:"uptimeSeconds"`
	RequestsTotal int64   `json:"requestsTotal"`
	ErrorsTotal   int64   `json:"errorsTotal"`
	ErrorRate     float64 `json:"errorRate"`
	AvgLatencyMs  float64 `json:"avgLatencyMs"`
	RoutesTracked int     `json:"routesTracked"`
	ModelsTracked int     `json:"modelsTracked"`
}

type RouteSnapshot struct {
	Method         string  `json:"method"`
	Path           string  `json:"path"`
	RequestsTotal  int64   `json:"requestsTotal"`
	ErrorsTotal    int64   `json:"errorsTotal"`
	ErrorRate      float64 `json:"errorRate"`
	AvgLatencyMs   float64 `json:"avgLatencyMs"`
	LastLatencyMs  float64 `json:"lastLatencyMs"`
	LastHTTPStatus int     `json:"lastHttpStatus"`
}

type ModelSnapshot struct {
	ModelType     string    `json:"modelType"`
	Operation     string    `json:"operation"`
	RequestsTotal int64     `json:"requestsTotal"`
	ErrorsTotal   int64     `json:"errorsTotal"`
	ErrorRate     float64   `json:"errorRate"`
	AvgLatencyMs  float64   `json:"avgLatencyMs"`
	LastLatencyMs float64   `json:"lastLatencyMs"`
	LastSuccessAt time.Time `json:"lastSuccessAt,omitempty"`
	LastFailureAt time.Time `json:"lastFailureAt,omitempty"`
}

type OverviewArchivePoint struct {
	Timestamp     time.Time `json:"timestamp"`
	RequestsTotal int64     `json:"requestsTotal"`
	ErrorsTotal   int64     `json:"errorsTotal"`
	ErrorRate     float64   `json:"errorRate"`
	AvgLatencyMs  float64   `json:"avgLatencyMs"`
}

type AllMetricsSnapshot struct {
	GeneratedAt    time.Time              `json:"generatedAt"`
	Overview       Overview               `json:"overview"`
	Routes         []RouteSnapshot        `json:"routes"`
	Models         []ModelSnapshot        `json:"models"`
	Archives       []OverviewArchivePoint `json:"archives"`
	ArchiveWindowS int64                  `json:"archiveWindowSeconds"`
}

type routeState struct {
	Method         string
	Path           string
	RequestsTotal  int64
	ErrorsTotal    int64
	TotalLatency   time.Duration
	LastLatency    time.Duration
	LastHTTPStatus int
	LastSeenAt     time.Time
}

type modelState struct {
	ModelType     string
	Operation     string
	RequestsTotal int64
	ErrorsTotal   int64
	TotalLatency  time.Duration
	LastLatency   time.Duration
	LastSuccessAt time.Time
	LastFailureAt time.Time
	LastSeenAt    time.Time
}

type Collector struct {
	mu                   sync.RWMutex
	startedAt            time.Time
	requestsTotal        int64
	errorsTotal          int64
	totalRequestDuration time.Duration
	routes               map[string]*routeState
	models               map[string]*modelState
	archives             []OverviewArchivePoint
	lastArchiveAt        time.Time
	lastCleanupAt        time.Time
	archiveSampleEvery   time.Duration
	retentionWindow      time.Duration
}

var (
	globalCollector *Collector
	once            sync.Once
)

const (
	maxOverviewArchives = 720
	cleanupInterval     = 5 * time.Minute
	archiveSampleEvery  = 30 * time.Second
	retentionWindow     = 6 * time.Hour
)

func GetCollector() *Collector {
	once.Do(func() {
		globalCollector = &Collector{
			startedAt:          time.Now(),
			routes:             make(map[string]*routeState),
			models:             make(map[string]*modelState),
			archiveSampleEvery: archiveSampleEvery,
			retentionWindow:    retentionWindow,
		}
	})
	return globalCollector
}

func (c *Collector) cleanupLocked(now time.Time) {
	if !c.lastCleanupAt.IsZero() && now.Sub(c.lastCleanupAt) < cleanupInterval {
		return
	}

	staleBefore := now.Add(-c.retentionWindow)

	for key, routeMetric := range c.routes {
		if !routeMetric.LastSeenAt.IsZero() && routeMetric.LastSeenAt.Before(staleBefore) {
			delete(c.routes, key)
		}
	}

	for key, modelMetric := range c.models {
		if !modelMetric.LastSeenAt.IsZero() && modelMetric.LastSeenAt.Before(staleBefore) {
			delete(c.models, key)
		}
	}

	filteredArchives := c.archives[:0]
	for _, archive := range c.archives {
		if archive.Timestamp.After(staleBefore) || archive.Timestamp.Equal(staleBefore) {
			filteredArchives = append(filteredArchives, archive)
		}
	}
	c.archives = filteredArchives
	if len(c.archives) > maxOverviewArchives {
		c.archives = c.archives[len(c.archives)-maxOverviewArchives:]
	}

	c.lastCleanupAt = now
}

func (c *Collector) appendArchiveLocked(now time.Time) {
	if !c.lastArchiveAt.IsZero() && now.Sub(c.lastArchiveAt) < c.archiveSampleEvery {
		return
	}

	avgLatency := 0.0
	if c.requestsTotal > 0 {
		avgLatency = float64(c.totalRequestDuration.Milliseconds()) / float64(c.requestsTotal)
	}

	errorRate := 0.0
	if c.requestsTotal > 0 {
		errorRate = float64(c.errorsTotal) / float64(c.requestsTotal)
	}

	c.archives = append(c.archives, OverviewArchivePoint{
		Timestamp:     now,
		RequestsTotal: c.requestsTotal,
		ErrorsTotal:   c.errorsTotal,
		ErrorRate:     errorRate,
		AvgLatencyMs:  avgLatency,
	})
	c.lastArchiveAt = now

	if len(c.archives) > maxOverviewArchives {
		c.archives = c.archives[len(c.archives)-maxOverviewArchives:]
	}
}

func (c *Collector) RecordRequest(method, path, userName string, latency time.Duration, businessCode int64, httpStatus int) {
	if path == "" {
		path = "unknown"
	}

	isError := httpStatus >= 400
	if businessCode != 0 && businessCode != int64(code.CodeSuccess) {
		isError = true
	}
	now := time.Now()

	c.mu.Lock()
	defer c.mu.Unlock()

	c.requestsTotal++
	c.totalRequestDuration += latency
	if isError {
		c.errorsTotal++
	}

	routeKey := method + " " + path
	routeMetric, exists := c.routes[routeKey]
	if !exists {
		routeMetric = &routeState{
			Method: method,
			Path:   path,
		}
		c.routes[routeKey] = routeMetric
	}

	routeMetric.RequestsTotal++
	routeMetric.TotalLatency += latency
	routeMetric.LastLatency = latency
	routeMetric.LastHTTPStatus = httpStatus
	routeMetric.LastSeenAt = now
	if isError {
		routeMetric.ErrorsTotal++
	}

	c.appendArchiveLocked(now)
	c.cleanupLocked(now)
}

func (c *Collector) RecordModel(modelType, operation, userName string, latency time.Duration, err error) {
	if modelType == "" {
		modelType = "unknown"
	}
	if operation == "" {
		operation = "unknown"
	}

	now := time.Now()

	c.mu.Lock()
	defer c.mu.Unlock()

	key := modelType + ":" + operation
	modelMetric, exists := c.models[key]
	if !exists {
		modelMetric = &modelState{
			ModelType: modelType,
			Operation: operation,
		}
		c.models[key] = modelMetric
	}

	modelMetric.RequestsTotal++
	modelMetric.TotalLatency += latency
	modelMetric.LastLatency = latency
	modelMetric.LastSeenAt = now

	if err == nil {
		modelMetric.LastSuccessAt = now
		c.cleanupLocked(now)
		return
	}

	modelMetric.ErrorsTotal++
	modelMetric.LastFailureAt = now

	c.cleanupLocked(now)
}

func (c *Collector) Overview() Overview {
	c.mu.RLock()
	defer c.mu.RUnlock()

	avgLatency := 0.0
	if c.requestsTotal > 0 {
		avgLatency = float64(c.totalRequestDuration.Milliseconds()) / float64(c.requestsTotal)
	}

	errorRate := 0.0
	if c.requestsTotal > 0 {
		errorRate = float64(c.errorsTotal) / float64(c.requestsTotal)
	}

	return Overview{
		UptimeSeconds: int64(time.Since(c.startedAt).Seconds()),
		RequestsTotal: c.requestsTotal,
		ErrorsTotal:   c.errorsTotal,
		ErrorRate:     errorRate,
		AvgLatencyMs:  avgLatency,
		RoutesTracked: len(c.routes),
		ModelsTracked: len(c.models),
	}
}

func (c *Collector) RouteSnapshots() []RouteSnapshot {
	c.mu.RLock()
	defer c.mu.RUnlock()

	snapshots := make([]RouteSnapshot, 0, len(c.routes))
	for _, routeMetric := range c.routes {
		errorRate := 0.0
		avgLatency := 0.0
		if routeMetric.RequestsTotal > 0 {
			errorRate = float64(routeMetric.ErrorsTotal) / float64(routeMetric.RequestsTotal)
			avgLatency = float64(routeMetric.TotalLatency.Milliseconds()) / float64(routeMetric.RequestsTotal)
		}

		snapshots = append(snapshots, RouteSnapshot{
			Method:         routeMetric.Method,
			Path:           routeMetric.Path,
			RequestsTotal:  routeMetric.RequestsTotal,
			ErrorsTotal:    routeMetric.ErrorsTotal,
			ErrorRate:      errorRate,
			AvgLatencyMs:   avgLatency,
			LastLatencyMs:  float64(routeMetric.LastLatency.Milliseconds()),
			LastHTTPStatus: routeMetric.LastHTTPStatus,
		})
	}

	sort.Slice(snapshots, func(i, j int) bool {
		if snapshots[i].RequestsTotal == snapshots[j].RequestsTotal {
			return snapshots[i].Path < snapshots[j].Path
		}
		return snapshots[i].RequestsTotal > snapshots[j].RequestsTotal
	})

	return snapshots
}

func (c *Collector) ModelSnapshots() []ModelSnapshot {
	c.mu.RLock()
	defer c.mu.RUnlock()

	snapshots := make([]ModelSnapshot, 0, len(c.models))
	for _, modelMetric := range c.models {
		errorRate := 0.0
		avgLatency := 0.0
		if modelMetric.RequestsTotal > 0 {
			errorRate = float64(modelMetric.ErrorsTotal) / float64(modelMetric.RequestsTotal)
			avgLatency = float64(modelMetric.TotalLatency.Milliseconds()) / float64(modelMetric.RequestsTotal)
		}

		snapshots = append(snapshots, ModelSnapshot{
			ModelType:     modelMetric.ModelType,
			Operation:     modelMetric.Operation,
			RequestsTotal: modelMetric.RequestsTotal,
			ErrorsTotal:   modelMetric.ErrorsTotal,
			ErrorRate:     errorRate,
			AvgLatencyMs:  avgLatency,
			LastLatencyMs: float64(modelMetric.LastLatency.Milliseconds()),
			LastSuccessAt: modelMetric.LastSuccessAt,
			LastFailureAt: modelMetric.LastFailureAt,
		})
	}

	sort.Slice(snapshots, func(i, j int) bool {
		if snapshots[i].RequestsTotal == snapshots[j].RequestsTotal {
			if snapshots[i].ModelType == snapshots[j].ModelType {
				return snapshots[i].Operation < snapshots[j].Operation
			}
			return snapshots[i].ModelType < snapshots[j].ModelType
		}
		return snapshots[i].RequestsTotal > snapshots[j].RequestsTotal
	})

	return snapshots
}

func (c *Collector) ArchiveSnapshots() []OverviewArchivePoint {
	c.mu.RLock()
	defer c.mu.RUnlock()

	snapshots := make([]OverviewArchivePoint, len(c.archives))
	copy(snapshots, c.archives)
	return snapshots
}

func (c *Collector) AllMetricsSnapshot() AllMetricsSnapshot {
	c.mu.Lock()
	c.cleanupLocked(time.Now())
	c.mu.Unlock()

	return AllMetricsSnapshot{
		GeneratedAt:    time.Now(),
		Overview:       c.Overview(),
		Routes:         c.RouteSnapshots(),
		Models:         c.ModelSnapshots(),
		Archives:       c.ArchiveSnapshots(),
		ArchiveWindowS: int64(c.retentionWindow.Seconds()),
	}
}
