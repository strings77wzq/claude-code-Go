<script setup lang="ts">
import { ref, onMounted, watch } from 'vue'

const starCount = ref(0)
const isLoading = ref(true)
const error = ref('')

const fetchStars = async () => {
  try {
    const response = await fetch('https://api.github.com/repos/strings77wzq/claude-code-Go')
    if (!response.ok) throw new Error('Failed to fetch')
    const data = await response.json()
    starCount.value = data.stargazers_count
  } catch (err) {
    error.value = 'Unable to load'
  } finally {
    isLoading.value = false
  }
}

const formatNumber = (num: number): string => {
  if (num >= 1000) {
    return (num / 1000).toFixed(1) + 'k'
  }
  return num.toString()
}

onMounted(() => {
  fetchStars()
  // Refresh every 5 minutes
  setInterval(fetchStars, 300000)
})
</script>

<template>
  <a 
    href="https://github.com/strings77wzq/claude-code-Go" 
    target="_blank" 
    rel="noopener noreferrer"
    class="github-stars"
    :class="{ 'is-loading': isLoading }"
  >
    <span class="stars-icon">
      <svg viewBox="0 0 24 24" fill="currentColor" width="16" height="16">
        <path d="M12 0c-6.626 0-12 5.373-12 12 0 5.302 3.438 9.8 8.207 11.387.599.111.793-.261.793-.577v-2.234c-3.338.726-4.033-1.416-4.033-1.416-.546-1.387-1.333-1.756-1.333-1.756-1.089-.745.083-.729.083-.729 1.205.084 1.839 1.237 1.839 1.237 1.07 1.834 2.807 1.304 3.492.997.107-.775.418-1.305.762-1.604-2.665-.305-5.467-1.334-5.467-5.931 0-1.311.469-2.381 1.236-3.221-.124-.303-.535-1.524.117-3.176 0 0 1.008-.322 3.301 1.23.957-.266 1.983-.399 3.003-.404 1.02.005 2.047.138 3.006.404 2.291-1.552 3.297-1.23 3.297-1.23.653 1.653.242 2.874.118 3.176.77.84 1.235 1.911 1.235 3.221 0 4.609-2.807 5.624-5.479 5.921.43.372.823 1.102.823 2.222v3.293c0 .319.192.694.801.576 4.765-1.589 8.199-6.086 8.199-11.386 0-6.627-5.373-12-12-12z"/>
      </svg>
    </span>
    <span class="stars-icon-star">⭐</span>
    <span v-if="isLoading" class="stars-count">...</span>
    <span v-else-if="error" class="stars-count">Star</span>
    <span v-else class="stars-count">{{ formatNumber(starCount) }}</span>
  </a>
</template>

<style scoped>
.github-stars {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  padding: 6px 14px;
  background: linear-gradient(135deg, #fafbfc 0%, #f6f8fa 100%);
  border: 1px solid #d0d7de;
  border-radius: 20px;
  font-size: 14px;
  font-weight: 600;
  color: #24292f;
  text-decoration: none;
  transition: all 0.2s ease;
  box-shadow: 0 1px 2px rgba(0, 0, 0, 0.05);
}

.dark .github-stars {
  background: linear-gradient(135deg, #21262d 0%, #161b22 100%);
  border-color: #30363d;
  color: #e6edf3;
}

.github-stars:hover {
  transform: translateY(-1px);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
  border-color: var(--vp-c-brand);
}

.stars-icon {
  display: flex;
  align-items: center;
  color: #57606a;
}

.dark .stars-icon {
  color: #8b949e;
}

.stars-icon-star {
  font-size: 14px;
}

.stars-count {
  min-width: 24px;
  text-align: center;
}

.is-loading {
  opacity: 0.7;
}
</style>