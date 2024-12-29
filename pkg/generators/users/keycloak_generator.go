package users

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	templatesv1 "github.com/gitops-tools/gitopssets-controller/api/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
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
	Enabled       bool   `json:"enabled"`

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

// // FetchConfig configures how to load the results from Keycloak.
// type FetchConfig struct {
// 	// TODO: Mark this as default to true.
// 	AllPages bool `json:"allPages"`
// 	PageSize int  `json:"pageSize"`
// }

func (k *UsersGenerator) generateKeycloakUsers(ctx context.Context, sg *templatesv1.GitOpsSetGenerator, ks *templatesv1.GitOpsSet) ([]map[string]any, error) {
	// TODO: Standard validation checks

	secretName := types.NamespacedName{
		Namespace: ks.GetNamespace(),
		Name:      sg.Users.Keycloak.SecretRef.Name,
	}
	authToken, err := getSecretToken(ctx, secretName, k.Client)
	if err != nil {
		// TODO: Improve this error
		return nil, err
	}

	// TODO: Improve the URL generation
	req, err := http.NewRequest(http.MethodGet, sg.Users.Keycloak.Endpoint+"/users?"+sg.Users.Keycloak.QueryConfig.ToValues().Encode(), nil)
	if err != nil {
		// TODO: Improve error
		return nil, err
	}
	req.Header.Add("Authorization", "Bearer "+authToken)
	client := k.ClientFactory(nil)

	resp, err := client.Do(req)
	if err != nil {
		// TODO: Improve error
		return nil, err
	}

	if resp.StatusCode > 300 {
		return nil, fmt.Errorf("invalid response from %s: %v", sg.Users.Keycloak.Endpoint, resp.StatusCode)
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

func getSecretToken(ctx context.Context, secretName types.NamespacedName, secretClient client.Reader) (string, error) {
	var secret corev1.Secret
	if err := secretClient.Get(ctx, secretName, &secret); err != nil {
		return "", fmt.Errorf("failed to load repository generator credentials: %w", err)
	}
	data, ok := secret.Data["token"]
	if !ok {
		return "", fmt.Errorf("secret %s does not contain required field 'token'", secretName)
	}

	return string(data), nil
}
