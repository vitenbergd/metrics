package metrics

import (
	"fmt"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/shaj13/go-guardian/v2/auth"
	"github.com/shaj13/go-guardian/v2/auth/strategies/token"
	"github.com/shaj13/go-guardian/v2/auth/strategies/union"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestMetricsMiddlewareNoToken(t *testing.T) {
	req, err := http.NewRequest("GET", "/metrics", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	tok := "90d64460d14870c08c81352a05dedd3465940a7"
	strategy := union.New(token.NewStatic(map[string]auth.Info{
		tok: auth.NewDefaultUser("example", "1", nil, nil),
	}))

	handler := middleware(strategy, promhttp.Handler())
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusUnauthorized {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
}

func TestMetricsMiddlewareCorrectToken(t *testing.T) {
	req, err := http.NewRequest("GET", "/metrics", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	tok := "90d64460d14870c08c81352a05dedd3465940a7"
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", tok))
	strategy := union.New(token.NewStatic(map[string]auth.Info{
		tok: auth.NewDefaultUser("example", "1", nil, nil),
	}))

	handler := middleware(strategy, promhttp.Handler())
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
}

func TestMetricsMiddlewareWrongToken(t *testing.T) {
	req, err := http.NewRequest("GET", "/metrics", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	tok := "90d64460d14870c08c81352a05dedd3465940a7"
	wrongTok := "wrong!"
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", wrongTok))
	strategy := union.New(token.NewStatic(map[string]auth.Info{
		tok: auth.NewDefaultUser("example", "1", nil, nil),
	}))

	handler := middleware(strategy, promhttp.Handler())
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusUnauthorized {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
}
