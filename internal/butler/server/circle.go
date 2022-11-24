package server

import (
	context "context"

	"github.com/octopipe/charlescd/internal/butler/sync"
	pbv1 "github.com/octopipe/charlescd/pb/v1"
	charlescdiov1alpha1 "github.com/octopipe/charlescd/pkg/api/v1alpha1"
	anypb "google.golang.org/protobuf/types/known/anypb"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type CircleServer struct {
	k8sCache   client.Client
	circleSync sync.CircleSync
}

func NewCircleServer(
	k8sCache client.Client,
	circleSync sync.CircleSync,
) pbv1.CircleServiceServer {
	return CircleServer{
		k8sCache:   k8sCache,
		circleSync: circleSync,
	}
}

// Tree implements v1.CircleServiceServer
func (s CircleServer) Sync(ctx context.Context, request *pbv1.SyncRequest) (*anypb.Any, error) {
	circle := &charlescdiov1alpha1.Circle{}
	namespacedName := types.NamespacedName{
		Name:      request.CircleName,
		Namespace: request.CircleNamespace,
	}
	err := s.k8sCache.Get(ctx, namespacedName, circle)
	if err != nil {
		return nil, err
	}

	err = s.circleSync.Sync(circle)
	if err != nil {
		return nil, err
	}

	return &anypb.Any{}, nil
}
