<template>
  <div class="preview-frame">
    <div class="preview-topbar">
      <div class="preview-lights">
        <span></span>
        <span></span>
        <span></span>
      </div>
      <div class="preview-chip">LIVE AI CHAT</div>
    </div>

    <div ref="threadRef" class="preview-thread">
      <div
        v-for="(message, index) in messages"
        :key="`${message.role}-${index}`"
        class="message-row"
        :class="[message.role, { 'first-entry': index === 0 }]"
      >
        <div class="message-bubble">
          <span>{{ message.text }}</span>
          <span
            v-if="message.role === 'ai' && index === loopingMessageIndex"
            class="streaming-status"
          >
            {{ streamingHint }}
          </span>
          <span v-if="message.role === 'ai' && !message.done" class="typing-cursor"></span>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { nextTick, onBeforeUnmount, onMounted, ref } from 'vue'

const threadRef = ref(null)
const messages = ref([])
const isActive = ref(true)
const streamingHint = ref('')
const loopingMessageIndex = ref(-1)
const timers = new Set()

const script = [
  {
    role: 'user',
    text: '',
    fullText: '周末想看点轻松的电影，有推荐吗？',
    done: true
  },
  {
    role: 'ai',
    text: '',
    fullText: '如果你想放松一点，我会先推荐《帕丁顿熊2》，温暖、好笑，而且节奏很舒服。',
    done: false
  },
  {
    role: 'user',
    text: '',
    fullText: '最好不要太烧脑，晚上只想随便看看。',
    done: true
  },
  {
    role: 'ai',
    text: '',
    fullText: '那《落魄大厨》也很适合，氛围轻松，还会让你看着看着就想点夜宵。',
    done: false
  },
  {
    role: 'user',
    text: '',
    fullText: '哈哈，那再来一部适合雨天宅家的。',
    done: true
  },
  {
    role: 'ai',
    text: '',
    fullText: '如果是雨天宅家，我会想继续给你推《爱在黎明破晓前》，聊天感很强，而且越看越安静……',
    done: false
  }
]

const wait = (ms) =>
  new Promise((resolve) => {
    const timer = window.setTimeout(() => {
      timers.delete(timer)
      resolve()
    }, ms)
    timers.add(timer)
  })

const randomDelay = (min, max) => Math.floor(Math.random() * (max - min + 1)) + min

const scrollToBottom = async () => {
  await nextTick()
  if (threadRef.value) {
    threadRef.value.scrollTop = threadRef.value.scrollHeight
  }
}

const typeWriter = async (messageIndex, options = {}) => {
  const { keepGenerating = false } = options
  const target = messages.value[messageIndex]
  if (!target) {
    return
  }

  const finalLength = keepGenerating
    ? Math.max(target.fullText.length - 6, Math.floor(target.fullText.length * 0.8))
    : target.fullText.length

  for (let i = 0; i < finalLength; i += 1) {
    if (!isActive.value) {
      return
    }

    target.text += target.fullText[i]
    await scrollToBottom()
    await wait(randomDelay(28, 52))
  }

  if (!keepGenerating) {
    target.done = true
  }
}

const loopStreamingStatus = async (messageIndex) => {
  loopingMessageIndex.value = messageIndex
  const stages = ['继续生成中', '继续生成中.', '继续生成中..', '继续生成中...']
  let stageIndex = 0

  while (isActive.value && loopingMessageIndex.value === messageIndex) {
    streamingHint.value = stages[stageIndex % stages.length]
    stageIndex += 1
    await wait(420)
  }
}

const playConversation = async () => {
  messages.value = []
  streamingHint.value = ''
  loopingMessageIndex.value = -1

  for (let i = 0; i < script.length; i += 1) {
    if (!isActive.value) {
      return
    }

    const source = script[i]
    const item = {
      role: source.role,
      text: source.role === 'user' ? source.fullText : '',
      fullText: source.fullText,
      done: source.role === 'user'
    }

    messages.value.push(item)
    await scrollToBottom()

    if (source.role === 'user') {
      continue
    }

    await wait(randomDelay(400, 800))

    const isLastAI = i === script.length - 1
    await typeWriter(messages.value.length - 1, { keepGenerating: isLastAI })

    if (!isLastAI) {
      await wait(500)
    } else {
      void loopStreamingStatus(messages.value.length - 1)
    }
  }
}

