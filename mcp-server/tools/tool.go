package tools

// Tool interface defines the contract for all MCP tools
type Tool interface {
	// Schema returns the JSON schema definition for the tool
	Schema() map[string]interface{}
	// Execute runs the tool with the given arguments
	Execute(arguments map[string]interface{}) (string, error)
}

// BaseTool provides common functionality for all tools
type BaseTool struct {
	name        string
	description string
	inputSchema map[string]interface{}
}

// Schema returns the tool's JSON schema
func (t *BaseTool) Schema() map[string]interface{} {
	return map[string]interface{}{
		"name":        t.name,
		"description": t.description,
		"inputSchema": t.inputSchema,
	}
}
