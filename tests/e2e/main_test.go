package tests

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	imagev1 "github.com/fluxcd/image-reflector-controller/api/v1beta2"
	kustomizev1 "github.com/fluxcd/kustomize-controller/api/v1beta2"
	"github.com/fluxcd/pkg/http/fetch"
	"github.com/fluxcd/pkg/runtime/testenv"
	"github.com/fluxcd/pkg/tar"
	sourcev1 "github.com/fluxcd/source-controller/api/v1"
	sourcev1beta2 "github.com/fluxcd/source-controller/api/v1beta2"
	"github.com/gitops-tools/gitopssets-controller/test"
	clustersv1 "github.com/weaveworks/cluster-controller/api/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/client-go/kubernetes/scheme"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client/apiutil"

	gitopssetsv1 "github.com/gitops-tools/gitopssets-controller/api/v1alpha1"
	"github.com/gitops-tools/gitopssets-controller/controllers"
	"github.com/gitops-tools/gitopssets-controller/pkg/generators"
	"github.com/gitops-tools/gitopssets-controller/pkg/generators/cluster"
	"github.com/gitops-tools/gitopssets-controller/pkg/generators/config"
	"github.com/gitops-tools/gitopssets-controller/pkg/generators/gitrepository"
	"github.com/gitops-tools/gitopssets-controller/pkg/generators/imagepolicy"
	"github.com/gitops-tools/gitopssets-controller/pkg/generators/list"
	"github.com/gitops-tools/gitopssets-controller/pkg/generators/matrix"
	"github.com/gitops-tools/gitopssets-controller/pkg/generators/ocirepository"
	"github.com/gitops-tools/gitopssets-controller/pkg/generators/pullrequests"
)

const (
	timeout = 10 * time.Second
)

var (
	testEnv       *testenv.Environment
	ctx           = ctrl.SetupSignalHandler()
	eventRecorder *test.FakeEventRecorder
)

func TestMain(m *testing.M) {
	builder := runtime.NewSchemeBuilder(
		gitopssetsv1.AddToScheme,
		clustersv1.AddToScheme,
		sourcev1.AddToScheme,
		sourcev1beta2.AddToScheme,
		imagev1.AddToScheme,
		kustomizev1.AddToScheme,
		metav1.AddMetaToScheme,
		clientgoscheme.AddToScheme,
	)
	utilruntime.Must(builder.AddToScheme(scheme.Scheme))

	fetcher := fetch.NewArchiveFetcher(1, tar.UnlimitedUntarSize, tar.UnlimitedUntarSize, "")

	testEnv = testenv.New(testenv.WithCRDPath(filepath.Join("..", "..", "config", "crd", "bases"),
		filepath.Join("..", "..", "controllers", "testdata", "crds"),
		filepath.Join("testdata", "crds"),
	))

	httpClient, err := rest.HTTPClientFor(testEnv.GetConfig())
	if err != nil {
		panic(fmt.Sprintf("failed to get an HTTP client for the server: %s", err))
	}
	mapper, err := apiutil.NewDynamicRESTMapper(testEnv.GetConfig(), httpClient)
	if err != nil {
		panic(fmt.Sprintf("failed to create RESTMapper: %s", err))
	}

	eventRecorder = &test.FakeEventRecorder{}
	if err := (&controllers.GitOpsSetReconciler{
		Client: testEnv,
		Config: testEnv.GetConfig(),
		Scheme: testEnv.GetScheme(),
		Mapper: mapper,
		Generators: map[string]generators.GeneratorFactory{
			"List": list.GeneratorFactory,
			"Matrix": matrix.GeneratorFactory(map[string]generators.GeneratorFactory{
				"List":          list.GeneratorFactory,
				"GitRepository": gitrepository.GeneratorFactory(fetcher),
				"OCIRepository": ocirepository.GeneratorFactory(fetcher),
				"PullRequests":  pullrequests.GeneratorFactory,
				"ImagePolicy":   imagepolicy.GeneratorFactory,
				"Config":        config.GeneratorFactory,
			}),
			"PullRequests":  pullrequests.GeneratorFactory,
			"OCIRepository": ocirepository.GeneratorFactory(fetcher),
			"Cluster":       cluster.GeneratorFactory(gitopssetsv1.GitopsClusterListGVK),
			"ImagePolicy":   imagepolicy.GeneratorFactory,
			"Config":        config.GeneratorFactory,
		},
		EventRecorder: eventRecorder,
	}).SetupWithManager(testEnv); err != nil {
		panic(fmt.Sprintf("Failed to start GitOpsSetReconciler: %v", err))
	}

	go func() {
		fmt.Println("Starting the test environment")
		if err := testEnv.Start(ctx); err != nil {
			panic(fmt.Sprintf("Failed to start the test environment manager: %v", err))
		}
	}()
	<-testEnv.Manager.Elected()

	code := m.Run()

	fmt.Println("Stopping the test environment")
	if err := testEnv.Stop(); err != nil {
		panic(fmt.Sprintf("Failed to stop the test environment: %v", err))
	}

	os.Exit(code)
}
