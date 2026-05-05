package workflows

import (
	"time"

	"github.com/davidwrenner/replay-2026-hackathon/activities"
	"go.temporal.io/sdk/workflow"
)

// ResearchWorkflowOld represents our workflow without using the workflow package. This shows how verbose this pattern would
// have been without this project
func ResearchWorkflowOld(ctx workflow.Context, input ResearchWorkflowInput) (*ResearchWorkflowOutput, error) {
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

	// Step 2: Fetch data from sources
	if input.BloombergConfig != nil {
		err = workflow.ExecuteActivity(ctx, "Bloomberg", input.BloombergConfig).Get(ctx, nil)
		if err != nil {
			return nil, err
		}
	}

	if input.DowJonesConfig != nil {
		err = workflow.ExecuteActivity(ctx, "DowJones", input.DowJonesConfig).Get(ctx, nil)
		if err != nil {
			return nil, err
		}
	}

	if input.LexisNexisConfig != nil {
		err = workflow.ExecuteActivity(ctx, "LexisNexis", input.LexisNexisConfig).Get(ctx, nil)
		if err != nil {
			return nil, err
		}
	}

	if input.NYTimesConfig != nil {
		err = workflow.ExecuteActivity(ctx, "NYTimes", input.NYTimesConfig).Get(ctx, nil)
		if err != nil {
			return nil, err
		}
	}

	if input.PolymarketConfig != nil {
		err = workflow.ExecuteActivity(ctx, "Polymarket", input.PolymarketConfig).Get(ctx, nil)
		if err != nil {
			return nil, err
		}
	}

	if input.RedditConfig != nil {
		err = workflow.ExecuteActivity(ctx, "Reddit", input.RedditConfig).Get(ctx, nil)
		if err != nil {
			return nil, err
		}
	}

	if input.RefinitivConfig != nil {
		err = workflow.ExecuteActivity(ctx, "Refinitiv", input.RefinitivConfig).Get(ctx, nil)
		if err != nil {
			return nil, err
		}
	}

	if input.TwitterConfig != nil {
		err = workflow.ExecuteActivity(ctx, "Twitter", input.TwitterConfig).Get(ctx, nil)
		if err != nil {
			return nil, err
		}
	}

	if input.WallStreetJournalConfig != nil {
		err = workflow.ExecuteActivity(ctx, "WallStreetJournal", input.WallStreetJournalConfig).Get(ctx, nil)
		if err != nil {
			return nil, err
		}
	}

	if input.YouTubeConfig != nil {
		err = workflow.ExecuteActivity(ctx, "YouTube", input.YouTubeConfig).Get(ctx, nil)
		if err != nil {
			return nil, err
		}
	}

	// Step 3: Perform research
	var result activities.ResearchOutput
	err = workflow.ExecuteActivity(ctx, "Research", input.ResearchConfig).Get(ctx, &result)
	if err != nil {
		return nil, err
	}

	return &ResearchWorkflowOutput{
		Research: result.Research,
	}, nil
}
