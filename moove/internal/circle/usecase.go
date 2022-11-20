package circle

import (
	"context"

	"github.com/octopipe/charlescd/moove/internal/core/listoptions"
	"github.com/octopipe/charlescd/moove/internal/workspace"
)

type UseCase struct {
	circleRepository CircleRepository
	workspaceUseCase workspace.WorkspaceUseCase
}

func NewUseCase(workspaceUseCase workspace.WorkspaceUseCase, circleRepository CircleRepository) CircleUseCase {
	return UseCase{
		circleRepository: circleRepository,
		workspaceUseCase: workspaceUseCase,
	}
}

// Create implements CircleUseCase
func (u UseCase) Create(ctx context.Context, workspaceId string, circle Circle) (Circle, error) {
	namespace, err := u.workspaceUseCase.GetKebabCaseNameById(workspaceId)
	if err != nil {
		return Circle{}, err
	}

	createdCircle, err := u.circleRepository.Create(ctx, namespace, circle)
	if err != nil {
		return Circle{}, err
	}

	return createdCircle, nil
}

// Delete implements CircleUseCase
func (u UseCase) Delete(ctx context.Context, workspaceId string, name string) error {
	namespace, err := u.workspaceUseCase.GetKebabCaseNameById(workspaceId)
	if err != nil {
		return err
	}

	err = u.circleRepository.Delete(ctx, name, namespace)
	if err != nil {
		return err
	}

	return nil
}

// FindAll implements CircleUseCase
func (u UseCase) FindAll(ctx context.Context, workspaceId string, options listoptions.Request) (listoptions.Response, error) {
	namespace, err := u.workspaceUseCase.GetKebabCaseNameById(workspaceId)
	if err != nil {
		return listoptions.Response{}, err
	}

	circles, err := u.circleRepository.FindAll(ctx, namespace, options)
	if err != nil {
		return listoptions.Response{}, err
	}

	return circles, nil
}

// FindByName implements CircleUseCase
func (u UseCase) FindByName(ctx context.Context, workspaceId string, name string) (Circle, error) {
	namespace, err := u.workspaceUseCase.GetKebabCaseNameById(workspaceId)
	if err != nil {
		return Circle{}, err
	}

	circle, err := u.circleRepository.FindByName(ctx, namespace, name)
	if err != nil {
		return Circle{}, err
	}

	return circle, nil
}

// Update implements CircleUseCase
func (u UseCase) Update(ctx context.Context, workspaceId string, name string, circle Circle) (Circle, error) {
	namespace, err := u.workspaceUseCase.GetKebabCaseNameById(workspaceId)
	if err != nil {
		return Circle{}, err
	}

	updatedCircle, err := u.circleRepository.Update(ctx, namespace, name, circle)
	if err != nil {
		return Circle{}, err
	}

	return updatedCircle, nil
}
