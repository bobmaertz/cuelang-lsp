// Kubernetes resource definitions example
package kubernetes

// Pod defines a Kubernetes Pod
#Pod: {
	apiVersion: "v1"
	kind:       "Pod"
	metadata:   #ObjectMeta
	spec:       #PodSpec
}

// ObjectMeta defines common metadata
#ObjectMeta: {
	name:      string
	namespace: string | *"default"
	labels?: [string]: string
	annotations?: [string]: string
}

// PodSpec defines the Pod specification
#PodSpec: {
	containers: [...#Container]
	restartPolicy: *"Always" | "OnFailure" | "Never"
	volumes?: [...#Volume]
}

// Container defines a container
#Container: {
	name:  string
	image: string
	ports?: [...#ContainerPort]
	env?: [...#EnvVar]
	resources?: #ResourceRequirements
}

// ContainerPort defines a container port
#ContainerPort: {
	containerPort: int & >=1 & <=65535
	protocol:      *"TCP" | "UDP"
	name?:         string
}

// EnvVar defines an environment variable
#EnvVar: {
	name:  string
	value: string
}

// ResourceRequirements defines resource limits
#ResourceRequirements: {
	limits?: {
		cpu?:    string
		memory?: string
	}
	requests?: {
		cpu?:    string
		memory?: string
	}
}

// Volume defines a volume
#Volume: {
	name: string
	configMap?: {
		name: string
	}
	secret?: {
		secretName: string
	}
}

// Example Pod instance
examplePod: #Pod & {
	metadata: {
		name:      "nginx-pod"
		namespace: "production"
		labels: {
			app:     "nginx"
			version: "1.0"
		}
	}
	spec: {
		containers: [{
			name:  "nginx"
			image: "nginx:1.21"
			ports: [{
				containerPort: 80
				protocol:      "TCP"
			}]
			resources: {
				limits: {
					cpu:    "500m"
					memory: "512Mi"
				}
				requests: {
					cpu:    "250m"
					memory: "256Mi"
				}
			}
		}]
	}
}
