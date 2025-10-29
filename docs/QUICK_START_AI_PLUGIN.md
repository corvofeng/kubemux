# Kubemux AI Plugin - Quick Start Guide

## 概述 / Overview

本指南帮助您快速将 kubemux 集成到 AI 助手中，实现通过自然语言管理 Kubernetes 集群。

This guide helps you quickly integrate kubemux with AI assistants to manage Kubernetes clusters using natural language.

---

## 5分钟快速开始 / 5-Minute Quick Start

### 步骤 1: 构建 MCP 服务器 / Step 1: Build MCP Server

```bash
cd /path/to/kubemux
make build-mcp
```

### 步骤 2: 安装 / Step 2: Install

```bash
# Linux/macOS
sudo cp kubemux-mcp-server /usr/local/bin/
sudo chmod +x /usr/local/bin/kubemux-mcp-server

# 验证安装 / Verify installation
which kubemux-mcp-server
```

### 步骤 3: 配置 AI 助手 / Step 3: Configure AI Assistant

#### For Claude Desktop (推荐 / Recommended)

**macOS**:
```bash
# 编辑配置 / Edit configuration
nano ~/Library/Application\ Support/Claude/claude_desktop_config.json
```

**Linux**:
```bash
# 创建目录 / Create directory
mkdir -p ~/.config/Claude

# 编辑配置 / Edit configuration
nano ~/.config/Claude/claude_desktop_config.json
```

**配置内容 / Configuration**:
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

### 步骤 4: 重启 AI 助手 / Step 4: Restart AI Assistant

完全退出并重新启动 Claude Desktop（或您的 AI 助手）

Completely quit and restart Claude Desktop (or your AI assistant)

### 步骤 5: 测试 / Step 5: Test

在 AI 助手中输入 / Type in your AI assistant:

```
列出我的 Kubernetes 集群
或
List my Kubernetes clusters
```

---

## 支持的功能 / Supported Features

### 1. 列出集群 / List Clusters
```
Show me all my Kubernetes clusters
显示我所有的 Kubernetes 集群
```

### 2. 列出会话 / List Sessions
```
What kubemux sessions are currently active?
当前有哪些活动的 kubemux 会话？
```

### 3. 创建会话 / Create Session
```
Create a session for my production cluster
为我的生产集群创建一个会话
```

### 4. 附加到会话 / Attach to Session
```
How do I reconnect to the dev-workspace session?
如何重新连接到 dev-workspace 会话？
```

---

## 常见问题 / Common Issues

### 问题 1: MCP 服务器未找到 / MCP Server Not Found

**解决方案 / Solution**:
```bash
# 检查是否在 PATH 中 / Check if in PATH
which kubemux-mcp-server

# 如果未找到，使用完整路径 / If not found, use full path
{
  "command": "/usr/local/bin/kubemux-mcp-server"
}
```

### 问题 2: 权限错误 / Permission Error

**解决方案 / Solution**:
```bash
# 添加执行权限 / Add execute permission
chmod +x /usr/local/bin/kubemux-mcp-server
```

### 问题 3: AI 无法识别工具 / AI Doesn't Recognize Tools

**解决方案 / Solution**:
1. 确保配置文件格式正确 / Ensure config file is valid JSON
2. 完全重启 AI 助手 / Completely restart AI assistant
3. 检查日志 / Check logs

---

## 下一步 / Next Steps

- 📖 查看完整文档 / Read full documentation: [AI_PLUGIN_INTEGRATION.md](AI_PLUGIN_INTEGRATION.md)
- 💡 学习使用示例 / Learn usage examples: [mcp-usage-examples.md](examples/mcp-usage-examples.md)
- 🔧 了解技术细节 / Technical details: [mcp-server/README.md](../mcp-server/README.md)

---

## 架构图 / Architecture Diagram

```
┌─────────────────────────────────────────┐
│   AI 助手 / AI Assistant                │
│   (Claude, Copilot, Gemini)            │
└────────────────┬────────────────────────┘
                 │ JSON-RPC over stdio
                 │
┌────────────────▼────────────────────────┐
│   Kubemux MCP Server                    │
│   - list_clusters                       │
│   - list_sessions                       │
│   - create_session                      │
│   - attach_session                      │
└────────────────┬────────────────────────┘
                 │
┌────────────────▼────────────────────────┐
│   Kubemux Core                          │
│   - K8s cluster management              │
│   - Tmux/Zellij integration            │
└─────────────────────────────────────────┘
```

---

## 支持的 AI 平台 / Supported AI Platforms

| Platform | Status | Notes |
|----------|--------|-------|
| Claude Desktop | ✅ 完全支持 / Fully Supported | 推荐使用 / Recommended |
| GitHub Copilot | 🧪 实验性 / Experimental | 通过 MCP 扩展 / Via MCP extension |
| Google Gemini | 🧪 实验性 / Experimental | 需要适配器 / Requires adapter |
| 其他 MCP 兼容客户端 / Other MCP Clients | ✅ 支持 / Supported | 遵循标准协议 / Standard protocol |

---

## 获取帮助 / Getting Help

- 📧 提交 Issue / Submit Issue: https://github.com/corvofeng/kubemux/issues
- 📚 查看文档 / Documentation: https://kubemux.corvo.fun
- 💬 讨论区 / Discussions: https://github.com/corvofeng/kubemux/discussions

---

**提示 / Tip**: 将此文档添加到书签，方便快速参考！

**Tip**: Bookmark this document for quick reference!
