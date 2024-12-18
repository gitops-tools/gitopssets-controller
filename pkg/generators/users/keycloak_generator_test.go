package users

import (
	"context"
	"os"
	"testing"

	"github.com/go-logr/logr"
	"github.com/google/go-cmp/cmp"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
)

func TestGenerate(t *testing.T) {
	secret := newSecret(map[string]string{"token": os.Getenv("BEARER_TOKEN")})
	config := KeycloakGeneratorConfig{
		APIEndpoint: "http://localhost:8080/admin/realms/master",
		SecretRef: corev1.LocalObjectReference{
			Name: secret.Name,
		},
		Users: &KeycloakGeneratorUsers{
			Enabled: true,
		},
	}

	generator := NewKeycloakGenerator(logr.Discard(), newFakeClient(t, secret), DefaultClientFactory)
	generated, err := generator.Generate(context.TODO(), config)
	if err != nil {
		t.Fatal(err)
	}

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
			Namespace: "testing",
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
