package client

import (
	"context"
	"fmt"
	"gitlab.com/mediasoft-internship/final-task/contracts/pkg/contracts/restaurant"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"mediasoft-statistics/internal/config"
)

type RestaurantClient struct {
	client restaurant.ProductServiceClient
}

func New(cfg config.Config) (*RestaurantClient, error) {
	conn, err := grpc.Dial(fmt.Sprintf("%s:%d", cfg.RestaurantGRPC.IP, cfg.RestaurantGRPC.Port),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, err
	}
	client := restaurant.NewProductServiceClient(conn)

	return &RestaurantClient{client: client}, err
}

func (r *RestaurantClient) Close() error {
	return r.Close()
}

func (r *RestaurantClient) GetProductList(ctx context.Context) (*restaurant.GetProductListResponse, error) {
	res, err := r.client.GetProductList(ctx, &restaurant.GetProductListRequest{})
	if err != nil {
		return nil, err
	}
	return res, err
}
