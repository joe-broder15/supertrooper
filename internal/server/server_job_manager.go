package server

import (
	"sync"

	"github.com/joe-broder15/supertrooper/internal/common"
)

// entries in the job manager, this tracks which jobs have been accepted and completed
type agentJobQueues struct {
	completedJobs map[string]common.JobRsp //jobs that have returned from an agent
	taskedJobs    map[string]common.JobReq //jobs that have been sent to an agent but have not returned
	untaskedJobs  map[string]common.JobReq //jobs that have not been sent to an agent yet
}

func newAgentJobQueues() *agentJobQueues {
	return &agentJobQueues{
		completedJobs: make(map[string]common.JobRsp),
		taskedJobs:    make(map[string]common.JobReq),
		untaskedJobs:  make(map[string]common.JobReq),
	}
}

// the job manager is just a map of agent ids to the above structs
type jobManager struct {
	lock        sync.RWMutex
	agentQueues map[string]*agentJobQueues
}

func newJobManager() *jobManager {
	return &jobManager{
		agentQueues: make(map[string]*agentJobQueues),
	}
}

// add a job to an agent's queue of unaccepted jovs
func (jm *jobManager) addJob(agentID string, job common.JobReq) {
	// acquire lock
	jm.lock.Lock()
	defer jm.lock.Unlock()

	// check if the agent aleady has job queues
	if _, ok := jm.agentQueues[agentID]; !ok {
		jm.agentQueues[agentID] = newAgentJobQueues()
	}

	// add the job to the agent's untasked jobs
	jm.agentQueues[agentID].untaskedJobs[job.JobID] = job
}

// task unaccepted jobs and return the newly tasked jobs
func (jm *jobManager) taskUntaskedJobs(agentID string) []common.JobReq {
	// acquire lock
	jm.lock.Lock()
	defer jm.lock.Unlock()

	// make a slice to store all of the jobs we move over to tasked status
	newlyTasked := make([]common.JobReq, 0)

	// check if the agent aleady has job queues, otherwise return an empty slice
	if _, ok := jm.agentQueues[agentID]; !ok {
		return newlyTasked
	}

	// get the queues for this agent
	queues := jm.agentQueues[agentID]

	// promote all of the untasked jobs to tasked jovs
	for jobID, jobReq := range queues.untaskedJobs {
		// copy the untasked job to the output slice
		newlyTasked = append(newlyTasked, jobReq)
		// add the new job to the tasked queue
		queues.taskedJobs[jobID] = jobReq
		// delete the new job from the untasked queue
		delete(queues.untaskedJobs, jobID)
	}

	return newlyTasked
}

// register completed jobs
func (jm *jobManager) registerCompletedJobs(agentID string, jobs *[]common.JobRsp) {
	// acquire lock
	jm.lock.Lock()
	defer jm.lock.Unlock()

	// check if the agent aleady has job queues, otherwise return an empty slice
	if _, ok := jm.agentQueues[agentID]; !ok {
		return
	}

	// get the queues for this agent
	queues := jm.agentQueues[agentID]

	// move all of the matching tasked jobs to completed status
	for _, job := range *jobs {
		if _, ok := queues.taskedJobs[job.JobReq.JobID]; ok {
			queues.completedJobs[job.JobReq.JobID] = job
			delete(queues.taskedJobs, job.JobReq.JobID)
		}
	}
}

// get completed jobs
func (jm *jobManager) getCompletedJobs(agentID string) []common.JobRsp {
	// acquire lock
	jm.lock.RLock()
	defer jm.lock.RUnlock()

	// check if the agent aleady has job queues
	if _, ok := jm.agentQueues[agentID]; !ok {
		return make([]common.JobRsp, 0)
	}

	// convert map to slice and return
	jobs := make([]common.JobRsp, 0, len(jm.agentQueues[agentID].completedJobs))
	for _, job := range jm.agentQueues[agentID].completedJobs {
		jobs = append(jobs, job)
	}
	return jobs
}

// get tasked jobs
func (jm *jobManager) getTaskedJobs(agentID string) []common.JobReq {
	// acquire lock
	jm.lock.RLock()
	defer jm.lock.RUnlock()

	// check if the agent aleady has job queues
	if _, ok := jm.agentQueues[agentID]; !ok {
		return make([]common.JobReq, 0)
	}

	// convert map to slice and return
	jobs := make([]common.JobReq, 0, len(jm.agentQueues[agentID].taskedJobs))
	for _, job := range jm.agentQueues[agentID].taskedJobs {
		jobs = append(jobs, job)
	}
	return jobs
}

// get untasked jobs
func (jm *jobManager) getUntaskedJobs(agentID string) []common.JobReq {
	// acquire lock
	jm.lock.RLock()
	defer jm.lock.RUnlock()

	// check if the agent aleady has job queues
	if _, ok := jm.agentQueues[agentID]; !ok {
		return make([]common.JobReq, 0)
	}

	// convert map to slice and return
	jobs := make([]common.JobReq, 0, len(jm.agentQueues[agentID].untaskedJobs))
	for _, job := range jm.agentQueues[agentID].untaskedJobs {
		jobs = append(jobs, job)
	}
	return jobs
}

// is empty function
func (jm *jobManager) isEmpty(agentID string) bool {
	// acquire lock
	jm.lock.RLock()
	defer jm.lock.RUnlock()

	// check if the agent aleady has job queues
	if _, ok := jm.agentQueues[agentID]; !ok {
		return true
	}

	// get the queues for this agent
	queues := jm.agentQueues[agentID]

	// check if all queues are empty
	return len(queues.completedJobs) == 0 &&
		len(queues.taskedJobs) == 0 &&
		len(queues.untaskedJobs) == 0
}

// clean empty queues function
func (jm *jobManager) cleanEmptyQueues() {
	// acquire lock
	jm.lock.Lock()
	defer jm.lock.Unlock()

	// iterate through all agent queues
	for agentID, queues := range jm.agentQueues {
		// if all queues are empty, delete the agent's entry
		if len(queues.completedJobs) == 0 &&
			len(queues.taskedJobs) == 0 &&
			len(queues.untaskedJobs) == 0 {
			delete(jm.agentQueues, agentID)
		}
	}
}
