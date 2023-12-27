package asset

const BashScriptTemplate = `#!/bin/bash

{{ if $.Debug }}
set -ex
{{ end }}

{{$.Tmux}} start-server;

{{if not (TmuxHasSession $.Name)}}

{{- range $i, $cmd:= $.OnProjectStart }}
{{$cmd}}
{{- end }}

{{- range $i, $window := $.Windows }}
#====================== start ================================
# Window: {{$window.Name}}
{{$winId := Inc $i}}

{{if eq $i 0}}
{{$.Tmux}} new-session -d -s {{$.Name}} -n {{$window.Name}}
{{ else }}
{{$.Tmux}} new-window -c {{$.Root}} -t {{$.Name}}:{{$winId}} -n {{$window.Name}}
{{- end }}
# 设置窗口的根目录
# {{- if $window.Root}}
# {{$.Tmux}} send-keys -t {{$.Name}}:{{$i}} "cd {{$window.Root}}" C-m
# {{- end}}

{{$.Tmux}} select-layout -t {{$.Name}}:{{$winId}} main-vertical


{{- range $j, $pane := $window.Panes}}
{{$panelId := Inc $j}}
{{$.Tmux}} split-window -c {{$window.Root}} -t {{$.Name}}:{{$winId}}
{{$.Tmux}} select-layout -t {{$.Name}}:{{$winId}} tiled

{{- range $k, $cmd := $pane.Commands}}
{{$.Tmux}} send-keys -t {{$.Name}}:{{$winId}}.{{$panelId}} {{$cmd | Safe }} C-m
{{- end}}

{{- end}}

# 关闭最后一个多余的pane
{{$.Tmux}} kill-pane -t {{$.Name}}:{{$winId}}.0

#======================  end  ================================

{{$.Tmux}} list-panes -a
{{- end}}
{{ else }}
# Already have a session
{{- end }}


# 附加到 tmux 会话
#{{$.Tmux}} attach-session -t {{.Name}}
`
