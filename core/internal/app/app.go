package app

import (
	"fmt"
	"net"

	dbconn "github.com/Launchkit-org/LaunchKit/db"
	db "github.com/Launchkit-org/LaunchKit/db/sqlc"
	"github.com/Launchkit-org/LaunchKit/shared/config"
	"github.com/Launchkit-org/LaunchKit/shared/logger"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
)

type App struct {
	Server   *grpc.Server
	Listener net.Listener
	DBPool   *pgxpool.Pool
	Queries  *db.Queries
	Logger   zerolog.Logger
}

func StartApp(cfg *config.Config) (*App, error) {
	appLogger := logger.NewLogger(&cfg.Log)

	pool, err := dbconn.ConnectDB(&dbconn.Config{
		DSN:             cfg.Database.DSN,
		MaxOpenConns:    cfg.Database.MaxOpenConns,
		MinOpenConns:    cfg.Database.MinOpenConns,
		MaxConnLifetime: cfg.Database.MaxConnLifetime,
		MaxIdleTime:     cfg.Database.MaxIdleTime,
	})
	if err != nil {
		return nil, fmt.Errorf("connect db: %w", err)
	}

	queries := db.New(pool)

	// cfg.Core.HTTPAddr is a holdover name from the old REST stub — it's just
	// a listen address (":8081") and works fine for gRPC's TCP listener too.
	lis, err := net.Listen("tcp", cfg.Core.HTTPAddr)
	if err != nil {
		pool.Close()
		return nil, fmt.Errorf("listen on %s: %w", cfg.Core.HTTPAddr, err)
	}

	grpcServer := grpc.NewServer()

	healthServer := health.NewServer()
	healthpb.RegisterHealthServer(grpcServer, healthServer)
	healthServer.SetServingStatus("", healthpb.HealthCheckResponse_SERVING)

	reflection.Register(grpcServer)

	return &App{
		Server:   grpcServer,
		Listener: lis,
		DBPool:   pool,
		Queries:  queries,
		Logger:   appLogger,
	}, nil
}