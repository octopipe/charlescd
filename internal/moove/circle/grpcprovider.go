package circle

import (
	"context"

	"github.com/octopipe/charlescd/internal/moove/core/grpcclient"
	pbv1 "github.com/octopipe/charlescd/pb/v1"
)

type grpcProvider struct {
	grpcClient grpcclient.Client
}

func NewGrpcProvider(grpcClient grpcclient.Client) CircleProvider {
	return grpcProvider{grpcClient: grpcClient}
}

func (p grpcProvider) Sync(ctx context.Context, namespace string, name string) error {
	req := &pbv1.GetCircle{
		CircleName:      name,
		CircleNamespace: namespace,
	}
	_, err := p.grpcClient.CircleClient.Sync(ctx, req)
	return err
}

func (p grpcProvider) Status(ctx context.Context, namespace string, name string) (*pbv1.StatusResponse, error) {
	req := &pbv1.GetCircle{
		CircleName:      name,
		CircleNamespace: namespace,
	}
	circleStatus, err := p.grpcClient.CircleClient.Status(ctx, req)
	return circleStatus, err
}
