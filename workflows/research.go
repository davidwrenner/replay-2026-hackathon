package workflows

import (
	"time"

	"github.com/davidwrenner/replay-2026-hackathon/activities"
	"github.com/davidwrenner/replay-2026-hackathon/pkg/workflowext"
	"go.temporal.io/sdk/workflow"
)

// DataSource represents an external data provider
type DataSource struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	APIKey      string `json:"apiKey"`
}

// ResearchWorkflowInput is the input for the research workflow.
type ResearchWorkflowInput struct {
	Prompt        string       `json:"prompt"`
	RiskTolerance string       `json:"riskTolerance"`
	MaxBudget     string       `json:"maxBudget"`
	DataSources   []DataSource `json:"dataSources"`
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

	jobConfig := &activities.ResearchInput{
		Prompt:        input.Prompt,
		RiskTolerance: input.RiskTolerance,
		MaxBudget:     input.MaxBudget,
		DataSources:   convertDataSources(input.DataSources),
	}
	// jobConfig = nil

	// Step 2: Fetch data from sources
	err = workflowext.ExecuteOptional(ctx, "Bloomberg", jobConfig).Get(ctx, nil)
	if err != nil {
		return nil, err
	}

	err = workflowext.ExecuteOptional(ctx, "DowJones", jobConfig).Get(ctx, nil)
	if err != nil {
		return nil, err
	}

	err = workflowext.ExecuteOptional(ctx, "LexisNexis", jobConfig).Get(ctx, nil)
	if err != nil {
		return nil, err
	}

	err = workflowext.ExecuteOptional(ctx, "NYTimes", jobConfig).Get(ctx, nil)
	if err != nil {
		return nil, err
	}

	err = workflowext.ExecuteOptional(ctx, "Polymarket", jobConfig).Get(ctx, nil)
	if err != nil {
		return nil, err
	}

	err = workflowext.ExecuteOptional(ctx, "Reddit", jobConfig).Get(ctx, nil)
	if err != nil {
		return nil, err
	}

	err = workflowext.ExecuteOptional(ctx, "Refinitiv", jobConfig).Get(ctx, nil)
	if err != nil {
		return nil, err
	}

	err = workflowext.ExecuteOptional(ctx, "Twitter", jobConfig).Get(ctx, nil)
	if err != nil {
		return nil, err
	}

	err = workflowext.ExecuteOptional(ctx, "WallStreetJournal", jobConfig).Get(ctx, nil)
	if err != nil {
		return nil, err
	}

	err = workflowext.ExecuteOptional(ctx, "YouTube", jobConfig).Get(ctx, nil)
	if err != nil {
		return nil, err
	}

	// Step 3: Perform research
	var result activities.ResearchOutput
	err = workflowext.ExecuteOptional(ctx, "Research", jobConfig).Get(ctx, &result)
	if err != nil {
		return nil, err
	}

	return &ResearchWorkflowOutput{
		Research: result.Research,
	}, nil
}

func convertDataSources(ds []DataSource) []activities.DataSource {
	res := make([]activities.DataSource, len(ds))
	for i, d := range ds {
		res[i] = activities.DataSource{Name: d.Name, Description: d.Description, APIKey: d.APIKey}
	}
	return res
}
