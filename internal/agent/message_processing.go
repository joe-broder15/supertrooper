package agent

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/joe-broder15/supertrooper/internal/messages"
)

// processPrintJob executes a print job defined in the JSON payload and returns a JSON-encoded response.
func processPrintJob(printJobJSON json.RawMessage) (json.RawMessage, error) {
	// Decode the incoming JSON into the JobPayloadPrint structure.
	var job messages.JobPayloadPrint
	if err := json.Unmarshal(printJobJSON, &job); err != nil {
		return nil, fmt.Errorf("processPrintJob: unable to unmarshal job payload: %w", err)
	}

	// Attempt to print the job content to standard output.
	_, err := fmt.Println(job.Content)

	// Build the response: set Success to true if printing did not error.
	response := messages.JobPayloadPrintRsp{
		Success: err == nil,
	}

	// Encode and return the response into JSON.
	responseJSON, err := json.Marshal(response)
	if err != nil {
		return nil, fmt.Errorf("processPrintJob: failed to marshal response: %w", err)
	}
	return responseJSON, nil
}

func processJobDispatch(payload json.RawMessage) (json.RawMessage, error) {
	// Initialize the default response with an error status.
	response := messages.C2MessageTypeJobDispatchRspPayload{
		JobID:     "none",
		JobType:   0,
		JobStatus: messages.JobStatusError,
	}

	// Parse the incoming dispatch payload.
	jobDispatch, err := messages.ParseJobDispatchPayload(payload)
	if err != nil {
		// Record the error in the response.
		response.JobError = err.Error()
	} else {
		// Update the response job ID.
		response.JobID = jobDispatch.JobID

		// Process the job based on its type.
		switch jobDispatch.JobType {
		case messages.JobTypePrint:
			// Process the print job and update the response accordingly.
			processedPayload, err := processPrintJob(jobDispatch.JobPayload)
			if err != nil {
				response.JobError = err.Error()
				response.JobStatus = messages.JobStatusError
			} else {
				response.JobType = messages.JobTypePrint
				response.JobStatus = messages.JobStatusSuccess
				response.JobPayload = processedPayload
			}
		default:
			// Unknown job type: record an error.
			response.JobError = fmt.Sprintf("unknown job type: %v", jobDispatch.JobType)
			response.JobStatus = messages.JobStatusError
		}
	}

	// Encode the response structure into JSON.
	responseJSON, err := json.Marshal(response)
	if err != nil {
		return nil, fmt.Errorf("processJobDispatch: failed to marshal response: %w", err)
	}
	return responseJSON, nil
}

// processC2Messages routes C2 messages to the appropriate handler based on the message type.
func processC2Messages(agentState *AgentState, message messages.C2MessageBase) (messages.C2MessageBase, error) {
	// Initialize the response with the agent ID and current timestamp.
	response := messages.C2MessageBase{
		AgentID:   agentState.id, // Consider renaming agentState.id -> agentState.ID for idiomatic Go.
		Timestamp: time.Now().Unix(),
	}
	var processErr error

	// Dispatch processing based on the message Type.
	switch message.Type {
	case messages.C2MessageTypeJobDispatch:
		processedPayload, err := processJobDispatch(message.C2MessagePayload)
		if err != nil {
			processErr = err
		}
		response.C2MessagePayload = processedPayload
		response.Type = messages.C2MessageTypeJobDispatchRsp

	default:
		return messages.C2MessageBase{}, fmt.Errorf("processC2Messages: unhandled message type: %v", message.Type)
	}

	return response, processErr
}
