<script setup lang="ts">
import { ref, onMounted } from 'vue'

interface TerminalLine {
  text: string
  type: 'prompt' | 'command' | 'output' | 'success' | 'thinking' | 'tool' | 'cursor'
  delay?: number
}

const lines: TerminalLine[] = [
  { text: '$ go-code', type: 'prompt', delay: 500 },
  { text: 'claude-code-Go v0.1.0', type: 'output', delay: 300 },
  { text: 'Type /help for commands, /exit to quit.', type: 'output', delay: 200 },
  { text: '', type: 'output', delay: 100 },
  { text: '> 帮我写个 HTTP 服务器', type: 'command', delay: 800 },
  { text: '⠿ Thinking...', type: 'thinking', delay: 1500 },
  { text: '🛠️ Tool call: Write → main.go', type: 'tool', delay: 800 },
  { text: '✓ File written', type: 'success', delay: 600 },
  { text: '', type: 'output', delay: 100 },
  { text: '> ', type: 'cursor', delay: 0 }
]

const displayedLines = ref<TerminalLine[]>([])
const currentIndex = ref(0)

onMounted(() => {
  const processLine = () => {
    if (currentIndex.value < lines.length) {
      const line = lines[currentIndex.value]
      displayedLines.value.push(line)
      currentIndex.value++
      if (line.delay && line.type !== 'cursor') {
        setTimeout(processLine, line.delay)
      } else if (line.type !== 'cursor') {
        setTimeout(processLine, 100)
      }
    }
  }
  setTimeout(processLine, 500)
})
</script>

<template>
  <div class="typewriter-terminal">
    <div class="terminal-window">
      <div class="terminal-header">
        <span class="terminal-dot red"></span>
        <span class="terminal-dot yellow"></span>
        <span class="terminal-dot green"></span>
        <span class="terminal-title">claude-code-Go</span>
      </div>
      <div class="terminal-body">
        <template v-for="(line, index) in displayedLines" :key="index">
          <div v-if="line.type === 'cursor'" class="terminal-line cursor-line">
            <span class="terminal-prompt">go-code></span>
            <span class="terminal-cursor"></span>
          </div>
          <div v-else-if="line.text" class="terminal-line" :class="'terminal-' + line.type">
            <span v-if="line.type === 'prompt'" class="terminal-prompt">$ </span>
            <span v-if="line.type === 'command'" class="terminal-prompt">> </span>
            <span :class="line.type === 'command' ? 'terminal-cmd' : ''">{{ line.text }}</span>
          </div>
          <div v-else class="terminal-line">&nbsp;</div>
        </template>
      </div>
    </div>
  </div>
</template>

<style scoped>
.typewriter-terminal {
  width: 100%;
  max-width: 640px;
  margin: 0 auto;
}

.terminal-window {
  background: #0d1117;
  border: 1px solid #30363d;
  border-radius: 8px;
  overflow: hidden;
  font-family: 'JetBrains Mono', 'Fira Code', 'SF Mono', Consolas, monospace;
  font-size: 13px;
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.4);
}

.terminal-header {
  background: #161b22;
  padding: 10px 14px;
  display: flex;
  align-items: center;
  gap: 8px;
  border-bottom: 1px solid #30363d;
}

.terminal-dot {
  width: 12px;
  height: 12px;
  border-radius: 50%;
}

.terminal-dot.red { background: #ff5f56; }
.terminal-dot.yellow { background: #ffbd2e; }
.terminal-dot.green { background: #27c93f; }

.terminal-title {
  color: #8b949e;
  font-size: 12px;
  margin-left: auto;
}

.terminal-body {
  padding: 16px;
  color: #e6edf3;
  line-height: 1.7;
  min-height: 280px;
}

.terminal-line {
  margin-bottom: 2px;
}

.cursor-line {
  display: flex;
  align-items: center;
}

.terminal-prompt {
  color: #7ee787;
}

.terminal-cmd {
  color: #e6edf3;
}

.terminal-output {
  color: #8b949e;
}

.terminal-success {
  color: #7ee787;
}

.terminal-thinking {
  color: #d2a8ff;
}

.terminal-tool {
  color: #79c0ff;
}

.terminal-cursor {
  display: inline-block;
  width: 8px;
  height: 16px;
  background: #e6edf3;
  animation: blink 1s step-end infinite;
  vertical-align: text-bottom;
}

@keyframes blink {
  0%, 100% { opacity: 1; }
  50% { opacity: 0; }
}
</style>