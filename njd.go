// References:
// * https://github.com/denji/golang-tls
// * http://www.bite-code.com/2015/06/25/tls-mutual-auth-in-golang/

package main

import (
    // "fmt"
    //"io"
    "io/ioutil"
    "crypto/tls"
    "crypto/x509"
    "net/http"
    "log"
)

func HelloServer(w http.ResponseWriter, req *http.Request) {
    w.Header().Set("Content-Type", "text/plain")
    w.Write([]byte("This is an example server.\n"))
    // fmt.Fprintf(w, "This is an example server.\n")
    // io.WriteString(w, "This is an example server.\n")
}

func main() {
    http.HandleFunc("/hello", HelloServer)

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
