<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue'

const copied = ref(false)
const copyTimeout = ref<ReturnType<typeof setTimeout> | null>(null)

const copyCode = async (event: MouseEvent) => {
  const button = event.target as HTMLElement
  const pre = button.closest('div[class*="language-"]')?.querySelector('code') 
    || button.parentElement?.querySelector('code')
    || button.closest('pre')?.querySelector('code')
  
  if (!pre) return
  
  const code = pre.textContent || ''
  
  try {
    await navigator.clipboard.writeText(code)
    copied.value = true
    
    if (copyTimeout.value) {
      clearTimeout(copyTimeout.value)
    }
    
    copyTimeout.value = setTimeout(() => {
      copied.value = false
    }, 2000)
  } catch (err) {
    console.error('Failed to copy:', err)
  }
}

onMounted(() => {
  document.addEventListener('click', copyCode)
})

onUnmounted(() => {
  document.removeEventListener('click', copyCode)
  if (copyTimeout.value) {
    clearTimeout(copyTimeout.value)
  }
})
</script>

<template>
  <Teleport to="body">
    <div class="copy-button-container" style="display: none;">
      <!-- This component injects copy buttons into code blocks via CSS/JS -->
    </div>
  </Teleport>
</template>

<style>
/* Copy button styles for all code blocks */
.vp-doc .vp-code-group .tabs,
.vp-doc div[class*="language-"] {
  position: relative;
}

.vp-doc div[class*="language-"]:hover .copy-button,
.vp-doc .vp-code-group:hover .copy-button {
  opacity: 1;
}

.copy-button {
  position: absolute;
  top: 8px;
  right: 8px;
  z-index: 10;
  display: flex;
  align-items: center;
  justify-content: center;
  width: 32px;
  height: 32px;
  padding: 0;
  background: var(--vp-c-bg-soft, #161b22);
  border: 1px solid var(--vp-c-divider, #30363d);
  border-radius: 6px;
  color: var(--vp-c-text-2, #8b949e);
  cursor: pointer;
  opacity: 0;
  transition: all 0.2s ease;
}

.copy-button:hover {
  background: var(--vp-c-bg, #0d1117);
  color: var(--vp-c-text-1, #e6edf3);
  border-color: var(--vp-c-brand-1, #00add8);
}

.copy-button svg {
  width: 16px;
  height: 16px;
}

.copy-button.copied {
  color: #7ee787;
  border-color: #7ee787;
}

/* Ensure code groups have proper positioning */
.vp-doc .vp-code-group {
  position: relative;
}
</style>