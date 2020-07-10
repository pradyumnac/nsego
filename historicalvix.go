package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

var client = &http.Client{
	Timeout: 20 * time.Second,
}

var LINKVIX string = "https://www1.nseindia.com/products/dynaContent/equities/indices/hist_vix_data.jsp?=&fromDate=%s&toDate=%s"

type VixRecord struct {
	Date      string
	Open      string
	High      string
	Low       string
	Close     string
	PrevClose string
	Change    string
	ChangePct string
}

// ReadCsv accepts a file and returns its content as a multi-dimentional type
// with lines and each column. Only parses to string type.
func ReadCsv(filename string) ([][]string, error) {

	// Open CSV file
	f, err := os.Open(filename)
	if err != nil {
		return [][]string{}, err
	}
	defer f.Close()

	// Read File into a Variable
	lines, err := csv.NewReader(f).ReadAll()
	if err != nil {
		return [][]string{}, err
	}

	return lines, nil
}

func getHistoricalVix(startDate string, endDate string) []VixRecord {
	//Both date inputs are str objects
	//Returns []{ []{'Date', 'Open', 'High', 'Low', 'Close', 'Prev. Close', 'Change', '% Change'} }

	link := fmt.Sprintf(LINKVIX, startDate, endDate)
	// fmt.Println(link)

	vixRows := []VixRecord{}

	request, err := http.NewRequest("GET", link, nil)
	request.Header.Set("Host", "www1.nseindia.com")
	request.Header.Set("Referer", "https://www1.nseindia.com/products/content/equities/indices/historical_vix.htm")
	request.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:78.0) Gecko/20100101 Firefox/78.0")
	request.Header.Set("X-Requested-With", "XMLHttpRequest")

	res, err := client.Do(request)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
	}

	fmt.Println(res.StatusCode)
	log.Println(res.Body)

	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	// Find the review items
	doc.Find("div#csvContentDiv").Each(func(i int, s *goquery.Selection) {
		// For each item found, get the band and title

		// Read File into a Variable
		csvText := strings.ReplaceAll(s.Text(), ":", "\r\n")
		r := csv.NewReader(strings.NewReader(csvText))
		r.LazyQuotes = true
		lines, err := r.ReadAll()
		if err != nil {
			panic(err)
		}

		// Loop through lines & turn into object
		for _, line := range lines {
			data := VixRecord{
				Date:      line[0],
				Open:      line[1],
				High:      line[2],
				Low:       line[3],
				Close:     line[4],
				PrevClose: line[5],
				Change:    line[6],
				ChangePct: line[7],
			}
			vixRows = append(vixRows, data)
		}

	})

	return vixRows[1:]
}

func main() {
	//overwriting fetchg interval here
	StartDate := "01-JUL-2020"
	EndDate := "10-JUL-2020"

	// insertQuery :=
	vixRows := getHistoricalVix(StartDate, EndDate)

	for index := range vixRows {
		d := vixRows[index]
		fmt.Println(d.Open, d.High, d.Low, d.Close)
	}
}
