package analytics

import (
	"sort"
	"tea17/internal/domain"
	"time"
)

type DailyPoint struct {
	Day           string `json:"day"`
	Registrations int    `json:"registrations"`
	Reviews       int    `json:"reviews"`
	Queries       int    `json:"queries"`
}

func Daily(events []domain.AuditEvent, from, to time.Time) []DailyPoint {
	points := map[string]*DailyPoint{}
	for day := from.UTC(); !day.After(to.UTC()); day = day.AddDate(0, 0, 1) {
		key := day.Format("2006-01-02")
		points[key] = &DailyPoint{Day: key}
	}
	for _, event := range events {
		key := event.OccurredAt.UTC().Format("2006-01-02")
		point := points[key]
		if point == nil {
			continue
		}
		switch event.Kind {
		case "record.registered":
			point.Registrations++
		case "record.reviewed":
			point.Reviews++
		case "recommendation.queried":
			point.Queries++
		}
	}
	out := make([]DailyPoint, 0, len(points))
	for _, point := range points {
		out = append(out, *point)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Day < out[j].Day })
	return out
}
func Conversion(records []domain.BeverageRecord) float64 {
	if len(records) == 0 {
		return 0
	}
	approved := 0
	for _, record := range records {
		if record.Status == domain.StatusApproved {
			approved++
		}
	}
	return float64(approved) / float64(len(records))
}
