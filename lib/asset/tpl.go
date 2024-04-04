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

{{- if $.Debug }}
{{$.Tmux}} list-panes -a
{{- end }}
{{- end}}
{{ else }}
# Already have a session
{{- end }}

# show host name and IP address on left side of status bar
# {{$.Tmux}} set-option -g status-left-length 60
# {{$.Tmux}} set-option -g status-left "#[fg=colour198]: #h : #[fg=brightblue]#(curl ipv4.ip.sb) #(ifconfig eno1 | grep 'inet ' | awk '{print \"eno1 \" \$2}')"

# {{$.Tmux}} set-option -g status-right-length 60
# {{$.Tmux}} set-option -g status-right "#[fg=blue]#S #I:#P #[fg=yellow]: %d %b %Y #[fg=green]: %l:%M %p : #(date -u | awk '{print $4}') :"
# {{$.Tmux}} set-option -g status-right "#[fg=blue]#(tmux-cpu --no-color)"

# 附加到 tmux 会话
#{{$.Tmux}} attach-session -t {{.Name}}
`
