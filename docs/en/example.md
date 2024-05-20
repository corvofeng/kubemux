# Best Practice


## kube-ps1 

> If you want to use kubectl with zsh, I recommend to use `kube-ps1`

https://github.com/jonmosco/kube-ps1

![kube-ps1](https://raw.githubusercontent.com/jonmosco/kube-ps1/master/img/kube-ps1.gif)

```bash
# For zsh
plugins=(
  kube-ps1
)

PROMPT='$(kube_ps1)'$PROMPT # or RPROMPT='$(kube_ps1)'
```

## fzf

> A good search engine in the terminal

https://github.com/junegunn/fzf

```bash
git clone --depth 1 https://github.com/junegunn/fzf.git ~/.fzf && \
    ~/.fzf/install --bin  && ~/.fzf/install --completion --update-rc --key-bindings --no-bash --no-fish  && \
    grep  -q 'fzf.zsh' ~/.zshrc || echo '[ -f ~/.fzf.zsh ] && source ~/.fzf.zsh' >> ~/.zshrc
```

## k9s

https://k9scli.io/topics/commands/

k9s will use the kubemux environment, and will not impact other context.

<img width="1183" alt="image" src="https://github.com/corvofeng/kubemux/assets/12025071/36c4aa71-f30b-42dd-b487-5291d17166ff">

```
bin install https://github.com/derailed/k9s
   • Getting latest release for derailed/k9s

Multiple matches found, please select one:

 [1] k9s_Linux_amd64.tar.gz
 [2] k9s_Linux_amd64.tar.gz.sbom
 [3] k9s_linux_amd64.apk
 Select an option: 1
   • Starting download of https://api.github.com/repos/derailed/k9s/releases/assets/157727973
28.58 MiB / 28.58 MiB [------------------------------------------------------------------------] 100.00% 7.94 MiB p/s 4s

Multiple matches found, please select one:

 [1] LICENSE
 [2] README.md
 [3] k9s
 Select an option: 3
   • Copying for k9s@v0.32.4 into /home/corvo/.local/bin/k9s
   • Done installing k9s v0.32.4
```



## tmuxinator example template

```yaml
name: <%= @settings["project"] %>
root: ~/GitRepo

socket_name: <%= @settings["project"] %>
on_project_start:

  # Since a jump host is needed to connect to the API,
  # I used ssh to open a socks5 proxy on my local machine.
  # The following method is used to obtain a locally available port and the corresponding jump host:
  - export KUBECONFIG=~/.kube/config-<%= @settings["project"] %>
  - export TMUX_SSH_PORT="$(python3 -c 'import socket; s=socket.socket(); s.bind(("", 0)); print(s.getsockname()[1])')"
  - export TMUX_SSH_HOST="<%= @settings["host"] %>"
startup_window: kubectl

windows:
  - proxy:
      layout: main-vertical
      panes:
        - startup:
          - ls -alh
          - yq -i e '.current-context |= "tpl-<%= @settings["project"] %>"' $KUBECONFIG
          - yq -i e '.contexts[0].name |= "tpl-<%= @settings["project"] %>"' $KUBECONFIG
          - yq -i e ".clusters[0].cluster.proxy-url |=\"socks5://127.0.0.1:$TMUX_SSH_PORT\"" $KUBECONFIG
          - ssh -D $TMUX_SSH_PORT $TMUX_SSH_HOST
  - kubectl: ls # kubectl get pods
```

In addition to modifying the context, you can also add prompts to the terminal PS1, similar to this:

![image](https://github.com/corvofeng/kubemux/assets/12025071/e3a5b879-5af0-41ca-b2bf-91496ab8bcd8)
