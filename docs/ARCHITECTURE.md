# Kubemux AI Plugin - Architecture Visualization

## System Architecture

### High-Level Overview

```
┌─────────────────────────────────────────────────────────────────┐
│                         User Layer                              │
│  Natural Language: "Show me all my Kubernetes clusters"        │
└────────────────────────────┬────────────────────────────────────┘
                             │
                             ▼
┌─────────────────────────────────────────────────────────────────┐
│                      AI Assistant Layer                         │
│  ┌─────────────┐  ┌──────────────┐  ┌─────────────────┐       │
│  │   Claude    │  │GitHub Copilot│  │ Google Gemini   │       │
│  │  Desktop    │  │              │  │                 │       │
│  └──────┬──────┘  └──────┬───────┘  └────────┬────────┘       │
│         │                │                    │                 │
│         └────────────────┴────────────────────┘                 │
│                          │                                      │
│              Understands intent, calls tools                    │
└──────────────────────────┬──────────────────────────────────────┘
                           │
                           │ MCP Protocol (JSON-RPC 2.0)
                           │ Transport: stdio (stdin/stdout)
                           │
┌──────────────────────────▼──────────────────────────────────────┐
│                  MCP Server (kubemux-mcp-server)                │
│  ┌───────────────────────────────────────────────────────────┐  │
│  │  Protocol Handler (main.go)                               │  │
│  │  - initialize: Return capabilities                        │  │
│  │  - tools/list: Return available tools                     │  │
│  │  - tools/call: Execute requested tool                     │  │
│  └───────────────────────┬───────────────────────────────────┘  │
│                          │                                      │
│  ┌───────────────────────▼───────────────────────────────────┐  │
│  │  Tool Manager (server/server.go)                          │  │
│  │  - Register tools                                         │  │
│  │  - Route tool calls                                       │  │
│  │  - Format responses                                       │  │
│  └───────────────────────┬───────────────────────────────────┘  │
│                          │                                      │
│  ┌───────────────────────▼───────────────────────────────────┐  │
│  │  Tool Implementations (tools/*.go)                        │  │
│  │  ┌─────────────┐ ┌──────────────┐ ┌──────────────────┐   │  │
│  │  │list_clusters│ │list_sessions │ │create_session    │   │  │
│  │  └─────────────┘ └──────────────┘ └──────────────────┘   │  │
│  │  ┌─────────────┐ ┌──────────────┐                        │  │
│  │  │attach_session│ │list_kubeconfs│                       │  │
│  │  └─────────────┘ └──────────────┘                        │  │
│  └───────────────────────┬───────────────────────────────────┘  │
└──────────────────────────┼──────────────────────────────────────┘
                           │
                           │ Function calls
                           │
┌──────────────────────────▼──────────────────────────────────────┐
│                    Kubemux Core Library (lib/)                  │
│  ┌──────────────┐  ┌───────────┐  ┌─────────────────────┐     │
│  │ config.go    │  │ tmux.go   │  │ cloud_provider/     │     │
│  │ - ParseConfig│  │ - RunTmux │  │ - AWS EKS          │     │
│  │ - RenderERB  │  │ - HasTmux │  │ - GetClusters      │     │
│  └──────────────┘  └───────────┘  └─────────────────────┘     │
│  ┌──────────────┐  ┌───────────┐  ┌─────────────────────┐     │
│  │ zellij.go    │  │ window.go │  │ common/             │     │
│  │ - HasZellij  │  │ - Window  │  │ - Cluster           │     │
│  └──────────────┘  └───────────┘  └─────────────────────┘     │
└──────────────────────────┬──────────────────────────────────────┘
                           │
                           │ File system access
                           │ Process execution
                           │
┌──────────────────────────▼──────────────────────────────────────┐
│                        System Resources                         │
│  ┌──────────────────┐  ┌──────────────────┐  ┌──────────────┐  │
│  │ ~/.kube/         │  │ tmux/zellij      │  │ K8s Clusters │  │
│  │ - config files   │  │ - Sessions       │  │ - API access │  │
│  └──────────────────┘  └──────────────────┘  └──────────────┘  │
└─────────────────────────────────────────────────────────────────┘
```

## Data Flow: Example Interaction

### Scenario: User asks "Show me all my Kubernetes clusters"

