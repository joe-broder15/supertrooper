package server

import (
	"log"
	"sync"

	"github.com/joe-broder15/supertrooper/internal/common"
)

// entries in the agent manager map, this contains all the info the server will track about each agent
type AgentData struct {
	AgentConfig   common.AgentConfig
	AgentHostInfo common.AgentHostInfo
	DeathDate     int
	NextBeacon    int
	PublicIP      string
}

// a glorified locked hashmap
type agentManager struct {
	agentMapLock sync.RWMutex
	agentMap     map[string]*AgentData
}

// constructor for agentManager
func newAgentManager() *agentManager {
	return &agentManager{
		agentMap: make(map[string]*AgentData),
	}
}

// takes in a beacon request and updates the agentdata accordingly
func (am *agentManager) ProcessBeacon(beaconReq common.BeaconReq, remoteAddr string) {

	// acquire the lock
	am.agentMapLock.Lock()
	defer am.agentMapLock.Unlock()

	// check if the agent exists, otherwise register it
	_, ok := am.agentMap[beaconReq.AgentConfig.AgentID]
	if !ok {
		am.agentMap[beaconReq.AgentConfig.AgentID] = &AgentData{}
		log.Printf("registered new agent %v", beaconReq.AgentConfig.AgentID)
	}

	// get and update the agent data
	data := am.agentMap[beaconReq.AgentConfig.AgentID]

	data.AgentConfig = beaconReq.AgentConfig
	data.AgentHostInfo = beaconReq.AgentHostInfo
	data.NextBeacon = beaconReq.NextBeaconExpected
	data.DeathDate = beaconReq.AgentConfig.BeaconTries*(beaconReq.AgentConfig.BeaconInterval-1) + beaconReq.NextBeaconExpected
	data.PublicIP = remoteAddr
}

// takes an agent id and returns an agent data
func (am *agentManager) GetAgentData(agentID string) (AgentData, bool) {
	am.agentMapLock.RLock()
	defer am.agentMapLock.RUnlock()
	if data, ok := am.agentMap[agentID]; !ok {
		return AgentData{}, false
	} else {
		return *data, true
	}
}
