package server

import (
	"sync"
	"time"
)

type AgentManager struct {
	mutex  sync.Mutex
	agents map[string]AgentInfo
}

type AgentInfo struct {
	ID                string
	IpAddress         string
	Hostname          string
	IsLive            bool
	IsPersistent      bool
	BeaconSeconds     int
	MissesBeforeDeath int
	LastContact       time.Duration
}

func NewAgentManager() *AgentManager {
	return &AgentManager{
		agents: make(map[string]AgentInfo),
	}
}
