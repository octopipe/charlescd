package circle

import (
	"context"
	"encoding/json"

	"github.com/octopipe/charlescd/moove/internal/core/grpcclient"
	pbv1 "github.com/octopipe/charlescd/moove/pb/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
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
	circle, err := r.grpcClient.CircleClient.Get(context.Background(), &pbv1.GetCircleRequest{
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

func (r GrpcRepository) GetDiagram(namespace string, name string) ([]*pbv1.Resource, error) {
	hierarchy, err := r.grpcClient.ResourceClient.Hierarchy(context.Background(), &pbv1.HierarchyRequest{
		Name:      name,
		Namespace: namespace,
	})
	if err != nil {
		return nil, err
	}

	return hierarchy.Items, nil
}

// GetEvents implements CircleRepository
func (r GrpcRepository) GetEvents(namespace string, resourceName string, kind string) ([]*pbv1.Event, error) {
	eventsResponse, err := r.grpcClient.ResourceClient.Events(context.Background(), &pbv1.EventsRequest{
		Name:      resourceName,
		Namespace: namespace,
		Kind:      kind,
	})
	if err != nil {
		return nil, err
	}

	return eventsResponse.Items, nil
}

// GetLogs implements CircleRepository
func (r GrpcRepository) GetLogs(circleName string, resourceName string, group string, kind string) (interface{}, error) {
	return nil, nil
}

// GetResource implements CircleRepository
func (r GrpcRepository) GetResource(namespace string, resourceName string, group string, kind string) (*pbv1.Resource, *unstructured.Unstructured, error) {
	resource, err := r.grpcClient.ResourceClient.Get(context.Background(), &pbv1.GetResourceRequest{
		Namespace: namespace,
		Name:      resourceName,
		Group:     group,
		Kind:      kind,
	})
	if err != nil {
		return nil, nil, err
	}

	manifest := &unstructured.Unstructured{}
	err = json.Unmarshal(resource.Manifest, manifest)
	if err != nil {
		return nil, nil, err
	}

	return resource.Metadata, manifest, nil
}
