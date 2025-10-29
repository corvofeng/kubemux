# Kubemux AI Plugin Integration Guide

## 将 Kubemux 作为 AI 助手插件使用指南

本指南详细说明如何将 kubemux 集成到各种 AI 助手中，使其能够作为插件/扩展使用。

---

## 目录

1. [架构概述](#架构概述)
2. [支持的 AI 平台](#支持的-ai-平台)
3. [快速开始](#快速开始)
4. [详细集成步骤](#详细集成步骤)
5. [使用场景](#使用场景)
6. [高级配置](#高级配置)
7. [故障排除](#故障排除)

---

## 架构概述

### MCP 协议架构

Kubemux 通过实现 Model Context Protocol (MCP) 来支持 AI 助手集成：

```
┌──────────────────────────────────────────────────────────┐
│                     AI 助手层                              │
│  Claude Desktop / GitHub Copilot / Google Gemini / etc.  │
└──────────────────────┬───────────────────────────────────┘
                       │
                       │ JSON-RPC 2.0 over stdio
                       │
┌──────────────────────▼───────────────────────────────────┐
│                 Kubemux MCP Server                        │
│  - 工具注册与管理                                           │
│  - 请求解析与路由                                           │
│  - 响应格式化                                              │
└──────────────────────┬───────────────────────────────────┘
                       │
        ┌──────────────┼──────────────┐
        │              │              │
┌───────▼──────┐ ┌────▼─────┐ ┌─────▼──────┐
│ 集群工具      │ │ 会话工具  │ │ 配置工具    │
│ - 列出集群    │ │ - 创建会话│ │ - 读取配置  │
│ - 切换上下文  │ │ - 附加会话│ │ - 验证配置  │
└───────┬──────┘ └────┬─────┘ └─────┬──────┘
        │              │              │
        └──────────────┼──────────────┘
                       │
┌──────────────────────▼───────────────────────────────────┐
│                 Kubemux 核心库                            │
│  - Kubernetes 集群管理                                    │
│  - Tmux/Zellij 会话管理                                  │
│  - 配置解析与渲染                                         │
└──────────────────────────────────────────────────────────┘
```

### 关键组件

1. **MCP Server** (`mcp-server/main.go`)
   - 实现 JSON-RPC 2.0 协议
   - 处理 `initialize`, `tools/list`, `tools/call` 等方法
   - 通过 stdin/stdout 与 AI 助手通信

2. **Tool Layer** (`mcp-server/tools/`)
   - 定义可供 AI 调用的工具
   - 每个工具包含 schema 和 execute 方法
   - 工具类型：cluster、session、attach

3. **Kubemux Library** (`lib/`)
   - 提供核心功能实现
   - 处理实际的 Kubernetes 和终端复用器操作

---

## 支持的 AI 平台

### ✅ 完全支持

1. **Claude Desktop** (Anthropic)
   - 原生 MCP 支持
   - 最佳集成体验
   - 推荐用于生产环境

2. **MCP Compatible Clients**
   - 任何支持 MCP 协议的客户端
   - 标准 JSON-RPC 2.0 接口

### 🚧 实验性支持

3. **GitHub Copilot** (Microsoft)
   - 通过 VS Code MCP 扩展
   - 需要额外配置

4. **Google Gemini**
   - 通过 API 集成
   - 可能需要自定义适配器

### 📋 计划支持

- OpenAI ChatGPT 插件
- JetBrains AI Assistant
- Cursor IDE

---

## 快速开始

### 前置要求

- Go 1.22 或更高版本
- Kubemux 已安装
- Tmux 或 Zellij（可选）
- 目标 AI 助手已安装

### 5 分钟快速设置

```bash
# 1. 克隆或进入 kubemux 项目目录
cd /path/to/kubemux

# 2. 构建 MCP 服务器
make build-mcp

# 3. 安装到系统路径
sudo cp kubemux-mcp-server /usr/local/bin/
sudo chmod +x /usr/local/bin/kubemux-mcp-server

# 4. 验证安装
which kubemux-mcp-server
kubemux-mcp-server --help 2>&1 | head -1
```

### 快速测试

```bash
# 测试初始化
echo '{"jsonrpc":"2.0","id":1,"method":"initialize","params":{}}' | kubemux-mcp-server

# 测试工具列表
echo '{"jsonrpc":"2.0","id":2,"method":"tools/list"}' | kubemux-mcp-server
```

---

## 详细集成步骤

### 1. Claude Desktop 集成

#### 步骤 1: 找到配置文件

**macOS**:
```bash
open ~/Library/Application\ Support/Claude/
# 编辑 claude_desktop_config.json
```

**Linux**:
```bash
mkdir -p ~/.config/Claude
# 编辑 ~/.config/Claude/claude_desktop_config.json
```

**Windows**:
```powershell
explorer %APPDATA%\Claude
# 编辑 claude_desktop_config.json
```

#### 步骤 2: 添加配置

```json
{
  "mcpServers": {
    "kubemux": {
      "command": "kubemux-mcp-server",
      "args": [],
      "env": {
        "HOME": "/Users/yourname"
      }
    }
  }
}
```

#### 步骤 3: 重启 Claude Desktop

完全退出并重新启动 Claude Desktop 应用。

#### 步骤 4: 验证集成

在 Claude 中输入：
```
请列出我的 Kubernetes 集群
```

Claude 应该能够调用 `list_kubeconfigs` 工具并返回结果。

### 2. VS Code + GitHub Copilot 集成

#### 步骤 1: 安装 MCP 扩展

```bash
# 查找并安装 MCP 相关扩展
# 注意：这是实验性功能，可能需要特定版本的 VS Code
```

#### 步骤 2: 配置 settings.json

打开 VS Code 设置 (`Ctrl+,` 或 `Cmd+,`)，编辑 `settings.json`:

```json
{
  "mcp.servers": {
    "kubemux": {
      "command": "kubemux-mcp-server",
      "type": "stdio"
    }
  }
}
```

#### 步骤 3: 重载窗口

按 `Ctrl+Shift+P` (或 `Cmd+Shift+P`)，运行 "Developer: Reload Window"

### 3. 自定义 AI 工具集成

对于其他 AI 工具，您可以使用标准的 MCP 客户端库：

#### Python 示例

```python
import subprocess
import json

class KubemuxMCPClient:
    def __init__(self):
        self.process = subprocess.Popen(
            ['kubemux-mcp-server'],
            stdin=subprocess.PIPE,
            stdout=subprocess.PIPE,
            stderr=subprocess.PIPE
        )
        self.request_id = 0
    
    def call_tool(self, tool_name, arguments=None):
        self.request_id += 1
        request = {
            "jsonrpc": "2.0",
            "id": self.request_id,
            "method": "tools/call",
            "params": {
                "name": tool_name,
                "arguments": arguments or {}
            }
        }
        
        self.process.stdin.write(json.dumps(request).encode() + b'\n')
        self.process.stdin.flush()
        
        response = json.loads(self.process.stdout.readline())
        return response
    
    def list_tools(self):
        self.request_id += 1
        request = {
            "jsonrpc": "2.0",
            "id": self.request_id,
            "method": "tools/list"
        }
        
        self.process.stdin.write(json.dumps(request).encode() + b'\n')
        self.process.stdin.flush()
        
        response = json.loads(self.process.stdout.readline())
        return response

# 使用示例
client = KubemuxMCPClient()
tools = client.list_tools()
print("Available tools:", tools)

clusters = client.call_tool("list_kubeconfigs")
print("Clusters:", clusters)
```

---

## 使用场景

### 场景 1: 日常集群管理

**对话示例**:

```
用户: "我有哪些 Kubernetes 集群？"
AI: 正在查询您的集群...
    [调用 list_kubeconfigs 工具]
    您有以下集群：
    - prod-cluster.config
    - dev-cluster.config
    - staging-cluster.config

用户: "帮我连接到生产集群"
AI: [调用 create_session 工具]
    好的，请运行以下命令连接到生产集群：
    kubemux kube --kube prod-cluster.config
```

### 场景 2: 多集群会话管理

**对话示例**:

```
用户: "创建三个会话，分别连接到我的三个集群"
AI: [调用 create_session 工具三次]
    好的，我为您生成了三个命令：
    
    1. 生产环境:
       kubemux kube --kube prod-cluster.config --plexer tmux
    
    2. 开发环境:
       kubemux kube --kube dev-cluster.config --plexer tmux
    
    3. 测试环境:
       kubemux kube --kube staging-cluster.config --plexer tmux
    
    您可以在不同的终端窗口中运行这些命令。

用户: "显示当前活动的会话"
AI: [调用 list_sessions 工具]
    当前有 2 个活动会话：
    - prod-cluster
    - dev-cluster
```

### 场景 3: 故障排查

**对话示例**:

```
用户: "我想检查生产环境的问题，但忘记了会话名称"
AI: [调用 list_sessions 工具]
    您有以下活动会话：
    - prod-cluster (tmux)
    - dev-env (zellij)
    
    看起来 "prod-cluster" 就是生产环境的会话。

用户: "怎么连接到这个会话？"
AI: [调用 attach_session 工具]
    运行以下命令即可连接：
    tmux -L prod-cluster attach-session -t prod-cluster
```

---

## 高级配置

### 环境变量配置

在 MCP 配置中可以设置环境变量：

```json
{
  "mcpServers": {
    "kubemux": {
      "command": "kubemux-mcp-server",
      "env": {
        "KUBECONFIG": "/custom/path/to/kubeconfig",
        "KUBEMUX_DEBUG": "true",
        "TMUX_TMPDIR": "/tmp/tmux-custom"
      }
    }
  }
}
```

### 自定义工具开发

添加新工具的步骤：

1. 在 `mcp-server/tools/` 创建新文件
2. 实现 `Tool` 接口
3. 在 `server.go` 中注册

示例：

```go
// mcp-server/tools/custom.go
package tools

type CustomTool struct {
    BaseTool
}

func NewCustomTool() *CustomTool {
    return &CustomTool{
        BaseTool: BaseTool{
            name:        "custom_operation",
            description: "Execute custom operation",
            inputSchema: map[string]interface{}{
                "type": "object",
                "properties": map[string]interface{}{
                    "param": map[string]interface{}{
                        "type": "string",
                        "description": "Custom parameter",
                    },
                },
            },
        },
    }
}

func (t *CustomTool) Execute(arguments map[string]interface{}) (string, error) {
    // 实现您的逻辑
    return "result", nil
}
```

### 日志和调试

启用调试日志：

```bash
# 设置环境变量
export KUBEMUX_MCP_DEBUG=1

# 或在 AI 配置中添加
{
  "mcpServers": {
    "kubemux": {
      "command": "bash",
      "args": ["-c", "KUBEMUX_MCP_DEBUG=1 kubemux-mcp-server 2>/tmp/kubemux-mcp.log"]
    }
  }
}

# 查看日志
tail -f /tmp/kubemux-mcp.log
```

---

## 故障排除

### 常见问题

#### 1. AI 助手找不到 kubemux-mcp-server

**症状**: 配置后 AI 无法识别 kubemux 工具

**解决方案**:
```bash
# 检查二进制文件位置
which kubemux-mcp-server

# 如果未找到，添加到 PATH
export PATH=$PATH:/path/to/kubemux-mcp-server

# 或使用完整路径
{
  "command": "/usr/local/bin/kubemux-mcp-server"
}
```

#### 2. 权限问题

**症状**: "Permission denied" 错误

**解决方案**:
```bash
# 添加执行权限
chmod +x /usr/local/bin/kubemux-mcp-server

# 检查文件权限
ls -la /usr/local/bin/kubemux-mcp-server
```

#### 3. JSON 解析错误

**症状**: AI 无法解析 MCP 响应

**解决方案**:
```bash
# 手动测试 JSON 输出
echo '{"jsonrpc":"2.0","id":1,"method":"tools/list"}' | \
  kubemux-mcp-server | jq .

# 应该返回有效的 JSON
```

#### 4. Kubeconfig 未找到

**症状**: list_kubeconfigs 返回空列表

**解决方案**:
```bash
# 确保 kubeconfig 文件存在
ls -la ~/.kube/

# 设置正确的 HOME 环境变量
{
  "env": {
    "HOME": "/Users/yourname"
  }
}
```

### 调试步骤

1. **验证安装**
   ```bash
   kubemux-mcp-server --version
   ```

2. **测试基本功能**
   ```bash
   echo '{"jsonrpc":"2.0","id":1,"method":"initialize","params":{}}' | \
     kubemux-mcp-server
   ```

3. **检查日志**
   - Claude Desktop: 查看应用日志
   - VS Code: 打开输出面板，选择 MCP

4. **验证环境**
   ```bash
   # 检查必需的依赖
   which tmux
   which zellij
   ls ~/.kube/
   ```

### 获取帮助

如果问题仍未解决：

1. 查看 GitHub Issues: https://github.com/corvofeng/kubemux/issues
2. 提交新 Issue，包含：
   - 操作系统和版本
   - AI 助手类型和版本
   - 错误信息和日志
   - MCP 配置文件

---

## 最佳实践

### 1. 安全性

- 不要在 MCP 配置中暴露敏感信息
- 使用环境变量管理凭证
- 定期更新 kubemux-mcp-server

### 2. 性能

- 使用 `list_sessions` 前先检查是否需要
- 批量操作时考虑使用脚本而非逐个调用工具
- 定期清理不再使用的会话

### 3. 维护

- 记录您的 MCP 配置
- 版本控制自定义工具代码
- 测试工具更新后再部署

---

## 下一步

- 查看 [MCP Server README](README.md) 了解更多技术细节
- 浏览 [示例用例](../docs/examples/mcp-examples.md)（如果可用）
- 参与社区讨论和贡献

---

## 贡献

欢迎为 MCP 集成做出贡献：

1. Fork 项目
2. 创建功能分支
3. 提交您的更改
4. 推送到分支
5. 创建 Pull Request

---

## 许可证

本项目采用 MIT 许可证 - 详见 [LICENSE](../LICENSE) 文件
