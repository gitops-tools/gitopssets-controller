package v1alpha1

import (
	"bytes"
	"context"
	"errors"
	"fmt"

	"google.golang.org/grpc/codes"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"knative.dev/pkg/apis"
	duckv1 "knative.dev/pkg/apis/duck/v1"
)

// ErrNilURL is returned when calculating the URL for an external generator if
// the URL to access it is nil.
var ErrNilURL = errors.New("generator URL is nil")

// +genclient
// +genreconciler:krshapedlogic=false
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +k8s:openapi-gen=true
// Generator describes a pluggable generator including configuration
// such as the fields it accepts and its deployment address.
type Generator struct {
	metav1.TypeMeta `json:",inline"`
	// +optional
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec GeneratorSpec `json:"spec"`
	// +optional
	Status GeneratorStatus `json:"status"`
}

// GeneratorSpec describes the Spec for an Generator
type GeneratorSpec struct {
	ClientConfig ClientConfig `json:"clientConfig"`
}

// GeneratorStatus holds the status of the Generator
// +k8s:deepcopy-gen=true
type GeneratorStatus struct {
	duckv1.Status `json:",inline"`

	// Generator is Addressable and exposes the URL where the Generator is running
	duckv1.AddressStatus `json:",inline"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// GeneratorList contains a list of Generator
// We don't use this but it's required for certain codegen features.
type GeneratorList struct {
	metav1.TypeMeta `json:",inline"`
	// +optional
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Generator `json:"items"`
}

// ResolveAddress returns the URL where the generator is running using its clientConfig
func (it *Generator) ResolveAddress() (*apis.URL, error) {
	if url := it.Spec.ClientConfig.URL; url != nil {
		return url, nil
	}
	svc := it.Spec.ClientConfig.Service
	if svc == nil {
		return nil, ErrNilURL
	}
	var (
		port *int32
	)

	if svc.Port != nil {
		port = svc.Port
	}

	if bytes.Equal(it.Spec.ClientConfig.CaBundle, []byte{}) {
		if port == nil {
			port = &defaultHTTPPort
		}
		return formURL("http", svc, port), nil
	}

	if port == nil {
		port = &defaultHTTPSPort
	}

	return formURL("https", svc, port), nil
}

// +k8s:deepcopy-gen=false
type GeneratorInterface interface {
	// Process executes the given GeneratorRequest.
	// Simply getting a non-nil GeneratorResponse back is not sufficient
	// to determine if the generator processing was successful.
	Process(ctx context.Context, r *GeneratorRequest) *GeneratorResponse
}

// +k8s:deepcopy-gen=false
type GeneratorRequest struct {
	// GeneratorParams are the user specified params for generator in the
	// GitOpsSet.
	GeneratorParams map[string]interface{} `json:"generator_params,omitempty"`
}

// +k8s:deepcopy-gen=false
type GeneratorResponse struct {
	// Generated is the result of the generation, a slice of maps containing the
	// result data.
	Generated []map[string]interface{} `json:"generated,omitempty"`
	// Status is an Error status containing details on any generator processing
	// errors.
	Status Status `json:"status"`
}

type Status struct {
	// The status code, which should be an enum value of [google.rpc.Code][google.rpc.Code].
	Code codes.Code `json:"code,omitempty"`
	// A developer-facing error message, which should be in English.
	Message string `json:"message,omitempty"`
}

func (s Status) Err() StatusError {
	return StatusError{s: s}
}

type StatusError struct {
	s Status
}

func (s StatusError) Error() string {
	return fmt.Sprintf("rpc error: code = %s desc = %s", s.s.Code, s.s.Message)
}

// ClientConfig describes how a client can communicate with the Interceptor
type ClientConfig struct {
	// CaBundle is a PEM encoded CA bundle which will be used to validate the clusterinterceptor server certificate
	CaBundle []byte `json:"caBundle,omitempty"`
	// URL is a fully formed URL pointing to the interceptor
	// Mutually exclusive with Service
	URL *apis.URL `json:"url,omitempty"`

	// Service is a reference to a Service object where the interceptor is running
	// Mutually exclusive with URL
	Service *ServiceReference `json:"service,omitempty"`
}

var (
	defaultHTTPSPort = int32(8443)
	defaultHTTPPort  = int32(80)
)

// ServiceReference is a reference to a Service object
// with an optional path
type ServiceReference struct {
	// Name is the name of the service
	Name string `json:"name"`

	// Namespace is the namespace of the service
	Namespace string `json:"namespace"`

	// Path is an optional URL path
	// +optional
	Path string `json:"path,omitempty"`

	// Port is a valid port number
	Port *int32 `json:"port,omitempty"`
}

func formURL(scheme string, svc *ServiceReference, port *int32) *apis.URL {
	return &apis.URL{
		Scheme: scheme,
		Host:   fmt.Sprintf("%s.%s.svc:%d", svc.Name, svc.Namespace, *port),
		Path:   svc.Path,
	}
}
