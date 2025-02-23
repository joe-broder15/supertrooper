// Package messages defines the various message types, payload formats, and helper functions
// for building and parsing messages for the C2 communication system.
package messages

import (
	"encoding/json"
)

// ---------------------------------------------------------------------
// Type Definitions and Constants
// ---------------------------------------------------------------------

// JobType represents the type of job.
type JobType int

// JobStatus represents the status of a job.
type JobStatus int

// Definitions of job types.
const (
	JobTypePrint JobType = iota
)

// Definitions of job statuses.
const (
	JobStatusPending JobStatus = iota
	JobStatusSuccess
	JobStatusFailure
)

// ---------------------------------------------------------------------
// Agent and Server Information Types
// ---------------------------------------------------------------------

// AgentInfo holds information about the agent.
// JSON tags are provided for correct marshalling/unmarshalling.
type AgentInfo struct {
	AgentID           string `json:"agent_id"`
	BeaconInterval    int    `json:"beacon_interval"`
	MissesBeforeDeath int    `json:"misses_before_death"`
	NextBeacon        int    `json:"next_beacon"`
	Persist           bool   `json:"persist"`
}

// ServerInfo holds information about the server.
type ServerInfo struct {
	ServerID      string `json:"server_id"`
	ServerVersion string `json:"server_version"`
	ServerTime    int    `json:"server_time"`
}

// ---------------------------------------------------------------------
// Reconfiguration Information
// ---------------------------------------------------------------------

// ReconfigureInfo holds reconfiguration parameters sent from the server to the agent.
type ReconfigureInfo struct {
	AgentID           string `json:"agent_id"`
	BeaconInterval    int    `json:"beacon_interval"`
	MissesBeforeDeath int    `json:"misses_before_death"`
	Persist           bool   `json:"persist"`
}

// ---------------------------------------------------------------------
// Job Payload Types
// ---------------------------------------------------------------------

// JobPayloadPrint is a specific type of job payload for printing a message.
type JobPayloadPrint struct {
	Content string `json:"content"` // The message to print.
}

// JobPayloadPrintRsp is a specific type of job payload for printing a message response.
type JobPayloadPrintRsp struct {
	Success bool `json:"success"` // Whether the print job was successful.
}

// ---------------------------------------------------------------------
// Job Message Types
// ---------------------------------------------------------------------

// JobReq represents a job request message.
type JobReq struct {
	JobID      string          `json:"job_id"`
	JobType    JobType         `json:"job_type"`
	JobPayload json.RawMessage `json:"job_payload"`
}

// JobRsp represents a job response message.
type JobRsp struct {
	JobID      string          `json:"job_id"`
	JobType    JobType         `json:"job_type"`
	JobStatus  JobStatus       `json:"job_status"`
	JobPayload json.RawMessage `json:"job_payload"`
}

// ---------------------------------------------------------------------
// Beacon Message Types
// ---------------------------------------------------------------------

// BeaconReq represents a beacon request from the agent to the server.
type BeaconReq struct {
	AgentInfo AgentInfo `json:"agent_info"`
	Errors    int       `json:"errors"`
	JobRsp    []JobRsp  `json:"job_rsp"`
}

// BeaconRsp represents a beacon response from the server to the agent.
type BeaconRsp struct {
	ServerInfo  ServerInfo      `json:"server_info"`
	Reconfigure ReconfigureInfo `json:"reconfigure"`
	JobReq      []JobReq        `json:"job_req"`
}

// ---------------------------------------------------------------------
// Parsing Functions
// ---------------------------------------------------------------------

// ParseBeaconReq parses a JSON-encoded beacon request message and returns a BeaconReq struct.
// It returns an error if the parsing fails.
func ParseBeaconReq(data []byte) (*BeaconReq, error) {
	var req BeaconReq
	err := json.Unmarshal(data, &req)
	if err != nil {
		return nil, err
	}
	return &req, nil
}

// ParseBeaconRsp parses a JSON-encoded beacon response message and returns a BeaconRsp struct.
// It returns an error if the parsing fails.
func ParseBeaconRsp(data []byte) (*BeaconRsp, error) {
	var rsp BeaconRsp
	err := json.Unmarshal(data, &rsp)
	if err != nil {
		return nil, err
	}
	return &rsp, nil
}
