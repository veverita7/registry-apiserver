package server

import (
	apiserver "k8s.io/apiserver/pkg/server"
	"k8s.io/client-go/rest"
)

type Config struct {
	Apiserver *apiserver.Config
	Rest      *rest.Config
}
