package analytics

import (
	"sort"
	"tea17/internal/domain"
	"time"
)

type StatusCount struct {
	Status domain.Status `json:"status"`
	Count  int           `json:"count"`
}
type ReviewActivity struct {
	ViewerID string    `json:"viewer_id"`
	Total    int       `json:"total"`
	Approved int       `json:"approved"`
	Rejected int       `json:"rejected"`
	LastSeen time.Time `json:"last_seen"`
}
type PortfolioReport struct {
	TotalRecords    int           `json:"total_records"`
	AveragePrice    int           `json:"average_price"`
	AverageCalories int           `json:"average_calories"`
	Statuses        []StatusCount `json:"statuses"`
	PopularTags     []TagCount    `json:"popular_tags"`
}
type TagCount struct {
	Tag   string `json:"tag"`
	Count int    `json:"count"`
}

func Portfolio(records []domain.BeverageRecord) PortfolioReport {
	report := PortfolioReport{TotalRecords: len(records)}
	statuses := map[domain.Status]int{}
	tags := map[string]int{}
	for _, r := range records {
		report.AveragePrice += r.PriceCents
		report.AverageCalories += r.Calories
		statuses[r.Status]++
		for _, tag := range r.Tags {
			tags[tag]++
		}
	}
	if len(records) > 0 {
		report.AveragePrice /= len(records)
		report.AverageCalories /= len(records)
	}
	for status, count := range statuses {
		report.Statuses = append(report.Statuses, StatusCount{Status: status, Count: count})
	}
	sort.Slice(report.Statuses, func(i, j int) bool { return report.Statuses[i].Status < report.Statuses[j].Status })
	for tag, count := range tags {
		report.PopularTags = append(report.PopularTags, TagCount{Tag: tag, Count: count})
	}
	sort.Slice(report.PopularTags, func(i, j int) bool {
		if report.PopularTags[i].Count == report.PopularTags[j].Count {
			return report.PopularTags[i].Tag < report.PopularTags[j].Tag
		}
		return report.PopularTags[i].Count > report.PopularTags[j].Count
	})
	if len(report.PopularTags) > 10 {
		report.PopularTags = report.PopularTags[:10]
	}
	return report
}
func Reviewers(reviews []domain.ReviewResult) []ReviewActivity {
	byViewer := map[string]*ReviewActivity{}
	for _, review := range reviews {
		item := byViewer[review.ViewerID]
		if item == nil {
			item = &ReviewActivity{ViewerID: review.ViewerID}
			byViewer[review.ViewerID] = item
		}
		item.Total++
		if review.Decision == "approve" {
			item.Approved++
		}
		if review.Decision == "reject" {
			item.Rejected++
		}
		if review.CreatedAt.After(item.LastSeen) {
			item.LastSeen = review.CreatedAt
		}
	}
	out := make([]ReviewActivity, 0, len(byViewer))
	for _, item := range byViewer {
		out = append(out, *item)
	}
	sort.Slice(out, func(i, j int) bool {
		if out[i].Total == out[j].Total {
			return out[i].ViewerID < out[j].ViewerID
		}
		return out[i].Total > out[j].Total
	})
	return out
}
