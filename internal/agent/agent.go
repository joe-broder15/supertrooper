package agent

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

// Start is the entry point for the agent. It initializes the state and
// enters a loop to repeatedly beacon the server.
func Start(caCertPEM []byte, agentCertPEM []byte, agentKeyPEM []byte) error {

	// Load the agent's certificate and key for client identity.
	agentCert, err := tls.X509KeyPair(agentCertPEM, agentKeyPEM)
	if err != nil {
		log.Fatalf("agent: error loading key pair: %v", err)
	}

	// Create a CA pool with the provided CA certificate to verify the server.
	caCertPool := x509.NewCertPool()
	if ok := caCertPool.AppendCertsFromPEM(caCertPEM); !ok {
		log.Fatal("agent: failed to append CA certificate")
	}

	// Configure mTLS settings for the client.
	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{agentCert},
		RootCAs:      caCertPool,
		MinVersion:   tls.VersionTLS12,
	}

	// Create transport and client
	transport := &http.Transport{
		TLSClientConfig: tlsConfig,
	}
	client := &http.Client{Transport: transport}

	// Make the request
	url := "https://localhost:443"
	resp, err := client.Get(url)
	if err != nil {
		log.Fatalf("Error making request: %v", err)
	}
	defer resp.Body.Close()

	// Read and print the response
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Error reading response: %v", err)
	}

	fmt.Printf("Status: %s\n", resp.Status)
	fmt.Printf("Body: %s\n", string(body))

	return nil
}
