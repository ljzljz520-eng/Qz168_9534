package recommend

import (
	"fmt"
	"strings"
	"tea17/internal/domain"
)

func Explain(r domain.Recommendation) string {
	parts := append([]string{}, r.Reasons...)
	if len(r.Warnings) > 0 {
		parts = append(parts, "notes: "+strings.Join(r.Warnings, ", "))
	}
	return fmt.Sprintf("%s scored %d: %s", r.Name, r.Score, strings.Join(parts, "; "))
}
func GroupByStrength(items []domain.Recommendation) map[string][]domain.Recommendation {
	groups := map[string][]domain.Recommendation{"excellent": {}, "good": {}, "other": {}}
	for _, item := range items {
		switch {
		case item.Score >= 80:
			groups["excellent"] = append(groups["excellent"], item)
		case item.Score >= 60:
			groups["good"] = append(groups["good"], item)
		default:
			groups["other"] = append(groups["other"], item)
		}
	}
	return groups
}
