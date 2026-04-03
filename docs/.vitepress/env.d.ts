/// <reference types="vite/client" />

interface UseCaseItem {
  icon: string
  title: string
  details: string
}

interface CTAAction {
  text: string
  link: string
}

interface CTAItem {
  title: string
  details: string
  actions: CTAAction[]
}

interface CustomFrontMatter {
  useCases?: UseCaseItem[]
  cta?: CTAItem[]
}

declare module 'vitepress' {
  interface PageData {
    frontmatter: Record<string, any> & CustomFrontMatter
  }
}
