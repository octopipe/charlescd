package circle

import (
	"fmt"

	"github.com/iancoleman/strcase"
	"github.com/octopipe/charlescd/moove/internal/workspace"
	pbv1 "github.com/octopipe/charlescd/moove/pb/v1"
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
func (UseCase) Create(circle Circle) (*pbv1.Circle, error) {
	panic("unimplemented")
}

// Delete implements CircleUseCase
func (UseCase) Delete(id string) error {
	panic("unimplemented")
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
	fmt.Println(filter)
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

	circle, err := u.circleRepository.FindByName(workspace.Name, name)
	if err != nil {
		return nil, err
	}

	return circle, nil
}

// GetDiagram implements CircleUseCase
func (UseCase) GetDiagram(circleName string) (interface{}, error) {
	panic("unimplemented")
}

// GetEvents implements CircleUseCase
func (UseCase) GetEvents(circleName string, resourceName string, group string, kind string) (interface{}, error) {
	panic("unimplemented")
}

// GetLogs implements CircleUseCase
func (UseCase) GetLogs(circleName string, resourceName string, group string, kind string) (interface{}, error) {
	panic("unimplemented")
}

// GetResource implements CircleUseCase
func (UseCase) GetResource(circleName string, resourceName string, group string, kind string) (interface{}, error) {
	panic("unimplemented")
}

// Update implements CircleUseCase
func (UseCase) Update(id string, circle Circle) (CircleProvider, error) {
	panic("unimplemented")
}
