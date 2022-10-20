package networking

import (
	"context"
	"fmt"

	"github.com/argoproj/gitops-engine/pkg/health"
	charlescdiov1alpha1 "github.com/octopipe/charlescd/butler/api/v1alpha1"
	"github.com/octopipe/charlescd/butler/internal/utils"
	networkingv1alpha3 "istio.io/api/networking/v1alpha3"
	"istio.io/client-go/pkg/apis/networking/v1alpha3"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/util/retry"
)

func getMatchForDefaultRouting(circle charlescdiov1alpha1.Circle) []*networkingv1alpha3.HTTPMatchRequest {
	match := []*networkingv1alpha3.HTTPMatchRequest{}
	defaultRouting := circle.Spec.Routing.Default

	if len(defaultRouting.CustomMatch) > 0 {
		for _, customMatch := range defaultRouting.CustomMatch {

			headers := map[string]*networkingv1alpha3.StringMatch{}
			for key, value := range customMatch.Headers {
				headers[key] = &networkingv1alpha3.StringMatch{
					MatchType: &networkingv1alpha3.StringMatch_Regex{Regex: value},
				}
			}

			newMatch := &networkingv1alpha3.HTTPMatchRequest{
				Name:    circle.GetName(),
				Headers: headers,
			}

			match = append(match, newMatch)
		}

	}

	return match
}

func getMatch(circle charlescdiov1alpha1.Circle) []*networkingv1alpha3.HTTPMatchRequest {
	match := []*networkingv1alpha3.HTTPMatchRequest{
		{
			Uri: &networkingv1alpha3.StringMatch{
				MatchType: &networkingv1alpha3.StringMatch_Exact{
					Exact: "/",
				},
			},
		},
	}

	// if circle.Spec.Routing.Strategy == DefaultRouting {
	// 	return getMatchForDefaultRouting(circle)
	// }

	return match
}

func getModuleService(module charlescdiov1alpha1.CircleModuleStatus) *charlescdiov1alpha1.CircleModuleResource {
	for _, res := range module.Resources {
		if res.Kind == "Service" {
			return &res
		}
	}

	return nil
}

func getRoutes(circle charlescdiov1alpha1.Circle) []*networkingv1alpha3.HTTPRoute {
	routes := []*networkingv1alpha3.HTTPRoute{}

	for name, module := range circle.Status.Modules {
		if module.Status != string(health.HealthStatusHealthy) {
			continue
		}

		service := getModuleService(module)
		if service == nil {
			continue
		}

		newRoute := &networkingv1alpha3.HTTPRoute{
			Name:  name,
			Match: getMatch(circle),
			Route: []*networkingv1alpha3.HTTPRouteDestination{
				{
					Destination: &networkingv1alpha3.Destination{
						Host:   fmt.Sprintf("%s.%s.svc.cluster.local", service.Name, service.Namespace),
						Subset: fmt.Sprintf("%s-%s", circle.GetName(), name),
					},
				},
			},
		}

		routes = append(routes, newRoute)
	}

	return routes
}

func newVirtualService(module charlescdiov1alpha1.CircleModule, circle charlescdiov1alpha1.Circle) *v1alpha3.VirtualService {
	newVirtualService := &v1alpha3.VirtualService{
		ObjectMeta: metav1.ObjectMeta{
			Name:      module.ModuleRef,
			Namespace: circle.Spec.Namespace,
			Annotations: map[string]string{
				utils.AnnotationManagedBy: utils.ManagedBy,
			},
		},
		Spec: networkingv1alpha3.VirtualService{
			Hosts: []string{
				"*",
			},
			Gateways: []string{
				"myapp/charlescd-guestbook-simple-gateway",
			},
			Http: getRoutes(circle),
		},
	}

	newVirtualService.SetName(module.ModuleRef)

	return newVirtualService
}

func newDestinationRule(module charlescdiov1alpha1.CircleModule, circle charlescdiov1alpha1.Circle) *v1alpha3.DestinationRule {
	service := getModuleService(circle.Status.Modules[module.ModuleRef])
	destinationRule := &v1alpha3.DestinationRule{
		ObjectMeta: metav1.ObjectMeta{
			Name:      module.ModuleRef,
			Namespace: circle.Spec.Namespace,
			Annotations: map[string]string{
				utils.AnnotationManagedBy: utils.ManagedBy,
			},
		},
		Spec: networkingv1alpha3.DestinationRule{
			Host: fmt.Sprintf("%s.%s.svc.cluster.local", service.Name, service.Namespace),
			Subsets: []*networkingv1alpha3.Subset{
				{
					Name: fmt.Sprintf("%s-%s", circle.GetName(), module.ModuleRef),
					Labels: map[string]string{
						"version": fmt.Sprintf("%s-%s", circle.GetName(), module.ModuleRef),
					},
				},
			},
		},
	}

	destinationRule.SetName(module.ModuleRef)

	return destinationRule
}

func (n networkingLayer) createOrUpdateVirtualService(circle charlescdiov1alpha1.Circle, module charlescdiov1alpha1.CircleModule, namespace string) error {
	currVirtualService, err := n.istioClient.NetworkingV1alpha3().VirtualServices(namespace).Get(context.TODO(), module.ModuleRef, metav1.GetOptions{})
	if errors.IsNotFound(err) {
		newVirtualService := newVirtualService(module, circle)

		_, err = n.istioClient.NetworkingV1alpha3().VirtualServices(namespace).Create(context.TODO(), newVirtualService, metav1.CreateOptions{})
	}

	if !errors.IsNotFound(err) {
		return err
	}

	retryErr := retry.RetryOnConflict(retry.DefaultRetry, func() error {
		_, err = n.istioClient.NetworkingV1alpha3().VirtualServices(namespace).Update(context.TODO(), currVirtualService, metav1.UpdateOptions{})
		return err
	})
	if retryErr != nil {
		return retryErr
	}

	return nil
}

func (n networkingLayer) createOrUpdateDestinationRules(circle charlescdiov1alpha1.Circle, module charlescdiov1alpha1.CircleModule, namespace string) error {
	currDestinationRules, err := n.istioClient.NetworkingV1alpha3().DestinationRules(namespace).Get(context.TODO(), module.ModuleRef, metav1.GetOptions{})
	if errors.IsNotFound(err) {
		newDestinationRule := newDestinationRule(module, circle)

		_, err = n.istioClient.NetworkingV1alpha3().DestinationRules(namespace).Create(context.TODO(), newDestinationRule, metav1.CreateOptions{})
	}

	if !errors.IsNotFound(err) {
		return err
	}

	retryErr := retry.RetryOnConflict(retry.DefaultRetry, func() error {
		_, err = n.istioClient.NetworkingV1alpha3().DestinationRules(namespace).Update(context.TODO(), currDestinationRules, metav1.UpdateOptions{})
		return err
	})
	if retryErr != nil {
		return retryErr
	}

	return nil
}

func (n networkingLayer) SyncIstio(circle charlescdiov1alpha1.Circle) error {
	namespace := "default"
	if circle.Spec.Namespace != "" {
		namespace = circle.Spec.Namespace
	}

	for _, module := range circle.Spec.Modules {
		currentModule := circle.Status.Modules[module.ModuleRef]
		if currentModule.Status != string(health.HealthStatusHealthy) {
			continue
		}

		service := getModuleService(currentModule)
		if service == nil {
			continue
		}

		err := n.createOrUpdateVirtualService(circle, module, namespace)
		if err != nil {
			return err
		}

		err = n.createOrUpdateDestinationRules(circle, module, namespace)
		if err != nil {
			return err
		}
	}

	return nil
}
