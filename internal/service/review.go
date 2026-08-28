package service

import (
	"fmt"
	"tea17/internal/domain"
	"tea17/internal/notify"
)

type ReviewInput struct {
	RecordID string
	ViewerID string
	Decision string
	Note     string
}

func (s *Service) Review(input ReviewInput) (domain.ReviewResult, error) {
	now := s.clock()
	sequence := int(s.sequence.Add(1))
	review := domain.ReviewResult{ID: fmt.Sprintf("review-%s-%06d", input.RecordID, sequence), RecordID: input.RecordID, ViewerID: input.ViewerID, Decision: input.Decision, Note: input.Note, Sequence: sequence, CreatedAt: now.UTC()}

	if s.reviewBarrier != nil {
		s.reviewBarrier()
	}

	// Persist the review and update the record inside a single bbolt write
	// transaction. bbolt serializes write transactions, so concurrent reviewers
	// each read the most recent committed record state: a later reviewer sees the
	// earlier review already applied and appends to it instead of overwriting it.
	// The review itself is keyed by its unique sequence-derived ID, so even
	// concurrent submissions never collide.
	updated, err := s.store.SaveReviewTx(review, func(record domain.BeverageRecord) (domain.BeverageRecord, error) {
		return domain.ApplyReview(record, review, now)
	})
	if err != nil {
		return review, err
	}

	event := domain.AuditEvent{ID: fmt.Sprintf("audit-%s", review.ID), Kind: "record.reviewed", SubjectID: input.RecordID, ActorID: input.ViewerID, Payload: map[string]string{"decision": input.Decision, "version": fmt.Sprint(updated.Version)}, OccurredAt: now.UTC()}
	if err = s.store.SaveAudit(event); err != nil {
		return review, err
	}
	notice := notify.ReviewNotification(updated, review, now)
	if s.sender != nil {
		if err = s.sender.Send(notice); err != nil {
			return review, err
		}
	}
	return review, nil
}
