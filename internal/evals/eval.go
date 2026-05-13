package evals

import "strings"

type Fixture struct {
	Name         string
	Input        string
	ExpectedKind string
	MustEscalate bool
}

type Result struct {
	Name   string
	Pass   bool
	Reason string
}

func Run(fixtures []Fixture, classify func(string) (kind string, escalate bool)) []Result {
	results := make([]Result, 0, len(fixtures))
	for _, fixture := range fixtures {
		kind, escalate := classify(fixture.Input)
		result := Result{Name: fixture.Name, Pass: true}
		switch {
		case kind != fixture.ExpectedKind:
			result.Pass = false
			result.Reason = "kind mismatch"
		case escalate != fixture.MustEscalate:
			result.Pass = false
			result.Reason = "escalation mismatch"
		}
		results = append(results, result)
	}
	return results
}

func KeywordClassifier(input string) (kind string, escalate bool) {
	normalized := strings.ToLower(input)
	if strings.Contains(normalized, "secret") || strings.Contains(normalized, "billing") || strings.Contains(normalized, "auth") {
		escalate = true
	}
	if strings.Contains(normalized, "feature") || strings.Contains(normalized, "add ") {
		return "feature", escalate
	}
	return "bug", escalate
}
