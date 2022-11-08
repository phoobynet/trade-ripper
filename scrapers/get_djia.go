package scrapers

import (
	"fmt"
	"github.com/gocolly/colly"
	"strings"
)

// GetDJIA scrapes the DJIA from https://en.wikipedia.org/wiki/Dow_Jones_Industrial_Average (probably not the best source)
func GetDJIA() ([]IndexConstituent, error) {
	c := colly.NewCollector()

	headerMap := make(map[string]int)
	indexConstituents := make([]IndexConstituent, 0)

	// Find and visit all links
	c.OnHTML("table#constituents", func(e *colly.HTMLElement) {
		// parse the columns headers to figure out indexes of the columns we want
		e.ForEach("tbody", func(i int, el *colly.HTMLElement) {
			// header
			el.ForEach("tr", func(rowIndex int, el *colly.HTMLElement) {
				if rowIndex == 0 {
					el.ForEach("th", func(headerIndex int, el *colly.HTMLElement) {
						header := strings.TrimSpace(strings.TrimSuffix(el.Text, "\n"))
						headerMap[header] = headerIndex
					})
				} else {
					// data
					ticker := el.ChildText(fmt.Sprintf("td:nth-child(%d)", headerMap["Symbol"]+1))

					if ticker != "" {
						indexConstituent := IndexConstituent{
							Ticker:     ticker,
							Company:    el.ChildText(fmt.Sprintf("th:nth-child(%d)", headerMap["Company"]+1)),
							GICSSector: el.ChildText(fmt.Sprintf("td:nth-child(%d)", headerMap["Industry"]+1)),
						}
						indexConstituents = append(indexConstituents, indexConstituent)
					}
				}
			})
		})
	})

	err := c.Visit("https://en.wikipedia.org/wiki/Dow_Jones_Industrial_Average")
	if err != nil {
		return nil, err
	}

	return indexConstituents, nil
}
