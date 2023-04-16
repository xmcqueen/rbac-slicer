package kubevar

import (
	"fmt"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	hd "k8s.io/client-go/util/homedir"
)

type Kubeconfig struct {
	Filepath  string
	Config    *rest.Config
	Clientset *kubernetes.Clientset
}

func Default() string {
	return fmt.Sprintf("%s/.kube/config", hd.HomeDir())
}

func (k *Kubeconfig) Set(v string) error {
	k.Filepath = v
	config, err := clientcmd.BuildConfigFromFlags("", k.Filepath)
	if err != nil {
		return err
	}
	k.Config = config

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return err
	}
	k.Clientset = clientset

	return nil
}
func (k *Kubeconfig) String() string {
	return k.Filepath
}
func (k *Kubeconfig) Get() *kubernetes.Clientset {
	return k.Clientset
}
