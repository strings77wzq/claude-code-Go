<script setup lang="ts">
import { useData } from 'vitepress'
import DefaultTheme from 'vitepress/theme'
import './custom.css'

const { page } = useData()

// Extract custom sections from frontmatter
const useCases = page.value.frontmatter.useCases as Array<{
  icon: string
  title: string
  details: string
}> | undefined

const cta = page.value.frontmatter.cta as Array<{
  title: string
  details: string
  actions: Array<{ text: string; link: string }>
}> | undefined

const stats = page.value.frontmatter.stats as Array<{
  label: string
  value: string
}> | undefined

const learningOutcomes = page.value.frontmatter.learningOutcomes as Array<{
  title: string
  description: string
}> | undefined

// Check if this is the home page
const isHome = page.value.frontmatter.layout === 'home'
</script>

<template>
  <DefaultTheme.Layout>
    <!-- Custom sections after the default home content -->
    <template #doc-after>
      <!-- Stats Section -->
      <div v-if="isHome && stats && stats.length > 0" class="stats-section">
        <h2 class="section-title">By The Numbers</h2>
        <div class="stats-grid">
          <div v-for="(stat, index) in stats" :key="index" class="stat-card">
            <div class="stat-value">{{ stat.value }}</div>
            <div class="stat-label">{{ stat.label }}</div>
          </div>
        </div>
      </div>

      <!-- Learning Outcomes Section -->
      <div v-if="isHome && learningOutcomes && learningOutcomes.length > 0" class="learning-section">
        <h2 class="section-title">What You'll Learn</h2>
        <div class="learning-grid">
          <div v-for="(item, index) in learningOutcomes" :key="index" class="learning-card">
            <h3 class="learning-title">{{ item.title }}</h3>
            <p class="learning-description">{{ item.description }}</p>
          </div>
        </div>
      </div>

      <div v-if="isHome && useCases && useCases.length > 0" class="use-cases-section">
        <h2 class="section-title">Use Cases</h2>
        <div class="use-cases-grid">
          <div v-for="(item, index) in useCases" :key="index" class="use-case-card">
            <div class="icon">{{ item.icon }}</div>
            <h3 class="title">{{ item.title }}</h3>
            <p class="details">{{ item.details }}</p>
          </div>
        </div>
      </div>

      <div v-if="isHome && cta && cta.length > 0" class="cta-section">
        <div v-for="(item, index) in cta" :key="index" class="cta-card">
          <h2 class="cta-title">{{ item.title }}</h2>
          <p class="cta-details">{{ item.details }}</p>
          <div class="cta-actions">
            <a
              v-for="(action, idx) in item.actions"
              :key="idx"
              :href="action.link"
              class="cta-button"
              :class="{ primary: idx === 0, secondary: idx !== 0 }"
            >
              {{ action.text }}
            </a>
          </div>
        </div>
      </div>
    </template>
  </DefaultTheme.Layout>
</template>

<style scoped>
.section-title {
  font-size: 2rem;
  font-weight: 700;
  text-align: center;
  margin-bottom: 2rem;
  color: var(--vp-c-text-1);
}

/* Stats Section */
.stats-section {
  margin-top: 3rem;
  padding-top: 3rem;
  border-top: 1px solid var(--vp-c-divider);
}

.stats-grid {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 1.5rem;
}

@media (max-width: 1024px) {
  .stats-grid {
    grid-template-columns: repeat(2, 1fr);
  }
}

@media (max-width: 640px) {
  .stats-grid {
    grid-template-columns: 1fr;
  }
}

.stat-card {
  background: linear-gradient(135deg, var(--vp-c-brand-1) 0%, var(--vp-c-brand-2) 100%);
  border-radius: 12px;
  padding: 2rem 1.5rem;
  text-align: center;
  color: white;
  transition: transform 0.2s ease;
}

.stat-card:hover {
  transform: translateY(-4px);
}

