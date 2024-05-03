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

# LEFT STATUS
{{$.Tmux}} set -g status-left-length 100
# {{$.Tmux}} set -g status-left-style default
# {{$.Tmux}} set -g status-left '#(byobu-status tmux_left)'
# {{$.Tmux}} set -g status-left "#[fg=colour220]#h\
#    #[fg=colour196] #(ip addr | grep -e 'state UP' -A 2 | awk '/inet /{printf $2}')\
#    #[fg=colour39] #(sensors | awk '/CPU/{printf $2}')\
#    #[fg=colour40] #(free -m -h | awk '/Mem/{printf \$3\"\/\"\$2\}')\
#    #[fg=colour128] #(free -m | awk '/Mem{printf 100*$2/$3" %"}')\
#    #[fg=colour202] #([ $(cat /sys/class/power_supply/AC/online) == 1 ] && printf %s'🗲') #(cat /sys/class/power_supply/BAT0/capacity)\%\
#    #[fg=colour7] #([ ! -z $(ip a | egrep 'ppp0|tun0' -A 2 | awk '/inet /{printf $2}') ] && echo $(ip a | egrep 'ppp0|tun0' -A 2 | awk '/inet /{printf \"VPN \"$2}'))\
#    #[default]"

# RIGHT STATUS
{{$.Tmux}} set -g status-right-length 100
# {{$.Tmux}} set -g status-right-style default
# {{$.Tmux}} set -g status-right '#(byobu-status tmux_right)'
# {{$.Tmux}} set -g status-right "#[fg=colour39] #(uptime | awk '{printf \$(NF-2)\" \"\$(NF-1)\" \"\$(NF)}' | tr -d ',')\
#   #[fg=colour40] %F\
#   #[fg=colour128] %T\
#   #[fg=colour202] %Z\
#   #[default]"


# 附加到 tmux 会话
#{{$.Tmux}} attach-session -t {{.Name}}
# echo {{$.Tmux}} {{ StringsJoin .TmuxArgs " "}}
# {{$.Tmux}} {{ StringsJoin .TmuxArgs " "}}
`
