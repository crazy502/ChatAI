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
          <p class="brand-kicker">OBSERVABILITY</p>
          <h1 class="brand-title">GopherAI 管理监控面板</h1>
          <p class="brand-subtitle">聚焦请求量、错误率、接口健康和模型运行状态的实时视图</p>
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

        <select
          v-model.number="refreshIntervalMs"
          class="refresh-select"
          :disabled="!autoRefresh"
        >
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
          <h2 class="hero-title">先看风险，再下钻到接口和模型细节</h2>
          <p class="hero-desc">
            当前数据来自后端现有 `/admin/metrics/*` 快照接口。页面会在本地维护一个短时观察窗口，用于展示最近刷新周期内的活动变化。
          </p>
        </div>

        <div class="hero-side">
          <div class="health-pill" :class="healthStatus.tone">{{ healthStatus.label }}</div>
          <p class="hero-meta">最近刷新：{{ lastUpdatedLabel }}</p>
          <p class="hero-meta">刷新策略：{{ autoRefresh ? refreshIntervalLabel : '手动刷新' }}</p>
          <p class="hero-meta">当前异常：{{ abnormalRouteCount }} 个接口 / {{ abnormalModelCount }} 个模型项</p>
        </div>
      </section>

      <section class="overview-grid">
        <article v-for="card in overviewCards" :key="card.label" class="metric-card panel-card">
          <span class="metric-label">{{ card.label }}</span>
          <strong class="metric-value">{{ card.value }}</strong>
          <p class="metric-desc">{{ card.desc }}</p>
        </article>
      </section>

      <section class="insight-grid">
        <article class="panel-card insight-panel">
          <div class="panel-head">
            <div>
              <p class="section-kicker">PRIORITY INSIGHTS</p>
              <h3 class="panel-title">优先关注项</h3>
            </div>
            <span class="panel-note">自动提炼当前最值得优先检查的对象</span>
          </div>

          <div v-if="insightCards.length" class="insight-list">
            <article v-for="item in insightCards" :key="item.key" class="insight-item">
              <div class="insight-topline">
                <span class="severity-badge" :class="item.tone">{{ item.level }}</span>
                <span class="insight-kind">{{ item.kind }}</span>
              </div>
              <h4 class="insight-title">{{ item.title }}</h4>
              <p class="insight-desc">{{ item.desc }}</p>
              <p class="insight-meta">{{ item.meta }}</p>
            </article>
          </div>
          <div v-else class="panel-empty">等待系统产生更多监控数据后，这里会自动生成优先级洞察。</div>
        </article>

        <article class="panel-card insight-panel">
          <div class="panel-head">
            <div>
              <p class="section-kicker">MODEL HEALTH</p>
              <h3 class="panel-title">模型健康矩阵</h3>
            </div>
            <span class="panel-note">按模型类型聚合 generate / stream 调用表现</span>
          </div>

          <div v-if="modelFamilies.length" class="family-grid">
            <article v-for="family in modelFamilies" :key="family.modelType" class="family-card">
              <div class="family-head">
                <div>
                  <p class="family-name">{{ family.modelType }}</p>
                  <p class="family-ops">{{ family.operationsLabel }}</p>
                </div>
                <span class="severity-badge" :class="family.severity.tone">{{ family.severity.label }}</span>
              </div>

              <div class="family-metrics">
                <span>请求 {{ formatNumber(family.requestsTotal) }}</span>
                <span>错误率 {{ formatPercent(family.errorRate) }}</span>
                <span>平均 {{ formatDuration(family.avgLatencyMs) }}</span>
              </div>

              <div class="family-track">
                <span class="family-fill" :style="{ width: `${family.loadPercent}%` }"></span>
              </div>

              <div class="family-foot">
                <span>最近成功 {{ formatTimestamp(family.lastSuccessAt) }}</span>
                <span>最近失败 {{ formatTimestamp(family.lastFailureAt) }}</span>
              </div>
            </article>
          </div>
          <div v-else class="panel-empty">当前还没有模型调用数据。</div>
        </article>
      </section>

      <section class="trend-grid">
        <article class="panel-card trend-card">
          <div class="panel-head">
            <div>
              <p class="section-kicker">PAGE WINDOW</p>
              <h3 class="panel-title">请求活跃度</h3>
            </div>
            <span class="panel-note">按本页刷新周期统计请求增量</span>
          </div>

          <div v-if="requestTrendBars.length" class="trend-bars">
            <div
              v-for="bar in requestTrendBars"
              :key="bar.key"
              class="trend-bar"
              :title="`${bar.label}：${bar.value}`"
            >
              <span class="trend-bar-fill success" :style="{ height: `${bar.height}%` }"></span>
            </div>
          </div>
          <div v-else class="panel-empty">等待采样中，刷新几次后会显示请求活动走势。</div>

          <div class="trend-summary">
            <span>当前增量 {{ formatNumber(latestRequestDelta) }}</span>
            <span>峰值 {{ formatNumber(requestPeak) }}</span>
          </div>
        </article>

        <article class="panel-card trend-card">
          <div class="panel-head">
            <div>
              <p class="section-kicker">LATENCY WINDOW</p>
              <h3 class="panel-title">平均延迟走势</h3>
            </div>
            <span class="panel-note">观察本页周期内的平均响应时间</span>
          </div>

          <div v-if="latencyTrendBars.length" class="trend-bars">
            <div
              v-for="bar in latencyTrendBars"
              :key="bar.key"
              class="trend-bar"
              :title="`${bar.label}：${formatDuration(bar.value)}`"
            >
              <span class="trend-bar-fill warning" :style="{ height: `${bar.height}%` }"></span>
            </div>
          </div>
          <div v-else class="panel-empty">等待采样中，刷新几次后会显示延迟变化。</div>

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
              <p class="section-kicker">ROUTE WATCHLIST</p>
              <h3 class="panel-title">接口重点监控</h3>
            </div>
            <span class="panel-note">优先展示异常接口，没有异常时展示最繁忙接口</span>
          </div>

          <div v-if="routeWatchlist.length" class="watch-list">
            <article v-for="route in routeWatchlist" :key="`${route.method}-${route.path}`" class="watch-item">
              <div class="watch-head">
                <div class="route-cell">
                  <span class="method-chip" :class="methodTone(route.method)">{{ route.method }}</span>
                  <span class="path-text">{{ route.path }}</span>
                </div>
                <span class="severity-badge" :class="route.severity.tone">{{ route.severity.label }}</span>
              </div>

              <div class="watch-meta">
                <span>请求 {{ formatNumber(route.requestsTotal) }}</span>
                <span>错误率 {{ formatPercent(route.errorRate) }}</span>
                <span>平均 {{ formatDuration(route.avgLatencyMs) }}</span>
              </div>

              <div class="watch-track">
                <span class="watch-fill" :style="{ width: `${route.loadPercent}%` }"></span>
              </div>
            </article>
          </div>
          <div v-else class="panel-empty">当前还没有接口热点数据。</div>
        </article>

        <article class="panel-card">
          <div class="panel-head">
            <div>
              <p class="section-kicker">MODEL WATCHLIST</p>
              <h3 class="panel-title">模型重点监控</h3>
            </div>
            <span class="panel-note">优先展示异常模型，没有异常时展示负载最高模型</span>
          </div>

          <div v-if="modelWatchlist.length" class="watch-list">
            <article v-for="model in modelWatchlist" :key="model.modelType" class="watch-item">
              <div class="watch-head">
                <div>
                  <p class="family-name">{{ model.modelType }}</p>
                  <p class="family-ops">{{ model.operationsLabel }}</p>
                </div>
                <span class="severity-badge" :class="model.severity.tone">{{ model.severity.label }}</span>
              </div>

              <div class="watch-meta">
                <span>请求 {{ formatNumber(model.requestsTotal) }}</span>
                <span>错误率 {{ formatPercent(model.errorRate) }}</span>
                <span>平均 {{ formatDuration(model.avgLatencyMs) }}</span>
              </div>

              <div class="watch-track">
                <span class="watch-fill warning" :style="{ width: `${model.loadPercent}%` }"></span>
              </div>
            </article>
          </div>
          <div v-else class="panel-empty">当前还没有模型热点数据。</div>
        </article>
      </section>

      <section class="summary-grid">
        <article v-for="card in summaryCards" :key="card.label" class="summary-card panel-card">
          <span class="metric-label">{{ card.label }}</span>
          <strong class="summary-value">{{ card.value }}</strong>
          <p class="metric-desc">{{ card.desc }}</p>
        </article>
      </section>

      <section class="watch-grid">
        <article class="panel-card">
          <div class="panel-head panel-head-wrap">
            <div>
              <p class="section-kicker">USER BREAKDOWN</p>
              <h3 class="panel-title">按用户拆分</h3>
            </div>

            <div class="panel-controls">
              <input
                v-model="userQuery"
                class="panel-input"
                type="search"
                placeholder="搜索用户名"
              />
            </div>
          </div>

          <div v-if="loading && !users.length" class="table-loading">正在读取用户维度...</div>
          <div v-else-if="!filteredUsers.length" class="panel-empty">当前没有符合条件的用户维度数据。</div>
          <div v-else class="table-shell">
            <table class="metrics-table">
              <thead>
                <tr>
                  <th>用户</th>
                  <th>请求量</th>
                  <th>错误数</th>
                  <th>错误率</th>
                  <th>平均延迟</th>
                  <th>最近访问</th>
                </tr>
              </thead>
              <tbody>
                <tr v-for="user in filteredUsers" :key="user.userName">
                  <td>{{ user.userName || 'anonymous' }}</td>
                  <td>{{ formatNumber(user.requestsTotal) }}</td>
                  <td>{{ formatNumber(user.errorsTotal) }}</td>
                  <td>{{ formatPercent(user.errorRate) }}</td>
                  <td>{{ formatDuration(user.avgLatencyMs) }}</td>
                  <td>{{ formatTimestamp(user.lastSeenAt) }}</td>
                </tr>
              </tbody>
            </table>
          </div>
        </article>

        <article class="panel-card">
          <div class="panel-head panel-head-wrap">
            <div>
              <p class="section-kicker">BUSINESS CODE</p>
              <h3 class="panel-title">按接口业务码拆分</h3>
            </div>

            <div class="panel-controls">
              <input
                v-model="businessCodeQuery"
                class="panel-input"
                type="search"
                placeholder="搜索路径或业务码"
              />
            </div>
          </div>

          <div v-if="loading && !businessCodes.length" class="table-loading">正在读取业务码维度...</div>
          <div v-else-if="!filteredBusinessCodes.length" class="panel-empty">当前没有符合条件的业务码数据。</div>
          <div v-else class="table-shell">
            <table class="metrics-table">
              <thead>
                <tr>
                  <th>接口</th>
                  <th>业务码</th>
                  <th>请求量</th>
                  <th>错误率</th>
                  <th>平均延迟</th>
                  <th>最近 HTTP</th>
                </tr>
              </thead>
              <tbody>
                <tr v-for="item in filteredBusinessCodes" :key="`${item.method}-${item.path}-${item.businessCode}`">
                  <td>
                    <div class="route-cell">
                      <span class="method-chip" :class="methodTone(item.method)">{{ item.method }}</span>
                      <span class="path-text">{{ item.path }}</span>
                    </div>
                  </td>
                  <td>{{ formatBusinessCode(item.businessCode) }}</td>
                  <td>{{ formatNumber(item.requestsTotal) }}</td>
                  <td>{{ formatPercent(item.errorRate) }}</td>
                  <td>{{ formatDuration(item.avgLatencyMs) }}</td>
                  <td>{{ item.lastHttpStatus }}</td>
                </tr>
              </tbody>
            </table>
          </div>
        </article>
      </section>

      <section class="panel-card">
        <div class="panel-head panel-head-wrap">
          <div>
            <p class="section-kicker">MODEL FAILURE DETAILS</p>
            <h3 class="panel-title">按模型操作拆分失败明细</h3>
          </div>

          <div class="panel-controls">
            <input
              v-model="failureQuery"
              class="panel-input"
              type="search"
              placeholder="搜索模型、操作、用户或错误"
            />
          </div>
        </div>

        <div v-if="loading && !modelFailures.length" class="table-loading">正在读取模型失败明细...</div>
        <div v-else-if="!filteredModelFailures.length" class="panel-empty">当前没有模型失败明细，说明最近窗口内模型调用比较稳定。</div>
        <div v-else class="table-shell">
          <table class="metrics-table">
            <thead>
              <tr>
                <th>发生时间</th>
                <th>模型</th>
                <th>操作</th>
                <th>用户</th>
                <th>失败耗时</th>
                <th>错误信息</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="failure in filteredModelFailures" :key="`${failure.modelType}-${failure.operation}-${failure.occurredAt}-${failure.userName}`">
                <td>{{ formatTimestamp(failure.occurredAt) }}</td>
                <td>{{ failure.modelType || '-' }}</td>
                <td>{{ failure.operation || '-' }}</td>
                <td>{{ failure.userName || 'anonymous' }}</td>
                <td>{{ formatDuration(failure.latencyMs) }}</td>
                <td>{{ failure.errorMessage || '-' }}</td>
              </tr>
            </tbody>
          </table>
        </div>
      </section>

      <section class="panel-card">
        <div class="panel-head panel-head-wrap">
          <div>
            <p class="section-kicker">ROUTE METRICS</p>
            <h3 class="panel-title">接口表现</h3>
          </div>

          <div class="panel-controls">
            <input
              v-model="routeQuery"
              class="panel-input"
              type="search"
              placeholder="搜索接口路径"
            />

            <select v-model="routeFocus" class="panel-select">
              <option value="all">全部接口</option>
              <option value="attention">只看需关注</option>
              <option value="danger">只看高风险</option>
            </select>

            <select v-model="routeSortKey" class="panel-select">
              <option value="requests">按请求量</option>
              <option value="latency">按平均延迟</option>
              <option value="errors">按错误数</option>
              <option value="errorRate">按错误率</option>
            </select>
          </div>
        </div>

        <div v-if="loading && !routes.length" class="table-loading">正在读取接口指标...</div>
        <div v-else-if="!filteredRoutes.length" class="panel-empty">当前没有符合筛选条件的接口指标。</div>
        <div v-else class="table-shell">
          <table class="metrics-table">
            <thead>
              <tr>
                <th>接口</th>
                <th>健康</th>
                <th>请求量</th>
                <th>错误数</th>
                <th>错误率</th>
                <th>平均延迟</th>
                <th>最近延迟</th>
                <th>最近状态</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="route in filteredRoutes" :key="`${route.method}-${route.path}`">
                <td>
                  <div class="route-cell">
                    <span class="method-chip" :class="methodTone(route.method)">{{ route.method }}</span>
                    <span class="path-text">{{ route.path }}</span>
                  </div>
                </td>
                <td>
                  <span class="severity-badge" :class="routeSeverity(route).tone">{{ routeSeverity(route).label }}</span>
                </td>
                <td>{{ formatNumber(route.requestsTotal) }}</td>
                <td>{{ formatNumber(route.errorsTotal) }}</td>
                <td>{{ formatPercent(route.errorRate) }}</td>
                <td>{{ formatDuration(route.avgLatencyMs) }}</td>
                <td>{{ formatDuration(route.lastLatencyMs) }}</td>
                <td>{{ route.lastHttpStatus }} / {{ route.lastBusinessCode }}</td>
              </tr>
            </tbody>
          </table>
        </div>
      </section>

      <section class="panel-card">
        <div class="panel-head panel-head-wrap">
          <div>
            <p class="section-kicker">MODEL METRICS</p>
            <h3 class="panel-title">模型调用表现</h3>
          </div>

          <div class="panel-controls">
            <input
              v-model="modelQuery"
              class="panel-input"
              type="search"
              placeholder="搜索模型或操作"
            />

            <select v-model="modelFocus" class="panel-select">
              <option value="all">全部模型项</option>
              <option value="attention">只看需关注</option>
              <option value="danger">只看高风险</option>
            </select>

            <select v-model="modelSortKey" class="panel-select">
              <option value="requests">按请求量</option>
              <option value="latency">按平均延迟</option>
              <option value="errors">按错误数</option>
              <option value="errorRate">按错误率</option>
            </select>
          </div>
        </div>

        <div v-if="loading && !models.length" class="table-loading">正在读取模型指标...</div>
        <div v-else-if="!filteredModels.length" class="panel-empty">当前没有符合筛选条件的模型指标。</div>
        <div v-else class="table-shell">
          <table class="metrics-table">
            <thead>
              <tr>
                <th>模型</th>
                <th>操作</th>
                <th>健康</th>
                <th>请求量</th>
                <th>错误率</th>
                <th>平均延迟</th>
                <th>最近延迟</th>
                <th>最近成功</th>
                <th>最近失败</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="model in filteredModels" :key="`${model.modelType}-${model.operation}`">
                <td>{{ model.modelType || '-' }}</td>
                <td>{{ model.operation || '-' }}</td>
                <td>
                  <span class="severity-badge" :class="modelSeverity(model).tone">{{ modelSeverity(model).label }}</span>
                </td>
                <td>{{ formatNumber(model.requestsTotal) }}</td>
                <td>{{ formatPercent(model.errorRate) }}</td>
                <td>{{ formatDuration(model.avgLatencyMs) }}</td>
                <td>{{ formatDuration(model.lastLatencyMs) }}</td>
                <td>{{ formatTimestamp(model.lastSuccessAt) }}</td>
                <td>{{ formatTimestamp(model.lastFailureAt) }}</td>
              </tr>
            </tbody>
          </table>
        </div>
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

