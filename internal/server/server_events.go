package server

import (
	"crypto/tls"

	"github.com/joe-broder15/supertrooper/internal/messages"
)

type ServerEventType int

const (
	ServerEventTypeError ServerEventType = iota
	ServerEventTypeAgentC2Message
)

type ServerEvent struct {
	Type ServerEventType
	Body any
}

type ServerEventAgentC2Message struct {
	Message   messages.C2MessageBase
	AgentConn *tls.Conn
}

type ServerEventError struct {
	Error error
}
