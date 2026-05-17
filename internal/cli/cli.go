package cli

import (
	"context"
	"flag"
	"fmt"
	"io"
	"os"
	"strconv"
	"time"

	"hateb/internal/hateb"
)

const dateLayout = "2006-01-02"
const defaultTerminalWidth = 80

type bookmarkFetcher interface {
	FetchBookmarks(ctx context.Context, user string) ([]hateb.Bookmark, error)
}

func Run(args []string, stdout, stderr io.Writer, fetcher bookmarkFetcher) int {
	fs := flag.NewFlagSet("hateb", flag.ContinueOnError)
	fs.SetOutput(stderr)

	user := fs.String("user", "", "Hatena user ID")
	sinceValue := fs.String("since", "", "Since date (YYYY-MM-DD)")

	if err := fs.Parse(args); err != nil {
		return 2
	}

	if *user == "" {
		fmt.Fprintln(stderr, "--user is required")
		return 2
	}

	var since *time.Time
	if *sinceValue != "" {
		parsed, err := time.ParseInLocation(dateLayout, *sinceValue, time.Local)
		if err != nil {
			fmt.Fprintln(stderr, "--since must be YYYY-MM-DD")
			return 2
		}
		since = &parsed
	}

	bookmarks, err := fetcher.FetchBookmarks(context.Background(), *user)
	if err != nil {
		fmt.Fprintf(stderr, "failed to fetch bookmarks: %v\n", err)
		return 1
	}

	if since != nil {
		bookmarks = hateb.FilterBookmarksSince(bookmarks, *since)
	}

	width := terminalWidth(stdout)
	for _, bookmark := range bookmarks {
		fmt.Fprint(stdout, formatBookmark(bookmark, width))
	}

	return 0
}

func terminalWidth(stdout io.Writer) int {
	if file, ok := stdout.(*os.File); ok {
		if width, ok := fileTerminalWidth(file); ok {
			return width
		}
	}

	if value := os.Getenv("COLUMNS"); value != "" {
		width, err := strconv.Atoi(value)
		if err == nil && width > 0 {
			return width
		}
	}

	return defaultTerminalWidth
}
