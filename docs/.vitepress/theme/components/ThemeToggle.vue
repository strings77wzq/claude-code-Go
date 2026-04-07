<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue'

const isDark = ref(false)
const isAnimating = ref(false)

const toggleTheme = () => {
  isAnimating.value = true
  isDark.value = !isDark.value
  
  // Apply theme
  if (isDark.value) {
    document.documentElement.classList.add('dark')
    localStorage.setItem('vitepress-theme-appearance', 'dark')
  } else {
    document.documentElement.classList.remove('dark')
    localStorage.setItem('vitepress-theme-appearance', 'light')
  }
  
  setTimeout(() => {
    isAnimating.value = false
  }, 300)
}

const systemPrefersDark = () => {
  return window.matchMedia && window.matchMedia('(prefers-color-scheme: dark)').matches
}

onMounted(() => {
  // Check saved preference or system preference
  const saved = localStorage.getItem('vitepress-theme-appearance')
  if (saved) {
    isDark.value = saved === 'dark'
  } else {
    isDark.value = systemPrefersDark()
  }
  
  // Apply initial theme
  if (isDark.value) {
    document.documentElement.classList.add('dark')
  }
  
  // Listen for system preference changes
  const mediaQuery = window.matchMedia('(prefers-color-scheme: dark)')
  const handleChange = (e: MediaQueryListEvent) => {
    if (!localStorage.getItem('vitepress-theme-appearance')) {
      isDark.value = e.matches
      if (e.matches) {
        document.documentElement.classList.add('dark')
      } else {
        document.documentElement.classList.remove('dark')
      }
    }
  }
  
  mediaQuery.addEventListener('change', handleChange)
  
  onUnmounted(() => {
    mediaQuery.removeEventListener('change', handleChange)
  })
})
</script>

<template>
  <button
    class="theme-toggle"
    @click="toggleTheme"
    :class="{ 'is-animating': isAnimating, 'is-dark': isDark }"
    :title="isDark ? 'Switch to light mode' : 'Switch to dark mode'"
  >
    <span class="toggle-track">
      <span class="toggle-thumb">
        <span class="icon-sun">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <circle cx="12" cy="12" r="5"/>
            <line x1="12" y1="1" x2="12" y2="3"/>
            <line x1="12" y1="21" x2="12" y2="23"/>
            <line x1="4.22" y1="4.22" x2="5.64" y2="5.64"/>
            <line x1="18.36" y1="18.36" x2="19.78" y2="19.78"/>
            <line x1="1" y1="12" x2="3" y2="12"/>
            <line x1="21" y1="12" x2="23" y2="12"/>
            <line x1="4.22" y1="19.78" x2="5.64" y2="18.36"/>
            <line x1="18.36" y1="5.64" x2="19.78" y2="4.22"/>
          </svg>
        </span>
        <span class="icon-moon">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <path d="M21 12.79A9 9 0 1 1 11.21 3 7 7 0 0 0 21 12.79z"/>
          </svg>
        </span>
      </span>
    </span>
  </button>
</template>

<style scoped>
.theme-toggle {
  background: none;
  border: none;
  padding: 4px;
  cursor: pointer;
  position: relative;
  display: flex;
  align-items: center;
}

.toggle-track {
  width: 48px;
  height: 26px;
  background: #e2e8f0;
  border-radius: 13px;
  position: relative;
  transition: background 0.3s ease;
  box-shadow: inset 0 1px 3px rgba(0, 0, 0, 0.1);
}

.dark .toggle-track {
  background: #334155;
}

.toggle-thumb {
  position: absolute;
  top: 2px;
  left: 2px;
  width: 22px;
  height: 22px;
  background: white;
  border-radius: 50%;
  transition: transform 0.3s cubic-bezier(0.16, 1, 0.3, 1);
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.2);
  display: flex;
  align-items: center;
  justify-content: center;
}

.is-dark .toggle-thumb {
  transform: translateX(22px);
}

.icon-sun,
.icon-moon {
  position: absolute;
  width: 14px;
  height: 14px;
  transition: all 0.3s ease;
}

.icon-sun {
  color: #f59e0b;
  opacity: 1;
  transform: scale(1);
}

.icon-moon {
  color: #94a3b8;
  opacity: 0;
  transform: scale(0.5);
}

.is-dark .icon-sun {
  opacity: 0;
  transform: scale(0.5);
}

.is-dark .icon-moon {
  opacity: 1;
  transform: scale(1);
  color: #cbd5e1;
}

.is-animating .toggle-thumb {
  animation: bounce 0.3s ease;
}

@keyframes bounce {
  0%, 100% { transform: translateX(0) scale(1); }
  50% { transform: translateX(11px) scale(0.9); }
}

.is-animating.is-dark .toggle-thumb {
  animation: bounceDark 0.3s ease;
}

@keyframes bounceDark {
  0%, 100% { transform: translateX(22px) scale(1); }
  50% { transform: translateX(11px) scale(0.9); }
}
</style>