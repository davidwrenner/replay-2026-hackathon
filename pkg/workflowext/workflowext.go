package workflowext

import (
	"context"

	"go.temporal.io/sdk/workflow"
)

// NoOpActivity is a traceable activity that does nothing.
// It is exported so it can be registered with the Temporal worker.
func NoOpActivity(_ context.Context) error {
	return nil
}

// ExecuteOptional wraps workflow.ExecuteActivity.
// If input is nil, it executes a NoOpActivity to maintain a visible execution history.
// If input is non-nil, it proceeds with the provided activity and input.
func ExecuteOptional[S, T any](ctx workflow.Context, activity S, input *T) workflow.Future {
	if input == nil {
		return workflow.ExecuteActivity(ctx, NoOpActivity)
	}

	return workflow.ExecuteActivity(ctx, activity, input)
}
