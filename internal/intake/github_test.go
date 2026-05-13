package intake

import (
	"testing"

	"github.com/jonathanpopham/estate-agent/internal/work"
)

func TestGitHubIssueToWorkItemFeature(t *testing.T) {
	item, ok := GitHubIssueToWorkItem(GitHubIssueEvent{
		Action:     "opened",
		Owner:      "jonathanpopham",
		Repo:       "estate-agent",
		Number:     5,
		Title:      "Add deployment target",
		LabelNames: []string{"feature request"},
	})

	if !ok {
		t.Fatal("event ignored")
	}
	if item.Kind != work.KindFeature {
		t.Fatalf("Kind = %q, want %q", item.Kind, work.KindFeature)
	}
	if item.ID != "github:jonathanpopham/estate-agent#5" {
		t.Fatalf("ID = %q", item.ID)
	}
}

func TestGitHubIssueToWorkItemIgnoresUnlabeledIssue(t *testing.T) {
	_, ok := GitHubIssueToWorkItem(GitHubIssueEvent{
		Action: "opened",
		Title:  "Question",
	})

	if ok {
		t.Fatal("unlabeled issue should be ignored")
	}
}
