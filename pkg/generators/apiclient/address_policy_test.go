package apiclient

import (
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/go-logr/logr"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"

	templatesv1 "github.com/gitops-tools/gitopssets-controller/api/v1alpha1"
)

func TestValidateEndpointAddress(t *testing.T) {
	testCases := []struct {
		name               string
		endpoint           string
		allowCluster       bool
		wantErrorSubstring string
	}{
		{name: "public address", endpoint: "https://8.8.8.8"},
		{name: "private address denied", endpoint: "https://10.0.0.1", wantErrorSubstring: "private"},
		{name: "private address allowed explicitly", endpoint: "https://10.0.0.1", allowCluster: true},
		{name: "loopback denied", endpoint: "http://127.0.0.1", wantErrorSubstring: "loopback"},
		{name: "loopback allowed explicitly", endpoint: "http://127.0.0.1", allowCluster: true},
		{name: "link local always denied", endpoint: "http://169.254.169.254", allowCluster: true, wantErrorSubstring: "restricted"},
		{name: "unspecified always denied", endpoint: "http://0.0.0.0", allowCluster: true, wantErrorSubstring: "restricted"},
		{name: "multicast always denied", endpoint: "http://224.0.0.1", allowCluster: true, wantErrorSubstring: "restricted"},
		{name: "cluster DNS denied", endpoint: "https://kubernetes.default.svc", wantErrorSubstring: "cluster"},
		{name: "cluster DNS allowed explicitly", endpoint: "https://kubernetes.default.svc", allowCluster: true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			endpoint, err := url.Parse(tc.endpoint)
			if err != nil {
				t.Fatal(err)
			}
			err = validateEndpointAddress(context.Background(), endpoint, tc.allowCluster)
			if tc.wantErrorSubstring == "" {
				if err != nil {
					t.Fatalf("validateEndpointAddress() error = %v", err)
				}
				return
			}
			if err == nil || !strings.Contains(err.Error(), tc.wantErrorSubstring) {
				t.Fatalf("validateEndpointAddress() error = %v, want containing %q", err, tc.wantErrorSubstring)
			}
		})
	}
}

func TestAllowsClusterNetworkDefaultsToFalse(t *testing.T) {
	if allowsClusterNetwork(nil) {
		t.Fatal("nil API client allowed cluster network access")
	}
	if allowsClusterNetwork(&templatesv1.APIClientGenerator{}) {
		t.Fatal("omitted allowClusterNetwork allowed cluster network access")
	}
	allow := false
	if allowsClusterNetwork(&templatesv1.APIClientGenerator{AllowClusterNetwork: &allow}) {
		t.Fatal("explicit false allowed cluster network access")
	}
	allow = true
	if !allowsClusterNetwork(&templatesv1.APIClientGenerator{AllowClusterNetwork: &allow}) {
		t.Fatal("explicit true did not allow cluster network access")
	}
}

func TestGenerateRejectsPrivateEndpointBeforeRequest(t *testing.T) {
	requests := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		requests++
		_, _ = w.Write([]byte("[]"))
	}))
	defer server.Close()

	gen := NewGenerator(logr.Discard(), fake.NewFakeClient(), DefaultClientFactory)
	_, err := gen.Generate(t.Context(), &templatesv1.GitOpsSetGenerator{
		APIClient: &templatesv1.APIClientGenerator{Endpoint: server.URL},
	}, &templatesv1.GitOpsSet{})
	if err == nil || !strings.Contains(err.Error(), "private or loopback") {
		t.Fatalf("Generate() error = %v, want private endpoint error", err)
	}
	if requests != 0 {
		t.Fatalf("server received %d requests, want none", requests)
	}
}

func TestAddressPolicyRejectsPrivateDialPeer(t *testing.T) {
	requests := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		requests++
		_, _ = w.Write([]byte("[]"))
	}))
	defer server.Close()

	httpClient := &http.Client{Transport: http.DefaultTransport.(*http.Transport).Clone()}
	if err := applyAddressPolicy(httpClient, false); err != nil {
		t.Fatal(err)
	}
	_, err := httpClient.Get(server.URL)
	if err == nil || !strings.Contains(err.Error(), "private or loopback") {
		t.Fatalf("Get() error = %v, want private peer error", err)
	}
	if requests != 0 {
		t.Fatalf("server received %d requests, want none", requests)
	}
}

func TestAddressPolicyRejectsCrossHostRedirect(t *testing.T) {
	targetRequests := 0
	target := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		targetRequests++
		_, _ = w.Write([]byte("[]"))
	}))
	defer target.Close()

	redirect := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		http.Redirect(w, req, target.URL, http.StatusFound)
	}))
	defer redirect.Close()

	allowClusterNetwork := true
	gen := NewGenerator(logr.Discard(), fake.NewFakeClient(), DefaultClientFactory)
	_, err := gen.Generate(t.Context(), &templatesv1.GitOpsSetGenerator{
		APIClient: &templatesv1.APIClientGenerator{Endpoint: redirect.URL, AllowClusterNetwork: &allowClusterNetwork},
	}, &templatesv1.GitOpsSet{})
	if err == nil || !strings.Contains(err.Error(), "different host") {
		t.Fatalf("Generate() error = %v, want cross-host redirect error", err)
	}
	if targetRequests != 0 {
		t.Fatalf("redirect target received %d requests, want none", targetRequests)
	}
}
