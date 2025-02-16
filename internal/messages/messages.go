// Package messages defines the various message types, payload formats, and helper functions
// for building and parsing messages for the C2 communication system.
package messages

import (
	"encoding/json"
	"log"
	"time"
)

// ---------------------------------------------------------------------
// Type Definitions and Constants
// ---------------------------------------------------------------------

// C2MessageType represents the type of a C2 message.
type C2MessageType int

// JobType represents the type of job.
type JobType int

// JobStatus represents the status of a job.
type JobStatus string

// Definitions of C2 message types.
const (
	C2MessageTypeHello C2MessageType = iota
	// C2MessageTypeHelloRsp
	C2MessageTypeGoodbye
	C2MessageTypeJobDispatch
	C2MessageTypeJobDispatchRsp
	C2MessageTypeKill
	C2MessageTypeKillRsp
)

// Definitions of job types.
const (
	JobTypePrint JobType = iota
)

// Definitions of job statuses.
const (
	JobStatusSuccess JobStatus = "success"
	JobStatusError   JobStatus = "error"
)

// C2MessageBase is the base structure for all C2 messages.
type C2MessageBase struct {
	Type             C2MessageType   // Type of the message.
	Timestamp        int64           // Unix timestamp when the message was created.
	AgentID          string          // Identifier for the agent.
	C2MessagePayload json.RawMessage // Encoded payload specific to the message type.
}

// C2MessageTypeHelloPayload is the payload for the "Hello" message (sent from agent to server).
type C2MessageTypeHelloPayload struct {
	IsPersistent      bool // Indicates if the agent should maintain a persistent connection.
	BeaconSeconds     int  // Interval (in seconds) for beaconing.
	MissesBeforeDeath int  // Number of missed beacons before considering the agent dead.
	WakeUpSeconds     int  // Time until the agent will be awake after a beacon.
	NextBeacon        int  // Time until the next beacon should be sent.
	DeathDate         int  // Time until the agent should be considered dead.
}

// C2MessageTypeGoodbyePayload is the payload for the "Goodbye" message (sent from agent to server).
type C2MessageTypeGoodbyePayload struct {
	MissesBeforeDeath int // Number of missed beacons before considering the agent dead.
	NextBeacon        int // Time until the next beacon should be sent.
	DeathDate         int // Time until the agent should be considered dead.
}

// C2MessageTypeKillPayload is the payload for "Kill" and "Kill Response" messages.
type C2MessageTypeKillPayload struct {
	LastWords string // Final message or reason for termination.
}

// C2MessageTypeJobDispatchPayload is the payload for the "Job Dispatch" message.
type C2MessageTypeJobDispatchPayload struct {
	JobID      string          // Unique identifier for the job.
	JobType    JobType         // Type of job.
	JobPayload json.RawMessage // Encoded job-specific payload.
}

// C2MessageTypeJobDispatchRspPayload is the payload for the "Job Dispatch Response" message.
type C2MessageTypeJobDispatchRspPayload struct {
	JobID      string          // Unique identifier for the job.
	JobType    JobType         // Type of job.
	JobStatus  JobStatus       // Status of the job execution.
	JobError   string          // Error message (if any) encountered during job execution.
	JobPayload json.RawMessage // Encoded response-specific payload.
}

// ---------------------------------------------------------------------
// Job Payload Types
// ---------------------------------------------------------------------

// JobPayloadPrint is a specific type of job payload for printing a message.
type JobPayloadPrint struct {
	Content string // The message to print.
}

// JobPayloadPrintRsp is a specific type of job payload for printing a message.
type JobPayloadPrintRsp struct {
	Success bool // Whether the print job was successful.
}

// ---------------------------------------------------------------------
// Helper Functions
// ---------------------------------------------------------------------

// mustMarshal marshals a value to JSON or logs a fatal error if marshalling fails.
// This helper function is used by builder functions to reduce duplicate error checking.
func mustMarshal(v interface{}) json.RawMessage {
	data, err := json.Marshal(v)
	if err != nil {
		log.Fatalf("failed to marshal %+v: %v", v, err)
	}
	return data
}

// BuildMessage constructs a basic C2MessageBase with the given type, agent ID, and JSON payload.
func BuildMessage(messageType C2MessageType, agentID string, payload json.RawMessage) C2MessageBase {
	return C2MessageBase{
		Type:             messageType,
		Timestamp:        time.Now().Unix(),
		AgentID:          agentID,
		C2MessagePayload: payload,
	}
}

//////////////////////////
// Builder Functions    //
//////////////////////////

// BuildHelloMessage creates a "Hello" message with the provided parameters.
func BuildHelloMessage(agentID string, isPersistent bool, beaconSeconds int, missesBeforeDeath int, deathDate int, wakeUpSeconds int, nextBeacon int) C2MessageBase {
	payload := C2MessageTypeHelloPayload{
		IsPersistent:      isPersistent,
		BeaconSeconds:     beaconSeconds,
		MissesBeforeDeath: missesBeforeDeath,
		WakeUpSeconds:     wakeUpSeconds,
		NextBeacon:        nextBeacon,
		DeathDate:         deathDate,
	}
	return BuildMessage(C2MessageTypeHello, agentID, mustMarshal(payload))
}

