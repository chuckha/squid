package main

import (
	"log"
	"github.com/ChuckHa/Squid/robots"
	"github.com/ChuckHa/Squid/crawler"
	"github.com/ChuckHa/Squid/db"
	"strings"
)

const (
	userAgent = "Squidbot"
	maxRequests = 2
)

var (
	seedURL = "http://www.michaelnielsen.org/ddi/how-to-crawl-a-quarter-billion-webpages-in-40-hours/"
	sem = make(chan int, maxRequests)
	urlc = make(chan string)
)

func main () {
	// Load the throttler
	for i := 0; i < maxRequests; i ++ {
		sem <- 1
	}

	// seed URL
	go func () {
		urlc <- seedURL
	}()

	// When there is an avaialble slot, run a new bot
	for url := range urlc {
		<-sem
		go func() {
			spawnBot(url)
		}()
	}
}

// Make a new bot
func spawnBot(url string) {
	// Get the robots.txt url
	robot := robots.NewRobotsTxtFromUrl(url)
	// Make sure we're allowed to crawl this site
	if robot.Allowed(userAgent, url) {
		// This builds the crawler for this particular url
		crawler := crawler.NewCrawler(userAgent, url)
		// Get the actual data
		metadata, err := crawler.Crawl()
		// Now that we've finished dowloading the site, we can unlock a slot
		sem <- 1
		// Some error handling, just to get a sense of what kind of errors we hit
		if err != nil {
			if !strings.HasPrefix(err.Error(), "Already visited") {
				log.Printf("Error crawling page %s", err)
			}
		} else {
			log.Printf("Saving data from %s", url)
			// Let's save the data
			metadata.Save()
			// And for each link we found, we put the URL on the channel.
			for _, link := range metadata.Links {

				// only put values on that we haven't seen before
				if !db.Exists(link) {
					urlc <- link
				}
			}
		}
	} else {
		sem <- 1
		log.Println("Url is not crawlable.")
	}
}
