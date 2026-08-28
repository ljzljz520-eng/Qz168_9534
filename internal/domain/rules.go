package domain

import (
	"errors"
	"fmt"
	"sort"
	"strings"
	"time"
)

func (r BeverageRecord) Validate() error {
	if strings.TrimSpace(r.ID) == "" {
		return errors.New("record id is required")
	}
	if len([]rune(strings.TrimSpace(r.Name))) < 2 {
		return errors.New("record name is too short")
	}
	if r.Category == "" {
		return errors.New("category is required")
	}
	if len(r.Ingredients) == 0 {
		return errors.New("ingredients are required")
	}
	if r.PriceCents < 100 || r.PriceCents > 10000 {
		return errors.New("price is outside store range")
	}
	if r.Calories < 0 || r.Calories > 2000 {
		return errors.New("calories are outside supported range")
	}
	for _, ingredient := range r.Ingredients {
		if strings.TrimSpace(ingredient) == "" {
			return errors.New("ingredient cannot be blank")
		}
	}
	return nil
}
func (p CustomerProfile) Validate() error {
	if p.ID == "" {
		return errors.New("profile id is required")
	}
	if p.MaxPriceCents < 0 {
		return errors.New("maximum price cannot be negative")
	}
	if p.MaxCalories < 0 {
		return errors.New("maximum calories cannot be negative")
	}
	switch p.CaffeinePreference {
	case "", "any", "low", "regular", "high":
		return nil
	default:
		return fmt.Errorf("unknown caffeine preference %q", p.CaffeinePreference)
	}
}
func NormalizeRecord(r BeverageRecord, now time.Time) BeverageRecord {
	r.Name = strings.TrimSpace(r.Name)
	r.Category = strings.ToLower(strings.TrimSpace(r.Category))
	r.Description = strings.TrimSpace(r.Description)
	r.Ingredients = uniqueSorted(r.Ingredients)
	r.Tags = uniqueSorted(r.Tags)
	if r.Status == "" {
		r.Status = StatusDraft
	}
	if r.Version == 0 {
		r.Version = 1
	}
	if r.CreatedAt.IsZero() {
		r.CreatedAt = now.UTC()
	}
	r.UpdatedAt = now.UTC()
	return r
}
func ApplyReview(r BeverageRecord, review ReviewResult, now time.Time) (BeverageRecord, error) {
	if review.RecordID != r.ID {
		return r, errors.New("review belongs to another record")
	}
	if review.ViewerID == "" {
		return r, errors.New("viewer is required")
	}
	switch review.Decision {
	case "approve":
		r.Status = StatusApproved
	case "reject":
		r.Status = StatusRejected
	case "hold":
		r.Status = StatusPending
	default:
		return r, errors.New("unsupported review decision")
	}
	r.ReviewIDs = append(r.ReviewIDs, review.ID)
	r.Version++
	r.UpdatedAt = now.UTC()
	return r, nil
}
func CanTransition(from, to Status) bool {
	if from == to {
		return true
	}
	switch from {
	case StatusDraft:
		return to == StatusPending
	case StatusPending:
		return to == StatusApproved || to == StatusRejected
	case StatusRejected:
		return to == StatusDraft
	default:
		return false
	}
}
func uniqueSorted(values []string) []string {
	seen := map[string]bool{}
	out := make([]string, 0, len(values))
	for _, value := range values {
		value = strings.ToLower(strings.TrimSpace(value))
		if value != "" && !seen[value] {
			seen[value] = true
			out = append(out, value)
		}
	}
	sort.Strings(out)
	return out
}
