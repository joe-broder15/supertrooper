package common

import "encoding/json"

// ====================================================
// DATASTRUCTURES FOR BEACON REQUESTS AND RESPONSE
// ====================================================

// information about host operating system for an agent
type AgentHostInfo struct {
	HostName   string
	InternalIP string
	Locale     string
	Version    string
}

// information about the agent itself to report or update
type AgentConfig struct {
	AgentID        string
	ServerAddr     string
	BeaconInterval int
	BeaconTries    int
	Persist        bool
}

// the entire structure sent to the server
type BeaconReq struct {
	AgentConfig        AgentConfig
	AgentHostInfo      AgentHostInfo
	NextBeaconExpected int
	JobResponses       []JobRsp
}

type BeaconRsp struct {
	ServerID    string
	JobRequests []JobReq
}

// ====================================================
// DATASTRUCTURES FOR JOB REQUESTS AND RESPONSE
// ====================================================

// specific job types that can be executed by different engines
type JobType = int

const (
	JobTypeReconfigure JobType = iota
	JobTypeUninstall
	JobTypeExecPSPlugin
)

// job priorities that indicate whether early callbacks are needed
type JobPriority = int

const (
	JobPriorityLow = iota
	JobPriorityHigh
)

// job statuses for completed jobs
type JobStatus = int

const (
	JobStatusFailed = iota
	JobStatusSuccess
)

// actual structures for job requests and responses
type JobReq struct {
	JobID       string
	JobPriority JobPriority
	JobType     JobType
	JobPayload  json.RawMessage
}

// job args for reconfigure
type JobPayloadReconfigure struct {
	AgentConfig AgentConfig
}

// job args for executing powershell plugins
type JobPayloadExecPSPlugin struct {
	PluginType int
	Args       string
	Script     string
}

type JobRsp struct {
	JobReq    JobReq
	Results   string
	JobStatus JobStatus
	Errors    []error
}
