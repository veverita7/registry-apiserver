package server

import (
	genericapiserver "k8s.io/apiserver/pkg/server"
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

	secrets := informer.Core().V1().Secrets()

	apiserver, err := cnf.Apiserver.Complete(informer).
		New("registry-apiserver", genericapiserver.NewEmptyDelegate())
	if err != nil {
		return nil, err
	}

	if err := api.Install(apiserver, secrets.Lister()); err != nil {
		return nil, err
	}
	return &Server{
		GenericAPIServer: apiserver,
		secrets:          secrets.Informer(),
	}, nil
}

func (s *Server) RunUntil(stopCh <-chan struct{}) error {
	go s.secrets.Run(stopCh)

	if ok := cache.WaitForCacheSync(stopCh, s.secrets.HasSynced); !ok {
		return nil
	}

	return s.GenericAPIServer.PrepareRun().Run(stopCh)
}
