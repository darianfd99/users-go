package server

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/darianfd99/geo/pkg/service"
	"google.golang.org/grpc"

	"github.com/darianfd99/users-go/pkg/proto"
)

type Server struct {
	port            string
	shutdownTimeout time.Duration
}

func NewServer(ctx context.Context, port string, shutdownTimeout time.Duration, services service.Service) (context.Context, *Server) {
	srv := &Server{
		port:            port,
		shutdownTimeout: shutdownTimeout,
	}

	return serverContext(ctx), srv
}

func (s *Server) Run(ctx context.Context) {
	addr := fmt.Sprintf(":%s", s.port)
	listen, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatal(err)
	}

	srv := grpc.NewServer()
	proto.RegisterUsersServiceServer(srv, s)

	if err := srv.Serve(listen); err != nil {
		log.Fatal(err)
	}

	log.Println("app started")

	<-ctx.Done()
	log.Println("app shutting down")
	srv.GracefulStop()
	return
}

func serverContext(ctx context.Context) context.Context {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	ctx, cancel := context.WithCancel(ctx)
	go func() {
		<-quit
		cancel()
	}()

	return ctx

}