// BuildGoodbyeMessage creates a "Goodbye" message with the given agent ID and beacon interval.
func BuildGoodbyeMessage(agentID string, missesBeforeDeath int, nextBeacon int, deathDate int) C2MessageBase {
	payload := C2MessageTypeGoodbyePayload{
		MissesBeforeDeath: missesBeforeDeath,
		NextBeacon:        nextBeacon,
		DeathDate:         deathDate,
	}
	return BuildMessage(C2MessageTypeGoodbye, agentID, mustMarshal(payload))
}

// BuildJobDispatchMessage creates a "Job Dispatch" message for scheduling a new job.
func BuildJobDispatchMessage(agentID string, jobID string, jobType JobType, jobPayload json.RawMessage) C2MessageBase {
	dispatchPayload := C2MessageTypeJobDispatchPayload{
		JobID:      jobID,
		JobType:    jobType,
		JobPayload: jobPayload,
	}
	return BuildMessage(C2MessageTypeJobDispatch, agentID, mustMarshal(dispatchPayload))
}

// BuildJobDispatchRspMessage creates a "Job Dispatch Response" message to report job execution results.
func BuildJobDispatchRspMessage(agentID, jobID string, jobType JobType, jobStatus JobStatus, jobError string, jobPayload json.RawMessage) C2MessageBase {
	rspPayload := C2MessageTypeJobDispatchRspPayload{
		JobID:      jobID,
		JobType:    jobType,
		JobStatus:  jobStatus,
		JobError:   jobError,
		JobPayload: jobPayload,
	}
	return BuildMessage(C2MessageTypeJobDispatchRsp, agentID, mustMarshal(rspPayload))
}

// BuildKillMessage creates a "Kill" message to request termination (sent from agent to server).
func BuildKillMessage(agentID, lastWords string) C2MessageBase {
	killPayload := C2MessageTypeKillPayload{
		LastWords: lastWords,
	}
	return BuildMessage(C2MessageTypeKill, agentID, mustMarshal(killPayload))
}

// BuildKillRspMessage creates a "Kill Response" message acknowledging termination (sent from server to agent).
func BuildKillRspMessage(agentID, lastWords string) C2MessageBase {
	killRspPayload := C2MessageTypeKillPayload{
		LastWords: lastWords,
	}
	return BuildMessage(C2MessageTypeKillRsp, agentID, mustMarshal(killRspPayload))
}

// BuildJobPayloadPrint creates a job payload for printing a message.
func BuildJobPayloadPrint(content string) json.RawMessage {
	payload := JobPayloadPrint{
		Content: content,
	}
	return mustMarshal(payload)
}

//////////////////////////
// Parser Functions     //
//////////////////////////

// ParseHelloPayload decodes a JSON payload into a C2MessageTypeHelloPayload struct.
func ParseHelloPayload(message []byte) (C2MessageTypeHelloPayload, error) {
	var payload C2MessageTypeHelloPayload
	err := json.Unmarshal(message, &payload)
	return payload, err
}

// ParseGoodbyePayload decodes a JSON payload into a C2MessageTypeGoodbyePayload struct.
func ParseGoodbyePayload(message []byte) (C2MessageTypeGoodbyePayload, error) {
	var payload C2MessageTypeGoodbyePayload
	err := json.Unmarshal(message, &payload)
	return payload, err
}

// ParseJobDispatchPayload decodes a JSON payload into a C2MessageTypeJobDispatchPayload struct.
func ParseJobDispatchPayload(message []byte) (C2MessageTypeJobDispatchPayload, error) {
	var payload C2MessageTypeJobDispatchPayload
	err := json.Unmarshal(message, &payload)
	return payload, err
}

// ParseJobDispatchRspPayload decodes a JSON payload into a C2MessageTypeJobDispatchRspPayload struct.
func ParseJobDispatchRspPayload(message []byte) (C2MessageTypeJobDispatchRspPayload, error) {
	var payload C2MessageTypeJobDispatchRspPayload
	err := json.Unmarshal(message, &payload)
	return payload, err
}

// ParseKillPayload decodes a JSON payload into a C2MessageTypeKillPayload struct.
// This parser is used for both "Kill" and "Kill Response" messages.
func ParseKillPayload(message []byte) (C2MessageTypeKillPayload, error) {
	var payload C2MessageTypeKillPayload
	err := json.Unmarshal(message, &payload)
	return payload, err
}

// ParseJobPayloadPrint decodes a JSON payload into a JobPayloadPrint struct.
func ParseJobPayloadPrint(message []byte) (JobPayloadPrint, error) {
	var payload JobPayloadPrint
	err := json.Unmarshal(message, &payload)
	return payload, err
}
