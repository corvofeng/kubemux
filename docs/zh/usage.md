# 如何使用 kubemux

## tmuxinator

kubemux 支持 tmuxinator 的配置和模板功能，这也是 kubemux 最初的灵感来源。

<script src="https://asciinema.org/a/6kYCveJwVr4Sggj8QhqlsCKLm.js" id="asciicast-658053" async="true"></script>

```bash
mkdir ~/.tmuxinator

echo '
name: kubemux
root: "~/"
windows:
  - p1:
    - ls
    - pwd
  - p2:
    - pwd
    - echo "hello world"
  - p3: htop
' > ~/.tmuxinator/kubemux.yml

kubemux -p kubemux
```

完整文档请参考 [Tmuxinator](./usage/tmuxinator.md)

## kubeconfig

随着项目的使用，我发现并不需要对 tmux 进行过多的自定义，而是只需要打开我想要的 kubeconfig。因此，我扩展了项目本身以更好地支持 kubeconfig 配置。我还添加了自动补全支持，现在使用 Kubernetes 集群变得更加快捷。

<script async src="https://asciinema.org/a/9lB50c5mndYfl0jBZLaG8ymdg.js" id="asciicast-658052" async="true"></script>

```bash
ls ~/.kube
# pve-kube.config xxx

kubemux kube --kube pve-kube.config

# 建议添加补全支持
#   source <(kubemux completion bash)
#   source <(kubemux completion zsh)
# 或者您可以将命令添加到 .bashrc 或 .zshrc 中。
kubemux kube --kube <tab>
```

完整文档请参考 [kubeconfig](./usage/kubeconfig.md) 