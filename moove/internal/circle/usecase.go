package circle

import (
	"github.com/iancoleman/strcase"
	"github.com/octopipe/charlescd/moove/internal/workspace"
	pbv1 "github.com/octopipe/charlescd/moove/pb/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

type UseCase struct {
	circleRepository   CircleRepository
	worksaceRepository workspace.WorkspaceRepository
}

func NewUseCase(workspaceRepository workspace.WorkspaceRepository, circleRepository CircleRepository) CircleUseCase {
	return UseCase{
		circleRepository:   circleRepository,
		worksaceRepository: workspaceRepository,
	}
}

// Create implements CircleUseCase
func (u UseCase) Create(circle *pbv1.CreateCircleRequest) (*pbv1.Circle, error) {
	_, err := u.circleRepository.Create(circle)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

// Delete implements CircleUseCase
func (u UseCase) Delete(workspaceId string, name string) error {
	workspace, err := u.worksaceRepository.FindById(workspaceId)
	if err != nil {
		return err
	}

	namespace := strcase.ToKebab(workspace.Name)
	err = u.circleRepository.Delete(name, namespace)
	if err != nil {
		return err
	}

	return nil
}

// FindAll implements CircleUseCase
func (u UseCase) FindAll(workspaceId string) ([]*pbv1.CircleMetadata, error) {
	workspace, err := u.worksaceRepository.FindById(workspaceId)
	if err != nil {
		return nil, err
	}
	namespace := workspace.Name
	filter := &pbv1.ListRequest{
		Namespace: strcase.ToKebab(namespace),
	}
	circles, err := u.circleRepository.FindAll(filter)
	if err != nil {
		return nil, err
	}

	return circles, nil
}

// FindByName implements CircleUseCase
func (u UseCase) FindByName(workspaceId string, name string) (*pbv1.Circle, error) {
	workspace, err := u.worksaceRepository.FindById(workspaceId)
	if err != nil {
		return nil, err
	}

	namespace := strcase.ToKebab(workspace.Name)
	circle, err := u.circleRepository.FindByName(namespace, name)
	if err != nil {
		return nil, err
	}

	return circle, nil
}

// GetDiagram implements CircleUseCase
func (u UseCase) GetDiagram(workspaceId string, name string) ([]*pbv1.Resource, error) {
	workspace, err := u.worksaceRepository.FindById(workspaceId)
	if err != nil {
		return nil, err
	}

	namespace := strcase.ToKebab(workspace.Name)
	items, err := u.circleRepository.GetDiagram(namespace, name)
	if err != nil {
		return nil, err
	}

	return items, nil
}

// GetEvents implements CircleUseCase
func (u UseCase) GetEvents(workspaceId string, resourceName string, kind string) ([]*pbv1.Event, error) {
	workspace, err := u.worksaceRepository.FindById(workspaceId)
	if err != nil {
		return nil, err
	}

	namespace := strcase.ToKebab(workspace.Name)
	events, err := u.circleRepository.GetEvents(namespace, resourceName, kind)
	return events, err
}

// GetLogs implements CircleUseCase
func (UseCase) GetLogs(circleName string, resourceName string, group string, kind string) (interface{}, error) {
	panic("unimplemented")
}

// GetResource implements CircleUseCase
func (u UseCase) GetResource(workspaceId string, resourceName string, group string, kind string) (*pbv1.Resource, *unstructured.Unstructured, error) {
	workspace, err := u.worksaceRepository.FindById(workspaceId)
	if err != nil {
		return nil, nil, err
	}

	namespace := strcase.ToKebab(workspace.Name)
	resource, manifest, err := u.circleRepository.GetResource(namespace, resourceName, group, kind)
	return resource, manifest, err
}

// Update implements CircleUseCase
func (u UseCase) Update(workspaceId string, name string, circle *pbv1.CreateCircleRequest) error {
	workspace, err := u.worksaceRepository.FindById(workspaceId)
	if err != nil {
		return err
	}

	namespace := strcase.ToKebab(workspace.Name)
	circle.Name = name
	circle.Namespace = namespace
	_, err = u.circleRepository.Update(circle)
	if err != nil {
		return err
	}

	return nil
}
