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
	test.AssertNoError(t, keycloakContainer.CreateUser(ctx, token, keycloak.CreateUserRequest{Username: "testing1", Enabled: false, Firstname: "Test", Lastname: "User1"}))
	test.AssertNoError(t, keycloakContainer.CreateUser(ctx, token, keycloak.CreateUserRequest{Username: "testing2", Enabled: true, Firstname: "Test", Lastname: "User2"}))

	realmEndpoint, err := keycloakContainer.EndpointPath(ctx, "/admin/realms/master")
	test.AssertNoError(t, err)

	secret := newSecret(map[string]string{"token": token})

	queryTests := map[string]struct {
		keycloakUsers *templatesv1.KeycloakUsersGeneration
		want          []map[string]any
	}{
		"querying enabled users": {
			keycloakUsers: &templatesv1.KeycloakUsersGeneration{
				Endpoint: realmEndpoint,
				SecretRef: &templatesv1.LocalObjectReference{
					Name: secret.Name,
				},
				QueryConfig: &templatesv1.KeycloakUsersConfig{
					Enabled: true,
				},
			},
			want: []map[string]any{
				{
					"emailVerified": false,
					"firstname":     "",
					"lastname":      "",
					"username":      "administrator",
				},
				{
					"emailVerified": false,
					"firstname":     "Test",
					"lastname":      "User2",
					"username":      "testing2",
				},
			},
		},
		"querying not enabled users": {
			keycloakUsers: &templatesv1.KeycloakUsersGeneration{
				Endpoint: realmEndpoint,
				SecretRef: &templatesv1.LocalObjectReference{
					Name: secret.Name,
				},
				QueryConfig: &templatesv1.KeycloakUsersConfig{
					Enabled: false,
					// Enabled defaults to true!
				},
			},
			want: []map[string]any{
				{
					"emailVerified": false,
					"firstname":     "Test",
					"lastname":      "User1",
					"username":      "testing1",
				},
			},
		},
	}

	for name, tt := range queryTests {
		t.Run(name, func(t *testing.T) {
			generator := NewGenerator(logr.Discard(), newFakeClient(t, secret), generators.DefaultClientFactory)

			gsg := templatesv1.GitOpsSetGenerator{
				Users: &templatesv1.UsersGenerator{
					Keycloak: tt.keycloakUsers,
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
