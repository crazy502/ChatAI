<template>
  <div class="chat-shell">
    <transition name="overlay-fade">
      <div
        v-if="isMobile && sidebarOpen"
        class="sidebar-overlay"
        @click="sidebarOpen = false"
      ></div>
    </transition>

    <header class="chat-header">
      <div class="header-left">
        <button class="menu-toggle" type="button" aria-label="切换会话列表" @click="toggleSidebar">
          <span class="btn-icon">菜单</span>
        </button>
        <div class="brand-block">
          <div class="logo-badge">AI</div>
          <div>
            <p class="brand-title">GopherAI</p>
            <p class="brand-subtitle">{{ currentSessionLabel }}</p>
          </div>
        </div>
      </div>

      <div class="header-right">
        <label class="control-group">
          <span class="control-label">模型</span>
          <select v-model="selectedModel" class="control-select">
            <option value="qwen">通义千问</option>
            <option value="deepseek">DeepSeek</option>
          </select>
        </label>

        <label class="toggle-group">
          <input v-model="isStreaming" type="checkbox" class="toggle-input" />
          <span class="toggle-slider"></span>
          <span class="toggle-text">流式响应</span>
        </label>

        <button class="ghost-btn" type="button" :disabled="syncing" @click="syncSessions()">
          {{ syncing ? '同步中...' : '同步会话' }}
        </button>
        <button class="danger-btn" type="button" @click="logout">
          退出登录
        </button>
      </div>
    </header>

    <div class="chat-body">
      <aside class="chat-sidebar" :class="{ active: sidebarOpen }">
        <div class="sidebar-header">
          <div>
            <p class="sidebar-kicker">SESSION MATRIX</p>
            <h2 class="sidebar-title">会话列表</h2>
          </div>
          <button class="new-chat-btn" type="button" @click="startNewChat">
            新对话
          </button>
        </div>

        <div class="sidebar-summary">
          <span>共 {{ sessionCount }} 个会话</span>
          <span>{{ tempSession ? '当前为草稿模式' : '当前会话已连接' }}</span>
        </div>

        <div class="sessions-list">
          <div v-if="!orderedSessions.length" class="sessions-empty">
            暂无历史会话，发送第一条消息后会自动创建。
          </div>

          <div
            v-for="session in orderedSessions"
            :key="session.id"
            class="session-item"
            :class="{ active: currentSessionId === session.id }"
          >
            <button class="session-select" type="button" @click="switchSession(session.id)">
              <span class="session-name">{{ session.name }}</span>
              <span class="session-meta">{{ sessionMetaText(session) }}</span>
            </button>
            <button class="delete-btn" type="button" @click.stop="deleteSession(session.id)">
              移除
            </button>
          </div>
        </div>
      </aside>

      <main class="chat-main">
        <div class="status-bar">
          <span class="status-chip">{{ selectedModelLabel }}</span>
          <span class="status-chip muted">{{ isStreaming ? '流式响应已开启' : '标准响应模式' }}</span>
          <span class="status-chip" :class="{ warn: historyLoading || loading }">
            {{ historyLoading ? '载入历史中' : loading ? 'AI 正在生成回复' : '系统就绪' }}
          </span>
        </div>

        <section class="messages-panel" ref="messagesContainer">
          <div v-if="historyLoading && !currentMessages.length" class="messages-skeleton">
            <div v-for="row in 3" :key="row" class="skeleton-row" :class="{ short: row === 2 }"></div>
          </div>

          <div v-else-if="!currentMessages.length" class="empty-state">
            <div class="empty-badge">GopherAI</div>
            <h2 class="empty-title">从一个问题开始新的协作</h2>
            <p class="empty-desc">
              你可以直接提问、让 AI 总结内容，或者让它帮你梳理实现思路。
            </p>
            <div class="prompt-list">
              <button
                v-for="prompt in starterPrompts"
                :key="prompt"
                class="prompt-chip"
                type="button"
                @click="fillSuggestion(prompt)"
              >
                {{ prompt }}
              </button>
            </div>
          </div>

          <transition-group v-else name="message-list" tag="div" class="message-list">
            <article
              v-for="message in currentMessages"
              :key="message.id"
              class="message-item"
              :class="message.role"
            >
              <div class="message-avatar">{{ message.role === 'user' ? '我' : 'AI' }}</div>

              <div class="message-body">
                <div class="message-meta">
                  <span class="sender-name">{{ message.role === 'user' ? '我的提问' : 'AI 助手' }}</span>
                  <span v-if="message.meta?.status === 'streaming'" class="state-tag streaming">生成中</span>
                  <span v-else-if="message.meta?.status === 'error'" class="state-tag error">失败</span>
                </div>

                <div class="message-bubble" :class="bubbleClass(message)">
                  <div
                    v-if="message.content"
                    class="bubble-content"
                    v-html="renderMarkdown(message.content)"
                  ></div>
                  <div v-else class="stream-placeholder">
                    <span class="placeholder-dot"></span>
                    <span class="placeholder-dot"></span>
                    <span class="placeholder-dot"></span>
                  </div>
                </div>

                <div v-if="message.role === 'assistant'" class="message-actions">
                  <button
                    class="action-btn"
                    type="button"
                    :disabled="!hasMessageContent(message.content)"
                    @click="copyMessage(message.content)"
                  >
                    复制
                  </button>
                  <button
                    class="action-btn"
                    type="button"
                    :disabled="!hasMessageContent(message.content)"
                    @click="playTTS(message.content)"
                  >
                    朗读
                  </button>
                </div>
              </div>
            </article>
          </transition-group>
        </section>

        <footer class="input-area">
          <div class="input-shell" :class="{ focused: inputFocused, busy: loading }">
            <textarea
              ref="messageInput"
              v-model="inputMessage"
              class="message-input"
              rows="1"
              :disabled="loading"
              placeholder="输入你的问题，回车发送，Shift + Enter 换行"
              @input="handleInput"
              @focus="inputFocused = true"
              @blur="inputFocused = false"
              @keydown.enter.exact.prevent="sendMessage"
            ></textarea>
            <button
              class="send-btn"
              type="button"
              :disabled="!inputMessage.trim() || loading"
              @click="sendMessage"
            >
              {{ loading ? '发送中' : '发送' }}
            </button>
          </div>
          <div class="input-footer">
            <span>{{ tempSession ? '发送后会自动创建新会话' : '消息会同步到当前会话' }}</span>
            <span>Enter 发送，Shift + Enter 换行</span>
          </div>
        </footer>
      </main>
    </div>
  </div>
