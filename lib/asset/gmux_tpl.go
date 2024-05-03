package asset

const KubemuxKubeconfig = `
name: <%= @settings["name"] %>
root: ~/

socket_name: <%= @settings["name"] %>
on_project_start:
  - export KUBECONFIG=~/.kube/<%= @settings["kubeconfig"] %>
startup_window: kubectl

windows:
  - pods: kubectl get pods
  - namespaces: kubectl get ns
`
