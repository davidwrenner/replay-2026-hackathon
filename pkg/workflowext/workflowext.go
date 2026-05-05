package workflowext

import (
	"context"

	"go.temporal.io/sdk/workflow"
)

// NoOpActivity is a traceable activity that does nothing.
// It is exported so it can be registered with the Temporal worker.
func NoOpActivity(_ context.Context, _ *any) error {
	return nil
}

// ExecuteOptional wraps workflow.ExecuteActivity.
// If args is nil, it executes a NoOpActivity to maintain a visible execution history.
// If args is non-nil, it proceeds with the provided activity and arguments.
func ExecuteOptional(ctx workflow.Context, activity any, args *any) workflow.Future {
	if args == nil {
		return workflow.ExecuteActivity(ctx, NoOpActivity)
	}

	return workflow.ExecuteActivity(ctx, activity, args)
}
