//go:build !ignore_autogenerated
// +build !ignore_autogenerated

/*
Copyright 2023.
*/

// Code generated by controller-gen. DO NOT EDIT.

package v1alpha1

import (
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *APIClientGenerator) DeepCopyInto(out *APIClientGenerator) {
	*out = *in
	out.Interval = in.Interval
	if in.HeadersRef != nil {
		in, out := &in.HeadersRef, &out.HeadersRef
		*out = new(HeadersReference)
		**out = **in
	}
	if in.Body != nil {
		in, out := &in.Body, &out.Body
		*out = new(v1.JSON)
		(*in).DeepCopyInto(*out)
	}
	if in.SecretRef != nil {
		in, out := &in.SecretRef, &out.SecretRef
		*out = new(corev1.LocalObjectReference)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new APIClientGenerator.
func (in *APIClientGenerator) DeepCopy() *APIClientGenerator {
	if in == nil {
		return nil
	}
	out := new(APIClientGenerator)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ClusterGenerator) DeepCopyInto(out *ClusterGenerator) {
	*out = *in
	in.Selector.DeepCopyInto(&out.Selector)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ClusterGenerator.
func (in *ClusterGenerator) DeepCopy() *ClusterGenerator {
	if in == nil {
		return nil
	}
	out := new(ClusterGenerator)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ConfigGenerator) DeepCopyInto(out *ConfigGenerator) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ConfigGenerator.
func (in *ConfigGenerator) DeepCopy() *ConfigGenerator {
	if in == nil {
		return nil
	}
	out := new(ConfigGenerator)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *GitOpsSet) DeepCopyInto(out *GitOpsSet) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	in.Status.DeepCopyInto(&out.Status)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new GitOpsSet.
func (in *GitOpsSet) DeepCopy() *GitOpsSet {
	if in == nil {
		return nil
	}
	out := new(GitOpsSet)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *GitOpsSet) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *GitOpsSetGenerator) DeepCopyInto(out *GitOpsSetGenerator) {
	*out = *in
	if in.List != nil {
		in, out := &in.List, &out.List
		*out = new(ListGenerator)
		(*in).DeepCopyInto(*out)
	}
	if in.PullRequests != nil {
		in, out := &in.PullRequests, &out.PullRequests
		*out = new(PullRequestGenerator)
		(*in).DeepCopyInto(*out)
	}
	if in.GitRepository != nil {
		in, out := &in.GitRepository, &out.GitRepository
		*out = new(GitRepositoryGenerator)
		(*in).DeepCopyInto(*out)
	}
	if in.OCIRepository != nil {
		in, out := &in.OCIRepository, &out.OCIRepository
		*out = new(OCIRepositoryGenerator)
		(*in).DeepCopyInto(*out)
	}
	if in.Matrix != nil {
		in, out := &in.Matrix, &out.Matrix
		*out = new(MatrixGenerator)
		(*in).DeepCopyInto(*out)
	}
	if in.Cluster != nil {
		in, out := &in.Cluster, &out.Cluster
		*out = new(ClusterGenerator)
		(*in).DeepCopyInto(*out)
	}
	if in.APIClient != nil {
		in, out := &in.APIClient, &out.APIClient
		*out = new(APIClientGenerator)
		(*in).DeepCopyInto(*out)
	}
	if in.ImagePolicy != nil {
		in, out := &in.ImagePolicy, &out.ImagePolicy
		*out = new(ImagePolicyGenerator)
		**out = **in
	}
	if in.Config != nil {
		in, out := &in.Config, &out.Config
		*out = new(ConfigGenerator)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new GitOpsSetGenerator.
func (in *GitOpsSetGenerator) DeepCopy() *GitOpsSetGenerator {
	if in == nil {
		return nil
	}
	out := new(GitOpsSetGenerator)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *GitOpsSetList) DeepCopyInto(out *GitOpsSetList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]GitOpsSet, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new GitOpsSetList.
func (in *GitOpsSetList) DeepCopy() *GitOpsSetList {
	if in == nil {
		return nil
	}
	out := new(GitOpsSetList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *GitOpsSetList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *GitOpsSetNestedGenerator) DeepCopyInto(out *GitOpsSetNestedGenerator) {
	*out = *in
	if in.List != nil {
		in, out := &in.List, &out.List
		*out = new(ListGenerator)
		(*in).DeepCopyInto(*out)
	}
	if in.GitRepository != nil {
		in, out := &in.GitRepository, &out.GitRepository
		*out = new(GitRepositoryGenerator)
		(*in).DeepCopyInto(*out)
	}
	if in.OCIRepository != nil {
		in, out := &in.OCIRepository, &out.OCIRepository
		*out = new(OCIRepositoryGenerator)
		(*in).DeepCopyInto(*out)
	}
	if in.PullRequests != nil {
		in, out := &in.PullRequests, &out.PullRequests
		*out = new(PullRequestGenerator)
		(*in).DeepCopyInto(*out)
	}
	if in.Cluster != nil {
		in, out := &in.Cluster, &out.Cluster
		*out = new(ClusterGenerator)
		(*in).DeepCopyInto(*out)
	}
	if in.APIClient != nil {
		in, out := &in.APIClient, &out.APIClient
		*out = new(APIClientGenerator)
		(*in).DeepCopyInto(*out)
	}
	if in.ImagePolicy != nil {
		in, out := &in.ImagePolicy, &out.ImagePolicy
		*out = new(ImagePolicyGenerator)
		**out = **in
	}
	if in.Config != nil {
		in, out := &in.Config, &out.Config
		*out = new(ConfigGenerator)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new GitOpsSetNestedGenerator.
func (in *GitOpsSetNestedGenerator) DeepCopy() *GitOpsSetNestedGenerator {
	if in == nil {
		return nil
	}
	out := new(GitOpsSetNestedGenerator)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *GitOpsSetSpec) DeepCopyInto(out *GitOpsSetSpec) {
	*out = *in
	if in.Generators != nil {
		in, out := &in.Generators, &out.Generators
		*out = make([]GitOpsSetGenerator, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	if in.Templates != nil {
		in, out := &in.Templates, &out.Templates
		*out = make([]GitOpsSetTemplate, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new GitOpsSetSpec.
func (in *GitOpsSetSpec) DeepCopy() *GitOpsSetSpec {
	if in == nil {
		return nil
	}
	out := new(GitOpsSetSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *GitOpsSetStatus) DeepCopyInto(out *GitOpsSetStatus) {
	*out = *in
	out.ReconcileRequestStatus = in.ReconcileRequestStatus
	if in.Conditions != nil {
		in, out := &in.Conditions, &out.Conditions
		*out = make([]metav1.Condition, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	if in.Inventory != nil {
		in, out := &in.Inventory, &out.Inventory
		*out = new(ResourceInventory)
		(*in).DeepCopyInto(*out)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new GitOpsSetStatus.
func (in *GitOpsSetStatus) DeepCopy() *GitOpsSetStatus {
	if in == nil {
		return nil
	}
	out := new(GitOpsSetStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *GitOpsSetTemplate) DeepCopyInto(out *GitOpsSetTemplate) {
	*out = *in
	in.Content.DeepCopyInto(&out.Content)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new GitOpsSetTemplate.
func (in *GitOpsSetTemplate) DeepCopy() *GitOpsSetTemplate {
	if in == nil {
		return nil
	}
	out := new(GitOpsSetTemplate)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *GitRepositoryGenerator) DeepCopyInto(out *GitRepositoryGenerator) {
	*out = *in
	if in.Files != nil {
		in, out := &in.Files, &out.Files
		*out = make([]RepositoryGeneratorFileItem, len(*in))
		copy(*out, *in)
	}
	if in.Directories != nil {
		in, out := &in.Directories, &out.Directories
		*out = make([]RepositoryGeneratorDirectoryItem, len(*in))
		copy(*out, *in)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new GitRepositoryGenerator.
func (in *GitRepositoryGenerator) DeepCopy() *GitRepositoryGenerator {
	if in == nil {
		return nil
	}
	out := new(GitRepositoryGenerator)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *HeadersReference) DeepCopyInto(out *HeadersReference) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new HeadersReference.
func (in *HeadersReference) DeepCopy() *HeadersReference {
	if in == nil {
		return nil
	}
	out := new(HeadersReference)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ImagePolicyGenerator) DeepCopyInto(out *ImagePolicyGenerator) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ImagePolicyGenerator.
func (in *ImagePolicyGenerator) DeepCopy() *ImagePolicyGenerator {
	if in == nil {
		return nil
	}
	out := new(ImagePolicyGenerator)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ListGenerator) DeepCopyInto(out *ListGenerator) {
	*out = *in
	if in.Elements != nil {
		in, out := &in.Elements, &out.Elements
		*out = make([]v1.JSON, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ListGenerator.
func (in *ListGenerator) DeepCopy() *ListGenerator {
	if in == nil {
		return nil
	}
	out := new(ListGenerator)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *MatrixGenerator) DeepCopyInto(out *MatrixGenerator) {
	*out = *in
	if in.Generators != nil {
		in, out := &in.Generators, &out.Generators
		*out = make([]GitOpsSetNestedGenerator, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new MatrixGenerator.
func (in *MatrixGenerator) DeepCopy() *MatrixGenerator {
	if in == nil {
		return nil
	}
	out := new(MatrixGenerator)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *OCIRepositoryGenerator) DeepCopyInto(out *OCIRepositoryGenerator) {
	*out = *in
	if in.Files != nil {
		in, out := &in.Files, &out.Files
		*out = make([]RepositoryGeneratorFileItem, len(*in))
		copy(*out, *in)
	}
	if in.Directories != nil {
		in, out := &in.Directories, &out.Directories
		*out = make([]RepositoryGeneratorDirectoryItem, len(*in))
		copy(*out, *in)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new OCIRepositoryGenerator.
func (in *OCIRepositoryGenerator) DeepCopy() *OCIRepositoryGenerator {
	if in == nil {
		return nil
	}
	out := new(OCIRepositoryGenerator)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *PullRequestGenerator) DeepCopyInto(out *PullRequestGenerator) {
	*out = *in
	out.Interval = in.Interval
	if in.SecretRef != nil {
		in, out := &in.SecretRef, &out.SecretRef
		*out = new(corev1.LocalObjectReference)
		**out = **in
	}
	if in.Labels != nil {
		in, out := &in.Labels, &out.Labels
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new PullRequestGenerator.
func (in *PullRequestGenerator) DeepCopy() *PullRequestGenerator {
	if in == nil {
		return nil
	}
	out := new(PullRequestGenerator)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *RepositoryGeneratorDirectoryItem) DeepCopyInto(out *RepositoryGeneratorDirectoryItem) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new RepositoryGeneratorDirectoryItem.
func (in *RepositoryGeneratorDirectoryItem) DeepCopy() *RepositoryGeneratorDirectoryItem {
	if in == nil {
		return nil
	}
	out := new(RepositoryGeneratorDirectoryItem)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *RepositoryGeneratorFileItem) DeepCopyInto(out *RepositoryGeneratorFileItem) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new RepositoryGeneratorFileItem.
func (in *RepositoryGeneratorFileItem) DeepCopy() *RepositoryGeneratorFileItem {
	if in == nil {
		return nil
	}
	out := new(RepositoryGeneratorFileItem)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ResourceInventory) DeepCopyInto(out *ResourceInventory) {
	*out = *in
	if in.Entries != nil {
		in, out := &in.Entries, &out.Entries
		*out = make([]ResourceRef, len(*in))
		copy(*out, *in)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ResourceInventory.
func (in *ResourceInventory) DeepCopy() *ResourceInventory {
	if in == nil {
		return nil
	}
	out := new(ResourceInventory)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ResourceRef) DeepCopyInto(out *ResourceRef) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ResourceRef.
func (in *ResourceRef) DeepCopy() *ResourceRef {
	if in == nil {
		return nil
	}
	out := new(ResourceRef)
	in.DeepCopyInto(out)
	return out
}
