package server

import (
	"crypto/rand"
	"fmt"
	"time"
)

// information about an agent that the server keeps track of
type ServerAgentInfo struct {
	ID                string
	IpAddress         string
	IsLive            bool
	IsPersistent      bool
	BeaconSeconds     int
	MissesBeforeDeath int
	LastContact       time.Time
	NextBeacon        time.Time
	DeathDate         time.Time
}

// a manager for agents that the server keeps track of
type AgentManager struct {
	m map[string]ServerAgentInfo // map of agents. keys are all 256 bit hex strings
}

// register a new agent with the agent manager. generate a new agent id and add it to the map. this id must be unique and 256 bits hex
func (am *AgentManager) RegisterAgent(ipAddress string, isPersistent bool, beaconSeconds int, missesBeforeDeath int) string {

	// generate a new agent id
	var agentID string
	for {
		b := make([]byte, 32) // 256 bits = 32 bytes
		if _, err := rand.Read(b); err != nil {
			panic("failed to generate agent id: " + err.Error())
		}
		agentID = fmt.Sprintf("%x", b)
		if _, exists := am.m[agentID]; !exists {
			break
		}
	}

	// add the agent to the map
	am.m[agentID] = ServerAgentInfo{
		ID:                agentID,
		IpAddress:         ipAddress,
		IsPersistent:      isPersistent,
		BeaconSeconds:     beaconSeconds,
		MissesBeforeDeath: missesBeforeDeath,
		LastContact:       time.Now(),
		NextBeacon:        time.Now().Add(time.Duration(beaconSeconds) * time.Second),
		DeathDate:         time.Now().Add(time.Duration(missesBeforeDeath*beaconSeconds) * time.Second),
	}

	return agentID
}

// update an existing agent with new info
func (am *AgentManager) UpdateAgent(agentID string, ipAddress string, isPersistent bool, beaconSeconds int, missesBeforeDeath int, nextBeacon int, deathDate int) {
	am.m[agentID] = ServerAgentInfo{
		ID:                agentID,
		IpAddress:         ipAddress,
		IsPersistent:      isPersistent,
		BeaconSeconds:     beaconSeconds,
		MissesBeforeDeath: missesBeforeDeath,
		LastContact:       time.Now(),
		NextBeacon:        time.Unix(int64(nextBeacon), 0),
		DeathDate:         time.Unix(int64(deathDate), 0),
	}
}

// create a new agent manager
func NewAgentManager() *AgentManager {
	return &AgentManager{
		m: make(map[string]ServerAgentInfo),
	}
}
