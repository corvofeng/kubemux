# kubeconfig

<script async src="https://asciinema.org/a/9lB50c5mndYfl0jBZLaG8ymdg.js" id="asciicast-658052" async="true"></script>

## 基本用法

kubemux 提供了一种简单的方式来管理多个 kubeconfig 文件。您可以使用 `kubemux kube` 命令来启动一个新的会话，该会话将使用指定的 kubeconfig 文件。

## 命令行选项

- `--kube`: 指定要使用的 kubeconfig 文件
- `--plexer`: 选择使用的终端复用器（tmux 或 zellij）

## 示例

```bash
# 使用特定的 kubeconfig 文件
kubemux kube --kube ~/.kube/config-dev

# 使用 zellij 作为终端复用器
kubemux kube --kube ~/.kube/config-dev --plexer zellij
```

## 自动补全

为了更方便地使用 kubemux，建议启用命令行自动补全功能。您可以将以下命令添加到您的 shell 配置文件中：

```bash
# 对于 bash 用户（添加到 ~/.bashrc）
source <(kubemux completion bash)

# 对于 zsh 用户（添加到 ~/.zshrc）
source <(kubemux completion zsh)
``` 