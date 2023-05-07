package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"sync"

	"github.com/gocolly/colly"
)

const ELEMENT_TYPES = "h1, h2, h3, h4, h5, h6, p, ul, ol, span, pre"
const OUTPUT_DIRECTORY_NAME = "./outputs"

var FORBIDDEN_TEXT_CHARACTERS = [2]string{"\n", "\t"}

type RouteDocument struct {
	ElementTag string `json:"elementTag"`
	InnerText  string `json:"innerText"`
}

type Scraper struct {
	protocol        string
	domain          string
	collector       *colly.Collector
	contentRouteMap map[string][]*RouteDocument
	mutex           *sync.Mutex
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

func (scraper *Scraper) SaveScrapedDataToJSON(fileName string) bool {

	defer scraper.mutex.Unlock()
	scraper.mutex.Lock()

	data, marshalErr := json.Marshal(scraper.contentRouteMap)
	if marshalErr != nil {
		panic(marshalErr)
	}

	dirCreateErr := os.MkdirAll(OUTPUT_DIRECTORY_NAME, os.ModePerm)
	if dirCreateErr != nil {
		panic(dirCreateErr)
	}

	filePath := OUTPUT_DIRECTORY_NAME + "/" + fileName
	writeErr := ioutil.WriteFile(filePath, data, 0644)
	if writeErr != nil {
		panic(writeErr)
	}

	return true
}

func (scraper *Scraper) AggregateData() {

	scraper.collector.OnHTML(ELEMENT_TYPES, func(el *colly.HTMLElement) {

		relativeURLPath := el.Request.URL.String()

		defer scraper.mutex.Unlock()
		scraper.mutex.Lock()

		_, ok := scraper.contentRouteMap[relativeURLPath]
		if !ok {
			scraper.contentRouteMap[relativeURLPath] = []*RouteDocument{}
		}

		scraper.contentRouteMap[relativeURLPath] = append(scraper.contentRouteMap[relativeURLPath], CreateRouteDocument(el))

	})

	scraper.collector.OnHTML("a", func(el *colly.HTMLElement) {
		relativeURLPath := el.Attr("href")
		el.Request.Visit(el.Request.AbsoluteURL(relativeURLPath))
	})

	scraper.collector.OnRequest(func(r *colly.Request) {
		log.Println("Visiting:", r.URL.String())
	})

	scraper.collector.OnScraped(func(r *colly.Response) {
		log.Println("Finished scraping content for the url:", r.Request.URL.String())
	})

	baseURL := scraper.protocol + "://" + scraper.domain
	visitErr := scraper.collector.Visit(baseURL)
	if visitErr != nil {
		panic(visitErr)
	}

	scraper.collector.Wait()

}

func NewScraper(domain string, protocol string, maxExplorationDepth uint) *Scraper {

	collector := colly.NewCollector(colly.AllowedDomains(domain), colly.MaxDepth(int(maxExplorationDepth)), colly.Async(true))

	return &Scraper{
		protocol:        protocol,
		domain:          domain,
		collector:       collector,
		contentRouteMap: make(map[string][]*RouteDocument),
		mutex:           &sync.Mutex{},
	}
}
