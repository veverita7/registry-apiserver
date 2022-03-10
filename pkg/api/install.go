package api

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	apiserver "k8s.io/apiserver/pkg/server"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"

	registryv1alpha1 "github.com/veverita7/registry-server/api/v1alpha1"
)

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

func Install(server *apiserver.GenericAPIServer) error {
	apiGroupInfo := apiserver.NewDefaultAPIGroupInfo(
		registryv1alpha1.GroupVersion.Group, Scheme, metav1.ParameterCodec, Codecs)
	return server.InstallAPIGroup(&apiGroupInfo)
}
