# Kubemux MCP Server

[English](#english) | [中文](#中文)

## English

### Overview

The Kubemux MCP (Model Context Protocol) Server enables AI assistants like Claude, GitHub Copilot, and Google Gemini to interact with your Kubernetes clusters through kubemux. This allows you to manage Kubernetes environments using natural language commands.

### What is MCP?

Model Context Protocol (MCP) is an open standard that enables AI assistants to safely interact with local tools and data sources. By implementing an MCP server, kubemux can be used as a plugin/extension for various AI assistants.

### Features

The MCP server exposes the following capabilities:

- **List Kubernetes Clusters**: View all available kubeconfig files
- **List Sessions**: See active tmux/zellij sessions
- **Create Sessions**: Generate commands to create new kubemux sessions
- **Attach to Sessions**: Get commands to attach to existing sessions
- **Multi-plexer Support**: Works with both tmux and zellij

### Installation

#### 1. Build the MCP Server

```bash
cd /path/to/kubemux
go build -o kubemux-mcp-server ./mcp-server
```

#### 2. Install the Binary

```bash
# Linux/macOS
sudo cp kubemux-mcp-server /usr/local/bin/
sudo chmod +x /usr/local/bin/kubemux-mcp-server
```

#### 3. Configure Your AI Assistant

##### For Claude Desktop

Edit your Claude Desktop configuration file:

**macOS**: `~/Library/Application Support/Claude/claude_desktop_config.json`
**Windows**: `%APPDATA%\Claude\claude_desktop_config.json`

Add the kubemux server:

```json
{
  "mcpServers": {
    "kubemux": {
      "command": "kubemux-mcp-server",
      "args": [],
      "env": {}
    }
  }
}
```

##### For GitHub Copilot (VS Code)

1. Install the MCP extension for VS Code (if available)
2. Configure the MCP server in VS Code settings:

```json
{
  "mcp.servers": {
    "kubemux": {
      "command": "kubemux-mcp-server"
    }
  }
}
```

##### For Google Gemini / Other AI Tools

The MCP server follows the standard MCP protocol. Configure it according to your AI tool's MCP integration documentation, using:

- **Command**: `kubemux-mcp-server`
- **Protocol**: JSON-RPC over stdio

### Usage Examples

Once configured, you can use natural language to interact with kubemux:

#### Example Conversations

**User**: "List all my Kubernetes clusters"
**AI**: *Calls `list_kubeconfigs` tool and returns the list*

**User**: "Show me active kubemux sessions"
**AI**: *Calls `list_sessions` tool and displays current sessions*

**User**: "Create a new kubemux session for my production cluster"
**AI**: *Calls `create_session` tool with appropriate parameters and provides the command*

**User**: "Connect to the dev-cluster session"
**AI**: *Calls `attach_session` tool and provides the attachment command*

### Available Tools

#### 1. list_kubeconfigs
Lists all kubeconfig files in `~/.kube` directory.

**Parameters**: None

**Returns**: JSON with list of kubeconfig files

#### 2. list_clusters
Alias for `list_kubeconfigs`.

#### 3. list_sessions
Lists all active kubemux/tmux/zellij sessions.

**Parameters**:
- `plexer` (optional): "tmux" or "zellij"

**Returns**: JSON with session list

#### 4. create_session
Generates a command to create a new kubemux session.

**Parameters**:
- `project` (optional): Project name/configuration (default: "default")
- `kubeconfig` (optional): Kubeconfig file to use
- `plexer` (optional): "tmux" or "zellij"
- `directory` (optional): Configuration directory (default: "~/.tmuxinator")

**Returns**: Command to execute

#### 5. attach_session
Generates a command to attach to an existing session.

**Parameters**:
- `session_name` (required): Name of the session
- `plexer` (optional): "tmux" or "zellij"

**Returns**: Command to execute

### Development

#### Testing the MCP Server

You can test the MCP server manually using stdio:

```bash
echo '{"jsonrpc":"2.0","id":1,"method":"initialize","params":{}}' | kubemux-mcp-server
echo '{"jsonrpc":"2.0","id":2,"method":"tools/list"}' | kubemux-mcp-server
```

#### Adding New Tools

1. Create a new tool in `mcp-server/tools/`
2. Implement the `Tool` interface
3. Register the tool in `mcp-server/server/server.go`

### Troubleshooting

**MCP server not responding**:
- Check that `kubemux-mcp-server` is in your PATH
- Verify the binary has execute permissions
- Check AI assistant logs for errors

**Tools not appearing**:
- Restart your AI assistant after configuration changes
- Verify the MCP configuration file syntax

**Permission errors**:
- Ensure kubemux has access to `~/.kube` directory
- Check that tmux/zellij are installed if using session tools

---

## 中文

### 概述

Kubemux MCP（模型上下文协议）服务器使 Claude、GitHub Copilot 和 Google Gemini 等 AI 助手能够通过 kubemux 与您的 Kubernetes 集群交互。这允许您使用自然语言命令管理 Kubernetes 环境。

### 什么是 MCP？

模型上下文协议（MCP）是一个开放标准，使 AI 助手能够安全地与本地工具和数据源交互。通过实现 MCP 服务器，kubemux 可以作为各种 AI 助手的插件/扩展使用。

### 功能特性

MCP 服务器提供以下功能：

- **列出 Kubernetes 集群**：查看所有可用的 kubeconfig 文件
- **列出会话**：查看活动的 tmux/zellij 会话
- **创建会话**：生成创建新 kubemux 会话的命令
- **附加到会话**：获取附加到现有会话的命令
- **多路复用器支持**：同时支持 tmux 和 zellij

### 安装步骤

#### 1. 构建 MCP 服务器

```bash
cd /path/to/kubemux
go build -o kubemux-mcp-server ./mcp-server
```

#### 2. 安装二进制文件

```bash
# Linux/macOS
sudo cp kubemux-mcp-server /usr/local/bin/
sudo chmod +x /usr/local/bin/kubemux-mcp-server
```

#### 3. 配置您的 AI 助手

##### Claude Desktop 配置

编辑 Claude Desktop 配置文件：

**macOS**: `~/Library/Application Support/Claude/claude_desktop_config.json`
**Windows**: `%APPDATA%\Claude\claude_desktop_config.json`

添加 kubemux 服务器：

```json
{
  "mcpServers": {
    "kubemux": {
      "command": "kubemux-mcp-server",
      "args": [],
      "env": {}
    }
  }
}
```

##### GitHub Copilot (VS Code) 配置

1. 安装 VS Code 的 MCP 扩展（如果可用）
2. 在 VS Code 设置中配置 MCP 服务器：

```json
{
  "mcp.servers": {
    "kubemux": {
      "command": "kubemux-mcp-server"
    }
  }
}
```

##### Google Gemini / 其他 AI 工具

MCP 服务器遵循标准的 MCP 协议。根据您的 AI 工具的 MCP 集成文档进行配置，使用：

- **命令**: `kubemux-mcp-server`
- **协议**: JSON-RPC over stdio

### 使用示例

配置完成后，您可以使用自然语言与 kubemux 交互：

#### 示例对话

**用户**："列出我所有的 Kubernetes 集群"
**AI**：*调用 `list_kubeconfigs` 工具并返回列表*

**用户**："显示活动的 kubemux 会话"
**AI**：*调用 `list_sessions` 工具并显示当前会话*

**用户**："为我的生产集群创建一个新的 kubemux 会话"
**AI**：*使用适当的参数调用 `create_session` 工具并提供命令*

**用户**："连接到 dev-cluster 会话"
**AI**：*调用 `attach_session` 工具并提供附加命令*

### 可用工具

#### 1. list_kubeconfigs
列出 `~/.kube` 目录中的所有 kubeconfig 文件。

**参数**：无

**返回**：包含 kubeconfig 文件列表的 JSON

#### 2. list_clusters
`list_kubeconfigs` 的别名。

#### 3. list_sessions
列出所有活动的 kubemux/tmux/zellij 会话。

**参数**：
- `plexer`（可选）："tmux" 或 "zellij"

**返回**：包含会话列表的 JSON

#### 4. create_session
生成创建新 kubemux 会话的命令。

**参数**：
- `project`（可选）：项目名称/配置（默认："default"）
- `kubeconfig`（可选）：要使用的 kubeconfig 文件
- `plexer`（可选）："tmux" 或 "zellij"
- `directory`（可选）：配置目录（默认："~/.tmuxinator"）

**返回**：要执行的命令

#### 5. attach_session
生成附加到现有会话的命令。

**参数**：
- `session_name`（必需）：会话名称
- `plexer`（可选）："tmux" 或 "zellij"

**返回**：要执行的命令

### 开发

#### 测试 MCP 服务器

您可以使用 stdio 手动测试 MCP 服务器：

```bash
echo '{"jsonrpc":"2.0","id":1,"method":"initialize","params":{}}' | kubemux-mcp-server
echo '{"jsonrpc":"2.0","id":2,"method":"tools/list"}' | kubemux-mcp-server
```

#### 添加新工具

1. 在 `mcp-server/tools/` 中创建新工具
2. 实现 `Tool` 接口
3. 在 `mcp-server/server/server.go` 中注册工具

### 故障排除

**MCP 服务器无响应**：
- 检查 `kubemux-mcp-server` 是否在您的 PATH 中
- 验证二进制文件具有执行权限
- 检查 AI 助手日志中的错误

**工具未显示**：
- 配置更改后重启您的 AI 助手
- 验证 MCP 配置文件语法

**权限错误**：
- 确保 kubemux 可以访问 `~/.kube` 目录
- 如果使用会话工具，请检查是否安装了 tmux/zellij

### 架构说明

MCP 服务器使用以下架构：

```
┌─────────────────┐
│  AI Assistant   │
│ (Claude/Copilot)│
└────────┬────────┘
         │ JSON-RPC
         │ (stdio)
┌────────▼────────┐
│  MCP Server     │
│ (kubemux-mcp)   │
├─────────────────┤
│   Tool Layer    │
│  - list_*       │
│  - create_*     │
│  - attach_*     │
└────────┬────────┘
         │
┌────────▼────────┐
│   Kubemux Lib   │
│   + CLI Tool    │
└─────────────────┘
```

### 贡献

欢迎提交 Pull Request！对于重大更改，请先开一个 issue 讨论您想要改变的内容。

### 许可证

[MIT](../LICENSE)
