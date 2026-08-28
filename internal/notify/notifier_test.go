package notify

import (
	"tea17/internal/domain"
	"testing"
	"time"
)

func TestMemorySender(t *testing.T) {
	sender := NewMemorySender()
	notice := ReviewNotification(domain.BeverageRecord{Name: "Tea", Version: 2}, domain.ReviewResult{ID: "r", ViewerID: "v", Decision: "approve"}, time.Unix(1, 0))
	if err := sender.Send(notice); err != nil {
		t.Fatal(err)
	}
	if len(sender.Sent()) != 1 || sender.Sent()[0].Recipient != "v" {
		t.Fatalf("unexpected notices: %+v", sender.Sent())
	}
}
