package server

import (
	"context"
	"fmt"
	"github.com/go-nunu/nunu-layout-mcp/internal/handler"
	"github.com/go-nunu/nunu-layout-mcp/internal/model"
	"github.com/go-nunu/nunu-layout-mcp/pkg/log"
	servermcp "github.com/go-nunu/nunu-layout-mcp/pkg/server/mcp"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

func NewMCPServer(
	conf *viper.Viper,
	logger *log.Logger,
	exampleHandler handler.ExampleHandler,
) *servermcp.Server {
	s := setupSrv(conf, logger)

	s.AddResource(mcp.NewResource("test://static/resource",
		"Static Resource",
		mcp.WithMIMEType("text/plain"),
	), exampleHandler.ReadResource)
	s.AddResourceTemplate(
		mcp.NewResourceTemplate(
			"test://dynamic/resource/{id}",
			"Dynamic Resource",
		),
		exampleHandler.ResourceTemplate,
	)

	resources := generateResources()
	for _, resource := range resources {
		s.AddResource(resource, exampleHandler.GeneratedResource)
	}

	s.AddPrompt(mcp.NewPrompt(string(model.SIMPLE),
		mcp.WithPromptDescription("A simple prompt"),
	), exampleHandler.SimplePrompt)
	s.AddPrompt(mcp.NewPrompt(string(model.COMPLEX),
		mcp.WithPromptDescription("A complex prompt"),
		mcp.WithArgument("temperature",
			mcp.ArgumentDescription("The temperature parameter for generation"),
			mcp.RequiredArgument(),
		),
		mcp.WithArgument("style",
			mcp.ArgumentDescription("The style to use for the response"),
			mcp.RequiredArgument(),
		),
	), exampleHandler.ComplexPrompt)
	s.AddTool(mcp.NewTool(string(model.ECHO),
		mcp.WithDescription("Echoes back the input"),
		mcp.WithString("message",
			mcp.Description("Message to echo"),
			mcp.Required(),
		),
	), exampleHandler.EchoTool)
	s.AddTool(mcp.NewTool("http_request",
		mcp.WithDescription("Make HTTP requests to external APIs"),
		mcp.WithString("method",
			mcp.Required(),
			mcp.Description("HTTP method to use"),
			mcp.Enum("GET", "POST", "PUT", "DELETE"),
		),
		mcp.WithString("url",
			mcp.Required(),
			mcp.Description("URL to send the request to"),
			mcp.Pattern("^https?://.*"),
		),
		mcp.WithString("body",
			mcp.Description("Request body (for POST/PUT)"),
		),
	), exampleHandler.HttpTool)
	s.AddTool(
		mcp.NewTool("notify"),
		exampleHandler.SendNotification,
	)

	s.AddTool(mcp.NewTool(string(model.ADD),
		mcp.WithDescription("Adds two numbers"),
		mcp.WithNumber("a",
			mcp.Description("First number"),
			mcp.Required(),
		),
		mcp.WithNumber("b",
			mcp.Description("Second number"),
			mcp.Required(),
		),
	), exampleHandler.AddTool)
	s.AddTool(mcp.NewTool(
		string(model.LONG_RUNNING_OPERATION),
		mcp.WithDescription(
			"Demonstrates a long running operation with progress updates",
		),
		mcp.WithNumber("duration",
			mcp.Description("Duration of the operation in seconds"),
			mcp.DefaultNumber(10),
		),
		mcp.WithNumber("steps",
			mcp.Description("Number of steps in the operation"),
			mcp.DefaultNumber(5),
		),
	), exampleHandler.LongRunningOperationTool)

	s.AddTool(mcp.Tool{
		Name:        string(model.SAMPLE_LLM),
		Description: "Samples from an LLM using MCP's sampling feature",
		InputSchema: mcp.ToolInputSchema{
			Type: "object",
			Properties: map[string]any{
				"prompt": map[string]any{
					"type":        "string",
					"description": "The prompt to send to the LLM",
				},
				"maxTokens": map[string]any{
					"type":        "number",
					"description": "Maximum number of tokens to generate",
					"default":     100,
				},
			},
		},
	}, exampleHandler.SampleLLMTool)
	s.AddTool(mcp.NewTool(string(model.GET_TINY_IMAGE),
		mcp.WithDescription("Returns the MCP_TINY_IMAGE"),
	), exampleHandler.GetTinyImageTool)

	s.AddNotificationHandler("notification", exampleHandler.Notification)

	return s
}

func setupSrv(conf *viper.Viper, logger *log.Logger) *servermcp.Server {
	mcpServer := server.NewMCPServer(
		conf.GetString("mcp.name"),
		conf.GetString("mcp.version"),
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
		servermcp.WithSSESrv(conf.GetString("mcp.sse_addr"), server.NewSSEServer(
			mcpServer,
			server.WithSSEEndpoint("/sse"),
		)),
		// StreamableHTTP
		servermcp.WithStreamableHTTPSrv(conf.GetString("mcp.http_addr"), server.NewStreamableHTTPServer(
			mcpServer,
			server.WithEndpointPath("/mcp"),
		)),
	)
}
func newHooks(logger *log.Logger) *server.Hooks {
	hooks := &server.Hooks{}

	hooks.AddBeforeAny(func(ctx context.Context, id any, method mcp.MCPMethod, message any) {
		logger.Info("AddBeforeAny",
			zap.Any("method", method),
			zap.Any("id", id),
			zap.Any("message", message),
		)
	})

	hooks.AddOnSuccess(func(ctx context.Context, id any, method mcp.MCPMethod, message any, result any) {
		logger.Info("AddOnSuccess",
			zap.Any("method", method),
			zap.Any("id", id),
			zap.Any("message", message),
			zap.Any("result", result),
		)
	})

	hooks.AddOnError(func(ctx context.Context, id any, method mcp.MCPMethod, message any, err error) {
		logger.Info("AddOnError",
			zap.Any("method", method),
			zap.Any("id", id),
			zap.Any("message", message),
			zap.Error(err),
		)
	})

	hooks.AddBeforeInitialize(func(ctx context.Context, id any, message *mcp.InitializeRequest) {
		logger.Info("AddBeforeInitialize",
			zap.Any("id", id),
			zap.Any("message", message),
		)
	})

	hooks.AddOnRequestInitialization(func(ctx context.Context, id any, message any) error {
		logger.Info("AddOnRequestInitialization: ", zap.Any("id", id), zap.Any("message", message))
		// authorization verification and other preprocessing tasks are performed.
		return nil
	})
	hooks.AddAfterInitialize(func(ctx context.Context, id any, message *mcp.InitializeRequest, result *mcp.InitializeResult) {
		logger.Info("AddAfterInitialize",
			zap.Any("id", id),
			zap.Any("message", message),
			zap.Any("result", result),
		)
	})

	hooks.AddAfterCallTool(func(ctx context.Context, id any, message *mcp.CallToolRequest, result *mcp.CallToolResult) {
		logger.Info("AddAfterCallTool",
			zap.Any("id", id),
			zap.Any("message", message),
			zap.Any("result", result),
		)
	})

	hooks.AddBeforeCallTool(func(ctx context.Context, id any, message *mcp.CallToolRequest) {
		logger.Info("AddBeforeCallTool",
			zap.Any("id", id),
			zap.Any("message", message),
		)
	})

	return hooks
}
func generateResources() []mcp.Resource {
	resources := make([]mcp.Resource, 100)
	for i := 0; i < 100; i++ {
		uri := fmt.Sprintf("test://static/resource/%d", i+1)
		if i%2 == 0 {
			resources[i] = mcp.Resource{
				URI:      uri,
				Name:     fmt.Sprintf("Resource %d", i+1),
				MIMEType: "text/plain",
			}
		} else {
			resources[i] = mcp.Resource{
				URI:      uri,
				Name:     fmt.Sprintf("Resource %d", i+1),
				MIMEType: "application/octet-stream",
			}
		}
	}
	return resources
}
