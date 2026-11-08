package scrapper

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"golang.org/x/net/html"
)

// Some sites block requests without a browser-like User-Agent.
const userAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/117.0.0.0 Safari/537.36"

const maxChars = 2000

// Python's requests has no default timeout; in Go always set one.
var httpClient = &http.Client{Timeout: 15 * time.Second}

// FetchWebsiteContents returns the page title and visible body text,
// truncated to maxChars — equivalent of fetch_website_contents(url).
func FetchWebsiteContents(url string) (string, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return "", fmt.Errorf("build request: %w", err)
	}
	req.Header.Set("User-Agent", userAgent)

	resp, err := httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("fetch %s: %w", url, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return "", fmt.Errorf("fetch %s: status %d", url, resp.StatusCode)
	}

	// BeautifulSoup(response.content, "html.parser")
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return "", fmt.Errorf("parse html: %w", err)
	}

	// soup.title.string if soup.title else "No title found"
	title := strings.TrimSpace(doc.Find("title").First().Text())
	if title == "" {
		title = "No title found"
	}

	text := ""
	body := doc.Find("body")
	if body.Length() > 0 {
		// for irrelevant in soup.body(["script", "style", "img", "input"]): irrelevant.decompose()
		body.Find("script, style, img, input").Remove()
		// soup.body.get_text(separator="\n", strip=True)
		text = getText(body, "\n")
	}

	return truncateRunes(title+"\n\n"+text, maxChars), nil
}

// getText mimics BeautifulSoup's get_text(separator=sep, strip=True):
// collect every text node, trim it, drop empties, join with sep.
func getText(sel *goquery.Selection, sep string) string {
	var parts []string
	var walk func(*html.Node)
	walk = func(n *html.Node) {
		if n.Type == html.TextNode {
			if t := strings.TrimSpace(n.Data); t != "" {
				parts = append(parts, t)
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			walk(c)
		}
	}
	for _, n := range sel.Nodes {
		walk(n)
	}
	return strings.Join(parts, sep)
}

// truncateRunes slices by characters, not bytes — like Python's s[:n].
// Plain s[:n] in Go cuts bytes and can break multi-byte chars (Hindi, emoji).
func truncateRunes(s string, n int) string {
	r := []rune(s)
	if len(r) > n {
		return string(r[:n])
	}
	return s
}
