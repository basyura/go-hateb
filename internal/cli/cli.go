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

const dateLayout = "20060102"
const defaultTerminalWidth = 80

type bookmarkFetcher interface {
	FetchBookmarks(ctx context.Context, user string) ([]hateb.Bookmark, error)
}

func Run(args []string, stdout, stderr io.Writer, fetcher bookmarkFetcher) int {
	fs := flag.NewFlagSet("hateb", flag.ContinueOnError)
	fs.SetOutput(stderr)

	user := fs.String("user", "", "Hatena user ID")
	fromValue := fs.String("from", "", "From date (yyyyMMdd)")

	if err := fs.Parse(args); err != nil {
		return 2
	}

	if *user == "" {
		fmt.Fprintln(stderr, "--user is required")
		return 2
	}

	var from *time.Time
	if *fromValue != "" {
		parsed, err := time.ParseInLocation(dateLayout, *fromValue, time.Local)
		if err != nil {
			fmt.Fprintln(stderr, "--from must be yyyyMMdd")
			return 2
		}
		from = &parsed
	}

	bookmarks, err := fetcher.FetchBookmarks(context.Background(), *user)
	if err != nil {
		fmt.Fprintf(stderr, "failed to fetch bookmarks: %v\n", err)
		return 1
	}

	if from != nil {
		bookmarks = hateb.FilterBookmarksSince(bookmarks, *from)
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
