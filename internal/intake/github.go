package intake

import (
	"fmt"
	"strings"

	"github.com/jonathanpopham/estate-agent/internal/work"
)

type GitHubIssueEvent struct {
	Action     string
	Owner      string
	Repo       string
	Number     int
	Title      string
	Body       string
	LabelNames []string
}

func GitHubIssueToWorkItem(event GitHubIssueEvent) (work.Item, bool) {
	if event.Action != "opened" && event.Action != "edited" && event.Action != "labeled" && event.Action != "reopened" {
		return work.Item{}, false
	}

	kind, ok := kindFromLabels(event.LabelNames)
	if !ok {
		return work.Item{}, false
	}

	return work.Item{
		ID:          fmt.Sprintf("github:%s/%s#%d", event.Owner, event.Repo, event.Number),
		Kind:        kind,
		Source:      "github_issue",
		RepoOwner:   event.Owner,
		RepoName:    event.Repo,
		IssueNumber: event.Number,
		Title:       event.Title,
		Body:        event.Body,
		Labels:      event.LabelNames,
	}, true
}

func kindFromLabels(labels []string) (work.Kind, bool) {
	for _, label := range labels {
		switch strings.ToLower(strings.TrimSpace(label)) {
		case "bug", "estate:bug", "bug:autofix":
			return work.KindBug, true
		case "feature", "feature request", "estate:feature":
			return work.KindFeature, true
		}
	}
	return "", false
}
