package workflows

import (
	"time"

	"github.com/davidwrenner/replay-2026-hackathon/activities"
	"github.com/davidwrenner/replay-2026-hackathon/pkg/workflowext"
	"go.temporal.io/sdk/workflow"
)

// DataSource represents a data source from the request.
type DataSource struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	APIKey      string `json:"apiKey"`
}

// ResearchWorkflowInput is the input for the research workflow.
type ResearchWorkflowInput struct {
	DataSources   []DataSource `json:"dataSources"`
	Prompt        string       `json:"prompt"`
	RiskTolerance string       `json:"riskTolerance"`
	MaxBudget     string       `json:"maxBudget"`
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

	// Step 1: Check for prompt injection
	err := workflow.ExecuteActivity(ctx, "CheckInjection", activities.InjectionCheckInput{
		Prompt:        input.Prompt,
		RiskTolerance: input.RiskTolerance,
	}).Get(ctx, nil)
	if err != nil {
		return nil, err
	}

	// Step 2: Perform research
	jobConfig := &activities.ResearchInput{Query: input.Prompt}

	var result activities.ResearchOutput
	err = workflowext.ExecuteOptional(ctx, "Research", jobConfig).Get(ctx, &result)
	if err != nil {
		return nil, err
	}

	return &ResearchWorkflowOutput{
		Research: result.Research,
	}, nil
}
