package users

import (
	"context"
	"testing"

	"github.com/bigkevmcd/testcontainer-modules/keycloak"
	templatesv1 "github.com/gitops-tools/gitopssets-controller/api/v1alpha1"
	"github.com/gitops-tools/gitopssets-controller/test"
	"github.com/go-logr/logr"
	"github.com/google/go-cmp/cmp"
	"github.com/testcontainers/testcontainers-go"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
)

func TestGenerate(t *testing.T) {
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

	realmEndpoint, err := keycloakContainer.EndpointPath(ctx, "/admin/realms/master")
	test.AssertNoError(t, err)

	secret := newSecret(map[string]string{"token": token})

	generator := NewKeycloakGenerator(logr.Discard(), newFakeClient(t, secret), DefaultClientFactory)
	gsg := templatesv1.GitOpsSetGenerator{
		Users: &templatesv1.UsersGenerator{
			Keycloak: &templatesv1.KeycloakUsersGeneration{
				Endpoint: realmEndpoint,
				SecretRef: &templatesv1.LocalObjectReference{
					Name: secret.Name,
				},
				QueryConfig: &templatesv1.KeycloakUsersConfig{
					Enabled: true,
				},
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

	want := []map[string]any{
		{
			"emailVerified": false,
			"firstname":     "",
			"id":            "746968cd-bf8b-476c-bc27-780478cbdbf6",
			"lastname":      "",
			"username":      "admin",
		},
		{
			"emailVerified": true,
			"firstname":     "Kevin",
			"id":            "09871cb2-b64b-4cc3-8b7d-577467921139",
			"lastname":      "McDermott",
			"username":      "kevin",
		},
	}
	if diff := cmp.Diff(want, generated); diff != "" {
		t.Fatalf("failed to generate users:\n %s", diff)
	}
}

func TestGenerate_with_no_generator(t *testing.T) {
	t.Skip()
}

func TestGenerate_with_no_config(t *testing.T) {
	t.Skip()
}

func newFakeClient(t *testing.T, objs ...runtime.Object) client.WithWatch {
	scheme := runtime.NewScheme()
	if err := clientgoscheme.AddToScheme(scheme); err != nil {
		t.Fatal(err)
	}
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
