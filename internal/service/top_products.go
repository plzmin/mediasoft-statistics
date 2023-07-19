package service

import (
	"context"
	"gitlab.com/mediasoft-internship/final-task/contracts/pkg/contracts/statistics"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"sort"
)

func (s *Service) TopProducts(ctx context.Context,
	req *statistics.TopProductsRequest) (*statistics.TopProductsResponse, error) {
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

	var result []*statistics.Product
	for _, product := range products {
		for _, p := range restaurantProducts.Result {
			if product.ProductUuid.String() == p.Uuid {
				result = append(result, &statistics.Product{
					Uuid:        p.Uuid,
					Name:        p.Name,
					Count:       product.Count,
					ProductType: statistics.StatisticsProductType(p.Type),
				})
			}
		}
	}

	topProducts := make(map[string]int64)
	if req.ProductType != nil {
		for _, el := range result {
			if el.ProductType == *req.ProductType {
				topProducts[el.Uuid] += el.Count
			}
		}
	} else {
		for _, el := range result {
			topProducts[el.Uuid] += el.Count
		}
	}

	sort.Slice(result, func(i, j int) bool {
		return result[i].Count > result[j].Count
	})

	return &statistics.TopProductsResponse{Result: result[:3]}, nil
}
