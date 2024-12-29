package v1alpha1

import (
	"fmt"
	"net/url"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// KeycloakUsersConfig provides the configuration for querying Keycloak API for
// users.
type KeycloakUsersConfig struct {
	// Enabled is used to filter only users who are enabled.
	// +optional
	// +kubebuilder:default=true
	Enabled bool `json:"enabled"`
}

// ToValues() converts the config to URL Values for communicating with the
// Keycloak HTTP API.
func (k KeycloakUsersConfig) ToValues() url.Values {
	v := url.Values{}
	v.Set("enabled", fmt.Sprintf("%v", k.Enabled))

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
