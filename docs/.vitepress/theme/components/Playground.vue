<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import AnimatedButton from './AnimatedButton.vue'

const code = ref(`// Welcome to claude-code-Go Playground!
// Try editing this code and see the result

func greet(name string) string {
    return "Hello, " + name + "!"
}

func main() {
    message := greet("World")
    println(message)
}`)

const output = ref('')
const isRunning = ref(false)
const copied = ref(false)
const selectedExample = ref('hello')

const examples = {
  hello: {
    name: 'Hello World',
    code: `// Welcome to claude-code-Go Playground!
// Try editing this code and see the result

func greet(name string) string {
    return "Hello, " + name + "!"
}

func main() {
    message := greet("World")
    println(message)
}`,
    output: 'Hello, World!'
  },
  http: {
    name: 'HTTP Server',
    code: `package main

import (
    "fmt"
    "net/http"
)

func handler(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "Hello from Go!")
}

func main() {
    http.HandleFunc("/", handler)
    fmt.Println("Server on :8080")
    http.ListenAndServe(":8080", nil)
}`,
    output: 'Server on :8080'
  },
  agent: {
    name: 'Agent Loop',
    code: `// Simulating claude-code-Go agent loop

type Agent struct {
    Context []Message
}

func (a *Agent) Process(input string) string {
    // 1. Think: Analyze the request
    // 2. Act: Select and execute tools
    // 3. Observe: Process results
    return "Processing: " + input
}

func main() {
    agent := &Agent{}
    result := agent.Process("Create a file")
    println(result)
}`,
    output: 'Processing: Create a file'
  }
}

const lineNumbers = computed(() => {
  return code.value.split('\n').length
})

const runCode = async () => {
  isRunning.value = true
  output.value = ''
  
  // Simulate execution with typing effect
  const lines = [
    '> Compiling...',
    '> Build successful',
    '> Running...',
    ''
  ]
  
  for (const line of lines) {
    output.value += line + '\n'
    await new Promise(r => setTimeout(r, 200))
  }
  
  // Show actual output based on selected example
  const example = examples[selectedExample.value as keyof typeof examples]
  output.value += example.output
  
  isRunning.value = false
}

const copyCode = async () => {
  try {
    await navigator.clipboard.writeText(code.value)
    copied.value = true
    setTimeout(() => copied.value = false, 2000)
  } catch (err) {
    console.error('Failed to copy:', err)
  }
}

const loadExample = (key: string) => {
  selectedExample.value = key
  code.value = examples[key as keyof typeof examples].code
  output.value = ''
}

const handleKeydown = (e: KeyboardEvent) => {
  if (e.key === 'Tab') {
    e.preventDefault()
    const start = (e.target as HTMLTextAreaElement).selectionStart
    const end = (e.target as HTMLTextAreaElement).selectionEnd
    code.value = code.value.substring(0, start) + '    ' + code.value.substring(end)
    setTimeout(() => {
      (e.target as HTMLTextAreaElement).selectionStart = (e.target as HTMLTextAreaElement).selectionEnd = start + 4
    }, 0)
  }
}

onMounted(() => {
  // Auto-run on mount for demo effect
  setTimeout(() => {
    runCode()
  }, 1000)
})
</script>

<template>
  <div class="playground-container">
    <div class="playground-header">
      <div class="playground-title">
        <span class="title-icon">⚡</span>
        <h3>Interactive Playground</h3>
      </div>
      <div class="example-selector">
        <button
          v-for="(example, key) in examples"
          :key="key"
          class="example-btn"
          :class="{ active: selectedExample === key }"
          @click="loadExample(key)"
        >
          {{ example.name }}
        </button>
      </div>
    </div>
    
    <div class="playground-body">
      <div class="editor-section">
        <div class="editor-header">
          <span class="editor-label">main.go</span>
          <button class="copy-btn" @click="copyCode" :class="{ copied }">
            <svg v-if="!copied" viewBox="0 0 24 24" width="14" height="14" fill="none" stroke="currentColor" stroke-width="2">
              <rect x="9" y="9" width="13" height="13" rx="2" ry="2"></rect>
              <path d="M5 15H4a2 2 0 0 1-2-2V4a2 2 0 0 1 2-2h9a2 2 0 0 1 2 2v1"></path>
            </svg>
            <svg v-else viewBox="0 0 24 24" width="14" height="14" fill="none" stroke="currentColor" stroke-width="2">
              <polyline points="20 6 9 17 4 12"></polyline>
            </svg>
            <span>{{ copied ? 'Copied!' : 'Copy' }}</span>
          </button>
        </div>
        <div class="editor-wrapper">
          <div class="line-numbers">
            <span v-for="n in lineNumbers" :key="n" class="line-num">{{ n }}</span>
          </div>
          <textarea
            v-model="code"
            class="code-editor"
            spellcheck="false"
            @keydown="handleKeydown"
          ></textarea>
        </div>
      </div>
      
      <div class="output-section">
        <div class="output-header">
          <span class="output-label">Output</span>
          <AnimatedButton
            variant="primary"
            size="sm"
            :icon="isRunning ? undefined : 'play'"
            @click="runCode"
            :disabled="isRunning"
          >
            <span v-if="isRunning" class="run-spinner"></span>
            <span v-else>Run</span>
          </AnimatedButton>
        </div>
        <div class="output-content">
          <pre v-if="output">{{ output }}</pre>
          <div v-else class="output-placeholder">
            Click "Run" to see the output
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.playground-container {
  background: linear-gradient(145deg, var(--vp-c-bg-soft) 0%, var(--vp-c-bg) 100%);
  border: 1px solid var(--vp-c-divider);
  border-radius: 16px;
  overflow: hidden;
  box-shadow: 0 20px 60px rgba(0, 0, 0, 0.1);
}

