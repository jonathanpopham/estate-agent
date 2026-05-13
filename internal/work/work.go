package work

type Kind string

const (
	KindBug     Kind = "bug"
	KindFeature Kind = "feature"
)

type Item struct {
	ID          string
	Kind        Kind
	Source      string
	RepoOwner   string
	RepoName    string
	IssueNumber int
	Title       string
	Body        string
	Labels      []string
	Severity    string
}
