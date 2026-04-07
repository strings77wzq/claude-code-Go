import DefaultTheme from 'vitepress/theme'
import Layout from './Layout.vue'
import CustomHome from './CustomHome.vue'
import TerminalTypewriter from './components/TerminalTypewriter.vue'
import TypedText from './components/TypedText.vue'
import CodePreview from './components/CodePreview.vue'
import AnimatedButton from './components/AnimatedButton.vue'
import Playground from './components/Playground.vue'
import GitHubStars from './components/GitHubStars.vue'
import ThemeToggle from './components/ThemeToggle.vue'
import Testimonials from './components/Testimonials.vue'
import './custom.css'

export default {
  extends: DefaultTheme,
  Layout,
  enhanceApp({ app }) {
    app.component('CustomHome', CustomHome)
    app.component('TerminalTypewriter', TerminalTypewriter)
    app.component('TypedText', TypedText)
    app.component('CodePreview', CodePreview)
    app.component('AnimatedButton', AnimatedButton)
    app.component('Playground', Playground)
    app.component('GitHubStars', GitHubStars)
    app.component('ThemeToggle', ThemeToggle)
    app.component('Testimonials', Testimonials)
    
    if (typeof window !== 'undefined') {
      window.addEventListener('DOMContentLoaded', () => {
        const addCopyButtons = () => {
          document.querySelectorAll('div[class*="language-"]:not(.has-copy-button)').forEach((block) => {
            block.classList.add('has-copy-button')
            const button = document.createElement('button')
            button.className = 'copy-button'
            button.innerHTML = `<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <rect x="9" y="9" width="13" height="13" rx="2" ry="2"></rect>
              <path d="M5 15H4a2 2 0 0 1-2-2V4a2 2 0 0 1 2-2h9a2 2 0 0 1 2 2v1"></path>
            </svg>`
            button.addEventListener('click', async (e) => {
              const code = block.querySelector('code')?.textContent || ''
              try {
                await navigator.clipboard.writeText(code)
                button.classList.add('copied')
                button.innerHTML = `<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <polyline points="20 6 9 17 4 12"></polyline>
                </svg>`
                setTimeout(() => {
                  button.classList.remove('copied')
                  button.innerHTML = `<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                    <rect x="9" y="9" width="13" height="13" rx="2" ry="2"></rect>
                    <path d="M5 15H4a2 2 0 0 1-2-2V4a2 2 0 0 1 2-2h9a2 2 0 0 1 2 2v1"></path>
                  </svg>`
                }, 2000)
              } catch (err) {
                console.error('Failed to copy:', err)
              }
            })
            block.appendChild(button)
          })
        }

        setTimeout(addCopyButtons, 100)

        const mutationObserver = new MutationObserver(addCopyButtons)
        mutationObserver.observe(document.body, { childList: true, subtree: true })

        const scrollObserver = new IntersectionObserver((entries) => {
          entries.forEach(entry => {
            if (entry.isIntersecting) {
              entry.target.classList.add('is-visible')
            }
          })
        }, { threshold: 0.1 })

        document.querySelectorAll('.fade-in-section').forEach(el => {
          scrollObserver.observe(el)
        })
      })
    }
  }
}