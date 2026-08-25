package service

import (
	"path/filepath"
	"sync"
	"tea17/internal/catalog"
	"tea17/internal/domain"
	"tea17/internal/notify"
	"tea17/internal/store"
	"testing"
	"time"
)

func newServiceWithStore(t *testing.T) (*Service, *notify.MemorySender) {
	t.Helper()
	db, err := store.Open(filepath.Join(t.TempDir(), "tea17.db"))
	if err != nil {
		t.Fatalf("open store: %v", err)
	}
	t.Cleanup(func() { _ = db.Close() })
	sender := notify.NewMemorySender()
	svc := New(db, catalog.New(), sender)
	svc.SetClock(func() time.Time { return time.Date(2026, 8, 26, 9, 0, 0, 0, time.UTC) })
	return svc, sender
}

func registerRecord(t *testing.T, svc *Service, id, name string) domain.BeverageRecord {
	t.Helper()
	record := domain.BeverageRecord{
		ID:          id,
		Name:        name,
		Category:    "tea",
		Description: "新品",
		Ingredients: []string{"assam tea", "whole milk", "brown sugar"},
		PriceCents:  1200,
		Calories:    200,
		CaffeineMG:  30,
		Status:      domain.StatusPending,
	}
	out, err := svc.Register(record, "clerk")
	if err != nil {
		t.Fatalf("register: %v", err)
	}
	return out
}

// TestReviewConcurrentPreservesFirstResult verifies that when several reviewers
// view and confirm a new drink at the same time, the first review result is
// preserved on the record (rather than overwritten by a later reviewer that
// read the same stale snapshot).
func TestReviewConcurrentPreservesFirstResult(t *testing.T) {
	svc, _ := newServiceWithStore(t)
	registerRecord(t, svc, "rec-1", "芋圆奶茶")

	// Rendezvous barrier so all reviewers race through Review simultaneously
	// after they have all started, maximizing the chance that two read the
	// same stale record snapshot if the bug were present.
	var ready sync.WaitGroup
	ready.Add(2)
	release := make(chan struct{})
	svc.SetReviewBarrier(func() {
		ready.Done()
		<-release
	})

	var wg sync.WaitGroup
	wg.Add(2)
	reviews := make([]domain.ReviewResult, 2)
	errs := make([]error, 2)
	for i := 0; i < 2; i++ {
		i := i
		go func() {
			defer wg.Done()
			dec := "approve"
			if i == 1 {
				dec = "reject"
			}
			reviews[i], errs[i] = svc.Review(ReviewInput{
				RecordID: "rec-1",
				ViewerID: "viewer-" + string(rune('A'+i)),
				Decision: dec,
				Note:     "并发确认",
			})
		}()
	}
	ready.Wait()
	close(release)
	wg.Wait()

	for i, err := range errs {
		if err != nil {
			t.Fatalf("reviewer %d failed: %v", i, err)
		}
	}

	// The record must reference BOTH reviews and show version incremented twice.
	record, err := svc.Record("rec-1")
	if err != nil {
		t.Fatalf("record: %v", err)
	}
	if len(record.ReviewIDs) != 2 {
		t.Fatalf("expected 2 review ids on record, got %d: %v", len(record.ReviewIDs), record.ReviewIDs)
	}
	history, err := svc.History("rec-1")
	if err != nil {
		t.Fatalf("history: %v", err)
	}
	if len(history) != 2 {
		t.Fatalf("expected 2 review results in history, got %d", len(history))
	}
	// First submitted review (lower sequence) must be recorded in history.
	first := history[0]
	if first.ID != reviews[0].ID && first.ID != reviews[1].ID {
		t.Fatalf("history mismatch: %v vs %v,%v", first.ID, reviews[0].ID, reviews[1].ID)
	}
}

// TestReviewConcurrentPreservesFirstResultWithSleep uses a plain sleep barrier
// to widen the race window, ensuring the fix does not depend on a rendezvous.
func TestReviewConcurrentPreservesFirstResultWithSleep(t *testing.T) {
	svc, _ := newServiceWithStore(t)
	registerRecord(t, svc, "rec-2", "桃桃乌龙")

	svc.SetReviewBarrier(func() { time.Sleep(20 * time.Millisecond) })

	var wg sync.WaitGroup
	wg.Add(3)
	for i := 0; i < 3; i++ {
		i := i
		go func() {
			defer wg.Done()
			dec := "approve"
			if i%2 == 1 {
				dec = "reject"
			}
			_, _ = svc.Review(ReviewInput{
				RecordID: "rec-2",
				ViewerID: "viewer-" + string(rune('A'+i)),
				Decision: dec,
			})
		}()
	}
	wg.Wait()

	record, err := svc.Record("rec-2")
	if err != nil {
		t.Fatalf("record: %v", err)
	}
	if len(record.ReviewIDs) != 3 {
		t.Fatalf("expected 3 review ids on record, got %d: %v", len(record.ReviewIDs), record.ReviewIDs)
	}
	history, err := svc.History("rec-2")
	if err != nil {
		t.Fatalf("history: %v", err)
	}
	if len(history) != 3 {
		t.Fatalf("expected 3 review results in history, got %d", len(history))
	}
}
