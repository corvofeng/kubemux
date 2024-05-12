# How to use kubemux


## kubeconfig

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


## tmuxinator

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

