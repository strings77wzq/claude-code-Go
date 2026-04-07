<script setup lang="ts">
import { useData, useRoute } from 'vitepress'
import DefaultTheme from 'vitepress/theme'
import { computed } from 'vue'

const { site, frontmatter } = useData()
const route = useRoute()

const isHome = computed(() => route.path === '/' || route.path === '/zh/')
</script>

<template>
  <DefaultTheme.Layout>
    <template #nav-bar-content-after>
      <div class="nav-actions">
        <ThemeToggle />
        <GitHubStars />
      </div>
    </template>

    <template #footer-before>
      <div class="custom-footer">
        <div class="custom-footer-content">
          <div class="custom-footer-links">
            <a href="https://github.com/strings77wzq/claude-code-Go" class="custom-footer-link">GitHub</a>
            <a href="https://github.com/strings77wzq/claude-code-Go/releases" class="custom-footer-link">Releases</a>
            <a href="/guide/introduction" class="custom-footer-link">Documentation</a>
            <a href="/feedback" class="custom-footer-link">Feedback</a>
          </div>
          <div class="custom-footer-badges">
            <span class="custom-footer-badge">MIT License</span>
            <span class="custom-footer-badge">Go 1.24+</span>
          </div>
          <div class="custom-footer-copyright">
            Copyright © 2026-present claude-code-Go Contributors
          </div>
        </div>
      </div>
    </template>
  </DefaultTheme.Layout>
</template>

<style scoped>
.nav-actions {
  display: flex;
  align-items: center;
  gap: 12px;
  margin-left: 12px;
  padding-left: 12px;
  border-left: 1px solid var(--vp-c-divider);
}

@media (max-width: 768px) {
  .nav-actions {
    gap: 8px;
    margin-left: 8px;
    padding-left: 8px;
  }
}

.custom-footer {
  background: var(--vp-c-bg-soft);
  border-top: 1px solid var(--vp-c-divider);
  padding: 2rem 1.5rem;
  margin-top: auto;
}

.custom-footer-content {
  max-width: 960px;
  margin: 0 auto;
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 1.5rem;
}

.custom-footer-links {
  display: flex;
  gap: 1.5rem;
  flex-wrap: wrap;
  justify-content: center;
}

.custom-footer-link {
  color: var(--vp-c-text-2);
  text-decoration: none;
  font-size: 0.875rem;
  transition: color 0.2s ease;
}

.custom-footer-link:hover {
  color: var(--vp-c-brand-1);
}

.custom-footer-badges {
  display: flex;
  gap: 0.75rem;
}

.custom-footer-badge {
  display: inline-flex;
  align-items: center;
  gap: 0.375rem;
  padding: 0.375rem 0.75rem;
  background: var(--vp-c-brand-soft);
  color: var(--vp-c-brand-1);
  border-radius: 9999px;
  font-size: 0.75rem;
  font-weight: 500;
}

.custom-footer-copyright {
  width: 100%;
  text-align: center;
  padding-top: 1.5rem;
  border-top: 1px solid var(--vp-c-divider);
  color: var(--vp-c-text-2);
  font-size: 0.8125rem;
}

@media (max-width: 768px) {
  .custom-footer {
    padding: 1.5rem 1rem;
  }

  .custom-footer-links {
    gap: 1rem;
  }

  .custom-footer-link {
    font-size: 0.8125rem;
  }
}
</style>
