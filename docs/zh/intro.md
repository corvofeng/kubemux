# 中文文档

## 安装

### MacOS

```bash
brew install corvofeng/tap/kubemux
```

### Linux

> 使用 bin 安装: https://github.com/marcosnils/bin

```bash
bin install https://github.com/corvofeng/kubemux ~/usr/bin
# bin ls
# Path                  Version  URL                                                       Status
# ~/usr/bin/kubemux     v1.1.2   https://github.com/corvofeng/kubemux/releases/tag/v1.1.2  OK
```

> 使用二进制文件安装

```bash
cd /tmp
rm kubemux_linux_amd64.tar.gz
wget https://github.com/corvofeng/kubemux/releases/latest/download/kubemux_linux_amd64.tar.gz
tar -zxvf kubemux_linux_amd64.tar.gz
sudo install -v kubemux /usr/local/bin
```

### 命令补全

要启用shell补全功能，可以使用以下命令：

对于 Bash：
bash
source <(kubemux completion bash)
```

对于 Zsh：
```bash
source <(kubemux completion zsh)
```

## 项目简介

kubemux 是一个用于管理多个 Kubernetes 集群的 tmux 封装工具，并集成了 [tmuxinator](https://github.com/tmuxinator/tmuxinator) 的功能。

本项目在开发过程中大量使用了 ChatGPT，主要用于编写逻辑和单元测试。原本预计需要 2-3 天的工作量被缩短到 1 天。最终的结果非常好，重写后的 kubemux 没有依赖，安装非常轻量。

项目已在 GitHub 开源，并且应该与 tmuxinator 的配置兼容。如果您发现任何缺失的功能，欢迎提出 issue。

### 主要功能

- 支持多集群管理
- 兼容 tmuxinator 配置
- 支持 kubeconfig 配置
- 提供命令行补全
- 支持 zellij 终端复用器

### 文档导航

- [使用说明](usage.md)
- [最佳实践](example.md)
- [为什么创建这个工具](why.md)
```
