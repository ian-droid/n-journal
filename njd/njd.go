package main

import (
	"./journal"
	"crypto/tls"
	"crypto/x509"
	"flag"
	"fmt"
	"github.com/shopspring/decimal"
	"html/template"
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
	svrPort          = flag.String("port", "8086", "Listening port, default: 80.")
)

type DiaryForm struct {
	JournalDB    *journal.DB
	Mode         string
	Message      string
	Diaries      []journal.Diary
	Diary2Update journal.Diary
	StartDate    string
	EndDate      string
	DayCount     int
	RowCount     int
}

func (diaryForm *DiaryForm) Form(w http.ResponseWriter, r *http.Request) {
	diaryForm.Mode = "New"

	r.ParseForm()

	if diaryForm.StartDate = strings.Join(r.Form["s_date"], ""); diaryForm.StartDate == "" {
		diaryForm.StartDate = getDateByDays(-7)
	}
	if diaryForm.EndDate = strings.Join(r.Form["e_date"], ""); diaryForm.EndDate == "" {
		diaryForm.EndDate = getDateByDays(0)
	}
	diaryForm.DayCount = getDaysByDate(diaryForm.StartDate, diaryForm.EndDate)

	if r.Method == "POST" {
		nd := journal.Diary{}

		if oid, ok := r.Form["oid"]; ok {
			fmt.Sscanf(strings.Join(oid, ""), "%d", &nd.Oid)
			fmt.Printf("Update existing diary %d\n", nd.Oid)
		} else {
			nd.Oid = 0
		}
		if nd.Date = strings.Join(r.Form["date"], ""); nd.Date == "" {
			nd.Date = getDateByDays(0)
			fmt.Println("Using system date for diary, please make sure TZ is correct.")
		}
		nd.Content = strings.Join(r.Form["content"], "")
		if strings.Join(r.Form["highlighted"], "") == "on" {
			nd.Highlighted = true
		} else {
			nd.Highlighted = false
		}
		diaryForm.JournalDB.SaveDiary(nd)
		diaryForm.Message = "Diary of " + strings.Join(r.Form["date"], "") + " saved."
	}
	if oid, ok := r.Form["Edit"]; r.Method == "GET" && ok {
		diaryForm.Mode = "Update"
		fmt.Sscanf(strings.Join(oid, ""), "%d", &diaryForm.Diary2Update.Oid)
		diaryForm.JournalDB.GetDiary(&diaryForm.Diary2Update)
	}

	diaryForm.Diaries, diaryForm.RowCount = diaryForm.JournalDB.GetDiariesByDateRange(diaryForm.StartDate, diaryForm.EndDate)

	tmpl := template.Must(template.ParseFiles("diary.gtpl"))
	tmpl.Execute(w, diaryForm)
	diaryForm.Diary2Update = journal.Diary{}
	diaryForm.Message = ""
}

type TransactionForm struct {
	JournalDB    *journal.DB
	Mode         string
	Transactions []journal.Transaction
	Currency     []journal.Currency
	Payment      []journal.Payment
	Bank         []journal.Bank
	StartDate    string
	EndDate      string
	DayCount     int
	RowCount     int
}

func (transactionForm *TransactionForm) Form(w http.ResponseWriter, r *http.Request) {
	transactionForm.Mode = "Insert"

	r.ParseForm()

	if transactionForm.StartDate = strings.Join(r.Form["s_date"], ""); transactionForm.StartDate == "" {
		transactionForm.StartDate = getDateByDays(-7)
	}
	if transactionForm.EndDate = strings.Join(r.Form["e_date"], ""); transactionForm.EndDate == "" {
		transactionForm.EndDate = getDateByDays(0)
	}
	transactionForm.DayCount = getDaysByDate(transactionForm.StartDate, transactionForm.EndDate)

	if r.Method == "POST" {
		nt := journal.Transaction{}

		if oid, ok := r.Form["oid"]; ok {
			fmt.Sscanf(strings.Join(oid, ""), "%d", &nt.Oid)
			fmt.Println("Update existing transaction")
		} else {
			nt.Oid = 0
		}
		if nt.Date = strings.Join(r.Form["date"], ""); nt.Date == "" {
			nt.Date = getDateByDays(0)
			fmt.Println("Using system date for transaction, please make sure TZ is correct.")
		}
		nt.Item = strings.Join(r.Form["item"], "")
		nt.Description = strings.Join(r.Form["description"], "")
		fmt.Sscanf(strings.Join(r.Form["currency"], ""), "%d", &nt.Currency)

		nt.Amount, _ = decimal.NewFromString(strings.Join(r.Form["amount"], ""))

		if strings.Join(r.Form["direction"], "") == "pay" {
			nt.Pay = true
			nt.Income = false
		} else {
			nt.Pay = false
			nt.Income = true
		}

		fmt.Sscanf(strings.Join(r.Form["payment"], ""), "%d", &nt.Payment)
		fmt.Sscanf(strings.Join(r.Form["bank"], ""), "%d", &nt.Bank)

		transactionForm.JournalDB.SaveTransaction(nt)
	}

	transactionForm.Transactions, transactionForm.RowCount = transactionForm.JournalDB.GetTransactionsByDateRange(transactionForm.StartDate, transactionForm.EndDate)

	transactionForm.Currency = transactionForm.JournalDB.GetCurrencies()
	transactionForm.Payment = transactionForm.JournalDB.GetPayments()
	transactionForm.Bank = transactionForm.JournalDB.GetBanks()

	tmpl := template.Must(template.ParseFiles("transaction.gtpl"))
	tmpl.Execute(w, transactionForm)
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
	diaryForm := DiaryForm{JournalDB: jdb}
	transactionForm := TransactionForm{JournalDB: jdb}

	fs := http.FileServer(http.Dir("static"))
	http.Handle("/", fs)
	http.HandleFunc("/diary", diaryForm.Form)
	http.HandleFunc("/transaction", transactionForm.Form)

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
