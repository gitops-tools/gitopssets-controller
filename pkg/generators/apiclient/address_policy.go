package apiclient

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"
)

func validateEndpointAddress(ctx context.Context, endpoint *url.URL, allowClusterNetwork bool) error {
	host := strings.TrimSuffix(strings.ToLower(endpoint.Hostname()), ".")
	if host == "" {
		return fmt.Errorf("refusing API endpoint with an empty host")
	}
	if isClusterDNS(host) {
		if !allowClusterNetwork {
			return fmt.Errorf("refusing cluster address %s", host)
		}
		return nil
	}
	if ip := net.ParseIP(host); ip != nil {
		return checkEndpointIP(ip, allowClusterNetwork, host)
	}

	ips, err := net.DefaultResolver.LookupIP(ctx, "ip", host)
	if err != nil {
		return fmt.Errorf("failed to resolve API endpoint host %s: %w", host, err)
	}
	if len(ips) == 0 {
		return fmt.Errorf("API endpoint host %s resolved to no addresses", host)
	}
	for _, ip := range ips {
		if err := checkEndpointIP(ip, allowClusterNetwork, host); err != nil {
			return err
		}
	}
	return nil
}

func isClusterDNS(host string) bool {
	return host == "kubernetes" || strings.HasSuffix(host, ".svc") || strings.HasSuffix(host, ".cluster.local")
}

func checkEndpointIP(ip net.IP, allowClusterNetwork bool, host string) error {
	if ip == nil || ip.IsUnspecified() || ip.IsMulticast() || ip.IsLinkLocalUnicast() || ip.IsLinkLocalMulticast() {
		return fmt.Errorf("refusing restricted API endpoint address %s", host)
	}
	if ip.IsPrivate() || ip.IsLoopback() {
		if !allowClusterNetwork {
			return fmt.Errorf("refusing private or loopback API endpoint address %s", host)
		}
	}
	if !ip.IsGlobalUnicast() && !ip.IsPrivate() && !ip.IsLoopback() {
		return fmt.Errorf("refusing non-unicast API endpoint address %s", host)
	}
	return nil
}

func applyAddressPolicy(httpClient *http.Client, allowClusterNetwork bool) error {
	if httpClient == nil {
		return fmt.Errorf("API client factory returned a nil HTTP client")
	}

	transport := http.DefaultTransport.(*http.Transport)
	if httpClient.Transport != nil {
		var ok bool
		transport, ok = httpClient.Transport.(*http.Transport)
		if !ok {
			return fmt.Errorf("API client transport %T cannot enforce destination policy", httpClient.Transport)
		}
	}
	transport = transport.Clone()
	// A proxy hides the request destination from the dialer, so policy-enforced
	// clients connect directly to the validated endpoint.
	transport.Proxy = nil

	baseDial := transport.DialContext
	if baseDial == nil {
		dialer := &net.Dialer{Timeout: 30 * time.Second}
		baseDial = dialer.DialContext
	}
	transport.DialContext = func(ctx context.Context, network, address string) (net.Conn, error) {
		conn, err := baseDial(ctx, network, address)
		if err != nil {
			return nil, err
		}
		remote, ok := conn.RemoteAddr().(*net.TCPAddr)
		if !ok {
			conn.Close()
			return nil, fmt.Errorf("refusing non-TCP API endpoint connection")
		}
		if err := checkEndpointIP(remote.IP, allowClusterNetwork, remote.IP.String()); err != nil {
			conn.Close()
			return nil, err
		}
		return conn, nil
	}
	httpClient.Transport = transport

	previousCheckRedirect := httpClient.CheckRedirect
	httpClient.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		if len(via) > 0 && !strings.EqualFold(req.URL.Host, via[0].URL.Host) {
			return fmt.Errorf("refusing API endpoint redirect to a different host")
		}
		if err := validateEndpointAddress(req.Context(), req.URL, allowClusterNetwork); err != nil {
			return err
		}
		if previousCheckRedirect != nil {
			return previousCheckRedirect(req, via)
		}
		if len(via) >= 10 {
			return fmt.Errorf("stopped after 10 API endpoint redirects")
		}
		return nil
	}
	return nil
}
