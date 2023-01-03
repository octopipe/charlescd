package metric

import (
	"context"
	"fmt"
	"time"

	"github.com/octopipe/charlescd/internal/moove/circle"
	"github.com/octopipe/charlescd/internal/moove/core/listoptions"
	"github.com/octopipe/charlescd/internal/moove/workspace"
	"github.com/octopipe/charlescd/internal/utils/id"
	charlescdiov1alpha1 "github.com/octopipe/charlescd/pkg/api/v1alpha1"
)

type UseCase struct {
	metricProvider   MetricProvider
	metricRepository MetricRepository
	workspaceUseCase workspace.WorkspaceUseCase
	circleUseCase    circle.CircleUseCase
}

// Query implements MetricUseCase

func NewUseCase(workspaceUseCase workspace.WorkspaceUseCase, circleUseCase circle.CircleUseCase, metricProvider MetricProvider, metricRepository MetricRepository) MetricUseCase {
	return UseCase{
		metricProvider:   metricProvider,
		metricRepository: metricRepository,
		workspaceUseCase: workspaceUseCase,
		circleUseCase:    circleUseCase,
	}
}

// Create implements MetricUseCase
func (u UseCase) Create(ctx context.Context, workspaceId string, circleId string, circle Metric) (MetricModel, error) {
	namespace, err := u.workspaceUseCase.GetKebabCaseNameById(workspaceId)
	if err != nil {
		return MetricModel{}, err
	}

	circleModel, err := u.circleUseCase.FindById(ctx, workspaceId, circleId)
	if err != nil {
		return MetricModel{}, err
	}

	createdMetric, err := u.metricRepository.Create(ctx, namespace, circleModel, circle)
	if err != nil {
		return MetricModel{}, err
	}

	return createdMetric, nil
}

// Delete implements MetricUseCase
func (u UseCase) Delete(ctx context.Context, workspaceId string, metricId string) error {
	namespace, err := u.workspaceUseCase.GetKebabCaseNameById(workspaceId)
	if err != nil {
		return err
	}

	err = u.metricRepository.Delete(ctx, namespace, metricId)
	if err != nil {
		return err
	}

	return nil
}

func (u UseCase) basicMatrics(circleModel circle.CircleModel) (map[string]MetricModel, error) {
	circleName, err := id.DecodeID(circleModel.ID)
	if err != nil {
		return nil, err
	}

	return map[string]MetricModel{
		"cpu": {
			ID: id.ToID("cpu"),
			Metric: Metric{
				Name:     "CPU",
				Category: BasicMetricCategory,
				Provider: PrometheusMetricProvider,
				MetricSpec: charlescdiov1alpha1.MetricSpec{
					Query: fmt.Sprintf(`container_cpu_usage_seconds_total{pod=~".*%s.*"}`, circleName),
				},
			},
			CreatedAt: circleModel.CreatedAt,
		},
		"memory": {
			ID: id.ToID("memory"),
			Metric: Metric{
				Name:     "Memory",
				Category: BasicMetricCategory,
				Provider: PrometheusMetricProvider,
				MetricSpec: charlescdiov1alpha1.MetricSpec{
					Query: fmt.Sprintf(`container_memory_usage_bytes{pod=~".*%s.*"}`, circleName),
				},
			},
			CreatedAt: circleModel.CreatedAt,
		},
	}, nil
}

func (u UseCase) getBasicMetrics(circleModel circle.CircleModel) ([]MetricModel, error) {
	metricModels := []MetricModel{}

	basicMetrics, err := u.basicMatrics(circleModel)
	if err != nil {
		return nil, err
	}

	for _, m := range basicMetrics {
		metricModels = append(metricModels, m)
	}

	return metricModels, nil
}

