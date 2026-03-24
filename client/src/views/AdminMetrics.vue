<template>
  <div class="admin-shell">
    <div class="admin-bg" aria-hidden="true">
      <div class="glow glow-a"></div>
      <div class="glow glow-b"></div>
      <div class="grid"></div>
    </div>

    <header class="admin-header">
      <div class="header-brand">
        <div class="brand-badge">OPS</div>
        <div>
          <p class="brand-kicker">LIGHTWEIGHT MONITORING</p>
          <h1 class="brand-title">AgentGo 管理概览</h1>
          <p class="brand-subtitle">只保留请求、错误、延迟、接口健康和模型健康这几项必要信息。</p>
        </div>
      </div>

      <section class="header-status">
        <p class="status-kicker">GLOBAL STATUS</p>
        <div class="status-summary">
          <article class="status-item">
            <span class="status-label">当前状态</span>
            <strong class="status-value" :class="healthStatus.tone">{{ healthStatus.label }}</strong>
          </article>
          <article class="status-item">
            <span class="status-label">最近刷新</span>
            <strong class="status-value">{{ lastUpdatedLabel }}</strong>
          </article>
          <article class="status-item">
            <span class="status-label">自动刷新</span>
            <strong class="status-value">{{ autoRefresh ? refreshIntervalLabel : '手动刷新' }}</strong>
          </article>
        </div>
      </section>

      <section class="header-actions">
        <div class="refresh-controls">
          <label class="refresh-toggle">
            <input v-model="autoRefresh" type="checkbox" />
            <span>自动刷新</span>
          </label>
          <select v-model.number="refreshIntervalMs" class="refresh-select" :disabled="!autoRefresh">
            <option v-for="option in refreshOptions" :key="option.value" :value="option.value">
              {{ option.label }}
            </option>
          </select>
        </div>
        <div class="action-row">
          <button class="header-btn weak" type="button" @click="goMenu">返回控制台</button>
          <button class="header-btn secondary" type="button" @click="goChat">进入对话</button>
          <button class="header-btn primary" type="button" :disabled="refreshing" @click="fetchMetrics()">
            {{ refreshing ? '刷新中...' : '立即刷新' }}
          </button>
          <button class="header-btn danger" type="button" @click="logout">退出登录</button>
        </div>
      </section>
    </header>

    <main class="admin-main">
      <section class="hero-panel panel-card primary-card">
        <div class="hero-main">
          <div class="hero-copy">
            <p class="section-kicker">SYSTEM SNAPSHOT</p>
            <h2 class="hero-title">先判断系统状态，再决定是否继续下钻排查</h2>
            <p class="hero-desc">
              这个区域只承担一件事：给出当前最值得相信的整体判断，避免一上来就淹没在细节里。
            </p>
            <div class="diagnosis-box">
              <span class="diagnosis-label">当前结论</span>
              <strong class="diagnosis-title">{{ snapshotTitle }}</strong>
              <p class="diagnosis-text">{{ snapshotSummary }}</p>
            </div>
          </div>

          <div class="hero-side">
            <div class="health-pill" :class="healthStatus.tone">{{ healthStatus.label }}</div>
            <div class="hero-stats">
              <article class="hero-stat">
                <span class="hero-stat-label">最近刷新</span>
                <strong class="hero-stat-value">{{ lastUpdatedLabel }}</strong>
              </article>
              <article class="hero-stat">
                <span class="hero-stat-label">异常接口数</span>
                <strong class="hero-stat-value">{{ abnormalRouteCount }}</strong>
              </article>
              <article class="hero-stat">
                <span class="hero-stat-label">异常模型数</span>
                <strong class="hero-stat-value">{{ abnormalModelCount }}</strong>
              </article>
              <article class="hero-stat">
                <span class="hero-stat-label">刷新策略</span>
                <strong class="hero-stat-value">{{ autoRefresh ? refreshIntervalLabel : '手动刷新' }}</strong>
              </article>
            </div>
          </div>
        </div>

        <div class="hero-status-strip" :class="healthStatus.tone">
          <span class="strip-dot"></span>
          <span>{{ statusStripText }}</span>
        </div>
      </section>

      <section class="overview-grid">
        <article v-for="card in overviewCards" :key="card.label" class="metric-card panel-card" :class="card.tone">
          <span class="metric-label">{{ card.label }}</span>
          <strong class="metric-value">{{ card.value }}</strong>
          <p class="metric-desc">{{ card.desc }}</p>
        </article>
      </section>

      <section class="trend-grid">
        <article class="panel-card secondary-card">
          <div class="panel-head">
            <div>
              <p class="section-kicker">REQUEST TREND</p>
              <h3 class="panel-title">请求增量趋势</h3>
            </div>
            <span class="panel-note">绿色折线 / 低透明度面积</span>
          </div>

          <div v-if="requestTrendPoints.length" class="chart-shell">
            <svg viewBox="0 0 360 180" class="trend-chart" role="img" aria-label="请求趋势图">
              <g class="chart-grid">
                <line v-for="y in chartGuideLines" :key="`req-${y}`" x1="18" :y1="y" x2="342" :y2="y"></line>
              </g>
              <path class="chart-area request-area" :d="requestTrendAreaPath"></path>
              <path class="chart-line request-line" :d="requestTrendLinePath"></path>
              <circle v-for="point in requestTrendPoints" :key="point.key" class="chart-point request-point" :cx="point.x" :cy="point.y" r="3.5"></circle>
            </svg>
          </div>
          <div v-else class="panel-empty">等待更多采样数据后显示趋势。</div>

          <div class="trend-summary">
            <span>当前增量 {{ formatNumber(latestRequestDelta) }}</span>
            <span>峰值 {{ formatNumber(requestPeak) }}</span>
            <span>样本 {{ historySamples.length }} 个</span>
          </div>
        </article>

        <article class="panel-card secondary-card">
          <div class="panel-head">
            <div>
              <p class="section-kicker">LATENCY TREND</p>
              <h3 class="panel-title">平均延迟趋势</h3>
            </div>
            <span class="panel-note">橙色折线 / 低透明度面积</span>
          </div>

          <div v-if="latencyTrendPoints.length" class="chart-shell">
            <svg viewBox="0 0 360 180" class="trend-chart" role="img" aria-label="延迟趋势图">
              <g class="chart-grid">
                <line v-for="y in chartGuideLines" :key="`lat-${y}`" x1="18" :y1="y" x2="342" :y2="y"></line>
              </g>
              <path class="chart-area latency-area" :d="latencyTrendAreaPath"></path>
              <path class="chart-line latency-line" :d="latencyTrendLinePath"></path>
              <circle v-for="point in latencyTrendPoints" :key="point.key" class="chart-point latency-point" :cx="point.x" :cy="point.y" r="3.5"></circle>
            </svg>
          </div>
          <div v-else class="panel-empty">等待更多采样数据后显示趋势。</div>

          <div class="trend-summary">
            <span>当前平均 {{ formatDuration(overview.avgLatencyMs) }}</span>
            <span>峰值 {{ formatDuration(latencyPeak) }}</span>
            <span>样本 {{ historySamples.length }} 个</span>
          </div>
        </article>
      </section>

      <section class="watch-grid">
        <article class="panel-card secondary-card">
          <div class="panel-head">
            <div>
              <p class="section-kicker">ROUTES</p>
              <h3 class="panel-title">接口健康</h3>
            </div>
            <span class="panel-note">优先展示需要关注的接口</span>
          </div>

          <div v-if="loading && !routeRows.length" class="panel-empty">正在加载接口指标...</div>
          <div v-else-if="!routeRows.length" class="panel-empty">当前还没有接口请求数据。</div>
          <div v-else class="table-shell">
            <table class="metrics-table">
              <thead>
                <tr>
                  <th>接口</th>
                  <th>状态</th>
                  <th>请求量</th>
                  <th>错误率</th>
                  <th>平均延迟</th>
                </tr>
              </thead>
              <tbody>
                <tr v-for="route in routeRows" :key="`${route.method}-${route.path}`">
                  <td>
                    <div class="route-cell">
                      <span class="method-chip" :class="methodTone(route.method)">{{ route.method }}</span>
                      <span class="path-text">{{ route.path }}</span>
                    </div>
                  </td>
                  <td><span class="severity-badge" :class="route.severity.tone">{{ route.severity.label }}</span></td>
                  <td>{{ formatNumber(route.requestsTotal) }}</td>
                  <td>{{ formatPercent(route.errorRate) }}</td>
                  <td>{{ formatDuration(route.avgLatencyMs) }}</td>
                </tr>
              </tbody>
            </table>
          </div>
        </article>

        <article class="panel-card secondary-card">
          <div class="panel-head">
            <div>
              <p class="section-kicker">MODELS</p>
              <h3 class="panel-title">模型健康</h3>
            </div>
            <span class="panel-note">按模型聚合展示调用表现</span>
          </div>

          <div v-if="loading && !modelRows.length" class="panel-empty">正在加载模型指标...</div>
          <div v-else-if="!modelRows.length" class="panel-empty">当前还没有模型调用数据。</div>
          <div v-else class="table-shell">
            <table class="metrics-table">
              <thead>
                <tr>
                  <th>模型</th>
                  <th>状态</th>
                  <th>请求量</th>
                  <th>错误率</th>
                  <th>平均延迟</th>
                </tr>
              </thead>
              <tbody>
                <tr v-for="model in modelRows" :key="model.modelType">
                  <td>
                    <div class="model-cell">
                      <strong>{{ model.modelType }}</strong>
                      <span class="cell-subtext">{{ model.operationsLabel }}</span>
                    </div>
                  </td>
                  <td><span class="severity-badge" :class="model.severity.tone">{{ model.severity.label }}</span></td>
                  <td>{{ formatNumber(model.requestsTotal) }}</td>
                  <td>{{ formatPercent(model.errorRate) }}</td>
                  <td>{{ formatDuration(model.avgLatencyMs) }}</td>
                </tr>
              </tbody>
            </table>
          </div>
        </article>
      </section>
    </main>
  </div>
