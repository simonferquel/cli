package command

import (
	"github.com/docker/docker/pkg/contextstore"
)

// ContextMetadata is a typed representation of what we put in Context metadata
type ContextMetadata struct {
	Description             string
	StackOrchestrator       Orchestrator
	HelperEnabledDockerHost string
}

// SetContextMetadata set the metadata inside a stored context
func SetContextMetadata(ctx *contextstore.ContextMetadata, metadata ContextMetadata) {
	ctx.Metadata = map[string]interface{}{
		"description":              metadata.Description,
		"defaultStackOrchestrator": string(metadata.StackOrchestrator),
		"helperEnabledDockerHost":  metadata.HelperEnabledDockerHost,
	}
}

// GetContextMetadata extract metadata from stored context metadata
func GetContextMetadata(ctx contextstore.ContextMetadata) (ContextMetadata, error) {
	if ctx.Metadata == nil {
		return ContextMetadata{}, nil
	}
	result := ContextMetadata{}
	var err error
	if val, ok := ctx.Metadata["description"]; ok {
		result.Description, _ = val.(string)
	}
	if val, ok := ctx.Metadata["defaultStackOrchestrator"]; ok {
		v, _ := val.(string)
		if result.StackOrchestrator, err = normalize(v); err != nil {
			return result, err
		}
	}
	if val, ok := ctx.Metadata["helperEnabledDockerHost"]; ok {
		result.HelperEnabledDockerHost, _ = val.(string)

	}
	return result, nil
}
