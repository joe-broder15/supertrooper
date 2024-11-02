package server

import (
	"crypto/tls"
	"crypto/x509"
	"log"
	"net/http"
	"os"
)

func Start(serverCertFile string, serverKeyFile string) {

	// Load server certificate and private key
	serverCert, err := tls.LoadX509KeyPair(serverCertFile, serverKeyFile)
	if err != nil {
		log.Fatalf("Failed to load server certificate: %v", err)
	}

	// Load CA certificate to validate client certificates
	caCert, err := os.ReadFile(serverCertFile)
	if err != nil {
		log.Fatalf("Failed to load CA certificate: %v", err)
	}
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	// Configure TLS with client certificate verification
	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{serverCert},
		ClientCAs:    caCertPool,
		ClientAuth:   tls.RequireAndVerifyClientCert, // Enforce client certificate verification
	}

	// Set up the HTTPS server using the above TLS config
	server := &http.Server{
		Addr:      ":443",
		TLSConfig: tlsConfig,
	}

	// Define a simple handler
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, mutual TLS client!"))
	})

	// Start the HTTPS server
	log.Println("Starting HTTPS server on port 443")
	err = server.ListenAndServeTLS("", "") // Certificates are already configured
	if err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
