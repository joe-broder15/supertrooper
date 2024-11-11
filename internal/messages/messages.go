package messages

import (
	"encoding/json"
	"fmt"
)

type MessageType int

const (
	MessageTypeError MessageType = iota
	MessageTypeAgentChallengeRequest
	MessageTypeServerChallengeResponse
	MessageTypeLast //MUST BE LAST ITEM IN ENUM
)

type Message struct {
	Type  MessageType
	Body  string
	Token string
}

type AgentChallengeRequestBody struct {
	AgentNonce string
}

type ServerChallengeResponseBody struct {
	SignedAgentNonce string
	ServerNonce      string
}

/*
-------------------------
MESSAGE AND BODY PARSERS
-------------------------
*/

func ParseMessage(messageBytes []byte) (Message, error) {

	// attempt to unmarshal the message bytes
	var message Message
	err := json.Unmarshal(messageBytes, &message)
	if err != nil {
		return message, fmt.Errorf("failed to unmarshal message struct: %v", err)
	}

	// check the message type
	if !(message.Type < MessageTypeLast && message.Type >= MessageTypeError) {
		return message, fmt.Errorf("invalid message type: %v", message.Type)
	}

	return message, nil
}

// Parses the AgentChallengeRequestBody from a byte slice
func ParseAgentChallengeRequestBody(messageBytes []byte) (AgentChallengeRequestBody, error) {
	message, err := ParseMessage(messageBytes)
	if err != nil {
		return AgentChallengeRequestBody{}, err
	}

	if message.Type != MessageTypeAgentChallengeRequest {
		return AgentChallengeRequestBody{}, fmt.Errorf("invalid message type: expected MessageTypeAgentChallengeRequest, got %v", message.Type)
	}

	var body AgentChallengeRequestBody
	err = json.Unmarshal([]byte(message.Body), &body)
	if err != nil {
		return AgentChallengeRequestBody{}, fmt.Errorf("failed to unmarshal AgentChallengeRequestBody: %v", err)
	}
	return body, nil
}

// Parses the ServerChallengeResponseBody from a byte slice
func ParseServerChallengeResponseBody(messageBytes []byte) (ServerChallengeResponseBody, error) {
	message, err := ParseMessage(messageBytes)
	if err != nil {
		return ServerChallengeResponseBody{}, err
	}

	if message.Type != MessageTypeServerChallengeResponse {
		return ServerChallengeResponseBody{}, fmt.Errorf("invalid message type: expected MessageTypeServerChallengeResponse, got %v", message.Type)
	}

	var body ServerChallengeResponseBody
	err = json.Unmarshal([]byte(message.Body), &body)
	if err != nil {
		return ServerChallengeResponseBody{}, fmt.Errorf("failed to unmarshal ServerChallengeResponseBody: %v", err)
	}
	return body, nil
}

/*
-------------------------
MESSAGE AND BODY BUILDERS
-------------------------
*/

// creates and marshals a message struct with a specified body, token, and type
func NewMessage(messageType MessageType, body any, token string) ([]byte, error) {
	// marshal the body and check for errors
	bodyBytes, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal body: %v", err)
	}

	// create the message
	message := Message{
		Type:  messageType,
		Body:  string(bodyBytes),
		Token: token,
	}

	// marshal the message and check for errors
	messageBytes, err := json.Marshal(message)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal Message: %v", err)
	}

	return messageBytes, nil
}

// creates a new AgentChallengeRequest message
func NewAgentChallengeRequest(agentNonce string) ([]byte, error) {
	body := AgentChallengeRequestBody{
		AgentNonce: agentNonce,
	}
	return NewMessage(MessageTypeAgentChallengeRequest, body, "")
}

// creates a new ServerChallengeResponse message
func NewServerChallengeResponse(signedAgentNonce string, serverNonce string) ([]byte, error) {
	body := ServerChallengeResponseBody{
		SignedAgentNonce: signedAgentNonce,
		ServerNonce:      serverNonce,
	}
	return NewMessage(MessageTypeServerChallengeResponse, body, "")
}
