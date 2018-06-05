package journaldb

import (
    "fmt"
    "database/sql"
    _ "github.com/mattn/go-sqlite3"
    "github.com/shopspring/decimal"
    //"encoding/json"
)

type Diary struct {
    Oid int
    Date string
    Content string
    Highlighted bool
    Added int
    Updated int
}

func (diary Diary) Save (dbconn *sql.DB) {
    if diary.Oid == 0 {
      // INSERT
      stmt, err := dbconn.Prepare("INSERT INTO diary(date, content, highlighted) VALUES(?,?,?)")
      checkErr(err)
      res, err := stmt.Exec(diary.Date, diary.Content, diary.Highlighted)
      checkErr(err)
      id, err := res.LastInsertId()
      checkErr(err)
      fmt.Printf("New diary %d saved to database.\n", id)
    }
    // ToDo: Update
}


type Transaction struct {
    Oid int
    Date string
    Item string
    Description string
    Amount decimal.Decimal
    Pay bool
    Income bool
    Direction string
    Currency int
    CurrencyName string
    Payment int
    PaymentName string
    Bank int
    BankName string
    Added int
    Updated int
}

func (transaction Transaction) Save (dbconn *sql.DB) {
    if transaction.Oid == 0 {
      // INSERT
      stmt, err := dbconn.Prepare("INSERT INTO transactions (date, item, description, currency, amount, pay, income, payment, bank) VALUES(?,?,?,?,?,?,?,?,?)")
      checkErr(err)
      res, err := stmt.Exec(transaction.Date, transaction.Item, transaction.Description, transaction.Currency, transaction.Amount, transaction.Pay, transaction.Income, transaction.Payment, transaction.Bank)
      checkErr(err)
      id, err := res.LastInsertId()
      checkErr(err)
      fmt.Printf("New transaction %d saved to database.\n", id)
    }
    // ToDo: Update
}

type Currency struct {
    Id int
    Name string
    Prefix string
}

func GetCurrencies (dbconn *sql.DB) []Currency {
    rows, err := dbconn.Query("SELECT id, name, prefix from currency WHERE id <> 0")
    checkErr(err)
    var currencies []Currency
    for rows.Next() {
      var currency Currency
      err = rows.Scan(&currency.Id, &currency.Name, &currency.Prefix)
      checkErr(err)
      currencies = append(currencies, currency)
    }
    return currencies
}

type Payment struct {
    Id int
    Name string
    Desc string
    Priority bool
}

func GetPayments (dbconn *sql.DB) []Payment {
    rows, err := dbconn.Query("SELECT id, name, description, priority FROM payment WHERE id <> 0")
    checkErr(err)
    var payments []Payment
    for rows.Next() {
      var payment Payment
      err = rows.Scan(&payment.Id, &payment.Name, &payment.Desc, &payment.Priority)
      checkErr(err)
      payments = append(payments, payment)
    }
    return payments
}

type Bank struct {
    Id int
    Name string
    Desc string
    Priority bool
}

func GetBanks (dbconn *sql.DB) []Bank {
    rows, err := dbconn.Query("SELECT id, name, description, priority from bank WHERE active")
    checkErr(err)
    var banks []Bank
    for rows.Next() {
      var bank Bank
      err = rows.Scan(&bank.Id, &bank.Name, &bank.Desc, &bank.Priority)
      checkErr(err)
      banks = append(banks, bank)
    }
    return banks
}


func Open(DSN string) *sql.DB {
    conn, err := sql.Open("sqlite3", DSN)
    checkErr(err)
    fmt.Printf("Dateabase '%s' opened for journaling.\n", DSN)
    return conn
}

func Close(conn *sql.DB) {
    err := conn.Close()
    checkErr(err)
    fmt.Printf("Dateabase closed.\n")
}

func checkErr(err error) {
    if err != nil {
        panic(err)
    }
}
