# 快速构建高性能 Go MCP Server

本项目是一个基于 [Nunu](https://github.com/go-nunu/nunu) 框架和 [MCP-GO](https://github.com/model-context-protocol/mcp-go) 实现的 **MCP Server 示例**。通过本示例，你可以快速搭建一个 Golang 版本的 MCP Server，并使用 MCP Inspector 工具进行测试与调试。


![Nunu](https://github.com/go-nunu/nunu/blob/main/.github/assets/banner.png)

---


## 🚀 快速开始：三分钟跑通 MCP Server 示例

### 一、创建项目

<details>
<summary>方式一：直接 clone 模板仓库</summary>

```bash
git clone https://github.com/go-nunu/nunu-layout-mcp.git
# 注意：此方式项目包名默认为 nunu-layout-mcp
```

</details>

<details>
<summary>方式二：使用 Nunu CLI 创建新项目（推荐）</summary>

```bash
go install github.com/go-nunu/nunu@latest

nunu new mcp-demo -r https://github.com/go-nunu/nunu-layout-mcp.git
```

</details>

---

### 二、编译与启动

#### 编译 MCP Server

```bash
cd mcp-demo

go build -ldflags="-s -w" -o ./bin/server ./cmd/server
```

#### 启动 MCP Inspector

MCP Inspector 是官方提供的交互式开发者工具，用于测试和调试 MCP 服务：

```bash
npx -y @modelcontextprotocol/inspector ./bin/server
# 需要预先安装 Node.js 才能使用 npx
```

---

### 三、测试服务

访问浏览器：`http://127.0.0.1:6274`，即可测试不同传输协议。

| 传输方式           | 地址/参数                       | 使用场景          | 优点                | 缺点                      |
| -------------- |-----------------------------| ------------- | ----------------- | ----------------------- |
| STDIO          | `./bin/server`              | CLI 工具、桌面应用   | 简洁、安全、无需网络        | 仅支持本地、单一客户端             |
| SSE            | `http://localhost:3001/sse` | Web 实时通信      | 多客户端支持、实时推送、浏览器友好 | 有 HTTP 开销，仅支持服务器到客户端的推送 |
| StreamableHTTP | `http://localhost:3002/mcp` | Web 服务、API 通信 | 标准协议、支持缓存与负载均衡    | 不支持实时通信，逻辑实现稍复杂         |
| In-Process     | （无外部地址）                     | 嵌入式系统、测试      | 无需序列化、性能极高        | 仅支持进程内调用                |

![STDIO](https://github.com/go-nunu/nunu/blob/main/.github/assets/mcp/stdio.png)
![SSE](https://github.com/go-nunu/nunu/blob/main/.github/assets/mcp/sse.png)
![StreamableHTTP](https://github.com/go-nunu/nunu/blob/main/.github/assets/mcp/http.png)

---

## 🛠️ 开发说明

由于本项目基于 [Nunu 架构](https://github.com/go-nunu/nunu/blob/main/docs/zh/architecture.md) 实现，建议在开发前先了解该框架的基本原理与使用方法。

---

### 📡 MCP Server 开发

查看 [MCP-GO Server 文档](https://mcp-go.dev/server)

项目默认启用了三种协议：STDIO、SSE、StreamableHTTP。可根据实际需求修改或禁用：

```go
// 文件: internal/server/mcp.go

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
		// STDIO 通信
		servermcp.WithStdioSrv(true),

		// SSE 通信
		servermcp.WithSSESrv(":3001", server.NewSSEServer(
			mcpServer,
			server.WithSSEEndpoint("/sse"),
		)),

		// StreamableHTTP 通信
		servermcp.WithStreamableHTTPSrv(":3002", server.NewStreamableHTTPServer(
			mcpServer,
			server.WithEndpointPath("/mcp"),
		)),
	)
}
```

#### 调用链路示意

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
DB / 第三方服务
```

**注意事项：**

如果启用了`MCP STDIO`协议，那么命令行终端中是不允许出现任何日志输出的，这个时候项目必须关闭终端日志输出，也就是服务日志只写入到文件。
```
// 文件：config/local.yml，重点是mode字段
log:
  log_level: debug
  mode: file               # file or console or both
  encoding: console           # json or console
  log_file_name: "./storage/logs/server.log"
  max_backups: 30
  max_age: 7
  max_size: 1024
  compress: true
```

Unix开发环境可以使用如下命令查看服务器日志：
```
tail -f storage/logs/server.log
```

---

### 🤝 MCP Client 集成

查看 [MCP-GO Client 文档](https://mcp-go.dev/clients)

可在 `internal/repository/repository.go` 中注册客户端，如：

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

如需集成其它协议或客户端，只需仿照 `redis`、`gorm` 的注入方式进行即可。

---

## 📚 相关链接

* [Nunu 项目架构](https://github.com/go-nunu/nunu/blob/main/docs/zh/architecture.md)
* [MCP 协议说明](https://modelcontextprotocol.io/docs/protocol/overview)
* [MCP Inspector 使用指南](https://modelcontextprotocol.io/docs/tools/inspector)
* [MCP-GO 官方文档](https://mcp-go.dev/getting-started)

---

## 后话
由于目前MCP官方并没有提供Golang的SDK。而开源社区相对稳定完善的`Golang MCP SDK`只有 `https://github.com/mark3labs/mcp-go` 这一个项目。


后续我们可以关注期待下Golang官方的实现，等待官方实现正式发布后，本项目也会同步更新。

Go官方MCP SDK：https://github.com/golang/tools/tree/master/internal/mcp

## 📄 License

Nunu 使用 [MIT License](LICENSE) 开源发布，欢迎自由使用与贡献。

