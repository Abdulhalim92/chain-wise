package middleware

import (
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRequestID_SetsHeader(t *testing.T) {
	var capturedID string
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedID = r.Context().Value(RequestIDKey).(string)
		w.WriteHeader(http.StatusOK)
	})
	handler := RequestID(next)
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)
	if rec.Header().Get("X-Request-ID") == "" {
		t.Error("X-Request-ID header not set")
	}
	if capturedID != rec.Header().Get("X-Request-ID") {
		t.Error("context request_id != header")
	}
}

func TestRequestID_PropagatesExisting(t *testing.T) {
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := r.Context().Value(RequestIDKey).(string)
		if id != "custom-id" {
			t.Errorf("want custom-id, got %q", id)
		}
		w.WriteHeader(http.StatusOK)
	})
	handler := RequestID(next)
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("X-Request-ID", "custom-id")
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)
	if rec.Header().Get("X-Request-ID") != "custom-id" {
		t.Errorf("want custom-id, got %q", rec.Header().Get("X-Request-ID"))
	}
}

func TestRecovery_Returns500OnPanic(t *testing.T) {
	log := slog.Default()
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		panic("test panic")
	})
	handler := Recovery(log, next)
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req = req.WithContext(req.Context())
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)
	if rec.Code != http.StatusInternalServerError {
		t.Errorf("status = %d, want 500", rec.Code)
	}
}

func TestLogging_DoesNotPanic(t *testing.T) {
	log := slog.Default()
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	handler := Logging(log, next)
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Errorf("status = %d, want 200", rec.Code)
	}
}
