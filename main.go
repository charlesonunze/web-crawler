package main

import (
	"flag"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/charlesonunze/monzo/utils"
)

func main() {
	var startingURL string
	flag.StringVar(&startingURL, "url", "https://monzo.com", "a starting url from where the crawler should start crawling")
	flag.Parse()

	visitedLinks := make(map[string]bool)
	wg := &sync.WaitGroup{}
	mtx := &sync.RWMutex{}
	queue := make(chan string)

	if startingURL == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}

	start := time.Now()

	visitedLinks[startingURL] = true

	go func() {
		queue <- startingURL
	}()

	links, err := crawl(startingURL, startingURL)
	if err != nil {
		fmt.Printf("error %v", err)
	}

	wg.Add(1)
	go func() {
		defer wg.Done()
		for _, link := range links {
			queue <- link
		}
	}()

	terminate := false
	for !terminate {
		select {
		case l := <-queue:
			fmt.Println("Found -> ", l)
			fmt.Println("")

			mtx.RLock()
			_, found := visitedLinks[l]
			mtx.RUnlock()

			if !found {
				mtx.Lock()
				visitedLinks[l] = true
				mtx.Unlock()
				links, err := crawl(l, startingURL)
				if err != nil {
					fmt.Printf("error getting page %s %s\n", l, err)
				}

				wg.Add(1)
				go func() {
					defer wg.Done()
					for _, link := range links {
						queue <- link
					}
				}()
			}

		case <-time.After(3 * time.Second):
			fmt.Println("byeeeeee!")
			close(queue)
			terminate = true
		}
	}

	wg.Wait()

	fmt.Println("Links visited ->", len(visitedLinks))
	fmt.Println("Time taken ->", time.Since(start))
}

func crawl(url, subdomain string) ([]string, error) {
	var links []string
	fmt.Println("Visiting -> ", url)
	fmt.Println("")

	page, err := utils.GetHTMLPage(url)
	if err != nil {
		fmt.Printf("error getting page %s %s\n", url, err)
		return links, err
	}

	for _, l := range utils.ExtractLinks(nil, page) {
		if utils.BelongsToSubdomain(l, subdomain) {
			link := utils.FormatURL(l, subdomain)
			links = append(links, link)
		}
	}

	return links, nil
}
