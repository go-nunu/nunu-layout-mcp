# 🚀 Quickly Build a High-Performance Go MCP Server

[简体中文介绍](https://github.com/go-nunu/nunu-layout-mcp/blob/main/README_zh.md)

This project is a **sample MCP Server** built using the [Nunu](https://github.com/go-nunu/nunu) framework and [MCP-GO](https://github.com/model-context-protocol/mcp-go). It allows you to rapidly set up a Golang-based MCP Server and test/debug it using the MCP Inspector tool.

![Nunu](https://github.com/go-nunu/nunu/blob/main/.github/assets/banner.png)

---

## 🚀 Quick Start: Run the MCP Server in 3 Minutes

### 1. Create a Project

<details>
<summary>Option 1: Clone the Template Repository</summary>

```bash
git clone https://github.com/go-nunu/nunu-layout-mcp.git
# Note: The default project name is nunu-layout-mcp
```

</details>

<details>
<summary>Option 2: Create a New Project via Nunu CLI (Recommended)</summary>

```bash
go install github.com/go-nunu/nunu@latest

nunu new mcp-demo -r https://github.com/go-nunu/nunu-layout-mcp.git
```

</details>

---

### 2. Build and Start

#### Build the MCP Server

```bash
cd mcp-demo

go build -ldflags="-s -w" -o ./bin/server ./cmd/server
```

#### Start MCP Inspector

MCP Inspector is an interactive developer tool provided by the MCP community for testing and debugging:

```bash
npx -y @modelcontextprotocol/inspector ./bin/server
# Requires Node.js to be installed
```

---

### 3. Test the Service

Open your browser: `http://127.0.0.1:6274` to test different transport protocols.

| Transport Type | Address / Parameter         | Use Case                | Pros                                        | Cons                                   |
| -------------- | --------------------------- | ----------------------- | ------------------------------------------- | -------------------------------------- |
| STDIO          | `./bin/server`              | CLI tools, desktop apps | Simple, secure, no network needed           | Local only, single client              |
| SSE            | `http://localhost:3001/sse` | Web real-time comm.     | Multi-client, real-time, browser friendly   | HTTP overhead, server-to-client only   |
| StreamableHTTP | `http://localhost:3002/mcp` | Web services, APIs      | Standard protocol, caching & load balancing | No real-time support, slightly complex |
| In-Process     | (no external address)       | Embedded, testing       | No serialization, ultra-fast                | In-process only                        |

![STDIO](https://github.com/go-nunu/nunu/blob/main/.github/assets/mcp/stdio.png)
![SSE](https://github.com/go-nunu/nunu/blob/main/.github/assets/mcp/sse.png)
![StreamableHTTP](https://github.com/go-nunu/nunu/blob/main/.github/assets/mcp/http.png)

---

## 🛠️ Development Guide

As this project is built on the [Nunu architecture](https://github.com/go-nunu/nunu/blob/main/docs/zh/architecture.md), it’s recommended to understand the framework before development.

---

### 📡 MCP Server Development

See the [MCP-GO Server Docs](https://mcp-go.dev/server)

This project enables three protocols by default: STDIO, SSE, StreamableHTTP. You can modify or disable them as needed:

```go
// File: internal/server/mcp.go

func setupSrv(logger *log.Logger) *servermcp.Server {
	mcpServer := server.NewMCPServer(
		"example-servers/everything",
		"1.0.0",
		server.WithResourceCapabilities(true, true),
		server.WithPromptCapabilities(true),
		server.WithToolCapabilities(true),
		server.WithLogging(),
		server.WithHooks(newHooks(logger)),
	)

	return servermcp.NewServer(logger,
		servermcp.WithMCPSrv(mcpServer),

		// STDIO
		servermcp.WithStdioSrv(true),

		// SSE
		servermcp.WithSSESrv(":3001", server.NewSSEServer(
			mcpServer,
			server.WithSSEEndpoint("/sse"),
		)),

		// StreamableHTTP
		servermcp.WithStreamableHTTPSrv(":3002", server.NewStreamableHTTPServer(
			mcpServer,
			server.WithEndpointPath("/mcp"),
		)),
	)
}
```

#### Call Flow Diagram

```txt
client
  ↓
internal/server/mcp.go
  ↓
internal/handler/example.go
  ↓
internal/service/example.go
  ↓
internal/repository/example.go
  ↓
DB / Third-party services
```

**Note:**

If `MCP STDIO` is enabled, **no logs should be printed to the terminal**, or the communication will break. You must configure logs to write only to a file:

```yaml
# File: config/local.yml
log:
  log_level: debug
  mode: file               # file, console, or both
  encoding: console        # json or console
  log_file_name: "./storage/logs/server.log"
  max_backups: 30
  max_age: 7
  max_size: 1024
  compress: true
```

To view logs in Unix systems:

```bash
tail -f storage/logs/server.log
```

---

### 🤝 Integrating MCP Client

See [MCP-GO Client Docs](https://mcp-go.dev/clients)

You can register a client in `internal/repository/repository.go`, for example:

```go
func NewStdioMCPClient() *client.Client {
	c, err := client.NewStdioMCPClient(
		"go", []string{}, "run", "/path/to/server/main.go",
	)
	if err != nil {
		panic(err)
	}
	defer c.Close()
	return c
}
```

To integrate other protocols or clients, follow similar patterns used in `redis`, `gorm`, etc.

---

## 📚 Resources

* [Nunu Architecture](https://github.com/go-nunu/nunu/blob/main/docs/zh/architecture.md)
* [MCP Protocol Overview](https://modelcontextprotocol.io/docs/protocol/overview)
* [MCP Inspector Guide](https://modelcontextprotocol.io/docs/tools/inspector)
* [MCP-GO Documentation](https://mcp-go.dev/getting-started)

---

## Final Notes

Currently, there is **no official Golang SDK** provided by the MCP organization. The most mature open-source project is [`mark3labs/mcp-go`](https://github.com/mark3labs/mcp-go).

We’re looking forward to an official Golang SDK, and this project will be updated accordingly once it’s released.

Official Go MCP SDK: [https://github.com/golang/tools/tree/master/internal/mcp](https://github.com/golang/tools/tree/master/internal/mcp)

---

## 📄 License

Nunu is released under the [MIT License](LICENSE) — free to use and contribute!

