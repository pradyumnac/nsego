package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

// const BASEPATH string = "D:/Projects/Active/Gotut/tmp"

// var linkLanding string = "https://www.google.com/"

var linkLanding string = "https://www.nseindia.com/market-data/live-equity-market"

var linkN50 string = "https://www.nseindia.com/api/equity-stockIndices?index=NIFTY%2050"
var linkNn50 string = "https://www.nseindia.com/api/equity-stockIndices?index=NIFTY%20NEXT%2050"
var linkM400 string = "https://www.nseindia.com/api/equity-stockIndices?index=NIFTY%20MIDSMALLCAP%20400"
var linkN100 string = "https://www.nseindia.com/api/equity-stockIndices?index=NIFTY%20100"

var grabLinks = [4]string{linkN50, linkNn50, linkN100, linkM400}
var grabLinksNames = [4]string{"N50", "NN50", "N100", "M400"}

func main() {
	fmt.Println("Running Program: ")

	client := &http.Client{
		Timeout: 20 * time.Second,
	}

	// Create and modify HTTP request before sending
	request, err := http.NewRequest("GET", linkLanding, nil)
	if err != nil {
		log.Fatal(err)
	}
	// request.Header.Set("Accept", "application/json, text/javascript, */*; q=0.01")
	request.Header.Set("Referer", linkLanding)
	// request.Header.Set("Upgrade-Insecure-Requests", "1")
	request.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/80.0.3987.163 Safari/537.36")

	response, err := client.Do(request)
	if err != nil {
		log.Fatal(err)
	}

	// Get Cookies
	cookies := response.Cookies()
	response.Body.Close()
	// fmt.Print(cookies)

	for linkIndex := range grabLinks {
		fmt.Println("Running for: " + grabLinksNames[linkIndex])
		request, err = http.NewRequest("GET", grabLinks[linkIndex], nil)
		for i := range cookies {
			request.AddCookie(cookies[i])
		}
		request.Header.Set("Accept", "application/json, text/javascript, */*; q=0.01")
		request.Header.Set("Referer", "https://www1.nseindia.com/live_market/dynaContent/live_watch/equities_stock_watch.htm")
		request.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/80.0.3987.163 Safari/537.36")

		response, err = client.Do(request)
		if err != nil {
			log.Fatal(err)
		}

		// Copy data from HTTP response to file
		now := time.Now()
		dt := now.Format("02.01.2006")
		ts := now.Format("20060102150405")
		BASEPATH, _ := os.Getwd()
		BASEPATH = BASEPATH + "/export"
		os.MkdirAll(BASEPATH+"/"+dt, os.ModePerm)
		filepath := BASEPATH + "/" + dt + "/" + grabLinksNames[linkIndex] + "-" + ts + ".json"

		outFile, err := os.Create(filepath)
		if err != nil {
			log.Fatal(err)
		}
		_, err = io.Copy(outFile, response.Body)
		if err != nil {
			log.Fatal(err)
		}

		//get latest cookies
		cookies = response.Cookies()
		response.Body.Close()
		outFile.Close()
	}

}
