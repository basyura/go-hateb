package cli

import (
	"testing"
	"time"

	"hateb/internal/hateb"
)

func TestFormatBookmarkTruncatesTitleToWidth(t *testing.T) {
	bookmark := hateb.Bookmark{
		Title: "abcdefghijklmnopqrstuvwxyz",
		URL:   "https://example.com/",
		Date:  time.Date(2026, 5, 17, 0, 0, 0, 0, time.Local),
	}

	got := formatBookmark(bookmark, 20)
	want := "2026/05/17 abcdef...\n           https://example.com/\n\n"

	if got != want {
		t.Fatalf("formatBookmark = %q, want %q", got, want)
	}
}

func TestTruncateTitleCountsWideRunes(t *testing.T) {
	got := truncateTitle("日本語タイトル", 9)
	want := "日本語..."

	if got != want {
		t.Fatalf("truncateTitle = %q, want %q", got, want)
	}
	if width := stringWidth(got); width > 9 {
		t.Fatalf("width = %d, want <= 9", width)
	}
}
