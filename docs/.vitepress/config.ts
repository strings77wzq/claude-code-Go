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
          { text: 'Guide', link: '/guide/introduction' },
          { text: 'Architecture', link: '/architecture/overview' },
          { text: 'Extensions', link: '/extension/skills' },
          { text: 'Tools', link: '/tools/overview' },
          { text: 'Resources', link: '/roadmap' },
          { text: 'GitHub', link: 'https://github.com/strings77wzq/claude-code-Go' }
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
          { text: '扩展', link: '/zh/extension/skills' },
          { text: '工具', link: '/zh/tools/overview' },
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
    appearance: { label: '自动', value: 'auto' },

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
            { text: 'Agent Loop', link: '/architecture/agent-loop' }
          ]
        },
        {
          text: 'Extensions',
          collapsed: false,
          items: [
            { text: 'Skills Tutorial', link: '/extension/skills' },
            { text: 'MCP Protocol', link: '/extension/mcp' },
            { text: 'Hooks System', link: '/extension/hooks' }
          ]
        },
        {
          text: 'Tools',
          collapsed: false,
          items: [
            { text: 'Overview', link: '/tools/overview' }
          ]
        },
        {
          text: 'Resources',
          collapsed: false,
          items: [
            { text: 'Roadmap', link: '/roadmap' },
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
            { text: 'Agent Loop', link: '/zh/architecture/agent-loop' }
          ]
        },
        {
          text: '扩展',
          collapsed: false,
          items: [
            { text: 'Skills 教程', link: '/zh/extension/skills' },
            { text: 'MCP 协议', link: '/zh/extension/mcp' },
            { text: 'Hooks 系统', link: '/zh/extension/hooks' }
          ]
        },
        {
          text: '工具',
          collapsed: false,
          items: [
            { text: '概览', link: '/zh/tools/overview' }
          ]
        },
        {
          text: '资源',
          collapsed: false,
          items: [
            { text: 'Roadmap', link: '/zh/roadmap' },
            { text: '反馈', link: '/zh/feedback' }
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