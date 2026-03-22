<template>
  <div class="sci-fi-container">
    <div class="sci-fi-header">
      <div class="logo">
        <span class="logo-icon">◈</span>
        <span class="logo-text">Agent Go</span>
      </div>
      
      <div class="header-actions">
        <div class="model-selector">
          <span class="selector-label">◉ MODEL</span>
          <select v-model="selectedModel" class="sci-fi-select">
            <option value="qwen">通义千问 (QWEN)</option>
            <option value="deepseek">深度求索 (DeepSeek)</option>
          </select>
        </div>
        
        <div class="stream-toggle">
          <label class="toggle-label">
            <input type="checkbox" v-model="isStreaming" class="toggle-input" />
            <span class="toggle-slider"></span>
            <span class="toggle-text">流式响应</span>
          </label>
        </div>
      </div>

      <div class="bar-right">
        <button 
          class="sync-btn" 
          @click="syncSessions" 
          :disabled="syncing"
          :class="{ 'syncing': syncing }"
        >
          <span class="btn-icon">⟳</span>
          <span v-if="!syncing">同步</span>
          <span v-else>同步中...</span>
        </button>
        
        <button class="logout-btn" @click="logout">
          <span class="btn-icon">⊟</span>
          <span>退出</span>
        </button>
      </div>
    </div>

    <div class="sci-fi-body">
      <div class="sci-fi-sidebar">
        <div class="sidebar-header">
          <h2 class="sidebar-title">对话</h2>
          <button class="new-chat-btn" @click="startNewChat">
            <span class="btn-icon">⊕</span>
            <span>新对话</span>
          </button>
        </div>
        
        <div class="sessions-list">
          <div 
            v-for="(session, id) in sessions" 
            :key="id"
            class="session-item"
            :class="{ 'active': currentSessionId === id }"
            @click="switchSession(id)"
          >
            <div class="session-title">{{ session.name }}</div>
            <button 
              class="delete-btn" 
              @click.stop="deleteSession(id)"
              title="删除会话"
            >
              <span class="btn-icon">⊗</span>
            </button>
          </div>
        </div>
      </div>

      <div class="sci-fi-main">
        <div class="messages-container" ref="messagesContainer">
          <div 
            v-for="(message, index) in currentMessages" 
            :key="index"
            :class="['message-wrapper', message.role === 'user' ? 'user' : 'ai']"
          >
            <div class="message-avatar">
              <span class="avatar-icon">{{ message.role === 'user' ? '◉' : '◈' }}</span>
            </div>
            
            <div class="message-content-wrapper">
              <div class="message-header">
                <span class="sender-name">{{ message.role === 'user' ? '用户' : 'AI助手' }}</span>
                <span v-if="message.meta && message.meta.status === 'streaming'" class="streaming-dot"></span>
              </div>
              
              <div class="message-bubble">
                <div class="bubble-content" v-html="renderMarkdown(message.content)"></div>
                <div v-if="message.role === 'assistant'" class="message-actions">
                  <button class="action-btn tts-btn" @click="playTTS(message.content)" title="语音播放">
                    <span class="btn-icon">◊</span>
                  </button>
                </div>
              </div>
            </div>
          </div>
        </div>

        <div class="input-area">
          <div class="input-wrapper">
            <div class="input-frame">
              <div class="frame-corner tl"></div>
              <div class="frame-corner tr"></div>
              <div class="frame-corner bl"></div>
              <div class="frame-corner br"></div>
            </div>
            
            <textarea
              v-model="inputMessage"
              class="sci-fi-textarea"
              placeholder="输入您的问题..."
              @keydown.enter.exact.prevent="sendMessage"
              :disabled="loading"
              ref="messageInput"
              rows="1"
            ></textarea>
            
            <button
              type="button"
              class="send-btn"
              :disabled="!inputMessage.trim() || loading"
              :class="{ 'loading': loading }"
              @click="sendMessage"
            >
              <span class="btn-icon">⟩</span>
            </button>
          </div>
          
          <div class="input-footer">
            <span class="model-badge">{{ selectedModel === 'qwen' ? 'QWEN' : 'DEEPSEEK' }}</span>
            <div class="input-info">
              <span v-if="loading" class="loading-text">AI正在思考...</span>
              <span v-else class="ready-text">就绪</span>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, watch, nextTick } from 'vue'