</template>

<script setup>
import { computed, nextTick, onBeforeUnmount, onMounted, ref, watch } from 'vue'
import { useRouter } from 'vue-router'
import { marked } from 'marked'
import hljs from 'highlight.js/lib/core'
import bash from 'highlight.js/lib/languages/bash'
import css from 'highlight.js/lib/languages/css'
import go from 'highlight.js/lib/languages/go'
import javascript from 'highlight.js/lib/languages/javascript'
import json from 'highlight.js/lib/languages/json'
import plaintext from 'highlight.js/lib/languages/plaintext'
import python from 'highlight.js/lib/languages/python'
import xml from 'highlight.js/lib/languages/xml'
import 'highlight.js/styles/github.css'
import api, { buildApiUrl } from '../utils/api.js'
import { useUi } from '../composables/useUi'

hljs.registerLanguage('bash', bash)
hljs.registerLanguage('sh', bash)
hljs.registerLanguage('css', css)
hljs.registerLanguage('go', go)
hljs.registerLanguage('javascript', javascript)
hljs.registerLanguage('js', javascript)
hljs.registerLanguage('json', json)
hljs.registerLanguage('plaintext', plaintext)
hljs.registerLanguage('text', plaintext)
hljs.registerLanguage('python', python)
hljs.registerLanguage('py', python)
hljs.registerLanguage('html', xml)
hljs.registerLanguage('xml', xml)

marked.setOptions({
  gfm: true,
  breaks: true,
  highlight(code, lang) {
    const language = (lang || '').trim().toLowerCase()
    if (language && hljs.getLanguage(language)) {
      return hljs.highlight(code, { language }).value
    }
    return hljs.highlightAuto(code).value
  }
})

const starterPrompts = [
  '帮我总结这段需求并拆出实现步骤',
  '给我一个可执行的前端优化清单',
  '帮我分析这个接口调用为什么失败',
  '把这段文档整理成结构化要点'
]

const router = useRouter()
const { showToast, confirmAction } = useUi()

const inputMessage = ref('')
const sessions = ref({})
const currentSessionId = ref('')
const currentMessages = ref([])
const loading = ref(false)
const syncing = ref(false)
const historyLoading = ref(false)
const sidebarOpen = ref(false)
const inputFocused = ref(false)
const isStreaming = ref(true)
const selectedModel = ref('qwen')
const tempSession = ref(false)
const messageInput = ref(null)
const messagesContainer = ref(null)
const viewportWidth = ref(typeof window === 'undefined' ? 1280 : window.innerWidth)

let messageSeed = 0

const isMobile = computed(() => viewportWidth.value <= 768)
const sessionCount = computed(() => Object.keys(sessions.value).length)
const orderedSessions = computed(() => {
  return Object.values(sessions.value).sort((left, right) => (right.updatedAt || 0) - (left.updatedAt || 0))
})
const activeSession = computed(() => sessions.value[currentSessionId.value] || null)
const currentSessionLabel = computed(() => {
  if (tempSession.value) {
    return '新对话草稿'
  }
  return activeSession.value?.name || '会话未选择'
})
const selectedModelLabel = computed(() => (selectedModel.value === 'deepseek' ? 'DEEPSEEK' : 'QWEN'))

const createMessage = (role, content = '', meta = {}) => ({
  id: `msg-${Date.now()}-${++messageSeed}`,
  role,
  content,
  meta: Object.keys(meta).length ? { ...meta } : undefined
})

const cloneMessages = (messages = []) => messages.map((message) => ({
  ...message,
  meta: message.meta ? { ...message.meta } : undefined
}))

const normalizeSessionName = (name, fallback = '未命名会话') => {
  const normalized = (name || '').trim()
  if (!normalized) {
    return fallback
  }
  return normalized.length > 26 ? `${normalized.slice(0, 26)}...` : normalized
}

const hasMessageContent = (content) => Boolean(content && content.trim())

const bubbleClass = (message) => ({
  streaming: message.meta?.status === 'streaming',
  error: message.meta?.status === 'error'
})

const renderMarkdown = (content) => marked.parse(content || '')

const updateViewportWidth = () => {
  viewportWidth.value = window.innerWidth
  if (!isMobile.value) {
    sidebarOpen.value = false
  }
}

const scrollToBottom = () => {
  if (messagesContainer.value) {
    messagesContainer.value.scrollTop = messagesContainer.value.scrollHeight
  }
}

