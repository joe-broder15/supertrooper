package server

import (
	"crypto/rand"
	"encoding/hex"
	"sync"
	"time"
)

// information about an agent that the server keeps track of
type ServerAgentInfo struct {
	ID                string
	IpAddress         string
	Persist           bool
	BeaconInterval    int
	MissesBeforeDeath int
	LastContact       time.Time
	NextBeacon        time.Time
	DeathDate         time.Time
}

// a manager for agents that the server keeps track of
type AgentManager struct {
	m  map[string]ServerAgentInfo // map of agents. keys are all 256 bit hex strings
	mu sync.Mutex
}

func (am *AgentManager) RegisterAgent() string {
	am.mu.Lock()
	defer am.mu.Unlock()

	// generate a new agent id
	var hexString string
	for {
		// generate a random 256 bit value
		randomBytes := make([]byte, 32)
		_, err := rand.Read(randomBytes)
		if err != nil {
			return "none"
		}

		hexString = hex.EncodeToString(randomBytes)
		if _, exists := am.m[hexString]; !exists {
			break
		}
	}

	// add an empty agent info to the map
	am.m[hexString] = ServerAgentInfo{}

	// convert to hex string
	return hexString

}

func (am *AgentManager) UpdateAgent(agentID string, ipAddress string, persist bool, beaconInterval int, missesBeforeDeath int, lastContact time.Time, nextBeacon time.Time) {
	am.mu.Lock()
	defer am.mu.Unlock()

	am.m[agentID] = ServerAgentInfo{
		ID:                agentID,
		IpAddress:         ipAddress,
		Persist:           persist,
		BeaconInterval:    beaconInterval,
		MissesBeforeDeath: missesBeforeDeath,
		LastContact:       lastContact,
		NextBeacon:        nextBeacon,
		DeathDate:         time.Now().Add(time.Duration(missesBeforeDeath*beaconInterval) * time.Second),
	}
}

func (am *AgentManager) IsRegistered(agentID string) bool {
	am.mu.Lock()
	defer am.mu.Unlock()

	_, exists := am.m[agentID]
	return exists
}

// create a new agent manager
func NewAgentManager() *AgentManager {
	return &AgentManager{
		m: make(map[string]ServerAgentInfo),
	}
}
