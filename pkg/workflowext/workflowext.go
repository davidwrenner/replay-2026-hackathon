package workflowext

import (
	"context"
	"reflect"

	"go.temporal.io/sdk/workflow"
)

// NoOpActivity is a traceable activity that does nothing.
// It is exported so it can be registered with the Temporal worker.
func NoOpActivity(_ context.Context) error {
	return nil
}

type noOpArgs struct {
	SkippedActivity string
}

// ExecuteOptional wraps workflow.ExecuteActivity.
// If input is nil, it executes a NoOpActivity to maintain a visible execution history.
// If input is non-nil, it proceeds with the provided activity and input.
func ExecuteOptional[S, T any](ctx workflow.Context, activity S, input *T) workflow.Future {
	if input == nil {
		skipped := "unknown"
		if reflect.TypeOf(activity).Kind() == reflect.String {
			skipped = any(activity).(string)
		}

		return workflow.ExecuteActivity(ctx, NoOpActivity, noOpArgs{SkippedActivity: skipped})
	}

	return workflow.ExecuteActivity(ctx, activity, input)
}
