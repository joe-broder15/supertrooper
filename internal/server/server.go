package server

import (
	"crypto/tls"
	"crypto/x509"
	"log"
	"net/http"
	"os"
)

type ServerState struct {
	agentManager *agentManager
	jobManager   *jobManager
}

func NewServerState() *ServerState {
	return &ServerState{
		agentManager: newAgentManager(),
		jobManager:   newJobManager(),
	}
}

func Start(serverCertFile string, serverKeyFile string, caCertFile string, configFile string) {

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

	// Set up mTLS configuration for the server.
	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{cert},
		ClientAuth:   tls.RequireAndVerifyClientCert,
		ClientCAs:    caCertPool,
		MinVersion:   tls.VersionTLS12,
	}

	// Create the server
	server := &http.Server{
		Addr:      ":443",
		TLSConfig: tlsConfig,
	}

	// initialize server state
	serverState := NewServerState()

	// register functions
	http.HandleFunc("/", serverState.HandleBeacon)

	// Start the server
	log.Println("Starting mTLS server on :443")
	log.Fatal(server.ListenAndServeTLS("", ""))

}
