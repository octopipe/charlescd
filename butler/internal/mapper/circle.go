package mapper

import (
	charlescdiov1alpha1 "github.com/octopipe/charlescd/butler/api/v1alpha1"
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
	customMatch := &pbv1.CircleMatch{
		Headers: routing.Default.CustomMatch.Headers,
	}
	routingMessage := &pbv1.CircleRouting{
		Strategy: routing.Strategy,
	}

	if routing.Canary != nil {
		routingMessage.Canary = &pbv1.CanaryRouting{
			Weight: int64(routing.Canary.Weight),
		}
	}

	if routing.Default != nil {
		routingMessage.Default = &pbv1.DefaultRouting{
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
