package main

import (
	"flag"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/charlesonunze/monzo/utils"
)

var waitChan = make(chan struct{}, 1000)

func main() {
	var startingURL string
	flag.StringVar(&startingURL, "url", "https://monzo.com", "a starting url from where the crawler should start crawling")
	flag.Parse()

	if startingURL == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}

	visitedLinks := make(map[string]bool)
	mtx := &sync.RWMutex{}

	start := time.Now()

	visitedLinks[startingURL] = true

	crawl(startingURL, startingURL, mtx, visitedLinks)

	fmt.Println("Links visited ->", len(visitedLinks))
	fmt.Println("Time taken ->", time.Since(start))
}

func crawl(url, subdomain string, mtx *sync.RWMutex, visitedLinks map[string]bool) {
	fmt.Println("Visiting -> ", url)
	fmt.Println("")

	page, err := utils.GetHTMLPage(url)
	if err != nil {
		fmt.Printf("error getting page %s %s\n", url, err)
		return
	}

	wg := sync.WaitGroup{}

	for _, l := range utils.ExtractLinks(nil, page) {
		if utils.BelongsToSubdomain(l, subdomain) {
			link := utils.FormatURL(l, subdomain)

			// fmt.Println("Found -> ", link)
			// fmt.Println("")

			mtx.RLock()
			_, found := visitedLinks[link]
			mtx.RUnlock()

			if !found {
				mtx.Lock()
				visitedLinks[link] = true
				mtx.Unlock()

				select {
				case waitChan <- struct{}{}:
					wg.Add(1)
					go func() {
						defer wg.Done()
						crawl(link, subdomain, mtx, visitedLinks)
						<-waitChan
					}()

				default:
					wg.Add(1)
					defer wg.Done()
					crawl(link, subdomain, mtx, visitedLinks)
				}
			}
		}
	}

	wg.Wait()
	return
}
