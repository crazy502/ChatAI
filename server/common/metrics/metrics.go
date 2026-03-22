package metrics

import (
	"fmt"
	"sort"
	"strings"
	"sync"
	"time"

	"server/common/code"
)

type Overview struct {
	UptimeSeconds         int64   `json:"uptimeSeconds"`
	RequestsTotal         int64   `json:"requestsTotal"`
	ErrorsTotal           int64   `json:"errorsTotal"`
	ErrorRate             float64 `json:"errorRate"`
	AvgLatencyMs          float64 `json:"avgLatencyMs"`
	RoutesTracked         int     `json:"routesTracked"`
	ModelsTracked         int     `json:"modelsTracked"`
	UsersTracked          int     `json:"usersTracked"`
	BusinessCodesTracked  int     `json:"businessCodesTracked"`
	RecentFailuresTracked int     `json:"recentFailuresTracked"`
}

type RouteSnapshot struct {
	Method           string  `json:"method"`
	Path             string  `json:"path"`
	RequestsTotal    int64   `json:"requestsTotal"`
	ErrorsTotal      int64   `json:"errorsTotal"`
	ErrorRate        float64 `json:"errorRate"`
	AvgLatencyMs     float64 `json:"avgLatencyMs"`
	LastLatencyMs    float64 `json:"lastLatencyMs"`
	LastBusinessCode int64   `json:"lastBusinessCode"`
	LastHTTPStatus   int     `json:"lastHttpStatus"`
	LastUserName     string  `json:"lastUserName"`
}

type ModelSnapshot struct {
	ModelType        string    `json:"modelType"`
	Operation        string    `json:"operation"`
	RequestsTotal    int64     `json:"requestsTotal"`
	ErrorsTotal      int64     `json:"errorsTotal"`
	ErrorRate        float64   `json:"errorRate"`
	AvgLatencyMs     float64   `json:"avgLatencyMs"`
	LastLatencyMs    float64   `json:"lastLatencyMs"`
	LastSuccessAt    time.Time `json:"lastSuccessAt,omitempty"`
	LastFailureAt    time.Time `json:"lastFailureAt,omitempty"`
	LastUserName     string    `json:"lastUserName,omitempty"`
	LastErrorMessage string    `json:"lastErrorMessage,omitempty"`
}

type UserSnapshot struct {
	UserName      string    `json:"userName"`
	RequestsTotal int64     `json:"requestsTotal"`
	ErrorsTotal   int64     `json:"errorsTotal"`
	ErrorRate     float64   `json:"errorRate"`
	AvgLatencyMs  float64   `json:"avgLatencyMs"`
	LastSeenAt    time.Time `json:"lastSeenAt,omitempty"`
}

type BusinessCodeSnapshot struct {
	Method         string    `json:"method"`
	Path           string    `json:"path"`
	BusinessCode   int64     `json:"businessCode"`
	RequestsTotal  int64     `json:"requestsTotal"`
	ErrorsTotal    int64     `json:"errorsTotal"`
	ErrorRate      float64   `json:"errorRate"`
	AvgLatencyMs   float64   `json:"avgLatencyMs"`
	LastHTTPStatus int       `json:"lastHttpStatus"`
	LastSeenAt     time.Time `json:"lastSeenAt,omitempty"`
}

type ModelFailureSnapshot struct {
	ModelType    string    `json:"modelType"`
	Operation    string    `json:"operation"`
	UserName     string    `json:"userName"`
	ErrorMessage string    `json:"errorMessage"`
	LatencyMs    float64   `json:"latencyMs"`
	OccurredAt   time.Time `json:"occurredAt"`
}

type routeState struct {
	Method           string
	Path             string
	RequestsTotal    int64
	ErrorsTotal      int64
	TotalLatency     time.Duration
	LastLatency      time.Duration
	LastBusinessCode int64
	LastHTTPStatus   int
	LastUserName     string
}

type modelState struct {
	ModelType        string
	Operation        string
	RequestsTotal    int64
	ErrorsTotal      int64
	TotalLatency     time.Duration
	LastLatency      time.Duration
	LastSuccessAt    time.Time
	LastFailureAt    time.Time
	LastUserName     string
	LastErrorMessage string
}

type userState struct {
	UserName      string
	RequestsTotal int64
	ErrorsTotal   int64
	TotalLatency  time.Duration
	LastSeenAt    time.Time
}

type businessCodeState struct {
	Method         string
	Path           string
	BusinessCode   int64
	RequestsTotal  int64
	ErrorsTotal    int64
	TotalLatency   time.Duration
	LastHTTPStatus int
	LastSeenAt     time.Time
}

type modelFailureState struct {
	ModelType    string
	Operation    string
	UserName     string
	ErrorMessage string
	Latency      time.Duration
	OccurredAt   time.Time
}

