package command

import (
	"testing"

	"github.com/docker/cli/cli/context/store"
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
				Description:       "description",
				StackOrchestrator: OrchestratorKubernetes,
			},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			storeMeta := store.ContextMetadata{}
			SetContextMetadata(&storeMeta, c.data)
			result, err := GetContextMetadata(storeMeta)
			assert.NilError(t, err)
			assert.Equal(t, result.Description, c.data.Description)
			if c.data.StackOrchestrator == Orchestrator("") { // ensure orchestrator is normalized
				assert.Equal(t, result.StackOrchestrator, orchestratorUnset)
			} else {
				assert.Equal(t, result.StackOrchestrator, c.data.StackOrchestrator)
			}
		})
	}
}
