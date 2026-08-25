package domain

import "time"

type Status string

const (
	StatusDraft    Status = "draft"
	StatusPending  Status = "pending"
	StatusApproved Status = "approved"
	StatusRejected Status = "rejected"
)

type BeverageRecord struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Category    string    `json:"category"`
	Description string    `json:"description"`
	Ingredients []string  `json:"ingredients"`
	Tags        []string  `json:"tags"`
	PriceCents  int       `json:"price_cents"`
	Calories    int       `json:"calories"`
	CaffeineMG  int       `json:"caffeine_mg"`
	Status      Status    `json:"status"`
	Version     int       `json:"version"`
	ReviewIDs   []string  `json:"review_ids"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
type ReviewResult struct {
	ID        string    `json:"id"`
	RecordID  string    `json:"record_id"`
	ViewerID  string    `json:"viewer_id"`
	Decision  string    `json:"decision"`
	Note      string    `json:"note"`
	Sequence  int       `json:"sequence"`
	CreatedAt time.Time `json:"created_at"`
}
type CustomerProfile struct {
	ID                 string    `json:"id"`
	PreferredTags      []string  `json:"preferred_tags"`
	AvoidIngredients   []string  `json:"avoid_ingredients"`
	MaxPriceCents      int       `json:"max_price_cents"`
	MaxCalories        int       `json:"max_calories"`
	CaffeinePreference string    `json:"caffeine_preference"`
	CreatedAt          time.Time `json:"created_at"`
}
type AuditEvent struct {
	ID         string            `json:"id"`
	Kind       string            `json:"kind"`
	SubjectID  string            `json:"subject_id"`
	ActorID    string            `json:"actor_id"`
	Payload    map[string]string `json:"payload"`
	OccurredAt time.Time         `json:"occurred_at"`
}
type Recommendation struct {
	RecordID string   `json:"record_id"`
	Name     string   `json:"name"`
	Score    int      `json:"score"`
	Reasons  []string `json:"reasons"`
	Warnings []string `json:"warnings"`
}
type Notification struct {
	ID        string    `json:"id"`
	Recipient string    `json:"recipient"`
	Subject   string    `json:"subject"`
	Body      string    `json:"body"`
	State     string    `json:"state"`
	CreatedAt time.Time `json:"created_at"`
}