type Collector struct {
	mu                   sync.RWMutex
	startedAt            time.Time
	requestsTotal        int64
	errorsTotal          int64
	totalRequestDuration time.Duration
	routes               map[string]*routeState
	models               map[string]*modelState
	users                map[string]*userState
	businessCodes        map[string]*businessCodeState
	modelFailures        []modelFailureState
}

var (
	globalCollector *Collector
	once            sync.Once
)

const (
	maxRecentModelFailures = 80
	maxErrorMessageLength  = 240
)

func GetCollector() *Collector {
	once.Do(func() {
		globalCollector = &Collector{
			startedAt:     time.Now(),
			routes:        make(map[string]*routeState),
			models:        make(map[string]*modelState),
			users:         make(map[string]*userState),
			businessCodes: make(map[string]*businessCodeState),
		}
	})
	return globalCollector
}

func normalizeUserName(userName string) string {
	value := strings.TrimSpace(userName)
	if value == "" {
		return "anonymous"
	}
	return value
}

func trimErrorMessage(message string) string {
	value := strings.TrimSpace(strings.ReplaceAll(strings.ReplaceAll(message, "\n", " "), "\r", " "))
	if value == "" {
		return "unknown error"
	}

	runes := []rune(value)
	if len(runes) <= maxErrorMessageLength {
		return value
	}

	return string(runes[:maxErrorMessageLength]) + "..."
}

func (c *Collector) RecordRequest(method, path, userName string, latency time.Duration, businessCode int64, httpStatus int) {
	if path == "" {
		path = "unknown"
	}

	isError := httpStatus >= 400
	if businessCode != 0 && businessCode != int64(code.CodeSuccess) {
		isError = true
	}

	normalizedUserName := normalizeUserName(userName)
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
	routeMetric.LastBusinessCode = businessCode
	routeMetric.LastHTTPStatus = httpStatus
	routeMetric.LastUserName = normalizedUserName
	if isError {
		routeMetric.ErrorsTotal++
	}

	userMetric, exists := c.users[normalizedUserName]
	if !exists {
		userMetric = &userState{
			UserName: normalizedUserName,
		}
		c.users[normalizedUserName] = userMetric
	}

	userMetric.RequestsTotal++
	userMetric.TotalLatency += latency
	userMetric.LastSeenAt = now
	if isError {
		userMetric.ErrorsTotal++
	}

	codeKey := fmt.Sprintf("%s:%d", routeKey, businessCode)
	codeMetric, exists := c.businessCodes[codeKey]
	if !exists {
		codeMetric = &businessCodeState{
			Method:       method,
			Path:         path,
			BusinessCode: businessCode,
		}
		c.businessCodes[codeKey] = codeMetric
	}

	codeMetric.RequestsTotal++
	codeMetric.TotalLatency += latency
	codeMetric.LastHTTPStatus = httpStatus
	codeMetric.LastSeenAt = now
	if isError {
		codeMetric.ErrorsTotal++
	}
}