.stat-value {
  font-size: 2.5rem;
  font-weight: 700;
  margin-bottom: 0.5rem;
}

.stat-label {
  font-size: 0.875rem;
  opacity: 0.9;
}

/* Learning Outcomes Section */
.learning-section {
  margin-top: 3rem;
  padding-top: 3rem;
  border-top: 1px solid var(--vp-c-divider);
}

.learning-grid {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 1.5rem;
}

@media (max-width: 1024px) {
  .learning-grid {
    grid-template-columns: repeat(2, 1fr);
  }
}

@media (max-width: 640px) {
  .learning-grid {
    grid-template-columns: 1fr;
  }
}

.learning-card {
  background: var(--vp-c-bg-soft);
  border-radius: 12px;
  padding: 1.5rem;
  border: 1px solid var(--vp-c-divider);
  transition: transform 0.2s ease, box-shadow 0.2s ease;
}

.learning-card:hover {
  transform: translateY(-4px);
  box-shadow: 0 8px 24px rgba(0, 0, 0, 0.1);
}

.learning-title {
  font-size: 1.125rem;
  font-weight: 600;
  color: var(--vp-c-brand-1);
  margin: 0 0 0.75rem 0;
}

.learning-description {
  font-size: 0.875rem;
  color: var(--vp-c-text-2);
  line-height: 1.6;
  margin: 0;
}

/* Use Cases Section */
.use-cases-section {
  margin-top: 3rem;
  padding-top: 3rem;
  border-top: 1px solid var(--vp-c-divider);
}

.use-cases-grid {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 1.5rem;
}

@media (max-width: 1024px) {
  .use-cases-grid {
    grid-template-columns: repeat(2, 1fr);
  }
}

@media (max-width: 640px) {
  .use-cases-grid {
    grid-template-columns: 1fr;
  }
}

.use-case-card {
  background: var(--vp-c-bg-soft);
  border-radius: 12px;
  padding: 1.5rem;
  border: 1px solid var(--vp-c-divider);
  transition: transform 0.2s ease, box-shadow 0.2s ease;
}

.use-case-card:hover {
  transform: translateY(-4px);
  box-shadow: 0 8px 24px rgba(0, 0, 0, 0.1);
}

.use-case-card .icon {
  font-size: 2rem;
  margin-bottom: 0.75rem;
}

.use-case-card .title {
  font-size: 1.25rem;
  font-weight: 600;
  color: var(--vp-c-text-1);
  margin: 0 0 0.5rem 0;
}

.use-case-card .details {
  font-size: 0.875rem;
  color: var(--vp-c-text-2);
  line-height: 1.6;
  margin: 0;
}

/* CTA Section */
.cta-section {
  margin-top: 3rem;
}

.cta-card {
  background: linear-gradient(135deg, var(--vp-c-brand-1) 0%, var(--vp-c-brand-2) 100%);
  border-radius: 16px;
  padding: 3rem;
  text-align: center;
  color: white;
}

.cta-title {
  font-size: 2rem;
  font-weight: 700;
  margin: 0 0 1rem 0;
}

.cta-details {
  font-size: 1.125rem;
  opacity: 0.9;
  margin: 0 0 2rem 0;
}

.cta-actions {
  display: flex;
  gap: 1rem;
  justify-content: center;
}

.cta-button {
  display: inline-block;
  padding: 0.75rem 1.5rem;
  border-radius: 8px;
  font-weight: 600;
  text-decoration: none;
  transition: all 0.2s ease;
}

.cta-button.primary {
  background: white;
  color: var(--vp-c-brand-1);
}

.cta-button.primary:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
}

.cta-button.secondary {
  background: transparent;
  color: white;
  border: 2px solid white;
}

.cta-button.secondary:hover {
  background: rgba(255, 255, 255, 0.1);
}

@media (max-width: 640px) {
  .cta-card {
    padding: 2rem 1.5rem;
  }
  
  .cta-actions {
    flex-direction: column;
    align-items: center;
  }
}
</style>