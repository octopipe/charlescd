package circle

import (
	"context"

	"github.com/octopipe/charlescd/moove/internal/core/grpcclient"
	pbv1 "github.com/octopipe/charlescd/moove/pb/v1"
)

type GrpcRepository struct {
	grpcClient grpcclient.Client
}

func NewRepository(grpcClient grpcclient.Client) CircleRepository {
	return GrpcRepository{grpcClient: grpcClient}
}

// Create implements WorkspaceRepository
func (r GrpcRepository) Create(circle Circle) (*pbv1.Circle, error) {
	return nil, nil

}

// Delete implements WorkspaceRepository
func (r GrpcRepository) Delete(id string) error {

	return nil
}

// FindAll implements WorkspaceRepository
func (r GrpcRepository) FindAll(filter *pbv1.ListRequest) ([]*pbv1.CircleMetadata, error) {
	res, err := r.grpcClient.CircleClient.List(context.Background(), filter)
	if err != nil {
		return nil, err
	}

	return res.Items, err
}

// FindById implements WorkspaceRepository
func (r GrpcRepository) FindByName(namespace string, name string) (*pbv1.Circle, error) {
	circle, err := r.grpcClient.CircleClient.Get(context.Background(), &pbv1.GetRequest{
		Namespace: namespace,
		Name:      name,
	})
	if err != nil {
		return nil, err
	}

	return circle, nil
}

// Update implements WorkspaceRepository
func (r GrpcRepository) Update(id string, workspace Circle) (CircleProvider, error) {
	return CircleProvider{}, nil
}

func (r GrpcRepository) GetDiagram(circleName string) (interface{}, error) {

	return nil, nil
}

// GetEvents implements CircleRepository
func (r GrpcRepository) GetEvents(circleName string, resourceName string, group string, kind string) (interface{}, error) {
	return nil, nil
}

// GetLogs implements CircleRepository
func (r GrpcRepository) GetLogs(circleName string, resourceName string, group string, kind string) (interface{}, error) {
	return nil, nil
}

// GetResource implements CircleRepository
func (r GrpcRepository) GetResource(circleName string, resourceName string, group string, kind string) (interface{}, error) {
	return nil, nil
}