.playground-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 16px 20px;
  background: linear-gradient(180deg, var(--vp-c-bg) 0%, var(--vp-c-bg-soft) 100%);
  border-bottom: 1px solid var(--vp-c-divider);
  flex-wrap: wrap;
  gap: 12px;
}

.playground-title {
  display: flex;
  align-items: center;
  gap: 10px;
}

.title-icon {
  font-size: 24px;
}

.playground-title h3 {
  margin: 0;
  font-size: 18px;
  font-weight: 600;
  color: var(--vp-c-text-1);
}

.example-selector {
  display: flex;
  gap: 8px;
  flex-wrap: wrap;
}

.example-btn {
  padding: 6px 14px;
  border: 1px solid var(--vp-c-divider);
  background: var(--vp-c-bg);
  border-radius: 6px;
  font-size: 13px;
  color: var(--vp-c-text-2);
  cursor: pointer;
  transition: all 0.2s ease;
}

.example-btn:hover {
  border-color: var(--vp-c-brand);
  color: var(--vp-c-brand);
}

.example-btn.active {
  background: var(--vp-c-brand);
  border-color: var(--vp-c-brand);
  color: white;
}

.playground-body {
  display: grid;
  grid-template-columns: 1.2fr 1fr;
  min-height: 400px;
}

@media (max-width: 768px) {
  .playground-body {
    grid-template-columns: 1fr;
    grid-template-rows: 1fr auto;
  }
}

.editor-section {
  border-right: 1px solid var(--vp-c-divider);
  display: flex;
  flex-direction: column;
}

@media (max-width: 768px) {
  .editor-section {
    border-right: none;
    border-bottom: 1px solid var(--vp-c-divider);
  }
}

.editor-header,
.output-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 12px 16px;
  background: var(--vp-c-bg);
  border-bottom: 1px solid var(--vp-c-divider);
}

.editor-label,
.output-label {
  font-size: 12px;
  font-weight: 600;
  text-transform: uppercase;
  letter-spacing: 0.5px;
  color: var(--vp-c-text-2);
  font-family: 'JetBrains Mono', monospace;
}

.copy-btn {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 6px 12px;
  border: 1px solid var(--vp-c-divider);
  background: var(--vp-c-bg);
  border-radius: 6px;
  font-size: 12px;
  color: var(--vp-c-text-2);
  cursor: pointer;
  transition: all 0.2s ease;
}

.copy-btn:hover {
  border-color: var(--vp-c-brand);
  color: var(--vp-c-brand);
}

.copy-btn.copied {
  background: #22c55e;
  border-color: #22c55e;
  color: white;
}

.editor-wrapper {
  display: flex;
  flex: 1;
  overflow: hidden;
}

.line-numbers {
  padding: 16px 12px;
  background: var(--vp-c-bg-soft);
  border-right: 1px solid var(--vp-c-divider);
  font-family: 'JetBrains Mono', 'Fira Code', monospace;
  font-size: 13px;
  line-height: 1.7;
  color: var(--vp-c-text-3);
  user-select: none;
  text-align: right;
  display: flex;
  flex-direction: column;
  min-width: 48px;
}

.line-num {
  display: block;
}

.code-editor {
  flex: 1;
  padding: 16px;
  border: none;
  background: var(--vp-c-bg);
  font-family: 'JetBrains Mono', 'Fira Code', monospace;
  font-size: 13px;
  line-height: 1.7;
  color: var(--vp-c-text-1);
  resize: none;
  outline: none;
  tab-size: 4;
}

.code-editor::placeholder {
  color: var(--vp-c-text-3);
}

.output-section {
  display: flex;
  flex-direction: column;
  background: #0d1117;
}

.output-header {
  background: #161b22;
  border-bottom-color: #30363d;
}

.output-label {
  color: #8b949e;
}

.output-content {
  flex: 1;
  padding: 16px;
  overflow: auto;
  font-family: 'JetBrains Mono', 'Fira Code', monospace;
  font-size: 13px;
  line-height: 1.7;
}

.output-content pre {
  margin: 0;
  color: #e6edf3;
  white-space: pre-wrap;
}

.output-placeholder {
  color: #8b949e;
  font-style: italic;
}

.run-spinner {
  display: inline-block;
  width: 14px;
  height: 14px;
  border: 2px solid rgba(255, 255, 255, 0.3);
  border-top-color: white;
  border-radius: 50%;
  animation: spin 1s linear infinite;
}

@keyframes spin {
  to { transform: rotate(360deg); }
}
</style>