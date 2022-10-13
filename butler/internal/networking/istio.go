package networking

import (
	"context"

	charlescdiov1alpha1 "github.com/octopipe/charlescd/butler/api/v1alpha1"
	"github.com/octopipe/charlescd/butler/internal/utils"
	networkingv1alpha3 "istio.io/api/networking/v1alpha3"
	"istio.io/client-go/pkg/apis/networking/v1alpha3"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/util/retry"
)

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
			Http: []*networkingv1alpha3.HTTPRoute{
				{
					Name: "dsds",
					Match: []*networkingv1alpha3.HTTPMatchRequest{
						{
							Headers: map[string]*networkingv1alpha3.StringMatch{
								"x-circle-name": {
									MatchType: &networkingv1alpha3.StringMatch_Exact{
										Exact: "",
									},
								},
							},
						},
					},
				},
			},
		},
	}

	newVirtualService.SetName(module.ModuleRef)

	return newVirtualService
}

func (n networkingLayer) SyncIstio(circle charlescdiov1alpha1.Circle) error {

	namespace := "default"
	if circle.Spec.Namespace != "" {
		namespace = circle.Spec.Namespace
	}

	for _, module := range circle.Spec.Modules {
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
	}

	return nil
}
