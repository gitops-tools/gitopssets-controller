package users

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"

	templatesv1 "github.com/gitops-tools/gitopssets-controller/api/v1alpha1"
	"github.com/go-logr/logr"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// KeycloakUser represents users in the Keycloak API.
type KeycloakUser struct {
	ID            string              `json:"id"`
	Username      string              `json:"username"`
	Firstname     string              `json:"firstName"`
	Lastname      string              `json:"lastName"`
	EmailVerified bool                `json:"emailVerified"`
	Enabled       bool                `json:"enabled"`
	Email         string              `json:"email"`
	Attributes    map[string][]string `json:"attributes,omitempty"`

	// TODO: This is a unix-timestamp - fix this!
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
		"email":         u.Email,
		"emailVerified": u.EmailVerified,
		"enabled":       u.Enabled,
		"attributes":    u.Attributes,

		// TODO: CreatedTimestamp in a user-readable format.
	}
}

// https://www.keycloak.org/docs-api/latest/rest-api/index.html#_users
func (k *UsersGenerator) generateKeycloakUsers(ctx context.Context, sg *templatesv1.GitOpsSetGenerator, ks *templatesv1.GitOpsSet) ([]map[string]any, error) {
	secretName := types.NamespacedName{
		Namespace: ks.GetNamespace(),
		Name:      sg.Users.Keycloak.SecretRef.Name,
	}
	authToken, err := getSecretToken(ctx, secretName, k.Client)
	if err != nil {
		// TODO: Improve this error
		return nil, err
	}

	// TODO: Handle nil QueryConfig!

	query := url.Values{}
	if sg.Users.Keycloak.QueryConfig != nil {
		query = sg.Users.Keycloak.QueryConfig.ToValues()
	}

	// TODO: This should allow customisation of the TLS setup

	httpClient := k.ClientFactory(nil)

	pageNumber := 0
	var combinedResult []map[string]any
	for {
		if pageNumber > 0 {
			query["first"] = []string{strconv.Itoa(pageNumber)}
		}

		result, err := getUsers(ctx, k.Logger, httpClient, sg.Users.Keycloak.Endpoint, authToken, query)
		if err != nil {
			return nil, err
		}

		if len(result) == 0 {
			break
		}

		combinedResult = append(combinedResult, result...)

		if sg.Users.Keycloak.QueryConfig == nil || (sg.Users.Keycloak.QueryConfig != nil && !sg.Users.Keycloak.QueryConfig.AllPages) {
			break
		}

		pageNumber++
	}

	return combinedResult, nil
}

func getUsers(ctx context.Context, logger logr.Logger, client *http.Client, endpoint, authToken string, query url.Values) ([]map[string]any, error) {
	// TODO: Improve the URL generation
	queryURL := endpoint + "/users?" + query.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, queryURL, nil)
	if err != nil {
		// TODO: Improve error
		return nil, err
	}
	req.Header.Add("Authorization", "Bearer "+authToken)

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed Keycloak HTTP request: %w", err)
	}

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("reading response body: %w", err)
	}

	if resp.StatusCode > http.StatusOK {
		logger.Info("invalid response from Keycloak", "status", resp.StatusCode, "queryURL", queryURL, "body", string(b))
		return nil, fmt.Errorf("invalid response from %s: %v", endpoint, resp.StatusCode)
	}

	decoder := json.NewDecoder(bytes.NewReader(b))
	// TODO: Wrap this in a function that can report the error.
	defer resp.Body.Close()

	var users []KeycloakUser
	if err := decoder.Decode(&users); err != nil {
		return nil, fmt.Errorf("parsing JSON from response: %w", err)
	}

	var result []map[string]any
	for _, user := range users {
		result = append(result, user.ToMap())
	}

	return result, nil
}

func getSecretToken(ctx context.Context, secretName types.NamespacedName, secretClient client.Reader) (string, error) {
	var secret corev1.Secret
	if err := secretClient.Get(ctx, secretName, &secret); err != nil {
		return "", fmt.Errorf("failed to load keycloak credentials: %w", err)
	}
	data, ok := secret.Data["token"]
	if !ok {
		return "", fmt.Errorf("secret %s does not contain required field 'token'", secretName)
	}

	return string(data), nil
}
