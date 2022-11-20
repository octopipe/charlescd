package cmd

import (
	"github.com/octopipe/charlescd/compass/internal/server"
	metrics "k8s.io/metrics/pkg/client/clientset/versioned"
	ctrl "sigs.k8s.io/controller-runtime"
)

func main() {

	config := ctrl.GetConfigOrDie()
	clientset := metrics.NewForConfigOrDie(config)

	server := server.NewServer(logger, resourceServer)

}
