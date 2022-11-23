package workspace

import (
	"gorm.io/gorm"
)

type GormRepository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) WorkspaceRepository {
	return GormRepository{db: db}
}

// Create implements WorkspaceRepository
func (r GormRepository) Create(workspace Workspace) (WorkspaceModel, error) {
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
