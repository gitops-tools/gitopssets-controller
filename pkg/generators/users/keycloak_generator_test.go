package users

import (
	"context"
	"testing"

	"github.com/bigkevmcd/testcontainer-modules/keycloak"
	"github.com/gitops-tools/gitopssets-controller/pkg/generators"
	"github.com/gitops-tools/gitopssets-controller/test"
	"github.com/go-logr/logr"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/testcontainers/testcontainers-go"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	"k8s.io/utils/ptr"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"

	templatesv1 "github.com/gitops-tools/gitopssets-controller/api/v1alpha1"
)

func TestKeycloakUsersGeneration(t *testing.T) {
	ctx := context.TODO()
	keycloakContainer, err := keycloak.Run(ctx,
		"quay.io/keycloak/keycloak:26.0.6-0",
		keycloak.WithAdminCredentials("administrator", "secretpassword"),
	)
	test.AssertNoError(t, err)
	testcontainers.CleanupContainer(t, keycloakContainer)

	token, err := keycloakContainer.GetBearerToken(ctx, "administrator", "secretpassword")
	test.AssertNoError(t, err)
	if token == "" {
		t.Fatal("did not get a bearer token for communicating with Keycloak")
	}
	test.AssertNoError(t, keycloakContainer.CreateUser(ctx, token, keycloak.CreateUserRequest{
		Username: "testing1", Enabled: false, Firstname: "User1", Lastname: "Tested",
		Email: "testing1@example.com", EmailVerified: false}))
	test.AssertNoError(t, keycloakContainer.CreateUser(ctx, token, keycloak.CreateUserRequest{
		Username: "testing2", Enabled: true, Firstname: "User2", Lastname: "Testing",
		Email: "testing2@example.com", EmailVerified: true}))

	realmEndpoint, err := keycloakContainer.EndpointPath(ctx, "/admin/realms/master")
	test.AssertNoError(t, err)

	secret := newSecret(map[string]string{"token": token})

	administrator := map[string]any{
		"emailVerified": false,
		"firstname":     "",
		"lastname":      "",
		"email":         "",
		"username":      "administrator",
		"enabled":       true,
	}

	user1 := map[string]any{
		"emailVerified": false,
		"firstname":     "User1",
		"lastname":      "Tested",
		"email":         "testing1@example.com",
		"username":      "testing1",
		"enabled":       false,
	}

	user2 := map[string]any{
		"emailVerified": true,
		"firstname":     "User2",
		"lastname":      "Testing",
		"email":         "testing2@example.com",
		"username":      "testing2",
		"enabled":       true,
	}

	queryTests := map[string]struct {
		keycloakUsers templatesv1.KeycloakUsersConfig
		want          []map[string]any
	}{
		"querying all users": {
			keycloakUsers: templatesv1.KeycloakUsersConfig{},
			want: []map[string]any{
				administrator,
				user1,
				user2,
			},
		},
		"limiting users": {
			keycloakUsers: templatesv1.KeycloakUsersConfig{
				Limit: ptr.To(1),
			},
			want: []map[string]any{
				administrator,
			},
		},
		"limiting all users with all pages - queries all pages": {
			keycloakUsers: templatesv1.KeycloakUsersConfig{
				AllPages: true,
				Limit:    ptr.To(1),
			},
			want: []map[string]any{
				administrator,
				user1,
				user2,
			},
		},
		"querying enabled users": {
			keycloakUsers: templatesv1.KeycloakUsersConfig{
				Enabled: ptr.To(true),
			},
			want: []map[string]any{
				administrator,
				user2,
			},
		},
		"querying not enabled users": {
			keycloakUsers: templatesv1.KeycloakUsersConfig{
				Enabled: ptr.To(false),
			},
			want: []map[string]any{
				user1,
			},
		},
		"querying verified users": {
			keycloakUsers: templatesv1.KeycloakUsersConfig{
				EmailVerified: ptr.To(true),
			},
			want: []map[string]any{
				user2,
			},
		},
		"querying users by email": {
			keycloakUsers: templatesv1.KeycloakUsersConfig{
				Email: "example.com",
			},
			want: []map[string]any{
				user1,
				user2,
			},
		},
		"querying users by email with exact": {
			keycloakUsers: templatesv1.KeycloakUsersConfig{
				Email: "example.com",
				Exact: ptr.To(true),
			},
			want: nil,
		},
		"querying users by email with matching email": {
			keycloakUsers: templatesv1.KeycloakUsersConfig{
				Email: "testing1@example.com",
				Exact: ptr.To(true),
			},
			want: []map[string]any{
				user1,
			},
		},
		"querying users by first name": {
			keycloakUsers: templatesv1.KeycloakUsersConfig{
				Firstname: "User",
			},
			want: []map[string]any{
				user1,
				user2,
			},
		},
		"querying users by first name with exact": {
			keycloakUsers: templatesv1.KeycloakUsersConfig{
				Firstname: "User",
				Exact:     ptr.To(true),
			},
			want: nil,
		},
		"querying users by first name with matching name": {
			keycloakUsers: templatesv1.KeycloakUsersConfig{
				Firstname: "User1",
				Exact:     ptr.To(true),
			},
			want: []map[string]any{
				user1,
			},
		},
		"querying users by last name": {
			keycloakUsers: templatesv1.KeycloakUsersConfig{
				Lastname: "Test",
			},
			want: []map[string]any{
				user1,
				user2,
			},
		},
		"querying users by last name with exact": {
			keycloakUsers: templatesv1.KeycloakUsersConfig{
				Lastname: "Test",
				Exact:    ptr.To(true),
			},
			want: nil,
		},
		"querying users by last name with matching name": {
			keycloakUsers: templatesv1.KeycloakUsersConfig{
				Lastname: "Tested",
				Exact:    ptr.To(true),
			},
			want: []map[string]any{
				user1,
			},
		},
		"querying users by username": {
			keycloakUsers: templatesv1.KeycloakUsersConfig{
				Username: "testing",
			},
			want: []map[string]any{
				user1,
				user2,
			},
		},
		"querying users by username with exact": {
			keycloakUsers: templatesv1.KeycloakUsersConfig{
				Username: "testing",
				Exact:    ptr.To(true),
			},
			want: nil,
		},
		"querying users by username with matching name": {
			keycloakUsers: templatesv1.KeycloakUsersConfig{
				Username: "testing2",
				Exact:    ptr.To(true),
			},
			want: []map[string]any{
				user2,
			},
		},
		"searching users": {
			keycloakUsers: templatesv1.KeycloakUsersConfig{
				Search: "testing",
			},
			want: []map[string]any{
				user1,
				user2,
			},
		},
	}

	for name, tt := range queryTests {
		t.Run(name, func(t *testing.T) {
			generator := NewGenerator(logr.Discard(), newFakeClient(t, secret), generators.DefaultClientFactory)

			gsg := templatesv1.GitOpsSetGenerator{
				Users: &templatesv1.UsersGenerator{
					Keycloak: &templatesv1.KeycloakUsersGeneration{
						Endpoint: realmEndpoint,
						SecretRef: &templatesv1.LocalObjectReference{
							Name: secret.Name,
						},
						QueryConfig: &tt.keycloakUsers,
					},
				},
			}

			generated, err := generator.Generate(context.TODO(), &gsg,
				&templatesv1.GitOpsSet{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "demo-set",
						Namespace: "default",
					},
					Spec: templatesv1.GitOpsSetSpec{
						Generators: []templatesv1.GitOpsSetGenerator{
							gsg,
						},
					},
				})
			test.AssertNoError(t, err)

			if diff := cmp.Diff(tt.want, generated, cmpopts.IgnoreMapEntries(func(k string, _ any) bool {
				// We can't know the unique ID that's generated in advance so ignore it
				// for the purposes of comparison.
				return k == "id"
			})); diff != "" {
				t.Fatalf("failed to generate users: diff -want +got\n%s", diff)
			}
		})
	}
}

func newFakeClient(t *testing.T, objs ...runtime.Object) client.WithWatch {
	scheme := runtime.NewScheme()
	test.AssertNoError(t, clientgoscheme.AddToScheme(scheme))

	return fake.NewClientBuilder().
		WithScheme(scheme).
		WithRuntimeObjects(objs...).
		Build()
}

func newSecret(data map[string]string, opts ...func(*corev1.Secret)) *corev1.Secret {
	cm := &corev1.Secret{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "v1",
			Kind:       "Secret",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "demo-secret",
			Namespace: "default",
		},
		Data: dataToBytes(data),
	}

	for _, o := range opts {
		o(cm)
	}

	return cm
}

func dataToBytes(src map[string]string) map[string][]byte {
	result := map[string][]byte{}
	for k, v := range src {
		result[k] = []byte(v)
	}

	return result
}