```
Step 1: User Input
┌──────────────────────────────────────┐
│ User types in AI assistant:          │
│ "Show me all my Kubernetes clusters" │
└──────────────────┬───────────────────┘
                   │
                   ▼
Step 2: AI Processing
┌──────────────────────────────────────┐
│ AI Assistant analyzes intent:        │
│ - Recognizes: User wants cluster list│
│ - Finds tool: list_kubeconfigs       │
│ - Prepares JSON-RPC call             │
└──────────────────┬───────────────────┘
                   │
                   ▼
Step 3: MCP Request
┌──────────────────────────────────────┐
│ JSON-RPC 2.0 Message:                │
│ {                                    │
│   "jsonrpc": "2.0",                  │
│   "id": 123,                         │
│   "method": "tools/call",            │
│   "params": {                        │
│     "name": "list_kubeconfigs",      │
│     "arguments": {}                  │
│   }                                  │
│ }                                    │
└──────────────────┬───────────────────┘
                   │ stdin
                   ▼
Step 4: MCP Server Processing
┌──────────────────────────────────────┐
│ kubemux-mcp-server:                  │
│ 1. Parse JSON-RPC message            │
│ 2. Route to "list_kubeconfigs" tool  │
│ 3. Tool.Execute({})                  │
└──────────────────┬───────────────────┘
                   │
                   ▼
Step 5: Tool Execution
┌──────────────────────────────────────┐
│ list_kubeconfigs tool:               │
│ 1. Call lib.GetKubeConfigList()      │
│ 2. Read ~/.kube/ directory           │
│ 3. Filter .config files              │
│ 4. Format as JSON                    │
└──────────────────┬───────────────────┘
                   │
                   ▼
Step 6: Library Function
┌──────────────────────────────────────┐
│ lib.GetKubeConfigList():             │
│ homeDir := os.UserHomeDir()          │
│ files := os.ReadDir(homeDir/.kube)   │
│ return ["prod.config", "dev.config"] │
└──────────────────┬───────────────────┘
                   │
                   ▼
Step 7: Response Formatting
┌──────────────────────────────────────┐
│ Tool formats result:                 │
│ {                                    │
│   "kubeconfigs": [                   │
│     "prod-us-east.config",           │
│     "prod-eu-west.config",           │
│     "dev-local.config"               │
│   ],                                 │
│   "count": 3                         │
│ }                                    │
└──────────────────┬───────────────────┘
                   │
                   ▼
Step 8: MCP Response
┌──────────────────────────────────────┐
│ JSON-RPC Response:                   │
│ {                                    │
│   "jsonrpc": "2.0",                  │
│   "id": 123,                         │
│   "result": {                        │
│     "content": [{                    │
│       "type": "text",                │
│       "text": "{...JSON...}"         │
│     }]                               │
│   }                                  │
│ }                                    │
└──────────────────┬───────────────────┘
                   │ stdout
                   ▼
Step 9: AI Processing
┌──────────────────────────────────────┐
│ AI Assistant:                        │
│ - Parses JSON response               │
│ - Understands cluster list           │
│ - Formats for user                   │
└──────────────────┬───────────────────┘
                   │
                   ▼
Step 10: User Response
┌──────────────────────────────────────┐
│ AI Assistant displays:               │
│                                      │
│ "You have 3 Kubernetes clusters:     │
│ 1. prod-us-east.config               │
│ 2. prod-eu-west.config               │
│ 3. dev-local.config"                 │
└──────────────────────────────────────┘
```

## Component Interactions

### Tool Registration Flow

```
Program Start
     │
     ▼
┌─────────────────────┐
│ main()              │
│ - Create decoder    │
│ - Create encoder    │
│ - Initialize server │
└────────┬────────────┘
         │
         ▼
┌─────────────────────┐
│ NewMCPServer()      │
│ - Create tool map   │
│ - Call registerTools│
└────────┬────────────┘
         │
         ▼
┌─────────────────────────────────────────┐
│ registerTools()                         │
│ s.tools["list_clusters"] = NewListCl...│
│ s.tools["list_sessions"] = NewListSe...│
│ s.tools["create_session"] = NewCreate..│
│ s.tools["attach_session"] = NewAttach..│
│ s.tools["list_kubeconfigs"] = NewList..│
└─────────────────────────────────────────┘
```

### Request Handling Flow

```
JSON-RPC Request arrives via stdin
     │
     ▼
┌─────────────────────┐
│ decoder.Decode()    │
│ Parse JSON          │
└────────┬────────────┘
         │
         ▼
┌─────────────────────┐
│ handleRequest()     │
│ Switch on method    │
└────────┬────────────┘
         │
    ┌────┴────┬──────────┬───────────┐
    │         │          │           │
    ▼         ▼          ▼           ▼
┌────────┐ ┌──────┐ ┌─────────┐ ┌──────────┐
│initial-│ │tools/│ │tools/   │ │  other   │
│ize     │ │list  │ │call     │ │  error   │
└────────┘ └──────┘ └────┬────┘ └──────────┘
                         │
                         ▼
                  ┌──────────────┐
                  │CallTool()    │
                  │- Find tool   │
                  │- Execute     │
                  └──────┬───────┘
                         │
                         ▼
                  ┌──────────────┐
                  │Tool.Execute()│
                  │- Run logic   │
                  │- Return JSON │
                  └──────┬───────┘
                         │
                         ▼
                  ┌──────────────┐
                  │Format response│
                  │Encode JSON   │
                  └──────┬───────┘
                         │
                         ▼
                    Send to stdout
```

