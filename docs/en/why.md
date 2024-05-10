
## Kubernetes Multi-Cluster Management Solution

The official Kubernetes website provides a solution for switching between clusters using the context in KUBECONFIG.

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

Many existing tools for managing multiple Kubernetes clusters, such as kubecm, rely on maintaining a single KUBECONFIG file. This approach has several drawbacks:

1. KUBECONFIG file maintenance: Adding and deleting clusters requires manually editing the KUBECONFIG file, which can become tedious and error-prone, especially when managing a large number of clusters.
2. ingle-cluster operation: Since all clusters are configured in the same KUBECONFIG file, only one cluster can be operated on at a time. This can be inconvenient and inefficient in a multi-cluster management environment.

Solution using tmux configuration file splitting:

To address these problems, I propose a solution that utilizes tmux configuration file splitting. This approach involves creating a separate tmux session for each Kubernetes cluster. Each tmux session has its own KUBECONFIG file, which allows you to operate on multiple clusters simultaneously and independently.


## tmux Multiple Sessions

`tmux` is a powerful terminal multiplexer that allows you to create and manage multiple sessions in the same terminal window. This is very useful for managing multiple servers or clusters because you can easily switch between different sessions without opening multiple terminal windows.


To use tmux, you need to install it first. In most Linux distributions, you can install tmux using the following command:

```
sudo apt install tmux
```

After the installation is complete, you can start tmux using the following command:

```
tmux
```

This will create a new tmux session in your terminal window. You can use the following command to switch between different sessions:

```
tmux attach-session -t <session-name>
tmux new-session -s <session-name>
```

### tmux -L

In simple terms, the `-L socket-name` parameter allows you to specify the location of the tmux socket, and different sockets correspond to completely isolated sessions.

We can use different environment variables in different sessions to achieve environment separation.

For example, with the following two commands, you can create two completely independent terminals with their own environment variables:

```
KUBECONFIG=~/.kube/config-aa tmux -L aa
KUBECONFIG=~/.kube/config-bb tmux -L bb
```

The script here can already achieve multi-cluster management. So why do we introduce tmuxinator and my new tool kubemux?

1. The configuration for the production environment needs to go through a jump host. How can we use KUBECONFIG locally? (I have attached it at the end)
2. I want the tmux session to have multiple windows with their own functionalities.


## tmuxinator Usage and Limitations

https://github.com/tmuxinator/tmuxinator

![](https://user-images.githubusercontent.com/289949/44366875-1a6cee00-a49c-11e8-9322-76e70df0c88b.gif)

It is a tool written in Ruby that allows you to define tmux terminals in YAML format. It also supports templating. Here is an example of a templated YAML file:

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

During the months of managing cluster environment migration, some of the most frequently used commands were:

```bash
tmuxinator tpl project=ingame-pre-na
tmuxinator tpl project=ingame-pre-sg
tmuxinator tpl project=ingame-pre-fra
```

It can help me perfectly differentiate different environments, and because I use fzf, I can even fuzzy search for the environment I want to open.

![image](https://github.com/corvofeng/kubemux/assets/12025071/36c8a6ed-47e9-49cf-8a99-1389899b0091)


Limitations: Since it is written in Ruby, it requires a relatively new version of Ruby installed on the machine. AWS requires logging into a jump host for operations, but the machines we use are very old, and I don't want to compile and reinstall Ruby.
After a quick look at the code, I found that it doesn't use any advanced features. It would be very easy to completely rewrite it in Golang.

## Implementation of kubemux

For this project, I heavily relied on ChatGPT, mainly for writing the logic and unit tests. Originally, I thought it would take 2-3 days of work, but the development cycle was shortened to 1 day. The final result is also very good. After rewriting, kubemux has no dependencies and is very lightweight to install.

The project is open source here, and it should be compatible with tmuxinator configurations. If you find any missing features, feel free to open an issue.

https://github.com/corvofeng/kubemux
