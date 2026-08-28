package recommend

import (
	"sort"
	"strings"
	"tea17/internal/catalog"
	"tea17/internal/domain"
)

type Engine struct{ catalog *catalog.Catalog }

func New(c *catalog.Catalog) *Engine { return &Engine{catalog: c} }
func (e *Engine) Recommend(profile domain.CustomerProfile, records []domain.BeverageRecord, limit int) ([]domain.Recommendation, error) {
	if err := profile.Validate(); err != nil {
		return nil, err
	}
	ranked := []domain.Recommendation{}
	for _, record := range records {
		if record.Status != domain.StatusApproved {
			continue
		}
		candidate, ok := e.score(profile, record)
		if ok {
			ranked = append(ranked, candidate)
		}
	}
	sort.SliceStable(ranked, func(i, j int) bool {
		if ranked[i].Score == ranked[j].Score {
			return ranked[i].Name < ranked[j].Name
		}
		return ranked[i].Score > ranked[j].Score
	})
	if limit > 0 && len(ranked) > limit {
		ranked = ranked[:limit]
	}
	return ranked, nil
}
func (e *Engine) score(profile domain.CustomerProfile, record domain.BeverageRecord) (domain.Recommendation, bool) {
	result := domain.Recommendation{RecordID: record.ID, Name: record.Name, Score: 50}
	if profile.MaxPriceCents > 0 && record.PriceCents > profile.MaxPriceCents {
		return result, false
	}
	if profile.MaxCalories > 0 && record.Calories > profile.MaxCalories {
		return result, false
	}
	for _, name := range record.Ingredients {
		item, ok := e.catalog.Ingredient(name)
		if !ok {
			result.Warnings = append(result.Warnings, "unverified ingredient: "+name)
			continue
		}
		for _, avoid := range profile.AvoidIngredients {
			if strings.EqualFold(avoid, name) {
				return result, false
			}
			for _, allergen := range item.Allergens {
				if strings.EqualFold(avoid, allergen) {
					return result, false
				}
			}
		}
	}
	for _, preferred := range profile.PreferredTags {
		matched := false
		for _, tag := range record.Tags {
			if strings.EqualFold(preferred, tag) {
				matched = true
				break
			}
		}
		if matched {
			result.Score += 12
			result.Reasons = append(result.Reasons, "matches "+preferred)
		}
	}
	switch profile.CaffeinePreference {
	case "low":
		if record.CaffeineMG <= 20 {
			result.Score += 15
			result.Reasons = append(result.Reasons, "low caffeine")
		} else {
			result.Score -= 10
		}
	case "high":
		if record.CaffeineMG >= 50 {
			result.Score += 15
			result.Reasons = append(result.Reasons, "high caffeine")
		}
	case "regular":
		if record.CaffeineMG > 20 && record.CaffeineMG < 60 {
			result.Score += 10
			result.Reasons = append(result.Reasons, "balanced caffeine")
		}
	}
	if record.PriceCents <= profile.MaxPriceCents*8/10 {
		result.Score += 5
		result.Reasons = append(result.Reasons, "comfortably within budget")
	}
	if len(result.Reasons) == 0 {
		result.Reasons = append(result.Reasons, "available approved new drink")
	}
	return result, true
}
