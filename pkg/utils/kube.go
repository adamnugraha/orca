package utils

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

// CreateNamespace creates a namespace
func CreateNamespace(name, kubeContext string, print bool) error {
	clientset, err := getClientSet(kubeContext)
	if err != nil {
		return err
	}

	nsSpec := &v1.Namespace{ObjectMeta: metav1.ObjectMeta{
		Name: name,
	}}
	_, err = clientset.Core().Namespaces().Create(nsSpec)
	if err != nil {
		return err
	}
	if print {
		log.Printf("created namespace \"%s\"", name)
	}
	return nil
}

// GetNamespace gets a namespace
func GetNamespace(name, kubeContext string) (*v1.Namespace, error) {
	clientset, err := getClientSet(kubeContext)
	if err != nil {
		return nil, err
	}
	getOptions := metav1.GetOptions{}
	nsSpec, err := clientset.Core().Namespaces().Get(name, getOptions)
	if err != nil {
		return nil, err
	}
	return nsSpec, nil
}

// UpdateNamespace updates a namespace
func UpdateNamespace(name, kubeContext string, annotations map[string]string, print bool) error {
	clientset, err := getClientSet(kubeContext)
	if err != nil {
		return err
	}

	nsSpec := &v1.Namespace{ObjectMeta: metav1.ObjectMeta{
		Name:        name,
		Annotations: annotations,
	}}
	_, err = clientset.Core().Namespaces().Update(nsSpec)
	if err != nil {
		return err
	}
	if print {
		log.Printf("updated namespace \"%s\" with annotations (%s)", name, annotations)
	}
	return nil
}

// DeleteNamespace deletes a namespace
func DeleteNamespace(name, kubeContext string, print bool) error {
	clientset, err := getClientSet(kubeContext)
	if err != nil {
		return err
	}
	deleteOptions := &metav1.DeleteOptions{}
	err = clientset.Core().Namespaces().Delete(name, deleteOptions)
	if err != nil {
		return err
	}
	if print {
		log.Printf("deleted namespace \"%s\"", name)
	}
	return nil
}

// NamespaceExists returns true if the namespace exists
func NamespaceExists(name, kubeContext string) (bool, error) {
	clientset, err := getClientSet(kubeContext)
	if err != nil {
		return false, err
	}

	listOptions := metav1.ListOptions{}
	namespaces, err := clientset.Core().Namespaces().List(listOptions)
	if err != nil {
		return false, err
	}
	for _, ns := range namespaces.Items {
		if ns.Name == name {
			nsStatus := ns.Status.Phase
			if nsStatus != "Active" {
				return false, fmt.Errorf("Environment status is %s", nsStatus)
			}
			return true, nil
		}
	}
	return false, nil
}

func buildConfigFromFlags(context, kubeconfigPath string) (*rest.Config, error) {
	return clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
		&clientcmd.ClientConfigLoadingRules{ExplicitPath: kubeconfigPath},
		&clientcmd.ConfigOverrides{
			CurrentContext: context,
		}).ClientConfig()
}

func getClientSet(kubeContext string) (*kubernetes.Clientset, error) {
	var kubeconfig string
	if kubeConfigPath := os.Getenv("KUBECONFIG"); kubeConfigPath != "" {
		kubeconfig = kubeConfigPath
	} else {
		kubeconfig = filepath.Join(os.Getenv("HOME"), ".kube", "config")
	}

	config, err := buildConfigFromFlags(kubeContext, kubeconfig)
	if err != nil {
		return nil, err
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	return clientset, nil
}
