import { defineConfig } from 'vitepress'

export default defineConfig({
  title: 'claude-code-Go',
  description: 'Claude Code in Go — AI-powered coding assistant',
  base: '/claude-code-Go/',

  locales: {
    root: {
      label: 'English',
      lang: 'en',
      title: 'claude-code-Go',
      description: 'Claude Code in Go — AI-powered coding assistant',
      themeConfig: {
        nav: [
          { text: 'Home', link: '/' },
          { text: 'Guide', link: '/guide/introduction' },
          { 
            text: 'Tutorials',
            items: [
              { text: 'Quick Start', link: '/guide/tutorials/01-quickstart' },
              { text: 'First Tool Call', link: '/guide/tutorials/02-first-tool-call' },
              { text: 'Agent Loop', link: '/guide/tutorials/03-agent-loop' },
              { text: 'All Tutorials', link: '/guide/tutorials/' }
            ]
          },
          { text: 'Architecture', link: '/architecture/overview' },
          { text: 'API', link: '/api/tools' },
          { 
            text: 'Resources',
            items: [
              { text: 'Roadmap', link: '/roadmap' },
              { text: 'Troubleshooting', link: '/troubleshooting/common-issues' },
              { text: 'Benchmarks', link: '/benchmark' },
              { text: 'Showcase', link: '/showcase' }
            ]
          },
          { text: 'GitHub ⭐', link: 'https://github.com/strings77wzq/claude-code-Go' }
        ],
        footer: {
          message: 'Released under the MIT License.',
          copyright: 'Copyright © 2026-present claude-code-Go Contributors'
        }
      }
    },
    zh: {
      label: '中文',
      lang: 'zh-CN',
      title: 'claude-code-Go',
      description: 'Go 实现的 Claude Code — AI 编程助手',
      themeConfig: {
        nav: [
          { text: '指南', link: '/zh/guide/introduction' },
          { text: '架构', link: '/zh/architecture/overview' },
          { text: 'API', link: '/zh/api/tools' },
          { text: '资源', link: '/zh/roadmap' },
          { text: 'GitHub', link: 'https://github.com/strings77wzq/claude-code-Go' }
        ],
        footer: {
          message: '基于 MIT 许可证发布',
          copyright: 'Copyright © 2026-present claude-code-Go 贡献者'
        }
      }
    }
  },

  themeConfig: {
    logo: '/logo.svg',
    appearance: true,

    siteTitle: 'claude-code-Go',

    sidebar: {
      '/': [
        {
          text: 'Guide',
          collapsed: false,
          items: [
            { text: 'Introduction', link: '/guide/introduction' },
            { text: 'Quick Start', link: '/guide/quick-start' },
            { text: 'Project Structure', link: '/guide/project-structure' }
          ]
        },
        {
          text: 'Architecture',
          collapsed: false,
          items: [
            { text: 'Overview', link: '/architecture/overview' },
            { text: 'Design Philosophy', link: '/architecture/design-philosophy' },
            { text: 'Agent Loop', link: '/architecture/agent-loop' },
            { text: 'Tools', link: '/architecture/tools' },
            { text: 'Providers', link: '/architecture/providers' }
          ]
        },
        {
          text: 'API Reference',
          collapsed: false,
          items: [
            { text: 'Tools', link: '/api/tools' },
            { text: 'Commands', link: '/api/commands' },
            { text: 'Configuration', link: '/api/config' }
          ]
        },
        {
          text: 'Resources',
          collapsed: false,
          items: [
            { text: 'Roadmap', link: '/roadmap' },
            { text: 'Troubleshooting', link: '/troubleshooting' },
            { text: 'Contributing', link: '/contributing' },
            { text: 'Feedback', link: '/feedback' }
          ]
        }
      ],
      '/zh/': [
        {
          text: '指南',
          collapsed: false,
          items: [
            { text: '项目简介', link: '/zh/guide/introduction' },
            { text: '快速开始', link: '/zh/guide/quick-start' },
            { text: '项目结构', link: '/zh/guide/project-structure' }
          ]
        },
        {
          text: '架构',
          collapsed: false,
          items: [
            { text: '概览', link: '/zh/architecture/overview' },
            { text: '设计理念', link: '/zh/architecture/design-philosophy' },
            { text: 'Agent Loop', link: '/zh/architecture/agent-loop' },
            { text: '工具系统', link: '/zh/architecture/tools' },
            { text: 'Providers', link: '/zh/architecture/providers' }
          ]
        },
        {
          text: 'API 参考',
          collapsed: false,
          items: [
            { text: '工具', link: '/zh/api/tools' },
            { text: '命令', link: '/zh/api/commands' },
            { text: '配置', link: '/zh/api/config' }
          ]
        },
        {
          text: '资源',
          collapsed: false,
          items: [
            { text: 'Roadmap', link: '/zh/roadmap' },
            { text: '故障排除', link: '/zh/troubleshooting' },
            { text: '贡献指南', link: '/zh/contributing' },
            { text: '反馈', link: '/zh/feedback' }
          ]
        }
      ]
    },

    socialLinks: [
      {
        icon: 'github',
        link: 'https://github.com/strings77wzq/claude-code-Go',
        ariaLabel: 'GitHub'
      }
    ],

    outline: {
      level: [2, 3]
    }
  },

  markdown: {
    lineNumbers: true
  },

  lastUpdated: true,

  ignoreDeadLinks: true
})
