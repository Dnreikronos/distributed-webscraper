package main

import (
	"log"
	"log/slog"
	"time"

	"github.com/anthdm/hollywood/actor"
)

type VisitorRequest struct {
	links []string
}

type Manager struct{}

func NewManager() actor.Producer {
	return func() actor.Receiver {
		return &Manager{}
	}
}

func (m *Manager) Receive(c *actor.Context) {
	switch msg := c.Message().(type) {
	case VisitorRequest:
		m.handleVisitRequest(msg)
	case actor.Started:
		slog.Info("Manager started")
	case actor.Stopped:
	}
}

