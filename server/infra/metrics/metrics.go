package metrics

import (
	"context"
	"fmt"
	"sort"
	"sync"
	"time"

	"server/pkg/apperror"
	"server/pkg/code"
	"server/pkg/requestmeta"
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

type ErrorEvent struct {
	Timestamp time.Time         `json:"timestamp"`
	Code      int64             `json:"code"`
	Message   string            `json:"message"`
	RequestID string            `json:"requestId,omitempty"`
	UserID    int64             `json:"userId,omitempty"`
	UserName  string            `json:"userName,omitempty"`
	SessionID string            `json:"sessionId,omitempty"`
	Stack     []apperror.Frame  `json:"stack,omitempty"`
	Fields    map[string]string `json:"fields,omitempty"`
}

type AlertSnapshot struct {
	Timestamp     time.Time `json:"timestamp"`
	Code          int64     `json:"code"`
	Message       string    `json:"message"`
	Occurrences   int       `json:"occurrences"`
	WindowSeconds int64     `json:"windowSeconds"`
}

type AllMetricsSnapshot struct {
	GeneratedAt    time.Time              `json:"generatedAt"`
	Overview       Overview               `json:"overview"`
	Routes         []RouteSnapshot        `json:"routes"`
	Models         []ModelSnapshot        `json:"models"`
	Archives       []OverviewArchivePoint `json:"archives"`
	RecentErrors   []ErrorEvent           `json:"recentErrors"`
	Alerts         []AlertSnapshot        `json:"alerts"`
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
	recentErrors         []ErrorEvent
	alerts               []AlertSnapshot
	errorWindows         map[string][]time.Time
	lastAlertAt          map[string]time.Time
	lastArchiveAt        time.Time
	lastCleanupAt        time.Time
	archiveSampleEvery   time.Duration
	retentionWindow      time.Duration
	alertWindow          time.Duration
	alertSuppressWindow  time.Duration
	alertThreshold       int
}

var (
	globalCollector *Collector
	once            sync.Once
)

const (
	maxOverviewArchives = 720
	maxRecentErrors     = 100
	maxAlerts           = 50
	cleanupInterval     = 5 * time.Minute
	archiveSampleEvery  = 30 * time.Second
	retentionWindow     = 6 * time.Hour
	alertWindow         = 2 * time.Minute
	alertSuppressWindow = 2 * time.Minute
	alertThreshold      = 3
)

func NewCollector() *Collector {
	return &Collector{
		startedAt:           time.Now(),
		routes:              make(map[string]*routeState),
		models:              make(map[string]*modelState),
		errorWindows:        make(map[string][]time.Time),
		lastAlertAt:         make(map[string]time.Time),
		archiveSampleEvery:  archiveSampleEvery,
		retentionWindow:     retentionWindow,
		alertWindow:         alertWindow,
		alertSuppressWindow: alertSuppressWindow,
		alertThreshold:      alertThreshold,
	}
}

func GetCollector() *Collector {
	once.Do(func() {
		globalCollector = NewCollector()
	})
	return globalCollector
}

func (c *Collector) cleanupLocked(now time.Time) {
	if !c.lastCleanupAt.IsZero() && now.Sub(c.lastCleanupAt) < cleanupInterval {
		return
	}

	staleBefore := now.Add(-c.retentionWindow)
	alertStaleBefore := now.Add(-c.alertWindow)

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

	filteredErrors := c.recentErrors[:0]
	for _, event := range c.recentErrors {
		if event.Timestamp.After(staleBefore) || event.Timestamp.Equal(staleBefore) {
			filteredErrors = append(filteredErrors, event)
		}
	}
	c.recentErrors = filteredErrors
	if len(c.recentErrors) > maxRecentErrors {
		c.recentErrors = c.recentErrors[len(c.recentErrors)-maxRecentErrors:]
	}

	filteredAlerts := c.alerts[:0]
	for _, alert := range c.alerts {
		if alert.Timestamp.After(staleBefore) || alert.Timestamp.Equal(staleBefore) {
			filteredAlerts = append(filteredAlerts, alert)
		}
	}
	c.alerts = filteredAlerts
	if len(c.alerts) > maxAlerts {
		c.alerts = c.alerts[len(c.alerts)-maxAlerts:]
	}

	for key, window := range c.errorWindows {
		filteredWindow := window[:0]
		for _, occurredAt := range window {
			if occurredAt.After(alertStaleBefore) || occurredAt.Equal(alertStaleBefore) {
				filteredWindow = append(filteredWindow, occurredAt)
			}
		}

		if len(filteredWindow) == 0 {
			delete(c.errorWindows, key)
			delete(c.lastAlertAt, key)
			continue
		}

		c.errorWindows[key] = filteredWindow
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

func (c *Collector) RecordError(ctx context.Context, err error) {
	appErr := apperror.From(err)
	if appErr == nil {
		return
	}

	now := time.Now()
	event := ErrorEvent{
		Timestamp: now,
		Code:      appErr.Code.Code(),
		Message:   appErr.Message,
		RequestID: requestmeta.String(ctx, requestmeta.FieldRequestID),
		UserID:    requestmeta.Int64(ctx, requestmeta.FieldUserID),
		UserName:  requestmeta.String(ctx, requestmeta.FieldUserName),
		SessionID: requestmeta.String(ctx, requestmeta.FieldSessionID),
		Stack:     apperror.StackOf(appErr),
		Fields:    stringifyFields(apperror.FieldsOf(appErr)),
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	c.recentErrors = append(c.recentErrors, event)
	if len(c.recentErrors) > maxRecentErrors {
		c.recentErrors = c.recentErrors[len(c.recentErrors)-maxRecentErrors:]
	}

	c.recordAlertLocked(now, appErr)
	c.cleanupLocked(now)
}

func (c *Collector) recordAlertLocked(now time.Time, err *apperror.Error) {
	if err == nil || !shouldAlert(err.Code) {
		return
	}

	key := fmt.Sprintf("%d", err.Code.Code())
	window := c.errorWindows[key]
	staleBefore := now.Add(-c.alertWindow)

	filteredWindow := window[:0]
	for _, occurredAt := range window {
		if occurredAt.After(staleBefore) || occurredAt.Equal(staleBefore) {
			filteredWindow = append(filteredWindow, occurredAt)
		}
	}
	filteredWindow = append(filteredWindow, now)
	c.errorWindows[key] = filteredWindow

	if len(filteredWindow) < c.alertThreshold {
		return
	}

	lastAlertAt := c.lastAlertAt[key]
	if !lastAlertAt.IsZero() && now.Sub(lastAlertAt) < c.alertSuppressWindow {
		return
	}

	c.alerts = append(c.alerts, AlertSnapshot{
		Timestamp:     now,
		Code:          err.Code.Code(),
		Message:       fmt.Sprintf("错误码 %d 在最近 %d 秒内出现 %d 次", err.Code.Code(), int64(c.alertWindow.Seconds()), len(filteredWindow)),
		Occurrences:   len(filteredWindow),
		WindowSeconds: int64(c.alertWindow.Seconds()),
	})
	c.lastAlertAt[key] = now

	if len(c.alerts) > maxAlerts {
		c.alerts = c.alerts[len(c.alerts)-maxAlerts:]
	}
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

func (c *Collector) RecentErrors() []ErrorEvent {
	c.mu.RLock()
	defer c.mu.RUnlock()

	snapshots := make([]ErrorEvent, len(c.recentErrors))
	copy(snapshots, c.recentErrors)

	sort.Slice(snapshots, func(i, j int) bool {
		return snapshots[i].Timestamp.After(snapshots[j].Timestamp)
	})

	return snapshots
}

func (c *Collector) AlertSnapshots() []AlertSnapshot {
	c.mu.RLock()
	defer c.mu.RUnlock()

	snapshots := make([]AlertSnapshot, len(c.alerts))
	copy(snapshots, c.alerts)

	sort.Slice(snapshots, func(i, j int) bool {
		return snapshots[i].Timestamp.After(snapshots[j].Timestamp)
	})

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
		RecentErrors:   c.RecentErrors(),
		Alerts:         c.AlertSnapshots(),
		ArchiveWindowS: int64(c.retentionWindow.Seconds()),
	}
}

func shouldAlert(resultCode code.Code) bool {
	return resultCode.Code() >= code.CodeServerBusy.Code()
}

func stringifyFields(fields map[string]any) map[string]string {
	if len(fields) == 0 {
		return nil
	}

	result := make(map[string]string, len(fields))
	for key, value := range fields {
		if key == "" || value == nil {
			continue
		}
		result[key] = fmt.Sprintf("%v", value)
	}
	if len(result) == 0 {
		return nil
	}
	return result
}
