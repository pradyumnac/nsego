package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	_ "github.com/lib/pq"
)

// const BASEPATH string = "D:/Projects/Active/Gotut/tmp"

// var linkLanding string = "https://www.google.com/"

// FETCHINTERVAL - Frequency of price fetch
var FETCHINTERVAL = 15

// LINKLANDING -  Referrer page to get cookie for ultimate grab
var LINKLANDING string = "https://www.nseindia.com/market-data/live-equity-market"

type (
	Url     string
	UrlList map[string]Url
)

func GetAllEqIndices() UrlList {
	return UrlList{
		"n50":    "https://www.nseindia.com/api/equity-stockIndices?index=NIFTY%2050",
		"nn50":   "https://www.nseindia.com/api/equity-stockIndices?index=NIFTY%2050",
		"nm50":   "https://www.nseindia.com/api/equity-stockIndices?index=NIFTY%2050",
		"nm100":  "https://www.nseindia.com/api/equity-stockIndices?index=NIFTY%2050",
		"nm150":  "https://www.nseindia.com/api/equity-stockIndices?index=NIFTY%2050",
		"ns50":   "https://www.nseindia.com/api/equity-stockIndices?index=NIFTY%2050",
		"ns100":  "https://www.nseindia.com/api/equity-stockIndices?index=NIFTY%2050",
		"ns250":  "https://www.nseindia.com/api/equity-stockIndices?index=NIFTY%2050",
		"nms400": "https://www.nseindia.com/api/equity-stockIndices?index=NIFTY%2050",
		"n100":   "https://www.nseindia.com/api/equity-stockIndices?index=NIFTY%2050",
		"n200":   "https://www.nseindia.com/api/equity-stockIndices?index=NIFTY%2050",
		"n500mc": "https://www.nseindia.com/api/equity-stockIndices?index=NIFTY%2050",
		"nlm250": "https://www.nseindia.com/api/equity-stockIndices?index=NIFTY%2050",
		"nmc250": "https://www.nseindia.com/api/equity-stockIndices?index=NIFTY%2050",
	}
	// more index will be filled in
}

var (
	linkN50  string = "https://www.nseindia.com/api/equity-stockIndices?index=NIFTY%2050"
	linkNN50 string = "https://www.nseindia.com/api/equity-stockIndices?index=NIFTY%20NEXT%2050"
	linkM400 string = "https://www.nseindia.com/api/equity-stockIndices?index=NIFTY%20MIDSMALLCAP%20400"
	linkN100 string = "https://www.nseindia.com/api/equity-stockIndices?index=NIFTY%20100"
)

// GRABLINKNAMES - Default NSE Grablinks Names
var GRABLINKNAMES = [2]string{"M400", "N100"}

