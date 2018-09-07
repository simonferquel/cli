package command

import (
	"fmt"
	"io"
	"os"

	cliconfig "github.com/docker/cli/cli/config"
)

// Orchestrator type acts as an enum describing supported orchestrators.
type Orchestrator string

const (
	// OrchestratorKubernetes orchestrator
	OrchestratorKubernetes = Orchestrator("kubernetes")
	// OrchestratorSwarm orchestrator
	OrchestratorSwarm = Orchestrator("swarm")
	// OrchestratorAll orchestrator
	OrchestratorAll   = Orchestrator("all")
	orchestratorUnset = Orchestrator("unset")

	defaultOrchestrator           = OrchestratorSwarm
	envVarDockerStackOrchestrator = "DOCKER_STACK_ORCHESTRATOR"
	envVarDockerOrchestrator      = "DOCKER_ORCHESTRATOR"
)

// HasKubernetes returns true if defined orchestrator has Kubernetes capabilities.
func (o Orchestrator) HasKubernetes() bool {
	return o == OrchestratorKubernetes || o == OrchestratorAll
}

// HasSwarm returns true if defined orchestrator has Swarm capabilities.
func (o Orchestrator) HasSwarm() bool {
	return o == OrchestratorSwarm || o == OrchestratorAll
}

// HasAll returns true if defined orchestrator has both Swarm and Kubernetes capabilities.
func (o Orchestrator) HasAll() bool {
	return o == OrchestratorAll
}

func normalize(value string) (Orchestrator, error) {
	switch value {
	case "kubernetes":
		return OrchestratorKubernetes, nil
	case "swarm":
		return OrchestratorSwarm, nil
	case "", "unset":
		return orchestratorUnset, nil
	case "all":
		return OrchestratorAll, nil
	default:
		return defaultOrchestrator, fmt.Errorf("specified orchestrator %q is invalid, please use either kubernetes, swarm or all", value)
	}
}

// NormalizeOrchestrator parses an orchestrator value and check if it is valid
func NormalizeOrchestrator(value string) (Orchestrator, error) {
	return normalize(value)
}

// GetStackOrchestrator checks DOCKER_STACK_ORCHESTRATOR environment variable and configuration file
// orchestrator value and returns user defined Orchestrator.
func GetStackOrchestrator(flagValue string, dockerCli Cli, stderr io.Writer) (Orchestrator, error) {
	// Check flag
	if o, err := normalize(flagValue); o != orchestratorUnset {
		return o, err
	}
	// Check environment variable
	env := os.Getenv(envVarDockerStackOrchestrator)
	if env == "" && os.Getenv(envVarDockerOrchestrator) != "" {
		fmt.Fprintf(stderr, "WARNING: experimental environment variable %s is set. Please use %s instead\n", envVarDockerOrchestrator, envVarDockerStackOrchestrator)
	}
	if o, err := normalize(env); o != orchestratorUnset {
		return o, err
	}
	// if there is a current context, check context orchestrator
	ctxName := dockerCli.CurrentContext()
	if ctxName != "" && ctxName != ContextDockerHost {
		rawMeta, err := dockerCli.ContextStore().GetContextMetadata(ctxName)
		if err != nil {
			return defaultOrchestrator, err
		}
		meta, err := GetContextMetadata(rawMeta)
		if err != nil {
			return defaultOrchestrator, err
		}
		if meta.StackOrchestrator != orchestratorUnset {
			return meta.StackOrchestrator, nil
		}
	}
	// Check specified orchestrator in global config file
	configFile := dockerCli.ConfigFile()
	if configFile == nil {
		configFile = cliconfig.LoadDefaultConfigFile(dockerCli.Err())
	}
	if o, err := normalize(configFile.StackOrchestrator); o != orchestratorUnset {
		return o, err
	}
	// Nothing set, use default orchestrator
	return defaultOrchestrator, nil
}