import api from '../utils/api.js'
import { marked } from 'marked'
import hljs from 'highlight.js'
import 'highlight.js/styles/github.css'

const inputMessage = ref('')
const sessions = ref({})
const currentSessionId = ref('')
const loading = ref(false)
const syncing = ref(false)
const messagesContainer = ref(null)
const messageInput = ref(null)
const tempSession = ref(false)
const isStreaming = ref(true)
const selectedModel = ref('qwen')

const currentMessages = ref([])

const showMessage = (message, type = 'info') => {
  // 简单的消息提示，实际项目中可使用更完善的组件
  console.log(`${type}: ${message}`)
  // 可以在这里集成Element Plus的Message组件
}

const scrollToBottom = () => {
  if (messagesContainer.value) {
    messagesContainer.value.scrollTop = messagesContainer.value.scrollHeight
  }
}

const renderMarkdown = (content) => {
  if (!content) return ''
  
  marked.setOptions({
    highlight: function(code, lang) {
      if (lang && hljs.getLanguage(lang)) {
        return hljs.highlight(code, { language: lang }).value
      }
      return hljs.highlightAuto(code).value
    },
    breaks: true,
    gfm: true
  })
  
  return marked(content)
}

const startNewChat = () => {
  tempSession.value = true
  currentSessionId.value = ''
  currentMessages.value = []
  nextTick(() => {
    messageInput.value?.focus()
  })
}

const switchSession = (id) => {
  currentSessionId.value = id
  if (sessions.value[id]) {
    currentMessages.value = sessions.value[id].messages || []
  }
  nextTick(() => {
    scrollToBottom()
  })
}

const deleteSession = (id) => {
  if (confirm('确定要删除这个会话吗？')) {
    delete sessions.value[id]
    if (currentSessionId.value === id) {
      startNewChat()
    }
  }
}

const syncSessions = async () => {
  syncing.value = true
  try {
    const response = await api.get('/AI/chat/sessions')
    if (response.data && response.data.status_code === 1000) {
      const sessionList = response.data.sessions || []
      const newSessions = {}
      
      sessionList.forEach(session => {
        newSessions[session.SessionID] = {
          id: session.SessionID,
          name: session.Title || '未命名会话',
          messages: []
        }
      })
      
      sessions.value = newSessions
      
      if (Object.keys(newSessions).length > 0 && !currentSessionId.value) {
        const firstSessionId = Object.keys(newSessions)[0]
        switchSession(firstSessionId)
      }
      
      showMessage('会话同步成功')
    }
  } catch (error) {
    console.error('同步会话失败:', error)
    showMessage('同步会话失败', 'error')
  } finally {
    syncing.value = false
  }
}

const sendMessage = async () => {
  const question = inputMessage.value.trim()
  if (!question) return
  
  inputMessage.value = ''
  loading.value = true
  
  // 添加用户消息到当前会话
  currentMessages.value.push({ role: 'user', content: question })
  
  if (!tempSession.value && currentSessionId.value && sessions.value[currentSessionId.value]) {
    if (!sessions.value[currentSessionId.value].messages) sessions.value[currentSessionId.value].messages = []
    sessions.value[currentSessionId.value].messages.push({ role: 'user', content: question })
  }
  
  scrollToBottom()
  
  // 添加AI消息占位符
  const aiMessage = {
    role: 'assistant',
    content: '',
    meta: { status: 'streaming' }
  }
  
  currentMessages.value.push(aiMessage)
  const aiMessageIndex = currentMessages.value.length - 1

  if (!tempSession.value && currentSessionId.value && sessions.value[currentSessionId.value]) {
    if (!sessions.value[currentSessionId.value].messages) sessions.value[currentSessionId.value].messages = []
    sessions.value[currentSessionId.value].messages.push({ role: 'assistant', content: '' })
  }

  if (isStreaming.value) {
    await handleStreaming(question, aiMessageIndex)
  } else {
    await handleNormal(question, aiMessageIndex)
  }
}

