package workspace

import (
	"github.com/iancoleman/strcase"
)

type UseCase struct {
	repository WorkspaceRepository
}

func NewUseCase(repository WorkspaceRepository) WorkspaceUseCase {
	return UseCase{
		repository: repository,
	}
}

// Create implements WorkspaceModelUseCase
func (u UseCase) Create(workspace Workspace) (WorkspaceModel, error) {
	workspaceModel, err := u.repository.Create(workspace)
	if err != nil {
		return WorkspaceModel{}, err
	}

	return workspaceModel, nil
}

// Delete implements WorkspaceModelUseCase
func (u UseCase) Delete(id string) error {
	err := u.repository.Delete(id)
	if err != nil {
		return err
	}

	return nil
}

// FindAll implements WorkspaceModelUseCase
func (u UseCase) FindAll() ([]WorkspaceModel, error) {

	workspaceModels, err := u.repository.FindAll()
	if err != nil {
		return []WorkspaceModel{}, err
	}

	return workspaceModels, nil
}

// FindById implements WorkspaceModelUseCase
func (u UseCase) FindById(id string) (WorkspaceModel, error) {
	workspaceModel, err := u.repository.FindById(id)
	if err != nil {
		return WorkspaceModel{}, err
	}

	return workspaceModel, nil
}

// Update implements WorkspaceModelUseCase
func (u UseCase) Update(id string, workspace Workspace) (WorkspaceModel, error) {
	workspaceModel, err := u.repository.Update(id, workspace)
	if err != nil {
		return WorkspaceModel{}, err
	}

	return workspaceModel, nil
}

func (u UseCase) GetKebabCaseNameById(id string) (string, error) {
	workspace, err := u.repository.FindById(id)
	if err != nil {
		return "", err
	}

	return strcase.ToKebab(workspace.Name), nil
}
