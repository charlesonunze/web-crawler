package utils

import (
	"fmt"
	"net/http"
	"strings"

	"golang.org/x/net/html"
)

// GetHTMLPage return an HTML node of the URL
func GetHTMLPage(url string) (*html.Node, error) {
	r, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("cannot get page")
	}

	b, err := html.Parse(r.Body)
	if err != nil {
		return nil, fmt.Errorf("cannot parse page")
	}

	return b, err
}

// ExtractLinks returns all the links in an HTML page
func ExtractLinks(links []string, n *html.Node) []string {
	if n.Type == html.ElementNode && n.Data == "a" {
		for _, a := range n.Attr {
			if a.Key == "href" {
				links = append(links, a.Val)
			}
		}
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		links = ExtractLinks(links, c)
	}

	return links
}

// ToWWW adds www to a URL
func ToWWW(url string) string {
	if strings.Contains(url, "www") {
		return url
	}

	return fmt.Sprintf("%s%s", "https://www.", strings.Split(url, "https://")[1])
}

// BelongsToSubdomain checks if a link belongs to a subdomain
func BelongsToSubdomain(url, subdomain string) bool {
	return strings.HasPrefix(url, "/") || strings.HasPrefix(url, subdomain) || strings.HasPrefix(url, ToWWW(subdomain))
}

// RemoveTrailingSlash removes the "/" from a link
func RemoveTrailingSlash(url string) string {
	return strings.TrimSuffix(url, "/")
}

// FormatURL properly formats a link to a specific standard
func FormatURL(url, subdomain string) string {
	subdomain = RemoveTrailingSlash(subdomain)

	if strings.HasPrefix(url, "/") {
		if len(url) == 1 {
			return subdomain
		}

		return fmt.Sprintf("%s%s", subdomain, RemoveTrailingSlash(url))
	}

	return RemoveTrailingSlash(url)
}
