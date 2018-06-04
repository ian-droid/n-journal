/*
References:
* https://github.com/denji/golang-tls
* http://www.bite-code.com/2015/06/25/tls-mutual-auth-in-golang/
* https://www.alexedwards.net/blog/serving-static-sites-with-go
* https://legacy.gitbook.com/book/astaxie/build-web-application-with-golang/details
*/

package main

import (
    "fmt"
    "strings"
    //"io"
    "io/ioutil"
    "crypto/tls"
    "crypto/x509"
    "net/http"
    "database/sql"
    "html/template"
    "log"
    "./journaldb"
)


type DiaryForm struct {
    DBConn *sql.DB
    Diary []journaldb.Diary
}

func (diaryForm *DiaryForm) Form(w http.ResponseWriter, r *http.Request) {
    if r.Method == "POST" {
        r.ParseForm()
        nd := journaldb.Diary{}

        if oid, ok := r.Form["oid"]; ok {
          fmt.Sscanf(strings.Join(oid, ""),"%d", &nd.Oid)
          fmt.Println("Update existing ")
        } else {
          nd.Oid = 0
        }
        nd.Date = strings.Join(r.Form["date"], "")
        nd.Content = strings.Join(r.Form["content"], "")
        if strings.Join(r.Form["highlighted"], "") == "on" {
          nd.Highlighted = true
        } else {
          nd.Highlighted = false
        }
        nd.Save(diaryForm.DBConn)
    }

    rows, err := diaryForm.DBConn.Query("select oid, date, content, highlighted from diary order by date desc limit 7")
    checkErr(err)
    diaryForm.Diary = nil
    var diary journaldb.Diary
    for rows.Next() {
      err = rows.Scan(&diary.Oid, &diary.Date, &diary.Content, &diary.Highlighted)
      checkErr(err)
      diaryForm.Diary = append(diaryForm.Diary, diary)
    }

    tmpl := template.Must(template.ParseFiles("diary.gtpl"))
    tmpl.Execute(w, diaryForm)


}

func checkErr(err error) {
    if err != nil {
        panic(err)
    }
}

func main() {
    dbconn := journaldb.Open("ian_journal.db")
    diaryForm := DiaryForm{DBConn: dbconn}
//    banks := &journaldb.Banks{DBConn: dbconn}

    fs := http.FileServer(http.Dir("static"))
    http.Handle("/", fs)
    http.HandleFunc("/diary", diaryForm.Form)

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