onMounted(() => {
  playConversation()
})

onBeforeUnmount(() => {
  isActive.value = false
  timers.forEach((timer) => window.clearTimeout(timer))
  timers.clear()
})
</script>

<style scoped>
.preview-frame {
  position: relative;
  border-radius: 30px;
  background:
    linear-gradient(180deg, rgba(255, 255, 255, 0.11), rgba(255, 255, 255, 0.04)),
    var(--panel);
  border: 1px solid var(--border);
  box-shadow:
    0 26px 80px rgba(2, 8, 6, 0.42),
    inset 0 1px 0 rgba(243, 247, 244, 0.08);
  backdrop-filter: blur(22px);
  overflow: hidden;
}

.preview-frame::before {
  content: '';
  position: absolute;
  inset: 0;
  pointer-events: none;
  background:
    radial-gradient(circle at top right, rgba(168, 213, 187, 0.12), transparent 26%),
    linear-gradient(140deg, rgba(143, 191, 167, 0.08), transparent 34%);
}

.preview-topbar,
.preview-thread {
  position: relative;
  z-index: 1;
}

.preview-topbar {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 18px 22px;
  border-bottom: 1px solid var(--border);
  background: rgba(15, 22, 20, 0.44);
}

.preview-lights {
  display: flex;
  gap: 8px;
}

.preview-lights span {
  width: 10px;
  height: 10px;
  border-radius: 50%;
}

.preview-lights span:first-child {
  background: rgba(248, 113, 113, 0.7);
}

.preview-lights span:nth-child(2) {
  background: rgba(250, 204, 21, 0.72);
}

.preview-lights span:nth-child(3) {
  background: rgba(168, 213, 187, 0.72);
}

.preview-chip {
  padding: 8px 12px;
  border-radius: 999px;
  font-family: 'Orbitron', sans-serif;
  font-size: 10px;
  letter-spacing: 1.8px;
  color: var(--accent-strong);
  background: var(--accent-soft);
  border: 1px solid var(--border);
}

.preview-thread {
  display: flex;
  flex-direction: column;
  gap: 22px;
  max-height: 560px;
  padding: 26px 28px 30px;
  overflow: auto;
}

.message-row {
  display: flex;
}

.message-row.first-entry {
  animation: firstMessageFadeIn 0.72s cubic-bezier(0.22, 1, 0.36, 1) both;
}

.message-row.user {
  justify-content: flex-end;
}

.message-row.ai {
  justify-content: flex-start;
}

.message-bubble {
  max-width: 78%;
  padding: 18px 20px;
  border-radius: 26px;
  line-height: 1.8;
  color: var(--text-1);
  border: 1px solid var(--border);
  background: rgba(15, 22, 20, 0.52);
  box-shadow: inset 0 1px 0 rgba(243, 247, 244, 0.06);
}

.message-row.user .message-bubble {
  background:
    linear-gradient(135deg, rgba(143, 191, 167, 0.2), rgba(168, 213, 187, 0.08)),
    rgba(15, 22, 20, 0.56);
  color: var(--text-1);
}

.message-row.ai .message-bubble {
  background: rgba(255, 255, 255, 0.04);
  color: var(--text-2);
}

.streaming-status {
  margin-left: 8px;
  color: var(--accent-strong);
  white-space: nowrap;
}

.typing-cursor {
  display: inline-block;
  width: 9px;
  height: 1.1em;
  margin-left: 4px;
  vertical-align: -0.16em;
  border-radius: 999px;
  background: var(--accent-strong);
  box-shadow: 0 0 12px rgba(168, 213, 187, 0.42);
  animation: blink 1s infinite;
}

@keyframes firstMessageFadeIn {
  0% {
    opacity: 0;
    transform: translateY(12px);
  }

  100% {
    opacity: 1;
    transform: translateY(0);
  }
}

@keyframes blink {
  0%,
  45% {
    opacity: 1;
  }

  55%,
  100% {
    opacity: 0;
  }
}

@media (max-width: 560px) {
  .preview-thread {
    max-height: 460px;
    padding: 20px;
  }

  .message-bubble {
    max-width: 88%;
    padding: 16px 18px;
  }
}
</style>


