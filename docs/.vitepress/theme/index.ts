import DefaultTheme from 'vitepress/theme'
import CustomHome from './CustomHome.vue'
import './custom.css'

export default {
  extends: DefaultTheme,
  Layout: CustomHome,
  enhanceApp({ app }) {
    document.addEventListener('DOMContentLoaded', () => {
      const observer = new IntersectionObserver((entries) => {
        entries.forEach(entry => {
          if (entry.isIntersecting) {
            entry.target.classList.add('is-visible');
          }
        });
      }, { threshold: 0.1 });

      document.querySelectorAll('.fade-in-section').forEach(el => {
        observer.observe(el);
      });
    });
  }
}
