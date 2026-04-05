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
      <a
        href="https://github.com/strings77wzq/claude-code-Go"
        target="_blank"
        class="nav-github-link"
        aria-label="GitHub"
      >
        <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="currentColor">
          <path d="M12 0c-6.626 0-12 5.373-12 12 0 5.302 3.438 9.8 8.207 11.387.599.111.793-.261.793-.577v-2.234c-3.338.726-4.033-1.416-4.033-1.416-.546-1.387-1.333-1.756-1.333-1.756-1.089-.745.083-.729.083-.729 1.205.084 1.839 1.237 1.839 1.237 1.07 1.834 2.807 1.304 3.492.997.107-.775.418-1.305.762-1.604-2.665-.305-5.467-1.334-5.467-5.931 0-1.311.469-2.381 1.236-3.221-.124-.303-.535-1.524.117-3.176 0 0 1.008-.322 3.301 1.23.957-.266 1.983-.399 3.003-.404 1.02.005 2.047.138 3.006.404 2.291-1.552 3.297-1.23 3.297-1.23.653 1.653.242 2.874.118 3.176.77.84 1.235 1.911 1.235 3.221 0 4.609-2.807 5.624-5.479 5.921.43.372.823 1.102.823 2.222v3.293c0 .319.192.694.801.576 4.765-1.589 8.199-6.086 8.199-11.386 0-6.627-5.373-12-12-12z"/>
        </svg>
      </a>
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
.nav-github-link {
  display: flex;
  align-items: center;
  color: var(--vp-c-text-2);
  transition: color 0.2s ease;
}

.nav-github-link:hover {
  color: var(--vp-c-brand-1);
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
