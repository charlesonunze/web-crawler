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

	if startingURL == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}

	start := time.Now()

	visitedLinks[startingURL] = true

	wg.Add(1)
	go crawl(startingURL, startingURL, &visitedLinks, wg, mtx)
	wg.Wait()

	fmt.Println("Links visited ->", len(visitedLinks))
	fmt.Println("Time taken ->", time.Since(start))
}

func crawl(startingURL, currentLink string, visitedLinks *map[string]bool, wg *sync.WaitGroup, mtx *sync.RWMutex) {
	defer wg.Done()

	fmt.Println("")
	fmt.Println("Visiting -> ", currentLink)
	fmt.Println("")

	page, err := utils.GetHTMLPage(currentLink)
	if err != nil {
		// TODO implement a retry for failed links
		fmt.Printf("error getting html page %s %s\n", currentLink, err)
		return
	}

	wg2 := new(sync.WaitGroup)

	links := utils.ExtractLinks(nil, page)

	for _, l := range links {
		if utils.BelongsToSubdomain(l, startingURL) {
			link := utils.FormatURL(l, currentLink)
			fmt.Println("Found ->", link)
			fmt.Println("")

			mtx.RLock()
			_, found := (*visitedLinks)[link]
			mtx.RUnlock()

			if !found {
				mtx.Lock()
				(*visitedLinks)[link] = true
				mtx.Unlock()

				wg2.Add(1)
				go crawl(startingURL, link, visitedLinks, wg2, mtx)
			}
		}
	}
	wg2.Wait()

	return
}
