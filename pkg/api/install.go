package api

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apiserver/pkg/registry/rest"
	apiserver "k8s.io/apiserver/pkg/server"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	listers "k8s.io/client-go/listers/core/v1"

	registryv1alpha1 "github.com/veverita7/registry-server/api/v1alpha1"
)

const imageregistries = "imageregistries"

var (
	Scheme *runtime.Scheme
	Codecs serializer.CodecFactory
)

func init() {
	Scheme = runtime.NewScheme()
	utilruntime.Must(clientgoscheme.AddToScheme(Scheme))
	utilruntime.Must(registryv1alpha1.AddToScheme(Scheme))

	Codecs = serializer.NewCodecFactory(Scheme)
}

func Install(server *apiserver.GenericAPIServer, secretLister listers.SecretLister) error {
	apiGroupInfo := apiserver.NewDefaultAPIGroupInfo(
		registryv1alpha1.GroupVersion.Group, Scheme, metav1.ParameterCodec, Codecs)

	registry := newImageRegistry(groupResource(imageregistries), secretLister)
	storage := map[string]rest.Storage{imageregistries: registry}
	apiGroupInfo.VersionedResourcesStorageMap[registryv1alpha1.GroupVersion.Version] = storage

	return server.InstallAPIGroup(&apiGroupInfo)
}

func groupResource(rsc string) schema.GroupResource {
	return registryv1alpha1.GroupVersion.WithResource(rsc).GroupResource()
}
