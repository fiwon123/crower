package github

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestLatestTagReadsTheTag(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"tag_name":"v0.9.0"}`))
	}))
	defer server.Close()

	got, err := latestTagFrom(server.URL)
	if err != nil {
		t.Fatalf("latestTagFrom: %v", err)
	}

	if got != "v0.9.0" {
		t.Errorf("tag = %q, want v0.9.0", got)
	}
}

// A non-200 answer is an error, not an empty tag: an unauthenticated rate limit
// returns 403 with a JSON body that has no tag, and "no new version" would be a
// false answer.
func TestLatestTagRejectsANon200Response(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
	}))
	defer server.Close()

	if _, err := latestTagFrom(server.URL); err == nil {
		t.Fatal("a 403 response returned no error")
	}
}

func TestLatestTagRejectsABadBody(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("not json"))
	}))
	defer server.Close()

	if _, err := latestTagFrom(server.URL); err == nil {
		t.Fatal("a body that is not JSON returned no error")
	}
}

func TestLatestTagRejectsAnEmptyTag(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`{"tag_name":""}`))
	}))
	defer server.Close()

	if _, err := latestTagFrom(server.URL); err == nil {
		t.Fatal("an empty tag returned no error")
	}
}
