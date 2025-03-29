package server

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/joe-broder15/supertrooper/internal/common"
)

func (ss *ServerState) HandleBeacon(w http.ResponseWriter, r *http.Request) {

	// drop all requests that are not post and put
	if !(r.Method == "POST" || r.Method == "PUT") {
		return
	}

	// pre-emptively set header
	w.Header().Set("Content-Type", "application/json")

	// create a json decoder
	decoder := json.NewDecoder(r.Body)

	// attempt to decode the request
	var beaconReq common.BeaconReq
	err := decoder.Decode(&beaconReq)

	// check if the request was malformed and respond accordingly
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Invalid request format: " + err.Error(),
		})
	} else {

		// process the beacon request by updating the job manager as well as the agent manager
		log.Println("got request from agent: " + beaconReq.AgentConfig.AgentID)
		ss.agentManager.ProcessBeacon(beaconReq, r.RemoteAddr)
		ss.jobManager.registerCompletedJobs(beaconReq.AgentConfig.AgentID, &beaconReq.JobResponses)

		// then create and send the response
		w.WriteHeader(http.StatusOK)
		err = json.NewEncoder(w).Encode(common.BeaconRsp{
			ServerID:    uuid.New().String(),
			JobRequests: ss.jobManager.taskUntaskedJobs(beaconReq.AgentConfig.AgentID), // promote all untasked jobs to tasked status and send them in the response
		})
		if err != nil {
			// TODO: update job manager in case of things failing
			log.Printf("Error: %v", err)
		}

		log.Printf("send response to agent")
	}

}
