package tools

import (
	"encoding/json"
	"fmt"
	"kubemux/lib"
)

// ListClustersTool lists available Kubernetes clusters from kubeconfigs
type ListClustersTool struct {
	BaseTool
}

// NewListClustersToolE creates a new ListClustersTool
func NewListClustersToolE() *ListClustersTool {
	return &ListClustersTool{
		BaseTool: BaseTool{
			name:        "list_clusters",
			description: "List all available Kubernetes clusters from kubeconfig files in ~/.kube directory",
			inputSchema: map[string]interface{}{
				"type":       "object",
				"properties": map[string]interface{}{},
			},
		},
	}
}

// Execute lists all kubeconfig files
func (t *ListClustersTool) Execute(arguments map[string]interface{}) (string, error) {
	configs := lib.GetKubeConfigList()
	
	result := map[string]interface{}{
		"kubeconfigs": configs,
		"count":       len(configs),
	}
	
	jsonResult, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to marshal result: %w", err)
	}
	
	return string(jsonResult), nil
}

// ListKubeconfigsTool is an alias for ListClustersTool
type ListKubeconfigsTool struct {
	BaseTool
}

// NewListKubeconfigsTool creates a new ListKubeconfigsTool
func NewListKubeconfigsTool() *ListKubeconfigsTool {
	return &ListKubeconfigsTool{
		BaseTool: BaseTool{
			name:        "list_kubeconfigs",
			description: "List all kubeconfig files in ~/.kube directory",
			inputSchema: map[string]interface{}{
				"type":       "object",
				"properties": map[string]interface{}{},
			},
		},
	}
}

// Execute lists all kubeconfig files
func (t *ListKubeconfigsTool) Execute(arguments map[string]interface{}) (string, error) {
	configs := lib.GetKubeConfigList()
	
	result := map[string]interface{}{
		"kubeconfigs": configs,
		"count":       len(configs),
	}
	
	jsonResult, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to marshal result: %w", err)
	}
	
	return string(jsonResult), nil
}
