package api

import (
	"context"
	"fmt"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metainternalversion "k8s.io/apimachinery/pkg/apis/meta/internalversion"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	genericapiserverrequest "k8s.io/apiserver/pkg/endpoints/request"
	"k8s.io/apiserver/pkg/registry/rest"
	listers "k8s.io/client-go/listers/core/v1"

	registryv1alpha1 "github.com/veverita7/registry-server/api/v1alpha1"
)

type imageRegistry struct {
	groupResource schema.GroupResource
	secretLister  listers.SecretLister
}

func newImageRegistry(
	groupResource schema.GroupResource, secretLister listers.SecretLister) *imageRegistry {
	return &imageRegistry{
		groupResource: groupResource,
		secretLister:  secretLister,
	}
}

// KindProvider interface
func (r *imageRegistry) Kind() string {
	return "ImageRegistry"
}

// Scoper interface
func (r *imageRegistry) NamespaceScoped() bool {
	return true
}

// Storage interface
func (r *imageRegistry) New() runtime.Object {
	return &registryv1alpha1.ImageRegistry{}
}

// Lister interface
func (r *imageRegistry) NewList() runtime.Object {
	return &registryv1alpha1.ImageRegistryList{}
}

// Lister interface
func (r *imageRegistry) List(
	ctx context.Context, opts *metainternalversion.ListOptions) (runtime.Object, error) {
	labelSelector := labels.Everything()
	if opts != nil && opts.LabelSelector != nil {
		labelSelector = opts.LabelSelector
	}

	namespace := genericapiserverrequest.NamespaceValue(ctx)
	secrets, err := r.secretLister.Secrets(namespace).List(labelSelector)
	if err != nil {
		return &registryv1alpha1.ImageRegistryList{}, fmt.Errorf("failed listing secrets: %w", err)
	}

	registries, err := r.getImageRegistries(secrets...)
	if err != nil {
		return &registryv1alpha1.ImageRegistryList{},
			fmt.Errorf("failed converting secrets to imageregistries: %w", err)
	}

	return &registryv1alpha1.ImageRegistryList{Items: registries}, nil
}

// Getter interface
func (r *imageRegistry) Get(
	ctx context.Context, name string, opts *metav1.GetOptions) (runtime.Object, error) {
	namespace := genericapiserverrequest.NamespaceValue(ctx)
	secret, err := r.secretLister.Secrets(namespace).Get(name)
	if err != nil {
		if errors.IsNotFound(err) {
			return &registryv1alpha1.ImageRegistry{}, err
		}
		return &registryv1alpha1.ImageRegistry{}, fmt.Errorf("failed getting secret: %w", err)
	}
	if secret == nil {
		return &registryv1alpha1.ImageRegistry{},
			errors.NewNotFound(corev1.Resource("secrets"), fmt.Sprintf("%s/%s", namespace, name))
	}

	registries, err := r.getImageRegistries(secret)
	if err != nil {
		return &registryv1alpha1.ImageRegistry{},
			fmt.Errorf("failed converting secret to imageregistry: %w", err)
	}
	if len(registries) == 0 {
		return nil, errors.NewNotFound(r.groupResource, fmt.Sprintf("%s/%s", namespace, name))
	}

	return &registries[0], nil
}

func (r *imageRegistry) getImageRegistries(
	secrets ...*corev1.Secret) ([]registryv1alpha1.ImageRegistry, error) {
	registries := []registryv1alpha1.ImageRegistry{}
	for _, secret := range secrets {
		registries = append(registries, registryv1alpha1.ImageRegistry{
			ObjectMeta: metav1.ObjectMeta{
				Name:              secret.Name,
				Namespace:         secret.Namespace,
				CreationTimestamp: secret.CreationTimestamp,
				Labels:            secret.Labels,
			},
			Spec: registryv1alpha1.ImageRegistrySpec{
				Repositories: []string{},
			},
		})
	}
	return registries, nil
}

// TableConvertor interface
func (r *imageRegistry) ConvertToTable(
	ctx context.Context, obj runtime.Object, tableOptions runtime.Object) (*metav1.Table, error) {
	return rest.NewDefaultTableConvertor(r.groupResource).ConvertToTable(ctx, obj, tableOptions)
}
