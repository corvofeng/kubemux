### tmuxinator

<script src="https://asciinema.org/a/6kYCveJwVr4Sggj8QhqlsCKLm.js" id="asciicast-658053" async="true"></script>

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


