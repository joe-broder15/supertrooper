package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"time"

	"github.com/joe-broder15/supertrooper/internal/messages"
)

// build a beacon response for an agent
func buildBeaconRsp(agentManager *AgentManager, agentID string) (messages.BeaconRsp, error) {

	// build the server info
	serverInfo := messages.ServerInfo{
		ServerID:      "none",
		ServerVersion: "none",
		ServerTime:    int(time.Now().Unix()),
	}

	// build the reconfigure info
	reconfigureInfo := messages.ReconfigureInfo{
		AgentID:           agentID,
		BeaconInterval:    10,
		MissesBeforeDeath: 3,
		Persist:           false,
	}

	// build the job req
	jobReq := messages.JobReq{
		JobID:      "none",
		JobType:    messages.JobTypePrint,
		JobPayload: json.RawMessage(`"none"`),
	}

	// build the beacon rsp
	beaconRsp := messages.BeaconRsp{
		ServerInfo:  serverInfo,
		Reconfigure: reconfigureInfo,
		JobReq:      []messages.JobReq{jobReq},
	}

	return beaconRsp, nil
}

// send a beacon response to the agent
func sendBeaconRsp(beaconRsp messages.BeaconRsp, conn net.Conn) error {

	// encode and send the beacon rsp to the agent as json
	encoder := json.NewEncoder(conn)
	err := encoder.Encode(beaconRsp)
	if err != nil {
		return fmt.Errorf("server: error sending beacon rsp: %v", err)
	}

	return nil
}

// function where the server actually handles messages recieved from agents
func processBeaconReqEvent(event ServerEventBeaconReq, agentManager *AgentManager) error {

	// pretty print the event
	jsonBytes, err := json.MarshalIndent(event, "", "    ")
	if err != nil {
		return fmt.Errorf("server: error marshaling event to JSON: %v", err)
	}
	log.Println("server: pretty printed event JSON:")
	log.Println(string(jsonBytes))

	// get the agent id
	agentID := event.BeaconReq.AgentInfo.AgentID

	// check if the agent is registered, and register it if it is not
	if agentID == "none" {
		agentID = agentManager.RegisterAgent()
	}

	// check if the agent is registered with the agent manager
	if !agentManager.IsRegistered(agentID) {
		return fmt.Errorf("server: agent %s is not registered", agentID)
	}

	// update the agent info
	agentManager.UpdateAgent(agentID, event.AgentConn.RemoteAddr().String(), event.BeaconReq.AgentInfo.Persist, event.BeaconReq.AgentInfo.BeaconInterval, event.BeaconReq.AgentInfo.MissesBeforeDeath, time.Now(), time.Now().Add(time.Duration(event.BeaconReq.AgentInfo.BeaconInterval)*time.Second))

	//  TODO: process the job responses

	// build the beacon response
	beaconRsp, err := buildBeaconRsp(agentManager, agentID)
	if err != nil {
		return fmt.Errorf("server: error building beacon response: %v", err)
	}

	// send a beacon response to the agent
	err = sendBeaconRsp(beaconRsp, event.AgentConn)
	if err != nil {
		return fmt.Errorf("server: error sending beacon response: %v", err)
	}

	return nil

}

// function where the server actually handles errors
func processErrorEvent(event ServerEventError) {
	log.Printf("server: error: %v", event.Error)
}

func processEvent(event ServerEvent, agentManager *AgentManager) {
	switch event.Type {
	case ServerEventTypeBeaconReq:
		processBeaconReqEvent(event.Body.(ServerEventBeaconReq), agentManager)
	case ServerEventTypeError:
		processErrorEvent(event.Body.(ServerEventError))
	}
}
