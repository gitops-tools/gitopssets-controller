package generators

import (
	"crypto/tls"
	"net/http"
)

// This is used to create per API request http.Clients.
type HTTPClientFactory func(*tls.Config) *http.Client

// This is the default Client factory it returns a zero-value client.
var DefaultClientFactory = func(config *tls.Config) *http.Client {
	transport := http.DefaultTransport.(*http.Transport).Clone()
	transport.TLSClientConfig = config

	return &http.Client{
		Transport: transport,
	}
}