// NSEWatchlist Json Struct to get NSE response
type NSEWatchlist struct {
	Name    string `json:"name"`
	Advance struct {
		Declines  string `json:"declines"`
		Advances  string `json:"advances"`
		Unchanged string `json:"unchanged"`
	} `json:"advance"`
	Timestamp string `json:"timestamp"`
	Data      []struct {
		Priority          int     `json:"priority"`
		Symbol            string  `json:"symbol"`
		Identifier        string  `json:"identifier"`
		Open              float64 `json:"open"`
		DayHigh           float64 `json:"dayHigh"`
		DayLow            float64 `json:"dayLow"`
		LastPrice         float64 `json:"lastPrice"`
		PreviousClose     float64 `json:"previousClose"`
		Change            float64 `json:"change"`
		PChange           float64 `json:"pChange"`
		Ffmc              float64 `json:"ffmc"`
		YearHigh          float64 `json:"yearHigh"`
		YearLow           float64 `json:"yearLow"`
		TotalTradedVolume int     `json:"totalTradedVolume"`
		TotalTradedValue  float64 `json:"totalTradedValue"`
		LastUpdateTime    string  `json:"lastUpdateTime"`
		NearWKH           float64 `json:"nearWKH"`
		NearWKL           float64 `json:"nearWKL"`
		PerChange365D     float64 `json:"perChange365d"`
		Date365DAgo       string  `json:"date365dAgo"`
		Chart365DPath     string  `json:"chart365dPath"`
		Date30DAgo        string  `json:"date30dAgo"`
		PerChange30D      float64 `json:"perChange30d"`
		Chart30DPath      string  `json:"chart30dPath"`
		ChartTodayPath    string  `json:"chartTodayPath"`
		Series            string  `json:"series,omitempty"`
		Meta              struct {
			Symbol              string        `json:"symbol"`
			CompanyName         string        `json:"companyName"`
			Industry            string        `json:"industry"`
			ActiveSeries        []string      `json:"activeSeries"`
			DebtSeries          []interface{} `json:"debtSeries"`
			TempSuspendedSeries []string      `json:"tempSuspendedSeries"`
			IsFNOSec            bool          `json:"isFNOSec"`
			IsCASec             bool          `json:"isCASec"`
			IsSLBSec            bool          `json:"isSLBSec"`
			IsDebtSec           bool          `json:"isDebtSec"`
			IsSuspended         bool          `json:"isSuspended"`
			IsETFSec            bool          `json:"isETFSec"`
			IsDelisted          bool          `json:"isDelisted"`
			Isin                string        `json:"isin"`
		} `json:"meta,omitempty"`
	} `json:"data"`
	Metadata struct {
		IndexName         string  `json:"indexName"`
		Open              float64 `json:"open"`
		High              float64 `json:"high"`
		Low               float64 `json:"low"`
		PreviousClose     float64 `json:"previousClose"`
		Last              float64 `json:"last"`
		PercChange        float64 `json:"percChange"`
		Change            float64 `json:"change"`
		TimeVal           string  `json:"timeVal"`
		YearHigh          float64 `json:"yearHigh"`
		YearLow           float64 `json:"yearLow"`
		TotalTradedVolume int     `json:"totalTradedVolume"`
		TotalTradedValue  float64 `json:"totalTradedValue"`
		FfmcSum           float64 `json:"ffmc_sum"`
	} `json:"metadata"`
	MarketStatus struct {
		Market              string  `json:"market"`
		MarketStatus        string  `json:"marketStatus"`
		TradeDate           string  `json:"tradeDate"`
		Index               string  `json:"index"`
		Last                float64 `json:"last"`
		Variation           float64 `json:"variation"`
		PercentChange       float64 `json:"percentChange"`
		MarketStatusMessage string  `json:"marketStatusMessage"`
	} `json:"marketStatus"`
	Date30DAgo  string `json:"date30dAgo"`
	Date365DAgo string `json:"date365dAgo"`
}

var client = &http.Client{
	Timeout: 20 * time.Second,
}

// Database section
func connectPostgres() *sql.DB {
	const (
		dbHost     = "localhost"
		dbPort     = 5432
		dbUser     = "nsefetcher"
		dbPassword = "jw8s0F4"
		dbName     = "nsefetcher"
	)

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		dbHost, dbPort, dbUser, dbPassword, dbName)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	// fmt.Println("Successfully connected!")

	return db
}

func getInitialCookie() []*http.Cookie {
	// log.Println("Running Program: ")

	// Create and modify HTTP request before sending
	request, err := http.NewRequest("GET", LINKLANDING, nil)
	if err != nil {
		log.Fatal(err)
	}
	request.Header.Set("Referer", LINKLANDING)
	request.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/80.0.3987.163 Safari/537.36")

	response, err := client.Do(request)
	if err != nil {
		log.Fatal(err)
	}

	// Get Cookies
	cookies := response.Cookies()
	response.Body.Close()

	return cookies
}

