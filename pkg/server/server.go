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

	"google.golang.org/grpc"

	"github.com/darianfd99/users-go/pkg/proto"
	"github.com/darianfd99/users-go/pkg/service"
)

type Server struct {
	service.Service

	port            string
	shutdownTimeout time.Duration
}

func NewServer(ctx context.Context, port string, services service.Service, shutdownTimeout time.Duration) (context.Context, *Server) {
	srv := &Server{
		Service: services,

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

	go func() {
		if err := srv.Serve(listen); err != nil {
			log.Fatal(err)
		}
	}()

	log.Println("app started")

	<-ctx.Done()
	log.Println("app shutting down")
	srv.GracefulStop()
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
