# kubemux

![blog-test](https://github.com/corvofeng/kubemux/assets/12025071/375541b7-927f-485d-bd75-36edc39bbae2)

A terminal multiplexer wrapper designed for Kubernetes multi-cluster management, with support for both tmux and zellij.

## Features

- **Multi-Cluster Management**: Easily switch between different Kubernetes clusters
- **Terminal Multiplexer Support**: Works with both tmux and zellij
- **tmuxinator Compatible**: Supports existing tmuxinator configurations
- **Zero Dependencies**: Standalone binary with no external dependencies
- **Shell Completion**: Built-in completion for bash and zsh
- **Jump Host Support**: Seamlessly work with clusters behind jump hosts

## Installation

### MacOS

bash
```
brew install corvofeng/tap/kubemux
```

### Linux

Using [https://github.com/marcosnils/bin](https://github.com/marcosnils/bin):

```bash
bin install https://github.com/corvofeng/kubemux ~/usr/bin
```

Using binary:
```bash
cd /tmp
wget https://github.com/corvofeng/kubemux/releases/latest/download/kubemux_linux_amd64.tar.gz
tar -zxvf kubemux_linux_amd64.tar.gz
sudo install -v kubemux /usr/local/bin
```

## Quick Start

### Kubeconfig Management
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

### tmuxinator Configuration

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

## Documentation

Full documentation is available at: https://kubemux.corvo.fun

Blog posts:
- 英文版: https://corvo.fun/2023/12/27/2023-12-26-kubemux%E7%9A%84%E5%BC%80%E5%8F%91%E4%B8%8E%E4%BD%BF%E7%94%A8/
- 中文版: https://corvo.myseu.cn/2023/12/27/2023-12-26-kubemux%E7%9A%84%E5%BC%80%E5%8F%91%E4%B8%8E%E4%BD%BF%E7%94%A8/

## Demo

### Kubeconfig Management
[![asciicast](https://asciinema.org/a/9lB50c5mndYfl0jBZLaG8ymdg.svg)](https://asciinema.org/a/9lB50c5mndYfl0jBZLaG8ymdg)

### tmuxinator Support
[![asciicast](https://asciinema.org/a/6kYCveJwVr4Sggj8QhqlsCKLm.svg)](https://asciinema.org/a/6kYCveJwVr4Sggj8QhqlsCKLm)

### zellij Support
[![asciicast](https://asciinema.org/a/693805.svg)](https://asciinema.org/a/693805)


## Contributing

Pull requests are welcome! For major changes, please open an issue first to discuss what you would like to change.

## License

[MIT](LICENSE)