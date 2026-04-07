<script setup lang="ts">
import { ref, computed } from 'vue'

const props = defineProps<{
  variant?: 'primary' | 'secondary' | 'ghost'
  size?: 'sm' | 'md' | 'lg'
  href?: string
  to?: string
  external?: boolean
  icon?: string
}>()

const emit = defineEmits<{
  click: [event: MouseEvent]
}>()

const isHovered = ref(false)
const isPressed = ref(false)

const buttonClasses = computed(() => [
  'animated-button',
  `variant-${props.variant || 'primary'}`,
  `size-${props.size || 'md'}`,
  {
    'is-hovered': isHovered.value,
    'is-pressed': isPressed.value,
    'has-icon': !!props.icon
  }
])

const handleClick = (e: MouseEvent) => {
  emit('click', e)
}

const iconSvg = computed(() => {
  switch (props.icon) {
    case 'arrow-right':
      return `<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M5 12h14M12 5l7 7-7 7"/></svg>`
    case 'github':
      return `<svg viewBox="0 0 24 24" fill="currentColor"><path d="M12 0c-6.626 0-12 5.373-12 12 0 5.302 3.438 9.8 8.207 11.387.599.111.793-.261.793-.577v-2.234c-3.338.726-4.033-1.416-4.033-1.416-.546-1.387-1.333-1.756-1.333-1.756-1.089-.745.083-.729.083-.729 1.205.084 1.839 1.237 1.839 1.237 1.07 1.834 2.807 1.304 3.492.997.107-.775.418-1.305.762-1.604-2.665-.305-5.467-1.334-5.467-5.931 0-1.311.469-2.381 1.236-3.221-.124-.303-.535-1.524.117-3.176 0 0 1.008-.322 3.301 1.23.957-.266 1.983-.399 3.003-.404 1.02.005 2.047.138 3.006.404 2.291-1.552 3.297-1.23 3.297-1.23.653 1.653.242 2.874.118 3.176.77.84 1.235 1.911 1.235 3.221 0 4.609-2.807 5.624-5.479 5.921.43.372.823 1.102.823 2.222v3.293c0 .319.192.694.801.576 4.765-1.589 8.199-6.086 8.199-11.386 0-6.627-5.373-12-12-12z"/></svg>`
    case 'play':
      return `<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><polygon points="5 3 19 12 5 21 5 3"/></svg>`
    case 'download':
      return `<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M21 15v4a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2v-4M7 10l5 5 5-5M12 15V3"/></svg>`
    default:
      return ''
  }
})
</script>

<template>
  <component
    :is="href ? 'a' : 'button'"
    :class="buttonClasses"
    :href="href"
    :target="external ? '_blank' : undefined"
    :rel="external ? 'noopener noreferrer' : undefined"
    @mouseenter="isHovered = true"
    @mouseleave="isHovered = false, isPressed = false"
    @mousedown="isPressed = true"
    @mouseup="isPressed = false"
    @click="handleClick"
  >
    <span class="button-bg"></span>
    <span class="button-content">
      <span v-if="icon" class="button-icon" v-html="iconSvg"></span>
      <span class="button-text"><slot /></span>
    </span>
    <span class="button-shine"></span>
  </component>
</template>

<style scoped>
.animated-button {
  position: relative;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  border: none;
  border-radius: 10px;
  font-weight: 600;
  text-decoration: none;
  cursor: pointer;
  overflow: hidden;
  transition: all 0.3s cubic-bezier(0.16, 1, 0.3, 1);
  transform: translateY(0);
}

.button-bg {
  position: absolute;
  inset: 0;
  border-radius: inherit;
  transition: all 0.3s ease;
}

.button-content {
  position: relative;
  z-index: 1;
  display: flex;
  align-items: center;
  gap: 8px;
}

.button-icon {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 18px;
  height: 18px;
  transition: transform 0.3s ease;
}

.button-icon :deep(svg) {
  width: 100%;
  height: 100%;
}

.button-text {
  position: relative;
}

.button-shine {
  position: absolute;
  top: 0;
  left: -100%;
  width: 100%;
  height: 100%;
  background: linear-gradient(
    90deg,
    transparent,
    rgba(255, 255, 255, 0.2),
    transparent
  );
  transition: left 0.5s ease;
}

/* Primary Variant */
.variant-primary .button-bg {
  background: linear-gradient(135deg, var(--vp-c-brand) 0%, var(--vp-c-brand-light, #5c9aff) 100%);
  box-shadow: 
    0 4px 14px rgba(59, 130, 246, 0.4),
    0 1px 2px rgba(0, 0, 0, 0.1);
}

.variant-primary {
  color: white;
}

.variant-primary:hover .button-bg {
  box-shadow: 
    0 6px 20px rgba(59, 130, 246, 0.5),
    0 2px 4px rgba(0, 0, 0, 0.1);
}

/* Secondary Variant */
.variant-secondary .button-bg {
  background: var(--vp-c-bg-soft);
  border: 1px solid var(--vp-c-divider);
}

.variant-secondary {
  color: var(--vp-c-text-1);
}

.variant-secondary:hover .button-bg {
  background: var(--vp-c-bg-mute);
  border-color: var(--vp-c-brand);
}

/* Ghost Variant */
.variant-ghost .button-bg {
  background: transparent;
}

.variant-ghost {
  color: var(--vp-c-text-2);
}

.variant-ghost:hover .button-bg {
  background: var(--vp-c-bg-soft);
}

/* Sizes */
.size-sm {
  padding: 8px 16px;
  font-size: 14px;
}

.size-md {
  padding: 12px 24px;
  font-size: 15px;
}

.size-lg {
  padding: 16px 32px;
  font-size: 16px;
}

.size-lg .button-icon {
  width: 20px;
  height: 20px;
}

/* Hover Effects */
.animated-button.is-hovered {
  transform: translateY(-2px);
}

.animated-button.is-hovered .button-shine {
  left: 100%;
}

.animated-button.is-hovered.has-icon .button-icon {
  transform: translateX(2px);
}

.animated-button.is-pressed {
  transform: translateY(0) scale(0.98);
}

/* Focus States */
.animated-button:focus-visible {
  outline: 2px solid var(--vp-c-brand);
  outline-offset: 2px;
}

/* Reduced Motion */
@media (prefers-reduced-motion: reduce) {
  .animated-button,
  .button-bg,
  .button-icon,
  .button-shine {
    transition: none;
  }
  
  .animated-button.is-hovered {
    transform: none;
  }
  
  .animated-button.is-hovered .button-shine {
    left: -100%;
  }
}
</style>