<script setup lang="ts">
import { ref, onMounted } from 'vue'
import TypedText from './TypedText.vue'

const codeLines = ref<string[]>([])
const isVisible = ref(false)

const fullCode = `package main

import (
    "fmt"
    "net/http"
)

func main() {
    http.HandleFunc("/", handler)
    fmt.Println("Server starting on :8080")
    http.ListenAndServe(":8080", nil)
}

func handler(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "Hello from claude-code-Go!")
}`

onMounted(() => {
  const observer = new IntersectionObserver(
    (entries) => {
      entries.forEach((entry) => {
        if (entry.isIntersecting) {
          isVisible.value = true
          animateCode()
          observer.unobserve(entry.target)
        }
      })
    },
    { threshold: 0.2 }
  )

  const element = document.querySelector('.code-preview-container')
  if (element) {
    observer.observe(element)
  }
})

const animateCode = () => {
  const lines = fullCode.split('\n')
  let currentLine = 0

  const addLine = () => {
    if (currentLine < lines.length) {
      codeLines.value.push(lines[currentLine])
      currentLine++
      setTimeout(addLine, 80)
    }
  }

  setTimeout(addLine, 300)
}

const getLineClass = (line: string, index: number): string => {
  if (line.trim().startsWith('//') || line.trim().startsWith('/*')) return 'comment'
  if (line.trim().startsWith('package')) return 'keyword'
  if (line.trim().startsWith('import')) return 'keyword'
  if (line.trim().startsWith('func')) return 'keyword'
  if (line.includes('func ') && !line.trim().startsWith('func')) return 'function'
  if (line.includes('"')) return 'string'
  if (line.includes(':=') || line.includes('=')) return 'variable'
  if (line.trim().startsWith('fmt.')) return 'builtin'
  if (line.trim().startsWith('http.')) return 'builtin'
  return ''
}
</script>

<template>
  <div class="code-preview-container">
    <div class="code-window" :class="{ 'is-visible': isVisible }">
      <div class="code-header">
        <div class="code-dots">
          <span class="code-dot red"></span>
          <span class="code-dot yellow"></span>
          <span class="code-dot green"></span>
        </div>
        <span class="code-filename">main.go</span>
        <div class="code-actions">
          <button class="code-action-btn" title="Copy">
            <svg viewBox="0 0 24 24" width="14" height="14" fill="none" stroke="currentColor" stroke-width="2">
              <rect x="9" y="9" width="13" height="13" rx="2" ry="2"></rect>
              <path d="M5 15H4a2 2 0 0 1-2-2V4a2 2 0 0 1 2-2h9a2 2 0 0 1 2 2v1"></path>
            </svg>
          </button>
        </div>
      </div>
      <div class="code-body">
        <div class="code-content">
          <div
            v-for="(line, index) in codeLines"
            :key="index"
            class="code-line"
            :class="getLineClass(line, index)"
          >
            <span class="line-number">{{ index + 1 }}</span>
            <span class="line-content">{{ line || ' ' }}</span>
          </div>
        </div>
      </div>
    </div>
    
    <div class="code-preview-typing" v-if="isVisible">
      <span class="typing-label">$ go-code</span>
      <TypedText
        :strings="['Generate an HTTP server', 'Create a REST API', 'Build a CLI tool', 'Write unit tests']"
        :typeSpeed="60"
        :backSpeed="40"
        :backDelay="2500"
      />
    </div>
  </div>
</template>

<style scoped>
.code-preview-container {
  width: 100%;
  max-width: 700px;
  margin: 0 auto;
  position: relative;
}

.code-window {
  background: linear-gradient(145deg, #0d1117 0%, #161b22 100%);
  border: 1px solid #30363d;
  border-radius: 12px;
  overflow: hidden;
  box-shadow: 
    0 20px 60px rgba(0, 0, 0, 0.4),
    0 0 0 1px rgba(255, 255, 255, 0.05);
  opacity: 0;
  transform: translateY(30px);
  transition: all 0.6s cubic-bezier(0.16, 1, 0.3, 1);
}

.code-window.is-visible {
  opacity: 1;
  transform: translateY(0);
}

.code-header {
  background: linear-gradient(180deg, #21262d 0%, #161b22 100%);
  padding: 12px 16px;
  display: flex;
  align-items: center;
  gap: 12px;
  border-bottom: 1px solid #30363d;
}

.code-dots {
  display: flex;
  gap: 8px;
}

.code-dot {
  width: 12px;
  height: 12px;
  border-radius: 50%;
  transition: transform 0.2s ease;
}

.code-dot:hover {
  transform: scale(1.2);
}

.code-dot.red { background: #ff5f56; box-shadow: 0 0 6px rgba(255, 95, 86, 0.4); }
.code-dot.yellow { background: #ffbd2e; box-shadow: 0 0 6px rgba(255, 189, 46, 0.4); }
.code-dot.green { background: #27c93f; box-shadow: 0 0 6px rgba(39, 201, 63, 0.4); }

.code-filename {
  color: #8b949e;
  font-size: 13px;
  margin-left: 8px;
  font-family: 'JetBrains Mono', monospace;
}

.code-actions {
  margin-left: auto;
}

.code-action-btn {
  background: transparent;
  border: none;
  color: #8b949e;
  cursor: pointer;
  padding: 4px;
  border-radius: 4px;
  transition: all 0.2s ease;
  display: flex;
  align-items: center;
  justify-content: center;
}

.code-action-btn:hover {
  color: #e6edf3;
  background: rgba(255, 255, 255, 0.1);
}

.code-body {
  padding: 16px 0;
  font-family: 'JetBrains Mono', 'Fira Code', 'SF Mono', Consolas, monospace;
  font-size: 13px;
  line-height: 1.8;
  overflow-x: auto;
}

.code-content {
  min-height: 200px;
}

.code-line {
  display: flex;
  padding: 0 16px;
  white-space: pre;
  animation: fadeInLine 0.3s ease forwards;
}

@keyframes fadeInLine {
  from {
    opacity: 0;
    transform: translateX(-10px);
  }
  to {
    opacity: 1;
    transform: translateX(0);
  }
}

.line-number {
  color: #484f58;
  min-width: 32px;
  text-align: right;
  margin-right: 16px;
  user-select: none;
}

.line-content {
  color: #e6edf3;
  flex: 1;
}

.code-line.keyword .line-content {
  color: #ff7b72;
}

.code-line.function .line-content {
  color: #d2a8ff;
}

.code-line.string .line-content {
  color: #a5d6ff;
}

.code-line.comment .line-content {
  color: #8b949e;
}

.code-line.variable .line-content {
  color: #79c0ff;
}

.code-line.builtin .line-content {
  color: #ffa657;
}

.code-preview-typing {
  margin-top: 24px;
  padding: 16px 20px;
  background: linear-gradient(145deg, #f6f8fa 0%, #ffffff 100%);
  border: 1px solid #d0d7de;
  border-radius: 10px;
  font-family: 'JetBrains Mono', monospace;
  font-size: 14px;
  display: flex;
  align-items: center;
  gap: 12px;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.08);
  opacity: 0;
  transform: translateY(20px);
  animation: fadeInUp 0.5s ease 0.8s forwards;
}

.dark .code-preview-typing {
  background: linear-gradient(145deg, #21262d 0%, #161b22 100%);
  border-color: #30363d;
}

@keyframes fadeInUp {
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

.typing-label {
  color: var(--vp-c-text-2);
  font-weight: 500;
}

:global(.dark) .code-preview-typing {
  background: linear-gradient(145deg, #21262d 0%, #161b22 100%);
  border-color: #30363d;
}
</style>