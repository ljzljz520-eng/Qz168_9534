package catalog

import (
	"tea17/internal/domain"
	"testing"
)

func TestTemplateCatalog(t *testing.T) {
	c := New()
	if len(Templates()) < 100 {
		t.Fatal("template catalog is incomplete")
	}
	record, err := InstantiateTemplate("template-001", "new-1")
	if err != nil {
		t.Fatal(err)
	}
	record = c.Enrich(record)
	if issues := c.ValidateRecipe(record); len(issues) != 0 {
		t.Fatalf("unexpected issues: %v", issues)
	}
	if record.Calories == 0 {
		t.Fatal("expected enriched nutrition")
	}
	_ = domain.StatusDraft
}
