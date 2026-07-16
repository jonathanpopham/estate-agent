package work

import "testing"

func TestKindValid(t *testing.T) {
	for _, k := range []Kind{KindBug, KindFeature} {
		if !k.Valid() {
			t.Errorf("expected %q to be valid", k)
		}
	}
	for _, k := range []Kind{"", "chore", "Bug", "chore:autofix"} {
		if k.Valid() {
			t.Errorf("expected %q to be invalid", k)
		}
	}
}

func TestItemValidate(t *testing.T) {
	good := Item{ID: "github:o/r#1", Kind: KindBug, Source: "github_issue", Title: "boom"}
	if err := good.Validate(); err != nil {
		t.Fatalf("expected valid item, got %v", err)
	}
	bad := map[string]Item{
		"no id":     {Kind: KindBug, Source: "s", Title: "t"},
		"blank id":  {ID: "   ", Kind: KindBug, Source: "s", Title: "t"},
		"bad kind":  {ID: "x", Kind: "nope", Source: "s", Title: "t"},
		"no title":  {ID: "x", Kind: KindBug, Source: "s"},
		"no source": {ID: "x", Kind: KindBug, Title: "t"},
	}
	for name, it := range bad {
		if err := it.Validate(); err == nil {
			t.Errorf("%s: expected a validation error, got nil", name)
		}
	}
}

func TestDedupeKeyStable(t *testing.T) {
	a := Item{ID: "error:abc123", Kind: KindBug, Source: "error", Title: "boom"}
	b := Item{ID: "error:abc123", Kind: KindBug, Source: "error", Title: "boom (again)"}
	if a.DedupeKey() != b.DedupeKey() {
		t.Errorf("same ID must share a dedupe key: %q vs %q", a.DedupeKey(), b.DedupeKey())
	}
	if a.DedupeKey() == (Item{ID: "error:zzz999"}).DedupeKey() {
		t.Error("different IDs must not share a dedupe key")
	}
}

func TestTally(t *testing.T) {
	items := []Item{
		{ID: "error:a"}, {ID: "error:a"}, {ID: "error:a"},
		{ID: "error:b"},
		{ID: "github:o/r#1"},
	}
	counts := Tally(items)
	if counts["error:a"] != 3 {
		t.Errorf("error:a want 3, got %d", counts["error:a"])
	}
	if counts["error:b"] != 1 {
		t.Errorf("error:b want 1, got %d", counts["error:b"])
	}
	if len(counts) != 3 {
		t.Errorf("want 3 distinct keys, got %d", len(counts))
	}
}

func TestActionable(t *testing.T) {
	items := []Item{
		{ID: "error:a", Title: "first a"},
		{ID: "error:b", Title: "b"},
		{ID: "error:a", Title: "second a"},
		{ID: "error:a", Title: "third a"},
	}
	got := Actionable(items, 3)
	if len(got) != 1 {
		t.Fatalf("want 1 actionable item at threshold 3, got %d", len(got))
	}
	if got[0].ID != "error:a" {
		t.Errorf("want error:a, got %s", got[0].ID)
	}
	if got[0].Title != "first a" {
		t.Errorf("Actionable must return the first-seen item for determinism, got %q", got[0].Title)
	}

	all := Actionable(items, 1)
	if len(all) != 2 || all[0].ID != "error:a" || all[1].ID != "error:b" {
		t.Errorf("threshold 1 should return distinct items in first-seen order, got %+v", all)
	}

	if len(Actionable(items, 0)) != 2 {
		t.Error("a threshold below 1 should clamp to 1")
	}
	if len(Actionable(nil, 2)) != 0 {
		t.Error("no items should yield nothing actionable")
	}
}
