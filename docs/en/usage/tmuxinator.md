# tmuxinator

## Basic Usage

```
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

## ERB template support

```bash
mkdir -pv ~/.tmuxinator

echo '
name: <%= @settings["project"] %>
root: ~/

on_project_start:
  - export KUBECONFIG=~/.kube/config-<%= @settings["project"] %>

windows:
  - kubectl: ls
  - env: echo $KUBECONFIG
' > ~/.tmuxinator/example-tpl.yml

kubemux -p example-tpl --set project=hello
```

> You can use `-p example-tpl` to specify the template file we use
> The `--set project=hello` means we will replace `<%= @settings["project"] %>` with `hello`



