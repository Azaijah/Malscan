package elastic

import (
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"malscan/config"
)

func getUserName() (username string) {

	username = config.Values.Elasticsearch.Username

	return username
}

func getUserPassword() (password string) {

	password = config.Values.Elasticsearch.Password

	return password

}

func getHTTPSClient() *http.Client {

	// Read the key pair to create certificate
	cert, err := tls.LoadX509KeyPair(config.Values.Elasticsearch.Cert, config.Values.Elasticsearch.Key)
	if err != nil {
		log.Fatal(err)
	}

	// Create a CA certificate pool and add cert.pem to it
	caCert, err := ioutil.ReadFile(config.Values.Elasticsearch.Ca)
	if err != nil {
		log.Fatal(err)
	}
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	// Create a HTTPS client and supply the created CA pool and certificate
	httpClient := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				RootCAs:      caCertPool,
				Certificates: []tls.Certificate{cert},
			},
		}, Timeout: time.Second * 10,
	}

	return httpClient

}
