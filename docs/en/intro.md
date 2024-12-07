
## Installation

### MacOS

```bash
brew install corvofeng/tap/kubemux
```

### Linux

> Using bin: https://github.com/marcosnils/bin

```bash
bin install https://github.com/corvofeng/kubemux ~/usr/bin
# bin ls
# Path                  Version  URL                                                       Status
# ~/usr/bin/kubemux     v1.1.2   https://github.com/corvofeng/kubemux/releases/tag/v1.1.2  OK
```

> Using binary

```bash
cd /tmp
rm kubemux_linux_amd64.tar.gz
wget https://github.com/corvofeng/kubemux/releases/latest/download/kubemux_linux_amd64.tar.gz
tar -zxvf kubemux_linux_amd64.tar.gz
sudo install -v kubemux /usr/local/bin
```

### Completion

To enable shell completion, you can use the following commands:

For Bash:
```bash
source <(kubemux completion bash)
```

For Zsh:
```bash
source <(kubemux completion zsh)
```
