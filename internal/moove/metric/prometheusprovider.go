package metric

import (
	"context"
	"time"

	"github.com/octopipe/charlescd/internal/moove/circle"
	"github.com/prometheus/client_golang/api"
	v1 "github.com/prometheus/client_golang/api/prometheus/v1"
)

type prometheusProvider struct {
	prometheusClient api.Client
}

// Query implements MetricProvider

func NewPrometheusProvider(prometheusClient api.Client) MetricProvider {
	return prometheusProvider{prometheusClient: prometheusClient}
}

func (p prometheusProvider) Query(ctx context.Context, circleModel circle.CircleModel, metricModel MetricModel, metricRange MetricRange) (interface{}, error) {
	v1api := v1.NewAPI(p.prometheusClient)
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	r := v1.Range{
		Start: metricRange.Start,
		End:   metricRange.End,
		Step:  metricRange.Step,
	}

	result, _, err := v1api.QueryRange(ctx, metricModel.Query, r, v1.WithTimeout(5*time.Second))
	if err != nil {
		return nil, err
	}

	return result, nil
}
