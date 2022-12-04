package module

import (
	"context"

	"github.com/octopipe/charlescd/internal/moove/core/listoptions"
	"github.com/octopipe/charlescd/internal/moove/workspace"
)

type UseCase struct {
	moduleRepository ModuleRepository
	workspaceUseCase workspace.WorkspaceUseCase
}

func NewUseCase(workspaceUseCase workspace.WorkspaceUseCase, moduleRepository ModuleRepository) ModuleUseCase {
	return UseCase{
		moduleRepository: moduleRepository,
		workspaceUseCase: workspaceUseCase,
	}
}

// Create implements ModuleUseCase
func (u UseCase) Create(ctx context.Context, workspaceId string, module Module) (ModuleModel, error) {
	namespace, err := u.workspaceUseCase.GetKebabCaseNameById(workspaceId)
	if err != nil {
		return ModuleModel{}, err
	}

	createdModule, err := u.moduleRepository.Create(ctx, namespace, module)
	if err != nil {
		return ModuleModel{}, err
	}

	return createdModule, nil
}

// Delete implements ModuleUseCase
func (u UseCase) Delete(ctx context.Context, workspaceId string, name string) error {
	namespace, err := u.workspaceUseCase.GetKebabCaseNameById(workspaceId)
	if err != nil {
		return err
	}

	err = u.moduleRepository.Delete(ctx, name, namespace)
	if err != nil {
		return err
	}

	return nil
}

// FindAll implements ModuleUseCase
func (u UseCase) FindAll(ctx context.Context, workspaceId string, options listoptions.Request) (listoptions.Response, error) {
	namespace, err := u.workspaceUseCase.GetKebabCaseNameById(workspaceId)
	if err != nil {
		return listoptions.Response{}, err
	}

	modules, err := u.moduleRepository.FindAll(ctx, namespace, options)
	if err != nil {
		return listoptions.Response{}, err
	}

	return modules, nil
}

// FindByName implements ModuleUseCase
func (u UseCase) FindById(ctx context.Context, workspaceId string, name string) (ModuleModel, error) {
	namespace, err := u.workspaceUseCase.GetKebabCaseNameById(workspaceId)
	if err != nil {
		return ModuleModel{}, err
	}

	module, err := u.moduleRepository.FindById(ctx, namespace, name)
	if err != nil {
		return ModuleModel{}, err
	}

	return module, nil
}

// Update implements ModuleUseCase
func (u UseCase) Update(ctx context.Context, workspaceId string, name string, module Module) (ModuleModel, error) {
	namespace, err := u.workspaceUseCase.GetKebabCaseNameById(workspaceId)
	if err != nil {
		return ModuleModel{}, err
	}

	updatedModule, err := u.moduleRepository.Update(ctx, namespace, name, module)
	if err != nil {
		return ModuleModel{}, err
	}

	return updatedModule, nil
}
