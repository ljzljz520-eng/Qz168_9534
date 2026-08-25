package store

import (
	"encoding/json"
	"go.etcd.io/bbolt"
	"sort"
	"tea17/internal/domain"
)

func (s *Store) SaveRecord(record domain.BeverageRecord) error {
	return s.db.Update(func(tx *bbolt.Tx) error { return putJSON(tx, recordsBucket, record.ID, record) })
}
func (s *Store) Record(id string) (domain.BeverageRecord, error) {
	var out domain.BeverageRecord
	err := s.db.View(func(tx *bbolt.Tx) error { return getJSON(tx, recordsBucket, id, &out) })
	return out, err
}
func getRecordInTx(tx *bbolt.Tx, id string) (domain.BeverageRecord, error) {
	var out domain.BeverageRecord
	return out, getJSON(tx, recordsBucket, id, &out)
}
func (s *Store) Records() ([]domain.BeverageRecord, error) {
	out := []domain.BeverageRecord{}
	err := s.db.View(func(tx *bbolt.Tx) error {
		return scanJSON(tx, recordsBucket, func(data []byte) error {
			var v domain.BeverageRecord
			if err := json.Unmarshal(data, &v); err != nil {
				return err
			}
			out = append(out, v)
			return nil
		})
	})
	sort.Slice(out, func(i, j int) bool { return out[i].ID < out[j].ID })
	return out, err
}
func (s *Store) SaveReview(review domain.ReviewResult) error {
	return s.db.Update(func(tx *bbolt.Tx) error { return putJSON(tx, reviewsBucket, review.ID, review) })
}

// SaveReviewTx persists a review result and its updated record in a single
// transaction. Because bbolt serializes write transactions, each concurrent
// reviewer observes the most recent record state instead of a stale snapshot,
// so earlier review results are preserved rather than overwritten by a later
// reviewer that read the same stale record.
func (s *Store) SaveReviewTx(review domain.ReviewResult, update func(domain.BeverageRecord) (domain.BeverageRecord, error)) (domain.BeverageRecord, error) {
	var updated domain.BeverageRecord
	err := s.db.Update(func(tx *bbolt.Tx) error {
		record, err := getRecordInTx(tx, review.RecordID)
		if err != nil {
			return err
		}
		updated, err = update(record)
		if err != nil {
			return err
		}
		if err = putJSON(tx, reviewsBucket, review.ID, review); err != nil {
			return err
		}
		return putJSON(tx, recordsBucket, review.RecordID, updated)
	})
	return updated, err
}
func (s *Store) ReviewsFor(recordID string) ([]domain.ReviewResult, error) {
	out := []domain.ReviewResult{}
	err := s.db.View(func(tx *bbolt.Tx) error {
		return scanJSON(tx, reviewsBucket, func(data []byte) error {
			var v domain.ReviewResult
			if err := json.Unmarshal(data, &v); err != nil {
				return err
			}
			if v.RecordID == recordID {
				out = append(out, v)
			}
			return nil
		})
	})
	sort.Slice(out, func(i, j int) bool {
		if out[i].Sequence == out[j].Sequence {
			return out[i].ID < out[j].ID
		}
		return out[i].Sequence < out[j].Sequence
	})
	return out, err
}
func (s *Store) SaveProfile(profile domain.CustomerProfile) error {
	return s.db.Update(func(tx *bbolt.Tx) error { return putJSON(tx, profilesBucket, profile.ID, profile) })
}
func (s *Store) Profile(id string) (domain.CustomerProfile, error) {
	var out domain.CustomerProfile
	err := s.db.View(func(tx *bbolt.Tx) error { return getJSON(tx, profilesBucket, id, &out) })
	return out, err
}
func (s *Store) SaveAudit(event domain.AuditEvent) error {
	return s.db.Update(func(tx *bbolt.Tx) error { return putJSON(tx, auditsBucket, event.ID, event) })
}
func (s *Store) AuditsFor(subject string) ([]domain.AuditEvent, error) {
	out := []domain.AuditEvent{}
	err := s.db.View(func(tx *bbolt.Tx) error {
		return scanJSON(tx, auditsBucket, func(data []byte) error {
			var v domain.AuditEvent
			if err := json.Unmarshal(data, &v); err != nil {
				return err
			}
			if v.SubjectID == subject {
				out = append(out, v)
			}
			return nil
		})
	})
	sort.Slice(out, func(i, j int) bool { return out[i].OccurredAt.Before(out[j].OccurredAt) })
	return out, err
}
func (s *Store) SaveNotification(n domain.Notification) error {
	return s.db.Update(func(tx *bbolt.Tx) error { return putJSON(tx, notificationsBucket, n.ID, n) })
}
func (s *Store) Notification(id string) (domain.Notification, error) {
	var out domain.Notification
	err := s.db.View(func(tx *bbolt.Tx) error { return getJSON(tx, notificationsBucket, id, &out) })
	return out, err
}
