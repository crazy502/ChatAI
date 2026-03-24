<template>
  <div class="admin-shell">
    <div class="admin-bg">
      <div class="glow glow-a"></div>
      <div class="glow glow-b"></div>
      <div class="grid"></div>
    </div>

    <header class="admin-header">
      <div class="brand-block">
        <div class="brand-badge">OPS</div>
        <div>
          <p class="brand-kicker">LIGHTWEIGHT MONITORING</p>
          <h1 class="brand-title">GopherAI 管理概览</h1>
          <p class="brand-subtitle">只保留请求、错误、延迟和模型运行状态这几项核心信息。</p>
        </div>
      </div>

      <div class="header-actions">
        <button class="header-btn" type="button" @click="goMenu">返回控制台</button>
        <button class="header-btn" type="button" @click="goChat">进入对话</button>
        <button class="header-btn primary" type="button" :disabled="refreshing" @click="fetchMetrics()">
          {{ refreshing ? '刷新中...' : '立即刷新' }}
        </button>

        <label class="refresh-toggle">
          <input v-model="autoRefresh" type="checkbox" />
          <span>自动刷新</span>
        </label>

        <select v-model.number="refreshIntervalMs" class="refresh-select" :disabled="!autoRefresh">
          <option v-for="option in refreshOptions" :key="option.value" :value="option.value">
            {{ option.label }}
          </option>
        </select>

        <button class="header-btn danger" type="button" @click="logout">退出登录</button>
      </div>
    </header>

    <main class="admin-main">
      <section class="hero-panel panel-card">
        <div>
          <p class="section-kicker">SYSTEM SNAPSHOT</p>
          <h2 class="hero-title">先看是否稳定，再决定要不要继续排查</h2>
          <p class="hero-desc">
            当前只展示最有价值的管理员信息：请求量、错误率、平均延迟、接口健康和模型健康。
            不再展开用户维度、业务码维度和失败明细这类次要监控。
          </p>
        </div>

        <div class="hero-side">
          <div class="health-pill" :class="healthStatus.tone">{{ healthStatus.label }}</div>
          <p class="hero-meta">最近刷新：{{ lastUpdatedLabel }}</p>
          <p class="hero-meta">刷新策略：{{ autoRefresh ? refreshIntervalLabel : '手动刷新' }}</p>
          <p class="hero-meta">异常接口：{{ abnormalRouteCount }} 个</p>
          <p class="hero-meta">异常模型：{{ abnormalModelCount }} 个</p>
        </div>
      </section>

      <section class="overview-grid">
        <article v-for="card in overviewCards" :key="card.label" class="metric-card panel-card">
          <span class="metric-label">{{ card.label }}</span>
          <strong class="metric-value">{{ card.value }}</strong>
          <p class="metric-desc">{{ card.desc }}</p>
        </article>
      </section>

      <section class="trend-grid">
        <article class="panel-card">
          <div class="panel-head">
            <div>
              <p class="section-kicker">REQUEST TREND</p>
              <h3 class="panel-title">请求增量趋势</h3>
            </div>
            <span class="panel-note">来自后端归档快照</span>
          </div>

          <div v-if="requestTrendBars.length" class="trend-bars">
            <div v-for="bar in requestTrendBars" :key="bar.key" class="trend-bar" :title="bar.label">
              <span class="trend-bar-fill success" :style="{ height: bar.height + '%' }"></span>
            </div>
          </div>
          <div v-else class="panel-empty">等待更多采样数据后显示趋势。</div>

          <div class="trend-summary">
            <span>当前增量 {{ formatNumber(latestRequestDelta) }}</span>
            <span>峰值 {{ formatNumber(requestPeak) }}</span>
          </div>
        </article>

        <article class="panel-card">
          <div class="panel-head">
            <div>
              <p class="section-kicker">LATENCY TREND</p>
              <h3 class="panel-title">平均延迟趋势</h3>
            </div>
            <span class="panel-note">帮助判断系统是否变慢</span>
          </div>

          <div v-if="latencyTrendBars.length" class="trend-bars">
            <div v-for="bar in latencyTrendBars" :key="bar.key" class="trend-bar" :title="bar.label">
              <span class="trend-bar-fill warning" :style="{ height: bar.height + '%' }"></span>
            </div>
          </div>
          <div v-else class="panel-empty">等待更多采样数据后显示趋势。</div>

          <div class="trend-summary">
            <span>当前平均 {{ formatDuration(overview.avgLatencyMs) }}</span>
            <span>峰值 {{ formatDuration(latencyPeak) }}</span>
          </div>
        </article>
      </section>

      <section class="watch-grid">
        <article class="panel-card">
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
            <table class="metrics-table compact">
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
                  <td>
                    <span class="severity-badge" :class="route.severity.tone">{{ route.severity.label }}</span>
                  </td>
                  <td>{{ formatNumber(route.requestsTotal) }}</td>
                  <td>{{ formatPercent(route.errorRate) }}</td>
                  <td>{{ formatDuration(route.avgLatencyMs) }}</td>
                </tr>
              </tbody>
            </table>
          </div>
        </article>

        <article class="panel-card">
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
            <table class="metrics-table compact">
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
                  <td>
                    <span class="severity-badge" :class="model.severity.tone">{{ model.severity.label }}</span>
                  </td>
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

