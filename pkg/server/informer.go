package server

import (
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

const defaultResync = 0

func informerFactory(rest *rest.Config) (informers.SharedInformerFactory, error) {
	client, err := kubernetes.NewForConfig(rest)
	if err != nil {
		return nil, err
	}
	return informers.NewSharedInformerFactory(client, defaultResync), nil
}
