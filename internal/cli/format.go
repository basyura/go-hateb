package cli

import (
	"fmt"
	"unicode"

	"hateb/internal/hateb"
)

const (
	outputDateLayout = "2006/01/02"
	titlePrefixWidth = len("yyyy/MM/dd ")
	urlIndent        = "           "
	ellipsis         = "..."
)

func formatBookmark(bookmark hateb.Bookmark, width int) string {
	date := bookmark.Date.Format(outputDateLayout)
	title := truncateTitle(bookmark.Title, width-titlePrefixWidth)

	return fmt.Sprintf("%s %s\n%s%s\n\n", date, title, urlIndent, bookmark.URL)
}

func truncateTitle(title string, width int) string {
	if width <= 0 {
		return ""
	}
	if stringWidth(title) <= width {
		return title
	}
	if width <= len(ellipsis) {
		return ellipsis[:width]
	}

	limit := width - len(ellipsis)
	truncated := make([]rune, 0, len(title))
	used := 0
	for _, r := range title {
		runeWidth := displayWidth(r)
		if used+runeWidth > limit {
			break
		}
		truncated = append(truncated, r)
		used += runeWidth
	}

	return string(truncated) + ellipsis
}

func stringWidth(value string) int {
	width := 0
	for _, r := range value {
		width += displayWidth(r)
	}
	return width
}

func displayWidth(r rune) int {
	if r == '\t' {
		return 4
	}
	if r == '\n' || r == '\r' {
		return 0
	}
	if unicode.Is(unicode.Mn, r) || unicode.Is(unicode.Me, r) || unicode.Is(unicode.Cf, r) {
		return 0
	}
	if isWideRune(r) {
		return 2
	}
	return 1
}

func isWideRune(r rune) bool {
	return (r >= 0x1100 && r <= 0x115F) ||
		(r >= 0x2329 && r <= 0x232A) ||
		(r >= 0x2E80 && r <= 0xA4CF) ||
		(r >= 0xAC00 && r <= 0xD7A3) ||
		(r >= 0xF900 && r <= 0xFAFF) ||
		(r >= 0xFE10 && r <= 0xFE19) ||
		(r >= 0xFE30 && r <= 0xFE6F) ||
		(r >= 0xFF00 && r <= 0xFF60) ||
		(r >= 0xFFE0 && r <= 0xFFE6)
}
