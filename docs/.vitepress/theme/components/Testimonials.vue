<!--
  Testimonials component — hidden from production until real user stories exist.
  Re-enable by importing and registering Testimonials in ../index.ts.
-->
<template>
  <div class="testimonials-section">
    <div class="testimonials-placeholder">
      <p class="placeholder-text">
        Real user testimonials will appear here once available.
        <a href="https://github.com/strings77wzq/claude-code-Go/discussions/categories/showcase">Share your story</a>.
      </p>
    </div>
  </div>
</template>

<style scoped>
.testimonials-section {
  width: 100%;
  max-width: 800px;
  margin: 0 auto;
}

.testimonials-placeholder {
  background: var(--vp-c-bg-soft);
  border: 1px dashed var(--vp-c-divider);
  border-radius: 20px;
  padding: 40px;
  text-align: center;
}

.placeholder-text {
  font-size: 16px;
  color: var(--vp-c-text-2);
  line-height: 1.6;
}
</style>