func getNSE(grabLinks [2]string, grabLinksNames [2]string) ([]NSEWatchlist, error) {
	sliceWatchlist := []NSEWatchlist{}
	cookies := getInitialCookie()
	// log.Print(cookies)
	request, err := http.NewRequest("GET", LINKLANDING, nil)
	if err != nil {
		log.Fatal(err)
	}

	for linkIndex := range grabLinks {
		// log.Println("Running for: " + grabLinksNames[linkIndex])
		request, err = http.NewRequest("GET", grabLinks[linkIndex], nil)
		for i := range cookies {
			request.AddCookie(cookies[i])
		}
		request.Header.Set("Accept", "application/json, text/javascript, */*; q=0.01")
		request.Header.Set("Referer", "https://www.nseindia.com/market-data/live-equity-market")
		request.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/80.0.3987.163 Safari/537.36")

		response, err := client.Do(request)
		if err != nil {
			log.Fatal(err)
		}
		defer response.Body.Close()

		if response.StatusCode != http.StatusOK {
			log.Printf("%d", response.StatusCode)
			panic("Request not successful")

		}
		bodyBytes, err := ioutil.ReadAll(response.Body)
		if err != nil {
			log.Fatal(err)
		}

		var watchlist NSEWatchlist
		err2 := json.Unmarshal(bodyBytes, &watchlist)
		if err != nil {
			log.Fatal(err2)
		}

		log.Printf("%+s: %d\n", watchlist.Name, len(watchlist.Data))
		sliceWatchlist = append(sliceWatchlist, watchlist)

		// get latest cookies
		// cookies = response.Cookies()
		// using this messes up the subsequent request - 15/01/2021
		// hence commenting out

	}

	return sliceWatchlist, nil
}

func createTablesIfNotExist(db *sql.DB) {
	createTableTransMetaQuery := `CREATE TABLE IF NOT EXISTS NseTransactionMeta (
		id 				SERIAL PRIMARY KEY NOT NULL,
		wl_name 		varchar NOT NULL,
		advances 		INT NOT NULL,
		declines 		INT NOT NULL,
		unchanged 		INT NOT NULL,
		market_status	BOOLEAN NOT NULL,
		trade_date		TIMESTAMP NOT NULL,
		indexName 		varchar,
		ts TIMESTAMPTZ 	NOT NULL);`

	_, err := db.Exec(createTableTransMetaQuery)
	if err != nil {
		log.Println("Error in creating TransMeta Table ")
		panic(err)
	}
	createTablePriceQuery := `CREATE TABLE IF NOT EXISTS NsePrices (
		id 					SERIAL PRIMARY KEY NOT NULL,
		trans_meta_id 		SERIAL REFERENCES NseTransactionMeta(id),
		sym_name 			varchar NOT NULL,
		identifier 			varchar NOT NULL,
		o 					DECIMAL NOT NULL,
		h 					DECIMAL NOT NULL,
		l 					DECIMAL NOT NULL,
		c 					DECIMAL NOT NULL,
		prev_close 			DECIMAL NOT NULL,
		chg_d 				DECIMAL NOT NULL,
		chg_pct_d			DECIMAL NOT NULL,
		chg_pct_y 			DECIMAL NOT NULL,
		chg_pct_m			DECIMAL NOT NULL,
		week_h_pct			DECIMAL NOT NULL,
		week_l_pct			DECIMAL NOT NULL,
		y_high 				DECIMAL NOT NULL,
		y_low 				DECIMAL NOT NULL,
		day_volume 			DECIMAL NOT NULL,
		day_trade_val 		DECIMAL NOT Null, 
		trade_date 			TIMESTAMP NOT NULL,
		last_update_time 	TIMESTAMPTZ NOT NULL);`

	_, err = db.Exec(createTablePriceQuery)
	if err != nil {
		log.Println("Error in creating TransMeta Table ")
		panic(err)
	}
}

// replaces ? with $no in sql string
func ReplaceSQL(old, searchPattern string) string {
	tmpCount := strings.Count(old, searchPattern)
	for m := 1; m <= tmpCount; m++ {
		old = strings.Replace(old, searchPattern, "$"+strconv.Itoa(m), 1)
	}
	return old
}

