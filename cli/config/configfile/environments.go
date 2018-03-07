package configfile

import "github.com/docker/go-connections/tlsconfig"

// Orchestrator defines an env orchestrator
type Orchestrator string

const (
	// OrchestratorAuto specify auto-deduction rule
	OrchestratorAuto = Orchestrator("")

	// OrchestratorSwarm orchestrator
	OrchestratorSwarm = Orchestrator("swarm")

	// OrchestratorKubernetes orchestrator
	OrchestratorKubernetes = Orchestrator("kubernetes")
)

// Environment describes a connection to a Docker and/or Kubernetes server
type Environment struct {
	Docker              *DockerEnvironment     `json:"docker,omitempty"`
	Kubernetes          *KubernetesEnvironment `json:"kubernetes,omitempty"`
	DefaultOrchestrator Orchestrator           `json:"defaultOrchestrator,omitempty"`
}

//ResolveOrchestrator resolves the orchestrator value for the environment
func (e *Environment) ResolveOrchestrator() string {
	if e.DefaultOrchestrator != OrchestratorAuto {
		return string(e.DefaultOrchestrator)
	}
	if e.Docker != nil {
		return string(OrchestratorSwarm)
	}
	return string(OrchestratorKubernetes)
}

// DockerEnvironment describes a connection to a Docker server
type DockerEnvironment struct {
	Host            string `json:"host,omitempty"`
	CaData          []byte `json:"caData,omitempty"`
	Ca              string `json:"ca,omitempty"`
	CertData        []byte `json:"certData,omitempty"`
	Cert            string `json:"cert,omitempty"`
	KeyData         []byte `json:"keyData,omitempty"`
	Key             string `json:"key,omitempty"`
	SkipTLSCAVerify bool   `json:"skipTLSCAVerify,omitempty"`
	TLSEnabled      bool   `json:"tlsEnabled,omitempty"`
}

// TLSOptions returns TLS options for this environment
func (e *DockerEnvironment) TLSOptions() *tlsconfig.Options {
	if !e.TLSEnabled {
		return nil
	}
	return &tlsconfig.Options{
		CAPEM:              e.CaData,
		CertPEM:            e.CertData,
		KeyPEM:             e.KeyData,
		InsecureSkipVerify: e.SkipTLSCAVerify,
	}
}

// KubernetesEnvironment describes a connection to a Kubernetes cluster
// if Server is set, kubeconfig/kubeconfigContext are ignored and the k8s config is read directly from this file
type KubernetesEnvironment struct {
	KubeconfigFile    string `json:"kubeconfigFile,omitempty"`
	KubeconfigContext string `json:"kubeconfigContext,omitempty"`
	Server            string `json:"server,omitempty"`
	SkipTLSVerify     bool   `json:"skipTLSVerify,omitempty"`
	CaData            []byte `json:"caData,omitempty"`
	Ca                string `json:"ca,omitempty"`
	CertData          []byte `json:"certData,omitempty"`
	Cert              string `json:"cert,omitempty"`
	KeyData           []byte `json:"keyData,omitempty"`
	Key               string `json:"key,omitempty"`
}
