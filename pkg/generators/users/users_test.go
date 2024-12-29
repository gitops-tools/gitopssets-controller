package users

import (
	"context"
	"testing"
	"time"

	gitopssetsv1 "github.com/gitops-tools/gitopssets-controller/api/v1alpha1"
	"github.com/gitops-tools/gitopssets-controller/pkg/generators"
	"github.com/go-logr/logr"
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
