package api_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"microapp/pkg/api"
)

func TestAPI_RootHandler(t *testing.T) {
	appAPI := api.New()

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rr := httptest.NewRecorder()

	appAPI.Router.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rr.Code)
	}

	body := rr.Body.String()
	if !strings.Contains(body, "Microapp на машине") {
		t.Errorf("unexpected body: %s", body)
	}
}

func TestAPI_Healthz(t *testing.T) {
	appAPI := api.New()

	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	rr := httptest.NewRecorder()

	appAPI.Router.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rr.Code)
	}

	if rr.Body.String() != "OK" {
		t.Errorf("expected OK, got %s", rr.Body.String())
	}
}

func TestAPI_Reviews(t *testing.T) {
	appAPI := api.New()

	req := httptest.NewRequest(http.MethodGet, "/api/v1/reviews", nil)
	rr := httptest.NewRecorder()

	appAPI.Router.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rr.Code)
	}

	if rr.Header().Get("Content-Type") != "application/json; charset=utf-8" {
		t.Errorf("unexpected content type: %s", rr.Header().Get("Content-Type"))
	}
}

func BenchmarkAPI_RootHandler(b *testing.B) {
	appAPI := api.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		rr := httptest.NewRecorder()
		appAPI.Router.ServeHTTP(rr, req)
	}
}
