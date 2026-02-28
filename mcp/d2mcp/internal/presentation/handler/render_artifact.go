package handler

import (
	"context"
	"encoding/base64"
	"fmt"
	"io"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"

	"github.com/hmsoft0815/mlcgo_mcp/mcp/d2mcp/internal/domain/entity"
	"github.com/hmsoft0815/mlcgo_mcp/mcp/d2mcp/internal/usecase"
	"github.com/hmsoft0815/mlcartifact"
)

// RenderArtifactHandler handles rendering a D2 source artifact to an SVG artifact.
type RenderArtifactHandler struct {
	useCase *usecase.DiagramUseCase
}

// NewRenderArtifactHandler creates a new handler.
func NewRenderArtifactHandler(useCase *usecase.DiagramUseCase) *RenderArtifactHandler {
	return &RenderArtifactHandler{
		useCase: useCase,
	}
}

// GetTool returns the MCP tool definition.
func (h *RenderArtifactHandler) GetTool() mcp.Tool {
	return mcp.NewTool(
		"render_artifact",
		mcp.WithDescription("Renders a D2 diagram stored as an artifact and saves the result as a new SVG artifact. Ideal for visualizing previously saved D2 source files."),
		mcp.WithString("artifactId", mcp.Description("ID or filename of the D2 source artifact"), mcp.Required()),
	)
}

// GetHandler returns the tool handler function.
func (h *RenderArtifactHandler) GetHandler() server.ToolHandlerFunc {
	return h.Handle
}

// Handle processes the render request.
func (h *RenderArtifactHandler) Handle(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	artifactID := mcp.ParseString(request, "artifactId", "")
	if artifactID == "" {
		return mcp.NewToolResultError("artifactId is required"), nil
	}

	// 1. Fetch D2 source from artifact service
	cli, err := mlcartifact.NewClient()
	if err != nil {
		return mcp.NewToolResultErrorFromErr("Failed to connect to artifact service", err), nil
	}
	defer cli.Close()

	res, err := cli.Read(ctx, artifactID)
	if err != nil {
		return mcp.NewToolResultErrorFromErr("Failed to read D2 artifact", err), nil
	}

	// 2. Render to SVG
	reader, err := h.useCase.RenderDiagram(ctx, string(res.Content), entity.FormatSVG, nil)
	if err != nil {
		return mcp.NewToolResultErrorFromErr("Failed to render diagram", err), nil
	}

	data, err := io.ReadAll(reader)
	if err != nil {
		return mcp.NewToolResultErrorFromErr("Failed to read rendered data", err), nil
	}

	// 3. Save SVG back to artifact service
	filename := fmt.Sprintf("%s.svg", res.Filename)
	writeRes, err := cli.Write(ctx, filename, data, mlcartifact.WithSource("d2mcp"))
	if err != nil {
		// Return image anyway but report error
		imgResult := mcp.NewToolResultImage("svg", base64.StdEncoding.EncodeToString(data), "image/svg+xml")
		imgResult.Content = append(imgResult.Content, mcp.NewTextContent("\nWarning: Failed to save as artifact: "+err.Error()))
		return imgResult, nil
	}

	// 4. Return image preview and artifact reference
	imgResult := mcp.NewToolResultImage("svg", base64.StdEncoding.EncodeToString(data), "image/svg+xml")
	fileTag := fmt.Sprintf("<file id=\"%s\" type=\"image/svg+xml\">%s</file>", writeRes.Id, writeRes.Filename)
	imgResult.Content = append(imgResult.Content, mcp.TextContent{
		Type: "text",
		Text: fmt.Sprintf("\nArtifact saved: %s\nUse this tag in your response to the user.", fileTag),
	})

	return imgResult, nil
}
