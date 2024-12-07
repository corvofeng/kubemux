# 为什么创建这个工具

## Kubernetes 多集群管理方案

Kubernetes 官方网站提供了一个使用 KUBECONFIG 中的 context 来切换集群的解决方案。

https://kubernetes.io/docs/tasks/access-application-cluster/configure-access-multiple-clusters/

```yaml
apiVersion: v1
kind: Config
preferences: {}

clusters:
- cluster:
  name: development
- cluster:
  name: test

users:
- name: developer
- name: experimenter

contexts:
- context:
  name: dev-frontend
- context:
  name: dev-storage
- context:
  name: exp-test
```

许多现有的多 Kubernetes 集群管理工具（如 kubecm）都依赖于维护单个 KUBECONFIG 文件。这种方法有几个缺点：

1. KUBECONFIG 文件维护：添加和删除集群需要手动编辑 KUBECONFIG 文件，这在管理大量集群时可能会变得繁琐且容易出错。
2. 单集群操作：由于所有集群都配置在同一个 KUBECONFIG 文件中，一次只能操作一个集群。这在多集群管理环境中可能会不便且效率低下。

## 使用 tmux 配置文件分割的解决方案

为了解决这些问题，我提出了一个使用 tmux 配置文件分割的解决方案。这种方法为每个 Kubernetes 集群创建一个单独的 tmux 会话。每个 tmux 会话都有自己的 KUBECONFIG 文件，这样您就可以同时独立操作多个集群。

## tmux

`tmux` 是一个强大的终端复用器，允许您在同一个终端窗口中创建和管理多个会话。这对于管理多个服务器或集群非常有用，因为您可以在不打开多个终端窗口的情况下轻松切换不同的会话。

要使用 tmux，您需要先安装它。在大多数 Linux 发行版中，您可以使用以下命令安装 tmux：

```bash
sudo apt install tmux
```

安装完成后，您可以使用以下命令启动 tmux：

```bash
tmux
```

这将在您的终端窗口中创建一个新的 tmux 会话。您可以使用以下命令在不同的会话之间切换：

```bash
tmux attach-session -t <session-name>
tmux new-session -s <session-name>
```

### 多会话

简单来说，`-L socket-name` 参���允许您指定 tmux socket 的位置，不同的 socket 对应完全隔离的会话。

我们可以在不同的会话中使用不同的环境变量来实现环境隔离。

例如，通过以下两个命令，您可以创建两个完全独立的、拥有自己环境变量的终端：

```bash
KUBECONFIG=~/.kube/config-aa tmux -L aa
KUBECONFIG=~/.kube/config-bb tmux -L bb
```

这里的脚本已经可以实现多集群管理。那么为什么我们还要引入 tmuxinator 和我的新工具 kubemux 呢？

1. 生产环境的配置需要通过跳板机，如何在本地使用 KUBECONFIG？（我在最后附上了说明）
2. 我希望 tmux 会话有多个具有各自功能的窗口。

## tmuxinator

https://github.com/tmuxinator/tmuxinator

![](https://user-images.githubusercontent.com/289949/44366875-1a6cee00-a49c-11e8-9322-76e70df0c88b.gif)

这是一个用 Ruby 编写的工具，允许您以 YAML 格式定义 tmux 终端。它还支持模板功能。以下是一个模板化 YAML 文件的示例：

### 使用方法

```yaml
name: project
root: ~/<%= @settings["workspace"] %>
# tmuxinator start project workspace=~/workspace/todo

windows:
  - small_project:
      root: ~/projects/company/small_project
      panes:
        - start this
        - start that
```

在管理集群环境迁移的几个月里，一些���常用的命令是：

```bash
tmuxinator tpl project=ingame-pre-na
tmuxinator tpl project=ingame-pre-sg
tmuxinator tpl project=ingame-pre-fra
```

它可以帮助我完美地区分不同的环境，而且因为我使用 fzf，我甚至可以模糊搜索我想要打开的环境。

![image](https://github.com/corvofeng/kubemux/assets/12025071/36c8a6ed-47e9-49cf-8a99-1389899b0091)

### 局限性

由于它是用 Ruby 编写的，需要在机器上安装相对较新版本的 Ruby。AWS 需要登录跳板机进行操作，但我们使用的机器都很旧，我不想编译和重新安装 Ruby。
快速查看代码后，我发现它并没有使用任何高级特性。用 Golang 完全重写它会很容易。