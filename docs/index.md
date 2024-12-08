# kubemux

kubemux 是一个用于管理多个集群的 tmux 封装工具，并集成了 [tmuxinator](https://github.com/tmuxinator/tmuxinator) 的功能。

(A tmux wrapper for managing multiple clusters and incorporating [tmuxinator](https://github.com/tmuxinator/tmuxinator).)

<script async src="https://asciinema.org/a/9lB50c5mndYfl0jBZLaG8ymdg.js" id="asciicast-658052" async="true"></script>

## 项目简介 (Project Introduction)

在开发这个项目时，我大量使用了 ChatGPT，主要用于编写逻辑和单元测试。原本预计需要 2-3 天的工作量被缩短到 1 天。最终的结果非常好，重写后的 kubemux 没有依赖，安装非常轻量。

(For this project, I heavily relied on ChatGPT, mainly for writing the logic and unit tests. Originally, I thought it would take 2-3 days of work, but the development cycle was shortened to 1 day. The final result is also very good. After rewriting, kubemux has no dependencies and is very lightweight to install.)

项目已在 GitHub 开源，并且与 tmuxinator 的配置兼容。如果您发现任何缺失的功能，欢迎提出 issue。

(The project is open source here, and it should be compatible with tmuxinator configurations. If you find any missing features, feel free to open an issue.)

## 源代码 (Source Code)

https://github.com/corvofeng/kubemux

## 主要功能 (Main Features)

- 多集群管理 (Multiple cluster management)
- tmuxinator 配置兼容 (tmuxinator configuration compatibility)
- 轻量级安装 (Lightweight installation)
- 命令行自动补全 (Command-line auto-completion)
- 支持 zellij 终端复用器 (Support for zellij terminal multiplexer)

## 文档导航 (Documentation Navigation)

| 英文文档 (English Documents) | 中文文档 (Chinese Documents) |
| ---- | ---- |
| [Introduction](en/intro.md) | [简介](zh/intro.md) |
| [Why](en/why.md) | [原因](zh/why.md) |
| [Example](en/example.md) | [示例](zh/example.md) |
| [Usage](en/usage.md) | [使用说明](zh/usage.md) |
| [Usage - Tmuxinator](en/usage/tmuxinator.md) | [使用说明 - Tmuxinator](zh/usage/tmuxinator.md) |
| [Usage - Kubeconfig](en/usage/kubeconfig.md) | [使用说明 - Kubeconfig](zh/usage/kubeconfig.md) |
| [Usage - Zellij](en/usage/zellij.md) | [使用说明 - Zellij](zh/usage/zellij.md) |

## 快速开始 (Quick Start)

```bash
# MacOS 安装 (MacOS Installation)
brew install corvofeng/tap/kubemux

# Linux 安装 (Linux Installation)
bin install https://github.com/corvofeng/kubemux ~/usr/bin

# 启用命令补全 (Enable Command Completion)
source <(kubemux completion bash)  # Bash
source <(kubemux completion zsh)   # Zsh
```

## 贡献 (Contributing)

欢迎提交 Pull Requests 和 Issues！

(Pull Requests and Issues are welcome!)
