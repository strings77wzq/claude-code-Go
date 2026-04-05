---
title: 内置工具 API 参考
description: go-code 所有 10 个内置工具的完整 API 参考
---

# 内置工具 API 参考

go-code 提供 10 个内置工具，使代理能够与文件系统交互、执行 shell 命令并管理任务跟踪。本参考文档记录了每个工具的参数、权限和使用示例。

## 工具列表

| # | 工具名称 | 描述 | 所需权限 |
|---|-----------|-------------|---------------------|
| 1 | [Read](#read) | 读取文件内容，支持可选的偏移量/限制 | 否 (ReadOnly) |
| 2 | [Write](#write) | 创建或覆盖文件 | 是 (WorkspaceWrite) |
| 3 | [Edit](#edit) | 使用精确字符串匹配进行针对性代码编辑 | 是 (WorkspaceWrite) |
| 4 | [Glob](#glob) | 通过 glob 模式查找文件 | 否 (ReadOnly) |
| 5 | [Grep](#grep) | 使用正则表达式搜索文件内容 | 否 (ReadOnly) |
| 6 | [Bash](#bash) | 执行 shell 命令 | 是 (WorkspaceWrite) |
| 7 | [Diff](#diff) | 比较两个内容字符串 | 否 (ReadOnly) |
| 8 | [Tree](#tree) | 显示目录树结构 | 否 (ReadOnly) |
| 9 | [WebFetch](#webfetch) | 获取 URL 并返回可读文本 | 是 (WorkspaceWrite) |
| 10 | [TodoWrite](#todowrite) | 创建和管理待办事项 | 否 (WorkspaceWrite) |

---

## Read

读取文件并返回带行号的内容。

### 工具定义

```json
{
  "name": "Read",
  "description": "Reads a file and returns its contents with line numbers.",
  "input_schema": {
    "type": "object",
    "properties": {
      "file_path": {
        "type": "string",
        "description": "Path to the file to read"
      },
      "offset": {
        "type": "number",
        "description": "Line number to start reading from (0-based, default: 0)"
      },
      "limit": {
        "type": "number",
        "description": "Maximum number of lines to read (default: 2000)"
      }
    },
    "required": ["file_path"]
  }
}
```

### 参数

| 参数 | 类型 | 必填 | 默认值 | 描述 |
|-----------|------|----------|---------|-------------|
| `file_path` | string | 是 | - | 要读取的文件路径 |
| `offset` | number | 否 | 0 | 开始读取的行号（从 0 开始） |
| `limit` | number | 否 | 2000 | 最大读取行数 |

### 权限级别

- **需要权限**: 否
- **权限级别**: `ReadOnly`

### 限制条件

- 最大文件大小：200KB
- 不能读取目录（将返回错误）

### 示例

**读取整个文件：**
```
输入: { "file_path": "/home/user/project/main.go" }
输出:
1: package main
2:
3: func main() {
4:     fmt.Println("Hello, World!")
5: }
```

**读取指定行：**
```
输入: { "file_path": "/home/user/project/main.go", "offset": 10, "limit": 20 }
```

---

## Write

创建新文件或覆盖现有文件。

### 工具定义

```json
{
  "name": "Write",
  "description": "Creates a new file or overwrites an existing file.",
  "input_schema": {
    "type": "object",
    "properties": {
      "file_path": {
        "type": "string",
        "description": "Path to the file to write"
      },
      "content": {
        "type": "string",
        "description": "Content to write to the file"
      }
    },
    "required": ["file_path", "content"]
  }
}
```

### 参数

| 参数 | 类型 | 必填 | 描述 |
|-----------|------|----------|-------------|
| `file_path` | string | 是 | 要写入的文件路径 |
| `content` | string | 是 | 要写入文件的内容 |

### 权限级别

- **需要权限**: 是
- **权限级别**: `WorkspaceWrite`

### 行为

- 如果父目录不存在则创建
- 覆盖现有文件而不确认
- 返回带文件路径的成功消息

### 示例

**创建新文件：**
```
输入: {
  "file_path": "/home/user/project/README.md",
  "content": "# My Project\n\nThis is my project."
}
输出: Successfully wrote to /home/user/project/README.md
```

**创建带嵌套目录的文件：**
```
输入: {
  "file_path": "/home/user/project/src/lib/utils.go",
  "content": "package utils\n\nfunc Hello() string { return \"Hello\" }"
}
输出: Successfully wrote to /home/user/project/src/lib/utils.go
```

---

## Edit

在文件中执行精确字符串替换。

### 工具定义

```json
{
  "name": "Edit",
  "description": "Performs exact string replacement in a file.",
  "input_schema": {
    "type": "object",
    "properties": {
      "file_path": {
        "type": "string",
        "description": "Path to the file to edit"
      },
      "old_string": {
        "type": "string",
        "description": "The exact string to replace"
      },
      "new_string": {
        "type": "string",
        "description": "The replacement string"
      },
      "replace_all": {
        "type": "boolean",
        "description": "Replace all occurrences (default: false)"
      }
    },
    "required": ["file_path", "old_string", "new_string"]
  }
}
```

### 参数

| 参数 | 类型 | 必填 | 默认值 | 描述 |
|-----------|------|----------|---------|-------------|
| `file_path` | string | 是 | - | 要编辑的文件路径 |
| `old_string` | string | 是 | - | 要替换的精确字符串 |
| `new_string` | string | 是 | - | 替换字符串 |
| `replace_all` | boolean | 否 | false | 替换所有出现的位置 |

### 权限级别

- **需要权限**: 是
- **权限级别**: `WorkspaceWrite`

### 限制条件

- 需要精确字符串匹配（对空白敏感）
- 如果 `old_string` 出现多次，默认报错，除非 `replace_all=true`
- 最大文件大小：200KB

### 示例

**简单替换：**
```
输入: {
  "file_path": "/home/user/project/main.go",
  "old_string": "fmt.Println(\"Hello\")",
  "new_string": "fmt.Println(\"Hello, World!\")"
}
输出: Successfully edited /home/user/project/main.go
```

**替换所有出现的位置：**
```
输入: {
  "file_path": "/home/user/project/main.go",
  "old_string": "TODO",
  "new_string": "DONE",
  "replace_all": true
}
输出: Successfully edited /home/user/project/main.go
```

---

## Glob

查找与 glob 模式匹配的文件。

### 工具定义

```json
{
  "name": "Glob",
  "description": "Finds files matching glob patterns.",
  "input_schema": {
    "type": "object",
    "properties": {
      "path": {
        "type": "string",
        "description": "Glob pattern to match files (e.g., *.go, **/*.ts)"
      }
    },
    "required": ["path"]
  }
}
```

### 参数

| 参数 | 类型 | 必填 | 描述 |
|-----------|------|----------|-------------|
| `path` | string | 是 | 匹配文件的 glob 模式 |

### 权限级别

- **需要权限**: 否
- **权限级别**: `ReadOnly`

### 模式语法

| 模式 | 描述 | 示例 |
|---------|-------------|---------|
| `*` | 匹配单个目录中的任意字符 | `*.go` 匹配所有 `.go` 文件 |
| `**` | 跨目录递归匹配 | `**/*.ts` 匹配所有 `.ts` 文件 |
| `?` | 匹配单个字符 | `file?.txt` 匹配 `file1.txt` |
| `[abc]` | 匹配字符类 | `[abc].txt` 匹配 `a.txt`、`b.txt` |

### 示例

**查找所有 Go 文件：**
```
输入: { "path": "**/*.go" }
输出:
/home/user/project/main.go
/home/user/project/utils/helper.go
/home/user/project/cmd/app.go
```

**查找特定扩展名：**
```
输入: { "path": "*.md" }
输出:
README.md
CHANGELOG.md
```

---

## Grep

使用正则表达式搜索文件内容。

### 工具定义

```json
{
  "name": "Grep",
  "description": "Searches file contents using regular expressions.",
  "input_schema": {
    "type": "object",
    "properties": {
      "pattern": {
        "type": "string",
        "description": "Regular expression pattern to search for"
      },
      "path": {
        "type": "string",
        "description": "File or directory path to search in (glob patterns supported)"
      },
      "include": {
        "type": "string",
        "description": "File pattern to include (e.g., *.go, *.ts)"
      },
      "output_mode": {
        "type": "string",
        "description": "Output format: files_with_matches, content, or count"
      }
    },
    "required": ["pattern", "path"]
  }
}
```

### 参数

| 参数 | 类型 | 必填 | 默认值 | 描述 |
|-----------|------|----------|---------|-------------|
| `pattern` | string | 是 | - | 正则表达式模式 |
| `path` | string | 是 | - | 要搜索的文件或目录 |
| `include` | string | 否 | - | 要包含的文件模式 |
| `output_mode` | string | 否 | content | 输出格式 |

### 权限级别

- **需要权限**: 否
- **权限级别**: `ReadOnly`

### 输出模式

| 模式 | 描述 |
|------|-------------|
| `content` | 显示带上下文的匹配行（默认） |
| `files_with_matches` | 仅显示包含匹配的文件路径 |
| `count` | 显示每个文件的匹配数量 |

### 示例

**搜索函数定义：**
```
输入: { "pattern": "func.*main", "path": "/home/user/project", "include": "*.go" }
输出:
main.go:5: func main() {
main.go:10: func mainLogic() {
```

**查找所有 TODO 注释：**
```
输入: { "pattern": "TODO", "path": "/home/user/project", "include": "*.go", "output_mode": "files_with_matches" }
输出:
/home/user/project/main.go
/home/user/project/utils/helper.go
```

---

## Bash

执行 shell 命令。

### 工具定义

```json
{
  "name": "Bash",
  "description": "Executes shell commands in the project directory.",
  "input_schema": {
    "type": "object",
    "properties": {
      "command": {
        "type": "string",
        "description": "Shell command to execute"
      }
    },
    "required": ["command"]
  }
}
```

### 参数

| 参数 | 类型 | 必填 | 描述 |
|-----------|------|----------|-------------|
| `command` | string | 是 | 要执行的 shell 命令 |

### 权限级别

- **需要权限**: 是
- **权限级别**: `WorkspaceWrite`

### 行为

- 在项目根目录执行
- 默认超时：120 秒
- 输出在 100KB 处截断
- 工作目录：项目根目录

### 示例

**运行测试：**
```
输入: { "command": "go test ./..." }
输出: PASS
ok  	github.com/user/project	0.523s
```

**列出文件：**
```
输入: { "command": "ls -la" }
输出: total 48
drwxr-xr-x  4 user  staff   128 Apr  5 10:00 .
drwxr-xr-x  2 user  staff   256 Apr  5 10:01 src
```

**Git 操作：**
```
输入: { "command": "git status" }
输出: On branch main
Your branch is up to date with 'origin/main'.
```

---

## Diff

比较两个内容字符串并返回统一的 diff 输出。

### 工具定义

```json
{
  "name": "Diff",
  "description": "Compares two content strings and returns unified diff output.",
  "input_schema": {
    "type": "object",
    "properties": {
      "old_string": {
        "type": "string",
        "description": "Original content"
      },
      "new_string": {
        "type": "string",
        "description": "Modified content"
      }
    },
    "required": ["old_string", "new_string"]
  }
}
```

### 参数

| 参数 | 类型 | 必填 | 描述 |
|-----------|------|----------|-------------|
| `old_string` | string | 是 | 原始内容 |
| `new_string` | string | 是 | 修改后的内容 |

### 权限级别

- **需要权限**: 否
- **权限级别**: `ReadOnly`

### 行为

- 如果可用则使用系统 `diff` 命令
- 回退到纯 Go 实现
- 返回统一的 diff 格式

### 示例

**比较两个字符串：**
```
输入: {
  "old_string": "func Hello() string {\n    return \"Hello\"\n}",
  "new_string": "func Hello() string {\n    return \"Hello, World!\"\n}"
}
输出:
--- 
+++ 
@@ -1,3 +1,3 @@
 func Hello() string {
-    return "Hello"
+    return "Hello, World!"
 }
```

---

## Tree

将目录树结构显示为文本。

### 工具定义

```json
{
  "name": "Tree",
  "description": "Displays directory tree structure as text.",
  "input_schema": {
    "type": "object",
    "properties": {
      "path": {
        "type": "string",
        "description": "Root directory path (default: current directory)"
      },
      "limit": {
        "type": "number",
        "description": "Maximum directory depth to display (default: 3)"
      }
    },
    "required": []
  }
}
```

### 参数

| 参数 | 类型 | 必填 | 默认值 | 描述 |
|-----------|------|----------|---------|-------------|
| `path` | string | 否 | 当前目录 | 根目录路径 |
| `limit` | number | 否 | 3 | 最大显示深度 |

### 权限级别

- **需要权限**: 否
- **权限级别**: `ReadOnly`

### 示例

**显示项目结构：**
```
输入: { "path": "/home/user/project", "limit": 2 }
输出:
project/
├── cmd/
│   └── app/
├── internal/
│   └── handler/
├── pkg/
│   └── utils/
└── go.mod
```

---

## WebFetch

获取 URL 并返回可读文本（去除 HTML）。

### 工具定义

```json
{
  "name": "WebFetch",
  "description": "Fetches a URL and returns readable text with HTML tags stripped.",
  "input_schema": {
    "type": "object",
    "properties": {
      "url": {
        "type": "string",
        "description": "URL to fetch"
      },
      "goal": {
        "type": "string",
        "description": "Specific information to extract from the page"
      }
    },
    "required": ["url"]
  }
}
```

### 参数

| 参数 | 类型 | 必填 | 描述 |
|-----------|------|----------|-------------|
| `url` | string | 是 | 要获取的 URL |
| `goal` | string | 否 | 要从页面提取的特定信息 |

### 权限级别

- **需要权限**: 是
- **权限级别**: `WorkspaceWrite`

### 限制条件

- 输出限制为 50KB
- HTML 标签自动去除

### 示例

**获取页面：**
```
输入: { "url": "https://example.com" }
输出: Example Domain
This domain is for use in illustrative examples in documents...
```

**提取特定信息：**
```
输入: { "url": "https://api.github.com/repos/user/repo", "goal": "description" }
输出: A cool project that does stuff
```

---

## TodoWrite

创建、更新和管理待办事项以跟踪任务进度。

### 工具定义

```json
{
  "name": "TodoWrite",
  "description": "Create, update, and manage todo items for tracking task progress",
  "input_schema": {
    "type": "object",
    "properties": {
      "todos": {
        "type": "array",
        "description": "Array of todo items",
        "items": {
          "type": "object",
          "properties": {
            "id": {
              "type": "integer",
              "description": "Optional. If provided, update existing todo; if not, create new."
            },
            "content": {
              "type": "string",
              "description": "The content of the todo item"
            },
            "status": {
              "type": "string",
              "description": "Status: pending, in_progress, or completed",
              "enum": ["pending", "in_progress", "completed"]
            }
          }
        }
      }
    },
    "required": ["todos"]
  }
}
```

### 参数

| 参数 | 类型 | 必填 | 描述 |
|-----------|------|----------|-------------|
| `todos` | array | 是 | 待办事项数组 |

### 权限级别

- **需要权限**: 否
- **权限级别**: `WorkspaceWrite`

### 待办事项属性

| 属性 | 类型 | 必填 | 描述 |
|----------|------|----------|-------------|
| `id` | integer | 否 | 如果提供则更新现有；否则创建新的 |
| `content` | string | 是 | 待办事项内容 |
| `status` | string | 否 | 之一：pending, in_progress, completed |

### 示例

**创建待办：**
```
输入: {
  "todos": [
    { "content": "Fix the login bug", "status": "in_progress" },
    { "content": "Add unit tests", "status": "pending" }
  ]
}
输出:
Todo List:
🔄 [1] in_progress - Fix the login bug
📋 [2] pending - Add unit tests
```

**更新待办状态：**
```
输入: {
  "todos": [
    { "id": 1, "status": "completed" }
  ]
}
输出:
Todo List:
✅ [1] completed - Fix the login bug
📋 [2] pending - Add unit tests
```

---

## 相关文档

- [工具系统概览](../tools/overview.md) — 架构和扩展指南
- [配置指南](../guide/configuration.md) — 工具配置
- [MCP 集成](../extension/mcp.md) — 通过 MCP 的外部工具