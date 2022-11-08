package scrapers

import "fmt"

type IndexConstituent struct {
	Ticker          string
	Company         string
	GICSSector      string
	GICSSubIndustry string
}

func (s *IndexConstituent) String() string {
	return fmt.Sprintf("%s %s %s %s\n", s.Ticker, s.Company, s.GICSSector, s.GICSSubIndustry)
}
