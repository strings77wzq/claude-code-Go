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
</style>
