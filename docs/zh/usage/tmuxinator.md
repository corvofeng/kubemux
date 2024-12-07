# tmuxinator


<script src="https://asciinema.org/a/6kYCveJwVr4Sggj8QhqlsCKLm.js" id="asciicast-658053" async="true"></script>

## 基本用法

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

## ERB 模板支持

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

> 您可以使用 `-p example-tpl` 来指定我们使用的模板文件
> `--set project=hello` 表示我们将把 `<%= @settings["project"] %>` 替换为 `hello
`
