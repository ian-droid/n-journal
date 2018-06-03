/*
References:
* https://github.com/denji/golang-tls
* http://www.bite-code.com/2015/06/25/tls-mutual-auth-in-golang/
* https://www.alexedwards.net/blog/serving-static-sites-with-go
* https://legacy.gitbook.com/book/astaxie/build-web-application-with-golang/details
*/

package main

import (
    // "fmt"
    //"io"
    "io/ioutil"
    "crypto/tls"
    "crypto/x509"
    "net/http"
    "log"
    "./journaldb"
)


func main() {
    banks := &journaldb.Banks{DBConn: journaldb.Open("ian_journal.db")}

    fs := http.FileServer(http.Dir("static"))
    http.Handle("/", fs)
    http.HandleFunc("/banks/list", banks.List)

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
