package cluster

import (
	"context"
	"fmt"
	"time"

	templatesv1 "github.com/gitops-tools/gitopssets-controller/api/v1alpha1"
	"github.com/gitops-tools/gitopssets-controller/pkg/generators"
	"github.com/go-logr/logr"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// ClusterGenerator is a generator that queries for Cluster types and exposes
// the object metadata, labels and annotations.
type ClusterGenerator struct {
	Client client.Reader
	logr.Logger
	ClusterListGVK schema.GroupVersionKind
}

// GeneratorFactory is a function for creating per-reconciliation generators for
// the ClusterGenerator.
func GeneratorFactory(gvk schema.GroupVersionKind) generators.GeneratorFactory {
	return func(l logr.Logger, c client.Reader) generators.Generator {
		return NewGenerator(l, c, gvk)
	}
}

// NewGenerator creates and returns a new cluster generator.
func NewGenerator(l logr.Logger, c client.Reader, gvk schema.GroupVersionKind) *ClusterGenerator {
	return &ClusterGenerator{
		Client:         c,
		Logger:         l,
		ClusterListGVK: gvk,
	}
}

func (g *ClusterGenerator) Generate(ctx context.Context, sg *templatesv1.GitOpsSetGenerator, ks *templatesv1.GitOpsSet) ([]map[string]any, error) {
	if sg == nil {
		return nil, generators.ErrEmptyGitOpsSet
	}

	if sg.Cluster == nil {
		return nil, nil
	}
	g.Logger.Info("generating params from Cluster generator")

	selector, err := metav1.LabelSelectorAsSelector(&sg.Cluster.Selector)
	if err != nil {
		return nil, fmt.Errorf("unable to convert selector: %w", err)
	}

	listOptions := client.ListOptions{LabelSelector: selector}

	clusterList := metav1.PartialObjectMetadataList{}
	clusterList.SetGroupVersionKind(g.ClusterListGVK)
	err = g.Client.List(ctx, &clusterList, &listOptions)
	if err != nil {
		return nil, err
	}

	var paramsList []map[string]any
	for _, cluster := range clusterList.Items {
		params := map[string]any{
			"ClusterName":        cluster.Name,
			"ClusterNamespace":   cluster.Namespace,
			"ClusterLabels":      mapOrEmptyMap(cluster.Labels),
			"ClusterAnnotations": mapOrEmptyMap(cluster.Annotations),
		}
		paramsList = append(paramsList, params)
	}

	return paramsList, nil
}

func mapOrEmptyMap(src map[string]string) map[string]any {
	if src == nil {
		return map[string]any{}
	}

	result := map[string]any{}

	for k, v := range src {
		result[k] = v
	}

	return result
}

// Interval is an implementation of the Generator interface.
//
// This does not requeue because the resources should be watched.
func (g *ClusterGenerator) Interval(sg *templatesv1.GitOpsSetGenerator) time.Duration {
	return generators.NoRequeueInterval
}
