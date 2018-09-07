package command

import (
	"testing"

	"github.com/docker/docker/pkg/contextstore"
	"gotest.tools/assert"
)

func TestContextStorageParsingAndSaving(t *testing.T) {
	cases := []struct {
		name string
		data ContextMetadata
	}{
		{
			name: "empty",
			data: ContextMetadata{},
		},
		{
			name: "with-values",
			data: ContextMetadata{
				Description:             "description",
				HelperEnabledDockerHost: "test-host",
				StackOrchestrator:       OrchestratorKubernetes,
			},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			storeMeta := contextstore.ContextMetadata{}
			SetContextMetadata(&storeMeta, c.data)
			result, err := GetContextMetadata(storeMeta)
			assert.Check(t, err)
			assert.Equal(t, result.Description, c.data.Description)
			assert.Equal(t, result.HelperEnabledDockerHost, c.data.HelperEnabledDockerHost)
			if c.data.StackOrchestrator == Orchestrator("") { // ensure orchestrator is normalized
				assert.Equal(t, result.StackOrchestrator, orchestratorUnset)
			} else {
				assert.Equal(t, result.StackOrchestrator, c.data.StackOrchestrator)
			}
		})
	}
}