const refreshOptions = [
  { label: '60 秒刷新', value: 60000 },
  { label: '120 秒刷新', value: 120000 },
  { label: '300 秒刷新', value: 300000 }
]

const severityWeight = {
  neutral: 0,
  success: 1,
  warning: 2,
  danger: 3
}

const router = useRouter()
const { showToast, confirmAction } = useUi()

const overview = ref({
  requestsTotal: 0,
  errorsTotal: 0,
  errorRate: 0,
  avgLatencyMs: 0
})
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
  if (duration >= 1000) {
    return `${(duration / 1000).toFixed(duration >= 10000 ? 0 : 1)} s`
  }
  return `${Math.round(duration)} ms`
}

const toTimestamp = (value) => {
  if (!value) {
    return 0
  }
  const timestamp = new Date(value).getTime()
  return Number.isFinite(timestamp) ? timestamp : 0
}

const formatTimestamp = (value) => {
  const timestamp = toTimestamp(value)
  if (!timestamp) {
    return '-'
  }
  return new Date(timestamp).toLocaleString('zh-CN', {
    month: '2-digit',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit'
  })
}

const relativeTime = (timestamp) => {
  if (!timestamp) {
    return '尚未刷新'
  }
  const diff = Date.now() - timestamp
  if (diff < 5000) {
    return '刚刚'
  }
  if (diff < 60000) {
    return `${Math.max(1, Math.floor(diff / 1000))} 秒前`
  }
  if (diff < 3600000) {
    return `${Math.max(1, Math.floor(diff / 60000))} 分钟前`
  }
  return formatTimestamp(timestamp)
}

const formatRefreshInterval = (value) => `每 ${Math.floor((Number(value) || 0) / 1000)} 秒自动刷新`

const buildSeverity = (tone, label) => ({
  tone,
  label,
  rank: severityWeight[tone] || 0
})

const routeSeverity = (route) => {
  const requestsTotal = Number(route?.requestsTotal) || 0
  const errorRate = Number(route?.errorRate) || 0
  const errorsTotal = Number(route?.errorsTotal) || 0
  const avgLatencyMs = Number(route?.avgLatencyMs) || 0
  const lastLatencyMs = Number(route?.lastLatencyMs) || 0
  const lastHttpStatus = Number(route?.lastHttpStatus) || 0

  if (!requestsTotal) {
    return buildSeverity('neutral', '待采样')
  }
  if (lastHttpStatus >= 500 || errorRate >= 0.1 || errorsTotal >= 5) {
    return buildSeverity('danger', '高风险')
  }
  if (errorRate >= 0.03 || avgLatencyMs >= 1500 || lastLatencyMs >= 2500) {
    return buildSeverity('warning', '需关注')
  }
  return buildSeverity('success', '稳定')
}

