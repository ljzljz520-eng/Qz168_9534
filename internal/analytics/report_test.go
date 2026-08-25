package analytics

import (
	"tea17/internal/domain"
	"testing"
)

func TestPortfolio(t *testing.T) {
	report := Portfolio([]domain.BeverageRecord{{PriceCents: 400, Calories: 100, Status: domain.StatusApproved, Tags: []string{"fruit"}}, {PriceCents: 600, Calories: 200, Status: domain.StatusPending, Tags: []string{"fruit", "fresh"}}})
	if report.AveragePrice != 500 || report.AverageCalories != 150 {
		t.Fatalf("unexpected averages: %+v", report)
	}
	if len(report.PopularTags) == 0 || report.PopularTags[0].Tag != "fruit" {
		t.Fatalf("unexpected tags: %+v", report.PopularTags)
	}
}
