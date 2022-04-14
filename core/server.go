package core

import (
	context2 "github.com/hamster-shared/hamster-gateway/core/context"
	"github.com/hamster-shared/hamster-gateway/core/corehttp"
	"os"
)

type Server struct {
	ctx context2.CoreContext
}

func NewServer(ctx context2.CoreContext) *Server {
	return &Server{
		ctx: ctx,
	}
}

func (s *Server) Run() {

	// 1: start api
	err := corehttp.StartApi(&s.ctx)
	if err != nil {
		os.Exit(1)
	}
}
