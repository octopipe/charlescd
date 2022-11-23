package networking

import (
	"errors"

	charlescdiov1alpha1 "github.com/octopipe/charlescd/pkg/api/v1alpha1"
	versionedclient "istio.io/client-go/pkg/clientset/versioned"
)

const (
	IstioLayer = "istio"
	GateLayer  = "gate"

	DefaultRouting = "default"
	CanaryRouting  = "canary"
)

type NetworkingLayer interface {
	Sync(circle charlescdiov1alpha1.Circle) error
}

type networkingLayer struct {
	networkingType string
	istioClient    *versionedclient.Clientset
}

func NewNetworkingLayer(networkingType string, istioClient *versionedclient.Clientset) NetworkingLayer {
	return networkingLayer{
		networkingType: networkingType,
		istioClient:    istioClient,
	}
}

func (n networkingLayer) Sync(circle charlescdiov1alpha1.Circle) error {
	switch n.networkingType {
	case IstioLayer:
		return n.SyncIstio(circle)
	default:
		return errors.New("cannot support this networking layer")
	}
}