const resizeTextarea = () => {
  const textarea = messageInput.value
  if (!textarea) {
    return
  }

  textarea.style.height = 'auto'
  textarea.style.height = `${Math.min(textarea.scrollHeight, 220)}px`
}

const focusInput = () => {
  nextTick(() => {
    messageInput.value?.focus()
    resizeTextarea()
  })
}

const sessionMetaText = (session) => {
  if (session.historyLoaded) {
    return session.messages.length ? `${session.messages.length} 条消息` : '暂无历史消息'
  }
  if (session.historyUnavailable) {
    return '历史暂不可用'
  }
  return '等待载入历史'
}

const persistCurrentSessionMessages = () => {
  if (tempSession.value || !currentSessionId.value || !sessions.value[currentSessionId.value]) {
    return
  }

  const sessionId = currentSessionId.value
  sessions.value = {
    ...sessions.value,
    [sessionId]: {
      ...sessions.value[sessionId],
      messages: cloneMessages(currentMessages.value),
      historyLoaded: true,
      historyUnavailable: false,
      updatedAt: Date.now()
    }
  }
}

const updateMessage = (messageId, updater, persist = true) => {
  const nextMessages = cloneMessages(currentMessages.value)
  const target = nextMessages.find((message) => message.id === messageId)
  if (!target) {
    return
  }

  updater(target)
  currentMessages.value = nextMessages

  if (persist) {
    persistCurrentSessionMessages()
  }
}
const registerNewSession = (sessionId, seedQuestion) => {
  const normalizedId = String(sessionId)
  sessions.value = {
    ...sessions.value,
    [normalizedId]: {
      id: normalizedId,
      name: normalizeSessionName(seedQuestion, '新会话'),
      messages: cloneMessages(currentMessages.value),
      historyLoaded: true,
      historyUnavailable: false,
      updatedAt: Date.now()
    }
  }
  currentSessionId.value = normalizedId
  tempSession.value = false
}

const loadSessionHistory = async (sessionId, { silent = false } = {}) => {
  if (!sessions.value[sessionId]) {
    return false
  }

  historyLoading.value = true
  try {
    const response = await api.post('/AI/chat/history', { sessionId })
    if (response.data?.status_code === 1000) {
      const history = Array.isArray(response.data.history) ? response.data.history : []
      const mappedMessages = history.map((item) => createMessage(item.is_user ? 'user' : 'assistant', item.content || ''))

      sessions.value = {
        ...sessions.value,
        [sessionId]: {
          ...sessions.value[sessionId],
          messages: mappedMessages,
          historyLoaded: true,
          historyUnavailable: mappedMessages.length === 0,
          updatedAt: sessions.value[sessionId].updatedAt || Date.now()
        }
      }

      if (currentSessionId.value === sessionId) {
        currentMessages.value = cloneMessages(mappedMessages)
      }
      return true
    }

    sessions.value = {
      ...sessions.value,
      [sessionId]: {
        ...sessions.value[sessionId],
        historyLoaded: true,
        historyUnavailable: true
      }
    }

    if (!silent) {
      showToast(response.data?.status_msg || '当前会话暂无可加载历史', 'warning')
    }
  } catch (error) {
    console.error('Load history error:', error)
    sessions.value = {
      ...sessions.value,
      [sessionId]: {
        ...sessions.value[sessionId],
        historyUnavailable: true
      }
    }
    if (!silent) {
      showToast('加载会话历史失败', 'error')
    }
  } finally {
    historyLoading.value = false
    nextTick(() => {
      scrollToBottom()
    })
  }

  return false
}

const startNewChat = () => {
  tempSession.value = true
  currentSessionId.value = ''
  currentMessages.value = []
  historyLoading.value = false
  if (isMobile.value) {
    sidebarOpen.value = false
  }
  focusInput()
}

const switchSession = async (sessionId, options = {}) => {
  if (!sessions.value[sessionId]) {
    return
  }

  currentSessionId.value = sessionId
  tempSession.value = false
  currentMessages.value = cloneMessages(sessions.value[sessionId].messages || [])

  if (isMobile.value) {
    sidebarOpen.value = false
  }

  if (!sessions.value[sessionId].historyLoaded || options.forceRefresh) {
    await loadSessionHistory(sessionId, { silent: options.silent })
  } else {
    nextTick(() => {
      scrollToBottom()
    })
  }
}

const syncSessions = async ({ silent = false } = {}) => {
  syncing.value = true
  try {
    const response = await api.get('/AI/chat/sessions')
    if (response.data?.status_code === 1000) {
      const incomingSessions = Array.isArray(response.data.sessions) ? response.data.sessions : []
      const nextSessions = {}

      incomingSessions.forEach((session, index) => {
        const sessionId = String(session.SessionID || session.id || '')
        if (!sessionId) {
          return
        }

        const existing = sessions.value[sessionId]
        nextSessions[sessionId] = {
          id: sessionId,
          name: normalizeSessionName(session.Title || existing?.name || '未命名会话'),
          messages: existing?.messages ? cloneMessages(existing.messages) : [],
          historyLoaded: existing?.historyLoaded || false,
          historyUnavailable: existing?.historyUnavailable || false,
          updatedAt: existing?.updatedAt || Date.now() - index
        }
      })

      sessions.value = nextSessions

      if (!currentSessionId.value) {
        const firstSession = Object.keys(nextSessions)[0]
        if (firstSession) {
          await switchSession(firstSession, { silent: true })
        } else {
          startNewChat()
        }
      } else if (!nextSessions[currentSessionId.value]) {
        startNewChat()
      }

      if (!silent) {
        showToast('会话已同步', 'success', { duration: 1800 })
      }
    } else if (!silent) {
      showToast(response.data?.status_msg || '同步会话失败', 'error')
    }
  } catch (error) {
    console.error('Sync sessions error:', error)
    if (!silent) {
      showToast('同步会话失败', 'error')
    }
  } finally {
    syncing.value = false
  }
}