const DEFAULT_REFRESH_INTERVAL = 10000
const SUCCESS_CODE = 1000
const MAX_HISTORY_SAMPLES = 18
const refreshOptions = [
  { label: '10 秒刷新', value: 10000 },
  { label: '30 秒刷新', value: 30000 },
  { label: '60 秒刷新', value: 60000 }
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
  uptimeSeconds: 0,
  requestsTotal: 0,
  errorsTotal: 0,
  errorRate: 0,
  avgLatencyMs: 0,
  routesTracked: 0,
  modelsTracked: 0,
  usersTracked: 0,
  businessCodesTracked: 0,
  recentFailuresTracked: 0
})
const routes = ref([])
const models = ref([])
const users = ref([])
const businessCodes = ref([])
const modelFailures = ref([])
const loading = ref(true)
const refreshing = ref(false)
const autoRefresh = ref(true)
const refreshIntervalMs = ref(DEFAULT_REFRESH_INTERVAL)
const lastUpdatedAt = ref(0)
const routeQuery = ref('')
const routeSortKey = ref('requests')
const routeFocus = ref('all')
const modelQuery = ref('')
const modelSortKey = ref('requests')
const modelFocus = ref('all')
const userQuery = ref('')
const businessCodeQuery = ref('')
const failureQuery = ref('')
const historySamples = ref([])

