package domain

import (
	"testing"
	"time"
)

func TestNormalizeAndReview(t *testing.T) {
	now := time.Unix(1700000000, 0)
	r := NormalizeRecord(BeverageRecord{ID: "x", Name: "  New Tea ", Category: "FRUIT", Ingredients: []string{"Peach", "peach"}, PriceCents: 500}, now)
	if len(r.Ingredients) != 1 || r.Category != "fruit" {
		t.Fatalf("unexpected normalization: %+v", r)
	}
	updated, err := ApplyReview(r, ReviewResult{ID: "v1", RecordID: "x", ViewerID: "u", Decision: "approve"}, now)
	if err != nil || updated.Status != StatusApproved {
		t.Fatalf("unexpected review: %+v %v", updated, err)
	}
}
