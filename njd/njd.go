package main

import (
    "fmt"
    "strings"
    "time"
    "io/ioutil"
    "crypto/tls"
    "crypto/x509"
    "net/http"
    "database/sql"
    "html/template"
    "log"
    "github.com/shopspring/decimal"
    "./journaldb"
)


type DiaryForm struct {
    DBConn *sql.DB
    Mode string
    Message string
    Diary []journaldb.Diary
    Diary2Update journaldb.Diary
}

func (diaryForm *DiaryForm) Form(w http.ResponseWriter, r *http.Request) {
    diaryForm.Mode = "New"

    r.ParseForm()

    if r.Method == "POST" {
        nd := journaldb.Diary{}

        if oid, ok := r.Form["oid"]; ok {
          fmt.Sscanf(strings.Join(oid, ""),"%d", &nd.Oid)
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
        nd.Save(diaryForm.DBConn)
        diaryForm.Message = "Diary of " + strings.Join(r.Form["date"], "") + " saved."
    }
    if oid, ok := r.Form["Edit"]; r.Method == "GET" && ok {
        diaryForm.Mode = "Update"
        fmt.Sscanf(strings.Join(oid, ""),"%d", &diaryForm.Diary2Update.Oid)
        journaldb.GetDiary(diaryForm.DBConn, &diaryForm.Diary2Update)
    }

    qStr := "SELECT oid, date, content, highlighted FROM diary WHERE date >= '" + getDateByDays(-7) + "' ORDER BY date ASC"
    //fmt.Println(qStr)
    rows, err := diaryForm.DBConn.Query(qStr)
    checkErr(err)
    diaryForm.Diary = nil
    var diary journaldb.Diary
    for rows.Next() {
      var date string
      err = rows.Scan(&diary.Oid, &date, &diary.Content, &diary.Highlighted)
      checkErr(err)
      t, _ := time.Parse(time.RFC3339, date)
      diary.Date = t.Format("2006-01-02")
      diaryForm.Diary = append(diaryForm.Diary, diary)
    }

    tmpl := template.Must(template.ParseFiles("diary.gtpl"))
    tmpl.Execute(w, diaryForm)
    diaryForm.Diary2Update = journaldb.Diary{}
    diaryForm.Message = ""
}


type TransactionForm struct {
    DBConn *sql.DB
    Mode string
    Transaction []journaldb.Transaction
    Currency []journaldb.Currency
    Payment []journaldb.Payment
    Bank []journaldb.Bank
}

func (transactionForm *TransactionForm) Form(w http.ResponseWriter, r *http.Request) {
    transactionForm.Mode = "Insert"

    r.ParseForm()

    if r.Method == "POST" {
        nt := journaldb.Transaction{}

        if oid, ok := r.Form["oid"]; ok {
          fmt.Sscanf(strings.Join(oid, ""),"%d", &nt.Oid)
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
        fmt.Sscanf(strings.Join(r.Form["currency"], ""),"%d", &nt.Currency)

        nt.Amount, _ = decimal.NewFromString(strings.Join(r.Form["amount"], ""))

        if strings.Join(r.Form["direction"], "") == "pay" {
          nt.Pay = true
          nt.Income = false
        } else {
          nt.Pay = false
          nt.Income = true
        }

        fmt.Sscanf(strings.Join(r.Form["payment"], ""),"%d", &nt.Payment)
        fmt.Sscanf(strings.Join(r.Form["bank"], ""),"%d", &nt.Bank)

        nt.Save(transactionForm.DBConn)
    }

    qStr := "SELECT * FROM vTransaction WHERE date >= '" + getDateByDays(-5) + "' ORDER BY date ASC"
    rows, err := transactionForm.DBConn.Query(qStr)
    checkErr(err)
    transactionForm.Transaction = nil
    var transaction journaldb.Transaction
    for rows.Next() {
      var date,amount string
      err = rows.Scan(&transaction.Oid, &date, &transaction.Item, &transaction.Description, &transaction.Direction, &transaction.CurrencyPrefix, &amount, &transaction.PaymentName, &transaction.BankName)
      checkErr(err)
      t, _ := time.Parse(time.RFC3339, date)
      transaction.Date = t.Format("2006-01-02")
      transaction.Amount, _ = decimal.NewFromString(amount)
      transactionForm.Transaction = append(transactionForm.Transaction, transaction)
    }

    transactionForm.Currency = journaldb.GetCurrencies(transactionForm.DBConn)
    transactionForm.Payment = journaldb.GetPayments(transactionForm.DBConn)
    transactionForm.Bank = journaldb.GetBanks(transactionForm.DBConn)

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
    return n.Format("2006-01-02")
}


func main() {
    dbconn := journaldb.Open("ian_journal.db")
    diaryForm := DiaryForm {DBConn: dbconn}
    transactionForm := TransactionForm {DBConn: dbconn}

    fs := http.FileServer(http.Dir("static"))
    http.Handle("/", fs)
    http.HandleFunc("/diary", diaryForm.Form)
    http.HandleFunc("/transaction", transactionForm.Form)

    caCert, err := ioutil.ReadFile("ca.crt")
    if err != nil {
        log.Fatal(err)
    }

    caCertPool := x509.NewCertPool()
    caCertPool.AppendCertsFromPEM(caCert)

    // Setup HTTPS client
    tlsConfig := &tls.Config {
        ClientCAs: caCertPool,
        // NoClientCert
    		// RequestClientCert
    		// RequireAnyClientCert
    		// VerifyClientCertIfGiven
    		// RequireAndVerifyClientCert
    		ClientAuth: tls.RequireAndVerifyClientCert,
    }
    tlsConfig.BuildNameToCertificate()

    server := &http.Server {
		Addr:      ":8086",
		TLSConfig: tlsConfig,
    }

    server.ListenAndServeTLS("server.crt", "server.key")

}
