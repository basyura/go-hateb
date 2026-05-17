package hateb

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

const defaultBaseURL = "https://b.hatena.ne.jp"

type Client struct {
	httpClient *http.Client
	baseURL    string
}

func NewClient(httpClient *http.Client) *Client {
	if httpClient == nil {
		httpClient = http.DefaultClient
	}
	return &Client{
		httpClient: httpClient,
		baseURL:    defaultBaseURL,
	}
}

func NewClientWithBaseURL(httpClient *http.Client, baseURL string) *Client {
	client := NewClient(httpClient)
	client.baseURL = strings.TrimRight(baseURL, "/")
	return client
}

func FeedURL(baseURL, user string) (string, error) {
	user = strings.TrimSpace(user)
	if user == "" {
		return "", fmt.Errorf("user is required")
	}
	if strings.Contains(user, "/") {
		return "", fmt.Errorf("user must not contain slash")
	}

	return strings.TrimRight(baseURL, "/") + "/" + url.PathEscape(user) + "/rss", nil
}

func (c *Client) FetchBookmarks(ctx context.Context, user string) ([]Bookmark, error) {
	feedURL, err := FeedURL(c.baseURL, user)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, feedURL, nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("unexpected status: %s", resp.Status)
	}

	return ParseRSS(resp.Body)
}