</template>

<script setup>
import { computed, onBeforeUnmount, onMounted, ref, watch } from 'vue'
import { useRouter } from 'vue-router'
import api from '../utils/api'
import { useUi } from '../composables/useUi'
import { isAdminToken } from '../utils/auth'

const DEFAULT_REFRESH_INTERVAL = 60000
const SUCCESS_CODE = 1000
const MAX_HISTORY_SAMPLES = 16
const CHART = { left: 18, right: 342, top: 22, bottom: 150 }
const refreshOptions = [
  { label: '60 秒刷新', value: 60000 },
  { label: '120 秒刷新', value: 120000 },
  { label: '300 秒刷新', value: 300000 }
]
const severityWeight = { neutral: 0, success: 1, warning: 2, danger: 3 }
const chartGuideLines = [22, 54, 86, 118, 150]

const router = useRouter()
const { showToast, confirmAction } = useUi()

const overview = ref({ requestsTotal: 0, errorsTotal: 0, errorRate: 0, avgLatencyMs: 0 })
const routes = ref([])
const models = ref([])
const historySamples = ref([])
const loading = ref(true)
const refreshing = ref(false)
const autoRefresh = ref(true)
const refreshIntervalMs = ref(DEFAULT_REFRESH_INTERVAL)
const lastUpdatedAt = ref(0)