const handleStreaming = async (question, aiMessageIndex) => {
  const url = tempSession.value
    ? '/api/AI/chat/send-stream-new-session'
    : '/api/AI/chat/send-stream'

  const headers = {
    'Content-Type': 'application/json',
    'Authorization': `Bearer ${localStorage.getItem('token') || ''}`
  }

  const body = tempSession.value
    ? { question: question, modelType: selectedModel.value }
    : { question: question, modelType: selectedModel.value, sessionId: currentSessionId.value }

  try {
    const response = await fetch(url, {
      method: 'POST',
      headers,
      body: JSON.stringify(body)
    })

    if (!response.ok) {
      loading.value = false
      throw new Error('Network response was not ok')
    }

    const reader = response.body.getReader()
    const decoder = new TextDecoder()
    let buffer = ''

    for (;;) {
      const { done, value } = await reader.read()
      if (done) break

      const chunk = decoder.decode(value, { stream: true })
      buffer += chunk

      const lines = buffer.split('\n')
      buffer = lines.pop() || ''

      for (const line of lines) {
        const trimmedLine = line.trim()
        if (!trimmedLine) continue

        if (trimmedLine.startsWith('data:')) {
          const data = trimmedLine.slice(5).trim()

          if (data === '[DONE]') {
            loading.value = false
            currentMessages.value[aiMessageIndex].meta = { status: 'done' }
            currentMessages.value = [...currentMessages.value]
          } else if (data.startsWith('{')) {
            try {
              const parsed = JSON.parse(data)
              if (parsed.sessionId) {
                const newSid = String(parsed.sessionId)
                if (tempSession.value) {
                  sessions.value[newSid] = {
                    id: newSid,
                    name: '新会话',
                    messages: [...currentMessages.value]
                  }
                  currentSessionId.value = newSid
                  tempSession.value = false
                }
              }
            } catch (e) {
              currentMessages.value[aiMessageIndex].content += data
            }
          } else {
            currentMessages.value[aiMessageIndex].content += data
          }

          currentMessages.value = [...currentMessages.value]
          
          await new Promise(resolve => {
            requestAnimationFrame(() => {
              scrollToBottom()
              resolve()
            })
          })
        }
      }
    }

    loading.value = false
    currentMessages.value[aiMessageIndex].meta = { status: 'done' }
    currentMessages.value = [...currentMessages.value]

    if (!tempSession.value && currentSessionId.value && sessions.value[currentSessionId.value]) {
      const sessMsgs = sessions.value[currentSessionId.value].messages
      if (Array.isArray(sessMsgs) && sessMsgs.length) {
        const lastIndex = sessMsgs.length - 1
        if (sessMsgs[lastIndex] && sessMsgs[lastIndex].role === 'assistant') {
          sessMsgs[lastIndex].content = currentMessages.value[aiMessageIndex].content
        }
      }
    }
  } catch (err) {
    console.error('Stream error:', err)
    loading.value = false
    currentMessages.value[aiMessageIndex].meta = { status: 'error' }
    currentMessages.value = [...currentMessages.value]
    showMessage('流式传输出错', 'error')
  }
}

