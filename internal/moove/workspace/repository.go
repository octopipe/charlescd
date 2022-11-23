package workspace

import (
	"context"

	"github.com/iancoleman/strcase"
	"gorm.io/gorm"
	v1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type GormRepository struct {
	db        *gorm.DB
	clientset client.Client
}

func NewRepository(db *gorm.DB, clientset client.Client) WorkspaceRepository {
	return GormRepository{db: db, clientset: clientset}
}

// Create implements WorkspaceRepository
func (r GormRepository) Create(workspace Workspace) (WorkspaceModel, error) {
	newNamespace := v1.Namespace{}
	newNamespace.SetName(strcase.ToKebab(workspace.Name))

	err := r.clientset.Create(context.Background(), &newNamespace)
	if err != nil {
		return WorkspaceModel{}, err
	}

	workspaceModel := WorkspaceModel{Workspace: workspace}
	res := r.db.Table("workspaces").Save(&workspaceModel)
	if res.Error != nil {
		return WorkspaceModel{}, res.Error
	}

	return workspaceModel, nil
}

// Delete implements WorkspaceRepository
func (r GormRepository) Delete(id string) error {
	res := r.db.Delete(id)
	if res.Error != nil {
		return res.Error
	}

	return nil
}

// FindAll implements WorkspaceRepository
func (r GormRepository) FindAll() ([]WorkspaceModel, error) {
	workspaceModels := []WorkspaceModel{}
	res := r.db.Table("workspaces").Find(&workspaceModels)
	if res.Error != nil {
		return nil, res.Error
	}

	return workspaceModels, nil
}

// FindById implements WorkspaceRepository
func (r GormRepository) FindById(id string) (WorkspaceModel, error) {
	workspaceModel := new(WorkspaceModel)
	res := r.db.Table("workspaces").First(&workspaceModel, map[string]string{"id": id})
	if res.Error != nil {
		return WorkspaceModel{}, res.Error
	}

	return *workspaceModel, nil
}

// Update implements WorkspaceRepository
func (r GormRepository) Update(id string, workspace Workspace) (WorkspaceModel, error) {
	workspaceModel := WorkspaceModel{Workspace: workspace}
	res := r.db.Where("id = ?", id).Updates(&workspace)
	if res.Error != nil {
		return WorkspaceModel{}, res.Error
	}

	return workspaceModel, nil
}