let refreshTimer = null
let activeRequest = null

const numberFormatter = new Intl.NumberFormat('zh-CN')
const formatNumber = (value) => numberFormatter.format(Number(value) || 0)
const formatPercent = (value) => `${((Number(value) || 0) * 100).toFixed(1)}%`
const formatDuration = (value) => {
  const duration = Number(value) || 0
  return duration >= 1000 ? `${(duration / 1000).toFixed(duration >= 10000 ? 0 : 1)} s` : `${Math.round(duration)} ms`
}
const toTimestamp = (value) => {
  if (!value) return 0
  const timestamp = new Date(value).getTime()
  return Number.isFinite(timestamp) ? timestamp : 0
}
const formatTimestamp = (value) => {
  const timestamp = toTimestamp(value)
  if (!timestamp) return '-'
  return new Date(timestamp).toLocaleString('zh-CN', { month: '2-digit', day: '2-digit', hour: '2-digit', minute: '2-digit' })
}
const relativeTime = (timestamp) => {
  if (!timestamp) return '尚未刷新'
  const diff = Date.now() - timestamp
  if (diff < 5000) return '刚刚'
  if (diff < 60000) return `${Math.max(1, Math.floor(diff / 1000))} 秒前`
  if (diff < 3600000) return `${Math.max(1, Math.floor(diff / 60000))} 分钟前`
  return formatTimestamp(timestamp)
}
const formatRefreshInterval = (value) => `每 ${Math.floor((Number(value) || 0) / 1000)} 秒自动刷新`
const buildSeverity = (tone, label) => ({ tone, label, rank: severityWeight[tone] || 0 })

const routeSeverity = (route) => {
  const requestsTotal = Number(route?.requestsTotal) || 0
  const errorRate = Number(route?.errorRate) || 0
  const errorsTotal = Number(route?.errorsTotal) || 0
  const avgLatencyMs = Number(route?.avgLatencyMs) || 0
  const lastLatencyMs = Number(route?.lastLatencyMs) || 0
  const lastHttpStatus = Number(route?.lastHttpStatus) || 0
  if (!requestsTotal) return buildSeverity('neutral', '待采样')
  if (lastHttpStatus >= 500 || errorRate >= 0.1 || errorsTotal >= 5) return buildSeverity('danger', '异常')
  if (errorRate >= 0.03 || avgLatencyMs >= 1500 || lastLatencyMs >= 2500) return buildSeverity('warning', '警告')
  return buildSeverity('success', '正常')
}

const modelSeverity = (model) => {
  const requestsTotal = Number(model?.requestsTotal) || 0
  const errorRate = Number(model?.errorRate) || 0
  const errorsTotal = Number(model?.errorsTotal) || 0
  const avgLatencyMs = Number(model?.avgLatencyMs) || 0
  const lastLatencyMs = Number(model?.lastLatencyMs) || 0
  const lastFailureAt = toTimestamp(model?.lastFailureAt)
  const lastSuccessAt = toTimestamp(model?.lastSuccessAt)
  if (!requestsTotal) return buildSeverity('neutral', '待采样')
  if (errorRate >= 0.1 || errorsTotal >= 3) return buildSeverity('danger', '异常')
  if (avgLatencyMs >= 3000 || lastLatencyMs >= 4000 || errorRate >= 0.03 || lastFailureAt > lastSuccessAt) return buildSeverity('warning', '警告')
  return buildSeverity('success', '正常')
}

const routeStatusList = computed(() => routes.value.map((route) => ({ ...route, severity: routeSeverity(route) })))
const rawModelFamilies = computed(() => {
  const familyMap = new Map()
  for (const item of models.value) {
    const modelType = item.modelType || 'unknown'
    const current = familyMap.get(modelType) || {
      modelType, requestsTotal: 0, errorsTotal: 0, weightedLatency: 0, lastLatencyMs: 0, lastSuccessAt: '', lastFailureAt: '', operations: []
    }
    const requestsTotal = Number(item.requestsTotal) || 0
    current.requestsTotal += requestsTotal
    current.errorsTotal += Number(item.errorsTotal) || 0
    current.weightedLatency += (Number(item.avgLatencyMs) || 0) * requestsTotal
    current.lastLatencyMs = Math.max(current.lastLatencyMs, Number(item.lastLatencyMs) || 0)
    current.operations.push(item.operation || 'unknown')
    if (toTimestamp(item.lastSuccessAt) > toTimestamp(current.lastSuccessAt)) current.lastSuccessAt = item.lastSuccessAt
    if (toTimestamp(item.lastFailureAt) > toTimestamp(current.lastFailureAt)) current.lastFailureAt = item.lastFailureAt
    familyMap.set(modelType, current)
  }
  return Array.from(familyMap.values()).map((item) => {
    const requestsTotal = Number(item.requestsTotal) || 0
    return {
      modelType: item.modelType,
      requestsTotal,
      errorsTotal: item.errorsTotal,
      errorRate: requestsTotal ? item.errorsTotal / requestsTotal : 0,
      avgLatencyMs: requestsTotal ? item.weightedLatency / requestsTotal : 0,
      lastLatencyMs: item.lastLatencyMs,
      lastSuccessAt: item.lastSuccessAt,
      lastFailureAt: item.lastFailureAt,
      operationsLabel: [...new Set(item.operations)].join(' / ')
    }
  })
})