const toggleSidebar = () => {
  if (isMobile.value) {
    sidebarOpen.value = !sidebarOpen.value
  }
}

const deleteSession = async (sessionId) => {
  const session = sessions.value[sessionId]
  if (!session) {
    return
  }

  const confirmed = await confirmAction({
    title: '移除当前会话？',
    message: '这只会从当前前端列表移除，不会调用后端删除接口；下次同步时它仍可能重新出现。',
    confirmText: '确认移除',
    cancelText: '取消',
    intent: 'danger'
  })

  if (!confirmed) {
    return
  }

  const nextSessions = { ...sessions.value }
  delete nextSessions[sessionId]
  sessions.value = nextSessions

  if (currentSessionId.value === sessionId) {
    const nextSessionId = Object.keys(nextSessions)[0]
    if (nextSessionId) {
      await switchSession(nextSessionId, { silent: true })
    } else {
      startNewChat()
    }
  }

  showToast('已从当前列表移除', 'success', { duration: 1800 })
}

const buildDraftMessages = (question) => {
  const userMessage = createMessage('user', question)
  const aiMessage = createMessage('assistant', '', {
    status: isStreaming.value ? 'streaming' : 'pending'
  })

  currentMessages.value = [...currentMessages.value, userMessage, aiMessage]
  if (!tempSession.value) {
    persistCurrentSessionMessages()
  }

  return aiMessage.id
}

const handleStreaming = async (question, aiMessageId) => {
  const url = tempSession.value
    ? buildApiUrl('/AI/chat/send-stream-new-session')
    : buildApiUrl('/AI/chat/send-stream')

  const headers = {
    'Content-Type': 'application/json',
    Authorization: `Bearer ${localStorage.getItem('token') || ''}`
  }

  const payload = tempSession.value
    ? { question, modelType: selectedModel.value }
    : { question, modelType: selectedModel.value, sessionId: currentSessionId.value }

  try {
    const response = await fetch(url, {
      method: 'POST',
      headers,
      body: JSON.stringify(payload)
    })

    if (!response.ok || !response.body) {
      throw new Error('流式请求失败')
    }

    const reader = response.body.getReader()
    const decoder = new TextDecoder()
    let buffer = ''

    for (;;) {
      const { done, value } = await reader.read()
      if (done) {
        break
      }

      buffer += decoder.decode(value, { stream: true })
      const lines = buffer.split('\n')
      buffer = lines.pop() || ''

      for (const rawLine of lines) {
        const line = rawLine.trim()
        if (!line || line.startsWith('event:')) {
          continue
        }

        if (!line.startsWith('data:')) {
          continue
        }

        const data = line.slice(5).trim()
        if (!data) {
          continue
        }

        if (data === '[DONE]') {
          updateMessage(aiMessageId, (message) => {
            message.meta = { ...(message.meta || {}), status: 'done' }
          })
          continue
        }

        if (data.startsWith('{')) {
          try {
            const parsed = JSON.parse(data)
            if (parsed.sessionId) {
              registerNewSession(parsed.sessionId, question)
              continue
            }
            if (typeof parsed.content === 'string') {
              updateMessage(aiMessageId, (message) => {
                message.content += parsed.content
                message.meta = { ...(message.meta || {}), status: 'streaming' }
              })
              continue
            }
            if (typeof parsed.message === 'string') {
              throw new Error(parsed.message)
            }
          } catch (error) {
            if (error instanceof SyntaxError) {
              updateMessage(aiMessageId, (message) => {
                message.content += data
                message.meta = { ...(message.meta || {}), status: 'streaming' }
              })
            } else {
              throw error
            }
          }
        } else {
          updateMessage(aiMessageId, (message) => {
            message.content += data
            message.meta = { ...(message.meta || {}), status: 'streaming' }
          })
        }

        await nextTick()
        requestAnimationFrame(() => {
          scrollToBottom()
        })
      }
    }
  } catch (error) {
    console.error('Stream error:', error)
    updateMessage(aiMessageId, (message) => {
      if (!message.content) {
        message.content = '本次回复未完成，请重试。'
      }
      message.meta = { ...(message.meta || {}), status: 'error' }
    })
    showToast(error.message || '流式传输出错', 'error')
  }
}

const handleNormal = async (question, aiMessageId) => {
  try {
    if (tempSession.value) {
      const response = await api.post('/AI/chat/send-new-session', {
        question,
        modelType: selectedModel.value
      })

      if (response.data?.status_code !== 1000) {
        throw new Error(response.data?.status_msg || '发送失败')
      }

      registerNewSession(response.data.sessionId, question)
      updateMessage(aiMessageId, (message) => {
        message.content = response.data.Information || ''
        message.meta = { ...(message.meta || {}), status: 'done' }
      })
      persistCurrentSessionMessages()
      return
    }

    const response = await api.post('/AI/chat/send', {
      question,
      modelType: selectedModel.value,
      sessionId: currentSessionId.value
    })

    if (response.data?.status_code !== 1000) {
      throw new Error(response.data?.status_msg || '发送失败')
    }

    updateMessage(aiMessageId, (message) => {
      message.content = response.data.Information || ''
      message.meta = { ...(message.meta || {}), status: 'done' }
    })
  } catch (error) {
    console.error('Send message error:', error)
    updateMessage(aiMessageId, (message) => {
      message.content = message.content || '本次回复未完成，请重试。'
      message.meta = { ...(message.meta || {}), status: 'error' }
    })
    showToast(error.message || '发送失败', 'error')
  }
}