func (c *Collector) RecordModel(modelType, operation, userName string, latency time.Duration, err error) {
	if modelType == "" {
		modelType = "unknown"
	}
	if operation == "" {
		operation = "unknown"
	}

	normalizedUserName := normalizeUserName(userName)
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
	modelMetric.LastUserName = normalizedUserName

	if err == nil {
		modelMetric.LastSuccessAt = now
		return
	}

	errorMessage := trimErrorMessage(err.Error())
	modelMetric.ErrorsTotal++
	modelMetric.LastFailureAt = now
	modelMetric.LastErrorMessage = errorMessage

	c.modelFailures = append(c.modelFailures, modelFailureState{
		ModelType:    modelType,
		Operation:    operation,
		UserName:     normalizedUserName,
		ErrorMessage: errorMessage,
		Latency:      latency,
		OccurredAt:   now,
	})
	if len(c.modelFailures) > maxRecentModelFailures {
		c.modelFailures = c.modelFailures[len(c.modelFailures)-maxRecentModelFailures:]
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
		UptimeSeconds:         int64(time.Since(c.startedAt).Seconds()),
		RequestsTotal:         c.requestsTotal,
		ErrorsTotal:           c.errorsTotal,
		ErrorRate:             errorRate,
		AvgLatencyMs:          avgLatency,
		RoutesTracked:         len(c.routes),
		ModelsTracked:         len(c.models),
		UsersTracked:          len(c.users),
		BusinessCodesTracked:  len(c.businessCodes),
		RecentFailuresTracked: len(c.modelFailures),
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
			Method:           routeMetric.Method,
			Path:             routeMetric.Path,
			RequestsTotal:    routeMetric.RequestsTotal,
			ErrorsTotal:      routeMetric.ErrorsTotal,
			ErrorRate:        errorRate,
			AvgLatencyMs:     avgLatency,
			LastLatencyMs:    float64(routeMetric.LastLatency.Milliseconds()),
			LastBusinessCode: routeMetric.LastBusinessCode,
			LastHTTPStatus:   routeMetric.LastHTTPStatus,
			LastUserName:     routeMetric.LastUserName,
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
			ModelType:        modelMetric.ModelType,
			Operation:        modelMetric.Operation,
			RequestsTotal:    modelMetric.RequestsTotal,
			ErrorsTotal:      modelMetric.ErrorsTotal,
			ErrorRate:        errorRate,
			AvgLatencyMs:     avgLatency,
			LastLatencyMs:    float64(modelMetric.LastLatency.Milliseconds()),
			LastSuccessAt:    modelMetric.LastSuccessAt,
			LastFailureAt:    modelMetric.LastFailureAt,
			LastUserName:     modelMetric.LastUserName,
			LastErrorMessage: modelMetric.LastErrorMessage,
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

func (c *Collector) UserSnapshots() []UserSnapshot {
	c.mu.RLock()
	defer c.mu.RUnlock()

	snapshots := make([]UserSnapshot, 0, len(c.users))
	for _, userMetric := range c.users {
		errorRate := 0.0
		avgLatency := 0.0
		if userMetric.RequestsTotal > 0 {
			errorRate = float64(userMetric.ErrorsTotal) / float64(userMetric.RequestsTotal)
			avgLatency = float64(userMetric.TotalLatency.Milliseconds()) / float64(userMetric.RequestsTotal)
		}

		snapshots = append(snapshots, UserSnapshot{
			UserName:      userMetric.UserName,
			RequestsTotal: userMetric.RequestsTotal,
			ErrorsTotal:   userMetric.ErrorsTotal,
			ErrorRate:     errorRate,
			AvgLatencyMs:  avgLatency,
			LastSeenAt:    userMetric.LastSeenAt,
		})
	}

	sort.Slice(snapshots, func(i, j int) bool {
		if snapshots[i].RequestsTotal == snapshots[j].RequestsTotal {
			return snapshots[i].UserName < snapshots[j].UserName
		}
		return snapshots[i].RequestsTotal > snapshots[j].RequestsTotal
	})

	return snapshots
}

func (c *Collector) BusinessCodeSnapshots() []BusinessCodeSnapshot {
	c.mu.RLock()
	defer c.mu.RUnlock()

	snapshots := make([]BusinessCodeSnapshot, 0, len(c.businessCodes))
	for _, codeMetric := range c.businessCodes {
		errorRate := 0.0
		avgLatency := 0.0
		if codeMetric.RequestsTotal > 0 {
			errorRate = float64(codeMetric.ErrorsTotal) / float64(codeMetric.RequestsTotal)
			avgLatency = float64(codeMetric.TotalLatency.Milliseconds()) / float64(codeMetric.RequestsTotal)
		}

		snapshots = append(snapshots, BusinessCodeSnapshot{
			Method:         codeMetric.Method,
			Path:           codeMetric.Path,
			BusinessCode:   codeMetric.BusinessCode,
			RequestsTotal:  codeMetric.RequestsTotal,
			ErrorsTotal:    codeMetric.ErrorsTotal,
			ErrorRate:      errorRate,
			AvgLatencyMs:   avgLatency,
			LastHTTPStatus: codeMetric.LastHTTPStatus,
			LastSeenAt:     codeMetric.LastSeenAt,
		})
	}

	sort.Slice(snapshots, func(i, j int) bool {
		if snapshots[i].ErrorsTotal == snapshots[j].ErrorsTotal {
			if snapshots[i].RequestsTotal == snapshots[j].RequestsTotal {
				if snapshots[i].Path == snapshots[j].Path {
					return snapshots[i].BusinessCode < snapshots[j].BusinessCode
				}
				return snapshots[i].Path < snapshots[j].Path
			}
			return snapshots[i].RequestsTotal > snapshots[j].RequestsTotal
		}
		return snapshots[i].ErrorsTotal > snapshots[j].ErrorsTotal
	})

	return snapshots
}

func (c *Collector) ModelFailureSnapshots() []ModelFailureSnapshot {
	c.mu.RLock()
	defer c.mu.RUnlock()

	snapshots := make([]ModelFailureSnapshot, 0, len(c.modelFailures))
	for i := len(c.modelFailures) - 1; i >= 0; i-- {
		failure := c.modelFailures[i]
		snapshots = append(snapshots, ModelFailureSnapshot{
			ModelType:    failure.ModelType,
			Operation:    failure.Operation,
			UserName:     failure.UserName,
			ErrorMessage: failure.ErrorMessage,
			LatencyMs:    float64(failure.Latency.Milliseconds()),
			OccurredAt:   failure.OccurredAt,
		})
	}

	return snapshots
}
