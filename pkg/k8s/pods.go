package k8s

// Pod holds simplified information about a K8s pod.
type Pod struct {
	Name       string      `json:"name"`
	Age        string      `json:"age"`
	Containers []Container `json:"containers"`
}

// Container holds simplified information about a K8s container.
type Container struct {
	Image string `json:"image"`
}