const handleNormal = async (question, aiMessageIndex) => {
  if (tempSession.value) {
    const response = await api.post('/api/AI/chat/send-new-session', {
      question: question,
      modelType: selectedModel.value
    })
    if (response.data && response.data.status_code === 1000) {
      const sessionId = String(response.data.sessionId)
      const aiMessage = {
        role: 'assistant',
        content: response.data.Information || ''
      }

      sessions.value[sessionId] = {
        id: sessionId,
        name: '新会话',
        messages: [ { role: 'user', content: question }, aiMessage ]
      }
      currentSessionId.value = sessionId
      tempSession.value = false
      currentMessages.value = [...sessions.value[sessionId].messages]
    } else {
      showMessage(response.data?.status_msg || '发送失败', 'error')
      // 使用 aiMessageIndex 删除占位符消息
      if (aiMessageIndex >= 0 && aiMessageIndex < currentMessages.value.length) {
        currentMessages.value.splice(aiMessageIndex, 1)
      }
    }
  } else {
    const sessionMsgs = sessions.value[currentSessionId.value].messages
    sessionMsgs.push({ role: 'user', content: question })

    const response = await api.post('/api/AI/chat/send', {
      question: question,
      modelType: selectedModel.value,
      sessionId: currentSessionId.value
    })
    if (response.data && response.data.status_code === 1000) {
      const aiMessage = {
        role: 'assistant',
        content: response.data.Information || ''
      }
      sessionMsgs.push(aiMessage)
      // 使用 aiMessageIndex 更新占位符消息
      if (aiMessageIndex >= 0 && aiMessageIndex < currentMessages.value.length) {
        currentMessages.value[aiMessageIndex].content = aiMessage.content
        currentMessages.value[aiMessageIndex].meta = { status: 'done' }
      } else {
        currentMessages.value.push(aiMessage)
      }
    } else {
      showMessage(response.data?.status_msg || '发送失败', 'error')
      // 使用 aiMessageIndex 删除占位符消息
      if (aiMessageIndex >= 0 && aiMessageIndex < currentMessages.value.length) {
        currentMessages.value.splice(aiMessageIndex, 1)
      }
      sessionMsgs.pop()
    }
  }
  loading.value = false
  scrollToBottom()
}

const logout = () => {
  localStorage.removeItem('token')
  window.location.href = '/login'
}

const playTTS = (text) => {
  if ('speechSynthesis' in window) {
    const utterance = new SpeechSynthesisUtterance(text)
    utterance.lang = 'zh-CN'
    speechSynthesis.speak(utterance)
  } else {
    showMessage('浏览器不支持语音播放', 'warning')
  }
}

onMounted(() => {
  syncSessions()
  nextTick(() => {
    messageInput.value?.focus()
  })
})

watch(currentMessages, () => {
  nextTick(() => {
    scrollToBottom()
  })
}, { deep: true })
</script>

<style scoped>
.sci-fi-container {
  width: 100vw;
  height: 100vh;
  display: flex;
  flex-direction: column;
  background: linear-gradient(135deg, #f0f4f8 0%, #e9ecef 100%);
  font-family: 'Rajdhani', sans-serif;
  overflow: hidden;
  position: relative;
}

.sci-fi-container::before {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: 
    radial-gradient(circle at 25% 25%, rgba(16, 185, 129, 0.05) 0%, transparent 50%),
    radial-gradient(circle at 75% 75%, rgba(16, 185, 129, 0.05) 0%, transparent 50%);
  pointer-events: none;
  z-index: 0;
}

.sci-fi-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 16px 24px;
  background: rgba(255, 255, 255, 0.95);
  backdrop-filter: blur(10px);
  border-bottom: 1px solid rgba(16, 185, 129, 0.2);
  box-shadow: 0 2px 12px rgba(16, 185, 129, 0.1);
  z-index: 100;
  position: relative;
}

.logo {
  display: flex;
  align-items: center;
  gap: 12px;
  font-family: 'Orbitron', sans-serif;
}

.logo-icon {
  font-size: 24px;
  color: #10B981;
  text-shadow: 0 0 8px rgba(16, 185, 129, 0.4);
}

.logo-text {
  font-size: 20px;
  font-weight: 700;
  color: #10B981;
  letter-spacing: 2px;
  text-shadow: 0 1px 2px rgba(0, 0, 0, 0.1);
}

.header-actions {
  display: flex;
  align-items: center;
  gap: 24px;
}

.model-selector {
  display: flex;
  align-items: center;
  gap: 8px;
}

