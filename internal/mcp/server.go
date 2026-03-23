package mcp

import (
	"github.com/charmbracelet/log"
	"github.com/mark3labs/mcp-go/server"
)

type Server struct {
	logger *log.Logger
	server *server.MCPServer
}

func NewServer(logger *log.Logger) *Server {
	srv := server.NewMCPServer("hati", "0.1.0")

	s := &Server{
		logger: logger,
		server: srv,
	}

	s.registerAllTools()

	return s
}

func (s *Server) registerAllTools() {
	// Add tools specifically for Hati execution later
}

func (s *Server) RunStdio() error {
	return server.ServeStdio(s.server)
}
