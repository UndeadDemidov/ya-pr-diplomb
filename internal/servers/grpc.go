package servers

import (
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/UndeadDemidov/ya-pr-diplomb/config"
	"github.com/UndeadDemidov/ya-pr-diplomb/gen_pb"
	deliveryGRPC "github.com/UndeadDemidov/ya-pr-diplomb/internal/delivery/grpc"
	"github.com/UndeadDemidov/ya-pr-diplomb/internal/repo/postgres"
	"github.com/UndeadDemidov/ya-pr-diplomb/internal/services"
	"github.com/UndeadDemidov/ya-pr-diplomb/pkg/telemetry"
	grpc_ctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	"github.com/jackc/pgx/v5/pgxpool"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/reflection"
)

// https://github.com/grpc-ecosystem/go-grpc-middleware
// go-grpc-middleware/tracing/opentracing/ - add jaeger

// GRPC server.
type GRPC struct {
	logger telemetry.AppLogger
	cfg    *config.App
	db     *pgxpool.Pool
	srv    *grpc.Server
}

// NewGRPC creates new instance of GRPC server with given options.
func NewGRPC(log telemetry.AppLogger, cfg *config.App, db *pgxpool.Pool) *GRPC {
	g := GRPC{
		logger: log,
		cfg:    cfg,
		db:     db,
	}

	// ToDo add metrics
	// metrics, err := metric.CreateMetrics(server.cfg.Metrics.URL, server.cfg.Metrics.ServiceName)
	// if err != nil {
	// 	server.logger.Errorf("CreateMetrics Error: %server", err)
	// }
	// server.logger.Info(
	// 	"Metrics available URL: %server, ServiceName: %server",
	// 	server.cfg.Metrics.URL,
	// 	server.cfg.Metrics.ServiceName,
	// )
	g.srv = grpc.NewServer(grpc.KeepaliveParams(keepalive.ServerParameters{
		MaxConnectionIdle: g.cfg.Server.MaxConnectionIdle * time.Minute,
		Timeout:           g.cfg.Server.Timeout * time.Second,
		MaxConnectionAge:  g.cfg.Server.MaxConnectionAge * time.Minute,
		Time:              g.cfg.Server.Timeout * time.Minute,
	}),

		// ToDo implement brutforce protection.
		// grpc.UnaryInterceptor(im.Logger),
		grpc.ChainUnaryInterceptor(
			grpc_ctxtags.UnaryServerInterceptor(),
			// grpc_prometheus.UnaryServerInterceptor,
			// grpcrecovery.UnaryServerInterceptor(),
		),
	)

	if g.cfg.Server.Mode != "prod" {
		reflection.Register(g.srv)
	}

	gen_pb.RegisterUserServiceServer(g.srv,
		deliveryGRPC.NewUserServer(g.logger, g.cfg, *services.NewUser(postgres.NewUser(g.db, g.logger))))

	return &g
}

// Run starts GRPC server with graceful shutdown.
func (g GRPC) Run() error {
	l, err := net.Listen("tcp", g.cfg.Server.Port)
	if err != nil {
		return fmt.Errorf("listener starter got an error: %w", err)
	}
	defer l.Close()

	go func() {
		g.logger.Infof("Server is listening on port: %v", g.cfg.Server.Port)
		if err = g.srv.Serve(l); err != nil {
			g.logger.Fatal(err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	<-quit
	g.srv.GracefulStop()
	g.logger.Info("Server Exited Properly")

	return nil
}
