# gmux


![629109](https://github.com/corvofeng/gmux/assets/12025071/09293818-40d8-473e-8e6a-aa7b2a790a97)


An alternative for tmuxinator

```
mkdir ~/.tmuxinator

echo '
name: gmux
root: "~/"
windows:
  - p1:
    - ls
    - pwd
  - p2:
    - pwd
    - echo "hello world"
  - p3: htop
' > ~/.tmuxinator/gmux.yml

cd /tmp
wget https://github.com/corvofeng/gmux/releases/download/v0.3.2/gmux_linux_amd64.tar.gz
tar -zxvf gmux_linux_amd64.tar.gz
sudo install -v gmux /usr/local/bin

gmux -p gmux
```

Here is an example:

[![asciicast](https://asciinema.org/a/lVIIOwzWwFAL611IwUeZpohoy.svg)](https://asciinema.org/a/lVIIOwzWwFAL611IwUeZpohoy)

