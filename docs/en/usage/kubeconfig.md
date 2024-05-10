### kubeconfig

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

<script async src="https://asciinema.org/a/9lB50c5mndYfl0jBZLaG8ymdg.js" id="asciicast-658052" async="true"></script>

