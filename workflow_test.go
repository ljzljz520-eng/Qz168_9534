package tea17_test

import (
	"path/filepath"
	"sync"
	"tea17"
	"tea17/internal/catalog"
	"tea17/internal/domain"
	"tea17/internal/service"
	"testing"
)

func recordFixture() domain.BeverageRecord {
	return domain.BeverageRecord{ID: "drink-17", Name: "Peach Jasmine 17", Category: "fruit-tea", Ingredients: []string{"jasmine tea", "peach", "honey"}, Tags: []string{"fruity", "fresh"}, PriceCents: 520, Status: domain.StatusPending}
}
func TestWorkflowOne(t *testing.T) {
	app, err := tea17.Open(filepath.Join(t.TempDir(), "workflow.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer app.Close()
	created, err := app.Service.Register(recordFixture(), "clerk-1")
	if err != nil {
		t.Fatal(err)
	}
	found, err := app.Service.Record(created.ID)
	if err != nil {
		t.Fatal(err)
	}
	if found.Name != "Peach Jasmine 17" {
		t.Fatalf("unexpected record: %+v", found)
	}
}
func TestWorkflowTwo(t *testing.T) {
	app, err := tea17.Open(filepath.Join(t.TempDir(), "review.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer app.Close()
	if _, err = app.Service.Register(recordFixture(), "clerk-2"); err != nil {
		t.Fatal(err)
	}
	if _, err = app.Service.Review(service.ReviewInput{RecordID: "drink-17", ViewerID: "reviewer-1", Decision: "approve"}); err != nil {
		t.Fatal(err)
	}
	history, err := app.Service.History("drink-17")
	if err != nil {
		t.Fatal(err)
	}
	if len(history) != 1 {
		t.Fatalf("want one review, got %d", len(history))
	}
}
func TestWorkflowThree(t *testing.T) {
	app, err := tea17.Open(filepath.Join(t.TempDir(), "recommend.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer app.Close()
	record := recordFixture()
	record.Status = domain.StatusApproved
	if _, err = app.Service.Register(record, "clerk-3"); err != nil {
		t.Fatal(err)
	}
	profile := domain.CustomerProfile{ID: "guest-1", PreferredTags: []string{"fruity"}, AvoidIngredients: []string{"milk"}, MaxPriceCents: 700, MaxCalories: 300, CaffeinePreference: "regular"}
	if err = app.Service.SubmitProfile(profile); err != nil {
		t.Fatal(err)
	}
	items, err := app.Service.Recommendations(profile.ID, 3)
	if err != nil {
		t.Fatal(err)
	}
	if len(items) != 1 || items[0].RecordID != "drink-17" {
		t.Fatalf("unexpected recommendations: %+v", items)
	}
}
func TestRecordFlow17(t *testing.T) {
	app, err := tea17.Open(filepath.Join(t.TempDir(), "concurrent.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer app.Close()
	if _, err = app.Service.Register(recordFixture(), "clerk-17"); err != nil {
		t.Fatal(err)
	}
	arrived := make(chan struct{}, 2)
	release := make(chan struct{})
	app.Service.SetReviewBarrier(func() { arrived <- struct{}{}; <-release })
	inputs := []service.ReviewInput{{RecordID: "drink-17", ViewerID: "viewer-a", Decision: "approve"}, {RecordID: "drink-17", ViewerID: "viewer-b", Decision: "hold"}}
	var wg sync.WaitGroup
	wg.Add(2)
	errs := make(chan error, 2)
	for _, input := range inputs {
		go func(in service.ReviewInput) { defer wg.Done(); _, e := app.Service.Review(in); errs <- e }(input)
	}
	<-arrived
	<-arrived
	close(release)
	wg.Wait()
	close(errs)
	for e := range errs {
		if e != nil {
			t.Fatal(e)
		}
	}
	saved, err := app.Service.Record("drink-17")
	if err != nil {
		t.Fatal(err)
	}
	if len(saved.ReviewIDs) != 2 {
		t.Fatalf("record should retain both confirmation results, got %v", saved.ReviewIDs)
	}
}

var _ = catalog.Templates
