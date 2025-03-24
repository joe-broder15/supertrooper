package agent

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/joe-broder15/supertrooper/internal/common"
)

type AgentState struct {
	Config        common.AgentConfig
	CompletedJobs []common.JobRsp
	PendingJobs   []common.JobReq
}

func NewAgentState() *AgentState {
	return &AgentState{
		Config: common.AgentConfig{
			AgentID: uuid.New().String(),
		},
		CompletedJobs: []common.JobRsp{},
		PendingJobs:   []common.JobReq{},
	}
}

// initialize an https client with mtls using the provided CA cert, agent cert, and agent key
func NewHttpsClient(caCertPEM []byte, agentCertPEM []byte, agentKeyPEM []byte) (*http.Client, error) {

	// Load the agent's certificate and key for client identity.
	agentCert, err := tls.X509KeyPair(agentCertPEM, agentKeyPEM)
	if err != nil {
		log.Fatalf("agent: error loading key pair: %v", err)
	}

	// Create a CA pool with the provided CA certificate to verify the server.
	caCertPool := x509.NewCertPool()
	if ok := caCertPool.AppendCertsFromPEM(caCertPEM); !ok {
		return nil, fmt.Errorf("agent: failed to append CA certificate")
	}

	// create a tls config with the above certificates
	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{agentCert},
		RootCAs:      caCertPool,
		MinVersion:   tls.VersionTLS12,
	}

	// Create the https transport with the TLS config
	transport := &http.Transport{
		TLSClientConfig: tlsConfig,
	}

	// create and return the client
	return &http.Client{Transport: transport}, nil
}

// Start is the entry point for the agent. It initializes the state and
func Start(caCertPEM []byte, agentCertPEM []byte, agentKeyPEM []byte) error {

	// Initialize the HTTPS client with mTLS
	client, err := NewHttpsClient(caCertPEM, agentCertPEM, agentKeyPEM)

	// Make the request
	url := "https://localhost:443"
	resp, err := client.Get(url)
	if err != nil {
		log.Fatalf("Error making request: %v", err)
	}
	defer resp.Body.Close()

	// Read and print the response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Error reading response: %v", err)
	}

	fmt.Printf("Status: %s\n", resp.Status)
	fmt.Printf("Body: %s\n", string(body))

	return nil
}
