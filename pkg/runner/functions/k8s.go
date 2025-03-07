package functions

import (
	"context"

	"github.com/kyverno/chainsaw/pkg/client"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	ctrlclient "sigs.k8s.io/controller-runtime/pkg/client"
)

func jpKubernetesResourceExists(arguments []any) (any, error) {
	var client client.Client
	var apiVersion, kind string
	if err := getArg(arguments, 0, &client); err != nil {
		return false, err
	}
	if err := getArg(arguments, 1, &apiVersion); err != nil {
		return false, err
	}
	if err := getArg(arguments, 2, &kind); err != nil {
		return false, err
	}

	mapper := client.RESTMapper()

	gvk := schema.GroupVersionKind{Group: "", Version: apiVersion, Kind: kind}
	mapping, err := mapper.RESTMapping(gvk.GroupKind(), gvk.Version)
	if err != nil {
		if apierrors.IsNotFound(err) {
			// If the error is due to the resource not being found, return false without an error.
			return false, nil
		}
		// For any other error, return it.
		return false, err
	}
	// If a mapping for the resource is found, it means the resource exists.
	return mapping != nil, nil
}

func jpKubernetesExists(arguments []any) (any, error) {
	var client client.Client
	var apiVersion, kind string
	var key ctrlclient.ObjectKey
	if err := getArg(arguments, 0, &client); err != nil {
		return false, err
	}
	if err := getArg(arguments, 1, &apiVersion); err != nil {
		return false, err
	}
	if err := getArg(arguments, 2, &kind); err != nil {
		return false, err
	}
	if err := getArg(arguments, 3, &key.Namespace); err != nil {
		return false, err
	}
	if err := getArg(arguments, 4, &key.Name); err != nil {
		return false, err
	}

	err := client.Get(context.TODO(), key, &unstructured.Unstructured{})
	if err == nil {
		return true, nil
	}
	if apierrors.IsNotFound(err) {
		return false, nil
	}
	return false, err
}

func jpKubernetesGet(arguments []any) (any, error) {
	var client client.Client
	var apiVersion, kind string
	var key ctrlclient.ObjectKey
	if err := getArg(arguments, 0, &client); err != nil {
		return nil, err
	}
	if err := getArg(arguments, 1, &apiVersion); err != nil {
		return nil, err
	}
	if err := getArg(arguments, 2, &kind); err != nil {
		return nil, err
	}
	if err := getArg(arguments, 3, &key.Namespace); err != nil {
		return nil, err
	}
	if err := getArg(arguments, 4, &key.Name); err != nil {
		return nil, err
	}
	var obj unstructured.Unstructured
	obj.SetAPIVersion(apiVersion)
	obj.SetKind(kind)
	if err := client.Get(context.TODO(), key, &obj); err != nil {
		return nil, err
	}
	return obj.UnstructuredContent(), nil
}

func jpKubernetesList(arguments []any) (any, error) {
	var client client.Client
	var apiVersion, kind, namespace string
	if err := getArg(arguments, 0, &client); err != nil {
		return nil, err
	}
	if err := getArg(arguments, 1, &apiVersion); err != nil {
		return nil, err
	}
	if err := getArg(arguments, 2, &kind); err != nil {
		return nil, err
	}
	if len(arguments) == 4 {
		if err := getArg(arguments, 3, &namespace); err != nil {
			return nil, err
		}
	}
	var list unstructured.UnstructuredList
	list.SetAPIVersion(apiVersion)
	list.SetKind(kind)
	var listOptions []ctrlclient.ListOption
	if namespace != "" {
		listOptions = append(listOptions, ctrlclient.InNamespace(namespace))
	}
	if err := client.List(context.TODO(), &list, listOptions...); err != nil {
		return nil, err
	}
	return list.Items, nil
}
