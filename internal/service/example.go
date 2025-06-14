package service

import (
	"context"
	"fmt"
	v1 "github.com/go-nunu/nunu-layout-mcp/api/v1"
	"github.com/go-nunu/nunu-layout-mcp/internal/repository"
	"github.com/mark3labs/mcp-go/mcp"
	"io"
	"net/http"
	"strings"
)

type ExampleService interface {
	HttpTool(ctx context.Context, params *v1.HttpToolRequest) (*mcp.CallToolResult, error)
	GetTinyImageTool(ctx context.Context) (*mcp.CallToolResult, error)
}

func NewExampleService(
	service *Service,
	exampleRepo repository.ExampleRepository,
) ExampleService {
	return &exampleService{
		exampleRepo: exampleRepo,
		Service:     service,
	}
}

type exampleService struct {
	exampleRepo repository.ExampleRepository
	*Service
}

func (s *exampleService) GetTinyImageTool(ctx context.Context) (*mcp.CallToolResult, error) {
	img, _ := s.exampleRepo.FetchImage(ctx)
	return &mcp.CallToolResult{
		Content: []mcp.Content{
			mcp.TextContent{
				Type: "text",
				Text: "This is a tiny image:",
			},
			mcp.ImageContent{
				Type:     "image",
				Data:     img,
				MIMEType: "image/png",
			},
			mcp.TextContent{
				Type: "text",
				Text: "The image above is the MCP tiny image.",
			},
		},
	}, nil
}

func (s *exampleService) HttpTool(ctx context.Context, params *v1.HttpToolRequest) (*mcp.CallToolResult, error) {

	// Create and send request
	var req *http.Request
	var err error
	if params.Body != "" {
		req, err = http.NewRequest(params.Method, params.Url, strings.NewReader(params.Body))
	} else {
		req, err = http.NewRequest(params.Method, params.Url, nil)
	}
	if err != nil {
		return mcp.NewToolResultErrorFromErr("unable to create request", err), nil
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return mcp.NewToolResultErrorFromErr("unable to execute request", err), nil
	}
	defer resp.Body.Close()

	// Return response
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return mcp.NewToolResultErrorFromErr("unable to read request response", err), nil
	}

	return mcp.NewToolResultText(fmt.Sprintf("Status: %d\nBody: %s", resp.StatusCode, string(respBody))), nil
}
