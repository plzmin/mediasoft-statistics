package service

import (
	"context"
	"gitlab.com/mediasoft-internship/final-task/contracts/pkg/contracts/statistics"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"math"
)

func (s *Service) GetAmountOfProfit(ctx context.Context, req *statistics.GetAmountOfProfitRequest) (*statistics.GetAmountOfProfitResponse, error) {
	startDate := req.StartDate.AsTime()
	endDate := req.EndDate.AsTime()

	products, err := s.orderRepository.ListOrderItems(ctx, startDate, endDate)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	restaurantProducts, err := s.restaurant.GetProductList(ctx)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	var profit float64 = 0
	for _, product := range products {
		for _, p := range restaurantProducts.Result {
			if product.ProductUuid.String() == p.Uuid {
				profit += float64(product.Count) * p.Price
			}
		}
	}
	profit = math.Round(profit*100) / 100
	return &statistics.GetAmountOfProfitResponse{
		Profit: profit,
	}, nil
}
