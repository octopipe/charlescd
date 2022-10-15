package circle

import (
	"encoding/json"
	"fmt"

	"github.com/imroc/req/v3"
)

type HttpRepository struct {
	httpClient *req.Client
}

func NewRepository(httpClient *req.Client) CircleRepository {
	return HttpRepository{httpClient: httpClient}
}

// Create implements WorkspaceRepository
func (r HttpRepository) Create(circle Circle) (CircleProvider, error) {
	res, err := r.httpClient.R().
		SetBody(circle).
		Post("/circles")

	if err != nil {
		fmt.Errorf(err.Error())
		return CircleProvider{}, err
	}

	circleProvider := new(CircleProvider)
	err = json.Unmarshal(res.Bytes(), &circleProvider)
	if err != nil {
		return CircleProvider{}, err
	}

	return *circleProvider, err
}

// Delete implements WorkspaceRepository
func (r HttpRepository) Delete(id string) error {

	return nil
}

// FindAll implements WorkspaceRepository
func (r HttpRepository) FindAll() ([]CircleProvider, error) {
	res, err := r.httpClient.R().
		Get("/circles")

	if err != nil {
		fmt.Errorf("fail to find all circles", err.Error())
		return nil, err
	}

	circleProviders := new([]CircleProvider)
	err = json.Unmarshal(res.Bytes(), &circleProviders)
	if err != nil {
		return nil, err
	}

	return *circleProviders, err
}

// FindById implements WorkspaceRepository
func (r HttpRepository) FindById(id string) (CircleProvider, error) {
	return CircleProvider{}, nil
}

// Update implements WorkspaceRepository
func (r HttpRepository) Update(id string, workspace Circle) (CircleProvider, error) {
	return CircleProvider{}, nil
}

func (r HttpRepository) GetDiagram(circleName string) (interface{}, error) {
	res, err := r.httpClient.R().
		Get(fmt.Sprintf("/circles/%s/diagram", circleName))

	if err != nil {
		fmt.Errorf("fail to find all circles", err.Error())
		return nil, err
	}

	diagram := new(interface{})
	err = json.Unmarshal(res.Bytes(), &diagram)
	if err != nil {
		return nil, err
	}

	return *diagram, err
}

// GetEvents implements CircleRepository
func (r HttpRepository) GetEvents(circleName string, resourceName string, group string, kind string) (interface{}, error) {
	res, err := r.httpClient.R().
		Get(fmt.Sprintf("/circles/%s/resources/%s/events?group=%s&kind=%s", circleName, resourceName, group, kind))

	if err != nil {
		fmt.Errorf("fail to find all circles", err.Error())
		return nil, err
	}

	diagram := new(interface{})
	err = json.Unmarshal(res.Bytes(), &diagram)
	if err != nil {
		return nil, err
	}

	return *diagram, err
}

// GetLogs implements CircleRepository
func (r HttpRepository) GetLogs(circleName string, resourceName string, group string, kind string) (interface{}, error) {
	res, err := r.httpClient.R().
		Get(fmt.Sprintf("/circles/%s/resources/%s/logs?group=%s&kind=%s", circleName, resourceName, group, kind))

	if err != nil {
		fmt.Errorf("fail to find all circles", err.Error())
		return nil, err
	}

	diagram := new(interface{})
	err = json.Unmarshal(res.Bytes(), &diagram)
	if err != nil {
		return nil, err
	}

	return *diagram, err
}

// GetResource implements CircleRepository
func (r HttpRepository) GetResource(circleName string, resourceName string, group string, kind string) (interface{}, error) {
	res, err := r.httpClient.R().
		Get(fmt.Sprintf("/circles/%s/resources/%s?group=%s&kind=%s", circleName, resourceName, group, kind))

	if err != nil {
		fmt.Errorf("fail to find all circles", err.Error())
		return nil, err
	}

	diagram := new(interface{})
	err = json.Unmarshal(res.Bytes(), &diagram)
	if err != nil {
		return nil, err
	}

	return *diagram, err
}
