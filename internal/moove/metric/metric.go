package metric

import (
	"context"
	"time"

	"github.com/octopipe/charlescd/internal/moove/circle"
	"github.com/octopipe/charlescd/internal/moove/core/listoptions"
	charlescdiov1alpha1 "github.com/octopipe/charlescd/pkg/api/v1alpha1"
)

const (
	BasicMetricCategory  = "BASIC"
	CustomMetricCategory = "CUSTOM"
)

const (
	PrometheusMetricProvider = "PROMETHEUS"
)

type Metric struct {
	charlescdiov1alpha1.MetricSpec
	Name     string `json:"name,omitempty"`
	Category string `json:"category,omitempty"`
	Provider string `json:"provider,omitempty"`
}

type MetricModel struct {
	ID        string `json:"id"`
	CreatedAt string `json:"createdAt"`
	Metric
}

type MetricRange struct {
	Start time.Time
	End   time.Time
	Step  time.Duration
}

type MetricProvider interface {
	Query(ctx context.Context, circleModel circle.CircleModel, metricModel MetricModel, metricRange MetricRange) (interface{}, error)
}

type MetricRepository interface {
	FindAll(ctx context.Context, namespace string, circleModel circle.CircleModel, listoptions listoptions.Request) (listoptions.Response, error)
	FindById(ctx context.Context, namespace string, metricId string) (MetricModel, error)
	Create(ctx context.Context, namespace string, circleModel circle.CircleModel, circle Metric) (MetricModel, error)
	Update(ctx context.Context, namespace string, circleModel circle.CircleModel, metricId string, circle Metric) (MetricModel, error)
	Delete(ctx context.Context, namespace string, metricId string) error
}

type MetricUseCase interface {
	FindAll(ctx context.Context, workspaceId string, circleId string, listoptions listoptions.Request) (listoptions.Response, error)
	FindById(ctx context.Context, workspaceId string, metricId string) (MetricModel, error)
	Create(ctx context.Context, workspaceId string, circleId string, circle Metric) (MetricModel, error)
	Update(ctx context.Context, workspaceId string, circleId string, metricId string, circle Metric) (MetricModel, error)
	Delete(ctx context.Context, workspaceId string, metricId string) error
	Query(ctx context.Context, workspaceId string, circleId string, metricId string, metricRange MetricRange) (interface{}, error)
}
