package service

import (
	"path/filepath"
	"tea17/internal/catalog"
	"tea17/internal/domain"
	"tea17/internal/notify"
	"tea17/internal/store"
	"testing"
	"time"
)

func TestServiceCreatesAuditAndNotice(t *testing.T) {
	db, err := store.Open(filepath.Join(t.TempDir(), "s.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	sender := notify.NewMemorySender()
	svc := New(db, catalog.New(), sender)
	svc.SetClock(func() time.Time { return time.Unix(10, 0) })
	record := domain.BeverageRecord{ID: "r", Name: "Mango Tea", Category: "fruit", Ingredients: []string{"jasmine tea", "mango"}, PriceCents: 500, Status: domain.StatusPending}
	if _, err = svc.Register(record, "actor"); err != nil {
		t.Fatal(err)
	}
	if _, err = svc.Review(ReviewInput{RecordID: "r", ViewerID: "viewer", Decision: "approve"}); err != nil {
		t.Fatal(err)
	}
	if len(sender.Sent()) != 1 {
		t.Fatal("expected notification")
	}
	events, err := db.AuditsFor("r")
	if err != nil || len(events) != 2 {
		t.Fatalf("unexpected events: %+v %v", events, err)
	}
}
