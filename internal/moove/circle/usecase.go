package circle

import (
	"context"

	"github.com/octopipe/charlescd/internal/moove/core/listoptions"
	"github.com/octopipe/charlescd/internal/moove/workspace"
	"github.com/octopipe/charlescd/internal/utils/id"
	pbv1 "github.com/octopipe/charlescd/pb/v1"
)

type UseCase struct {
	circleProvider   CircleProvider
	circleRepository CircleRepository
	workspaceUseCase workspace.WorkspaceUseCase
}

func NewUseCase(workspaceUseCase workspace.WorkspaceUseCase, circleProvider CircleProvider, circleRepository CircleRepository) CircleUseCase {
	return UseCase{
		circleProvider:   circleProvider,
		circleRepository: circleRepository,
		workspaceUseCase: workspaceUseCase,
	}
}

// Create implements CircleUseCase
func (u UseCase) Create(ctx context.Context, workspaceId string, circle Circle) (CircleModel, error) {
	namespace, err := u.workspaceUseCase.GetKebabCaseNameById(workspaceId)
	if err != nil {
		return CircleModel{}, err
	}

	createdCircle, err := u.circleRepository.Create(ctx, namespace, circle)
	if err != nil {
		return CircleModel{}, err
	}

	return createdCircle, nil
}

// Delete implements CircleUseCase
func (u UseCase) Delete(ctx context.Context, workspaceId string, circleId string) error {
	namespace, err := u.workspaceUseCase.GetKebabCaseNameById(workspaceId)
	if err != nil {
		return err
	}

	err = u.circleRepository.Delete(ctx, namespace, circleId)
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
func (u UseCase) FindById(ctx context.Context, workspaceId string, circleId string) (CircleModel, error) {
	namespace, err := u.workspaceUseCase.GetKebabCaseNameById(workspaceId)
	if err != nil {
		return CircleModel{}, err
	}

	circle, err := u.circleRepository.FindById(ctx, namespace, circleId)
	if err != nil {
		return CircleModel{}, err
	}

	return circle, nil
}

func (u UseCase) Sync(ctx context.Context, workspaceId string, circleId string) error {
	namespace, err := u.workspaceUseCase.GetKebabCaseNameById(workspaceId)
	if err != nil {
		return err
	}

	_, err = u.circleRepository.FindById(ctx, namespace, circleId)
	if err != nil {
		return err
	}

	name, err := id.DecodeID(circleId)
	if err != nil {
		return err
	}

	err = u.circleProvider.Sync(ctx, namespace, name)
	return err
}

// Update implements CircleUseCase
func (u UseCase) Update(ctx context.Context, workspaceId string, circleId string, circle Circle) (CircleModel, error) {
	namespace, err := u.workspaceUseCase.GetKebabCaseNameById(workspaceId)
	if err != nil {
		return CircleModel{}, err
	}

	updatedCircle, err := u.circleRepository.Update(ctx, namespace, circleId, circle)
	if err != nil {
		return CircleModel{}, err
	}

	return updatedCircle, nil
}

func (u UseCase) Status(ctx context.Context, workspaceId string, circleId string) (*pbv1.StatusResponse, error) {
	namespace, err := u.workspaceUseCase.GetKebabCaseNameById(workspaceId)
	if err != nil {
		return nil, err
	}

	_, err = u.circleRepository.FindById(ctx, namespace, circleId)
	if err != nil {
		return nil, err
	}

	name, err := id.DecodeID(circleId)
	if err != nil {
		return nil, err
	}

	status, err := u.circleProvider.Status(ctx, namespace, name)
	return status, err
}
