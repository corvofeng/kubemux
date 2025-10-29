# Kubemux MCP Server Usage Examples

This document provides practical examples of using kubemux with AI assistants through the MCP server.

## Table of Contents

1. [Basic Operations](#basic-operations)
2. [Cluster Management](#cluster-management)
3. [Session Management](#session-management)
4. [Advanced Workflows](#advanced-workflows)

---

## Basic Operations

### Example 1: Discovering Available Clusters

**User Input (to AI Assistant):**
```
What Kubernetes clusters do I have available?
```

**AI Assistant Response:**
The AI will call the `list_kubeconfigs` tool and respond with something like:

```
You have the following Kubernetes clusters configured:
- prod-us-east.config
- prod-eu-west.config
- staging.config
- dev-local.config

Total: 4 clusters
```

**Behind the scenes:**
```json
{
  "method": "tools/call",
  "params": {
    "name": "list_kubeconfigs",
    "arguments": {}
  }
}
```

### Example 2: Listing Active Sessions

**User Input:**
```
Show me all my active kubemux sessions
```

**AI Response:**
```
You currently have 2 active tmux sessions:
- prod-monitoring
- dev-workspace

Would you like to attach to any of these?
```

---

## Cluster Management

### Example 3: Connecting to a Specific Cluster

**User Input:**
```
I need to connect to the production cluster in US East
```

**AI Response:**
```
To connect to the prod-us-east cluster, run this command:

kubemux kube --kube prod-us-east.config

This will create a new kubemux session with your production cluster context.
```

### Example 4: Creating Sessions for Multiple Clusters

**User Input:**
```
Create separate sessions for my prod and staging clusters, using zellij
```

**AI Response:**
```
I'll help you create two sessions with zellij:

1. For production:
   kubemux kube --kube prod-us-east.config --plexer zellij

2. For staging:
   kubemux kube --kube staging.config --plexer zellij

You can run these commands in different terminals to work with both clusters simultaneously.
```

---

## Session Management

### Example 5: Reconnecting to an Existing Session

**User Input:**
```
I was working on something in my prod-monitoring session. How do I get back to it?
```

**AI Response:**
```
To reconnect to your prod-monitoring session, use:

tmux -L prod-monitoring attach-session -t prod-monitoring

This will attach you to the existing session where you left off.
```

### Example 6: Session Overview Before Connecting

**User Input:**
```
I forgot what sessions I have running. Can you help me find the right one for my development work?
```

**AI Response:**
```
Let me check your active sessions...

You have these sessions:
- prod-monitoring (tmux) - likely for production
- dev-workspace (tmux) - this sounds like your development session
- staging-debug (zellij) - for staging environment

Based on the name, "dev-workspace" seems to be your development session.

To connect: tmux -L dev-workspace attach-session -t dev-workspace
```

---

## Advanced Workflows

### Example 7: Setting Up a New Project Environment

**User Input:**
```
I'm starting a new project. Set up a kubemux environment using the 'microservices' configuration and connect to the dev cluster
```

**AI Response:**
```
I'll help you set up your project environment. Run this command:

kubemux --project microservices kube --kube dev-local.config

This will:
1. Load your 'microservices' project configuration from ~/.tmuxinator/microservices.yml
2. Connect to your dev-local cluster
3. Set up the window/pane layout defined in your configuration
```

---

## Tips and Best Practices

### Naming Conventions

When the AI creates sessions, it uses smart naming:
- By kubeconfig: Uses the config filename (e.g., "prod-us-east" from "prod-us-east.config")
- By project: Uses the project name (e.g., "monitoring" from monitoring.yml)

### Session Persistence

Sessions persist even after disconnecting. The AI can help you:
- List all sessions: "What sessions are running?"
- Reconnect: "Connect me to session X"
- Clean up: "How do I close the old sessions?"

---

**Note**: The exact responses from AI assistants may vary, but the functionality remains consistent. These examples demonstrate the natural language capabilities enabled by the MCP server integration.
