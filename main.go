package main

import (
	"flag"
	"strconv"
	"time"
)

const ELEMENT_TAGS = "h1, h2, h3, h4, h5, h6, p, ul, ol, li, a, span, pre, em, figcaption"
const OUTPUT_DIRECTORY_NAME = "./outputs"

func main() {
	var domain string
	var protocol string
	var maxExplorationDepth uint
	var tags string
	var outputFolderPath string

	flag.StringVar(&domain, "domain", "go-colly.org", "The domain that the site scraper will be restricted to")
	flag.StringVar(&protocol, "protocol", "https", "The URI protocol of the base URL of the provided domain")
	flag.UintVar(&maxExplorationDepth, "maxExplorationDepth", 5, "The maximum recursive depth that the scraper will explore page hyperlinks through")
	flag.StringVar(&tags, "tags", ELEMENT_TAGS, "The comma seperated list of HTML tags/elements that will text-mined, and included in the output JSON file")
	flag.StringVar(&outputFolderPath, "outputFolderPath", OUTPUT_DIRECTORY_NAME, "The relative path to the output directory of the generated JSON file, that contains the scraped data")

	flag.Parse()

	scraper := NewScraper(domain, protocol, maxExplorationDepth, tags, outputFolderPath)
	scraper.AggregateData()

	siteDataFileName := domain + "_" + strconv.FormatInt(time.Now().Unix(), 10) + ".json"
	scraper.SaveScrapedDataToJSON(siteDataFileName)
}
