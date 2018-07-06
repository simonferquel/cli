package command

import (
	"github.com/docker/context-store"
)

// ContextMetadata is a typed representation of what we put in Context metadata
type ContextMetadata struct {
	Description       string
	Orchestrator      Orchestrator
	StackOrchestrator Orchestrator
}

// SetContextMetadata set the metadata inside a stored context
func SetContextMetadata(ctx *store.ContextMetadata, metadata ContextMetadata) {
	ctx.Metadata = map[string]interface{}{
		"description":              metadata.Description,
		"defaultOrchestrator":      string(metadata.Orchestrator),
		"defaultStackOrchestrator": string(metadata.StackOrchestrator),
	}
}

// GetContextMetadata extract metadata from stored context metadata
func GetContextMetadata(ctx store.ContextMetadata) (ContextMetadata, error) {
	if ctx.Metadata == nil {
		return ContextMetadata{}, nil
	}
	result := ContextMetadata{}
	var err error
	if val, ok := ctx.Metadata["description"]; ok {
		result.Description, _ = val.(string)
	}
	if val, ok := ctx.Metadata["defaultOrchestrator"]; ok {
		v, _ := val.(string)
		if result.Orchestrator, err = NormalizeOrchestrator(v); err != nil {
			return result, err
		}
	}
	if val, ok := ctx.Metadata["defaultStackOrchestrator"]; ok {
		v, _ := val.(string)
		if result.StackOrchestrator, err = NormalizeOrchestrator(v); err != nil {
			return result, err
		}
	}
	return result, nil
}
