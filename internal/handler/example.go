package handler

import (
	"context"
	"encoding/base64"
	"fmt"
	v1 "github.com/go-nunu/nunu-layout-mcp/api/v1"
	"github.com/go-nunu/nunu-layout-mcp/internal/model"
	"github.com/go-nunu/nunu-layout-mcp/internal/service"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"log"
	"strconv"
	"time"
)

type ExampleHandler interface {
	AddTool(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error)
	EchoTool(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error)
	HttpTool(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error)
	SampleLLMTool(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error)
	LongRunningOperationTool(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error)
	GetTinyImageTool(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error)

	SimplePrompt(ctx context.Context, request mcp.GetPromptRequest) (*mcp.GetPromptResult, error)
	ComplexPrompt(ctx context.Context, request mcp.GetPromptRequest) (*mcp.GetPromptResult, error)

	GeneratedResource(ctx context.Context, request mcp.ReadResourceRequest) ([]mcp.ResourceContents, error)
	ReadResource(ctx context.Context, request mcp.ReadResourceRequest) ([]mcp.ResourceContents, error)
	ResourceTemplate(ctx context.Context, request mcp.ReadResourceRequest) ([]mcp.ResourceContents, error)

	Notification(ctx context.Context, notification mcp.JSONRPCNotification)
	SendNotification(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error)
}

func NewExampleHandler(
	tool *Handler,
	exampleSvc service.ExampleService,
) ExampleHandler {
	return &exampleHandler{
		exampleSvc: exampleSvc,
		Handler:    tool,
	}
}

type exampleHandler struct {
	exampleSvc service.ExampleService
	*Handler
}

func (h exampleHandler) EchoTool(
	ctx context.Context,
	req mcp.CallToolRequest,
) (*mcp.CallToolResult, error) {
	params := v1.EchoToolRequest{}
	err := req.BindArguments(&params)
	if err != nil {
		return nil, err
	}
	return &mcp.CallToolResult{
		Content: []mcp.Content{
			mcp.TextContent{
				Type: "text",
				Text: fmt.Sprintf("Echo: %s", params.Message),
			},
		},
	}, nil
}
func (h exampleHandler) AddTool(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	params := v1.AddToolRequest{}
	err := req.BindArguments(&params)
	if err != nil {
		return nil, err
	}
	sum := params.A + params.B
	return &mcp.CallToolResult{
		Content: []mcp.Content{
			mcp.TextContent{
				Type: "text",
				Text: fmt.Sprintf("The sum of %f and %f is %f.", params.A, params.B, sum),
			},
		},
	}, nil
}

func (h exampleHandler) HttpTool(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	params := v1.HttpToolRequest{}
	err := req.BindArguments(&params)
	if err != nil {
		return nil, err
	}

	return h.exampleSvc.HttpTool(ctx, &params)
}

func (h exampleHandler) SampleLLMTool(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := req.GetArguments()
	prompt, _ := args["prompt"].(string)
	maxTokens, _ := args["maxTokens"].(float64)

	// This is a mock implementation. In a real scenario, you would use the server's RequestSampling method.
	result := fmt.Sprintf(
		"Sample LLM result for prompt: '%s' (max tokens: %d)",
		prompt,
		int(maxTokens),
	)

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			mcp.TextContent{
				Type: "text",
				Text: fmt.Sprintf("LLM sampling result: %s", result),
			},
		},
	}, nil
}

