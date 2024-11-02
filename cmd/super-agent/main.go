package main

import (
	"crypto/tls"
	"crypto/x509"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

func main() {
	fmt.Println("starting agent")
	// set up command line args
	agentCertPtr := flag.String("agent-cert", "certs/agent/agent.crt", "path to agent cert file")
	agentKeyPtr := flag.String("agent-key", "certs/agent/agent.key", "path to agent private key file")
	caCertPtrPtr := flag.String("ca-cert", "certs/server/server.key", "path to server cert file")
	flag.Parse()

	// Load server certificate and private key
	agentCert, err := tls.LoadX509KeyPair(*agentCertPtr, *agentKeyPtr)
	if err != nil {
		log.Fatalf("Failed to load server certificate: %v", err)
	}

	// Load CA certificate to validate agent certificates
	caCert, err := os.ReadFile(*caCertPtrPtr)
	if err != nil {
		log.Fatalf("Failed to load CA certificate: %v", err)
	}
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	// Configure TLS with agent certificate verification
	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{agentCert},
		RootCAs:      caCertPool,
	}

	transport := &http.Transport{
		TLSClientConfig: tlsConfig,
		MaxIdleConns:    10,
		IdleConnTimeout: 30 * time.Second,
	}
	client := &http.Client{Transport: transport}
	resp, err := client.Get("https://localhost/")
	if err != nil {
		log.Println("ERROR")
		log.Fatal(err)
	}
	log.Println(resp)

}
