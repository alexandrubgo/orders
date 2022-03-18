package server

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/bekzourdk/orders/internal/config"
	ordersGRPC "github.com/bekzourdk/orders/internal/delivery/grpc"
	"github.com/bekzourdk/orders/internal/service"
	"github.com/bekzourdk/orders/pb"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/reflection"
)

type Server struct {
	config        *config.Config
	logger        *zap.Logger
	ordersService *service.OrdersService
}

func NewServer(
	config *config.Config,
	logger *zap.Logger,
	ordersService *service.OrdersService,
) *Server {
	return &Server{
		config:        config,
		logger:        logger,
		ordersService: ordersService,
	}
}

func (s *Server) Run() error {
	mux := runtime.NewServeMux()

	srv := &http.Server{
		Addr:              s.config.Servers.Http,
		Handler:           mux,
		ReadTimeout:       1 * time.Second,
		ReadHeaderTimeout: 1 * time.Second,
		WriteTimeout:      1 * time.Second,
		IdleTimeout:       1 * time.Second,
	}

	l, err := net.Listen("tcp", s.config.Servers.Grpc)
	if err != nil {
		return fmt.Errorf("net.Listen: %w", err)
	}

	defer l.Close()

	grpcServer := grpc.NewServer(grpc.KeepaliveParams(keepalive.ServerParameters{
		MaxConnectionIdle: 5 * time.Minute,
		Timeout:           15 * time.Second,
		MaxConnectionAge:  5 * time.Minute,
		Time:              15 * time.Minute,
	}),
	)

	ordersServerGRPC := ordersGRPC.NewOrdersGRPC(*s.ordersService)
	pb.RegisterOrdersServiceServer(grpcServer, ordersServerGRPC)
	reflection.Register(grpcServer)

	ctx, stop := signal.NotifyContext(context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
		syscall.SIGQUIT,
		syscall.SIGINT,
		syscall.SIGHUP,
	)

	go func() {
		<-ctx.Done()
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
		defer func() {
			stop()
			cancel()
		}()

		grpcServer.GracefulStop()

		if err := srv.Shutdown(ctx); err != nil {
			s.logger.Fatal("failed shutting down", zap.Error(err))
		}
	}()

	if err := pb.RegisterOrdersServiceHandlerFromEndpoint(
		ctx,
		mux,
		s.config.Servers.Grpc,
		[]grpc.DialOption{
			grpc.WithInsecure(),
		}); err != nil {
		return err
	}

	go func() {
		if err := grpcServer.Serve(l); err != nil {
			s.logger.Fatal("failed starting grpc server", zap.Error(err))
		}
	}()

	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}

	return nil
}
