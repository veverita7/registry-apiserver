package server

import (
	genericapiserver "k8s.io/apiserver/pkg/server"
	"k8s.io/client-go/rest"
)

type Config struct {
	Apiserver *genericapiserver.Config
	Rest      *rest.Config
}
