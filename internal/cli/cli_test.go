package cli

import (
	"bytes"
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	"hateb/internal/hateb"
)

type stubFetcher struct {
	bookmarks []hateb.Bookmark
	err       error
	user      string
}

func (s *stubFetcher) FetchBookmarks(_ context.Context, user string) ([]hateb.Bookmark, error) {
	s.user = user
	return s.bookmarks, s.err
}

func TestRunOutputsBookmarksFromDate(t *testing.T) {
	fetcher := &stubFetcher{
		bookmarks: []hateb.Bookmark{
			{Title: "old", URL: "https://example.com/old", Date: time.Date(2026, 5, 16, 23, 59, 0, 0, time.Local)},
			{Title: "new", URL: "https://example.com/new", Date: time.Date(2026, 5, 17, 0, 0, 0, 0, time.Local)},
		},
	}
	var stdout, stderr bytes.Buffer

	code := Run([]string{"--user", "sample", "--from", "20260517"}, &stdout, &stderr, fetcher)

	if code != 0 {
		t.Fatalf("Run returned %d, want 0; stderr=%q", code, stderr.String())
	}
	if fetcher.user != "sample" {
		t.Fatalf("user = %q, want sample", fetcher.user)
	}
	if got, want := stdout.String(), "2026/05/17 new\n           https://example.com/new\n\n"; got != want {
		t.Fatalf("stdout = %q, want %q", got, want)
	}
}

func TestRunOutputsBookmarksWithoutFrom(t *testing.T) {
	fetcher := &stubFetcher{
		bookmarks: []hateb.Bookmark{
			{Title: "first", URL: "https://example.com/first", Date: time.Date(2026, 5, 16, 0, 0, 0, 0, time.Local)},
			{Title: "second", URL: "https://example.com/second", Date: time.Date(2026, 5, 17, 0, 0, 0, 0, time.Local)},
		},
	}
	var stdout, stderr bytes.Buffer

	code := Run([]string{"--user", "sample"}, &stdout, &stderr, fetcher)

	if code != 0 {
		t.Fatalf("Run returned %d, want 0; stderr=%q", code, stderr.String())
	}
	want := "2026/05/16 first\n           https://example.com/first\n\n2026/05/17 second\n           https://example.com/second\n\n"
	if got := stdout.String(); got != want {
		t.Fatalf("stdout = %q, want %q", got, want)
	}
}

func TestRunRequiresUser(t *testing.T) {
	var stdout, stderr bytes.Buffer

	code := Run(nil, &stdout, &stderr, &stubFetcher{})

	if code != 2 {
		t.Fatalf("Run returned %d, want 2", code)
	}
	if !strings.Contains(stderr.String(), "--user is required") {
		t.Fatalf("stderr = %q, want user error", stderr.String())
	}
}

func TestRunRejectsInvalidFrom(t *testing.T) {
	var stdout, stderr bytes.Buffer

	code := Run([]string{"--user", "sample", "--from", "2026-05-17"}, &stdout, &stderr, &stubFetcher{})

	if code != 2 {
		t.Fatalf("Run returned %d, want 2", code)
	}
	if !strings.Contains(stderr.String(), "--from must be yyyyMMdd") {
		t.Fatalf("stderr = %q, want from error", stderr.String())
	}
}

func TestRunReportsFetchError(t *testing.T) {
	var stdout, stderr bytes.Buffer

	code := Run(
		[]string{"--user", "sample", "--from", "20260517"},
		&stdout,
		&stderr,
		&stubFetcher{err: errors.New("network error")},
	)

	if code != 1 {
		t.Fatalf("Run returned %d, want 1", code)
	}
	if !strings.Contains(stderr.String(), "network error") {
		t.Fatalf("stderr = %q, want fetch error", stderr.String())
	}
}
