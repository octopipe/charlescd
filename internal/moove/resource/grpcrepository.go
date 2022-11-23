package resource

import (
	"context"

	"github.com/octopipe/charlescd/internal/moove/core/grpcclient"
	"github.com/octopipe/charlescd/internal/moove/errs"
	pbv1 "github.com/octopipe/charlescd/pb/v1"
)

type GrpcRepository struct {
	grpcClient grpcclient.Client
}

func NewRepository(grpcClient grpcclient.Client) GrpcRepository {
	return GrpcRepository{grpcClient: grpcClient}
}

func (r GrpcRepository) messageToResource(resourceMessage *pbv1.Resource) Resource {

	owner := ResourceOwner{}
	if resourceMessage.Owner != nil {
		owner.Name = resourceMessage.Name
		owner.Kind = resourceMessage.Kind
	}

	return Resource{
		Name:      resourceMessage.Name,
		Namespace: resourceMessage.Namespace,
		Kind:      resourceMessage.Kind,
		Group:     resourceMessage.Group,
		Owner:     owner,
		Status:    resourceMessage.Status,
		Message:   resourceMessage.Error,
	}
}

func (r GrpcRepository) GetTree(ctx context.Context, namespace string, name string) ([]Resource, error) {
	tree, err := r.grpcClient.ResourceClient.Tree(ctx, &pbv1.TreeRequest{
		CircleName:      name,
		CircleNamespace: namespace,
	})
	if err != nil {
		return nil, errs.ParseGrpcError(err)
	}

	resources := []Resource{}
	for _, i := range tree.Items {
		resources = append(resources, r.messageToResource(i))
	}

	return resources, nil
}

// GetEvents implements CircleRepository
func (r GrpcRepository) GetEvents(ctx context.Context, namespace string, resourceName string, kind string) ([]ResourceEvent, error) {
	eventsResponse, err := r.grpcClient.ResourceClient.Events(ctx, &pbv1.EventsRequest{
		Name:      resourceName,
		Namespace: namespace,
		Kind:      kind,
	})
	if err != nil {
		return nil, errs.ParseGrpcError(err)
	}

	resourceEvents := []ResourceEvent{}
	for _, i := range eventsResponse.Items {
		resourceEvents = append(resourceEvents, ResourceEvent{
			Reason:  i.Reason,
			Message: i.Message,
			Count:   i.Count,
			Type:    i.Type,
			Action:  i.Action,
		})
	}

	return resourceEvents, nil
}

// GetLogs implements CircleRepository
func (r GrpcRepository) GetLogs(ctx context.Context, circleName string, resourceName string, group string, kind string) (interface{}, error) {
	return nil, nil
}

// GetResource implements CircleRepository
func (r GrpcRepository) GetResource(ctx context.Context, namespace string, resourceName string, group string, kind string) (Resource, error) {
	resourceMessage, err := r.grpcClient.ResourceClient.Get(ctx, &pbv1.GetResourceRequest{
		Namespace: namespace,
		Name:      resourceName,
		Group:     group,
		Kind:      kind,
	})
	if err != nil {
		return Resource{}, errs.ParseGrpcError(err)
	}

	return r.messageToResource(resourceMessage), nil
}
