package service

import (
	"gitlab.com/mediasoft-internship/final-task/contracts/pkg/contracts/statistics"
	"mediasoft-statistics/internal/client"
	"mediasoft-statistics/internal/repository"
	"mediasoft-statistics/pkg/logger"
)

type Service struct {
	log             *logger.Logger
	restaurant      *client.RestaurantClient
	orderRepository repository.OrderRepository
	statistics.UnimplementedStatisticsServiceServer
}

func New(log *logger.Logger,
	restaurant *client.RestaurantClient,
	orderRepository repository.OrderRepository) *Service {
	return &Service{
		log:             log,
		restaurant:      restaurant,
		orderRepository: orderRepository,
	}
}
