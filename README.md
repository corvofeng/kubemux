# gmux


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

./gmux -p gmux
```



