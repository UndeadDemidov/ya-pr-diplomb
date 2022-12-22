package servers

import (
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/UndeadDemidov/ya-pr-diplomb/config"
	pbUser "github.com/UndeadDemidov/ya-pr-diplomb/gen_pb/user"
	deliveryGRPC "github.com/UndeadDemidov/ya-pr-diplomb/internal/delivery/grpc"
	user2 "github.com/UndeadDemidov/ya-pr-diplomb/internal/repo/postgres/user"
	"github.com/UndeadDemidov/ya-pr-diplomb/internal/services/user"
	"github.com/UndeadDemidov/ya-pr-diplomb/pkg/telemetry"
	grpc_ctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	"github.com/jackc/pgx/v5/pgxpool"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/reflection"
)

// https://github.com/grpc-ecosystem/go-grpc-middleware
// go-grpc-middleware/tracing/opentracing/ - add jaeger

type GRPC struct {
	logger telemetry.AppLogger
	cfg    *config.App
	db     *pgxpool.Pool
	srv    *grpc.Server
}

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

	pbUser.RegisterUserServiceServer(g.srv,
		deliveryGRPC.NewUserServer(g.logger, g.cfg, *user.NewService(user2.NewRepository(g.db, g.logger))))

	return &g
}

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
