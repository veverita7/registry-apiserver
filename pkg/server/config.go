package server

import (
	"net"

	"k8s.io/apimachinery/pkg/runtime/serializer"
	k8sapiserver "k8s.io/apiserver/pkg/server"
	"k8s.io/client-go/rest"
	"k8s.io/component-base/version"
	ctrl "sigs.k8s.io/controller-runtime"

	"github.com/veverita7/registry-server/cmd/registry-apiserver/options"
	"github.com/veverita7/registry-server/pkg/api"
)

type Config struct {
	Apiserver *k8sapiserver.Config
	Rest      *rest.Config
}

func NewConfig(opts *options.Options) (*Config, error) {
	apiserver, err := apiserverConfig(opts)
	if err != nil {
		return nil, err
	}

	rest, err := ctrl.GetConfig()
	if err != nil {
		return nil, err
	}

	cnf := &Config{
		Apiserver: apiserver,
		Rest:      rest,
	}
	return cnf, nil
}

func apiserverConfig(opts *options.Options) (*k8sapiserver.Config, error) {
	if err := opts.SecureServing.MaybeDefaultWithSelfSignedCerts(
		"localhost", nil, []net.IP{net.ParseIP("127.0.0.1")}); err != nil {
		return nil, err
	}

	cnf := k8sapiserver.NewConfig(serializer.NewCodecFactory(api.Scheme))
	if err := opts.SecureServing.ApplyTo(&cnf.SecureServing, &cnf.LoopbackClientConfig); err != nil {
		return nil, err
	}

	if !opts.DisableAuthForTesting {
		if err := opts.Authentication.ApplyTo(&cnf.Authentication, cnf.SecureServing, nil); err != nil {
			return nil, err
		}
		if err := opts.Authorization.ApplyTo(&cnf.Authorization); err != nil {
			return nil, err
		}
	}

	versionGet := version.Get()
	cnf.Version = &versionGet

	return cnf, nil
}
