package discovery

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

type K8sDiscoverer struct {
	clientset *kubernetes.Clientset
}

func NewK8sDiscoverer() (*K8sDiscoverer, error) {
	// creates the in-cluster config
	config, err := rest.InClusterConfig()
	if err != nil {
		return nil, err
	}
	// creates the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}
	return &K8sDiscoverer{
		clientset: clientset,
	}, nil

}

func (d *K8sDiscoverer) Discover() []Artifact {
	pods, err := d.clientset.CoreV1().Pods("").List(metav1.ListOptions{})
	if err != nil {
		panic(err.Error())
	}
	podInfos := make([]Artifact, 0, len(pods.Items))
	for _, p := range pods.Items {
		for _, c := range p.Spec.Containers {
			a := Artifact{}
			a.Metadata = p.Labels
			a.Metadata = make(map[string]string)
			a.Name = c.Image
			a.Metadata["pod"] = p.ObjectMeta.Name
			a.Metadata["namespace"] = p.ObjectMeta.Namespace
			podInfos = append(podInfos, a)
		}
	}
	return podInfos
}
