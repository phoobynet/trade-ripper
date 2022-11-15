package indexes

import (
	"fmt"
	"github.com/gocolly/colly"
	"strings"
)

// ScrapeSP500 scrapes the S&P 500 from https://en.wikipedia.org/wiki/List_of_S%26P_500_companies
func ScrapeSP500() ([]IndexConstituent, error) {
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
						header := strings.TrimSuffix(el.Text, "\n")
						headerMap[header] = headerIndex
					})
				} else {
					// data
					ticker := el.ChildText(fmt.Sprintf("td:nth-child(%d)", headerMap["Symbol"]+1))
					if ticker != "" {
						indexConstituent := IndexConstituent{
							Ticker:          ticker,
							Company:         el.ChildText(fmt.Sprintf("td:nth-child(%d)", headerMap["Security"]+1)),
							GICSSector:      el.ChildText(fmt.Sprintf("td:nth-child(%d)", headerMap["GICS Sector"]+1)),
							GICSSubIndustry: el.ChildText(fmt.Sprintf("td:nth-child(%d)", headerMap["GICS Sub-Industry"]+1)),
						}
						indexConstituents = append(indexConstituents, indexConstituent)
					}
				}
			})
		})
	})

	err := c.Visit("https://en.wikipedia.org/wiki/List_of_S%26P_500_companies")
	if err != nil {
		return nil, err
	}

	return indexConstituents, nil
}