const modelFamilies = computed(() => rawModelFamilies.value.map((item) => ({ ...item, severity: modelSeverity(item) })))
const abnormalRouteCount = computed(() => routeStatusList.value.filter((item) => item.severity.rank >= severityWeight.warning).length)
const abnormalModelCount = computed(() => modelFamilies.value.filter((item) => item.severity.rank >= severityWeight.warning).length)

const healthStatus = computed(() => {
  const requestsTotal = Number(overview.value.requestsTotal) || 0
  const errorRate = Number(overview.value.errorRate) || 0
  const avgLatencyMs = Number(overview.value.avgLatencyMs) || 0
  if (!requestsTotal) return buildSeverity('neutral', '待观察')
  if (errorRate >= 0.1 || avgLatencyMs >= 4000 || abnormalRouteCount.value >= 3) return buildSeverity('danger', '存在异常')
  if (errorRate >= 0.03 || avgLatencyMs >= 1500 || abnormalRouteCount.value > 0 || abnormalModelCount.value > 0) return buildSeverity('warning', '轻微波动')
  return buildSeverity('success', '运行稳定')
})

const snapshotTitle = computed(() => {
  const requestsTotal = Number(overview.value.requestsTotal) || 0
  if (!requestsTotal) return '系统刚启动，正在等待采样'
  if (healthStatus.value.tone === 'danger') return '当前存在明显异常，建议优先检查异常接口和模型'
  if (healthStatus.value.tone === 'warning') return '当前有轻微波动，适合继续观察趋势变化'
  return '当前整体运行稳定，可以继续观察后续变化'
})

const snapshotSummary = computed(() => {
  const errorRate = Number(overview.value.errorRate) || 0
  const avgLatencyMs = Number(overview.value.avgLatencyMs) || 0
  if (healthStatus.value.tone === 'danger') {
    return `当前错误率为 ${formatPercent(errorRate)}，平均延迟为 ${formatDuration(avgLatencyMs)}。如果异常项持续增加，建议立即排查上游模型或高风险接口。`
  }
  if (healthStatus.value.tone === 'warning') {
    return `当前错误率和延迟出现一定波动，异常接口 ${abnormalRouteCount.value} 个、异常模型 ${abnormalModelCount.value} 个，建议继续观察最近几个采样周期。`
  }
  if (healthStatus.value.tone === 'neutral') {
    return '监控刚开始采样，当前还不足以给出稳定结论，等待更多请求进入后再判断。'
  }
  return `当前错误率 ${formatPercent(errorRate)}，平均延迟 ${formatDuration(avgLatencyMs)}，系统整体处于稳定状态。`
})

const statusStripText = computed(() => {
  if (healthStatus.value.tone === 'danger') return `当前存在异常，异常接口 ${abnormalRouteCount.value} 个，异常模型 ${abnormalModelCount.value} 个，建议立即排查。`
  if (healthStatus.value.tone === 'warning') return `系统出现轻微波动，异常接口 ${abnormalRouteCount.value} 个，异常模型 ${abnormalModelCount.value} 个，建议继续观察。`
  if (healthStatus.value.tone === 'neutral') return '监控刚开始采样，当前没有足够数据形成结论。'
  return '系统运行稳定，当前无明显异常，建议继续保持自动刷新观察。'
})

const overviewCards = computed(() => [
  { label: '累计请求', value: formatNumber(overview.value.requestsTotal), desc: '当前监控窗口内累计请求量', tone: 'info' },
  { label: '累计错误', value: formatNumber(overview.value.errorsTotal), desc: 'HTTP 或业务异常累计次数', tone: 'danger' },
  { label: '错误率', value: formatPercent(overview.value.errorRate), desc: '快速判断系统是否稳定', tone: Number(overview.value.errorRate) >= 0.03 ? 'warning' : 'success' },
  { label: '平均延迟', value: formatDuration(overview.value.avgLatencyMs), desc: '请求的平均耗时', tone: Number(overview.value.avgLatencyMs) >= 1500 ? 'warning' : 'info' },
  { label: '异常接口数', value: formatNumber(abnormalRouteCount.value), desc: '需要优先关注的接口数量', tone: abnormalRouteCount.value > 0 ? 'warning' : 'success' },
  { label: '异常模型数', value: formatNumber(abnormalModelCount.value), desc: '需要关注的模型数量', tone: abnormalModelCount.value > 0 ? 'warning' : 'success' }
])

const routeRows = computed(() => [...routeStatusList.value].sort((left, right) => {
  if (right.severity.rank !== left.severity.rank) return right.severity.rank - left.severity.rank
  return (Number(right.requestsTotal) || 0) - (Number(left.requestsTotal) || 0)
}).slice(0, 6))

const modelRows = computed(() => [...modelFamilies.value].sort((left, right) => {
  if (right.severity.rank !== left.severity.rank) return right.severity.rank - left.severity.rank
  return (Number(right.requestsTotal) || 0) - (Number(left.requestsTotal) || 0)
}).slice(0, 6))