func (h exampleHandler) SimplePrompt(ctx context.Context, request mcp.GetPromptRequest) (*mcp.GetPromptResult, error) {
	return &mcp.GetPromptResult{
		Description: "A simple prompt without arguments",
		Messages: []mcp.PromptMessage{
			{
				Role: mcp.RoleUser,
				Content: mcp.TextContent{
					Type: "text",
					Text: "This is a simple prompt without arguments.",
				},
			},
		},
	}, nil
}
func (h exampleHandler) ReadResource(ctx context.Context, request mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
	return []mcp.ResourceContents{
		mcp.TextResourceContents{
			URI:      "test://static/resource",
			MIMEType: "text/plain",
			Text:     "This is a sample resource",
		},
	}, nil
}
func (h exampleHandler) ResourceTemplate(ctx context.Context, request mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
	return []mcp.ResourceContents{
		mcp.TextResourceContents{
			URI:      request.Params.URI,
			MIMEType: "text/plain",
			Text:     "This is a sample resource",
		},
	}, nil
}
func (h exampleHandler) GeneratedResource(ctx context.Context, request mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
	uri := request.Params.URI

	var resourceNumber string
	if _, err := fmt.Sscanf(uri, "test://static/resource/%s", &resourceNumber); err != nil {
		return nil, fmt.Errorf("invalid resource URI format: %w", err)
	}

	num, err := strconv.Atoi(resourceNumber)
	if err != nil {
		return nil, fmt.Errorf("invalid resource number: %w", err)
	}

	index := num - 1

	if index%2 == 0 {
		return []mcp.ResourceContents{
			mcp.TextResourceContents{
				URI:      uri,
				MIMEType: "text/plain",
				Text:     fmt.Sprintf("Text content for resource %d", num),
			},
		}, nil
	} else {
		return []mcp.ResourceContents{
			mcp.BlobResourceContents{
				URI:      uri,
				MIMEType: "application/octet-stream",
				Blob:     base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("Binary content for resource %d", num))),
			},
		}, nil
	}
}
func (h exampleHandler) Notification(ctx context.Context, notification mcp.JSONRPCNotification) {
	log.Printf("Received notification: %s", notification.Method)
}
func (h exampleHandler) GetTinyImageTool(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return h.exampleSvc.GetTinyImageTool(ctx)
}
func (h exampleHandler) LongRunningOperationTool(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	arguments := request.GetArguments()
	progressToken := request.Params.Meta.ProgressToken
	duration, _ := arguments["duration"].(float64)
	steps, _ := arguments["steps"].(float64)
	stepDuration := duration / steps
	server := server.ServerFromContext(ctx)

	for i := 1; i < int(steps)+1; i++ {
		time.Sleep(time.Duration(stepDuration * float64(time.Second)))
		if progressToken != nil {
			err := server.SendNotificationToClient(
				ctx,
				"notifications/progress",
				map[string]any{
					"progress":      i,
					"total":         int(steps),
					"progressToken": progressToken,
					"message":       fmt.Sprintf("Server progress %v%%", int(float64(i)*100/steps)),
				},
			)
			if err != nil {
				return nil, fmt.Errorf("failed to send notification: %w", err)
			}
		}
	}

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			mcp.TextContent{
				Type: "text",
				Text: fmt.Sprintf(
					"Long running operation completed. Duration: %f seconds, Steps: %d.",
					duration,
					int(steps),
				),
			},
		},
	}, nil
}
func (h exampleHandler) SendNotification(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {

	server := server.ServerFromContext(ctx)

	err := server.SendNotificationToClient(
		ctx,
		"notifications/progress",
		map[string]any{
			"progress":      10,
			"total":         10,
			"progressToken": 0,
		},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to send notification: %w", err)
	}

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			mcp.TextContent{
				Type: "text",
				Text: "notification sent successfully",
			},
		},
	}, nil
}

func (h exampleHandler) ComplexPrompt(ctx context.Context, request mcp.GetPromptRequest) (*mcp.GetPromptResult, error) {
	arguments := request.Params.Arguments
	return &mcp.GetPromptResult{
		Description: "A complex prompt with arguments",
		Messages: []mcp.PromptMessage{
			{
				Role: mcp.RoleUser,
				Content: mcp.TextContent{
					Type: "text",
					Text: fmt.Sprintf(
						"This is a complex prompt with arguments: temperature=%s, style=%s",
						arguments["temperature"],
						arguments["style"],
					),
				},
			},
			{
				Role: mcp.RoleAssistant,
				Content: mcp.TextContent{
					Type: "text",
					Text: "I understand. You've provided a complex prompt with temperature and style arguments. How would you like me to proceed?",
				},
			},
			{
				Role: mcp.RoleUser,
				Content: mcp.ImageContent{
					Type:     "image",
					Data:     model.MCP_TINY_IMAGE,
					MIMEType: "image/png",
				},
			},
		},
	}, nil
}
