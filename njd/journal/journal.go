package journal

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"github.com/shopspring/decimal"
	"time"
	//"encoding/json"
)

const dateFormatShort = "2006-01-02"

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

type DB struct {
	DBType string
	DSN    string
	conn   *sql.DB //dbconn in use
}

func (db *DB) Open() {
	var err error
	db.conn, err = sql.Open(db.DBType, db.DSN)
	checkErr(err)
	db.SaveMessage(fmt.Sprintf("%s dateabase '%s' opened for journaling.", db.DBType, db.DSN))
}

func (db *DB) CloseDB() {
	db.SaveMessage("Closing dateabase .")
	err := db.conn.Close()
	checkErr(err)
}

func (db *DB) SaveMessage(msg string) {
	now := time.Now()
	fullMsg := fmt.Sprintf("%s: %s", now.Format(time.UnixDate), msg)
	stmt, err := db.conn.Prepare("INSERT INTO message(added, content) VALUES(strftime('%s', 'now') ,?)")
	checkErr(err)
	_, err = stmt.Exec(fullMsg)
	checkErr(err)
	stmt.Close()
	fmt.Printf("%s\n", fullMsg)
}

func (db *DB) Write(p []byte) (n int, err error) {
	var msg = string(p)
	db.SaveMessage(msg)
	return len(msg), nil
}

func (db *DB) GetDiariesByDateRange(startDate string, endDate string) ([]Diary, int) {
	qStr := "SELECT oid, date, content, highlighted FROM diary WHERE date >= '" + startDate + "' and date <= '" + endDate + "' ORDER BY date ASC"
	//fmt.Println(qStr)
	rows, err := db.conn.Query(qStr)
	checkErr(err)
	var (
		diary   Diary
		diaries []Diary
	)
	count := 0
	for rows.Next() {
		var date string
		err = rows.Scan(&diary.Oid, &date, &diary.Content, &diary.Highlighted)
		checkErr(err)
		t, _ := time.Parse(time.RFC3339, date)
		diary.Date = t.Format(dateFormatShort)
		diaries = append(diaries, diary)
		count++
	}
	return diaries, count
}

func (db *DB) GetTransactionsByDateRange(startDate string, endDate string) ([]Transaction, int) {
	qStr := "SELECT * FROM vTransaction WHERE date >= '" + startDate + "' and date <= '" + endDate + "' ORDER BY date ASC"
	rows, err := db.conn.Query(qStr)
	checkErr(err)
	var (
		transaction  Transaction
		transactions []Transaction
	)
	count := 0
	for rows.Next() {
		var date, amount string
		err = rows.Scan(&transaction.Oid, &date, &transaction.Item, &transaction.Description, &transaction.Direction, &transaction.CurrencyPrefix, &amount, &transaction.PaymentName, &transaction.BankName)
		checkErr(err)
		t, _ := time.Parse(time.RFC3339, date)
		transaction.Date = t.Format(dateFormatShort)
		transaction.Amount, _ = decimal.NewFromString(amount)
		transactions = append(transactions, transaction)
		count++
	}
	return transactions, count
}

type Diary struct {
	Oid         int
	Date        string
	Content     string
	Highlighted bool
	Added       int
	Updated     int
}

func (db DB) SaveDiary(diary Diary) {
	if diary.Oid > 0 {
		// Update
		stmt, err := db.conn.Prepare("UPDATE diary set content = ?, highlighted = ?, updated = strftime('%s', 'now') WHERE oid = ? and date = ?")
		checkErr(err)
		_, err = stmt.Exec(diary.Content, diary.Highlighted, diary.Oid, diary.Date)
		checkErr(err)
		stmt.Close()
		db.SaveMessage(fmt.Sprintf("Diary %d updated.", diary.Oid))
	} else {
		// INSERT
		stmt, err := db.conn.Prepare("INSERT INTO diary(date, content, highlighted) VALUES(?,?,?)")
		checkErr(err)
		res, err := stmt.Exec(diary.Date, diary.Content, diary.Highlighted)
		checkErr(err)
		id, err := res.LastInsertId()
		checkErr(err)
		stmt.Close()
		db.SaveMessage(fmt.Sprintf("New diary %d saved to database.", id))
	}
}

func (db DB) GetDiary(diary *Diary) {
	var date string
	//fmt.Printf("SELECT date, content, highlighted FROM diary WHERE oid = %d \n", diary.Oid)
	rows, err := db.conn.Query("SELECT date, content, highlighted FROM diary WHERE oid = ?", diary.Oid)
	checkErr(err)
	for rows.Next() {
		err = rows.Scan(&date, &diary.Content, &diary.Highlighted)
	}
	t, _ := time.Parse(time.RFC3339, date)
	diary.Date = t.Format(dateFormatShort)
	checkErr(err)
	rows.Close()
}

type Transaction struct {
	Oid            int
	Date           string
	Item           string
	Description    string
	Amount         decimal.Decimal
	Pay            bool
	Income         bool
	Direction      string
	Currency       int
	CurrencyPrefix string
	Payment        int
	PaymentName    string
	Bank           int
	BankName       string
	Added          int
	Updated        int
}

func (db DB) SaveTransaction(transaction Transaction) {
	if transaction.Oid == 0 {
		// INSERT
		stmt, err := db.conn.Prepare("INSERT INTO transactions (date, item, description, currency, amount, pay, income, payment, bank) VALUES(?,?,?,?,?,?,?,?,?)")
		checkErr(err)
		res, err := stmt.Exec(transaction.Date, transaction.Item, transaction.Description, transaction.Currency, transaction.Amount, transaction.Pay, transaction.Income, transaction.Payment, transaction.Bank)
		checkErr(err)
		id, err := res.LastInsertId()
		checkErr(err)
		stmt.Close()
		db.SaveMessage(fmt.Sprintf("New transaction %d saved to database.", id))
	}
	// ToDo: Update
}

type Currency struct {
	Id      int
	Name    string
	Prefix  string
	Current bool
}

func (db DB) GetCurrencies() []Currency {
	rows, err := db.conn.Query("SELECT id, name, prefix, current from currency WHERE id <> 0")
	checkErr(err)
	var currencies []Currency
	for rows.Next() {
		var currency Currency
		err = rows.Scan(&currency.Id, &currency.Name, &currency.Prefix, &currency.Current)
		checkErr(err)
		currencies = append(currencies, currency)
	}
	rows.Close()
	return currencies
}

type Payment struct {
	Id       int
	Name     string
	Desc     string
	Priority bool
}

func (db DB) GetPayments() []Payment {
	rows, err := db.conn.Query("SELECT id, name, description, priority FROM payment WHERE id <> 0")
	checkErr(err)
	var payments []Payment
	for rows.Next() {
		var payment Payment
		err = rows.Scan(&payment.Id, &payment.Name, &payment.Desc, &payment.Priority)
		checkErr(err)
		payments = append(payments, payment)
	}
	rows.Close()
	return payments
}

type Bank struct {
	Id       int
	Name     string
	Desc     string
	Priority bool
}

func (db DB) GetBanks() []Bank {
	rows, err := db.conn.Query("SELECT id, name, description, priority from bank WHERE active")
	checkErr(err)
	var banks []Bank
	for rows.Next() {
		var bank Bank
		err = rows.Scan(&bank.Id, &bank.Name, &bank.Desc, &bank.Priority)
		checkErr(err)
		banks = append(banks, bank)
	}
	rows.Close()
	return banks
}
