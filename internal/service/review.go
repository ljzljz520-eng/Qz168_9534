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
	record, err := s.store.Record(input.RecordID)
	if err != nil {
		return domain.ReviewResult{}, err
	}
	if s.reviewBarrier != nil {
		s.reviewBarrier()
	}
	now := s.clock()
	sequence := int(s.sequence.Add(1))
	review := domain.ReviewResult{ID: fmt.Sprintf("review-%s-%06d", input.RecordID, sequence), RecordID: input.RecordID, ViewerID: input.ViewerID, Decision: input.Decision, Note: input.Note, Sequence: sequence, CreatedAt: now.UTC()}
	updated, err := domain.ApplyReview(record, review, now)
	if err != nil {
		return review, err
	}
	if err = s.store.SaveReview(review); err != nil {
		return review, err
	}
	if err = s.store.SaveRecord(updated); err != nil {
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
