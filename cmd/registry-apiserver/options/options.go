package options

import (
	serveroptions "k8s.io/apiserver/pkg/server/options"
	"k8s.io/component-base/cli/flag"
)

type Options struct {
	SecureServing  *serveroptions.SecureServingOptionsWithLoopback
	Authentication *serveroptions.DelegatingAuthenticationOptions
	Authorization  *serveroptions.DelegatingAuthorizationOptions
	Features       *serveroptions.FeatureOptions
	// Only to be used to for testing
	DisableAuthForTesting bool
}

func NewOptions() *Options {
	return &Options{
		SecureServing:  serveroptions.NewSecureServingOptions().WithLoopback(),
		Authentication: serveroptions.NewDelegatingAuthenticationOptions(),
		Authorization:  serveroptions.NewDelegatingAuthorizationOptions(),
		Features:       serveroptions.NewFeatureOptions(),
	}
}

func (o *Options) Flags() (fs flag.NamedFlagSets) {
	o.SecureServing.AddFlags(fs.FlagSet("apiserver secure serving"))
	o.Authentication.AddFlags(fs.FlagSet("apiserver authentication"))
	o.Authorization.AddFlags(fs.FlagSet("apiserver authorization"))
	o.Features.AddFlags(fs.FlagSet("apiserver features"))
	return fs
}
