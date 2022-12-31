package metric

import (
	"context"
	"fmt"
	"time"

	"github.com/iancoleman/strcase"
	"github.com/octopipe/charlescd/internal/moove/circle"
	"github.com/octopipe/charlescd/internal/moove/core/listoptions"
	"github.com/octopipe/charlescd/internal/moove/errs"
	"github.com/octopipe/charlescd/internal/utils/id"
	charlescdiov1alpha1 "github.com/octopipe/charlescd/pkg/api/v1alpha1"
	k8sErrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type k8sRepository struct {
	clientset client.Client
}

func NewK8sRepository(clientset client.Client) MetricRepository {
	return k8sRepository{clientset: clientset}
}

func (r k8sRepository) findById(metricId string, namespace string) (charlescdiov1alpha1.Metric, error) {
	name, err := id.DecodeID(metricId)
	if err != nil {
		return charlescdiov1alpha1.Metric{}, err
	}

	metric := charlescdiov1alpha1.Metric{}
	err = r.clientset.Get(context.Background(), types.NamespacedName{Name: name, Namespace: namespace}, &metric)
	if err != nil {
		if k8sErrors.IsNotFound(err) {
			return charlescdiov1alpha1.Metric{}, errs.E(errs.NotExist, errs.Code("METRIC_NOT_FOUND"), fmt.Errorf("metric %s not found", metricId))
		}

		return charlescdiov1alpha1.Metric{}, errs.E(errs.Internal, errs.Code("METRIC_FIND_BY_ID_FAILED"), err)
	}

	return metric, nil
}

func (r k8sRepository) fillMetric(target charlescdiov1alpha1.Metric, namespace string, circleId string, metric Metric) (charlescdiov1alpha1.Metric, error) {
	labels := map[string]string{
		"managed-by": "moove",
		"circleId":   circleId,
	}

	target.SetLabels(labels)
	target.SetAnnotations(map[string]string{
		"id":   id.ToID(target.Name),
		"name": metric.Name,
	})

	target.Spec = charlescdiov1alpha1.MetricSpec{
		Author:      metric.Author,
		Description: metric.Description,
	}
	return target, nil
}

func (r k8sRepository) toMetricModel(metric charlescdiov1alpha1.Metric) MetricModel {
	annotations := metric.GetAnnotations()
	return MetricModel{
		ID: annotations["id"],
		Metric: Metric{
			Name: annotations["name"],
		},
		CreatedAt: metric.CreationTimestamp.Format(time.RFC3339),
	}
}

// Create implements MetricModelRepository
func (r k8sRepository) Create(ctx context.Context, namespace string, circleModel circle.CircleModel, metric Metric) (MetricModel, error) {
	newMetric := charlescdiov1alpha1.Metric{}
	metricName := strcase.ToKebab(metric.Name)
	newMetric.SetName(metricName)
	newMetric.SetNamespace(namespace)

	newMetric, err := r.fillMetric(newMetric, namespace, circleModel.ID, metric)
	if err != nil {
		return MetricModel{}, err
	}

	err = r.clientset.Create(context.Background(), &newMetric)
	if err != nil {
		if k8sErrors.IsAlreadyExists(err) {
			return MetricModel{}, errs.E(errs.Exist, errs.Code("METRIC_ALREADY_EXIST"), fmt.Sprintf("%s metric already exist", metric.Name))
		}

		return MetricModel{}, errs.E(errs.Internal, errs.Code("METRIC_CREATE_ERROR"), err)
	}

	return r.toMetricModel(newMetric), nil
}

// Delete implements MetricModelRepository
func (r k8sRepository) Delete(ctx context.Context, namespace string, metricId string) error {
	metric, err := r.findById(metricId, namespace)
	if err != nil {
		return err
	}

	err = r.clientset.Delete(context.Background(), &metric)
	if err != nil {
		return errs.E(errs.Internal, errs.Code("METRIC_DELETE_FAILED"), err)
	}

	return nil
}

// FindAll implements MetricModelRepository
func (r k8sRepository) FindAll(ctx context.Context, namespace string, circleModel circle.CircleModel, options listoptions.Request) (listoptions.Response, error) {
	list := &charlescdiov1alpha1.MetricList{}
	labelSelector := labels.SelectorFromSet(labels.Set{"managed-by": "moove", "circleId": circleModel.ID})
	err := r.clientset.List(context.Background(), list, &client.ListOptions{
		LabelSelector: labelSelector,
		Namespace:     namespace,
		Limit:         options.Limit,
	})
	if err != nil {
		return listoptions.Response{}, errs.E(errs.Internal, errs.Code("METRIC_LIST_ERROR"), err)
	}

	models := []MetricModel{}
	for _, i := range list.Items {
		if i.DeletionTimestamp == nil {
			models = append(models, r.toMetricModel(i))
		}
	}

	return listoptions.Response{
		Continue: list.Continue,
		Items:    models,
	}, nil
}

// FindById implements MetricModelRepository
func (r k8sRepository) FindById(ctx context.Context, namespace string, metricId string) (MetricModel, error) {
	metric, err := r.findById(metricId, namespace)
	if err != nil {
		return MetricModel{}, err
	}

	return r.toMetricModel(metric), nil
}

// Update implements MetricModelRepository
func (r k8sRepository) Update(ctx context.Context, namespace string, circleModel circle.CircleModel, metricId string, metric Metric) (MetricModel, error) {
	metricObject, err := r.findById(metricId, namespace)
	if err != nil {
		return MetricModel{}, err
	}

	metricObject, err = r.fillMetric(metricObject, namespace, circleModel.ID, metric)
	if err != nil {
		return MetricModel{}, err
	}

	err = r.clientset.Update(context.Background(), &metricObject)
	if err != nil {
		return MetricModel{}, errs.E(errs.Internal, errs.Code("METRIC_UPDATE_FAILED"), err)
	}

	return r.toMetricModel(metricObject), nil
}
