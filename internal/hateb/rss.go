package hateb

import (
	"encoding/xml"
	"fmt"
	"io"
	"time"
)

type Bookmark struct {
	Title string
	URL   string
	Date  time.Time
}

type rssFeed struct {
	Items []rssItem `xml:"item"`
}

type rssItem struct {
	Title string `xml:"title"`
	Link  string `xml:"link"`
	Date  string `xml:"date"`
}

func ParseRSS(reader io.Reader) ([]Bookmark, error) {
	var feed rssFeed
	if err := xml.NewDecoder(reader).Decode(&feed); err != nil {
		return nil, err
	}

	bookmarks := make([]Bookmark, 0, len(feed.Items))
	for _, item := range feed.Items {
		createdAt, err := parseBookmarkDate(item.Date)
		if err != nil {
			return nil, fmt.Errorf("parse bookmark date %q: %w", item.Date, err)
		}
		bookmarks = append(bookmarks, Bookmark{
			Title: item.Title,
			URL:   item.Link,
			Date:  createdAt,
		})
	}

	return bookmarks, nil
}

func FilterBookmarksSince(bookmarks []Bookmark, since time.Time) []Bookmark {
	filtered := make([]Bookmark, 0, len(bookmarks))
	for _, bookmark := range bookmarks {
		if !bookmark.Date.Before(since) {
			filtered = append(filtered, bookmark)
		}
	}
	return filtered
}

func parseBookmarkDate(value string) (time.Time, error) {
	if value == "" {
		return time.Time{}, fmt.Errorf("date is empty")
	}

	layouts := []string{
		time.RFC3339,
		time.RFC1123Z,
		time.RFC1123,
		"2006-01-02",
	}
	for _, layout := range layouts {
		parsed, err := time.Parse(layout, value)
		if err == nil {
			return parsed, nil
		}
	}

	return time.Time{}, fmt.Errorf("unsupported date format")
}
