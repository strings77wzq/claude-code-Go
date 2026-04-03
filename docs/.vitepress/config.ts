import { defineConfig } from 'vitepress/config'

export default defineConfig({
  title: 'claude-code-Go',
  description: 'Claude Code in Go — AI-powered coding assistant',
  base: '/claude-code-Go/',

  locales: {
    root: {
      label: 'English',
      lang: 'en',
      title: 'claude-code-Go',
      description: 'Claude Code in Go — AI-powered coding assistant'
    },
    zh: {
      label: '中文',
      lang: 'zh-CN',
      title: 'claude-code-Go',
      description: 'Go 实现的 Claude Code — AI 编程助手'
    }
  },

  themeConfig: {
    logo: '/logo.svg',

    nav: [
      {
        text: 'Guide',
        link: '/guide/installation'
      },
      {
        text: 'Architecture',
        link: '/architecture/overview'
      },
      {
        text: 'GitHub',
        link: 'https://github.com/strings77wzq/claude-code-Go'
      }
    ],

    sidebar: {
      '/': [
        {
          text: 'Guide',
          collapsed: false,
          items: [
            { text: 'Installation', link: '/guide/installation' },
            { text: 'Quick Start', link: '/guide/quick-start' },
            { text: 'Configuration', link: '/guide/configuration' }
          ]
        },
        {
          text: 'Architecture',
          collapsed: false,
          items: [
            { text: 'Overview', link: '/architecture/overview' },
            { text: 'Agent Loop', link: '/architecture/agent-loop' },
            { text: 'Tools', link: '/architecture/tools' }
          ]
        }
      ],
      '/zh/': [
        {
          text: '指南',
          collapsed: false,
          items: [
            { text: '安装', link: '/zh/guide/installation' },
            { text: '快速开始', link: '/zh/guide/quick-start' },
            { text: '配置', link: '/zh/guide/configuration' }
          ]
        },
        {
          text: '架构',
          collapsed: false,
          items: [
            { text: '概述', link: '/zh/architecture/overview' },
            { text: 'Agent 循环', link: '/zh/architecture/agent-loop' },
            { text: '工具系统', link: '/zh/architecture/tools' }
          ]
        }
      ]
    },

    socialLinks: [
      {
        icon: 'github',
        link: 'https://github.com/strings77wzq/claude-code-Go'
      }
    ],

    footer: {
      message: 'Released under the MIT License.',
      copyright: 'Copyright © 2024-present claude-code-Go Contributors'
    },

    outline: {
      level: [2, 3]
    }
  },

  markdown: {
    lineNumbers: true
  },

  lastUpdated: true
})