const sendMessage = async () => {
  const question = inputMessage.value.trim()
  if (!question || loading.value) {
    return
  }

  inputMessage.value = ''
  resizeTextarea()
  loading.value = true

  const aiMessageId = buildDraftMessages(question)
  await nextTick()
  scrollToBottom()

  if (isStreaming.value) {
    await handleStreaming(question, aiMessageId)
  } else {
    await handleNormal(question, aiMessageId)
  }

  loading.value = false
  persistCurrentSessionMessages()
  focusInput()
}
const copyMessage = async (content) => {
  if (!hasMessageContent(content)) {
    return
  }

  if (!navigator.clipboard?.writeText) {
    showToast('当前浏览器不支持复制到剪贴板', 'warning')
    return
  }

  try {
    await navigator.clipboard.writeText(content)
    showToast('内容已复制', 'success', { duration: 1600 })
  } catch (error) {
    console.error('Copy error:', error)
    showToast('复制失败，请检查浏览器权限', 'error')
  }
}

const playTTS = (content) => {
  if (!hasMessageContent(content)) {
    return
  }

  if (!('speechSynthesis' in window)) {
    showToast('当前浏览器不支持语音播放', 'warning')
    return
  }

  window.speechSynthesis.cancel()
  const utterance = new SpeechSynthesisUtterance(content)
  utterance.lang = 'zh-CN'
  window.speechSynthesis.speak(utterance)
  showToast('开始语音播报', 'info', { duration: 1400 })
}

const fillSuggestion = (prompt) => {
  inputMessage.value = prompt
  focusInput()
}

const handleInput = () => {
  resizeTextarea()
}

const logout = async () => {
  const confirmed = await confirmAction({
    title: '退出当前账号？',
    message: '退出后将返回登录页，但不会影响后端已有会话数据。',
    confirmText: '退出登录',
    cancelText: '继续使用',
    intent: 'danger'
  })

  if (!confirmed) {
    return
  }

  localStorage.removeItem('token')
  router.push('/login')
}

watch(inputMessage, () => {
  nextTick(() => {
    resizeTextarea()
  })
})

watch(currentMessages, () => {
  nextTick(() => {
    scrollToBottom()
  })
}, { deep: true })

onMounted(async () => {
  window.addEventListener('resize', updateViewportWidth)
  await syncSessions({ silent: true })
  if (!Object.keys(sessions.value).length) {
    startNewChat()
  } else {
    focusInput()
  }
})

onBeforeUnmount(() => {
  window.removeEventListener('resize', updateViewportWidth)
  if ('speechSynthesis' in window) {
    window.speechSynthesis.cancel()
  }
})
</script>

<style scoped>
.chat-shell {
  min-height: 100vh;
  display: flex;
  flex-direction: column;
  background:
    radial-gradient(circle at top left, rgba(16, 185, 129, 0.12), transparent 30%),
    radial-gradient(circle at bottom right, rgba(45, 212, 191, 0.14), transparent 35%),
    linear-gradient(135deg, #f4fbf6 0%, #ecf7f0 45%, #f9fcfa 100%);
  position: relative;
  overflow: hidden;
}

.chat-shell::before {
  content: '';
  position: absolute;
  inset: 0;
  background-image:
    linear-gradient(rgba(16, 185, 129, 0.04) 1px, transparent 1px),
    linear-gradient(90deg, rgba(16, 185, 129, 0.04) 1px, transparent 1px);
  background-size: 48px 48px;
  opacity: 0.55;
  pointer-events: none;
}

.chat-header,
.chat-body {
  position: relative;
  z-index: 1;
}

.chat-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: 24px;
  padding: 18px 24px;
  background: rgba(255, 255, 255, 0.82);
  border-bottom: 1px solid rgba(16, 185, 129, 0.14);
  backdrop-filter: blur(18px);
}

.header-left,
.header-right {
  display: flex;
  align-items: center;
  gap: 18px;
}

.menu-toggle {
  display: none;
  min-width: 72px;
  padding: 10px 14px;
  border-radius: 12px;
  border: 1px solid rgba(16, 185, 129, 0.18);
  background: rgba(255, 255, 255, 0.9);
  color: var(--sci-fi-text-primary);
  cursor: pointer;
}

.brand-block {
  display: flex;
  align-items: center;
  gap: 14px;
}

.logo-badge {
  width: 44px;
  height: 44px;
  display: grid;
  place-items: center;
  border-radius: 14px;
  background: linear-gradient(135deg, var(--sci-fi-primary), var(--sci-fi-secondary));
  color: #fff;
  font-family: 'Orbitron', sans-serif;
  font-size: 15px;
  letter-spacing: 1px;
  box-shadow: 0 14px 24px rgba(16, 185, 129, 0.18);
}

.brand-title {
  margin: 0;
  font-family: 'Orbitron', sans-serif;
  font-size: 20px;
  letter-spacing: 2px;
  color: var(--sci-fi-text-primary);
}

