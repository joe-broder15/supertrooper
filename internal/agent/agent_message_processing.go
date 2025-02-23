package agent

import (
	"encoding/json"
	"log"

	"github.com/joe-broder15/supertrooper/internal/messages"
)

func processBeaconRsp(beaconRsp messages.BeaconRsp, agentState *AgentState) {
	log.Println("agent: received beacon response from server")

	prettyJSON, err := json.MarshalIndent(beaconRsp, "", "    ")
	if err != nil {
		log.Printf("agent: error pretty printing beacon response: %v", err)
	} else {
		log.Println("agent: pretty printed beacon response:")
		log.Println(string(prettyJSON))
	}

	// update the agent state with the reconfigure info
	agentState.agentID = beaconRsp.Reconfigure.AgentID
	agentState.beaconInterval = beaconRsp.Reconfigure.BeaconInterval
	agentState.missesBeforeDeath = beaconRsp.Reconfigure.MissesBeforeDeath
	agentState.persist = beaconRsp.Reconfigure.Persist

}
