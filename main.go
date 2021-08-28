package main

import (
	"flag"
	"fmt"
	"net/http"
	"strings"
	"sync"

	"golang.org/x/net/html"
)

func main() {
	var startingURL string
	flag.StringVar(&startingURL, "url", "https://monzo.com", "a starting url from where the crawler should start crawling")
	flag.Parse()

	visitedLinks := make(map[string]bool)
	wg := &sync.WaitGroup{}
	mtx := &sync.RWMutex{}

	url := removeTrailingSlash(startingURL)

	wg.Add(1)
	defer wg.Wait()
	go crawl(url, &visitedLinks, wg, mtx)
}

func crawl(url string, visitedLinks *map[string]bool, wg *sync.WaitGroup, mtx *sync.RWMutex) {
	defer wg.Done()
	fmt.Println("Visiting -> ", url)
	fmt.Println("")

	page, err := getHTMLPage(url)
	if err != nil {
		fmt.Printf("error getting page %s %s\n", url, err)
		return
	}

	mtx.Lock()
	(*visitedLinks)[url] = true
	mtx.Unlock()

	links := extractLinks(nil, page)
	// fmt.Printf("List of links inside -> %s \n", url)
	// fmt.Println("")

	for _, l := range links {
		// fmt.Printf("url -> %+v \n", l)
		// fmt.Println("")

		// crawl individual links in the same subdomain
		if belongsToSubdomain(l, url) {
			link := formatURL(l, url)

			mtx.RLock()
			visited := (*visitedLinks)[link]
			mtx.RUnlock()

			// prevent crawling a URL twice
			if !visited {
				wg.Add(1)
				go crawl(link, visitedLinks, wg, mtx)
			}
		}
	}

	return
}

func getHTMLPage(url string) (*html.Node, error) {
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

func extractLinks(links []string, n *html.Node) []string {
	if n.Type == html.ElementNode && n.Data == "a" {
		for _, a := range n.Attr {
			if a.Key == "href" {
				links = append(links, a.Val)
			}
		}
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		links = extractLinks(links, c)
	}

	return links
}

func toWWW(url string) string {
	if strings.Contains(url, "www") {
		return url
	}

	return fmt.Sprintf("%s%s", "https://www.", strings.Split(url, "https://")[1])
}

func belongsToSubdomain(url, subdomain string) bool {
	return strings.HasPrefix(url, "/") || strings.HasPrefix(url, subdomain) || strings.HasPrefix(url, toWWW(subdomain))
}

func removeTrailingSlash(url string) string {
	return strings.TrimSuffix(url, "/")
}

func formatURL(url, subdomain string) string {
	subdomain = removeTrailingSlash(subdomain)

	if strings.HasPrefix(url, "/") {
		if len(url) == 1 {
			return subdomain
		}

		return fmt.Sprintf("%s%s", subdomain, removeTrailingSlash(url))
	}

	return removeTrailingSlash(url)
}
