package main

import (
	"./journaldb"
	"crypto/tls"
	"crypto/x509"
	"database/sql"
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
	DBConn       *sql.DB
	Mode         string
	Message      string
	Diary        []journaldb.Diary
	Diary2Update journaldb.Diary
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
		nd := journaldb.Diary{}

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
		nd.Save()
		diaryForm.Message = "Diary of " + strings.Join(r.Form["date"], "") + " saved."
	}
	if oid, ok := r.Form["Edit"]; r.Method == "GET" && ok {
		diaryForm.Mode = "Update"
		fmt.Sscanf(strings.Join(oid, ""), "%d", &diaryForm.Diary2Update.Oid)
		journaldb.GetDiary(&diaryForm.Diary2Update)
	}

	qStr := "SELECT oid, date, content, highlighted FROM diary WHERE date >= '" + diaryForm.StartDate + "' and date <= '" + diaryForm.EndDate + "' ORDER BY date ASC"
	//fmt.Println(qStr)
	rows, err := diaryForm.DBConn.Query(qStr)
	checkErr(err)
	diaryForm.Diary = nil
	diaryForm.RowCount = 0
	var diary journaldb.Diary
	for rows.Next() {
		var date string
		err = rows.Scan(&diary.Oid, &date, &diary.Content, &diary.Highlighted)
		checkErr(err)
		t, _ := time.Parse(time.RFC3339, date)
		diary.Date = t.Format(dateFormatShort)
		diaryForm.Diary = append(diaryForm.Diary, diary)
		diaryForm.RowCount++
	}

	tmpl := template.Must(template.ParseFiles("diary.gtpl"))
	tmpl.Execute(w, diaryForm)
	diaryForm.Diary2Update = journaldb.Diary{}
	diaryForm.Message = ""
}

type TransactionForm struct {
	DBConn      *sql.DB
	Mode        string
	Transaction []journaldb.Transaction
	Currency    []journaldb.Currency
	Payment     []journaldb.Payment
	Bank        []journaldb.Bank
	StartDate   string
	EndDate     string
	DayCount    int
	RowCount    int
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
		nt := journaldb.Transaction{}

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

		nt.Save()
	}

	qStr := "SELECT * FROM vTransaction WHERE date >= '" + transactionForm.StartDate + "' and date <= '" + transactionForm.EndDate + "' ORDER BY date ASC"
	rows, err := transactionForm.DBConn.Query(qStr)
	checkErr(err)
	transactionForm.Transaction = nil
	transactionForm.RowCount = 0
	var transaction journaldb.Transaction
	for rows.Next() {
		var date, amount string
		err = rows.Scan(&transaction.Oid, &date, &transaction.Item, &transaction.Description, &transaction.Direction, &transaction.CurrencyPrefix, &amount, &transaction.PaymentName, &transaction.BankName)
		checkErr(err)
		t, _ := time.Parse(time.RFC3339, date)
		transaction.Date = t.Format(dateFormatShort)
		transaction.Amount, _ = decimal.NewFromString(amount)
		transactionForm.Transaction = append(transactionForm.Transaction, transaction)
		transactionForm.RowCount++
	}

	transactionForm.Currency = journaldb.GetCurrencies()
	transactionForm.Payment = journaldb.GetPayments()
	transactionForm.Bank = journaldb.GetBanks()

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

	dbconn := journaldb.Open(*dbFile)
	diaryForm := DiaryForm{DBConn: dbconn}
	transactionForm := TransactionForm{DBConn: dbconn}

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

	journaldb.SaveMessage(fmt.Sprintf("Serving requests on https://%s/.", addrStr))
	server.ListenAndServeTLS(*serverCertFile, *serverKeyFile)

}
