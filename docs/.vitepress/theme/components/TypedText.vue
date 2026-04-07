<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue'

interface TypedTextProps {
  strings: string[]
  typeSpeed?: number
  backSpeed?: number
  backDelay?: number
  startDelay?: number
  loop?: boolean
  showCursor?: boolean
  cursorChar?: string
}

const props = withDefaults(defineProps<TypedTextProps>(), {
  typeSpeed: 50,
  backSpeed: 30,
  backDelay: 2000,
  startDelay: 500,
  loop: true,
  showCursor: true,
  cursorChar: '|'
})

const displayText = ref('')
const currentStringIndex = ref(0)
const currentCharIndex = ref(0)
const isDeleting = ref(false)
const isPaused = ref(false)
let timeoutId: number | null = null

const type = () => {
  const currentString = props.strings[currentStringIndex.value]
  
  if (isDeleting.value) {
    displayText.value = currentString.substring(0, currentCharIndex.value - 1)
    currentCharIndex.value--
  } else {
    displayText.value = currentString.substring(0, currentCharIndex.value + 1)
    currentCharIndex.value++
  }

  let typeSpeed = props.typeSpeed

  if (isDeleting.value) {
    typeSpeed = props.backSpeed
  }

  if (!isDeleting.value && currentCharIndex.value === currentString.length) {
    if (!props.loop && currentStringIndex.value === props.strings.length - 1) {
      return
    }
    typeSpeed = props.backDelay
    isDeleting.value = true
  } else if (isDeleting.value && currentCharIndex.value === 0) {
    isDeleting.value = false
    currentStringIndex.value = (currentStringIndex.value + 1) % props.strings.length
    typeSpeed = props.startDelay
  }

  timeoutId = window.setTimeout(type, typeSpeed)
}

onMounted(() => {
  timeoutId = window.setTimeout(type, props.startDelay)
})

onUnmounted(() => {
  if (timeoutId) {
    clearTimeout(timeoutId)
  }
})
</script>

<template>
  <span class="typed-text">
    {{ displayText }}
    <span v-if="showCursor" class="typed-cursor" :class="{ 'is-typing': !isDeleting && !isPaused }">{{ cursorChar }}</span>
  </span>
</template>

<style scoped>
.typed-text {
  display: inline;
}

.typed-cursor {
  display: inline-block;
  color: var(--vp-c-brand);
  animation: blink 1s step-end infinite;
}

.typed-cursor.is-typing {
  animation: none;
  opacity: 1;
}

@keyframes blink {
  0%, 100% { opacity: 1; }
  50% { opacity: 0; }
}
</style>