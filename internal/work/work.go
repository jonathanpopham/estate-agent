package work

import (
	"errors"
	"fmt"
	"strings"
)

type Kind string

const (
	KindBug     Kind = "bug"
	KindFeature Kind = "feature"
)

// Valid reports whether k is a recognized work kind.
func (k Kind) Valid() bool {
	switch k {
	case KindBug, KindFeature:
		return true
	default:
		return false
	}
}

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

// DedupeKey returns the stable identity of a work item. Two signals that
// describe the same underlying work (the same error fingerprint, or the same
// GitHub issue) share a key, so repeated reports collapse to one item.
func (it Item) DedupeKey() string {
	return it.ID
}

// Validate returns the first reason an item is not safe to act on, or nil. The
// agent must refuse to mutate code for an item that does not validate, rather
// than guess at what a malformed signal meant.
func (it Item) Validate() error {
	if strings.TrimSpace(it.ID) == "" {
		return errors.New("work item has no ID")
	}
	if !it.Kind.Valid() {
		return fmt.Errorf("work item %s has invalid kind %q", it.ID, it.Kind)
	}
	if strings.TrimSpace(it.Title) == "" {
		return fmt.Errorf("work item %s has no title", it.ID)
	}
	if strings.TrimSpace(it.Source) == "" {
		return fmt.Errorf("work item %s has no source", it.ID)
	}
	return nil
}

// Tally groups items by dedupe key and counts how many signals map to each.
func Tally(items []Item) map[string]int {
	counts := make(map[string]int, len(items))
	for _, it := range items {
		counts[it.DedupeKey()]++
	}
	return counts
}

// Actionable returns the distinct items whose signal count reaches threshold,
// deduped by DedupeKey and returned in first-seen order so the output is
// deterministic. It is the guard the product describes: act once, and only when
// enough reports describe the same problem. A threshold below 1 is clamped to 1.
func Actionable(items []Item, threshold int) []Item {
	if threshold < 1 {
		threshold = 1
	}
	counts := make(map[string]int, len(items))
	first := make(map[string]Item, len(items))
	order := make([]string, 0, len(items))
	for _, it := range items {
		k := it.DedupeKey()
		if _, seen := first[k]; !seen {
			first[k] = it
			order = append(order, k)
		}
		counts[k]++
	}
	out := make([]Item, 0)
	for _, k := range order {
		if counts[k] >= threshold {
			out = append(out, first[k])
		}
	}
	return out
}
