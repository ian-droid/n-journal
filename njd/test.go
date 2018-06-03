package main

import ("./journaldb")

func main() {
    conn := journaldb.Open("ian_journal.db")
    defer journaldb.Close(conn)
    banks := &journaldb.Banks{DBConn: conn}
    banks.List()
}
