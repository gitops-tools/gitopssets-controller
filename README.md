# gitopssets-controller

**NOTE**: This is a fork of github.com/weaveworks/gitopssets-controller by the original author.


GitOpsSets provide a way to declaratively generate resources in a Kubernetes cluster, generating the values to template resources from multiple sources.

## Description

The gitopssets controller provides generators for creating the inputs to templates.

The `GitOpsSet` CRD declares `generators` which are Go code which generates JSON objects from a set of input parameters.

Creating of resources in the cluster is a two-phase process, _generate_ the template inputs and _render_ the templates with the inputs.

Resources are created, updated and deleted when they are no longer rendered by the templating mechanism.

There are plenty of examples in the [./examples](./examples) directory and full
documentation in [./docs/](./docs).

## Getting Started

You’ll need a Kubernetes cluster to run against. You can use [KIND](https://sigs.k8s.io/kind) to get a local cluster for testing, or run against a remote cluster.
**Note:** Your controller will automatically use the current context in your kubeconfig file (i.e. whatever cluster `kubectl cluster-info` shows).

### Running on the cluster

1. Install Instances of Custom Resources:

```sh
kubectl apply -f config/samples/
```

2. Build and push your image to the location specified by `IMG`:

```sh
make docker-build docker-push IMG=<some-registry>/gitopssets-controller:tag
```

3. Deploy the controller to the cluster with the image specified by `IMG`:

```sh
make deploy IMG=<some-registry>/gitopssets-controller:tag
```

### Uninstall CRDs

To delete the CRDs from the cluster:

```sh
make uninstall
```

### Undeploy controller

UnDeploy the controller to the cluster:

```sh
make undeploy
```

### Make a release file

```sh
IMG=<user/repo>:$(git rev-parse --short HEAD) make manifests generate docker-build docker-push release
```

This release file can be easily applied to a cluster:

```sh
kubectl apply -f release.yaml
```

### For development purposes

You will need a bare minimum of Flux installed

```sh
flux install --components source-controller,kustomize-controller
make run
```

## Contributing

Feel free to open issues against this repository https://github.com/weaveworks/gitopssets-controller

### How it works

This project aims to follow the Kubernetes [Operator pattern](https://kubernetes.io/docs/concepts/extend-kubernetes/operator/)

It uses [Controllers](https://kubernetes.io/docs/concepts/architecture/controller/)
which provides a reconcile function responsible for synchronizing resources untile the desired state is reached on the cluster

### Test It Out

1. Install the CRDs into the cluster:

```sh
make install
```

2. Run your controller (this will run in the foreground, so switch to a new terminal if you want to leave it running):

```sh
make run
```

**NOTE:** You can also run this in one step by running: `make install run`

### Modifying the API definitions

If you are editing the API definitions, generate the manifests such as CRs or CRDs using:

```sh
make manifests
```

### Generating the API reference docs

To generate API docs run:

```sh
make api-docs
```

**NOTE:** Run `make --help` for more information on all potential `make` targets

More information can be found via the [Kubebuilder Documentation](https://book.kubebuilder.io/introduction.html)

## License

Copyright 2023.
