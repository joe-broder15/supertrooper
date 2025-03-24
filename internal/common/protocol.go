package common

// ====================================================
// JSON DATASTRUCTURES FOR BEACON REQUESTS AND RESPONSE
// ====================================================

// information about host operating system for an agent
type BeaconReqSystemInfo struct {
	HostName   string
	InternalIP string
	Locale     string
	Version    string
}

// information about the agent itself to report or update
type AgentConfig struct {
	AgentID            string
	ServerAddr         string
	ServerPort         int
	BeaconInterval     int
	BeaconTries        int
	NextBeaconExpected int
	Persist            bool
}

// the entire structure sent to the server
type BeaconReq struct {
	AgentConfig  AgentConfig
	SystemInfo   BeaconReqSystemInfo
	JobResponses []any
}

type BeaconRsp struct {
	AgentConfig AgentConfig
	JobRequests []any
	Die         bool
}

// ====================================================
// DATASTRUCTURES FOR JOB REQUESTS AND RESPONSE
// ====================================================

// different job engines / methods of executing jobs. these are families
type JobEngine = int

const (
	JobEngineCore JobEngine = iota
	JobEnginePowerShell
)

// specific job types that can be executed by different engines
type JobType = int

const (
	JobTypeGetFile JobType = iota
)

// job priorities that indicate whether early callbacks are needed
type JobPriority = int

const (
	JobPriorityLow = iota
	JobPriorityHigh
)

// job statuses

type JobReq struct {
	JobID       string
	JobPriority JobPriority
	JobEngine   JobEngine
	JobType     JobType
	JobArgs     string
	JobPayload  string
}

type JobRsp struct {
	JobRequest JobReq
	Results    string
	Errors     []error
}
