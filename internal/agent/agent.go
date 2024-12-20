package agent

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/joe-broder15/supertrooper/internal/crypto"
	"github.com/joe-broder15/supertrooper/internal/messages"
)

func Start(agentCertPEM []byte, agentKeyPEM []byte, serverCertPEM []byte) {
	fmt.Println("starting agent")

	// Load agent certificate and private key from embedded data
	agentCert, err := tls.X509KeyPair(agentCertPEM, agentKeyPEM)
	if err != nil {
		log.Fatalf("Failed to load agent certificate: %v\n", err)
	}

	// Create a CertPool with the server certificate
	serverCertPool := x509.NewCertPool()
	if !serverCertPool.AppendCertsFromPEM(serverCertPEM) {
		log.Fatalln("Failed to append server certificate to pool")
	}

	// Configure TLS with custom verification
	tlsConfig := &tls.Config{
		Certificates:       []tls.Certificate{agentCert},
		RootCAs:            serverCertPool,
		InsecureSkipVerify: true,
	}

	// create https client
	transport := &http.Transport{
		TLSClientConfig: tlsConfig,
		MaxIdleConns:    10,
		IdleConnTimeout: 30 * time.Second,
	}
	client := &http.Client{Transport: transport}

	// get a nonce to send to the agent
	agentNonce, err := crypto.GenerateRandomBytes()
	if err != nil {
		log.Printf("failed to generate agent nonce: %v", err)
		return
	}

	// create payload
	challengeRequest, err := messages.NewAgentChallengeRequest(agentNonce)
	if err != nil {
		log.Fatal(err)
	}

	// post data
	resp, err := client.Post("https://localhost/challenge", "application/json", bytes.NewBuffer(challengeRequest))
	if err != nil {
		log.Println("ERROR")
		log.Fatalln(err)
	}

	// get response
	responseBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	// parse response
	message, err := messages.ParseServerChallengeResponseBody(responseBytes)
	if err != nil {
		log.Fatalln(err)
	}

	result, err := crypto.RSAVerifySignature(serverCertPEM, message.SignedAgentNonce, agentNonce)
	if err != nil {
		log.Fatalln(err)
	}

	println(result)
}
