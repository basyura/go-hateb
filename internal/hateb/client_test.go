package hateb

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestFeedURL(t *testing.T) {
	got, err := FeedURL("https://b.hatena.ne.jp/", "sample-user")
	if err != nil {
		t.Fatalf("FeedURL returned error: %v", err)
	}
	if want := "https://b.hatena.ne.jp/sample-user/rss"; got != want {
		t.Fatalf("FeedURL = %q, want %q", got, want)
	}
}

func TestFetchBookmarks(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/sample/rss" {
			t.Fatalf("path = %q, want /sample/rss", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/rss+xml")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`<?xml version="1.0"?>
<rdf:RDF xmlns:rdf="http://www.w3.org/1999/02/22-rdf-syntax-ns#"
         xmlns:dc="http://purl.org/dc/elements/1.1/">
  <item>
    <title>Example</title>
    <link>https://example.com/</link>
    <dc:date>2026-05-17T12:34:56+09:00</dc:date>
  </item>
</rdf:RDF>`))
	}))
	defer server.Close()

	client := NewClientWithBaseURL(server.Client(), server.URL)
	bookmarks, err := client.FetchBookmarks(context.Background(), "sample")
	if err != nil {
		t.Fatalf("FetchBookmarks returned error: %v", err)
	}
	if len(bookmarks) != 1 {
		t.Fatalf("len(bookmarks) = %d, want 1", len(bookmarks))
	}
	if bookmarks[0].Title != "Example" || bookmarks[0].URL != "https://example.com/" {
		t.Fatalf("bookmark = %#v", bookmarks[0])
	}
}

func TestFetchBookmarksReturnsHTTPError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		http.Error(w, "not found", http.StatusNotFound)
	}))
	defer server.Close()

	client := NewClientWithBaseURL(server.Client(), server.URL)
	if _, err := client.FetchBookmarks(context.Background(), "sample"); err == nil {
		t.Fatal("FetchBookmarks returned nil error, want HTTP error")
	}
}
