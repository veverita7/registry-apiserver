package server

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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

func secretInformerFactory(rest *rest.Config) (informers.SharedInformerFactory, error) {
	client, err := kubernetes.NewForConfig(rest)
	if err != nil {
		return nil, err
	}
	return informers.NewFilteredSharedInformerFactory(
		client, defaultResync, corev1.NamespaceAll, func(opts *metav1.ListOptions) {
			opts.FieldSelector = "type=kubernetes.io/dockerconfigjson"
		}), nil
}
