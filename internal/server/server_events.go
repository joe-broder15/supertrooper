package server

import (
	"crypto/tls"

	"github.com/joe-broder15/supertrooper/internal/messages"
)

type ServerEventType int

const (
	ServerEventTypeError ServerEventType = iota
	ServerEventTypeBeaconReq
)

type ServerEvent struct {
	Type ServerEventType
	Body any
}

type ServerEventBeaconReq struct {
	BeaconReq messages.BeaconReq
	AgentConn *tls.Conn
}

type ServerEventError struct {
	Error error
}
