package notify

import (
	"errors"
	"fmt"
	"sync"
	"tea17/internal/domain"
	"time"
)

type Sender interface {
	Send(domain.Notification) error
}
type MemorySender struct {
	mu            sync.Mutex
	sent          []domain.Notification
	failRecipient string
}

func NewMemorySender() *MemorySender { return &MemorySender{} }
func (s *MemorySender) Send(n domain.Notification) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if n.Recipient == "" {
		return errors.New("recipient is required")
	}
	if n.Recipient == s.failRecipient {
		return errors.New("delivery rejected")
	}
	s.sent = append(s.sent, n)
	return nil
}
func (s *MemorySender) Sent() []domain.Notification {
	s.mu.Lock()
	defer s.mu.Unlock()
	return append([]domain.Notification{}, s.sent...)
}
func ReviewNotification(record domain.BeverageRecord, review domain.ReviewResult, now time.Time) domain.Notification {
	return domain.Notification{ID: "notice-" + review.ID, Recipient: review.ViewerID, Subject: "Review recorded for " + record.Name, Body: fmt.Sprintf("Decision %s was recorded at version %d", review.Decision, record.Version), State: "queued", CreatedAt: now.UTC()}
}
