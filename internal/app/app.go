package app

import (
	"context"
	"fmt"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"gitlab.com/mediasoft-internship/final-task/contracts/pkg/contracts/statistics"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"mediasoft-statistics/internal/bootstrap"
	"mediasoft-statistics/internal/client"
	"mediasoft-statistics/internal/config"
	"mediasoft-statistics/internal/kafka"
	"mediasoft-statistics/internal/repository"
	"mediasoft-statistics/internal/service"
	"mediasoft-statistics/pkg/logger"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func Run(log *logger.Logger, cfg config.Config) error {
	s := grpc.NewServer()
	mux := runtime.NewServeMux()
	ctx, cancel := context.WithCancel(context.Background())

	go runGRPCServer(log, cfg, s)
	go runHTTPServer(log, ctx, cfg, mux)

	gracefulShutDown(log, s, cancel)

	return nil
}

func runGRPCServer(log *logger.Logger, cfg config.Config, s *grpc.Server) {
	db, err := bootstrap.InitSqlxDB(cfg)
	if err != nil {
		log.Fatal("Init db conn %v", err)
	}

	restaurant, err := client.New(cfg)
	if err != nil {
		log.Fatal("New restaurant client %v", err)
	}
	defer func(restaurant *client.RestaurantClient) {
		err := restaurant.Close()
		if err != nil {
			log.Fatal("Close restaurant client %v", err)
		}
	}(restaurant)

	consumer, err := kafka.New(cfg.Kafka, log, repository.New(db))
	if err != nil {
		log.Fatal("failed get conn kafka %v", err.Error())
	}
	defer func(consumer *kafka.Consumer) {
		err := consumer.Close()
		if err != nil {
			log.Fatal("failed close consumer %v", err)
		}
	}(consumer)

	consumer.Consume(cfg.Topic)

	st := service.New(log, restaurant, repository.New(db))
	statistics.RegisterStatisticsServiceServer(s, st)

	l, err := net.Listen("tcp", fmt.Sprintf("%s:%d", cfg.GRPC.IP, cfg.GRPC.Port))
	if err != nil {
		log.Fatal("failed to listen tcp %s:%d, %v", cfg.GRPC.IP, cfg.GRPC.Port, err)
	}

	log.Info("starting listening grpc server at %s:%d", cfg.GRPC.IP, cfg.GRPC.Port)
	if err := s.Serve(l); err != nil {
		log.Fatal("error service grpc server %v", err)
	}

}

func runHTTPServer(log *logger.Logger, ctx context.Context, cfg config.Config, mux *runtime.ServeMux) {
	dialOptions := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	endPoint := fmt.Sprintf("%s:%d", cfg.HTTP.IP, cfg.GRPC.Port)

	if err := statistics.RegisterStatisticsServiceHandlerFromEndpoint(ctx, mux, endPoint, dialOptions); err != nil {
		log.Fatal("failed to register office service %v", err)
	}

	log.Info("starting listening http server at %s:%d", cfg.HTTP.IP, cfg.HTTP.Port)
	if err := http.ListenAndServe(fmt.Sprintf("%s:%d", cfg.HTTP.IP, cfg.HTTP.Port), mux); err != nil {
		log.Fatal("error service http server %v", err)
	}

}

func gracefulShutDown(log *logger.Logger, s *grpc.Server, cancel context.CancelFunc) {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	defer signal.Stop(ch)
	sig := <-ch
	log.Info("Received shutdown signal: %v -  Graceful shutdown done", sig)
	s.GracefulStop()
	cancel()
}
