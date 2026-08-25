package api

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"tea17/internal/catalog"
	"tea17/internal/notify"
	"tea17/internal/service"
	"tea17/internal/store"
	"testing"
)

func TestCreateAndGetRecord(t *testing.T) {
	db, err := store.Open(filepath.Join(t.TempDir(), "api.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	server := New(service.New(db, catalog.New(), notify.NewMemorySender()))
	body := bytes.NewBufferString(`{"id":"api-1","name":"Lemon Jasmine","category":"fruit","ingredients":["jasmine tea","lemon"],"price_cents":500}`)
	request := httptest.NewRequest(http.MethodPost, "/records", body)
	response := httptest.NewRecorder()
	server.Handler().ServeHTTP(response, request)
	if response.Code != http.StatusCreated {
		t.Fatalf("create status %d: %s", response.Code, response.Body.String())
	}
	request = httptest.NewRequest(http.MethodGet, "/records/api-1", nil)
	response = httptest.NewRecorder()
	server.Handler().ServeHTTP(response, request)
	if response.Code != http.StatusOK {
		t.Fatalf("get status %d", response.Code)
	}
}