// FindAll implements MetricUseCase
func (u UseCase) FindAll(ctx context.Context, workspaceId string, circleId string, options listoptions.Request) (listoptions.Response, error) {
	namespace, err := u.workspaceUseCase.GetKebabCaseNameById(workspaceId)
	if err != nil {
		return listoptions.Response{}, err
	}

	circleModel, err := u.circleUseCase.FindById(ctx, workspaceId, circleId)
	if err != nil {
		return listoptions.Response{}, err
	}

	list, err := u.metricRepository.FindAll(ctx, namespace, circleModel, options)
	if err != nil {
		return listoptions.Response{}, err
	}

	basicMetrics, err := u.getBasicMetrics(circleModel)
	if err != nil {
		return listoptions.Response{}, err
	}

	metrics := list.Items.([]MetricModel)
	metrics = append(metrics, basicMetrics...)
	list.Items = metrics

	return list, nil
}

// FindByName implements MetricUseCase
func (u UseCase) FindById(ctx context.Context, workspaceId string, metricId string) (MetricModel, error) {
	namespace, err := u.workspaceUseCase.GetKebabCaseNameById(workspaceId)
	if err != nil {
		return MetricModel{}, err
	}

	circle, err := u.metricRepository.FindById(ctx, namespace, metricId)
	if err != nil {
		return MetricModel{}, err
	}

	return circle, nil
}

// Update implements MetricUseCase
func (u UseCase) Update(ctx context.Context, workspaceId string, circleId string, metricId string, circle Metric) (MetricModel, error) {
	namespace, err := u.workspaceUseCase.GetKebabCaseNameById(workspaceId)
	if err != nil {
		return MetricModel{}, err
	}

	circleModel, err := u.circleUseCase.FindById(ctx, workspaceId, circleId)
	if err != nil {
		return MetricModel{}, err
	}

	updatedMetric, err := u.metricRepository.Update(ctx, namespace, circleModel, metricId, circle)
	if err != nil {
		return MetricModel{}, err
	}

	return updatedMetric, nil
}

func getMetricRange(rangeTime string) MetricRange {
	timeToAdd := (time.Minute) * 3
	step := time.Minute

	if rangeTime == FiveMinutes {
		timeToAdd = time.Minute * 5
	}

	if rangeTime == ThirdyMinutes {
		timeToAdd = time.Minute * 30
	}

	if rangeTime == OneHour {
		timeToAdd = time.Hour
		step *= 3
	}

	if rangeTime == ThreeHours {
		timeToAdd = time.Hour * 3
		step *= 10
	}

	if rangeTime == OneDay {
		timeToAdd = time.Hour * 24
		step *= 30
	}

	if rangeTime == ThreeDays {
		timeToAdd = (time.Hour * 24) * 3
		step *= 60
	}

	if rangeTime == OneWeek {
		timeToAdd = (time.Hour * 24) * 7
		step *= 120
	}

	return MetricRange{
		Start: time.Now().Add(-timeToAdd),
		End:   time.Now(),
		Step:  step,
	}

}

func (u UseCase) Query(ctx context.Context, workspaceId string, circleId string, metricId string, rangeTime string) (interface{}, error) {
	namespace, err := u.circleUseCase.FindById(ctx, workspaceId, circleId)
	if err != nil {
		return nil, err
	}

	metricName, err := id.DecodeID(metricId)
	if err != nil {
		return nil, err
	}

	circleModel, err := u.circleUseCase.FindById(ctx, workspaceId, circleId)
	if err != nil {
		return MetricModel{}, err
	}

	metricModel := MetricModel{}
	basicMetricModels, err := u.basicMatrics(circleModel)
	if err != nil {
		return MetricModel{}, err
	}

	basicMetricModel, ok := basicMetricModels[metricName]
	if ok {
		metricModel = basicMetricModel
	} else {
		metricModel, err = u.metricRepository.FindById(ctx, namespace.Name, metricId)
		if err != nil {
			return MetricModel{}, err
		}
	}

	res, err := u.metricProvider.Query(ctx, circleModel, metricModel, getMetricRange(rangeTime))
	if err != nil {
		return nil, err
	}
	return res, nil
}
