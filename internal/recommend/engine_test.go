package recommend

import (
	"tea17/internal/catalog"
	"tea17/internal/domain"
	"testing"
)

func TestRecommendationFiltersAllergen(t *testing.T) {
	engine := New(catalog.New())
	profile := domain.CustomerProfile{ID: "p", AvoidIngredients: []string{"milk"}, MaxPriceCents: 800, MaxCalories: 500, CaffeinePreference: "any"}
	records := []domain.BeverageRecord{{ID: "a", Name: "Dairy", Status: domain.StatusApproved, Ingredients: []string{"whole milk"}, PriceCents: 500}, {ID: "b", Name: "Fruit", Status: domain.StatusApproved, Ingredients: []string{"jasmine tea", "peach"}, PriceCents: 500}}
	items, err := engine.Recommend(profile, records, 5)
	if err != nil {
		t.Fatal(err)
	}
	if len(items) != 1 || items[0].RecordID != "b" {
		t.Fatalf("unexpected items: %+v", items)
	}
}
