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

func (m *Manager) handleVisitRequest(msg VisitorRequest) error {
	for _, link := range msg.links {
		slog.Info("Visiting urls", "url", link)
	}
	return nil
}
func main() {
	e, err := actor.NewEngine(actor.NewEngineConfig())
	if err != nil {
		log.Fatal(err)
	}
	pid := e.Spawn(NewManager(), "manager")
	time.Sleep(time.Second * 10)
}
