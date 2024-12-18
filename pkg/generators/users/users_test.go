package users

import (
	"context"
	"testing"
	"time"

	gitopssetsv1 "github.com/gitops-tools/gitopssets-controller/api/v1alpha1"
	"github.com/gitops-tools/gitopssets-controller/pkg/generators"
	"github.com/gitops-tools/gitopssets-controller/test"
	"github.com/go-logr/logr"
	"github.com/google/go-cmp/cmp"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
)

var _ generators.Generator = (*UsersGenerator)(nil)

func TestGenerate_with_no_sources(t *testing.T) {
	gen := NewGenerator(logr.Discard(), fake.NewFakeClient(), generators.DefaultClientFactory)
	got, err := gen.Generate(context.TODO(), &gitopssetsv1.GitOpsSetGenerator{}, nil)

	if err != nil {
		t.Errorf("got an error with no users: %s", err)
	}
	if got != nil {
		t.Errorf("got %v, want %v with no List generator", got, nil)
	}
}

func TestGenerateForKeycloak(t *testing.T) {
	testCases := []struct {
		name   string
		config *gitopssetsv1.KeycloakUsersGeneration
		want   []map[string]any
	}{
		{
			name:   "querying users",
			config: &gitopssetsv1.KeycloakUsersGeneration{},
			want:   []map[string]any{{"cluster": "cluster", "url": "url"}},
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {

			gen := NewGenerator(logr.Discard(), fake.NewFakeClient(), generators.DefaultClientFactory)
			got, err := gen.Generate(context.TODO(), &gitopssetsv1.GitOpsSetGenerator{
				Users: &gitopssetsv1.UsersGenerator{
					Keycloak: tt.config,
				},
			}, nil)

			test.AssertNoError(t, err)
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Fatalf("failed to generate users elements:\n%s", diff)
			}
		})
	}
}

// func TestGenerate_errors(t *testing.T) {
// 	testCases := []struct {
// 		name      string
// 		generator *gitopssetsv1.GitOpsSetGenerator
// 		wantErr   string
// 	}{
// 		{
// 			name: "bad json",
// 			generator: &gitopssetsv1.GitOpsSetGenerator{
// 				List: &gitopssetsv1.UsersGenerator{
// 					Elements: []apiextensionsv1.JSON{{Raw: []byte(`{`)}},
// 				},
// 			},
// 			wantErr: "error unmarshaling users element: unexpected end of JSON input",
// 		},
// 		{
// 			name:      "no generator",
// 			generator: nil,
// 			wantErr:   "GitOpsSet is empty",
// 		},
// 	}

// 	for _, tt := range testCases {
// 		t.Run(tt.name, func(t *testing.T) {

// 			gen := GeneratorFactory(logr.Discard(), nil)
// 			_, err := gen.Generate(context.TODO(), tt.generator, nil)

//			test.AssertErrorMatch(t, tt.wantErr, err)
//		})
//	}

func TestUsersGenerator_Interval(t *testing.T) {
	interval := time.Minute * 10
	gen := NewGenerator(logr.Discard(), fake.NewFakeClient(), generators.DefaultClientFactory)
	sg := &gitopssetsv1.GitOpsSetGenerator{
		Users: &gitopssetsv1.UsersGenerator{
			Interval: metav1.Duration{Duration: interval},
		},
	}

	d := gen.Interval(sg)

	if d != interval {
		t.Fatalf("got %#v want %#v", d, interval)
	}
}