const normalizeHistorySamples = (archives = []) => {
  const normalized = [...archives]
    .map((item) => ({ timestamp: toTimestamp(item.timestamp), requestsTotal: Number(item.requestsTotal) || 0, avgLatencyMs: Number(item.avgLatencyMs) || 0 }))
    .filter((item) => item.timestamp > 0)
    .sort((left, right) => left.timestamp - right.timestamp)
    .slice(-MAX_HISTORY_SAMPLES)

  return normalized.map((item, index) => ({
    ...item,
    requestDelta: index === 0 ? item.requestsTotal : Math.max(0, item.requestsTotal - normalized[index - 1].requestsTotal)
  }))
}

const requestPeak = computed(() => Math.max(0, ...historySamples.value.map((item) => item.requestDelta)))
const latencyPeak = computed(() => Math.max(0, ...historySamples.value.map((item) => item.avgLatencyMs)))
const latestRequestDelta = computed(() => historySamples.value[historySamples.value.length - 1]?.requestDelta || 0)
const lastUpdatedLabel = computed(() => relativeTime(lastUpdatedAt.value))
const refreshIntervalLabel = computed(() => formatRefreshInterval(refreshIntervalMs.value))

const buildChartPoints = (samples, valueKey, peakValue) => {
  if (!samples.length) return []
  const peak = peakValue || 1
  const width = CHART.right - CHART.left
  const height = CHART.bottom - CHART.top
  return samples.map((item, index) => {
    const x = samples.length === 1 ? CHART.left + width / 2 : CHART.left + (width * index) / (samples.length - 1)
    const ratio = Math.max(0, (Number(item[valueKey]) || 0) / peak)
    const y = CHART.bottom - ratio * height
    return { key: `${valueKey}-${item.timestamp}-${index}`, x: Number(x.toFixed(2)), y: Number(y.toFixed(2)) }
  })
}

const buildLinePath = (points) => points.map((point, index) => `${index === 0 ? 'M' : 'L'} ${point.x} ${point.y}`).join(' ')
const buildAreaPath = (points) => {
  if (!points.length) return ''
  const first = points[0]
  const last = points[points.length - 1]
  return `${buildLinePath(points)} L ${last.x} ${CHART.bottom} L ${first.x} ${CHART.bottom} Z`
}

const requestTrendPoints = computed(() => buildChartPoints(historySamples.value, 'requestDelta', requestPeak.value))
const latencyTrendPoints = computed(() => buildChartPoints(historySamples.value, 'avgLatencyMs', latencyPeak.value))
const requestTrendLinePath = computed(() => buildLinePath(requestTrendPoints.value))
const requestTrendAreaPath = computed(() => buildAreaPath(requestTrendPoints.value))
const latencyTrendLinePath = computed(() => buildLinePath(latencyTrendPoints.value))
const latencyTrendAreaPath = computed(() => buildAreaPath(latencyTrendPoints.value))

const methodTone = (method) => {
  const value = (method || '').toUpperCase()
  if (value === 'GET') return 'get'
  if (value === 'POST') return 'post'
  return 'default'
}

const isSuccessResponse = (response) => response?.data?.status_code === SUCCESS_CODE
const isForbiddenResponse = (response) => response?.data?.status_code === 3001

const applyMetricsSnapshot = (snapshot) => {
  overview.value = { ...overview.value, ...(snapshot?.overview || {}) }
  routes.value = Array.isArray(snapshot?.routes) ? snapshot.routes : []
  models.value = Array.isArray(snapshot?.models) ? snapshot.models : []
  const nextHistory = normalizeHistorySamples(snapshot?.archives || [])
  historySamples.value = nextHistory.length ? nextHistory : [{
    timestamp: Date.now(),
    requestsTotal: Number(overview.value.requestsTotal) || 0,
    requestDelta: Number(overview.value.requestsTotal) || 0,
    avgLatencyMs: Number(overview.value.avgLatencyMs) || 0
  }]
  lastUpdatedAt.value = Date.now()
}

const fetchSnapshotMetrics = async () => {
  const response = await api.get('/admin/metrics/all')
  if (isForbiddenResponse(response)) {
    const error = new Error('FORBIDDEN')
    error.code = 'FORBIDDEN'
    throw error
  }
  if (!isSuccessResponse(response)) throw new Error('SNAPSHOT_FAILED')
  return response.data?.snapshot || {}
}

const fetchMetrics = async ({ silent = false } = {}) => {
  if (activeRequest) return activeRequest
  if (!lastUpdatedAt.value) loading.value = true
  refreshing.value = true

  activeRequest = fetchSnapshotMetrics()
    .then((snapshot) => applyMetricsSnapshot(snapshot))
    .catch((error) => {
      console.error('Fetch metrics error:', error)
      if (error?.code === 'FORBIDDEN' || !isAdminToken(localStorage.getItem('token'))) {
        stopAutoRefresh()
        showToast('当前账号没有管理员权限，已返回控制台。', 'error')
        router.replace('/menu')
        return
      }
      if (!silent || !lastUpdatedAt.value) {
        showToast('监控数据刷新失败，请检查后端服务状态。', 'error')
      }
    })
    .finally(() => {
      loading.value = false
      refreshing.value = false
      activeRequest = null
    })

  return activeRequest
}

const stopAutoRefresh = () => {
  if (refreshTimer) {
    window.clearInterval(refreshTimer)
    refreshTimer = null
  }
}

const startAutoRefresh = () => {
  stopAutoRefresh()
  if (!autoRefresh.value) return
  refreshTimer = window.setInterval(() => {
    fetchMetrics({ silent: true })
  }, refreshIntervalMs.value)
}

