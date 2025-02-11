package agent

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"log"
	"net/http"
	"time"
)

func Start(caCertPEM []byte, agentCertPEM []byte, agentKeyPEM []byte, serverCertPEM []byte) {
	fmt.Println("starting agent as mTLS client")

	// Load the agent's certificate and key for client identity
	agentCert, err := tls.X509KeyPair(agentCertPEM, agentKeyPEM)
	if err != nil {
		log.Fatalf("agent: error loading key pair: %v", err)
	}

	// Create a CA pool with the provided CA certificate to verify the server
	caCertPool := x509.NewCertPool()
	if ok := caCertPool.AppendCertsFromPEM(caCertPEM); !ok {
		log.Fatal("agent: failed to append CA certificate")
	}

	// Configure mTLS settings for the client
	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{agentCert},
		RootCAs:      caCertPool,
		MinVersion:   tls.VersionTLS13,
		// Optionally set ServerName if needed for hostname verification
		// ServerName: "your.server.domain",
	}

	// Create an HTTP client with the custom TLS transport
	httpClient := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: tlsConfig,
		},
		Timeout: 10 * time.Second,
	}

	// Make an example HTTPS request to the mTLS server
	resp, err := httpClient.Get("https://localhost:8443")
	if err != nil {
		log.Fatalf("agent: error making HTTPS request: %v", err)
	}
	defer resp.Body.Close()

	fmt.Printf("agent: received response with status: %s\n", resp.Status)
}