func insertPriceData(db *sql.DB, watchlist NSEWatchlist, dbMetaID int) {
	//
	insertQuery := `INSERT INTO NsePrices 
						(trans_meta_id, sym_name, identifier, o, h, l, c, prev_close,
						chg_d, chg_pct_d, y_high, y_low, day_volume, day_trade_val,
						last_update_time, week_h_pct, week_l_pct, chg_pct_y, chg_pct_m, trade_date)
					VALUES`
	vals := []interface{}{}

	for _, data := range watchlist.Data {
		insertQuery += ` (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, 
						  ?, ?, ?, ?, ?, ?, ?, ?, ?, ?),`
		vals = append(vals, dbMetaID, data.Symbol, data.Identifier, data.Open, data.DayHigh, data.DayLow, data.LastPrice, data.PreviousClose,
			data.Change, data.PChange, data.YearHigh, data.YearLow, data.TotalTradedVolume, data.TotalTradedValue,
			data.LastUpdateTime, data.NearWKH, data.NearWKL, data.PerChange365D, data.PerChange365D, watchlist.MarketStatus.TradeDate)
	}

	// trim the last ,
	insertQuery = insertQuery[0 : len(insertQuery)-1]

	// prepare the statement
	insertQuery = ReplaceSQL(insertQuery, "?")
	// log.Println(vals)
	stmt, _ := db.Prepare(insertQuery)

	// format all vals at once
	_, err := stmt.Exec(vals...)
	if err != nil {
		log.Println("Error in inserting to NsePrices table ")

		f, _ := os.Create("insert-vals.log")
		defer f.Close()

		insdatastr := fmt.Sprintf("%v", vals)
		n, _ := f.WriteString(insdatastr)
		log.Println("Wrote to file: insert-vals.log, bytes %d", n)

		panic(err)
	}
}

func insertTransMeta(db *sql.DB, watchlist NSEWatchlist, wlName string) int {
	var mktStatus bool
	if watchlist.MarketStatus.MarketStatus == "Closed" {
		mktStatus = false
	} else {
		mktStatus = true
	}
	// rowsInsert := []interface{}{}
	// fmt.Println(rowsInsert)
	insertQuery := `INSERT INTO NseTransactionMeta (wl_name, advances, declines, unchanged, market_status, trade_date, indexName, ts) 
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id;`

	// format all vals at once
	_, err := db.Exec(insertQuery, wlName, watchlist.Advance.Advances, watchlist.Advance.Declines, watchlist.Advance.Unchanged, mktStatus, watchlist.MarketStatus.TradeDate, watchlist.Name, watchlist.Timestamp)

	var id int

	if err != nil {
		log.Println("Error in creating TransMeta table ")
		panic(err)
	} else {

		err = db.QueryRow("select currval('NseTransactionMeta_id_seq')").Scan(&id)
		if err != nil {
			log.Println("Could not fetch lst inserted id")
			panic(err)
		}
		// log.Println(fmt.Sprintf("Inserted %d id in TransMeta table!", id))
	}

	if err != nil {
		log.Println("Error:", err.Error())
	}

	return id
}

func fetchLatest(db *sql.DB) error {
	insertedIds := []int{}
	sliceWatchlist, err := getNSE(GRABLINKS, GRABLINKNAMES)
	if err != nil {
		log.Fatal(err)
	}

	for i, watchlist := range sliceWatchlist {
		id := insertTransMeta(db, watchlist, watchlist.Name)
		insertedIds = append(insertedIds, id)
		insertPriceData(db, watchlist, insertedIds[i])
	}

	return nil
}

func nseFetchLoop() {
	log.Println("Starting fetcher. It will grab nse feed every " + strconv.Itoa(FETCHINTERVAL) + " seconds")
	db := connectPostgres()
	defer db.Close()
	createTablesIfNotExist(db)
	log.Println("Connected to database.")
	ticker := time.NewTicker(time.Duration(FETCHINTERVAL) * time.Second)

	for true {
		select {
		case t := <-ticker.C:
			// log.Println("Tick at", t)
			start := time.Now()
			log.Println("Starting Fetch at ", t)
			err := fetchLatest(db)
			if err != nil {
				log.Println("Failed to fetch latest")
				log.Fatal(err)
			} else {
				log.Println("Fetch successful!")
			}
			log.Println(time.Since(start))
		}
	}
}

func main() {
	// overwriting fetchg interval here
	FETCHINTERVAL = 30

	// insertQuery :=
	nseFetchLoop()
}
