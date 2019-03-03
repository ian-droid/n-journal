package main

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/ian-droid/njd/journal"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"
)

const dateFormatShort = "2006-01-02"

var (
	dbFile           = flag.String("db", "journal.db", "SQLite DB file for journaling.")
	serverCertFile   = flag.String("cert", "server.crt", "TLS server certification file.")
	serverKeyFile    = flag.String("key", "server.key", "TLS server key file.")
	clientCaCertFile = flag.String("ca", "ca.crt", "TLS client CA cert file.")
	svrAddress       = flag.String("address", "0.0.0.0", "Listening address, default: 0.0.0.0")
	svrPort          = flag.String("port", "8086", "Listening port, default: 8086.")
)

type Diaries struct {
	Code      int
	Message   string
	StartDate string
	EndDate   string
	DayCount  int
	RowCount  int
	List      []journal.Diary
}

type Transactions struct {
	Code      int
	Message   string
	StartDate string
	EndDate   string
	DayCount  int
	RowCount  int
	List      []journal.Transaction
}

type REST struct {
	JournalDB *journal.DB
}

func (rest *REST) route(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	// fmt.Println("Method:", r.Method)
	// fmt.Println("Path:", r.URL.Path)
	// fmt.Println(r.Form)
	//fmt.Fprintf(w, "REST request routed.")

	var (
		startDate, endDate string
		dayCount           int
	)
	if startDate = strings.Join(r.Form["s_date"], ""); startDate == "" {
		startDate = getDateByDays(-7)
	}
	if endDate = strings.Join(r.Form["e_date"], ""); endDate == "" {
		endDate = getDateByDays(0)
	}
	dayCount = getDaysByDate(startDate, endDate)

	switch r.URL.Path {
	case "/diary":
		switch r.Method {
		case "GET":
			rst := Diaries{
				Code:      0,
				StartDate: startDate,
				EndDate:   endDate,
				DayCount:  dayCount,
			}
			rst.List, rst.RowCount = rest.JournalDB.GetDiariesByDateRange(startDate, endDate)
			j, err := json.Marshal(rst)
			checkErr(err)
			fmt.Fprintf(w, string(j))
		case "POST":
			decoder := json.NewDecoder(r.Body)
			var diary journal.Diary
			err := decoder.Decode(&diary)
			checkErr(err)
			rid := rest.JournalDB.SaveDiary(diary)
			j, err := json.Marshal(rid)
			checkErr(err)
			fmt.Fprintf(w, string(j))
		}
	case "/transaction":
		switch r.Method {
		case "GET":
			rst := Transactions{
				Code:      0,
				StartDate: startDate,
				EndDate:   endDate,
				DayCount:  dayCount,
			}
			rst.List, rst.RowCount = rest.JournalDB.GetTransactionsByDateRange(startDate, endDate)
			j, err := json.Marshal(rst)
			checkErr(err)
			fmt.Fprintf(w, string(j))
		}
	}
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

func getDateByDays(d int) string {
	n := time.Now().AddDate(0, 0, d)
	return n.Format(dateFormatShort)
}

func getDaysByDate(ts1, ts2 string) int {
	t1, _ := time.Parse(dateFormatShort, ts1)
	t2, _ := time.Parse(dateFormatShort, ts2)
	return int(t2.Sub(t1)/(24*time.Hour)) + 1
}

func main() {
	flag.Parse()
	var addrStr = *svrAddress + ":" + *svrPort

	jdb := &journal.DB{
		DBType: "sqlite3",
		DSN:    *dbFile,
	}
	jdb.Open()

	rest := REST{JournalDB: jdb}

	fs := http.FileServer(http.Dir("static"))
	http.Handle("/", fs)
	http.HandleFunc("/diary", rest.route)
	http.HandleFunc("/transaction", rest.route)

	caCert, err := ioutil.ReadFile(*clientCaCertFile)
	if err != nil {
		log.Fatal(err)
	}

	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	// Setup HTTPS client
	tlsConfig := &tls.Config{
		ClientCAs: caCertPool,
		// NoClientCert
		// RequestClientCert
		// RequireAnyClientCert
		// VerifyClientCertIfGiven
		// RequireAndVerifyClientCert
		ClientAuth: tls.RequireAndVerifyClientCert,
	}
	tlsConfig.BuildNameToCertificate()

	server := &http.Server{
		Addr:      addrStr,
		TLSConfig: tlsConfig,
	}

	jdb.SaveMessage(fmt.Sprintf("Serving requests on https://%s/.", addrStr))
	server.ListenAndServeTLS(*serverCertFile, *serverKeyFile)

}
