package hateb

import (
	"strings"
	"testing"
	"time"
)

func TestParseRSS(t *testing.T) {
	bookmarks, err := ParseRSS(strings.NewReader(`<?xml version="1.0"?>
<rdf:RDF xmlns:rdf="http://www.w3.org/1999/02/22-rdf-syntax-ns#"
         xmlns:dc="http://purl.org/dc/elements/1.1/">
  <item>
    <title>First</title>
    <link>https://example.com/first</link>
    <dc:date>2026-05-17T12:34:56+09:00</dc:date>
  </item>
  <item>
    <title>Second</title>
    <link>https://example.com/second</link>
    <dc:date>2026-05-18T00:00:00+09:00</dc:date>
  </item>
</rdf:RDF>`))
	if err != nil {
		t.Fatalf("ParseRSS returned error: %v", err)
	}
	if len(bookmarks) != 2 {
		t.Fatalf("len(bookmarks) = %d, want 2", len(bookmarks))
	}
	if bookmarks[0].Title != "First" {
		t.Fatalf("first title = %q, want First", bookmarks[0].Title)
	}
}

func TestFilterBookmarksSince(t *testing.T) {
	since := time.Date(2026, 5, 17, 0, 0, 0, 0, time.UTC)
	bookmarks := []Bookmark{
		{Title: "old", Date: since.Add(-time.Nanosecond)},
		{Title: "same", Date: since},
		{Title: "new", Date: since.Add(time.Hour)},
	}

	filtered := FilterBookmarksSince(bookmarks, since)

	if len(filtered) != 2 {
		t.Fatalf("len(filtered) = %d, want 2", len(filtered))
	}
	if filtered[0].Title != "same" || filtered[1].Title != "new" {
		t.Fatalf("filtered = %#v", filtered)
	}
}
