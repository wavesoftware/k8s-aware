package k8s

import (
	"context"
	"errors"
	"fmt"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

// ErrNoKubernetesConnection if can't connect to Kube API server.
var ErrNoKubernetesConnection = errors.New("no Kubernetes connection")

// Info can retrieve information pods running in current namespace.
type Info interface {
	Pods() ([]Pod, error)
}

// Client implements Info.
type Client struct {
	restcfg *rest.Config
	ns      string
}

func NewClient() (*Client, error) {
	loadingRules := clientcmd.NewDefaultClientConfigLoadingRules()
	var configOverrides *clientcmd.ConfigOverrides
	cc := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(loadingRules, configOverrides)
	restcfg, err := cc.ClientConfig()
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrNoKubernetesConnection, err)
	}
	ns, _, err := cc.Namespace()
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrNoKubernetesConnection, err)
	}
	return &Client{restcfg: restcfg, ns: ns}, nil
}

func (c Client) Pods() ([]Pod, error) {
	t, err := c.typed()
	if err != nil {
		return nil, err
	}
	podList, err := t.CoreV1().Pods(c.ns).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrNoKubernetesConnection, err)
	}
	pods := make([]Pod, 0, len(podList.Items))
	for _, pod := range podList.Items {
		my := Pod{
			Name: pod.GetName(),
			Age:  calcAge(pod.GetCreationTimestamp()),
		}
		for _, container := range pod.Spec.Containers {
			my.Containers = append(my.Containers, Container{
				Image: container.Image,
			})
		}
		pods = append(pods, my)
	}
	return pods, nil
}

func calcAge(timestamp metav1.Time) string {
	return fmt.Sprintf("%v", time.Since(timestamp.Time))
}

func (c Client) typed() (kubernetes.Interface, error) {
	typed, err := kubernetes.NewForConfig(c.restcfg)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrNoKubernetesConnection, err)
	}
	return typed, nil
}
