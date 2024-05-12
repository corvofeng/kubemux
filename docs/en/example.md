# Best Practice


## kube-ps1 and fzf

> If you want to use kubectl with zsh, I recommend to use `kube-ps1`

https://github.com/jonmosco/kube-ps1

![kube-ps1](https://raw.githubusercontent.com/jonmosco/kube-ps1/master/img/kube-ps1.gif)

> A good search engine in the terminal

https://github.com/junegunn/fzf

```bash
git clone --depth 1 https://github.com/junegunn/fzf.git ~/.fzf && \
    ~/.fzf/install --bin  && ~/.fzf/install --completion --update-rc --key-bindings --no-bash --no-fish  && \
    grep  -q 'fzf.zsh' ~/.zshrc || echo '[ -f ~/.fzf.zsh ] && source ~/.fzf.zsh' >> ~/.zshrc
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
