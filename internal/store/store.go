package store

import (
	"encoding/json"
	"errors"
	"fmt"
	"go.etcd.io/bbolt"
	"os"
	"path/filepath"
	"time"
)

var recordsBucket = []byte("beverage_records")
var reviewsBucket = []byte("review_results")
var profilesBucket = []byte("customer_profiles")
var auditsBucket = []byte("audit_events")
var notificationsBucket = []byte("notifications")

type Store struct {
	db   *bbolt.DB
	path string
}

func Open(path string) (*Store, error) {
	if path == "" {
		return nil, errors.New("database path is required")
	}
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return nil, err
	}
	db, err := bbolt.Open(path, 0600, &bbolt.Options{Timeout: time.Second})
	if err != nil {
		return nil, err
	}
	s := &Store{db: db, path: path}
	err = db.Update(func(tx *bbolt.Tx) error {
		for _, name := range [][]byte{recordsBucket, reviewsBucket, profilesBucket, auditsBucket, notificationsBucket} {
			if _, e := tx.CreateBucketIfNotExists(name); e != nil {
				return e
			}
		}
		return nil
	})
	if err != nil {
		db.Close()
		return nil, err
	}
	return s, nil
}
func (s *Store) Close() error {
	if s == nil || s.db == nil {
		return nil
	}
	return s.db.Close()
}
func putJSON(tx *bbolt.Tx, bucket []byte, key string, value any) error {
	if key == "" {
		return errors.New("empty persistence key")
	}
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return tx.Bucket(bucket).Put([]byte(key), data)
}
func getJSON(tx *bbolt.Tx, bucket []byte, key string, dst any) error {
	data := tx.Bucket(bucket).Get([]byte(key))
	if data == nil {
		return fmt.Errorf("%s: %w", key, ErrNotFound)
	}
	return json.Unmarshal(data, dst)
}
func scanJSON(tx *bbolt.Tx, bucket []byte, fn func([]byte) error) error {
	return tx.Bucket(bucket).ForEach(func(_, v []byte) error { return fn(v) })
}

var ErrNotFound = errors.New("not found")
