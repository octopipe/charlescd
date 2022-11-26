package test

import (
	charlescdiov1alpha1 "github.com/octopipe/charlescd/pkg/api/v1alpha1"
)

func newCircleObject(name string, moduleName string) *charlescdiov1alpha1.Circle {
	newCircle := &charlescdiov1alpha1.Circle{}
	newCircle.SetName(name)
	newCircle.SetNamespace("default")
	newCircle.Spec = charlescdiov1alpha1.CircleSpec{
		Namespace: "default",
		Author:    "Test",
		IsDefault: true,
		Routing: charlescdiov1alpha1.CircleRouting{
			Match: &charlescdiov1alpha1.MatchRouteStrategy{
				CustomMatch: &charlescdiov1alpha1.CircleMatch{
					Headers: map[string]string{
						"x-test-id": "1111",
					},
				},
			},
		},
		Modules: []charlescdiov1alpha1.CircleModule{
			{
				Name:      moduleName,
				Namespace: "default",
				Revision:  "1",
				Overrides: []charlescdiov1alpha1.Override{
					{
						Key:   "$.spec.template.spec.containers[0].image",
						Value: "mayconjrpacheco/dragonboarding:goku",
					},
				},
			},
		},
	}

	return newCircle
}

func newModuleObject(name string) *charlescdiov1alpha1.Module {
	newModule := &charlescdiov1alpha1.Module{}
	newModule.SetName(name)
	newModule.SetNamespace("default")
	newModule.Spec = charlescdiov1alpha1.ModuleSpec{
		Path:         "guestbook",
		Url:          "https://github.com/octopipe/charlescd-samples",
		TemplateType: "simple",
		Author:       "test",
	}

	return newModule
}
