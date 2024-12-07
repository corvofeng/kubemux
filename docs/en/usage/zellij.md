# zellij support


> News: 2024-12-07 I add the support for zellij, we can reuse the tmuxinator config for zellij now!

<script src="https://asciinema.org/a/693805.js" id="asciicast-693805" async="true"></script>


Please install zellij first, I still use the `bin` 

```bash
# Using bin: https://github.com/marcosnils/bin
bin install https://github.com/zellij-org/zellij
   • Getting latest release for zellij-org/zellij

Multiple matches found, please select one:

 [1] zellij-x86_64-unknown-linux-musl.sha256sum
 [2] zellij-x86_64-unknown-linux-musl.tar.gz
 Select an option: 2
   • Starting download of https://api.github.com/repos/zellij-org/zellij/releases/assets/207566049
12.95 MiB / 12.95 MiB [----------------------------------------] 100.00% 2.43 MiB p/s 5s
   • Copying for zellij@v0.41.2 into /home/corvo/.local/bin/zellij
   • Done installing zellij v0.41.2
```

```bash
kubemux -h
# --plexer string      Specify the plexer we want to use, [tmux|zellij]
```

You may set `--plexer` to use the zellij or tmux. 
However, if you don't set it. the `kubemux` will detect if we can use zellij as a default plexer.



