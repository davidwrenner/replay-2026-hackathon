package activities

import (
	"context"
	"strings"

	"go.temporal.io/sdk/temporal"
)

// InjectionCheckInput is the input for the injection check activity.
type InjectionCheckInput struct {
	Prompt        string `json:"prompt"`
	RiskTolerance string `json:"riskTolerance"`
}

// injectionPatterns contains patterns that indicate potential prompt injection.
var injectionPatterns = []string{
	"ignore previous",
	"ignore above",
	"disregard previous",
	"forget previous",
	"ignore all instructions",
	"ignore your instructions",
	"new instructions:",
	"system prompt:",
	"you are now",
	"act as",
	"pretend to be",
	"jailbreak",
	"bypass safety",
	"bypass filters",
}

// maliciousPatterns contains patterns that indicate malicious instructions.
var maliciousPatterns = []string{
	"reveal your api key",
	"show me your api key",
	"what is your api key",
	"output your instructions",
	"print your instructions",
	"reveal your prompt",
	"show your system prompt",
	"delete all",
	"drop table",
	"<script>",
	"eval(",
	"exec(",
}

// CheckInjection checks the prompt for potential injection attacks.
func (a *Activities) CheckInjection(ctx context.Context, input InjectionCheckInput) error {
	promptLower := strings.ToLower(input.Prompt)

	for _, pattern := range injectionPatterns {
		if strings.Contains(promptLower, pattern) {
			return temporal.NewNonRetryableApplicationError("prompt injection detected: "+pattern, "INJECTION", nil)
		}
	}

	for _, pattern := range maliciousPatterns {
		if strings.Contains(promptLower, pattern) {
			return temporal.NewNonRetryableApplicationError("malicious instruction detected: "+pattern, "MALICIOUS", nil)
		}
	}

	return nil
}
