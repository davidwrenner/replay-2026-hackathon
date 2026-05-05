package workflowext

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go.temporal.io/sdk/activity"
	"go.temporal.io/sdk/testsuite"
	"go.temporal.io/sdk/workflow"
)

//func TestNoOpActivity(t *testing.T) {
//	cases := []struct {
//		name string
//		arg  *any
//	}{
//		{"nil argument", nil},
//		{"non-nil argument", new(any("payload"))},
//	}
//
//	for _, tc := range cases {
//		tc := tc
//		t.Run(tc.name, func(t *testing.T) {
//			ts := &testsuite.WorkflowTestSuite{}
//			env := ts.NewTestActivityEnvironment()
//			env.RegisterActivity(NoOpActivity)
//
//			_, err := env.ExecuteActivity(NoOpActivity, tc.arg)
//			require.NoError(t, err)
//		})
//	}
//}

type testStruct struct {
	Field1 string
	Field2 int
}

func TestExecuteOptional(t *testing.T) {
	const realActivityName = "TestRealActivity"

	cases := []struct {
		name            string
		args            *testStruct
		expectNoOpCalls int
		expectRealCalls int
	}{
		{"nil args dispatches NoOpActivity", nil, 1, 0},
		{"non-nil args dispatches provided activity", &testStruct{"123", 2}, 0, 1},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			var noOpCalls, realCalls int

			noOp := func(ctx context.Context, args *testStruct) error {
				noOpCalls++
				return NoOpActivity(ctx)
			}
			realFn := func(_ context.Context, _ *testStruct) error {
				realCalls++
				return nil
			}

			ts := &testsuite.WorkflowTestSuite{}
			env := ts.NewTestWorkflowEnvironment()
			env.RegisterActivityWithOptions(noOp, activity.RegisterOptions{Name: "NoOpActivity"})
			env.RegisterActivityWithOptions(realFn, activity.RegisterOptions{Name: realActivityName})

			wrapper := func(ctx workflow.Context, args *testStruct) error {
				ctx = workflow.WithActivityOptions(ctx, workflow.ActivityOptions{
					StartToCloseTimeout: time.Minute,
				})
				return ExecuteOptional(ctx, realActivityName, args).Get(ctx, nil)
			}
			env.RegisterWorkflow(wrapper)

			env.ExecuteWorkflow(wrapper, tc.args)

			require.True(t, env.IsWorkflowCompleted())
			require.NoError(t, env.GetWorkflowError())
			require.Equal(t, tc.expectNoOpCalls, noOpCalls)
			require.Equal(t, tc.expectRealCalls, realCalls)
		})
	}
}
