package networking

import (
	"context"
	"fmt"

	"github.com/octopipe/charlescd/internal/butler/utils"
	charlescdiov1alpha1 "github.com/octopipe/charlescd/pkg/api/v1alpha1"
	networkingv1alpha3 "istio.io/api/networking/v1alpha3"
	"istio.io/client-go/pkg/apis/networking/v1alpha3"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/util/retry"
)

func getMatchForDefaultRouting(circle charlescdiov1alpha1.Circle) []*networkingv1alpha3.HTTPMatchRequest {
	match := []*networkingv1alpha3.HTTPMatchRequest{}
	matchRouting := circle.Spec.Routing.Match

	if matchRouting != nil {
		headers := map[string]*networkingv1alpha3.StringMatch{}
		for key, value := range matchRouting.Headers {
			headers[key] = &networkingv1alpha3.StringMatch{
				MatchType: &networkingv1alpha3.StringMatch_Regex{Regex: value},
			}
		}

		newMatch := &networkingv1alpha3.HTTPMatchRequest{
			Headers: headers,
		}

		match = append(match, newMatch)

	}

	return match
}

func getMatch(circle charlescdiov1alpha1.Circle) []*networkingv1alpha3.HTTPMatchRequest {
	match := []*networkingv1alpha3.HTTPMatchRequest{}

	if circle.Spec.Routing.Strategy == charlescdiov1alpha1.MatchRoutingStrategy {
		return getMatchForDefaultRouting(circle)
	}

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

func getRoute(circle charlescdiov1alpha1.Circle, module charlescdiov1alpha1.CircleModule) *networkingv1alpha3.HTTPRoute {
	currentModule := circle.Status.Modules[module.Name]
	service := getModuleService(currentModule)

	if circle.Spec.Routing.Strategy != charlescdiov1alpha1.CanaryRoutingStrategy {
		return &networkingv1alpha3.HTTPRoute{
			Name:  circle.GetName(),
			Match: getMatch(circle),
			Route: []*networkingv1alpha3.HTTPRouteDestination{
				{
					Destination: &networkingv1alpha3.Destination{
						Host:   fmt.Sprintf("%s.%s.svc.cluster.local", service.Name, service.Namespace),
						Subset: circle.GetName(),
					},
				},
			},
		}
	}

	return &networkingv1alpha3.HTTPRoute{
		Name: circle.GetName(),
		Route: []*networkingv1alpha3.HTTPRouteDestination{
			{
				Destination: &networkingv1alpha3.Destination{
					Host:   fmt.Sprintf("%s.%s.svc.cluster.local", service.Name, service.Namespace),
					Subset: circle.GetName(),
				},
				Weight: int32(circle.Spec.Routing.Canary.Weight),
			},
		},
	}
}

func newVirtualService(module charlescdiov1alpha1.CircleModule, circle charlescdiov1alpha1.Circle) *v1alpha3.VirtualService {
	newVirtualService := &v1alpha3.VirtualService{
		ObjectMeta: metav1.ObjectMeta{
			Name:      module.Name,
			Namespace: circle.Spec.Namespace,
			Labels: map[string]string{
				utils.AnnotationManagedBy: utils.ManagedBy,
			},
		},
		Spec: networkingv1alpha3.VirtualService{
			Hosts: []string{
				"*",
			},
			Gateways: []string{
				"guestbook-gateway",
			},
			Http: []*networkingv1alpha3.HTTPRoute{getRoute(circle, module)},
		},
	}

	newVirtualService.SetName(module.Name)

	return newVirtualService
}

func mergeVirtualServices(module charlescdiov1alpha1.CircleModule, circle charlescdiov1alpha1.Circle, virtualService *v1alpha3.VirtualService) *v1alpha3.VirtualService {
	currentRoutes := []*networkingv1alpha3.HTTPRoute{getRoute(circle, module)}
	for _, r := range virtualService.Spec.Http {
		if r.Name != circle.GetName() {
			currentRoutes = append(currentRoutes, r)
		}
	}

	virtualService.Spec.Http = currentRoutes

	return virtualService
}

func (n networkingLayer) createOrUpdateVirtualService(circle charlescdiov1alpha1.Circle, module charlescdiov1alpha1.CircleModule, namespace string) (*v1alpha3.VirtualService, error) {
	currVirtualService, err := n.istioClient.NetworkingV1alpha3().VirtualServices(namespace).Get(context.TODO(), module.Name, metav1.GetOptions{})
	if errors.IsNotFound(err) {
		newVirtualService := newVirtualService(module, circle)
		newVirtualService, err = n.istioClient.NetworkingV1alpha3().VirtualServices(namespace).Create(context.TODO(), newVirtualService, metav1.CreateOptions{})
		return newVirtualService, err
	}

	if err != nil {
		return nil, err
	}

	currVirtualService = mergeVirtualServices(module, circle, currVirtualService)

	retryErr := retry.RetryOnConflict(retry.DefaultRetry, func() error {
		currVirtualService, err = n.istioClient.NetworkingV1alpha3().VirtualServices(namespace).Update(context.TODO(), currVirtualService, metav1.UpdateOptions{})
		return err
	})
	if retryErr != nil {
		return nil, retryErr
	}

	return currVirtualService, nil
}

func newDestinationRule(module charlescdiov1alpha1.CircleModule, circle charlescdiov1alpha1.Circle) *v1alpha3.DestinationRule {
	service := getModuleService(circle.Status.Modules[module.Name])
	destinationRule := &v1alpha3.DestinationRule{
		ObjectMeta: metav1.ObjectMeta{
			Name:      module.Name,
			Namespace: circle.Spec.Namespace,
			Labels: map[string]string{
				utils.AnnotationManagedBy: utils.ManagedBy,
			},
		},
		Spec: networkingv1alpha3.DestinationRule{
			Host: fmt.Sprintf("%s.%s.svc.cluster.local", service.Name, service.Namespace),
			Subsets: []*networkingv1alpha3.Subset{
				{
					Name: circle.GetName(),
					Labels: map[string]string{
						utils.AnnotationCircleMark: utils.GetCircleMark(types.NamespacedName{Namespace: circle.Namespace, Name: circle.Name}),
					},
				},
			},
		},
	}

	destinationRule.SetName(module.Name)

	return destinationRule
}

func mergeDestionRules(module charlescdiov1alpha1.CircleModule, circle charlescdiov1alpha1.Circle, destinationRule *v1alpha3.DestinationRule) *v1alpha3.DestinationRule {
	newSubset := &networkingv1alpha3.Subset{
		Name: circle.GetName(),
		Labels: map[string]string{
			utils.AnnotationCircleMark: utils.GetCircleMark(types.NamespacedName{Namespace: circle.Namespace, Name: circle.Name}),
		},
	}

	subsets := []*networkingv1alpha3.Subset{newSubset}
	for _, s := range destinationRule.Spec.Subsets {
		if s.Name != circle.GetName() {
			subsets = append(subsets, s)
		}
	}

	destinationRule.Spec.Subsets = subsets

	return destinationRule
}

func (n networkingLayer) createOrUpdateDestinationRules(circle charlescdiov1alpha1.Circle, module charlescdiov1alpha1.CircleModule, namespace string) (*v1alpha3.DestinationRule, error) {
	currDestinationRules, err := n.istioClient.NetworkingV1alpha3().DestinationRules(namespace).Get(context.TODO(), module.Name, metav1.GetOptions{})
	if errors.IsNotFound(err) {
		newDestinationRule := newDestinationRule(module, circle)

		newDestinationRule, err = n.istioClient.NetworkingV1alpha3().DestinationRules(namespace).Create(context.TODO(), newDestinationRule, metav1.CreateOptions{})
		return newDestinationRule, err
	}

	if err != nil && !errors.IsNotFound(err) {
		return nil, err
	}

	currDestinationRules = mergeDestionRules(module, circle, currDestinationRules)
	retryErr := retry.RetryOnConflict(retry.DefaultRetry, func() error {
		currDestinationRules, err = n.istioClient.NetworkingV1alpha3().DestinationRules(namespace).Update(context.TODO(), currDestinationRules, metav1.UpdateOptions{})
		return err
	})
	if retryErr != nil {
		return nil, retryErr
	}

	return currDestinationRules, nil
}

func (n networkingLayer) SyncIstio(circle charlescdiov1alpha1.Circle, circleModule charlescdiov1alpha1.CircleModule) ([]charlescdiov1alpha1.CircleModuleResource, error) {
	namespace := circle.Spec.Namespace
	res := []charlescdiov1alpha1.CircleModuleResource{}

	currentModule := circle.Status.Modules[circleModule.Name]
	// if currentModule.Status != string(health.HealthStatusHealthy) {
	// 	return nil, nil
	// }
	service := getModuleService(currentModule)
	if service == nil {
		return nil, nil
	}

	virtualService, err := n.createOrUpdateVirtualService(circle, circleModule, namespace)
	if err != nil {
		return nil, err
	}

	destinationRule, err := n.createOrUpdateDestinationRules(circle, circleModule, namespace)
	if err != nil {
		return nil, err
	}

	res = append(res, charlescdiov1alpha1.CircleModuleResource{
		Group:     v1alpha3.SchemeGroupVersion.Group,
		Kind:      "VirtualService",
		Name:      virtualService.Name,
		Namespace: virtualService.Namespace,
	})

	res = append(res, charlescdiov1alpha1.CircleModuleResource{
		Group:     v1alpha3.SchemeGroupVersion.Group,
		Kind:      "DestinationRule",
		Name:      destinationRule.Name,
		Namespace: destinationRule.Namespace,
	})

	return res, nil
}
