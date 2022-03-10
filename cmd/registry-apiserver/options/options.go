package options

import (
	"net"

	apiserver "k8s.io/apiserver/pkg/server"
	apiserveroptions "k8s.io/apiserver/pkg/server/options"
	"k8s.io/component-base/cli/flag"
	"k8s.io/component-base/version"
	ctrl "sigs.k8s.io/controller-runtime"

	"github.com/veverita7/registry-server/pkg/api"
	"github.com/veverita7/registry-server/pkg/server"
)

type Options struct {
	SecureServing  *apiserveroptions.SecureServingOptionsWithLoopback
	Authentication *apiserveroptions.DelegatingAuthenticationOptions
	Authorization  *apiserveroptions.DelegatingAuthorizationOptions
	Features       *apiserveroptions.FeatureOptions
}

func NewOptions() *Options {
	return &Options{
		SecureServing:  newSecureServingOptions(),
		Authentication: apiserveroptions.NewDelegatingAuthenticationOptions(),
		Authorization:  apiserveroptions.NewDelegatingAuthorizationOptions(),
		Features:       apiserveroptions.NewFeatureOptions(),
	}
}

func newSecureServingOptions() *apiserveroptions.SecureServingOptionsWithLoopback {
	secureServing := apiserveroptions.NewSecureServingOptions().WithLoopback()
	secureServing.BindPort = 8443
	secureServing.ServerCert = apiserveroptions.GeneratableKeyCert{
		PairName:      "registry-apiserver",
		CertDirectory: "pki",
	}
	return secureServing
}

func (o *Options) Flags() (fs flag.NamedFlagSets) {
	o.SecureServing.AddFlags(fs.FlagSet("secure serving"))
	o.Authentication.AddFlags(fs.FlagSet("authentication"))
	o.Authorization.AddFlags(fs.FlagSet("authorization"))
	o.Features.AddFlags(fs.FlagSet("features"))
	return fs
}

func (o *Options) ServerConfig() (*server.Config, error) {
	apiserver, err := o.apiserverConfig()
	if err != nil {
		return nil, err
	}

	rest, err := ctrl.GetConfig()
	if err != nil {
		return nil, err
	}

	cnf := &server.Config{
		Apiserver: apiserver,
		Rest:      rest,
	}
	return cnf, nil
}

func (o *Options) apiserverConfig() (*apiserver.Config, error) {
	if err := o.SecureServing.MaybeDefaultWithSelfSignedCerts(
		"localhost", nil, []net.IP{net.ParseIP("127.0.0.1")}); err != nil {
		return nil, err
	}

	cnf := apiserver.NewConfig(api.Codecs)
	if err := o.SecureServing.ApplyTo(&cnf.SecureServing, &cnf.LoopbackClientConfig); err != nil {
		return nil, err
	}
	if err := o.Authentication.ApplyTo(&cnf.Authentication, cnf.SecureServing, nil); err != nil {
		return nil, err
	}
	if err := o.Authorization.ApplyTo(&cnf.Authorization); err != nil {
		return nil, err
	}

	versionGet := version.Get()
	cnf.Version = &versionGet

	return cnf, nil
}
