package evals

import "testing"

func TestRunKeywordFixtures(t *testing.T) {
	fixtures := []Fixture{
		{Name: "panic is bug", Input: "panic in checkout", ExpectedKind: "bug"},
		{Name: "feature request", Input: "feature: add export button", ExpectedKind: "feature"},
		{Name: "auth escalates", Input: "bug in auth token validation", ExpectedKind: "bug", MustEscalate: true},
	}

	results := Run(fixtures, KeywordClassifier)

	for _, result := range results {
		if !result.Pass {
			t.Fatalf("%s failed: %s", result.Name, result.Reason)
		}
	}
}
