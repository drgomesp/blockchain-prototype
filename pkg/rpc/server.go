package rpc

import (
	"context"
)

type Server struct{}

func NewServer() (*Server, error) {
	return &Server{}, nil
}

func (s *Server) Name() string {
	return "rpc"
}

func (s *Server) Start(_ context.Context) error {
	return nil
}

func (s *Server) Stop(_ context.Context) error {
	return nil
}
