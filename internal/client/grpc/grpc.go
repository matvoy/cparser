package grpc

import (
	"context"

	txv1beta1 "github.com/cosmos/cosmos-sdk/api/cosmos/tx/v1beta1"
	"github.com/matvoy/cparser/internal/models"
	"google.golang.org/grpc"
)

type grpcClient struct {
	conn          *grpc.ClientConn
	serviceClient txv1beta1.ServiceClient
}

func NewClient(url string) (*grpcClient, error) {
	if len(url) == 0 {
		return nil, nil
	}
	conn, err := grpc.Dial(url, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	client := txv1beta1.NewServiceClient(conn)

	return &grpcClient{
		conn,
		client,
	}, nil
}

func (c *grpcClient) Close() {
	c.conn.Close()
}

func (c *grpcClient) GetBlockData(height uint64) (*models.Block, error) {
	if height == 0 {
		return nil, nil
	}

	resp, err := c.serviceClient.GetBlockWithTxs(context.Background(), &txv1beta1.GetBlockWithTxsRequest{
		Height: int64(height),
		// Pagination: &queryv1beta1.PageRequest{
		// 	Limit: 1,
		// },
	})
	if err != nil {
		return nil, err
	}
	_ = resp

	// TO DO
	//
	// I need GRPC endpoint.
	// If http client is not enough I can continue work on grpc, but I need endpoint

	return nil, nil
}
