<script setup lang="ts">
import { ref, onMounted } from 'vue'

interface Testimonial {
  id: number
  name: string
  role: string
  company: string
  avatar: string
  content: string
  rating: number
}

const testimonials = ref<Testimonial[]>([
  {
    id: 1,
    name: 'Alex Chen',
    role: 'Senior Backend Engineer',
    company: 'TechCorp',
    avatar: '👨‍💻',
    content: 'claude-code-Go has completely transformed our workflow. The single binary deployment makes it incredibly easy to distribute across our team. The permission system gives us confidence when running automated tasks.',
    rating: 5
  },
  {
    id: 2,
    name: 'Sarah Miller',
    role: 'DevOps Lead',
    company: 'CloudScale',
    avatar: '👩‍💼',
    content: 'The agent loop architecture is brilliant. We have integrated it into our CI/CD pipeline and it handles complex deployment tasks autonomously. The SSE streaming provides real-time feedback.',
    rating: 5
  },
  {
    id: 3,
    name: 'David Park',
    role: 'Open Source Contributor',
    company: 'GitHub',
    avatar: '🧑‍🔧',
    content: 'As someone who contributes to many Go projects, having an AI assistant that understands Go idioms is invaluable. The MCP integration means I can extend it with my own tools easily.',
    rating: 5
  },
  {
    id: 4,
    name: 'Emily Zhang',
    role: 'Software Architect',
    company: 'StartupXYZ',
    avatar: '👩‍🎨',
    content: 'The harness-first approach really resonates with me. Safety and reliability are built in, not afterthoughts. This is how AI coding tools should be designed.',
    rating: 5
  }
])

const currentIndex = ref(0)
const isAnimating = ref(false)

const nextTestimonial = () => {
  if (isAnimating.value) return
  isAnimating.value = true
  setTimeout(() => {
    currentIndex.value = (currentIndex.value + 1) % testimonials.value.length
    isAnimating.value = false
  }, 300)
}

const prevTestimonial = () => {
  if (isAnimating.value) return
  isAnimating.value = true
  setTimeout(() => {
    currentIndex.value = (currentIndex.value - 1 + testimonials.value.length) % testimonials.value.length
    isAnimating.value = false
  }, 300)
}

const goToTestimonial = (index: number) => {
  if (isAnimating.value || index === currentIndex.value) return
  isAnimating.value = true
  setTimeout(() => {
    currentIndex.value = index
    isAnimating.value = false
  }, 300)
}

// Auto-advance
let intervalId: number
onMounted(() => {
  intervalId = window.setInterval(nextTestimonial, 6000)
})
</script>

<template>
  <div class="testimonials-section">
    <div class="testimonials-container">
      <div class="testimonial-card" :class="{ 'is-animating': isAnimating }">
        <div class="quote-icon">❝</div>
        <p class="testimonial-content">{{ testimonials[currentIndex].content }}</p>
        <div class="testimonial-author">
          <span class="author-avatar">{{ testimonials[currentIndex].avatar }}</span>
          <div class="author-info">
            <span class="author-name">{{ testimonials[currentIndex].name }}</span>
            <span class="author-role">{{ testimonials[currentIndex].role }} @ {{ testimonials[currentIndex].company }}</span>
          </div>
          <div class="rating">
            <span v-for="n in testimonials[currentIndex].rating" :key="n" class="star">⭐</span>
          </div>
        </div>
      </div>
      
      <div class="testimonial-nav">
        <button class="nav-btn" @click="prevTestimonial" aria-label="Previous testimonial">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <polyline points="15 18 9 12 15 6"></polyline>
          </svg>
        </button>
        <div class="nav-dots">
          <button
            v-for="(_, index) in testimonials"
            :key="index"
            class="nav-dot"
            :class="{ active: currentIndex === index }"
            @click="goToTestimonial(index)"
            :aria-label="`Go to testimonial ${index + 1}`"
          />
        </div>
        <button class="nav-btn" @click="nextTestimonial" aria-label="Next testimonial">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <polyline points="9 18 15 12 9 6"></polyline>
          </svg>
        </button>
      </div>
    </div>
  </div>
</template>

<style scoped>
.testimonials-section {
  width: 100%;
  max-width: 800px;
  margin: 0 auto;
}

.testimonials-container {
  position: relative;
}

.testimonial-card {
  background: linear-gradient(145deg, var(--vp-c-bg-soft) 0%, var(--vp-c-bg) 100%);
  border: 1px solid var(--vp-c-divider);
  border-radius: 20px;
  padding: 40px;
  position: relative;
  transition: all 0.3s ease;
  opacity: 1;
  transform: translateY(0);
}

.testimonial-card.is-animating {
  opacity: 0;
  transform: translateY(10px);
}

.quote-icon {
  position: absolute;
  top: 20px;
  left: 30px;
  font-size: 60px;
  color: var(--vp-c-brand);
  opacity: 0.2;
  font-family: Georgia, serif;
  line-height: 1;
}

.testimonial-content {
  font-size: 18px;
  line-height: 1.8;
  color: var(--vp-c-text-1);
  margin: 20px 0 30px;
  position: relative;
  z-index: 1;
}

.testimonial-author {
  display: flex;
  align-items: center;
  gap: 16px;
}

.author-avatar {
  width: 56px;
  height: 56px;
  border-radius: 50%;
  background: linear-gradient(135deg, var(--vp-c-brand) 0%, var(--vp-c-brand-light, #5c9aff) 100%);
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 28px;
  flex-shrink: 0;
}

.author-info {
  display: flex;
  flex-direction: column;
  flex: 1;
}

.author-name {
  font-weight: 600;
  color: var(--vp-c-text-1);
  font-size: 16px;
}

.author-role {
  font-size: 14px;
  color: var(--vp-c-text-2);
}

.rating {
  display: flex;
  gap: 2px;
}

.star {
  font-size: 16px;
}

.testimonial-nav {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 20px;
  margin-top: 24px;
}

.nav-btn {
  width: 40px;
  height: 40px;
  border-radius: 50%;
  border: 1px solid var(--vp-c-divider);
  background: var(--vp-c-bg);
  color: var(--vp-c-text-2);
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: all 0.2s ease;
}

.nav-btn:hover {
  border-color: var(--vp-c-brand);
  color: var(--vp-c-brand);
  transform: scale(1.1);
}

.nav-btn svg {
  width: 20px;
  height: 20px;
}

.nav-dots {
  display: flex;
  gap: 8px;
}

.nav-dot {
  width: 10px;
  height: 10px;
  border-radius: 50%;
  border: none;
  background: var(--vp-c-divider);
  cursor: pointer;
  transition: all 0.2s ease;
  padding: 0;
}

.nav-dot:hover {
  background: var(--vp-c-brand-light, #5c9aff);
}

.nav-dot.active {
  background: var(--vp-c-brand);
  transform: scale(1.2);
}

@media (max-width: 640px) {
  .testimonial-card {
    padding: 30px 20px;
  }
  
  .testimonial-content {
    font-size: 16px;
  }
  
  .quote-icon {
    font-size: 40px;
    top: 15px;
    left: 20px;
  }
  
  .author-avatar {
    width: 48px;
    height: 48px;
    font-size: 24px;
  }
}
</style>