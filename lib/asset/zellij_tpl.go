package asset

// ZellijSessionCreateTemplate is the template for creating a new zellij session
const ZellijSessionCreateTemplate = `#!/bin/bash
{{ if $.Debug }}
set -ex
{{ end }}

{{- range $i, $cmd:= $.OnProjectStart }}
{{$cmd}}
{{- end }}

zellij -s  {{$.Name}}
`
const ZellijSessionAttachTemplate = `#!/bin/bash
{{ if $.Debug }}
set -ex
{{ end }}

zellij -s {{$.Name}} attach {{$.Name}}
`

const ZellijSessionPrepareTemplate = `#!/bin/bash
{{ if $.Debug }}
set -ex
{{ end }}

{{- range $i, $window := $.Windows }}

{{- if eq $i 0}}
# Rename the first window: {{$window.Name}}
zellij -s {{$.Name}} action go-to-tab 1
zellij -s {{$.Name}} action rename-tab {{$window.Name}}
{{ else }}
# Create Window: {{$window.Name}}
zellij -s {{$.Name}} action new-tab -n {{$window.Name}}
{{- end }}

{{- range $j, $pane := $window.Panes}}

{{- if eq $j 0}}
{{ else }}
# Create new pane in the same window
zellij -s {{$.Name}} action new-pane -d down
{{- end}}

{{- range $k, $cmd := $pane.Commands}}
zellij -s {{$.Name}} action write-chars {{$cmd | Safe }}
zellij -s {{$.Name}} action write 13  # 13 means Enter key in ASCII
{{- end}}

{{- end}}
{{- end}}

{{- if $.Debug }}
zellij -s {{$.Name}} action query-tab-names
{{- end }}
`
