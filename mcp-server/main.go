package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	log "github.com/sirupsen/logrus"
	"kubemux/mcp-server/server"
)

// MCP Protocol Messages
type MCPRequest struct {
	JSONRPC string          `json:"jsonrpc"`
	ID      interface{}     `json:"id"`
	Method  string          `json:"method"`
	Params  json.RawMessage `json:"params,omitempty"`
}

type MCPResponse struct {
	JSONRPC string      `json:"jsonrpc"`
	ID      interface{} `json:"id"`
	Result  interface{} `json:"result,omitempty"`
	Error   *MCPError   `json:"error,omitempty"`
}

type MCPError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    string `json:"data,omitempty"`
}

func main() {
	log.SetFormatter(&log.TextFormatter{
		DisableTimestamp: true,
	})
	log.SetOutput(os.Stderr) // Log to stderr to not interfere with JSON-RPC

	mcpServer := server.NewMCPServer()

	decoder := json.NewDecoder(os.Stdin)
	encoder := json.NewEncoder(os.Stdout)

	log.Info("Kubemux MCP Server started")

	for {
		var request MCPRequest
		if err := decoder.Decode(&request); err != nil {
			if err == io.EOF {
				break
			}
			log.Errorf("Error decoding request: %v", err)
			continue
		}

		log.Debugf("Received request: %s", request.Method)

		response := handleRequest(&request, mcpServer)
		if err := encoder.Encode(response); err != nil {
			log.Errorf("Error encoding response: %v", err)
		}
	}
}

func handleRequest(request *MCPRequest, mcpServer *server.MCPServer) MCPResponse {
	switch request.Method {
	case "initialize":
		return MCPResponse{
			JSONRPC: "2.0",
			ID:      request.ID,
			Result: map[string]interface{}{
				"protocolVersion": "2024-11-05",
				"serverInfo": map[string]interface{}{
					"name":    "kubemux-mcp-server",
					"version": "1.0.0",
				},
				"capabilities": map[string]interface{}{
					"tools": map[string]interface{}{},
				},
			},
		}

	case "tools/list":
		tools := mcpServer.ListTools()
		return MCPResponse{
			JSONRPC: "2.0",
			ID:      request.ID,
			Result: map[string]interface{}{
				"tools": tools,
			},
		}

	case "tools/call":
		var params struct {
			Name      string                 `json:"name"`
			Arguments map[string]interface{} `json:"arguments"`
		}
		if err := json.Unmarshal(request.Params, &params); err != nil {
			return MCPResponse{
				JSONRPC: "2.0",
				ID:      request.ID,
				Error: &MCPError{
					Code:    -32602,
					Message: "Invalid params",
					Data:    err.Error(),
				},
			}
		}

		result, err := mcpServer.CallTool(params.Name, params.Arguments)
		if err != nil {
			return MCPResponse{
				JSONRPC: "2.0",
				ID:      request.ID,
				Error: &MCPError{
					Code:    -32000,
					Message: "Tool execution error",
					Data:    err.Error(),
				},
			}
		}

		return MCPResponse{
			JSONRPC: "2.0",
			ID:      request.ID,
			Result: map[string]interface{}{
				"content": []map[string]interface{}{
					{
						"type": "text",
						"text": result,
					},
				},
			},
		}

	default:
		return MCPResponse{
			JSONRPC: "2.0",
			ID:      request.ID,
			Error: &MCPError{
				Code:    -32601,
				Message: fmt.Sprintf("Method not found: %s", request.Method),
			},
		}
	}
}
