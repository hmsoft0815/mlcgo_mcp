package handler

import (
	"context"
	"encoding/base64"
	"io"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"

	"github.com/hmsoft0815/mlcgo_mcp/mcp/d2mcp/internal/domain/entity"
	"github.com/hmsoft0815/mlcgo_mcp/mcp/d2mcp/internal/usecase"
)

// SaveHandler handles the d2_save tool.
type SaveHandler struct {
	useCase *usecase.DiagramUseCase
}

// NewSaveHandler creates a new save handler.
func NewSaveHandler(useCase *usecase.DiagramUseCase) *SaveHandler {
	return &SaveHandler{
		useCase: useCase,
	}
}

// GetTool returns the MCP tool definition.
func (h *SaveHandler) GetTool() mcp.Tool {
	return mcp.NewTool(
		"d2_save",
		mcp.WithDescription("Render an existing diagram as SVG and return it as an image result (base64 encoded). The diagram must be created first using d2_create."),
		mcp.WithString("diagramId", mcp.Description("ID of the diagram to render"), mcp.Required()),
	)
}

// GetHandler returns the tool handler function.
func (h *SaveHandler) GetHandler() server.ToolHandlerFunc {
	return h.Handle
}

// Handle processes the save request.
func (h *SaveHandler) Handle(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	// Extract arguments.
	diagramID := mcp.ParseString(request, "diagramId", "")
	if diagramID == "" {
		return mcp.NewToolResultError("diagramId is required"), nil
	}

	// Export the diagram as SVG.
	reader, err := h.useCase.ExportDiagram(ctx, diagramID, entity.FormatSVG)
	if err != nil {
		return mcp.NewToolResultErrorFromErr("Failed to export diagram", err), nil
	}

	// Read the output.
	data, err := io.ReadAll(reader)
	if err != nil {
		return mcp.NewToolResultErrorFromErr("Failed to read exported output", err), nil
	}

	// Return SVG as image result (base64 encoded).
	return mcp.NewToolResultImage("svg", base64.StdEncoding.EncodeToString(data), "image/svg+xml"), nil
}

