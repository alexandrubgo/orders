package service

import (
	"context"
	"fmt"
	"log"

	"github.com/bekzourdk/orders/internal/config"
	"github.com/bekzourdk/orders/internal/logger"
	"github.com/bekzourdk/orders/internal/repository"
	"github.com/bekzourdk/orders/internal/server"
	"github.com/bekzourdk/orders/internal/service"
	"github.com/jackc/pgx/v4"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

var Cmd = &cobra.Command{
	Use: "service",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := run(); err != nil {
			return fmt.Errorf("run: %w", err)
		}

		return nil
	},
}

/*
func main() {
	log.Println("service started")

	if err := run(); err != nil {
		log.Fatalf("run: %v", err)
	}
}
*/

func run() error {
	logger, err := logger.New()
	if err != nil {
		log.Fatalf("run: failed creating logger: %v", err)
	}

	config, err := config.LoadConfig()
	if err != nil {
		logger.Fatal("run: config.LoadConfig", zap.Error(err))
	}

	conn, err := pgx.Connect(context.Background(), config.DB.Source)
	if err != nil {
		logger.Fatal("cannot connect to db", zap.Error(err))
	}

	ordersRepo := repository.NewOrdersRepository(conn)
	ordersService := service.NewOrdersService(ordersRepo)

	server := server.NewServer(config, logger, ordersService)

	if err := server.Run(); err != nil {
		return fmt.Errorf("server.Run: %w", err)
	}

	return nil
}
