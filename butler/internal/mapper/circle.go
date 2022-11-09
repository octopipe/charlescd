package mapper

import (
	charlescdiov1alpha1 "github.com/octopipe/charlescd/butler/api/v1alpha1"
	"github.com/octopipe/charlescd/butler/internal/errs"
	pbv1 "github.com/octopipe/charlescd/butler/pb/v1"
)

func CircleToProtoMessage(circle charlescdiov1alpha1.Circle) *pbv1.Circle {
	routing := CircleRoutingToProtoMessage(circle.Spec.Routing)
	modules := CircleModulesToProtoMessage(circle.Spec.Modules)
	environments := CircleEnvironmentsToProtoMessage(circle.Spec.Environments)
	status := CircleStatusToProtoMessage(circle.Status)

	return &pbv1.Circle{
		Name:         circle.GetName(),
		Namespace:    circle.GetNamespace(),
		IsDefault:    circle.Spec.IsDefault,
		Routing:      routing,
		Modules:      modules,
		Environments: environments,
		Status:       status,
	}
}

func CircleRoutingToProtoMessage(routing charlescdiov1alpha1.CircleRouting) *pbv1.CircleRouting {
	strategy := 0
	switch routing.Strategy {
	case "MATCH":
		strategy = 1
		break
	case "CANARY":
		strategy = 2
		break
	}

	customMatch := &pbv1.CircleMatch{
		Headers: routing.Match.CustomMatch.Headers,
	}
	routingMessage := &pbv1.CircleRouting{
		Strategy: pbv1.RoutingType(strategy),
	}

	if routing.Canary != nil {
		routingMessage.Canary = &pbv1.CanaryRouting{
			Weight: int64(routing.Canary.Weight),
		}
	}

	if routing.Match != nil {
		routingMessage.Match = &pbv1.MatchRouting{
			CustomMatch: customMatch,
		}
	}
	return routingMessage
}

func CircleModulesToProtoMessage(modules []charlescdiov1alpha1.CircleModule) []*pbv1.CircleModule {
	modulesMessage := []*pbv1.CircleModule{}
	for _, module := range modules {
		overrides := []*pbv1.CircleModuleOverride{}

		for _, o := range module.Overrides {
			overrides = append(overrides, &pbv1.CircleModuleOverride{
				Key:   o.Key,
				Value: o.Value,
			})
		}

		m := &pbv1.CircleModule{
			Name:      module.Name,
			Revision:  module.Revision,
			Overrides: overrides,
		}

		modulesMessage = append(modulesMessage, m)
	}

	return modulesMessage
}

func CircleEnvironmentsToProtoMessage(environments []charlescdiov1alpha1.CircleEnvironments) []*pbv1.CircleEnvironment {
	environmentsMessage := []*pbv1.CircleEnvironment{}
	for _, environment := range environments {
		e := &pbv1.CircleEnvironment{
			Key:   environment.Key,
			Value: environment.Value,
		}

		environmentsMessage = append(environmentsMessage, e)
	}

	return environmentsMessage
}

func CircleStatusToProtoMessage(status charlescdiov1alpha1.CircleStatus) *pbv1.CircleStatus {
	moduleStatusMessage := map[string]*pbv1.CircleStatusModule{}
	for moduleName, module := range status.Modules {

		resourcesMessage := []*pbv1.CircleStatusResource{}
		for _, resource := range module.Resources {
			r := &pbv1.CircleStatusResource{
				Name:      resource.Name,
				Namespace: resource.Namespace,
				Group:     resource.Group,
				Kind:      resource.Kind,
				Health:    resource.Health,
				Error:     resource.Error,
			}
			resourcesMessage = append(resourcesMessage, r)
		}

		m := &pbv1.CircleStatusModule{
			Status:    module.Status,
			Error:     module.Error,
			Resources: resourcesMessage,
		}
		moduleStatusMessage[moduleName] = m
	}

	return &pbv1.CircleStatus{
		Modules: moduleStatusMessage,
	}
}

func ProtoMessageToCircle(message *pbv1.CreateCircleRequest) (*charlescdiov1alpha1.Circle, error) {
	if err := message.ValidateAll(); err != nil {
		return nil, errs.E(errs.Invalid, errs.Code("proto_to_circle"), err)
	}

	circle := &charlescdiov1alpha1.Circle{}
	circle.SetName(message.Name)
	circle.SetNamespace(message.Namespace)

	routing := charlescdiov1alpha1.CircleRouting{}
	if message.Routing.Canary != nil {
		routing.Canary = &charlescdiov1alpha1.CanaryDeployStrategy{
			Weight: int(message.Routing.Canary.Weight),
		}
	}

	if message.Routing.Match != nil {
		if message.Routing.Match.CustomMatch != nil {
			routing.Match = &charlescdiov1alpha1.MatchRouteStrategy{
				CustomMatch: &charlescdiov1alpha1.CircleMatch{
					Headers: message.Routing.Match.CustomMatch.Headers,
				},
			}
		}
	}

	environments := []charlescdiov1alpha1.CircleEnvironments{}
	if message.Environments != nil {
		for _, e := range message.Environments {
			environments = append(environments, charlescdiov1alpha1.CircleEnvironments{
				Key:   e.Key,
				Value: e.Value,
			})
		}
	}

	modules := []charlescdiov1alpha1.CircleModule{}
	if message.Modules != nil {
		for _, m := range message.Modules {
			overrides := []charlescdiov1alpha1.Override{}
			if m.Overrides != nil {
				for _, o := range m.Overrides {
					overrides = append(overrides, charlescdiov1alpha1.Override{
						Key:   o.Key,
						Value: o.Value,
					})
				}
			}

			modules = append(modules, charlescdiov1alpha1.CircleModule{
				Name:      m.Name,
				Namespace: m.Namespace,
				Overrides: overrides,
				Revision:  m.Revision,
			})
		}
	}

	circle.Spec = charlescdiov1alpha1.CircleSpec{
		IsDefault:    message.IsDefault,
		Routing:      routing,
		Environments: environments,
		Modules:      modules,
	}

	return circle, nil
}
