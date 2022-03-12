package options

import (
	"net"

	genericapiserver "k8s.io/apiserver/pkg/server"
	genericapiserveroptions "k8s.io/apiserver/pkg/server/options"
	"k8s.io/client-go/pkg/version"
	"k8s.io/component-base/cli/flag"
	ctrl "sigs.k8s.io/controller-runtime"

	"github.com/veverita7/registry-server/pkg/api"
	"github.com/veverita7/registry-server/pkg/server"
)

type Options struct {
	SecureServing  *genericapiserveroptions.SecureServingOptionsWithLoopback
	Authentication *genericapiserveroptions.DelegatingAuthenticationOptions
	Authorization  *genericapiserveroptions.DelegatingAuthorizationOptions
	Features       *genericapiserveroptions.FeatureOptions
}

func NewOptions() *Options {
	return &Options{
		SecureServing:  newSecureServingOptions(),
		Authentication: genericapiserveroptions.NewDelegatingAuthenticationOptions(),
		Authorization:  genericapiserveroptions.NewDelegatingAuthorizationOptions(),
		Features:       genericapiserveroptions.NewFeatureOptions(),
	}
}

func newSecureServingOptions() *genericapiserveroptions.SecureServingOptionsWithLoopback {
	secureServing := genericapiserveroptions.NewSecureServingOptions().WithLoopback()
	secureServing.BindPort = 8443
	secureServing.ServerCert = genericapiserveroptions.GeneratableKeyCert{
		CertDirectory: "pki",
		PairName:      "registry-apiserver",
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

func (o *Options) apiserverConfig() (*genericapiserver.Config, error) {
	if err := o.SecureServing.MaybeDefaultWithSelfSignedCerts(
		"localhost", nil, []net.IP{net.ParseIP("127.0.0.1")}); err != nil {
		return nil, err
	}

	cnf := genericapiserver.NewConfig(api.Codecs)
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
