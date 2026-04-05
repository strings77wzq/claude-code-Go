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
          { 
            text: 'Guide', 
            items: [
              { text: 'Introduction', link: '/guide/introduction' },
              { text: 'Quick Start', link: '/guide/quick-start' },
              { text: 'Project Structure', link: '/guide/project-structure' },
              { text: 'Homebrew', link: '/guide/installation-homebrew' },
              { text: 'Session Management', link: '/guide/session-management' }
            ]
          },
          { 
            text: 'Architecture', 
            items: [
              { text: 'Overview', link: '/architecture/overview' },
              { text: 'Design Philosophy', link: '/architecture/design-philosophy' },
              { text: 'Agent Loop', link: '/architecture/agent-loop' },
              { text: 'Tools', link: '/architecture/tools' },
              { text: 'Core Code', link: '/architecture/core-code' },
              { text: 'Providers', link: '/architecture/providers' }
            ]
          },
          {
            text: 'API',
            items: [
              { text: 'Tools', link: '/api/tools' },
              { text: 'Commands', link: '/api/commands' },
              { text: 'Configuration', link: '/api/config' }
            ]
          },
          {
            text: 'Resources',
            items: [
              { text: 'Roadmap', link: '/roadmap' },
              { text: 'Troubleshooting', link: '/troubleshooting/' },
              { text: 'Contributing', link: '/contributing/' },
              { text: 'Feedback', link: '/feedback' }
            ]
          },
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
          { 
            text: '架构', 
            items: [
              { text: '概览', link: '/zh/architecture/overview' },
              { text: '交互式图表', link: '/zh/architecture/interactive-diagram' },
              { text: '设计理念', link: '/zh/architecture/design-philosophy' },
              { text: 'Agent Loop', link: '/zh/architecture/agent-loop' },
              { text: '工具系统', link: '/zh/architecture/tools' },
              { text: 'Providers', link: '/zh/architecture/providers' }
            ]
          },
          { text: '扩展', link: '/zh/extension/skills' },
          { text: '工具', link: '/zh/tools/overview' },
          {
            text: 'API 参考',
            items: [
              { text: '工具', link: '/zh/api/tools' },
              { text: '命令', link: '/zh/api/commands' },
              { text: '配置', link: '/zh/api/config' }
            ]
          },
          { text: '故障排除', link: '/zh/troubleshooting' },
          { text: '贡献指南', link: '/zh/contributing' },
          { text: 'Why Go?', link: '/zh/#why-go' },
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
            { text: 'Project Structure', link: '/guide/project-structure' },
            { text: 'Homebrew', link: '/guide/installation-homebrew' },
            { text: 'Session Management', link: '/guide/session-management' }
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
            { text: 'Core Code - Entry Point', link: '/architecture/core-code-entry' },
            { text: 'Core Code - Agent Loop', link: '/architecture/core-code-agent-loop' },
            { text: 'Providers', link: '/architecture/providers' }
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
          text: 'Troubleshooting',
          collapsed: false,
          items: [
            { text: 'Troubleshooting Guide', link: '/troubleshooting/' }
          ]
        },
        {
          text: 'Contributing',
          collapsed: false,
          items: [
            { text: 'Contributor Guide', link: '/contributing/' }
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
            { text: '交互式图表', link: '/zh/architecture/interactive-diagram' },
            { text: '设计理念', link: '/zh/architecture/design-philosophy' },
            { text: 'Agent Loop', link: '/zh/architecture/agent-loop' },
            { text: '工具系统', link: '/zh/architecture/tools' },
            { text: 'Providers', link: '/zh/architecture/providers' }
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
          text: '故障排除',
          collapsed: false,
          items: [
            { text: '故障排除指南', link: '/zh/troubleshooting' }
          ]
        },
        {
          text: '贡献指南',
          collapsed: false,
          items: [
            { text: '贡献者指南', link: '/zh/contributing' }
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