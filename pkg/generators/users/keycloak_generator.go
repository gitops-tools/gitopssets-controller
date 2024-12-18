package users

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"

	"github.com/gitops-tools/gitopssets-controller/pkg/generators"
	"github.com/go-logr/logr"
	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// This is an interim representation of Keycloak users.
// TODO: Do we need this?
// What about additional fields?
type KeycloakUser struct {
	ID            string `json:"id"`
	Username      string `json:"username"`
	Firstname     string `json:"firstName"`
	Lastname      string `json:"lastName"`
	EmailVerified bool   `json:"emailVerified"`

	// This is a unix-timestamp - fix this!
	CreatedTimestamp int64 `json:"createdTimestamp"`
}

func (u KeycloakUser) ToMap() map[string]any {
	// Upper or Lower case for Keys?
	// Is there a package that can convert a struct to a map?!
	//    There must be something in mitchellh's code!
	return map[string]any{
		"id":            u.ID,
		"username":      u.Username,
		"firstname":     u.Firstname,
		"lastname":      u.Lastname,
		"emailVerified": u.EmailVerified,

		// TODO: CreatedTimestamp in a user-readable format.
	}
}

// This is the default Client factory it returns a client configured with the
// tls.Config.
var DefaultClientFactory = func(config *tls.Config) *http.Client {
	transport := http.DefaultTransport.(*http.Transport).Clone()
	transport.TLSClientConfig = config

	return &http.Client{
		Transport: transport,
	}
}

// TODO: Move this to the api definition

type KeycloakGeneratorUsers struct {
	// TODO: Mark this as default to true.
	Enabled bool `json:"enabled"`
}

func (k KeycloakGeneratorUsers) ToValues() url.Values {
	v := url.Values{}
	v.Set("enabled", fmt.Sprintf("%v", k.Enabled))

	return v
}

// FetchConfig configures how to load the results from Keycloak.
type FetchConfig struct {
	// TODO: Mark this as default to true.
	AllPages bool `json:"allPages"`
	PageSize int  `json:"pageSize"`
}

// KeycloakGenerator generates from a Keycloak service.
type KeycloakGeneratorConfig struct {
	// TODO: Easier to make one "realm-terminated" url?
	APIEndpoint string       `json:"apiEndpoint"`
	Realm       string       `json:"realm"`
	Fetch       *FetchConfig `json:"fetchConfig"`

	SecretRef corev1.LocalObjectReference `json:"credentials"`

	// Users is a set of rules for fetching users.
	Users *KeycloakGeneratorUsers `json:"users,omitempty"`
}

// Keycloak generator generates from a Keycloak API.
type KeycloakGenerator struct {
	ClientFactory generators.HTTPClientFactory
	Client        client.Reader
	logr.Logger
}

// // GeneratorFactory is a function for creating per-reconciliation generators for
// // the KeycloakGenerator.
// func GeneratorFactory(factory HTTPClientFactory) generators.GeneratorFactory {
// 	return func(l logr.Logger, c client.Reader) generators.Generator {
// 		return NewGenerator(l, c, factory)
// 	}
// }

// NewGenerator creates and returns a new Keycloak generator.
func NewKeycloakGenerator(l logr.Logger, c client.Reader, clientFactory HTTPClientFactory) *KeycloakGenerator {
	return &KeycloakGenerator{
		Client:        c,
		Logger:        l,
		ClientFactory: clientFactory,
	}
}

func (k *KeycloakGenerator) Generate(ctx context.Context, config KeycloakGeneratorConfig) ([]map[string]any, error) {
	// TODO: Improve the URL generation
	req, err := http.NewRequest(http.MethodGet, config.APIEndpoint+"/users/?"+config.Users.ToValues().Encode(), nil)
	if err != nil {
		// TODO: Improve error
		return nil, err
	}
	req.Header.Add("Authorization", "Bearer "+os.Getenv("BEARER_TOKEN"))
	client := k.ClientFactory(nil)

	resp, err := client.Do(req)
	if err != nil {
		// TODO: Improve error
		return nil, err
	}

	if resp.StatusCode > 300 {
		return nil, fmt.Errorf("invalid response from %s: %v", config.APIEndpoint, resp.StatusCode)
	}

	// Should we just return the result?
	decoder := json.NewDecoder(resp.Body)
	// TODO: Wrap this in a function that can report the error.
	defer resp.Body.Close()

	var users []KeycloakUser
	if err := decoder.Decode(&users); err != nil {
		// TODO: Improve error
		return nil, err
	}

	var result []map[string]any
	for _, user := range users {
		result = append(result, user.ToMap())
	}

	return result, nil
}
