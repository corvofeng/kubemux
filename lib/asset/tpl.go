package asset

const BashScriptTemplate = `#!/bin/bash
set -ex
# 创建 tmux 会话
{{$.Tmux}} start-server;

{{if not (TmuxHasSession $.Name)}}
{{- range $i, $window := $.Windows }}
#============================================================
# Window: {{$window.Name}}
{{$winId := inc $i}}

{{if eq $i 0}}
{{$.Tmux}} new-session -d -s {{$.Name}} -n {{$window.Name}}
{{ else }}
{{$.Tmux}} new-window -c {{$window.Root}} -t {{$.Name}}:{{$winId}} -n {{$window.Name}}
{{- end }}
# 设置窗口的根目录
# {{- if $window.Root}}
# {{$.Tmux}} send-keys -t {{$.Name}}:{{$i}} "cd {{$window.Root}}" C-m
# {{- end}}

# {{$.Tmux}} select-layout -t {{$.Name}}:{{$winId}} tiled
{{$.Tmux}} select-layout -t {{$.Name}}:{{$winId}} main-vertical


{{- range $j, $pane := $window.Panes}}
{{$panelId := inc $j}}

{{$.Tmux}} list-panes -a
# 分割窗口并运行命令
{{$.Tmux}} split-window -c {{$window.Root}} -t {{$.Name}}:{{$winId}}
{{$.Tmux}} send-keys -t {{$.Name}}:{{$winId}}.{{$panelId}} "{{$pane.Command}}" C-m
{{- end}}

# 关闭最后一个多余的pane
{{$.Tmux}} kill-pane -t {{$.Name}}:{{$winId}}.0
#================================================================
{{- end}}
{{ else }}
# Already have a session
{{- end }}


# 附加到 tmux 会话
#{{$.Tmux}} attach-session -t {{.Name}}
`