let refreshTimer = null
let activeRequest = null

const numberFormatter = new Intl.NumberFormat('zh-CN')

const formatNumber = (value) => numberFormatter.format(Number(value) || 0)

const formatPercent = (value) => `${((Number(value) || 0) * 100).toFixed(1)}%`

const formatBusinessCode = (value) => {
  const normalized = Number(value) || 0
  return normalized === 0 ? '0 / 未解析' : String(normalized)
}

const formatDuration = (value) => {
  const duration = Number(value) || 0
  if (duration >= 1000) {
    return `${(duration / 1000).toFixed(duration >= 10000 ? 0 : 1)} s`
  }
  return `${Math.round(duration)} ms`
}

const formatUptime = (seconds) => {
  const totalSeconds = Math.max(0, Math.floor(Number(seconds) || 0))
  const days = Math.floor(totalSeconds / 86400)
  const hours = Math.floor((totalSeconds % 86400) / 3600)
  const minutes = Math.floor((totalSeconds % 3600) / 60)
  const remainingSeconds = totalSeconds % 60

  if (days > 0) {
    return `${days} 天 ${String(hours).padStart(2, '0')} 小时`
  }
  return `${String(hours).padStart(2, '0')}:${String(minutes).padStart(2, '0')}:${String(remainingSeconds).padStart(2, '0')}`
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

  return new Date(timestamp).toLocaleString('zh-CN', {
    month: '2-digit',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit'
  })
}