.brand-subtitle {
  margin: 4px 0 0;
  font-size: 13px;
  color: var(--sci-fi-text-secondary);
}

.control-group {
  display: flex;
  align-items: center;
  gap: 10px;
}

.control-label,
.toggle-text,
.sidebar-kicker,
.sender-name,
.status-chip,
.new-chat-btn,
.ghost-btn,
.danger-btn,
.action-btn,
.send-btn {
  font-family: 'Orbitron', sans-serif;
}

.control-label {
  font-size: 11px;
  letter-spacing: 1px;
  color: var(--sci-fi-text-muted);
}

.control-select {
  min-width: 124px;
  padding: 10px 12px;
  border-radius: 12px;
  border: 1px solid rgba(16, 185, 129, 0.2);
  background: rgba(255, 255, 255, 0.88);
  color: var(--sci-fi-text-primary);
}

.toggle-group {
  display: flex;
  align-items: center;
  gap: 10px;
  cursor: pointer;
}

.toggle-input {
  display: none;
}

.toggle-slider {
  width: 44px;
  height: 24px;
  border-radius: 999px;
  background: rgba(148, 163, 184, 0.4);
  position: relative;
  transition: background 0.2s ease;
}

.toggle-slider::before {
  content: '';
  position: absolute;
  top: 2px;
  left: 2px;
  width: 20px;
  height: 20px;
  border-radius: 50%;
  background: #fff;
  box-shadow: 0 4px 10px rgba(15, 23, 42, 0.16);
  transition: transform 0.2s ease;
}

.toggle-input:checked + .toggle-slider {
  background: rgba(16, 185, 129, 0.88);
}

.toggle-input:checked + .toggle-slider::before {
  transform: translateX(20px);
}

.toggle-text {
  font-size: 11px;
  letter-spacing: 1px;
  color: var(--sci-fi-text-secondary);
}

.ghost-btn,
.danger-btn,
.new-chat-btn,
.action-btn,
.send-btn {
  border: none;
  cursor: pointer;
  transition: transform 0.2s ease, box-shadow 0.2s ease, background 0.2s ease;
}

.ghost-btn,
.danger-btn {
  padding: 10px 16px;
  border-radius: 12px;
  font-size: 12px;
  letter-spacing: 1px;
}

.ghost-btn {
  background: rgba(16, 185, 129, 0.12);
  color: var(--sci-fi-primary);
}

.danger-btn {
  background: rgba(239, 68, 68, 0.12);
  color: var(--sci-fi-danger);
}

.ghost-btn:hover,
.danger-btn:hover,
.new-chat-btn:hover,
.action-btn:hover,
.send-btn:hover:not(:disabled) {
  transform: translateY(-1px);
}

.ghost-btn:disabled,
.send-btn:disabled {
  cursor: not-allowed;
  opacity: 0.6;
}

.chat-body {
  flex: 1;
  display: flex;
  min-height: 0;
}

.chat-sidebar {
  width: 300px;
  display: flex;
  flex-direction: column;
  border-right: 1px solid rgba(16, 185, 129, 0.12);
  background: rgba(255, 255, 255, 0.72);
  backdrop-filter: blur(18px);
}

.sidebar-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: 12px;
  padding: 22px 20px 14px;
}

.sidebar-kicker {
  margin: 0 0 6px;
  font-size: 10px;
  letter-spacing: 2px;
  color: var(--sci-fi-primary);
}

.sidebar-title {
  margin: 0;
  font-size: 24px;
  color: var(--sci-fi-text-primary);
}

.new-chat-btn {
  padding: 10px 14px;
  border-radius: 12px;
  background: linear-gradient(135deg, var(--sci-fi-primary), var(--sci-fi-secondary));
  color: #fff;
  font-size: 11px;
  letter-spacing: 1px;
  box-shadow: 0 14px 24px rgba(16, 185, 129, 0.18);
}

.sidebar-summary {
  display: flex;
  justify-content: space-between;
  gap: 12px;
  padding: 0 20px 14px;
  font-size: 12px;
  color: var(--sci-fi-text-muted);
}

.sessions-list {
  flex: 1;
  overflow-y: auto;
  padding: 0 14px 18px;
}
.sessions-empty {
  padding: 18px;
  border-radius: 16px;
  background: rgba(255, 255, 255, 0.78);
  color: var(--sci-fi-text-secondary);
  font-size: 14px;
  line-height: 1.6;
}

.session-item {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-bottom: 10px;
  padding: 8px;
  border-radius: 16px;
  background: rgba(255, 255, 255, 0.7);
  border: 1px solid rgba(16, 185, 129, 0.08);
}

.session-item.active {
  border-color: rgba(16, 185, 129, 0.28);
  box-shadow: 0 14px 24px rgba(16, 185, 129, 0.12);
}

.session-select {
  flex: 1;
  min-width: 0;
  padding: 10px 12px;
  border: none;
  border-radius: 12px;
  background: transparent;
  text-align: left;
  cursor: pointer;
}

