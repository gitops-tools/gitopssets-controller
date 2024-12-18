package users

import (
	"context"
	"errors"
	"time"

	templatesv1 "github.com/gitops-tools/gitopssets-controller/api/v1alpha1"
	"github.com/gitops-tools/gitopssets-controller/pkg/generators"
	"github.com/go-logr/logr"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// GeneratorFactory is a function for creating per-reconciliation generators for
// the UserGenerator.
func GeneratorFactory(factory generators.HTTPClientFactory) generators.GeneratorFactory {
	return func(l logr.Logger, c client.Reader) generators.Generator {
		return NewGenerator(l, c, factory)
	}
}

// UsersGenerator is a generator that can generate user resources from different
// sources.
type UsersGenerator struct {
	ClientFactory generators.HTTPClientFactory
	Client        client.Reader
	Logger        logr.Logger
}

// NewGenerator creates and returns a new list generator.
func NewGenerator(l logr.Logger, c client.Reader, clientFactory generators.HTTPClientFactory) *UsersGenerator {
	return &UsersGenerator{
		Logger:        l,
		Client:        c,
		ClientFactory: clientFactory,
	}
}

func (g *UsersGenerator) Generate(ctx context.Context, sg *templatesv1.GitOpsSetGenerator, ks *templatesv1.GitOpsSet) ([]map[string]any, error) {
	if sg == nil {
		return nil, generators.ErrEmptyGitOpsSet
	}

	if sg.Users == nil {
		return nil, nil
	}

	g.Logger.Info("generating params from Users generator")

	switch {
	case sg.Users.Keycloak != nil:
		return g.generateKeycloakUsers(ctx, sg, ks)
	}

	return nil, errors.New("user generator is not confgiured")
}

// Interval is an implementation of the Generator interface.
//
// The UsersGenerator requires to poll regularly as there's nothing to drive
// watches.
func (g *UsersGenerator) Interval(sg *templatesv1.GitOpsSetGenerator) time.Duration {
	return sg.Users.Interval.Duration
}
