package workflows

import (
	"time"

	"github.com/davidwrenner/replay-2026-hackathon/activities"
	"github.com/davidwrenner/replay-2026-hackathon/pkg/workflowext"
	"go.temporal.io/sdk/workflow"
)

// AllDataSources is the list of all available data source activity names

// ResearchWorkflowInputV2 is a cleaner input structure
type ResearchWorkflowInputV2 struct {
	Prompt         string                    `json:"prompt"`
	RiskTolerance  string                    `json:"riskTolerance"`
	MaxBudget      string                    `json:"maxBudget"`
	EnabledSources []string                  `json:"enabledSources"`
	ResearchConfig *activities.ResearchInput `json:"researchConfig"`
}

// ResearchWorkflowV2 is a cleaner version using a for loop
func ResearchWorkflowV2(ctx workflow.Context, input ResearchWorkflowInputV2) (*ResearchWorkflowOutput, error) {
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

	// Build enabled set
	enabled := make(map[string]bool)
	for _, s := range input.EnabledSources {
		enabled[s] = true
	}

	var AllDataSources = []string{
		"Bloomberg",
		"DowJones",
		"LexisNexis",
		"NYTimes",
		"Polymarket",
		"Reddit",
		"Refinitiv",
		"Twitter",
		"WallStreetJournal",
		"YouTube",
	}

	for _, name := range AllDataSources {
		var cfg *activities.ResearchInput
		if enabled[name] {
			cfg = input.ResearchConfig
		}
		err = workflowext.ExecuteOptional(ctx, name, cfg).Get(ctx, nil)
		if err != nil {
			return nil, err
		}
	}

	// Step 3: Perform research
	var result activities.ResearchOutput
	err = workflowext.ExecuteOptional(ctx, "Research", input.ResearchConfig).Get(ctx, &result)
	if err != nil {
		return nil, err
	}

	return &ResearchWorkflowOutput{
		Research: result.Research,
	}, nil
}
