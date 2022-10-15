package circle

type UseCase struct {
	repository CircleRepository
}

func NewUseCase(repository CircleRepository) CircleUseCase {
	return UseCase{
		repository: repository,
	}
}

// Create implements CircleUseCase
func (u UseCase) Create(Circle Circle) (CircleProvider, error) {
	circleProvider, err := u.repository.Create(Circle)
	if err != nil {
		return CircleProvider{}, err
	}

	return circleProvider, nil
}

// Delete implements CircleUseCase
func (u UseCase) Delete(id string) error {
	err := u.repository.Delete(id)
	if err != nil {
		return err
	}

	return nil
}

// FindAll implements CircleUseCase
func (u UseCase) FindAll() ([]CircleProvider, error) {
	circleProviders, err := u.repository.FindAll()
	if err != nil {
		return []CircleProvider{}, err
	}

	return circleProviders, nil
}

// FindById implements CircleUseCase
func (u UseCase) FindById(id string) (CircleProvider, error) {
	circleProvider, err := u.repository.FindById(id)
	if err != nil {
		return CircleProvider{}, err
	}

	return circleProvider, nil
}

// Update implements CircleUseCase
func (u UseCase) Update(id string, Circle Circle) (CircleProvider, error) {
	circleProvider, err := u.repository.Update(id, Circle)
	if err != nil {
		return CircleProvider{}, err
	}

	return circleProvider, nil
}

func (u UseCase) GetDiagram(circleName string) (interface{}, error) {
	return u.repository.GetDiagram(circleName)
}

// GetEvents implements CircleUseCase
func (u UseCase) GetEvents(circleName string, resourceName string, group string, kind string) (interface{}, error) {
	return u.repository.GetEvents(circleName, resourceName, group, kind)
}

// GetLogs implements CircleUseCase
func (u UseCase) GetLogs(circleName string, resourceName string, group string, kind string) (interface{}, error) {
	return u.repository.GetLogs(circleName, resourceName, group, kind)
}

// GetResource implements CircleUseCase
func (u UseCase) GetResource(circleName string, resourceName string, group string, kind string) (interface{}, error) {
	return u.repository.GetResource(circleName, resourceName, group, kind)
}