const goMenu = () => router.push('/menu')
const goChat = () => router.push('/ai-chat')

const logout = async () => {
  const confirmed = await confirmAction({
    title: '退出当前账号？',
    message: '退出后将返回登录页，但不会影响后端已采集的监控数据。',
    confirmText: '退出登录',
    cancelText: '继续查看',
    intent: 'danger'
  })
  if (!confirmed) return
  localStorage.removeItem('token')
  localStorage.removeItem('isAdmin')
  router.push('/login')
}

watch([autoRefresh, refreshIntervalMs], () => {
  startAutoRefresh()
})

onMounted(async () => {
  if (!isAdminToken(localStorage.getItem('token'))) {
    showToast('当前账号没有管理员权限，无法进入管理页。', 'error')
    router.replace('/menu')
    return
  }
  await fetchMetrics()
  startAutoRefresh()
})

onBeforeUnmount(() => {
  stopAutoRefresh()
})
</script>

<style scoped>
.admin-shell {
  --admin-bg: #f4f8f6;
  --admin-grid: rgba(80, 120, 100, 0.06);
  --admin-card-primary: #ffffff;
  --admin-card-secondary: #fcfdfc;
  --admin-border: #e5efea;
  --admin-shadow: rgba(15, 23, 42, 0.08);
  --admin-primary: #2ec5a7;
  --admin-primary-dark: #20a88d;
  --admin-primary-soft: #eafbf4;
  --admin-success: #22c55e;
  --admin-warning: #f59e0b;
  --admin-danger: #ef4444;
  --admin-info: #3b82f6;
  --admin-text: #18312b;
  --admin-text-soft: #5d746d;
  --admin-text-muted: #7d948e;
  min-height: 100vh;
  position: relative;
  overflow: hidden;
  padding: 28px;
  background: var(--admin-bg);
}

