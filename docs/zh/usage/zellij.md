# zellij 支持

> 新闻：2024-12-07 我添加了对 zellij 的支持，现在我们可以在 zellij 中重用 tmuxinator 配置了！

<script src="https://asciinema.org/a/693805.js" id="asciicast-693805" async="true"></script>

请先安装 zellij，我仍然使用 `bin`

```bash
# 使用 bin：https://github.com/marcosnils/bin
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
# --plexer string      指定我们要使用的终端复用器，[tmux|zellij]
```

您可以设置 `--plexer` 来使用 zellij 或 tmux。
但是，如果您不设置它，`kubemux` 会检测是否可以使用 zellij 作为默认的终端复用器。 