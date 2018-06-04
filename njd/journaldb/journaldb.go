package journaldb

import (
    "fmt"
    "net/http"
    "database/sql"
    _ "github.com/mattn/go-sqlite3"
    "encoding/json"
)

/*
type Currency struct {
    Id int
    Prefix string
}

type Payment struct {
    Id int
    Name string
}
*/

type Diary struct {
    Oid int
    Date string
    Content string
    Highlighted bool
    Added int
    Updated int
}


type Bank struct {
    Id int
    Name string
    Desc string
}

type Banks struct {
    DBConn *sql.DB
    Bank []Bank
}

func (diary Diary) Save (dbconn *sql.DB) {
    if diary.Oid == 0 {
      // INSERT
      stmt, err := dbconn.Prepare("INSERT INTO diary(date, content, highlighted) values(?,?,?)")
      checkErr(err)
      res, err := stmt.Exec(diary.Date, diary.Content, diary.Highlighted)
      checkErr(err)
      id, err := res.LastInsertId()
      checkErr(err)
      fmt.Printf("New diary %d saved to database.\n", id)
    }
    // ToDo: Update
}


func (banks *Banks) List(w http.ResponseWriter, r *http.Request) {
    //conn, err := sql.Open("sqlite3", "ian_journal.db")
    //checkErr(err)
    rows, err := banks.DBConn.Query("select id, name, description from bank")
    checkErr(err)
    banks.Bank = nil
    var bank Bank
    for rows.Next() {
      err = rows.Scan(&bank.Id, &bank.Name, &bank.Desc)
      checkErr(err)
      banks.Bank = append(banks.Bank, bank)
    }

    j, err := json.Marshal(banks.Bank)
    checkErr(err)
    w.Header().Set("Content-Type", "text/plain; charset=utf-8")
    w.Write(j)
    // fmt.Fprintf(w, "This is an example server.\n")
    // io.WriteString(w, "This is an example server.\n")

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