.admin-bg {
  position: absolute;
  inset: 0;
  background:
    radial-gradient(circle at top right, rgba(46, 197, 167, 0.12), transparent 24%),
    radial-gradient(circle at bottom left, rgba(59, 130, 246, 0.06), transparent 20%),
    linear-gradient(180deg, #f4f8f6 0%, #f1f7f3 100%);
}

.glow {
  position: absolute;
  border-radius: 50%;
  filter: blur(100px);
  opacity: 0.42;
}

.glow-a {
  width: 380px;
  height: 380px;
  top: -120px;
  right: -90px;
  background: rgba(46, 197, 167, 0.18);
}

.glow-b {
  width: 320px;
  height: 320px;
  left: -100px;
  bottom: -120px;
  background: rgba(59, 130, 246, 0.08);
}

.grid {
  position: absolute;
  inset: 0;
  background-image:
    linear-gradient(var(--admin-grid) 1px, transparent 1px),
    linear-gradient(90deg, var(--admin-grid) 1px, transparent 1px);
  background-size: 48px 48px;
}

.admin-header,
.admin-main {
  position: relative;
  z-index: 1;
}

.admin-header,
.panel-card {
  border: 1px solid var(--admin-border);
  box-shadow: 0 20px 40px var(--admin-shadow);
}

.admin-header {
  max-width: 1320px;
  margin: 0 auto 24px;
  padding: 22px 24px;
  border-radius: 24px;
  background: rgba(255, 255, 255, 0.92);
  display: grid;
  grid-template-columns: minmax(280px, 1.08fr) minmax(300px, 0.92fr) minmax(320px, 1fr);
  gap: 18px;
  align-items: stretch;
}

.header-brand {
  display: flex;
  gap: 16px;
}

.brand-badge {
  width: 54px;
  height: 54px;
  display: grid;
  place-items: center;
  border-radius: 16px;
  background: linear-gradient(135deg, var(--admin-primary), #8ee3d3);
  color: #fff;
  font-family: var(--wa-font-label);
  font-size: 14px;
  letter-spacing: 2px;
  box-shadow: 0 14px 26px rgba(46, 197, 167, 0.24);
}

.brand-kicker,
.section-kicker,
.metric-label,
.panel-note,
.header-btn,
.refresh-toggle,
.refresh-select,
.method-chip,
.severity-badge,
.status-kicker,
.status-label {
  font-family: var(--wa-font-label);
}

.brand-kicker,
.section-kicker,
.status-kicker {
  margin: 0 0 8px;
  font-size: 11px;
  letter-spacing: 0.2em;
  color: var(--admin-primary-dark);
}

.brand-title,
.hero-title,
.panel-title {
  margin: 0;
  color: var(--admin-text);
}

.brand-title {
  font-size: 28px;
}

.brand-subtitle,
.hero-desc,
.metric-desc,
.panel-note,
.panel-empty,
.path-text,
.cell-subtext,
.diagnosis-text {
  color: var(--admin-text-soft);
}

.brand-subtitle {
  margin: 8px 0 0;
  font-size: 14px;
  line-height: 1.7;
}

.header-status {
  padding: 16px 18px;
  border-radius: 18px;
  background: var(--admin-card-secondary);
  border: 1px solid var(--admin-border);
}

.status-summary {
  display: grid;
  gap: 12px;
}

.status-item {
  display: flex;
  flex-direction: column;
  gap: 6px;
  padding-bottom: 12px;
  border-bottom: 1px solid rgba(229, 239, 234, 0.9);
}

.status-item:last-child {
  padding-bottom: 0;
  border-bottom: none;
}

.status-label {
  font-size: 10px;
  letter-spacing: 0.14em;
  color: var(--admin-text-muted);
}

.status-value {
  font-size: 14px;
  line-height: 1.5;
  color: var(--admin-text);
}

.status-value.success {
  color: var(--admin-success);
}

.status-value.warning {
  color: var(--admin-warning);
}

.status-value.danger {
  color: var(--admin-danger);
}

.status-value.neutral {
  color: var(--admin-info);
}

.header-actions {
  display: flex;
  flex-direction: column;
  justify-content: space-between;
  gap: 14px;
}

.refresh-controls,
.action-row {
  display: flex;
  flex-wrap: wrap;
  gap: 10px;
}

.header-btn,
.refresh-select {
  border: 1px solid transparent;
  border-radius: 14px;
  padding: 12px 16px;
  font-size: 11px;
  letter-spacing: 0.12em;
  transition: transform 0.3s ease, background 0.3s ease, border-color 0.3s ease, box-shadow 0.3s ease, opacity 0.3s ease;
}

.header-btn {
  cursor: pointer;
}

.header-btn.weak {
  color: var(--admin-text-soft);
  background: rgba(255, 255, 255, 0.96);
  border-color: var(--admin-border);
}

.header-btn.secondary {
  color: var(--admin-primary-dark);
  background: var(--admin-primary-soft);
  border-color: rgba(46, 197, 167, 0.16);
}

.header-btn.primary {
  color: #fff;
  background: linear-gradient(135deg, var(--admin-primary), #69d9c2);
  box-shadow: 0 12px 22px rgba(46, 197, 167, 0.22);
}

.header-btn.danger {
  color: var(--admin-danger);
  background: rgba(239, 68, 68, 0.08);
  border-color: rgba(239, 68, 68, 0.14);
}

.header-btn:hover:not(:disabled),
.refresh-select:hover:not(:disabled) {
  transform: translateY(-1px);
}

.header-btn:disabled,
.refresh-select:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

.refresh-toggle {
  display: inline-flex;
  align-items: center;
  gap: 8px;
  padding: 12px 14px;
  border-radius: 14px;
  background: #fff;
  border: 1px solid var(--admin-border);
  color: var(--admin-text-soft);
  font-size: 11px;
  letter-spacing: 0.12em;
}

.refresh-select {
  min-width: 148px;
  cursor: pointer;
  color: var(--admin-text);
  background: #fff;
  border-color: var(--admin-border);
}

.admin-main {
  max-width: 1320px;
  margin: 0 auto;
  display: grid;
  gap: 20px;
}

.panel-card {
  border-radius: 24px;
  padding: 24px;
}

.primary-card {
  background: var(--admin-card-primary);
}

.secondary-card,
.metric-card {
  background: var(--admin-card-secondary);
}

.hero-panel {
  display: grid;
  gap: 18px;
}

.hero-main {
  display: grid;
  grid-template-columns: minmax(0, 1.25fr) minmax(290px, 0.75fr);
  gap: 20px;
}

.hero-title {
  font-size: 34px;
  line-height: 1.24;
}

.hero-desc {
  max-width: 760px;
  margin: 16px 0 0;
  font-size: 15px;
  line-height: 1.82;
}

.diagnosis-box {
  margin-top: 22px;
  padding: 18px 20px;
  border-radius: 18px;
  background: rgba(244, 248, 246, 0.92);
  border: 1px solid var(--admin-border);
}

.diagnosis-label {
  display: inline-block;
  font-family: var(--wa-font-label);
  font-size: 10px;
  letter-spacing: 0.16em;
  color: var(--admin-text-muted);
}

.diagnosis-title {
  display: block;
  margin-top: 10px;
  font-size: 22px;
  line-height: 1.45;
  color: var(--admin-text);
}

.diagnosis-text {
  margin: 10px 0 0;
  font-size: 14px;
  line-height: 1.82;
}

.hero-side {
  display: flex;
  flex-direction: column;
  gap: 14px;
}

.hero-stats {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 12px;
}

.hero-stat {
  padding: 14px 14px 12px;
  border-radius: 16px;
  background: rgba(252, 253, 252, 0.98);
  border: 1px solid var(--admin-border);
}

.hero-stat-label {
  display: block;
  font-family: var(--wa-font-label);
  font-size: 10px;
  letter-spacing: 0.12em;
  color: var(--admin-text-muted);
}

.hero-stat-value {
  display: block;
  margin-top: 8px;
  font-size: 15px;
  line-height: 1.55;
  color: var(--admin-text);
}

.health-pill,
.severity-badge {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  padding: 8px 14px;
  border-radius: 999px;
  font-size: 11px;
  letter-spacing: 0.12em;
}

.health-pill.success,
.severity-badge.success {
  background: rgba(34, 197, 94, 0.12);
  color: var(--admin-success);
}

.health-pill.warning,
.severity-badge.warning {
  background: rgba(245, 158, 11, 0.12);
  color: var(--admin-warning);
}

.health-pill.danger,
.severity-badge.danger {
  background: rgba(239, 68, 68, 0.12);
  color: var(--admin-danger);
}

.health-pill.neutral,
.severity-badge.neutral {
  background: rgba(59, 130, 246, 0.12);
  color: var(--admin-info);
}

.hero-status-strip {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 14px 16px;
  border-radius: 16px;
  font-size: 14px;
  line-height: 1.7;
}

.hero-status-strip.success {
  background: rgba(34, 197, 94, 0.1);
  color: #17803f;
}

.hero-status-strip.warning {
  background: rgba(245, 158, 11, 0.1);
  color: #b45309;
}

.hero-status-strip.danger {
  background: rgba(239, 68, 68, 0.1);
  color: #b91c1c;
}

.hero-status-strip.neutral {
  background: rgba(59, 130, 246, 0.08);
  color: #2563eb;
}

.strip-dot {
  width: 10px;
  height: 10px;
  border-radius: 50%;
  background: currentColor;
  flex-shrink: 0;
}

.overview-grid,
.trend-grid,
.watch-grid {
  display: grid;
  gap: 18px;
}

.overview-grid {
  grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
}

.trend-grid,
.watch-grid {
  grid-template-columns: repeat(2, minmax(0, 1fr));
}

.metric-card {
  min-height: 150px;
}

.metric-card.info {
  background: linear-gradient(180deg, #fcfdfc 0%, #f5fbf8 100%);
}

.metric-card.success {
  background: linear-gradient(180deg, #fcfdfc 0%, #f3fcf5 100%);
}

.metric-card.warning {
  background: linear-gradient(180deg, #fcfdfc 0%, #fff9f0 100%);
}

.metric-card.danger {
  background: linear-gradient(180deg, #fcfdfc 0%, #fff5f5 100%);
}

.metric-label {
  display: block;
  margin-bottom: 12px;
  font-size: 11px;
  letter-spacing: 0.18em;
  color: var(--admin-text-muted);
}

.metric-value {
  display: block;
  font-size: 30px;
  color: var(--admin-text);
}

.metric-desc {
  margin: 12px 0 0;
  font-size: 14px;
  line-height: 1.72;
}

.panel-head {
  display: flex;
  justify-content: space-between;
  gap: 12px;
  align-items: flex-start;
}

.panel-title {
  font-size: 24px;
}

.panel-note {
  margin: 0;
  font-size: 13px;
  line-height: 1.6;
}

.chart-shell {
  margin-top: 22px;
  padding: 14px 10px 8px;
  border-radius: 18px;
  background: rgba(255, 255, 255, 0.96);
  border: 1px solid var(--admin-border);
}

.trend-chart {
  width: 100%;
  height: 180px;
  display: block;
}

.chart-grid line {
  stroke: rgba(93, 116, 109, 0.12);
  stroke-width: 1;
}

.chart-area {
  stroke: none;
}

.request-area {
  fill: rgba(34, 197, 94, 0.14);
}

.latency-area {
  fill: rgba(245, 158, 11, 0.14);
}

.chart-line {
  fill: none;
  stroke-width: 3;
  stroke-linecap: round;
  stroke-linejoin: round;
}

.request-line {
  stroke: var(--admin-success);
}

.latency-line {
  stroke: var(--admin-warning);
}

.request-point {
  fill: var(--admin-success);
}

.latency-point {
  fill: var(--admin-warning);
}

.trend-summary {
  margin-top: 18px;
  display: flex;
  justify-content: space-between;
  gap: 12px;
  color: var(--admin-text-soft);
  font-size: 13px;
}

.panel-empty {
  margin-top: 18px;
  padding: 18px;
  border-radius: 16px;
  background: rgba(255, 255, 255, 0.9);
  font-size: 14px;
  line-height: 1.7;
  border: 1px dashed var(--admin-border);
}

.table-shell {
  margin-top: 18px;
  overflow-x: auto;
}

.metrics-table {
  width: 100%;
  min-width: 680px;
  border-collapse: collapse;
}

.metrics-table th,
.metrics-table td {
  padding: 14px 12px;
  border-bottom: 1px solid rgba(229, 239, 234, 0.9);
  text-align: left;
  font-size: 14px;
  color: var(--admin-text);
  vertical-align: middle;
}

.metrics-table th {
  font-family: var(--wa-font-label);
  font-size: 11px;
  letter-spacing: 0.14em;
  color: var(--admin-text-muted);
}

.route-cell,
.model-cell {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.route-cell {
  flex-direction: row;
  align-items: center;
  gap: 10px;
}

.method-chip {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  min-width: 54px;
  padding: 6px 8px;
  border-radius: 999px;
  font-size: 10px;
  letter-spacing: 0.12em;
}

.method-chip.get {
  background: rgba(34, 197, 94, 0.12);
  color: var(--admin-success);
}

.method-chip.post {
  background: rgba(59, 130, 246, 0.12);
  color: var(--admin-info);
}

.method-chip.default {
  background: rgba(148, 163, 184, 0.16);
  color: var(--admin-text-soft);
}

.path-text {
  word-break: break-all;
}

.cell-subtext {
  font-size: 12px;
}

@media (max-width: 1180px) {
  .admin-header {
    grid-template-columns: 1fr;
  }

  .header-actions {
    align-items: flex-start;
  }
}

@media (max-width: 1100px) {
  .hero-main,
  .trend-grid,
  .watch-grid {
    grid-template-columns: 1fr;
  }
}

@media (max-width: 768px) {
  .admin-shell {
    padding: 16px;
  }

  .admin-header,
  .panel-card {
    padding: 20px 18px;
  }

  .brand-title {
    font-size: 24px;
  }

  .hero-title {
    font-size: 28px;
  }

  .hero-stats {
    grid-template-columns: 1fr;
  }

  .refresh-controls,
  .action-row,
  .trend-summary {
    flex-direction: column;
  }

  .header-btn,
  .refresh-toggle,
  .refresh-select {
    width: 100%;
  }
}
</style>
