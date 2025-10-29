package server

import (
	"fmt"
	"kubemux/mcp-server/tools"
)

// MCPServer represents the MCP server instance
type MCPServer struct {
	tools map[string]tools.Tool
}

// NewMCPServer creates a new MCP server instance
func NewMCPServer() *MCPServer {
	server := &MCPServer{
		tools: make(map[string]tools.Tool),
	}
	server.registerTools()
	return server
}

// registerTools registers all available tools
func (s *MCPServer) registerTools() {
	// Register kubemux tools
	s.tools["list_clusters"] = tools.NewListClustersToolE()
	s.tools["list_kubeconfigs"] = tools.NewListKubeconfigsTool()
	s.tools["attach_session"] = tools.NewAttachSessionTool()
	s.tools["list_sessions"] = tools.NewListSessionsTool()
	s.tools["create_session"] = tools.NewCreateSessionTool()
}

// ListTools returns all available tools with their schemas
func (s *MCPServer) ListTools() []map[string]interface{} {
	var toolList []map[string]interface{}
	for _, tool := range s.tools {
		toolList = append(toolList, tool.Schema())
	}
	return toolList
}

// CallTool executes a tool with given arguments
func (s *MCPServer) CallTool(name string, arguments map[string]interface{}) (string, error) {
	tool, exists := s.tools[name]
	if !exists {
		return "", fmt.Errorf("tool not found: %s", name)
	}
	return tool.Execute(arguments)
}
