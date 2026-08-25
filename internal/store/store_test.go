package store

import (
	"path/filepath"
	"tea17/internal/domain"
	"testing"
	"time"
)

func TestPersistenceSurvivesReopen(t *testing.T) {
	path := filepath.Join(t.TempDir(), "tea.db")
	s, err := Open(path)
	if err != nil {
		t.Fatal(err)
	}
	record := domain.BeverageRecord{ID: "persisted", Name: "Stored Tea", Category: "tea", Ingredients: []string{"oolong tea"}, PriceCents: 450, CreatedAt: time.Unix(1, 0)}
	if err = s.SaveRecord(record); err != nil {
		t.Fatal(err)
	}
	if err = s.SaveReview(domain.ReviewResult{ID: "review-p", RecordID: record.ID, ViewerID: "v", Decision: "approve"}); err != nil {
		t.Fatal(err)
	}
	if err = s.SaveProfile(domain.CustomerProfile{ID: "profile-p"}); err != nil {
		t.Fatal(err)
	}
	if err = s.SaveAudit(domain.AuditEvent{ID: "audit-p", SubjectID: record.ID}); err != nil {
		t.Fatal(err)
	}
	if err = s.Close(); err != nil {
		t.Fatal(err)
	}
	s, err = Open(path)
	if err != nil {
		t.Fatal(err)
	}
	defer s.Close()
	found, err := s.Record(record.ID)
	if err != nil {
		t.Fatal(err)
	}
	if found.Name != record.Name {
		t.Fatalf("unexpected record: %+v", found)
	}
	reviews, err := s.ReviewsFor(record.ID)
	if err != nil || len(reviews) != 1 {
		t.Fatalf("unexpected reviews: %+v %v", reviews, err)
	}
	if _, err = s.Profile("profile-p"); err != nil {
		t.Fatal(err)
	}
	audits, err := s.AuditsFor(record.ID)
	if err != nil || len(audits) != 1 {
		t.Fatalf("unexpected audits: %+v %v", audits, err)
	}
}
