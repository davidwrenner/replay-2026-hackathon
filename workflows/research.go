package workflows

import (
	"time"

	"github.com/davidwrenner/replay-2026-hackathon/activities"
	"go.temporal.io/sdk/workflow"
)

// ResearchWorkflowInput is the input for the research workflow.
type ResearchWorkflowInput struct {
	Query string `json:"query"`
}

// ResearchWorkflowOutput is the output from the research workflow.
type ResearchWorkflowOutput struct {
	Research string `json:"research"`
}

// ResearchWorkflow orchestrates the research process.
func ResearchWorkflow(ctx workflow.Context, input ResearchWorkflowInput) (*ResearchWorkflowOutput, error) {
	ao := workflow.ActivityOptions{
		StartToCloseTimeout: 10 * time.Second,
	}
	ctx = workflow.WithActivityOptions(ctx, ao)

	var result activities.ResearchOutput
	err := workflow.ExecuteActivity(ctx, "Research", activities.ResearchInput{Query: input.Query}).Get(ctx, &result)
	if err != nil {
		return nil, err
	}

	return &ResearchWorkflowOutput{
		Research: result.Research,
	}, nil
}