.session-name,
.session-meta {
  display: block;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.session-name {
  font-size: 14px;
  font-weight: 600;
  color: var(--sci-fi-text-primary);
}

.session-meta {
  margin-top: 6px;
  font-size: 12px;
  color: var(--sci-fi-text-muted);
}

.delete-btn {
  flex-shrink: 0;
  padding: 8px 10px;
  border: none;
  border-radius: 10px;
  background: rgba(239, 68, 68, 0.08);
  color: var(--sci-fi-danger);
  cursor: pointer;
}

.chat-main {
  flex: 1;
  min-width: 0;
  display: flex;
  flex-direction: column;
}

.status-bar {
  display: flex;
  flex-wrap: wrap;
  gap: 10px;
  padding: 18px 24px 0;
}

.status-chip {
  padding: 8px 12px;
  border-radius: 999px;
  background: rgba(16, 185, 129, 0.12);
  color: var(--sci-fi-primary);
  font-size: 11px;
  letter-spacing: 1px;
}

.status-chip.muted {
  color: var(--sci-fi-text-secondary);
  background: rgba(148, 163, 184, 0.12);
}

.status-chip.warn {
  color: #b45309;
  background: rgba(245, 158, 11, 0.16);
}

.messages-panel {
  flex: 1;
  overflow-y: auto;
  padding: 20px 24px 24px;
  scroll-behavior: smooth;
}

.messages-skeleton,
.empty-state {
  min-height: 100%;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
}

.messages-skeleton {
  gap: 18px;
}

.skeleton-row {
  width: min(100%, 720px);
  height: 22px;
  border-radius: 999px;
  background: linear-gradient(90deg, rgba(16, 185, 129, 0.08), rgba(16, 185, 129, 0.2), rgba(16, 185, 129, 0.08));
  background-size: 200% 100%;
  animation: shimmer 1.3s linear infinite;
}

.skeleton-row.short {
  width: min(76%, 560px);
}

@keyframes shimmer {
  from {
    background-position: 200% 0;
  }

  to {
    background-position: -200% 0;
  }
}

.empty-badge {
  padding: 10px 18px;
  border-radius: 999px;
  background: rgba(16, 185, 129, 0.12);
  color: var(--sci-fi-primary);
  font-family: 'Orbitron', sans-serif;
  font-size: 12px;
  letter-spacing: 2px;
}

.empty-title {
  margin: 22px 0 10px;
  font-family: 'Orbitron', sans-serif;
  font-size: 28px;
  letter-spacing: 2px;
  text-align: center;
  color: var(--sci-fi-text-primary);
}

.empty-desc {
  max-width: 520px;
  margin: 0;
  text-align: center;
  font-size: 15px;
  line-height: 1.8;
  color: var(--sci-fi-text-secondary);
}

.prompt-list {
  display: flex;
  flex-wrap: wrap;
  justify-content: center;
  gap: 10px;
  margin-top: 26px;
}

.prompt-chip {
  padding: 12px 16px;
  border-radius: 14px;
  border: 1px solid rgba(16, 185, 129, 0.16);
  background: rgba(255, 255, 255, 0.8);
  color: var(--sci-fi-text-primary);
  cursor: pointer;
}

.message-list {
  display: flex;
  flex-direction: column;
  gap: 18px;
}

.message-item {
  max-width: min(78%, 820px);
  display: flex;
  gap: 12px;
}

.message-item.user {
  align-self: flex-end;
  flex-direction: row-reverse;
}

.message-avatar {
  width: 42px;
  height: 42px;
  flex-shrink: 0;
  display: grid;
  place-items: center;
  border-radius: 14px;
  background: rgba(16, 185, 129, 0.14);
  color: var(--sci-fi-primary);
  font-family: 'Orbitron', sans-serif;
  font-size: 13px;
  letter-spacing: 1px;
}

.message-item.user .message-avatar {
  background: rgba(45, 212, 191, 0.18);
  color: #0f766e;
}

.message-body {
  min-width: 0;
}

.message-meta {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-bottom: 8px;
}

.sender-name {
  font-size: 10px;
  letter-spacing: 1.5px;
  color: var(--sci-fi-text-muted);
}

.state-tag {
  padding: 4px 8px;
  border-radius: 999px;
  font-size: 11px;
  color: #fff;
}

.state-tag.streaming {
  background: rgba(16, 185, 129, 0.78);
}

.state-tag.error {
  background: rgba(239, 68, 68, 0.84);
}

.message-bubble {
  padding: 16px 18px;
  border-radius: 18px;
  border: 1px solid rgba(16, 185, 129, 0.16);
  background: rgba(255, 255, 255, 0.9);
  box-shadow: 0 18px 32px rgba(15, 23, 42, 0.08);
}

.message-item.user .message-bubble {
  background: rgba(16, 185, 129, 0.12);
  border-color: rgba(16, 185, 129, 0.24);
}

.message-bubble.streaming {
  box-shadow: 0 18px 32px rgba(16, 185, 129, 0.12);
}

.message-bubble.error {
  border-color: rgba(239, 68, 68, 0.24);
}

.bubble-content {
  font-size: 15px;
  line-height: 1.75;
  color: var(--sci-fi-text-primary);
  word-break: break-word;
}

.bubble-content :deep(p) {
  margin: 0 0 10px;
}

.bubble-content :deep(p:last-child) {
  margin-bottom: 0;
}

.bubble-content :deep(code) {
  padding: 2px 6px;
  border-radius: 6px;
  background: rgba(16, 185, 129, 0.1);
  color: var(--sci-fi-primary);
  font-size: 0.92em;
}
.bubble-content :deep(pre) {
  margin: 12px 0;
  padding: 14px;
  border-radius: 12px;
  overflow-x: auto;
  background: rgba(15, 23, 42, 0.04);
  border: 1px solid rgba(16, 185, 129, 0.14);
}

.bubble-content :deep(pre code) {
  background: transparent;
  padding: 0;
  color: var(--sci-fi-text-primary);
}

.bubble-content :deep(blockquote) {
  margin: 12px 0;
  padding-left: 14px;
  border-left: 3px solid rgba(16, 185, 129, 0.4);
  color: var(--sci-fi-text-secondary);
}

.stream-placeholder {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  min-height: 24px;
}

.placeholder-dot {
  width: 8px;
  height: 8px;
  border-radius: 50%;
  background: rgba(16, 185, 129, 0.7);
  animation: dot-bounce 0.9s ease-in-out infinite;
}

.placeholder-dot:nth-child(2) {
  animation-delay: 0.15s;
}

.placeholder-dot:nth-child(3) {
  animation-delay: 0.3s;
}

@keyframes dot-bounce {
  0%,
  80%,
  100% {
    transform: translateY(0);
    opacity: 0.45;
  }

  40% {
    transform: translateY(-4px);
    opacity: 1;
  }
}

.message-actions {
  display: flex;
  gap: 8px;
  margin-top: 10px;
}

.action-btn {
  padding: 8px 12px;
  border-radius: 10px;
  background: rgba(16, 185, 129, 0.1);
  color: var(--sci-fi-primary);
  font-size: 11px;
  letter-spacing: 1px;
}

.action-btn:disabled {
  cursor: not-allowed;
  opacity: 0.45;
}

.input-area {
  padding: 18px 24px 24px;
  background: rgba(255, 255, 255, 0.82);
  border-top: 1px solid rgba(16, 185, 129, 0.12);
  backdrop-filter: blur(18px);
}

.input-shell {
  display: flex;
  align-items: flex-end;
  gap: 12px;
  padding: 12px;
  border-radius: 22px;
  background: rgba(255, 255, 255, 0.94);
  border: 1px solid rgba(16, 185, 129, 0.16);
  box-shadow: 0 16px 32px rgba(15, 23, 42, 0.08);
}

.input-shell.focused {
  border-color: rgba(16, 185, 129, 0.32);
  box-shadow: 0 18px 36px rgba(16, 185, 129, 0.12);
}

.input-shell.busy {
  opacity: 0.94;
}

.message-input {
  flex: 1;
  min-height: 52px;
  max-height: 220px;
  padding: 10px 6px 10px 8px;
  border: none;
  background: transparent;
  resize: none;
  outline: none;
  font-size: 15px;
  line-height: 1.65;
  color: var(--sci-fi-text-primary);
}

.send-btn {
  min-width: 104px;
  padding: 14px 18px;
  border-radius: 16px;
  background: linear-gradient(135deg, var(--sci-fi-primary), var(--sci-fi-secondary));
  color: #fff;
  font-size: 12px;
  letter-spacing: 1px;
  box-shadow: 0 18px 28px rgba(16, 185, 129, 0.16);
}

.input-footer {
  margin-top: 12px;
  display: flex;
  justify-content: space-between;
  gap: 12px;
  font-size: 12px;
  color: var(--sci-fi-text-muted);
}

.sidebar-overlay {
  position: fixed;
  inset: 0;
  z-index: 20;
  background: rgba(15, 23, 42, 0.28);
}

.overlay-fade-enter-active,
.overlay-fade-leave-active,
.message-list-enter-active,
.message-list-leave-active {
  transition: opacity 0.22s ease, transform 0.22s ease;
}

.overlay-fade-enter-from,
.overlay-fade-leave-to,
.message-list-enter-from,
.message-list-leave-to {
  opacity: 0;
}

.message-list-enter-from,
.message-list-leave-to {
  transform: translateY(10px);
}

@media (max-width: 1024px) {
  .chat-header {
    align-items: flex-start;
    flex-direction: column;
  }

  .header-right {
    flex-wrap: wrap;
  }

  .message-item {
    max-width: 88%;
  }
}

@media (max-width: 768px) {
  .menu-toggle {
    display: inline-flex;
    align-items: center;
    justify-content: center;
  }

  .chat-body {
    position: relative;
  }

  .chat-sidebar {
    position: absolute;
    top: 0;
    left: 0;
    bottom: 0;
    z-index: 30;
    width: min(86vw, 320px);
    transform: translateX(-100%);
    transition: transform 0.24s ease;
    box-shadow: 24px 0 48px rgba(15, 23, 42, 0.12);
  }

  .chat-sidebar.active {
    transform: translateX(0);
  }

  .message-item {
    max-width: 100%;
  }

  .input-footer {
    flex-direction: column;
  }
}

@media (max-width: 640px) {
  .chat-header,
  .status-bar,
  .messages-panel,
  .input-area {
    padding-left: 14px;
    padding-right: 14px;
  }

  .header-right {
    width: 100%;
    gap: 12px;
  }

  .control-group {
    width: 100%;
    justify-content: space-between;
  }

  .toggle-group,
  .ghost-btn,
  .danger-btn {
    width: 100%;
    justify-content: center;
  }

  .ghost-btn,
  .danger-btn {
    text-align: center;
  }

  .message-item.user,
  .message-item {
    flex-direction: column;
    align-self: stretch;
  }

  .message-avatar {
    width: 34px;
    height: 34px;
    border-radius: 12px;
  }

  .message-actions {
    flex-wrap: wrap;
  }

  .delete-btn {
    padding: 8px 10px;
  }
}

@media (prefers-reduced-motion: reduce) {
  *,
  *::before,
  *::after {
    animation: none !important;
    transition: none !important;
    scroll-behavior: auto !important;
  }
}
</style>