.selector-label {
  font-size: 12px;
  font-weight: 600;
  color: #10B981;
  letter-spacing: 1px;
  text-transform: uppercase;
}

.sci-fi-select {
  padding: 8px 12px;
  background: rgba(255, 255, 255, 0.9);
  border: 1px solid rgba(16, 185, 129, 0.3);
  border-radius: 6px;
  font-family: 'Rajdhani', sans-serif;
  font-size: 14px;
  color: #10B981;
  cursor: pointer;
  transition: all 0.3s ease;
  box-shadow: 0 2px 4px rgba(16, 185, 129, 0.1);
}

.sci-fi-select:hover {
  border-color: #10B981;
  box-shadow: 0 0 12px rgba(16, 185, 129, 0.2);
}

.stream-toggle {
  display: flex;
  align-items: center;
}

.toggle-label {
  display: flex;
  align-items: center;
  gap: 8px;
  cursor: pointer;
}

.toggle-input {
  display: none;
}

.toggle-slider {
  width: 44px;
  height: 24px;
  background: rgba(16, 185, 129, 0.2);
  border-radius: 12px;
  position: relative;
  transition: all 0.3s ease;
}

.toggle-slider::before {
  content: '';
  position: absolute;
  top: 2px;
  left: 2px;
  width: 20px;
  height: 20px;
  background: white;
  border-radius: 50%;
  transition: all 0.3s ease;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
}

.toggle-input:checked + .toggle-slider {
  background: #10B981;
  box-shadow: 0 0 12px rgba(16, 185, 129, 0.4);
}

.toggle-input:checked + .toggle-slider::before {
  transform: translateX(20px);
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.2);
}

.toggle-text {
  font-size: 14px;
  color: #10B981;
  font-weight: 500;
}

.bar-right {
  display: flex;
  align-items: center;
  gap: 12px;
}

.sync-btn, .logout-btn {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 8px 16px;
  background: rgba(255, 255, 255, 0.9);
  border: 1px solid rgba(16, 185, 129, 0.3);
  border-radius: 6px;
  font-family: 'Rajdhani', sans-serif;
  font-size: 14px;
  color: #10B981;
  cursor: pointer;
  transition: all 0.3s ease;
  box-shadow: 0 2px 4px rgba(16, 185, 129, 0.1);
}

.sync-btn:hover, .logout-btn:hover {
  border-color: #10B981;
  box-shadow: 0 0 12px rgba(16, 185, 129, 0.2);
  background: rgba(16, 185, 129, 0.05);
}

.sync-btn:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

.sync-btn.syncing {
  animation: pulse 1.5s infinite;
}

@keyframes pulse {
  0% {
    box-shadow: 0 0 0 0 rgba(16, 185, 129, 0.4);
  }
  70% {
    box-shadow: 0 0 0 10px rgba(16, 185, 129, 0);
  }
  100% {
    box-shadow: 0 0 0 0 rgba(16, 185, 129, 0);
  }
}

.btn-icon {
  font-size: 16px;
  font-weight: bold;
}

.sci-fi-body {
  flex: 1;
  display: flex;
  overflow: hidden;
  position: relative;
  z-index: 1;
}

.sci-fi-sidebar {
  width: 280px;
  background: rgba(255, 255, 255, 0.9);
  backdrop-filter: blur(10px);
  border-right: 1px solid rgba(16, 185, 129, 0.2);
  display: flex;
  flex-direction: column;
  box-shadow: 2px 0 12px rgba(16, 185, 129, 0.1);
  z-index: 10;
}

