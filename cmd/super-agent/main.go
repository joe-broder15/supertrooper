package main

import (
	_ "embed"

	"github.com/joe-broder15/supertrooper/internal/agent"
)

// EMBED CERTIFICATES AND KEYS INTO

//go:embed embed/agent_cert.pem
var agentCertPEM []byte

//go:embed embed/agent_private_key.pem
var agentKeyPEM []byte

//go:embed embed/server_cert.pem
var serverCertPEM []byte

func main() {
	agent.Start(agentCertPEM, agentKeyPEM, serverCertPEM)
}
