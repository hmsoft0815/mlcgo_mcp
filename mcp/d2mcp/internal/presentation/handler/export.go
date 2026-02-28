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

// ExportHandler handles diagram export operations.
type ExportHandler struct {
	useCase *usecase.DiagramUseCase
}

// NewExportHandler creates a new export handler.
func NewExportHandler(useCase *usecase.DiagramUseCase) *ExportHandler {
	return &ExportHandler{
		useCase: useCase,
	}
}

// GetTool returns the MCP tool definition.
func (h *ExportHandler) GetTool() mcp.Tool {
	return mcp.NewTool(
		"d2_export",
		mcp.WithDescription("Export an existing diagram to SVG format. The diagram must first be created using d2_create. Supports exporting all D2 features including SQL tables, UML classes, sequence diagrams, code blocks, and markdown-rich documentation."),
		mcp.WithString("diagramId", mcp.Description("ID of the diagram to export"), mcp.Required()),
	)
}

// GetHandler returns the tool handler function.
func (h *ExportHandler) GetHandler() server.ToolHandlerFunc {
	return h.Handle
}

// Handle processes the export request.
func (h *ExportHandler) Handle(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	// Extract arguments.
	diagramID := mcp.ParseString(request, "diagramId", "")
	if diagramID == "" {
		return mcp.NewToolResultError("diagramId is required"), nil
	}

	// 1. Export the diagram as SVG using UseCase
	reader, err := h.useCase.ExportDiagram(ctx, diagramID, entity.FormatSVG)
	if err != nil {
		return mcp.NewToolResultErrorFromErr("Failed to export diagram", err), nil
	}

	// 2. Read the output
	data, err := io.ReadAll(reader)
	if err != nil {
		return mcp.NewToolResultErrorFromErr("Failed to read exported data", err), nil
	}

	// 3. Save to Shared Artifact Service (Phase 2 Integration)
	artifactCli, err := mlcartifact.NewClient()
	if err == nil {
		defer artifactCli.Close()
		filename := fmt.Sprintf("%s.svg", diagramID)
		res, err := artifactCli.Write(ctx, filename, data, mlcartifact.WithSource("d2mcp"))
		if err == nil {
			// Return both image for preview and reference tag for persistence/sharing
			imgResult := mcp.NewToolResultImage("svg", base64.StdEncoding.EncodeToString(data), "image/svg+xml")

			// We append the artifact info as text content to the result
			fileTag := fmt.Sprintf("<file id=\"%s\" type=\"image/svg+xml\">%s</file>", res.Id, res.Filename)
			imgResult.Content = append(imgResult.Content, mcp.TextContent{
				Type: "text",
				Text: fmt.Sprintf("\nArtifact saved: %s\nUse this tag in your response to the user so they can access the file permanently.", fileTag),
			})
			return imgResult, nil
		}
	}

	// Fallback to just base64 if artifact service is unavailable
	return mcp.NewToolResultImage("svg", base64.StdEncoding.EncodeToString(data), "image/svg+xml"), nil
}

// getMimeType removed as we only support SVG now.
