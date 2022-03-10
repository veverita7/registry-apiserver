package server

import (
	apiserver "k8s.io/apiserver/pkg/server"

	"github.com/veverita7/registry-server/pkg/api"
)

type Server struct {
	*apiserver.GenericAPIServer
}

func NewServer(cnf *Config) (*Server, error) {
	informer, err := informerFactory(cnf.Rest)
	if err != nil {
		return nil, err
	}

	apiserver, err := cnf.Apiserver.Complete(informer).
		New("registry-apiserver", apiserver.NewEmptyDelegate())
	if err != nil {
		return nil, err
	}

	if err := api.Install(apiserver); err != nil {
		return nil, err
	}
	return &Server{GenericAPIServer: apiserver}, nil
}

func (s *Server) RunUntil(stopCh <-chan struct{}) error {
	return s.GenericAPIServer.PrepareRun().Run(stopCh)
}
