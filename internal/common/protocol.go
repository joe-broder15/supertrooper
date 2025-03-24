package common

// ====================================================
// DATASTRUCTURES FOR BEACON REQUESTS AND RESPONSE
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
	JobResponses []JobRsp
}

type BeaconRsp struct {
	AgentConfig AgentConfig
	JobRequests []JobReq
	Die         bool
}

// ====================================================
// DATASTRUCTURES FOR JOB REQUESTS AND RESPONSE
// ====================================================

// specific job types that can be executed by different engines
type JobType = int

const (
	JobTypeGetFile JobType = iota
	JobTypePutFile
	JobTypeDirList
	JobTypeSurvey
	JobTypeExecCmd
)

// job priorities that indicate whether early callbacks are needed
type JobPriority = int

const (
	JobPriorityLow = iota
	JobPriorityHigh
)

// actual structures for job requests and responses
type JobReq struct {
	JobID       string
	JobPriority JobPriority
	JobType     JobType
	JobArgs     string
	JobPayload  string
}

type JobRsp struct {
	JobRequest JobReq
	Results    string
	Errors     []error
}