const modelSeverity = (model) => {
  const requestsTotal = Number(model?.requestsTotal) || 0
  const errorRate = Number(model?.errorRate) || 0
  const errorsTotal = Number(model?.errorsTotal) || 0
  const avgLatencyMs = Number(model?.avgLatencyMs) || 0
  const lastLatencyMs = Number(model?.lastLatencyMs) || 0
  const lastFailureAt = toTimestamp(model?.lastFailureAt)
  const lastSuccessAt = toTimestamp(model?.lastSuccessAt)

  if (!requestsTotal) {
    return buildSeverity('neutral', '待采样')
  }
  if (errorRate >= 0.1 || errorsTotal >= 3) {
    return buildSeverity('danger', '高风险')
  }
  if (avgLatencyMs >= 3000 || lastLatencyMs >= 4000 || errorRate >= 0.03 || lastFailureAt > lastSuccessAt) {
    return buildSeverity('warning', '需关注')
  }
  return buildSeverity('success', '稳定')
}

const routeStatusList = computed(() => routes.value.map((route) => ({
  ...route,
  severity: routeSeverity(route)
})))

const rawModelFamilies = computed(() => {
  const familyMap = new Map()

  for (const item of models.value) {
    const modelType = item.modelType || 'unknown'
    const current = familyMap.get(modelType) || {
      modelType,
      requestsTotal: 0,
      errorsTotal: 0,
      weightedLatency: 0,
      lastLatencyMs: 0,
      lastSuccessAt: '',
      lastFailureAt: '',
      operations: []
    }

    const requestsTotal = Number(item.requestsTotal) || 0
    current.requestsTotal += requestsTotal
    current.errorsTotal += Number(item.errorsTotal) || 0
    current.weightedLatency += (Number(item.avgLatencyMs) || 0) * requestsTotal
    current.lastLatencyMs = Math.max(current.lastLatencyMs, Number(item.lastLatencyMs) || 0)
    current.operations.push(item.operation || 'unknown')

    if (toTimestamp(item.lastSuccessAt) > toTimestamp(current.lastSuccessAt)) {
      current.lastSuccessAt = item.lastSuccessAt
    }
    if (toTimestamp(item.lastFailureAt) > toTimestamp(current.lastFailureAt)) {
      current.lastFailureAt = item.lastFailureAt
    }

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

const modelFamilies = computed(() => rawModelFamilies.value.map((item) => ({
  ...item,
  severity: modelSeverity(item)
})))

const abnormalRouteCount = computed(() => routeStatusList.value.filter((item) => item.severity.rank >= severityWeight.warning).length)
const abnormalModelCount = computed(() => modelFamilies.value.filter((item) => item.severity.rank >= severityWeight.warning).length)

const healthStatus = computed(() => {
  const requestsTotal = Number(overview.value.requestsTotal) || 0
  const errorRate = Number(overview.value.errorRate) || 0
  const avgLatencyMs = Number(overview.value.avgLatencyMs) || 0

  if (!requestsTotal) {
    return buildSeverity('neutral', '待观察')
  }
  if (errorRate >= 0.1 || avgLatencyMs >= 4000 || abnormalRouteCount.value >= 3) {
    return buildSeverity('danger', '需重点关注')
  }
  if (errorRate >= 0.03 || avgLatencyMs >= 1500 || abnormalRouteCount.value > 0 || abnormalModelCount.value > 0) {
    return buildSeverity('warning', '轻度波动')
  }
  return buildSeverity('success', '运行稳定')
})

const overviewCards = computed(() => [
  {
    label: '累计请求',
    value: formatNumber(overview.value.requestsTotal),
    desc: '当前监控窗口内累计请求量'
  },
  {
    label: '累计错误',
    value: formatNumber(overview.value.errorsTotal),
    desc: 'HTTP 或业务异常的累计次数'
  },
  {
    label: '错误率',
    value: formatPercent(overview.value.errorRate),
    desc: '快速判断系统是否稳定'
  },
  {
    label: '平均延迟',
    value: formatDuration(overview.value.avgLatencyMs),
    desc: '请求平均耗时'
  }
])

const routeRows = computed(() => [...routeStatusList.value]
  .sort((left, right) => {
    if (right.severity.rank !== left.severity.rank) {
      return right.severity.rank - left.severity.rank
    }
    return (Number(right.requestsTotal) || 0) - (Number(left.requestsTotal) || 0)
  })
  .slice(0, 6))

const modelRows = computed(() => [...modelFamilies.value]
  .sort((left, right) => {
    if (right.severity.rank !== left.severity.rank) {
      return right.severity.rank - left.severity.rank
    }
    return (Number(right.requestsTotal) || 0) - (Number(left.requestsTotal) || 0)
  })
  .slice(0, 6))

const normalizeHistorySamples = (archives = []) => {
  const normalized = [...archives]
    .map((item) => ({
      timestamp: toTimestamp(item.timestamp),
      requestsTotal: Number(item.requestsTotal) || 0,
      avgLatencyMs: Number(item.avgLatencyMs) || 0
    }))
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

const requestTrendBars = computed(() => {
  const peak = requestPeak.value || 1
  return historySamples.value.map((item, index) => ({
    key: `req-${item.timestamp}-${index}`,
    label: formatTimestamp(item.timestamp),
    height: Math.max(10, Math.round((item.requestDelta / peak) * 100))
  }))
})

const latencyTrendBars = computed(() => {
  const peak = latencyPeak.value || 1
  return historySamples.value.map((item, index) => ({
    key: `lat-${item.timestamp}-${index}`,
    label: formatTimestamp(item.timestamp),
    height: Math.max(10, Math.round((item.avgLatencyMs / peak) * 100))
  }))
})

const methodTone = (method) => {
  const value = (method || '').toUpperCase()
  if (value === 'GET') {
    return 'get'
  }
  if (value === 'POST') {
    return 'post'
  }
  return 'default'
}

const isSuccessResponse = (response) => response?.data?.status_code === SUCCESS_CODE
const isForbiddenResponse = (response) => response?.data?.status_code === 3001

const applyMetricsSnapshot = (snapshot) => {
  overview.value = {
    ...overview.value,
    ...(snapshot?.overview || {})
  }
  routes.value = Array.isArray(snapshot?.routes) ? snapshot.routes : []
  models.value = Array.isArray(snapshot?.models) ? snapshot.models : []

  const nextHistory = normalizeHistorySamples(snapshot?.archives || [])
  if (nextHistory.length) {
    historySamples.value = nextHistory
  } else {
    historySamples.value = [
      {
        timestamp: Date.now(),
        requestsTotal: Number(overview.value.requestsTotal) || 0,
        requestDelta: Number(overview.value.requestsTotal) || 0,
        avgLatencyMs: Number(overview.value.avgLatencyMs) || 0
      }
    ]
  }

  lastUpdatedAt.value = Date.now()
}

const fetchSnapshotMetrics = async () => {
  const response = await api.get('/admin/metrics/all')
  if (isForbiddenResponse(response)) {
    const error = new Error('FORBIDDEN')
    error.code = 'FORBIDDEN'
    throw error
  }
  if (!isSuccessResponse(response)) {
    throw new Error('SNAPSHOT_FAILED')
  }
  return response.data?.snapshot || {}
}

const fetchMetrics = async ({ silent = false } = {}) => {
  if (activeRequest) {
    return activeRequest
  }

  if (!lastUpdatedAt.value) {
    loading.value = true
  }
  refreshing.value = true

  activeRequest = fetchSnapshotMetrics()
    .then((snapshot) => {
      applyMetricsSnapshot(snapshot)
    })
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
  if (!autoRefresh.value) {
    return
  }
  refreshTimer = window.setInterval(() => {
    fetchMetrics({ silent: true })
  }, refreshIntervalMs.value)
}

const goMenu = () => {
  router.push('/menu')
}

const goChat = () => {
  router.push('/ai-chat')
}

const logout = async () => {
  const confirmed = await confirmAction({
    title: '退出当前账号？',
    message: '退出后将返回登录页，但不会影响后端已采集的监控数据。',
    confirmText: '退出登录',
    cancelText: '继续查看',
    intent: 'danger'
  })

  if (!confirmed) {
    return
  }

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
  min-height: 100vh;
  position: relative;
  overflow: hidden;
  padding: 28px;
}

.admin-bg {
  position: absolute;
  inset: 0;
  background: linear-gradient(150deg, rgba(255, 255, 255, 0.72), rgba(236, 247, 240, 0.96));
}

.glow {
  position: absolute;
  border-radius: 50%;
  filter: blur(96px);
  opacity: 0.34;
}

.glow-a {
  width: 420px;
  height: 420px;
  top: -120px;
  right: -100px;
  background: rgba(16, 185, 129, 0.28);
}

.glow-b {
  width: 360px;
  height: 360px;
  left: -120px;
  bottom: -140px;
  background: rgba(45, 212, 191, 0.22);
}

.grid {
  position: absolute;
  inset: 0;
  background-image:
    linear-gradient(rgba(16, 185, 129, 0.04) 1px, transparent 1px),
    linear-gradient(90deg, rgba(16, 185, 129, 0.04) 1px, transparent 1px);
  background-size: 48px 48px;
}

.admin-header,
.admin-main {
  position: relative;
  z-index: 1;
}

.admin-header,
.panel-card {
  background: rgba(255, 255, 255, 0.86);
  border: 1px solid rgba(16, 185, 129, 0.14);
  box-shadow: 0 20px 48px rgba(15, 23, 42, 0.08);
  backdrop-filter: blur(18px);
}

.admin-header {
  max-width: 1320px;
  margin: 0 auto 24px;
  padding: 22px 24px;
  border-radius: 24px;
  display: flex;
  justify-content: space-between;
  gap: 18px;
}

.brand-block {
  display: flex;
  gap: 16px;
}

.brand-badge {
  width: 54px;
  height: 54px;
  display: grid;
  place-items: center;
  border-radius: 16px;
  background: linear-gradient(135deg, var(--sci-fi-primary), var(--sci-fi-secondary));
  color: #fff;
  font-family: 'Orbitron', sans-serif;
  font-size: 14px;
  letter-spacing: 2px;
  box-shadow: 0 16px 30px rgba(16, 185, 129, 0.2);
}

.brand-kicker,
.section-kicker,
.metric-label,
.panel-note,
.header-btn,
.refresh-toggle,
.refresh-select,
.method-chip,
.severity-badge {
  font-family: 'Orbitron', sans-serif;
}

.brand-kicker,
.section-kicker {
  margin: 0 0 8px;
  font-size: 11px;
  letter-spacing: 2px;
  color: var(--sci-fi-primary);
}

.brand-title,
.hero-title,
.panel-title {
  margin: 0;
  color: var(--sci-fi-text-primary);
}

.brand-title {
  font-size: 28px;
}

.brand-subtitle,
.hero-desc,
.metric-desc,
.hero-meta,
.panel-note,
.panel-empty,
.path-text,
.cell-subtext {
  color: var(--sci-fi-text-secondary);
}

.brand-subtitle,
.hero-meta,
.cell-subtext {
  margin: 6px 0 0;
  font-size: 14px;
}

.header-actions {
  display: flex;
  align-items: flex-start;
  justify-content: flex-end;
  flex-wrap: wrap;
  gap: 12px;
}

.header-btn,
.refresh-select {
  border: none;
  border-radius: 14px;
  padding: 12px 16px;
  font-size: 11px;
  letter-spacing: 1px;
  color: var(--sci-fi-text-primary);
  background: rgba(16, 185, 129, 0.1);
  transition: transform 0.2s ease, background 0.2s ease, opacity 0.2s ease;
}

.header-btn {
  cursor: pointer;
}

.header-btn.primary {
  color: #fff;
  background: linear-gradient(135deg, var(--sci-fi-primary), var(--sci-fi-secondary));
}

.header-btn.danger {
  color: var(--sci-fi-danger);
  background: rgba(239, 68, 68, 0.1);
}

.header-btn:hover:not(:disabled) {
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
  background: rgba(255, 255, 255, 0.84);
  color: var(--sci-fi-text-secondary);
  font-size: 11px;
  letter-spacing: 1px;
}

.refresh-select {
  min-width: 148px;
  cursor: pointer;
}

.admin-main {
  max-width: 1320px;
  margin: 0 auto;
  display: grid;
  gap: 20px;
}

.panel-card {
  border-radius: 26px;
  padding: 24px;
}

.hero-panel {
  display: flex;
  justify-content: space-between;
  gap: 18px;
}

.hero-title {
  font-size: 34px;
  line-height: 1.25;
}

.hero-desc {
  max-width: 760px;
  margin: 16px 0 0;
  font-size: 15px;
  line-height: 1.8;
}

.hero-side {
  min-width: 260px;
  display: flex;
  flex-direction: column;
  justify-content: center;
  align-items: flex-start;
  gap: 10px;
}

.health-pill,
.severity-badge {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  padding: 8px 14px;
  border-radius: 999px;
  font-size: 11px;
  letter-spacing: 1px;
}

.health-pill.success,
.severity-badge.success {
  background: rgba(16, 185, 129, 0.14);
  color: var(--sci-fi-primary);
}

.health-pill.warning,
.severity-badge.warning {
  background: rgba(245, 158, 11, 0.14);
  color: var(--sci-fi-warning);
}

.health-pill.danger,
.severity-badge.danger {
  background: rgba(239, 68, 68, 0.14);
  color: var(--sci-fi-danger);
}

.health-pill.neutral,
.severity-badge.neutral {
  background: rgba(148, 163, 184, 0.14);
  color: var(--sci-fi-text-secondary);
}

.overview-grid,
.trend-grid,
.watch-grid {
  display: grid;
  gap: 18px;
}

.overview-grid {
  grid-template-columns: repeat(auto-fit, minmax(220px, 1fr));
}

.trend-grid,
.watch-grid {
  grid-template-columns: repeat(2, minmax(0, 1fr));
}

.metric-card {
  min-height: 148px;
}

.metric-label {
  display: block;
  margin-bottom: 12px;
  font-size: 11px;
  letter-spacing: 2px;
  color: var(--sci-fi-text-muted);
}

.metric-value {
  display: block;
  font-size: 28px;
  color: var(--sci-fi-text-primary);
}

.metric-desc {
  margin: 12px 0 0;
  font-size: 14px;
  line-height: 1.7;
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

.trend-bars {
  height: 160px;
  display: flex;
  align-items: flex-end;
  gap: 10px;
  margin-top: 22px;
}

.trend-bar {
  flex: 1;
  min-width: 0;
  height: 100%;
  display: flex;
  align-items: flex-end;
}

.trend-bar-fill {
  width: 100%;
  border-radius: 14px 14px 6px 6px;
  min-height: 12px;
}

.trend-bar-fill.success {
  background: linear-gradient(180deg, rgba(45, 212, 191, 0.9), rgba(16, 185, 129, 0.78));
}

.trend-bar-fill.warning {
  background: linear-gradient(180deg, rgba(245, 158, 11, 0.9), rgba(217, 119, 6, 0.72));
}

.trend-summary {
  margin-top: 18px;
  display: flex;
  justify-content: space-between;
  gap: 12px;
  color: var(--sci-fi-text-secondary);
  font-size: 13px;
}

.panel-empty {
  margin-top: 18px;
  padding: 18px;
  border-radius: 16px;
  background: rgba(255, 255, 255, 0.72);
  font-size: 14px;
  line-height: 1.7;
}

.table-shell {
  margin-top: 18px;
  overflow-x: auto;
}

.metrics-table {
  width: 100%;
  border-collapse: collapse;
}

.metrics-table.compact {
  min-width: 680px;
}

.metrics-table th,
.metrics-table td {
  padding: 14px 12px;
  border-bottom: 1px solid rgba(16, 185, 129, 0.08);
  text-align: left;
  font-size: 14px;
  color: var(--sci-fi-text-primary);
  vertical-align: middle;
}

.metrics-table th {
  font-family: 'Orbitron', sans-serif;
  font-size: 11px;
  letter-spacing: 1px;
  color: var(--sci-fi-text-muted);
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
  letter-spacing: 1px;
}

.method-chip.get {
  background: rgba(16, 185, 129, 0.12);
  color: var(--sci-fi-primary);
}

.method-chip.post {
  background: rgba(59, 130, 246, 0.12);
  color: #2563eb;
}

.method-chip.default {
  background: rgba(148, 163, 184, 0.16);
  color: var(--sci-fi-text-secondary);
}

.path-text {
  word-break: break-all;
}

@media (max-width: 1100px) {
  .admin-header,
  .hero-panel,
  .trend-grid,
  .watch-grid {
    grid-template-columns: 1fr;
    flex-direction: column;
  }

  .header-actions {
    justify-content: flex-start;
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

  .header-actions {
    width: 100%;
  }

  .header-btn,
  .refresh-toggle,
  .refresh-select {
    width: 100%;
  }

  .trend-summary {
    flex-direction: column;
  }
}
</style>
