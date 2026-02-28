package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"github.com/i2y/d2mcp/internal/infrastructure/d2"
	"github.com/i2y/d2mcp/internal/infrastructure/mcp"
	"github.com/i2y/d2mcp/internal/presentation/handler"
	"github.com/i2y/d2mcp/internal/usecase"
	mcptypes "github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

const (
	// ServerName is the name of the MCP server.
	ServerName = "d2mcp"
	// ServerVersion is the version of the MCP server.
	ServerVersion = "0.6.0-mlc"
)

// transportConfig holds all transport-related CLI configuration.
type transportConfig struct {
	Transport         string
	Addr              string
	BaseURL           string
	BasePath          string
	KeepAlive         int
	EndpointPath      string
	HeartbeatInterval int
	Stateless         bool
}

// toolRegistration bundles a tool definition with its handler for batch registration.
type toolRegistration struct {
	tool    mcptypes.Tool
	handler server.ToolHandlerFunc
}

func main() {
	// Disable D2 logs.
	os.Setenv("D2_LOG_LEVEL", "NONE")

	// Parse command line flags.
	var (
		transport         string
		addr              string
		baseURL           string
		basePath          string
		keepAlive         int
		endpointPath      string
		heartbeatInterval int
		stateless         bool
	)
	flag.StringVar(&transport, "transport", "stdio", "Transport mode: stdio, sse, or streamable")
	flag.StringVar(&addr, "addr", ":3000", "Address to listen on for SSE/Streamable HTTP transport")
	flag.StringVar(&baseURL, "base-url", "", "Base URL for SSE transport (auto-generated if empty)")
	flag.StringVar(&basePath, "base-path", "/mcp", "Base path for SSE endpoints")
	flag.IntVar(&keepAlive, "keep-alive", 30, "Keep-alive interval in seconds for SSE")
	flag.StringVar(&endpointPath, "endpoint-path", "/mcp", "Endpoint path for Streamable HTTP transport")
	flag.IntVar(&heartbeatInterval, "heartbeat-interval", 30, "Heartbeat interval in seconds for Streamable HTTP")
	flag.BoolVar(&stateless, "stateless", false, "Enable stateless mode for Streamable HTTP")
	flag.Parse()

	// Validate transport mode.
	if transport != "stdio" && transport != "sse" && transport != "streamable" {
		fmt.Fprintf(os.Stderr, "Invalid transport mode: %s. Must be 'stdio', 'sse', or 'streamable'\n", transport)
		os.Exit(1)
	}

	configureLogging(transport, addr)

	ctx := context.Background()

	// Initialize domain layer.
	oracleRepo := d2.NewD2OracleRepository()
	diagramUseCase := usecase.NewDiagramUseCase(oracleRepo)
	oracleUseCase := usecase.NewOracleUseCase(oracleRepo)

	// Initialize MCP server with transport.
	srv, err := mcp.NewServer(ServerName, ServerVersion)
	if err != nil {
		log.Fatalf("Failed to create MCP server: %v", err)
	}

	configureTransport(srv, transportConfig{
		Transport:         transport,
		Addr:              addr,
		BaseURL:           baseURL,
		BasePath:          basePath,
		KeepAlive:         keepAlive,
		EndpointPath:      endpointPath,
		HeartbeatInterval: heartbeatInterval,
		Stateless:         stateless,
	})

	// Register all tools.
	tools := buildToolRegistrations(diagramUseCase, oracleUseCase)
	for _, t := range tools {
		if err := srv.RegisterTool(t.tool, t.handler); err != nil {
			log.Fatalf("Failed to register tool '%s': %v", t.tool.Name, err)
		}
	}

	// Start the server.
	log.Printf("Starting %s v%s (%s transport)...", ServerName, ServerVersion, transport)
	if err := srv.Start(ctx); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}

// configureLogging sets up the log output based on transport mode.
// In stdio mode, logs go to a file to avoid interfering with stdio communication.
func configureLogging(transport, addr string) {
	if transport == "stdio" {
		logFile, err := os.OpenFile("/tmp/d2mcp.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
		if err != nil {
			log.SetOutput(io.Discard)
		} else {
			log.SetOutput(logFile)
		}
	} else {
		log.SetOutput(os.Stderr)
		log.Printf("Starting in %s mode on %s", transport, addr)
	}
}

// configureTransport sets the transport and its configuration on the server.
func configureTransport(srv *mcp.Server, cfg transportConfig) {
	switch cfg.Transport {
	case "sse":
		baseURL := cfg.BaseURL
		if baseURL == "" {
			if cfg.Addr[0] == ':' {
				baseURL = fmt.Sprintf("http://localhost%s", cfg.Addr)
			} else {
				baseURL = fmt.Sprintf("http://%s", cfg.Addr)
			}
		}

		srv.WithTransport(mcp.TransportSSE)
		srv.WithSSEConfig(&mcp.SSEConfig{
			Addr:              cfg.Addr,
			BaseURL:           baseURL,
			StaticBasePath:    cfg.BasePath,
			KeepAliveInterval: time.Duration(cfg.KeepAlive) * time.Second,
		})

		log.Printf("  SSE stream: %s%s/sse", baseURL, cfg.BasePath)
		log.Printf("  Messages:   %s%s/message", baseURL, cfg.BasePath)

	case "streamable":
		srv.WithTransport(mcp.TransportStreamableHTTP)
		srv.WithStreamableHTTPConfig(&mcp.StreamableHTTPConfig{
			Addr:              cfg.Addr,
			EndpointPath:      cfg.EndpointPath,
			HeartbeatInterval: time.Duration(cfg.HeartbeatInterval) * time.Second,
			Stateless:         cfg.Stateless,
		})

		log.Printf("  Endpoint: http://localhost%s%s", cfg.Addr, cfg.EndpointPath)
		if cfg.Stateless {
			log.Printf("  Mode: stateless")
		}

	default:
		srv.WithTransport(mcp.TransportStdio)
	}
}

// buildToolRegistrations creates all handler instances and returns their tool registrations.
func buildToolRegistrations(diagramUC *usecase.DiagramUseCase, oracleUC *usecase.OracleUseCase) []toolRegistration {
	createHandler := handler.NewCreateHandler(diagramUC)
	exportHandler := handler.NewExportHandler(diagramUC)
	renderArtifactHandler := handler.NewRenderArtifactHandler(diagramUC)
	oracleCreate := handler.NewOracleCreateHandler(oracleUC)
	oracleSet := handler.NewOracleSetHandler(oracleUC)
	oracleDelete := handler.NewOracleDeleteHandler(oracleUC)
	oracleMove := handler.NewOracleMoveHandler(oracleUC)
	oracleRename := handler.NewOracleRenameHandler(oracleUC)
	oracleGet := handler.NewOracleGetHandler(oracleUC)
	oracleSerialize := handler.NewOracleSerializeHandler(oracleUC)

	return []toolRegistration{
		{createHandler.GetTool(), createHandler.GetHandler()},
		{exportHandler.GetTool(), exportHandler.GetHandler()},
		{renderArtifactHandler.GetTool(), renderArtifactHandler.GetHandler()},
		{oracleCreate.GetTool(), oracleCreate.GetHandler()},
		{oracleSet.GetTool(), oracleSet.GetHandler()},
		{oracleDelete.GetTool(), oracleDelete.GetHandler()},
		{oracleMove.GetTool(), oracleMove.GetHandler()},
		{oracleRename.GetTool(), oracleRename.GetHandler()},
		{oracleGet.GetTool(), oracleGet.GetHandler()},
		{oracleSerialize.GetTool(), oracleSerialize.GetHandler()},
	}
}