.sidebar-header {
  padding: 20px;
  border-bottom: 1px solid rgba(16, 185, 129, 0.2);
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.sidebar-title {
  font-family: 'Orbitron', sans-serif;
  font-size: 16px;
  font-weight: 700;
  color: #10B981;
  letter-spacing: 1px;
  margin: 0;
}

.new-chat-btn {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 6px 12px;
  background: #10B981;
  border: none;
  border-radius: 6px;
  font-family: 'Rajdhani', sans-serif;
  font-size: 14px;
  color: white;
  cursor: pointer;
  transition: all 0.3s ease;
  box-shadow: 0 2px 4px rgba(16, 185, 129, 0.3);
}

.new-chat-btn:hover {
  background: #059669;
  box-shadow: 0 4px 8px rgba(16, 185, 129, 0.4);
  transform: translateY(-1px);
}

.sessions-list {
  flex: 1;
  overflow-y: auto;
  padding: 12px;
}

.session-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 12px 16px;
  margin-bottom: 8px;
  background: rgba(255, 255, 255, 0.8);
  border: 1px solid rgba(16, 185, 129, 0.2);
  border-radius: 8px;
  cursor: pointer;
  transition: all 0.3s ease;
  box-shadow: 0 2px 4px rgba(16, 185, 129, 0.05);
}

.session-item:hover {
  border-color: #10B981;
  box-shadow: 0 0 12px rgba(16, 185, 129, 0.15);
  background: rgba(16, 185, 129, 0.02);
}

.session-item.active {
  background: rgba(16, 185, 129, 0.1);
  border-color: #10B981;
  box-shadow: 0 0 16px rgba(16, 185, 129, 0.2);
}

.session-title {
  flex: 1;
  font-size: 14px;
  color: #10B981;
  font-weight: 500;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  margin-right: 8px;
}

.delete-btn {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 24px;
  height: 24px;
  background: rgba(239, 68, 68, 0.1);
  border: 1px solid rgba(239, 68, 68, 0.3);
  border-radius: 4px;
  color: #ef4444;
  cursor: pointer;
  transition: all 0.3s ease;
  opacity: 0;
}

.session-item:hover .delete-btn {
  opacity: 1;
}

.delete-btn:hover {
  background: #ef4444;
  color: white;
  box-shadow: 0 0 8px rgba(239, 68, 68, 0.3);
}

.sci-fi-main {
  flex: 1;
  display: flex;
  flex-direction: column;
  overflow: hidden;
  position: relative;
}

.messages-container {
  flex: 1;
  overflow-y: auto;
  padding: 24px;
  display: flex;
  flex-direction: column;
  gap: 16px;
  scroll-behavior: smooth;
}

.messages-container::-webkit-scrollbar {
  width: 8px;
}

.messages-container::-webkit-scrollbar-track {
  background: rgba(16, 185, 129, 0.1);
  border-radius: 4px;
}

.messages-container::-webkit-scrollbar-thumb {
  background: rgba(16, 185, 129, 0.3);
  border-radius: 4px;
  transition: all 0.3s ease;
}

.messages-container::-webkit-scrollbar-thumb:hover {
  background: rgba(16, 185, 129, 0.5);
}

.message-wrapper {
  display: flex;
  gap: 12px;
  max-width: 80%;
  animation: slideIn 0.3s ease;
}

