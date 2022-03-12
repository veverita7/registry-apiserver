package server

import (
	corev1 "k8s.io/api/core/v1"
	genericapiserver "k8s.io/apiserver/pkg/server"
	listers "k8s.io/client-go/listers/core/v1"
	"k8s.io/client-go/tools/cache"

	"github.com/veverita7/registry-server/pkg/api"
)

type Server struct {
	*genericapiserver.GenericAPIServer
	secrets cache.Controller
}

func NewServer(cnf *Config) (*Server, error) {
	informer, err := informerFactory(cnf.Rest)
	if err != nil {
		return nil, err
	}

	apiserver, err := cnf.Apiserver.Complete(informer).
		New("registry-apiserver", genericapiserver.NewEmptyDelegate())
	if err != nil {
		return nil, err
	}

	secretInformerFactory, err := secretInformerFactory(cnf.Rest)
	if err != nil {
		return nil, err
	}
	secretInformer, err := secretInformerFactory.ForResource(corev1.SchemeGroupVersion.WithResource("secrets"))
	if err != nil {
		return nil, err
	}
	secretLister := listers.NewSecretLister(secretInformer.Informer().GetIndexer())

	if err := api.Install(apiserver, secretLister); err != nil {
		return nil, err
	}
	return &Server{
		GenericAPIServer: apiserver,
		secrets:          secretInformer.Informer(),
	}, nil
}

func (s *Server) RunUntil(stopCh <-chan struct{}) error {
	go s.secrets.Run(stopCh)

	if ok := cache.WaitForCacheSync(stopCh, s.secrets.HasSynced); !ok {
		return nil
	}

	return s.GenericAPIServer.PrepareRun().Run(stopCh)
}
