package workflows

import (
	"time"

	"github.com/davidwrenner/replay-2026-hackathon/activities"
	"github.com/davidwrenner/replay-2026-hackathon/pkg/workflowext"
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

	jobConfig := &activities.ResearchInput{Query: input.Query}
	//jobConfig = nil

	var result activities.ResearchOutput
	err := workflowext.ExecuteOptional(ctx, "Research", jobConfig).Get(ctx, &result)
	if err != nil {
		return nil, err
	}

	return &ResearchWorkflowOutput{
		Research: result.Research,
	}, nil
}
