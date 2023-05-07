package main

import (
	"flag"
	"strconv"
	"time"
)

func main() {
	var domain string
	var protocol string
	var maxExplorationDepth uint

	flag.StringVar(&domain, "domain", "go-colly.org", "The domain that the site scraper will be restricted to")
	flag.StringVar(&protocol, "protocol", "https", "The URI protocol of the base URL of the provided domain")
	flag.UintVar(&maxExplorationDepth, "maxExplorationDepth", 5, "The maximum recursive depth that the scraper will explore page hyperlinks through")

	flag.Parse()

	scraper := NewScraper(domain, protocol, maxExplorationDepth)
	scraper.AggregateData()

	siteDataFileName := domain + "_" + strconv.FormatInt(time.Now().Unix(), 10) + ".json"
	scraper.SaveScrapedDataToJSON(siteDataFileName)
}
