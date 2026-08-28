package service

import (
	"errors"
	"fmt"
	"sync/atomic"
	"tea17/internal/catalog"
	"tea17/internal/domain"
	"tea17/internal/notify"
	"tea17/internal/recommend"
	"tea17/internal/store"
	"time"
)

type Service struct {
	store         *store.Store
	catalog       *catalog.Catalog
	engine        *recommend.Engine
	sender        notify.Sender
	clock         func() time.Time
	sequence      atomic.Int64
	reviewBarrier func()
}

func New(s *store.Store, c *catalog.Catalog, sender notify.Sender) *Service {
	return &Service{store: s, catalog: c, engine: recommend.New(c), sender: sender, clock: time.Now}
}
func (s *Service) SetClock(clock func() time.Time) { s.clock = clock }
func (s *Service) SetReviewBarrier(barrier func()) { s.reviewBarrier = barrier }
func (s *Service) Register(record domain.BeverageRecord, actor string) (domain.BeverageRecord, error) {
	record = s.catalog.Enrich(record)
	record = domain.NormalizeRecord(record, s.clock())
	if err := record.Validate(); err != nil {
		return record, err
	}
	if issues := s.catalog.ValidateRecipe(record); len(issues) > 0 {
		return record, fmt.Errorf("recipe validation: %v", issues)
	}
	if _, err := s.store.Record(record.ID); err == nil {
		return record, errors.New("record already exists")
	}
	if err := s.store.SaveRecord(record); err != nil {
		return record, err
	}
	event := domain.AuditEvent{ID: fmt.Sprintf("audit-register-%s", record.ID), Kind: "record.registered", SubjectID: record.ID, ActorID: actor, Payload: map[string]string{"name": record.Name, "status": string(record.Status)}, OccurredAt: s.clock().UTC()}
	if err := s.store.SaveAudit(event); err != nil {
		return record, err
	}
	return record, nil
}
func (s *Service) SubmitProfile(profile domain.CustomerProfile) error {
	if err := profile.Validate(); err != nil {
		return err
	}
	if profile.CreatedAt.IsZero() {
		profile.CreatedAt = s.clock().UTC()
	}
	return s.store.SaveProfile(profile)
}
func (s *Service) Recommendations(profileID string, limit int) ([]domain.Recommendation, error) {
	profile, err := s.store.Profile(profileID)
	if err != nil {
		return nil, err
	}
	records, err := s.store.Records()
	if err != nil {
		return nil, err
	}
	items, err := s.engine.Recommend(profile, records, limit)
	if err != nil {
		return nil, err
	}
	event := domain.AuditEvent{ID: fmt.Sprintf("audit-recommend-%s-%d", profileID, s.sequence.Add(1)), Kind: "recommendation.queried", SubjectID: profileID, ActorID: profileID, Payload: map[string]string{"count": fmt.Sprint(len(items))}, OccurredAt: s.clock().UTC()}
	if err := s.store.SaveAudit(event); err != nil {
		return nil, err
	}
	return items, nil
}
func (s *Service) Record(id string) (domain.BeverageRecord, error)  { return s.store.Record(id) }
func (s *Service) History(id string) ([]domain.ReviewResult, error) { return s.store.ReviewsFor(id) }
