package v1alpha1

import (
	"fmt"
	"net/url"
	"strings"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// KeycloakUsersConfig provides the configuration for querying Keycloak API for
// users.
type KeycloakUsersConfig struct {
	// Enabled is used to filter only users who are enabled.
	// +optional
	Enabled *bool `json:"enabled"`

	// EmailVerified is used to filter users with verified emails.
	// +optional
	EmailVerified *bool `json:"emailVerified"`

	// Email filters users containing the string.
	// The Exact option modifies the search.
	// +optional
	Email string `json:"email"`

	// Firstname filters users containing the string.
	// The Exact option modifies the search.
	// +optional
	Firstname string `json:"firstName"`

	// Lastname filters users containing the string.
	// The Exact option modifies the search.
	// +optional
	Lastname string `json:"lastName"`

	// Username filters users containing the string.
	// The Exact option modifies the search.
	// +optional
	Username string `json:"userName"`

	// Search filters users with the username, first or last name, or email
	// containing the string.
	// +optional
	Search string `json:"search"`

	// Query filters users with attributes.
	// +optional
	Query map[string]string `json:"query,omitempty"`

	// Exact controls whether or not the Email, LastName, FirstName and Username
	// searches must match exactly.
	Exact *bool `json:"exact"`

	// Limit the number of results in the query.
	// +optional
	Limit *int `json:"limit"`

	// TODO: Change to a number of pages to query for with 0 being unlimited.
	// List all users page-by-page
	// +kubebuilder:default=true
	// +optional
	AllPages bool `json:"allPages"`
}

// ToValues() converts the config to URL Values for communicating with the
// Keycloak HTTP API.
func (k KeycloakUsersConfig) ToValues() url.Values {
	v := url.Values{}
	// TODO: Tests!

	if k.Enabled != nil {
		v.Set("enabled", fmt.Sprintf("%v", *k.Enabled))
	}

	if k.EmailVerified != nil {
		v.Set("emailVerified", fmt.Sprintf("%v", *k.EmailVerified))
	}

	if k.Exact != nil {
		v.Set("exact", fmt.Sprintf("%v", *k.Exact))
	}

	if k.Email != "" {
		v.Set("email", k.Email)
	}

	if k.Firstname != "" {
		v.Set("firstName", k.Firstname)
	}

	if k.Lastname != "" {
		v.Set("lastName", k.Lastname)
	}

	if k.Username != "" {
		v.Set("username", k.Username)
	}

	if k.Search != "" {
		v.Set("search", k.Search)
	}

	if k.Query != nil {
		var qa []string
		for k, v := range k.Query {
			qa = append(qa, fmt.Sprintf("%s:%s", k, v))
		}
		v.Set("q", strings.Join(qa, " "))
	}

	// The API defaults to 100 users.
	if k.Limit != nil {
		if *k.Limit > 0 && *k.Limit != 100 {
			v.Set("max", fmt.Sprintf("%v", *k.Limit))
		}
	}

	return v
}

// KeycloakUsersGeneration configures the Keycloak method for querying for
// users.
type KeycloakUsersGeneration struct {
	// This is the API endpoint to use.
	// +kubebuilder:validation:Pattern="^https://"
	Endpoint string `json:"endpoint"`
	// +required
	SecretRef *LocalObjectReference `json:"secretRef,omitempty"`

	// Control the users that are queried from Keycloak.
	// +required
	QueryConfig *KeycloakUsersConfig `json:"queryConfig"`
}

// UsersGenerator defines a generator that queries for Users from an upstream
// API.
type UsersGenerator struct {
	// The interval at which to poll the API endpoint.
	// +required
	Interval metav1.Duration `json:"interval"`

	// +optional
	Keycloak *KeycloakUsersGeneration `json:"keycloak,omitempty"`
}