@keyframes slideIn {
  from {
    opacity: 0;
    transform: translateY(10px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

.message-wrapper.user {
  align-self: flex-end;
  flex-direction: row-reverse;
}

.message-avatar {
  width: 36px;
  height: 36px;
  border-radius: 50%;
  background: rgba(16, 185, 129, 0.1);
  display: flex;
  align-items: center;
  justify-content: center;
  border: 1px solid rgba(16, 185, 129, 0.3);
  flex-shrink: 0;
  box-shadow: 0 2px 4px rgba(16, 185, 129, 0.1);
}

.message-wrapper.user .message-avatar {
  background: rgba(16, 185, 129, 0.2);
  border-color: #10B981;
}

.avatar-icon {
  font-size: 16px;
  font-weight: bold;
  color: #10B981;
  text-shadow: 0 0 4px rgba(16, 185, 129, 0.3);
}

.message-content-wrapper {
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.message-header {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-bottom: 4px;
}

.sender-name {
  font-size: 12px;
  font-weight: 600;
  color: #10B981;
  letter-spacing: 0.5px;
  text-transform: uppercase;
}

.streaming-dot {
  width: 8px;
  height: 8px;
  background: #10B981;
  border-radius: 50%;
  animation: pulse 1.5s infinite;
  box-shadow: 0 0 8px rgba(16, 185, 129, 0.4);
}

.message-bubble {
  padding: 16px 20px;
  border-radius: 16px;
  position: relative;
  backdrop-filter: blur(5px);
  box-shadow: 0 2px 8px rgba(16, 185, 129, 0.1);
  border: 1px solid rgba(16, 185, 129, 0.2);
}

.message-wrapper.user .message-bubble {
  background: rgba(16, 185, 129, 0.1);
  border-color: rgba(16, 185, 129, 0.3);
  border-bottom-right-radius: 4px;
}

.message-wrapper.ai .message-bubble {
  background: rgba(255, 255, 255, 0.98);
  border-color: rgba(16, 185, 129, 0.4);
  border-bottom-left-radius: 4px;
  box-shadow: 0 2px 12px rgba(16, 185, 129, 0.15);
}

.bubble-content {
  font-size: 16px;
  line-height: 1.6;
  color: #1f2937;
  word-wrap: break-word;
}

.bubble-content p {
  margin: 0 0 8px 0;
}

.bubble-content p:last-child {
  margin-bottom: 0;
}

.bubble-content code {
  background: rgba(16, 185, 129, 0.1);
  padding: 2px 4px;
  border-radius: 4px;
  font-family: 'Courier New', monospace;
  font-size: 0.9em;
  color: #10B981;
}

.bubble-content pre {
  background: rgba(16, 185, 129, 0.05);
  padding: 12px;
  border-radius: 8px;
  overflow-x: auto;
  border: 1px solid rgba(16, 185, 129, 0.2);
  margin: 8px 0;
}

.bubble-content pre code {
  background: transparent;
  padding: 0;
  color: #1f2937;
}

.bubble-content blockquote {
  border-left: 4px solid #10B981;
  padding-left: 12px;
  margin: 8px 0;
  color: #4b5563;
  font-style: italic;
}

.message-actions {
  display: flex;
  gap: 8px;
  margin-top: 8px;
  padding-top: 8px;
  border-top: 1px solid rgba(16, 185, 129, 0.1);
}

.action-btn {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 32px;
  height: 32px;
  background: rgba(16, 185, 129, 0.1);
  border: 1px solid rgba(16, 185, 129, 0.3);
  border-radius: 6px;
  color: #10B981;
  cursor: pointer;
  transition: all 0.3s ease;
}

.action-btn:hover {
  background: #10B981;
  color: white;
  box-shadow: 0 0 8px rgba(16, 185, 129, 0.3);
}

.input-area {
  padding: 24px;
  background: rgba(255, 255, 255, 0.95);
  backdrop-filter: blur(10px);
  border-top: 1px solid rgba(16, 185, 129, 0.2);
  box-shadow: 0 -2px 12px rgba(16, 185, 129, 0.1);
  position: relative;
  z-index: 10;
}

.input-wrapper {
  position: relative;
  margin-bottom: 12px;
}

.input-frame {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  border: 1px solid rgba(16, 185, 129, 0.3);
  border-radius: 12px;
  pointer-events: none;
  box-shadow: 0 0 12px rgba(16, 185, 129, 0.1);
}

.frame-corner {
  position: absolute;
  width: 12px;
  height: 12px;
  border: 2px solid #10B981;
  border-radius: 4px;
}

.frame-corner.tl {
  top: -2px;
  left: -2px;
  border-right: none;
  border-bottom: none;
  box-shadow: 0 0 8px rgba(16, 185, 129, 0.4);
}

.frame-corner.tr {
  top: -2px;
  right: -2px;
  border-left: none;
  border-bottom: none;
  box-shadow: 0 0 8px rgba(16, 185, 129, 0.4);
}

.frame-corner.bl {
  bottom: -2px;
  left: -2px;
  border-right: none;
  border-top: none;
  box-shadow: 0 0 8px rgba(16, 185, 129, 0.4);
}

.frame-corner.br {
  bottom: -2px;
  right: -2px;
  border-left: none;
  border-top: none;
  box-shadow: 0 0 8px rgba(16, 185, 129, 0.4);
}

.sci-fi-textarea {
  width: 100%;
  min-height: 80px;
  max-height: 200px;
  padding: 16px 60px 16px 16px;
  background: rgba(255, 255, 255, 0.9);
  border: 1px solid rgba(16, 185, 129, 0.2);
  border-radius: 12px;
  font-family: 'Rajdhani', sans-serif;
  font-size: 16px;
  line-height: 1.5;
  color: #1f2937;
  resize: none;
  transition: all 0.3s ease;
  box-shadow: inset 0 2px 4px rgba(16, 185, 129, 0.05);
}

.sci-fi-textarea:focus {
  outline: none;
  border-color: #10B981;
  box-shadow: inset 0 2px 4px rgba(16, 185, 129, 0.1), 0 0 12px rgba(16, 185, 129, 0.2);
  background: rgba(255, 255, 255, 0.95);
}

.sci-fi-textarea:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

.send-btn {
  position: absolute;
  bottom: 8px;
  right: 8px;
  width: 44px;
  height: 44px;
  background: #10B981;
  border: none;
  border-radius: 50%;
  color: white;
  cursor: pointer;
  transition: all 0.3s ease;
  box-shadow: 0 2px 8px rgba(16, 185, 129, 0.3);
  display: flex;
  align-items: center;
  justify-content: center;
}

.send-btn:hover:not(:disabled) {
  background: #059669;
  box-shadow: 0 4px 12px rgba(16, 185, 129, 0.4);
  transform: translateY(-1px);
}

.send-btn:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

.send-btn.loading {
  animation: pulse 1.5s infinite;
}

.btn-icon {
  font-size: 20px;
  font-weight: bold;
}

.input-footer {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.model-badge {
  padding: 4px 12px;
  background: rgba(16, 185, 129, 0.1);
  border: 1px solid rgba(16, 185, 129, 0.3);
  border-radius: 12px;
  font-size: 12px;
  font-weight: 600;
  color: #10B981;
  letter-spacing: 1px;
  text-transform: uppercase;
}

.input-info {
  font-size: 12px;
  color: #6b7280;
}

.loading-text {
  color: #10B981;
  font-weight: 500;
  animation: pulse 1.5s infinite;
}

.ready-text {
  color: #10B981;
  font-weight: 500;
}

@media (max-width: 768px) {
  .sci-fi-sidebar {
    width: 240px;
  }
  
  .message-wrapper {
    max-width: 90%;
  }
  
  .sci-fi-header {
    padding: 12px 16px;
  }
  
  .logo-text {
    font-size: 18px;
  }
  
  .header-actions {
    gap: 16px;
  }
  
  .model-selector {
    flex-direction: column;
    align-items: flex-start;
    gap: 4px;
  }
  
  .sci-fi-select {
    font-size: 12px;
    padding: 6px 10px;
  }
  
  .input-area {
    padding: 16px;
  }
  
  .messages-container {
    padding: 16px;
  }
}

@media (max-width: 480px) {
  .sci-fi-sidebar {
    width: 100%;
    position: absolute;
    left: -100%;
    transition: left 0.3s ease;
    z-index: 100;
  }
  
  .sci-fi-sidebar.active {
    left: 0;
  }
  
  .message-wrapper {
    max-width: 95%;
  }
  
  .sci-fi-header {
    padding: 10px 12px;
  }
  
  .logo-text {
    font-size: 16px;
  }
  
  .header-actions {
    gap: 12px;
  }
  
  .sync-btn, .logout-btn {
    padding: 6px 12px;
    font-size: 12px;
  }
  
  .input-area {
    padding: 12px;
  }
  
  .sci-fi-textarea {
    min-height: 60px;
    padding: 12px 50px 12px 12px;
  }
  
  .send-btn {
    width: 36px;
    height: 36px;
  }
  
  .btn-icon {
    font-size: 16px;
  }
}
</style>