package server

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"log"
	"net/http"
	"os"
)

func Start(serverCertFile string, serverKeyFile string, agentCertFile string, caCertFile string) {
	fmt.Println("starting server as mTLS server")

	// Load server certificate and private key
	cert, err := tls.LoadX509KeyPair(serverCertFile, serverKeyFile)
	if err != nil {
		log.Fatalf("server: error loading key pair: %v", err)
	}

	// Load CA certificate for verifying client's certificate
	caCert, err := os.ReadFile(caCertFile)
	if err != nil {
		log.Fatalf("server: error reading CA certificate: %v", err)
	}

	caCertPool := x509.NewCertPool()
	if ok := caCertPool.AppendCertsFromPEM(caCert); !ok {
		log.Fatal("server: failed to append CA certificate")
	}

	// Set up mTLS configuration for the server. The server will verify client certificates.
	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{cert},
		ClientAuth:   tls.RequireAndVerifyClientCert,
		ClientCAs:    caCertPool,
		MinVersion:   tls.VersionTLS13,
	}

	// Create an HTTP server with the TLS configuration
	server := &http.Server{
		Addr:      ":8443",
		TLSConfig: tlsConfig,
	}

	// Example handler for demonstration
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, mTLS client!")
	})

	log.Println("server: listening on https://localhost:8443")
	err = server.ListenAndServeTLS("", "")
	if err != nil {
		log.Fatalf("server: ListenAndServeTLS error: %v", err)
	}
}
