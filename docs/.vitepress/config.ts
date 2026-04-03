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
          { text: 'Core Code', link: '/core-code/entry-point' },
          { text: 'Tools', link: '/tools/overview' },
          { text: 'MCP', link: '/architecture/mcp' },
          { text: 'Roadmap', link: '/roadmap' },
          { text: 'Community', link: '/community' },
          { text: 'GitHub', link: 'https://github.com/strings77wzq/claude-code-Go' }
        ],
        footer: {
          message: 'Released under the MIT License.',
          copyright: 'Copyright © 2024-present claude-code-Go Contributors'
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
          { text: '核心代码', link: '/zh/core-code/entry-point' },
          { text: '工具', link: '/zh/tools/overview' },
          { text: 'MCP', link: '/zh/architecture/mcp' },
          { text: '路线图', link: '/zh/roadmap' },
          { text: '社区', link: '/zh/community' },
          { text: 'GitHub', link: 'https://github.com/strings77wzq/claude-code-Go' }
        ],
        footer: {
          message: '基于 MIT 许可证发布',
          copyright: 'Copyright © 2024-present claude-code-Go 贡献者'
        }
      }
    }
  },

  themeConfig: {
    logo: '/logo.svg',

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
            { text: 'Agent Loop', link: '/architecture/agent-loop' },
            { text: 'Tools', link: '/architecture/tools' }
          ]
        },
        {
          text: 'Core Code',
          collapsed: false,
          items: [
            { text: 'Entry Point', link: '/core-code/entry-point' },
            { text: 'Agent Loop Implementation', link: '/core-code/agent-loop-impl' }
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
          text: 'MCP',
          collapsed: false,
          items: [
            { text: 'Integration', link: '/architecture/mcp' }
          ]
        },
        {
          text: 'Resources',
          collapsed: false,
          items: [
            { text: 'Roadmap', link: '/roadmap' },
            { text: 'Community', link: '/community' }
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
            { text: 'Agent 循环', link: '/zh/architecture/agent-loop' },
            { text: '工具', link: '/zh/architecture/tools' }
          ]
        },
        {
          text: '核心代码',
          collapsed: false,
          items: [
            { text: '入口点', link: '/zh/core-code/entry-point' },
            { text: 'Agent Loop 实现', link: '/zh/core-code/agent-loop-impl' }
          ]
        },
        {
          text: '工具系统',
          collapsed: false,
          items: [
            { text: '概览', link: '/zh/tools/overview' }
          ]
        },
        {
          text: 'MCP 集成',
          collapsed: false,
          items: [
            { text: '协议详解', link: '/zh/architecture/mcp' }
          ]
        },
        {
          text: '资源',
          collapsed: false,
          items: [
            { text: '路线图', link: '/zh/roadmap' },
            { text: '社区', link: '/zh/community' }
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