const formatRefreshInterval = (value) => {
  const seconds = Math.floor((Number(value) || 0) / 1000)
  return `每 ${seconds} 秒自动刷新`
}

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
  const lastBusinessCode = Number(route?.lastBusinessCode) || 0

  if (!requestsTotal) {
    return buildSeverity('neutral', '待采样')
  }
  if (lastHttpStatus >= 500 || errorRate >= 0.1 || errorsTotal >= 5) {
    return buildSeverity('danger', '高风险')
  }
  if (errorRate >= 0.03 || avgLatencyMs >= 1500 || lastLatencyMs >= 2500 || (lastBusinessCode !== 0 && lastBusinessCode !== SUCCESS_CODE)) {
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

const matchesFocus = (severity, focus) => {
  if (focus === 'danger') {
    return severity.tone === 'danger'
  }
  if (focus === 'attention') {
    return severity.rank >= severityWeight.warning
  }
  return true
}

const compareBySeverity = (leftSeverity, rightSeverity) => rightSeverity.rank - leftSeverity.rank

const appendHistorySample = (nextOverview) => {
  const previous = historySamples.value[historySamples.value.length - 1]
  const requestsTotal = Number(nextOverview.requestsTotal) || 0
  const requestDelta = previous ? Math.max(0, requestsTotal - previous.requestsTotal) : requestsTotal

  const nextSample = {
    timestamp: Date.now(),
    requestsTotal,
    requestDelta,
    avgLatencyMs: Number(nextOverview.avgLatencyMs) || 0,
    errorRate: Number(nextOverview.errorRate) || 0
  }

  historySamples.value = [...historySamples.value, nextSample].slice(-MAX_HISTORY_SAMPLES)
}

const refreshIntervalLabel = computed(() => formatRefreshInterval(refreshIntervalMs.value))
const lastUpdatedLabel = computed(() => relativeTime(lastUpdatedAt.value))

const requestPeak = computed(() => Math.max(0, ...historySamples.value.map((item) => item.requestDelta)))
const latencyPeak = computed(() => Math.max(0, ...historySamples.value.map((item) => item.avgLatencyMs)))
const latestRequestDelta = computed(() => historySamples.value[historySamples.value.length - 1]?.requestDelta || 0)

const requestTrendBars = computed(() => {
  const peak = requestPeak.value || 1
  return historySamples.value.map((item, index) => ({
    key: `req-${item.timestamp}-${index}`,
    label: formatTimestamp(item.timestamp),
    value: item.requestDelta,
    height: Math.max(10, Math.round((item.requestDelta / peak) * 100))
  }))
})

const latencyTrendBars = computed(() => {
  const peak = latencyPeak.value || 1
  return historySamples.value.map((item, index) => ({
    key: `lat-${item.timestamp}-${index}`,
    label: formatTimestamp(item.timestamp),
    value: item.avgLatencyMs,
    height: Math.max(10, Math.round((item.avgLatencyMs / peak) * 100))
  }))
})

const routeStatusList = computed(() => routes.value.map((route) => ({
  ...route,
  severity: routeSeverity(route)
})))

const modelStatusList = computed(() => models.value.map((model) => ({
  ...model,
  severity: modelSeverity(model)
})))

const abnormalRouteCount = computed(() => routeStatusList.value.filter((item) => item.severity.rank >= severityWeight.warning).length)
const abnormalModelCount = computed(() => modelStatusList.value.filter((item) => item.severity.rank >= severityWeight.warning).length)

const healthStatus = computed(() => {
  const errorRate = Number(overview.value.errorRate) || 0
  const avgLatency = Number(overview.value.avgLatencyMs) || 0
  const requestsTotal = Number(overview.value.requestsTotal) || 0

  if (!requestsTotal) {
    return buildSeverity('neutral', '待观测')
  }
  if (errorRate >= 0.1 || avgLatency >= 4000 || abnormalRouteCount.value >= 3) {
    return buildSeverity('danger', '需重点关注')
  }
  if (errorRate >= 0.03 || avgLatency >= 1500 || abnormalRouteCount.value > 0 || abnormalModelCount.value > 0) {
    return buildSeverity('warning', '轻度波动')
  }
  return buildSeverity('success', '运行稳定')
})

const overviewCards = computed(() => [
  {
    label: '累计请求',
    value: formatNumber(overview.value.requestsTotal),
    desc: '从服务启动至今累计记录的请求数'
  },
  {
    label: '累计错误',
    value: formatNumber(overview.value.errorsTotal),
    desc: '包含 HTTP 错误与业务状态异常'
  },
  {
    label: '错误率',
    value: formatPercent(overview.value.errorRate),
    desc: '用于快速判断系统稳定性变化'
  },
  {
    label: '平均延迟',
    value: formatDuration(overview.value.avgLatencyMs),
    desc: '所有已采集请求的平均耗时'
  },
  {
    label: '异常接口数',
    value: formatNumber(abnormalRouteCount.value),
    desc: '健康等级为“需关注”或“高风险”的接口数'
  },
  {
    label: '异常模型项',
    value: formatNumber(abnormalModelCount.value),
    desc: '健康等级为“需关注”或“高风险”的模型调用项'
  },
  {
    label: '服务运行时长',
    value: formatUptime(overview.value.uptimeSeconds),
    desc: '从本次服务启动开始累计'
  },
  {
    label: '监控覆盖',
    value: `${formatNumber(overview.value.routesTracked)} / ${formatNumber(overview.value.modelsTracked)}`,
    desc: '当前已被采集的接口数 / 模型操作数'
  },
  {
    label: '用户维度',
    value: formatNumber(overview.value.usersTracked),
    desc: '当前已进入监控的用户身份数'
  },
  {
    label: '业务码桶',
    value: formatNumber(overview.value.businessCodesTracked),
    desc: '按接口业务码拆分后的统计桶数量'
  },
  {
    label: '最近失败样本',
    value: formatNumber(overview.value.recentFailuresTracked),
    desc: '当前保留的模型失败明细条数'
  }
])

const topBusyRoute = computed(() => [...routeStatusList.value].sort((a, b) => b.requestsTotal - a.requestsTotal)[0] || null)
const topSlowRoute = computed(() => [...routeStatusList.value].sort((a, b) => b.avgLatencyMs - a.avgLatencyMs)[0] || null)
const topErrorRoute = computed(() => [...routeStatusList.value].sort((a, b) => {
  if (b.errorsTotal === a.errorsTotal) {
    return b.errorRate - a.errorRate
  }
  return b.errorsTotal - a.errorsTotal
})[0] || null)
const topActiveUser = computed(() => [...users.value].sort((a, b) => {
  if ((Number(b.requestsTotal) || 0) === (Number(a.requestsTotal) || 0)) {
    return (Number(b.errorsTotal) || 0) - (Number(a.errorsTotal) || 0)
  }
  return (Number(b.requestsTotal) || 0) - (Number(a.requestsTotal) || 0)
})[0] || null)
const hottestBusinessCode = computed(() => [...businessCodes.value].sort((a, b) => {
  if ((Number(b.errorsTotal) || 0) === (Number(a.errorsTotal) || 0)) {
    return (Number(b.requestsTotal) || 0) - (Number(a.requestsTotal) || 0)
  }
  return (Number(b.errorsTotal) || 0) - (Number(a.errorsTotal) || 0)
})[0] || null)
const latestModelFailure = computed(() => modelFailures.value[0] || null)

const rawModelFamilies = computed(() => {
  const familyMap = new Map()

  for (const item of modelStatusList.value) {
    const modelType = item.modelType || 'unknown'
    const current = familyMap.get(modelType) || {
      modelType,
      requestsTotal: 0,
      errorsTotal: 0,
      weightedLatency: 0,
      operations: [],
      lastSuccessAt: '',
      lastFailureAt: '',
      lastLatencyMs: 0
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
    const avgLatencyMs = requestsTotal ? item.weightedLatency / requestsTotal : 0
    const operations = [...new Set(item.operations)]
    return {
      modelType: item.modelType,
      requestsTotal,
      errorsTotal: item.errorsTotal,
      errorRate: requestsTotal ? item.errorsTotal / requestsTotal : 0,
      avgLatencyMs,
      lastLatencyMs: item.lastLatencyMs,
      lastSuccessAt: item.lastSuccessAt,
      lastFailureAt: item.lastFailureAt,
      operationsLabel: operations.join(' / ')
    }
  })
})

const maxModelFamilyRequests = computed(() => Math.max(1, ...rawModelFamilies.value.map((item) => Number(item.requestsTotal) || 0)))

const modelFamilies = computed(() => rawModelFamilies.value
  .map((family) => ({
    ...family,
    severity: modelSeverity(family),
    loadPercent: Math.max(10, Math.round(((Number(family.requestsTotal) || 0) / maxModelFamilyRequests.value) * 100))
  }))
  .sort((left, right) => {
    const severityDiff = compareBySeverity(left.severity, right.severity)
    if (severityDiff !== 0) {
      return severityDiff
    }
    return (Number(right.requestsTotal) || 0) - (Number(left.requestsTotal) || 0)
  }))

const topBusyModelFamily = computed(() => modelFamilies.value[0] || null)

const abnormalRoutes = computed(() => [...routeStatusList.value]
  .filter((item) => item.severity.rank >= severityWeight.warning)
  .sort((left, right) => {
    const severityDiff = compareBySeverity(left.severity, right.severity)
    if (severityDiff !== 0) {
      return severityDiff
    }
    if (right.errorRate !== left.errorRate) {
      return right.errorRate - left.errorRate
    }
    return (Number(right.avgLatencyMs) || 0) - (Number(left.avgLatencyMs) || 0)
  }))

const abnormalModelFamilies = computed(() => modelFamilies.value.filter((item) => item.severity.rank >= severityWeight.warning))

const summaryCards = computed(() => [
  {
    label: '最繁忙接口',
    value: topBusyRoute.value ? `${topBusyRoute.value.method} ${topBusyRoute.value.path}` : '暂无数据',
    desc: topBusyRoute.value ? `累计 ${formatNumber(topBusyRoute.value.requestsTotal)} 次请求` : '等待接口流量进入监控窗口'
  },
  {
    label: '最慢接口',
    value: topSlowRoute.value ? `${topSlowRoute.value.method} ${topSlowRoute.value.path}` : '暂无数据',
    desc: topSlowRoute.value ? `平均 ${formatDuration(topSlowRoute.value.avgLatencyMs)}` : '等待延迟数据进入监控窗口'
  },
  {
    label: '错误最多接口',
    value: topErrorRoute.value ? `${topErrorRoute.value.method} ${topErrorRoute.value.path}` : '暂无数据',
    desc: topErrorRoute.value ? `错误率 ${formatPercent(topErrorRoute.value.errorRate)}` : '等待异常数据进入监控窗口'
  },
  {
    label: '最繁忙模型族',
    value: topBusyModelFamily.value ? topBusyModelFamily.value.modelType : '暂无数据',
    desc: topBusyModelFamily.value ? `累计 ${formatNumber(topBusyModelFamily.value.requestsTotal)} 次调用` : '等待模型调用进入监控窗口'
  },
  {
    label: '最活跃用户',
    value: topActiveUser.value ? topActiveUser.value.userName : '暂无数据',
    desc: topActiveUser.value ? `累计 ${formatNumber(topActiveUser.value.requestsTotal)} 次请求` : '等待用户维度数据进入监控窗口'
  },
  {
    label: '热点业务码',
    value: hottestBusinessCode.value ? `${hottestBusinessCode.value.method} ${hottestBusinessCode.value.path}` : '暂无数据',
    desc: hottestBusinessCode.value ? `业务码 ${formatBusinessCode(hottestBusinessCode.value.businessCode)} · ${formatNumber(hottestBusinessCode.value.requestsTotal)} 次` : '等待业务码维度数据进入监控窗口'
  },
  {
    label: '最近模型失败',
    value: latestModelFailure.value ? `${latestModelFailure.value.modelType} / ${latestModelFailure.value.operation}` : '暂无失败记录',
    desc: latestModelFailure.value ? `${latestModelFailure.value.userName || 'anonymous'} · ${formatTimestamp(latestModelFailure.value.occurredAt)}` : '当前还没有模型失败样本'
  }
])

const createRouteInsight = (route, kind, desc, meta) => ({
  key: `route-${kind}-${route.method}-${route.path}`,
  tone: route.severity.tone,
  level: route.severity.label,
  kind,
  title: `${route.method} ${route.path}`,
  desc,
  meta
})

const createModelInsight = (family, kind, desc, meta) => ({
  key: `model-${kind}-${family.modelType}`,
  tone: family.severity.tone,
  level: family.severity.label,
  kind,
  title: family.modelType,
  desc,
  meta
})

const insightCards = computed(() => {
  const cards = []
  const seen = new Set()
  const pushCard = (card) => {
    if (!card || seen.has(card.key)) {
      return
    }
    seen.add(card.key)
    cards.push(card)
  }

  if (!(Number(overview.value.requestsTotal) || 0)) {
    pushCard({
      key: 'system-pending',
      tone: 'neutral',
      level: '待采样',
      kind: '系统',
      title: '等待监控窗口积累数据',
      desc: '先触发登录、会话同步或聊天请求，这里就会开始自动提炼优先级洞察。',
      meta: '当前适合先验证基础链路是否畅通。'
    })
    return cards
  }

  if (healthStatus.value.tone === 'danger' || healthStatus.value.tone === 'warning') {
    pushCard({
      key: 'system-health',
      tone: healthStatus.value.tone,
      level: healthStatus.value.label,
      kind: '系统',
      title: '整体健康状态出现波动',
      desc: `当前错误率 ${formatPercent(overview.value.errorRate)}，平均延迟 ${formatDuration(overview.value.avgLatencyMs)}。`,
      meta: `异常接口 ${formatNumber(abnormalRouteCount.value)} 个，异常模型项 ${formatNumber(abnormalModelCount.value)} 个。`
    })
  }

  const criticalRoute = abnormalRoutes.value[0]
  if (criticalRoute) {
    pushCard(createRouteInsight(
      criticalRoute,
      '接口',
      `错误率 ${formatPercent(criticalRoute.errorRate)}，平均延迟 ${formatDuration(criticalRoute.avgLatencyMs)}。`,
      `最近状态 ${criticalRoute.lastHttpStatus} / ${criticalRoute.lastBusinessCode}。`
    ))
  }

  const slowRoute = [...routeStatusList.value]
    .sort((left, right) => (Number(right.avgLatencyMs) || 0) - (Number(left.avgLatencyMs) || 0))
    .find((item) => item.severity.rank >= severityWeight.warning)
  if (slowRoute) {
    pushCard(createRouteInsight(
      slowRoute,
      '延迟',
      `平均延迟 ${formatDuration(slowRoute.avgLatencyMs)}，最近一次 ${formatDuration(slowRoute.lastLatencyMs)}。`,
      `累计请求 ${formatNumber(slowRoute.requestsTotal)} 次。`
    ))
  }

  const criticalModel = abnormalModelFamilies.value[0]
  if (criticalModel) {
    pushCard(createModelInsight(
      criticalModel,
      '模型',
      `错误率 ${formatPercent(criticalModel.errorRate)}，平均延迟 ${formatDuration(criticalModel.avgLatencyMs)}。`,
      `覆盖操作 ${criticalModel.operationsLabel || '暂无'}。`
    ))
  }

  if (!cards.length) {
    pushCard({
      key: 'system-stable',
      tone: 'success',
      level: '稳定',
      kind: '系统',
      title: '当前没有明显异常热点',
      desc: '接口和模型指标目前都处于相对平稳区间，可以继续观察趋势变化。',
      meta: `当前累计请求 ${formatNumber(overview.value.requestsTotal)} 次。`
    })
  }

  return cards.slice(0, 4)
})

const routeWatchlist = computed(() => {
  const source = abnormalRoutes.value.length
    ? abnormalRoutes.value
    : [...routeStatusList.value].sort((left, right) => (Number(right.requestsTotal) || 0) - (Number(left.requestsTotal) || 0))

  const maxRequests = Math.max(1, ...source.map((item) => Number(item.requestsTotal) || 0))
  return source.slice(0, 5).map((item) => ({
    ...item,
    loadPercent: Math.max(10, Math.round(((Number(item.requestsTotal) || 0) / maxRequests) * 100))
  }))
})

const modelWatchlist = computed(() => {
  const source = abnormalModelFamilies.value.length
    ? abnormalModelFamilies.value
    : modelFamilies.value

  const maxRequests = Math.max(1, ...source.map((item) => Number(item.requestsTotal) || 0))
  return source.slice(0, 5).map((item) => ({
    ...item,
    loadPercent: Math.max(10, Math.round(((Number(item.requestsTotal) || 0) / maxRequests) * 100))
  }))
})

const normalizedRouteQuery = computed(() => routeQuery.value.trim().toLowerCase())
const normalizedModelQuery = computed(() => modelQuery.value.trim().toLowerCase())

const sortRoutes = (routeList) => {
  const list = [...routeList]
  switch (routeSortKey.value) {
    case 'latency':
      return list.sort((left, right) => (Number(right.avgLatencyMs) || 0) - (Number(left.avgLatencyMs) || 0))
    case 'errors':
      return list.sort((left, right) => (Number(right.errorsTotal) || 0) - (Number(left.errorsTotal) || 0))
    case 'errorRate':
      return list.sort((left, right) => (Number(right.errorRate) || 0) - (Number(left.errorRate) || 0))
    case 'requests':
    default:
      return list.sort((left, right) => (Number(right.requestsTotal) || 0) - (Number(left.requestsTotal) || 0))
  }
}

const sortModels = (modelList) => {
  const list = [...modelList]
  switch (modelSortKey.value) {
    case 'latency':
      return list.sort((left, right) => (Number(right.avgLatencyMs) || 0) - (Number(left.avgLatencyMs) || 0))
    case 'errors':
      return list.sort((left, right) => (Number(right.errorsTotal) || 0) - (Number(left.errorsTotal) || 0))
    case 'errorRate':
      return list.sort((left, right) => (Number(right.errorRate) || 0) - (Number(left.errorRate) || 0))
    case 'requests':
    default:
      return list.sort((left, right) => (Number(right.requestsTotal) || 0) - (Number(left.requestsTotal) || 0))
  }
}

const filteredRoutes = computed(() => {
  const query = normalizedRouteQuery.value
  const source = query
    ? routeStatusList.value.filter((route) => `${route.method} ${route.path}`.toLowerCase().includes(query))
    : routeStatusList.value

  return sortRoutes(source.filter((route) => matchesFocus(route.severity, routeFocus.value)))
})

const filteredModels = computed(() => {
  const query = normalizedModelQuery.value
  const source = query
    ? modelStatusList.value.filter((item) => `${item.modelType} ${item.operation}`.toLowerCase().includes(query))
    : modelStatusList.value

  return sortModels(source.filter((item) => matchesFocus(item.severity, modelFocus.value)))
})

const normalizedUserQuery = computed(() => userQuery.value.trim().toLowerCase())
const normalizedBusinessCodeQuery = computed(() => businessCodeQuery.value.trim().toLowerCase())
const normalizedFailureQuery = computed(() => failureQuery.value.trim().toLowerCase())

const filteredUsers = computed(() => {
  const query = normalizedUserQuery.value
  const source = query
    ? users.value.filter((item) => `${item.userName}`.toLowerCase().includes(query))
    : users.value

  return [...source].sort((left, right) => {
    if ((Number(right.requestsTotal) || 0) === (Number(left.requestsTotal) || 0)) {
      return (Number(right.errorsTotal) || 0) - (Number(left.errorsTotal) || 0)
    }
    return (Number(right.requestsTotal) || 0) - (Number(left.requestsTotal) || 0)
  })
})

const filteredBusinessCodes = computed(() => {
  const query = normalizedBusinessCodeQuery.value
  const source = query
    ? businessCodes.value.filter((item) => `${item.method} ${item.path} ${item.businessCode}`.toLowerCase().includes(query))
    : businessCodes.value

  return [...source].sort((left, right) => {
    if ((Number(right.errorsTotal) || 0) === (Number(left.errorsTotal) || 0)) {
      return (Number(right.requestsTotal) || 0) - (Number(left.requestsTotal) || 0)
    }
    return (Number(right.errorsTotal) || 0) - (Number(left.errorsTotal) || 0)
  })
})

const filteredModelFailures = computed(() => {
  const query = normalizedFailureQuery.value
  const source = query
    ? modelFailures.value.filter((item) => `${item.modelType} ${item.operation} ${item.userName} ${item.errorMessage}`.toLowerCase().includes(query))
    : modelFailures.value

  return [...source].sort((left, right) => toTimestamp(right.occurredAt) - toTimestamp(left.occurredAt))
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

const fetchMetrics = async ({ silent = false } = {}) => {
  if (activeRequest) {
    return activeRequest
  }

  if (!lastUpdatedAt.value) {
    loading.value = true
  }
  refreshing.value = true

  activeRequest = Promise.all([
    api.get('/admin/metrics/overview'),
    api.get('/admin/metrics/routes'),
    api.get('/admin/metrics/models'),
    api.get('/admin/metrics/users'),
    api.get('/admin/metrics/business-codes'),
    api.get('/admin/metrics/model-failures')
  ])
    .then(([overviewResponse, routesResponse, modelsResponse, usersResponse, businessCodesResponse, failuresResponse]) => {
      if ([overviewResponse, routesResponse, modelsResponse, usersResponse, businessCodesResponse, failuresResponse].some(isForbiddenResponse)) {
        const forbiddenError = new Error('FORBIDDEN')
        forbiddenError.code = 'FORBIDDEN'
        throw forbiddenError
      }

      if (!isSuccessResponse(overviewResponse) || !isSuccessResponse(routesResponse) || !isSuccessResponse(modelsResponse) || !isSuccessResponse(usersResponse) || !isSuccessResponse(businessCodesResponse) || !isSuccessResponse(failuresResponse)) {
        throw new Error('指标接口返回失败状态')
      }

      overview.value = {
        ...overview.value,
        ...(overviewResponse.data?.overview || {})
      }
      routes.value = Array.isArray(routesResponse.data?.routes) ? routesResponse.data.routes : []
      models.value = Array.isArray(modelsResponse.data?.models) ? modelsResponse.data.models : []
      users.value = Array.isArray(usersResponse.data?.users) ? usersResponse.data.users : []
      businessCodes.value = Array.isArray(businessCodesResponse.data?.businessCodes) ? businessCodesResponse.data.businessCodes : []
      modelFailures.value = Array.isArray(failuresResponse.data?.failures) ? failuresResponse.data.failures : []
      lastUpdatedAt.value = Date.now()
      appendHistorySample(overview.value)
    })
    .catch((error) => {
      console.error('Fetch metrics error:', error)
      if (error?.code === 'FORBIDDEN' || !isAdminToken(localStorage.getItem('token'))) {
        stopAutoRefresh()
        showToast('当前账号没有管理员权限，已返回控制台', 'error')
        router.replace('/menu')
        return
      }

      if (!silent || !lastUpdatedAt.value) {
        showToast('监控数据刷新失败，请检查后端服务状态', 'error')
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
    message: '退出后将返回登录页，但不会影响后端已采集的监控指标。',
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
    showToast('当前账号没有管理员权限，无法进入管理监控页', 'error')
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
  max-width: 1380px;
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
.panel-title,
.insight-title,
.family-name {
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
.table-loading,
.path-text,
.family-ops,
.family-foot,
.watch-meta,
.insight-desc,
.insight-meta,
.insight-kind {
  color: var(--sci-fi-text-secondary);
}

.brand-subtitle,
.hero-meta,
.family-ops,
.insight-meta,
.insight-kind {
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
  min-width: 138px;
  cursor: pointer;
}

.admin-main {
  max-width: 1380px;
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
  min-width: 280px;
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
.summary-grid,
.trend-grid,
.insight-grid,
.watch-grid {
  display: grid;
  gap: 18px;
}

.overview-grid {
  grid-template-columns: repeat(auto-fit, minmax(180px, 1fr));
}

.summary-grid {
  grid-template-columns: repeat(auto-fit, minmax(260px, 1fr));
}

.trend-grid,
.watch-grid {
  grid-template-columns: repeat(2, minmax(0, 1fr));
}

.insight-grid {
  grid-template-columns: minmax(0, 1.05fr) minmax(0, 0.95fr);
}

.metric-card,
.summary-card {
  min-height: 148px;
}

.metric-label {
  display: block;
  margin-bottom: 12px;
  font-size: 11px;
  letter-spacing: 2px;
  color: var(--sci-fi-text-muted);
}

.metric-value,
.summary-value {
  display: block;
  color: var(--sci-fi-text-primary);
}

.metric-value {
  font-size: 28px;
}

.summary-value {
  font-size: 22px;
  line-height: 1.4;
}

.metric-desc,
.insight-desc {
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

.panel-head-wrap {
  flex-wrap: wrap;
}

.panel-title {
  font-size: 24px;
}

.panel-note {
  margin: 0;
  font-size: 13px;
  line-height: 1.6;
}

.insight-panel {
  min-height: 100%;
}

.insight-list {
  display: grid;
  gap: 14px;
  margin-top: 18px;
}

.insight-item {
  padding: 16px;
  border-radius: 18px;
  background: rgba(255, 255, 255, 0.82);
  border: 1px solid rgba(16, 185, 129, 0.1);
}

.insight-topline,
.family-head,
.watch-head {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  gap: 10px;
}

.insight-kind {
  font-size: 12px;
  line-height: 1.5;
}

.insight-title {
  margin-top: 10px;
  font-size: 19px;
}

.insight-meta {
  margin-top: 10px;
  line-height: 1.6;
}

.family-grid {
  display: grid;
  gap: 14px;
  margin-top: 18px;
}

.family-card,
.watch-item {
  padding: 16px;
  border-radius: 18px;
  background: rgba(255, 255, 255, 0.82);
  border: 1px solid rgba(16, 185, 129, 0.1);
}

.family-name {
  font-size: 18px;
}

.family-metrics,
.watch-meta {
  display: flex;
  flex-wrap: wrap;
  gap: 10px 16px;
  margin-top: 12px;
  font-size: 13px;
  color: var(--sci-fi-text-secondary);
}

.family-track,
.watch-track {
  margin-top: 14px;
  height: 8px;
  border-radius: 999px;
  background: rgba(148, 163, 184, 0.12);
  overflow: hidden;
}

.family-fill,
.watch-fill {
  display: block;
  height: 100%;
  border-radius: inherit;
  background: linear-gradient(90deg, rgba(16, 185, 129, 0.75), rgba(45, 212, 191, 0.88));
}

.watch-fill.warning {
  background: linear-gradient(90deg, rgba(245, 158, 11, 0.78), rgba(251, 191, 36, 0.92));
}

.family-foot {
  margin-top: 12px;
  display: flex;
  flex-wrap: wrap;
  gap: 10px 16px;
  font-size: 12px;
  line-height: 1.6;
}

.watch-list {
  display: grid;
  gap: 14px;
  margin-top: 18px;
}

.trend-card {
  min-height: 264px;
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
  transition: height 0.24s ease;
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

.panel-controls {
  display: flex;
  flex-wrap: wrap;
  gap: 10px;
}

.panel-input,
.panel-select {
  min-height: 42px;
  border-radius: 12px;
  border: 1px solid rgba(16, 185, 129, 0.16);
  background: rgba(255, 255, 255, 0.88);
  color: var(--sci-fi-text-primary);
}

.panel-input {
  min-width: 220px;
  padding: 0 12px;
}

.panel-select {
  min-width: 140px;
  padding: 0 12px;
}

.panel-input:focus,
.panel-select:focus,
.refresh-select:focus {
  outline: none;
  border-color: rgba(16, 185, 129, 0.36);
  box-shadow: 0 0 0 3px rgba(16, 185, 129, 0.12);
}

.panel-empty,
.table-loading {
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
  min-width: 980px;
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

.route-cell {
  display: flex;
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

@media (max-width: 1200px) {
  .insight-grid,
  .trend-grid,
  .watch-grid {
    grid-template-columns: 1fr;
  }
}

@media (max-width: 1100px) {
  .admin-header,
  .hero-panel {
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

  .header-actions,
  .panel-controls {
    width: 100%;
  }

  .header-btn,
  .refresh-toggle,
  .refresh-select,
  .panel-input,
  .panel-select {
    width: 100%;
  }

  .trend-summary,
  .insight-topline,
  .family-head,
  .watch-head,
  .family-foot {
    flex-direction: column;
  }
}
</style>
