# How to use kubemux


## tmuxinator

kubemux supports the tmuxinator configuration and template, which was the original inspiration behind kubemux.

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

For full docs, please refer to [Tmuxinator](./usage/tmuxinator.md)

## kubeconfig

As I used this project more, I found that I didn't need excessive customization for tmux, but simply to open the kubeconfig I desired. Therefore, I extended the project itself to better support kubeconfig configurations. I also added support for auto-completion, making it quicker to use your Kubernetes cluster now.

<script async src="https://asciinema.org/a/9lB50c5mndYfl0jBZLaG8ymdg.js" id="asciicast-658052" async="true"></script>

```bash
ls ~/.kube
# pve-kube.config xxx

kubemux kube --kube pve-kube.config

# I suggest you add the completion support
#   source <(kubemux completion bash)
#   source <(kubemux completion zsh)
# or you can add the command into the .bashrc or .zshrc.
kubemux kube --kube <tab>
```


For full docs, please refer to [kubeconfig](./usage/kubeconfig.md)

