package catalog

import (
	"sort"
	"strings"
	"tea17/internal/domain"
)

type IngredientInfo struct {
	Name            string
	Family          string
	Allergens       []string
	Tags            []string
	CaloriesPerUnit int
	CaffeineMG      int
	SeasonalMonths  []int
	Description     string
}
type Catalog struct {
	ingredients map[string]IngredientInfo
	categories  map[string][]string
}

func New() *Catalog {
	c := &Catalog{ingredients: map[string]IngredientInfo{}, categories: map[string][]string{}}
	for _, item := range DefaultIngredients() {
		c.ingredients[item.Name] = item
		c.categories[item.Family] = append(c.categories[item.Family], item.Name)
	}
	return c
}
func (c *Catalog) Ingredient(name string) (IngredientInfo, bool) {
	item, ok := c.ingredients[strings.ToLower(strings.TrimSpace(name))]
	return item, ok
}
func (c *Catalog) ValidateRecipe(record domain.BeverageRecord) []string {
	issues := []string{}
	for _, name := range record.Ingredients {
		if _, ok := c.Ingredient(name); !ok {
			issues = append(issues, "unknown ingredient: "+name)
		}
	}
	if len(record.Ingredients) > 12 {
		issues = append(issues, "recipe exceeds preparation limit")
	}
	sort.Strings(issues)
	return issues
}
func (c *Catalog) Enrich(record domain.BeverageRecord) domain.BeverageRecord {
	tags := append([]string{}, record.Tags...)
	calories, caffeine := 0, 0
	for _, name := range record.Ingredients {
		if item, ok := c.Ingredient(name); ok {
			tags = append(tags, item.Tags...)
			calories += item.CaloriesPerUnit
			caffeine += item.CaffeineMG
		}
	}
	record.Tags = unique(tags)
	if record.Calories == 0 {
		record.Calories = calories
	}
	if record.CaffeineMG == 0 {
		record.CaffeineMG = caffeine
	}
	return record
}
func (c *Catalog) Families() []string {
	out := make([]string, 0, len(c.categories))
	for family := range c.categories {
		out = append(out, family)
	}
	sort.Strings(out)
	return out
}
func (c *Catalog) Seasonal(month int) []IngredientInfo {
	out := []IngredientInfo{}
	for _, item := range c.ingredients {
		for _, m := range item.SeasonalMonths {
			if m == month {
				out = append(out, item)
				break
			}
		}
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Name < out[j].Name })
	return out
}
func unique(values []string) []string {
	seen := map[string]bool{}
	out := []string{}
	for _, v := range values {
		v = strings.ToLower(strings.TrimSpace(v))
		if v != "" && !seen[v] {
			seen[v] = true
			out = append(out, v)
		}
	}
	sort.Strings(out)
	return out
}
