package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"sync"

	"github.com/gocolly/colly"
	"github.com/gocolly/colly/extensions"
)

var FORBIDDEN_TEXT_CHARACTERS = [2]string{"\n", "\t"}

type RouteDocument struct {
	ElementTag string `json:"elementTag"`
	InnerText  string `json:"innerText"`
}

type Scraper struct {
	Protocol         string
	Domain           string
	BaseURL          string
	OutputFolderPath string
	ScrapableTags    string
	Collector        *colly.Collector
	ContentRouteMap  map[string][]*RouteDocument
	VisitedURLS      map[string]bool
	Mutex            *sync.Mutex
}

func CreateRouteDocument(el *colly.HTMLElement) *RouteDocument {
	finalInnerText := strings.Join(strings.Fields(strings.ToLower(el.DOM.Text())), " ")

	for _, FORBIDDEN_CHAR := range FORBIDDEN_TEXT_CHARACTERS {
		finalInnerText = strings.ReplaceAll(finalInnerText, FORBIDDEN_CHAR, " ")
	}

	return &RouteDocument{
		ElementTag: el.Name,
		InnerText:  finalInnerText,
	}
}

func NormalizeURL(rawURL string) string {
	return strings.TrimSpace(strings.ToLower(rawURL))
}

func (scraper *Scraper) CanVisitURL(url string) bool {

	parsedURL := NormalizeURL(url)
	_, hasAlreadyVisitedParsedURL := scraper.VisitedURLS[parsedURL]

	if !strings.Contains(parsedURL, scraper.Domain) || hasAlreadyVisitedParsedURL {
		return false
	}

	return true
}

func (scraper *Scraper) SaveScrapedDataToJSON(fileName string) {

	defer scraper.Mutex.Unlock()
	scraper.Mutex.Lock()

	data, marshalErr := json.Marshal(scraper.ContentRouteMap)
	if marshalErr != nil {
		panic(marshalErr)
	}

	dirCreateErr := os.MkdirAll(OUTPUT_DIRECTORY_NAME, os.ModePerm)
	if dirCreateErr != nil {
		panic(dirCreateErr)
	}

	filePath := strings.TrimRight(scraper.OutputFolderPath, "/") + "/" + fileName
	writeErr := ioutil.WriteFile(filePath, data, 0644)
	if writeErr != nil {
		panic(writeErr)
	}

	log.Printf("Saved data for the domain: %s to the relative file path: %s", scraper.Domain, filePath)
}

func (scraper *Scraper) AggregateData() {

	scraper.Collector.OnHTML(scraper.ScrapableTags, func(el *colly.HTMLElement) {

		requestURL := NormalizeURL(el.Request.URL.String())

		defer scraper.Mutex.Unlock()
		scraper.Mutex.Lock()

		_, ok := scraper.ContentRouteMap[requestURL]
		if !ok {
			scraper.ContentRouteMap[requestURL] = []*RouteDocument{}
		}

		scraper.ContentRouteMap[requestURL] = append(scraper.ContentRouteMap[requestURL], CreateRouteDocument(el))

	})

	scraper.Collector.OnHTML("a[href]", func(el *colly.HTMLElement) {
		relativeURLPath := el.Attr("href")
		absoluteURLPath := el.Request.AbsoluteURL(relativeURLPath)
		el.Request.Visit(absoluteURLPath)
	})

	scraper.Collector.OnRequest(func(request *colly.Request) {

		defer scraper.Mutex.Unlock()
		scraper.Mutex.Lock()

		requestURL := NormalizeURL(request.URL.String())
		canVisitURL := scraper.CanVisitURL(requestURL)
		if !canVisitURL {
			request.Abort()
			return
		}

		log.Printf("Visiting: %s", requestURL)
	})

	scraper.Collector.OnScraped(func(response *colly.Response) {

		defer scraper.Mutex.Unlock()
		scraper.Mutex.Lock()

		requestURL := NormalizeURL(response.Request.URL.String())
		scraper.VisitedURLS[requestURL] = true
		log.Printf("Finished scraping content for: %s", requestURL)
	})

	scraper.Collector.Visit(scraper.BaseURL)
	scraper.Collector.Wait()

}

func NewScraper(domain string, protocol string, maxExplorationDepth uint, scrapableTags string, outputFolderPath string) *Scraper {

	collector := colly.NewCollector()
	collector.AllowedDomains = []string{domain}
	collector.AllowURLRevisit = false
	collector.MaxDepth = int(maxExplorationDepth)
	collector.Async = true

	extensions.RandomUserAgent(collector)
	extensions.Referer(collector)

	baseURL := protocol + "://" + domain + "/"

	return &Scraper{
		Protocol:         protocol,
		Domain:           domain,
		BaseURL:          baseURL,
		ScrapableTags:    scrapableTags,
		OutputFolderPath: outputFolderPath,
		Collector:        collector,
		ContentRouteMap:  make(map[string][]*RouteDocument),
		VisitedURLS:      make(map[string]bool),
		Mutex:            &sync.Mutex{},
	}
}
