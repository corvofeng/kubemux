# kubeconfig

<script async src="https://asciinema.org/a/9lB50c5mndYfl0jBZLaG8ymdg.js" id="asciicast-658052" async="true"></script>

As I used this project more, I found that I didn't need excessive customization for tmux, but simply to open the kubeconfig I desired. Therefore, I extended the project itself to better support kubeconfig configurations. I also added support for auto-completion, making it quicker to use your Kubernetes cluster now.


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


