package resource

import (
	"context"

	"github.com/octopipe/charlescd/moove/internal/workspace"
)

type UseCase struct {
	resourceRepository ResourceRepository
	workspaceUseCase   workspace.WorkspaceUseCase
}

func NewUseCase(workspaceUseCase workspace.WorkspaceUseCase, resourceRepository ResourceRepository) ResourceUseCase {
	return UseCase{
		resourceRepository: resourceRepository,
		workspaceUseCase:   workspaceUseCase,
	}
}

// GetDiagram implements CircleUseCase
func (u UseCase) GetTree(ctx context.Context, workspaceId string, name string) ([]Resource, error) {
	namespace, err := u.workspaceUseCase.GetKebabCaseNameById(workspaceId)
	if err != nil {
		return nil, err
	}

	items, err := u.resourceRepository.GetTree(ctx, namespace, name)
	if err != nil {
		return nil, err
	}

	return items, nil
}

// GetEvents implements CircleUseCase
func (u UseCase) GetEvents(ctx context.Context, workspaceId string, resourceName string, kind string) ([]ResourceEvent, error) {
	namespace, err := u.workspaceUseCase.GetKebabCaseNameById(workspaceId)
	if err != nil {
		return nil, err
	}

	events, err := u.resourceRepository.GetEvents(ctx, namespace, resourceName, kind)
	return events, err
}

// GetLogs implements CircleUseCase
func (UseCase) GetLogs(ctx context.Context, circleName string, resourceName string, group string, kind string) (interface{}, error) {
	return nil, nil
}

// GetResource implements CircleUseCase
func (u UseCase) GetResource(ctx context.Context, workspaceId string, resourceName string, group string, kind string) (Resource, error) {
	namespace, err := u.workspaceUseCase.GetKebabCaseNameById(workspaceId)
	if err != nil {
		return Resource{}, err
	}

	resource, err := u.resourceRepository.GetResource(ctx, namespace, resourceName, group, kind)
	return resource, err
}
