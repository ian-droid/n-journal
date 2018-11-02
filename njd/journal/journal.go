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

var db DB

type DB struct {
	DSN string
	Conn *sql.DB //dbconn in use
}

func OpenDB (dsn string) *DB {
	var err error
	db = DB{DSN:dsn}
	db.Conn, err = sql.Open("sqlite3", dsn)
	checkErr(err)
	SaveMessage(fmt.Sprintf("Dateabase '%s' opened for journaling.", dsn))
	return &db
}

func CloseDB() {
	SaveMessage("Closing dateabase .")
	err := db.Conn.Close()
	checkErr(err)
}

func (db DB) GetDiariesByDateRange(startDate string, endDate string) ([]Diary, int) {
	qStr := "SELECT oid, date, content, highlighted FROM diary WHERE date >= '" + startDate + "' and date <= '" + endDate + "' ORDER BY date ASC"
	//fmt.Println(qStr)
	rows, err := db.Conn.Query(qStr)
	checkErr(err)
	var (
		diary Diary
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

func (db DB) GetTransactionsByDateRange(startDate string, endDate string) ([]Transaction, int) {
	qStr := "SELECT * FROM vTransaction WHERE date >= '" + startDate + "' and date <= '" + endDate + "' ORDER BY date ASC"
		rows, err := db.Conn.Query(qStr)
		checkErr(err)
		var (
			transaction Transaction
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

func (diary Diary) Save() {
	if diary.Oid > 0 {
		// Update
		stmt, err := db.Conn.Prepare("UPDATE diary set content = ?, highlighted = ?, updated = strftime('%s', 'now') WHERE oid = ? and date = ?")
		checkErr(err)
		_, err = stmt.Exec(diary.Content, diary.Highlighted, diary.Oid, diary.Date)
		checkErr(err)
		stmt.Close()
		SaveMessage(fmt.Sprintf("Diary %d updated.", diary.Oid))
	} else {
		// INSERT
		stmt, err := db.Conn.Prepare("INSERT INTO diary(date, content, highlighted) VALUES(?,?,?)")
		checkErr(err)
		res, err := stmt.Exec(diary.Date, diary.Content, diary.Highlighted)
		checkErr(err)
		id, err := res.LastInsertId()
		checkErr(err)
		stmt.Close()
		SaveMessage(fmt.Sprintf("New diary %d saved to database.", id))
	}
}

func GetDiary(diary *Diary) {
	var date string
	//fmt.Printf("SELECT date, content, highlighted FROM diary WHERE oid = %d \n", diary.Oid)
	rows, err := db.Conn.Query("SELECT date, content, highlighted FROM diary WHERE oid = ?", diary.Oid)
	checkErr(err)
	for rows.Next() {
		err = rows.Scan(&date, &diary.Content, &diary.Highlighted)
	}
	t, _ := time.Parse(time.RFC3339, date)
	diary.Date = t.Format("2006-01-02")
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

func (transaction Transaction) Save() {
	if transaction.Oid == 0 {
		// INSERT
		stmt, err := db.Conn.Prepare("INSERT INTO transactions (date, item, description, currency, amount, pay, income, payment, bank) VALUES(?,?,?,?,?,?,?,?,?)")
		checkErr(err)
		res, err := stmt.Exec(transaction.Date, transaction.Item, transaction.Description, transaction.Currency, transaction.Amount, transaction.Pay, transaction.Income, transaction.Payment, transaction.Bank)
		checkErr(err)
		id, err := res.LastInsertId()
		checkErr(err)
		stmt.Close()
		SaveMessage(fmt.Sprintf("New transaction %d saved to database.", id))
	}
	// ToDo: Update
}

type Currency struct {
	Id      int
	Name    string
	Prefix  string
	Current bool
}

func GetCurrencies() []Currency {
	rows, err := db.Conn.Query("SELECT id, name, prefix, current from currency WHERE id <> 0")
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

func GetPayments() []Payment {
	rows, err := db.Conn.Query("SELECT id, name, description, priority FROM payment WHERE id <> 0")
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

func GetBanks() []Bank {
	rows, err := db.Conn.Query("SELECT id, name, description, priority from bank WHERE active")
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

func SaveMessage(msg string) {
	now := time.Now()
	fullMsg := fmt.Sprintf("%s: %s", now.Format(time.UnixDate), msg)
	stmt, err := db.Conn.Prepare("INSERT INTO message(added, content) VALUES(strftime('%s', 'now') ,?)")
	checkErr(err)
	_, err = stmt.Exec(fullMsg)
	checkErr(err)
	stmt.Close()
	fmt.Printf("%s\n", fullMsg)
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