## Tool Architecture

### Base Tool Structure

```
┌────────────────────────────────────────┐
│           Tool Interface               │
│ ┌────────────────────────────────────┐ │
│ │ Schema() map[string]interface{}    │ │
│ │ Execute(args) (string, error)      │ │
│ └────────────────────────────────────┘ │
└────────────────────────────────────────┘
                   △
                   │ implements
                   │
┌──────────────────┴─────────────────────┐
│           BaseTool                     │
│ ┌────────────────────────────────────┐ │
│ │ - name        string               │ │
│ │ - description string               │ │
│ │ - inputSchema map[string]interface{}│ │
│ └────────────────────────────────────┘ │
│ ┌────────────────────────────────────┐ │
│ │ Schema() - returns tool definition │ │
│ └────────────────────────────────────┘ │
└────────────────────────────────────────┘
                   △
                   │ embeds
                   │
    ┌──────────────┴───────────────┐
    │                              │
┌───▼──────────┐          ┌────────▼──────┐
│ListClusters  │          │ListSessions   │
│Tool          │          │Tool           │
│              │          │               │
│Execute():    │          │Execute():     │
│- GetKubeConf │          │- Check plexer │
│- Format JSON │          │- List sessions│
└──────────────┘          └───────────────┘
```

## Security Model

```
┌────────────────────────────────────────────┐
│              User Space                    │
│  ┌──────────────────────────────────────┐  │
│  │  AI Assistant                        │  │
│  │  - Runs as user                      │  │
│  │  - User's permissions                │  │
│  └────────────┬─────────────────────────┘  │
│               │ stdio (isolated)            │
│  ┌────────────▼─────────────────────────┐  │
│  │  kubemux-mcp-server                  │  │
│  │  - Same user process                 │  │
│  │  - Read-only file access             │  │
│  │  - No network                        │  │
│  │  - No privilege escalation           │  │
│  └────────────┬─────────────────────────┘  │
│               │ function calls              │
│  ┌────────────▼─────────────────────────┐  │
│  │  kubemux lib                         │  │
│  │  - File system reads                 │  │
│  │  - Process spawning (tmux/zellij)   │  │
│  │  - K8s API (via kubeconfig)          │  │
│  └──────────────────────────────────────┘  │
└────────────────────────────────────────────┘

Security Boundaries:
- ✅ No network exposure
- ✅ stdio-only communication
- ✅ User permission inheritance
- ✅ No credential storage
- ✅ Commands generated, not executed
- ✅ Read-only kubeconfig access
```

## Deployment Scenarios

### Scenario 1: Claude Desktop (macOS)

```
┌────────────────────────────────────────────┐
│  macOS System                              │
│  ┌──────────────────────────────────────┐  │
│  │ Claude.app                           │  │
│  │ Reads: ~/Library/Application Support/│  │
│  │        Claude/claude_desktop_config.json│
│  └────────────┬─────────────────────────┘  │
│               │                             │
│               │ spawns                      │
│               │                             │
│  ┌────────────▼─────────────────────────┐  │
│  │ /usr/local/bin/kubemux-mcp-server    │  │
│  │ (stdio connected to Claude.app)      │  │
│  └────────────┬─────────────────────────┘  │
│               │                             │
│               │ reads                       │
│               │                             │
│  ┌────────────▼─────────────────────────┐  │
│  │ ~/.kube/                             │  │
│  │ - prod-us-east.config                │  │
│  │ - dev-local.config                   │  │
│  └──────────────────────────────────────┘  │
└────────────────────────────────────────────┘
```

### Scenario 2: VS Code + Copilot

```
┌────────────────────────────────────────────┐
│  VS Code                                   │
│  ┌──────────────────────────────────────┐  │
│  │ Copilot Extension                    │  │
│  │ + MCP Extension                      │  │
│  │ Reads: .vscode/settings.json         │  │
│  └────────────┬─────────────────────────┘  │
│               │                             │
│               │ spawns                      │
│               │                             │
│  ┌────────────▼─────────────────────────┐  │
│  │ kubemux-mcp-server                   │  │
│  │ (stdio connected to VS Code)         │  │
│  └──────────────────────────────────────┘  │
└────────────────────────────────────────────┘
```

---

This architecture visualization helps understand how all components work together to enable AI-powered Kubernetes cluster management through natural language.
