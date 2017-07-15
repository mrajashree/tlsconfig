package main

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"time"

	log "github.com/Sirupsen/logrus"
)

func main() {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}
	_, err := client.Get("https://my.rancher.com/")
	if err != nil {
		fmt.Println(err)
		return
	}

	tlsconn, err := newConn()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("tlsconn Remote address: %v\n", tlsconn.RemoteAddr())
	fmt.Printf("tlsconn local address: %v\n", tlsconn.LocalAddr())
}

var certFile = "ca-certificates.crt"

func newConn() (*tls.Conn, error) {
	var err error
	var tlsConfig *tls.Config

	pool := x509.NewCertPool()
	certFound := false
	if _, err = os.Stat(certFile); os.IsNotExist(err) {
		log.Warningf("CA cert file %s is not present", certFile)
		return nil, fmt.Errorf("cannot find CA cert file '%s'", certFile)
	}
	caCert, err := ioutil.ReadFile(certFile)
	if err != nil {
		return nil, fmt.Errorf("cannot read CA cert file '%s'; err= %v", certFile, err)
	}

	if ok := pool.AppendCertsFromPEM(caCert); !ok {
		return nil, fmt.Errorf("cannot add CA cert file '%s'", certFile)
	}

	log.Infof("loaded CA cert file %s", certFile)
	certFound = true

	if certFound {
		tlsConfig = &tls.Config{RootCAs: pool, InsecureSkipVerify: false, ServerName: "my.rancher.com"}
		tlsConfig.BuildNameToCertificate()
	} else {
		tlsConfig = &tls.Config{InsecureSkipVerify: true}
	}

	dc, err := net.DialTimeout("tcp", "my.rancher.com:443", 60*time.Second)
	if err != nil {
		return nil, fmt.Errorf("Error during dial")
	}
	tlsConn := tls.Client(dc, tlsConfig)

	err = tlsConn.Handshake()
	if err != nil {
		// Handshake error, close the established connection before we return an error
		dc.Close()
		return nil, fmt.Errorf("Error during handshake")
	}

	log.Infof("TLS handshake completed")
	return tlsConn, nil
}
