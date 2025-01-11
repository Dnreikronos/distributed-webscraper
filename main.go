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

type Visitor struct {
	URL string
}

func NewVisitor(url string) actor.Producer {
	return func() actor.Receiver {
		return &Visitor{
			URL: url,
		}
	}
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
		m.handleVisitRequest(c, msg)
	case actor.Started:
		slog.Info("Manager started")
	case actor.Stopped:
	}
}

func (v *Visitor) Receive(c *actor.Context) {
	switch c.Message().(type) {
	case actor.Started:
		slog.Info("visitor started", "url", v.URL)
	case actor.Stopped:
	}
}

func (m *Manager) handleVisitRequest(c *actor.Context, msg VisitorRequest) error {
	for _, link := range msg.links {
		slog.Info("Visiting urls", "url", link)
		c.SpawnChild(NewVisitor(link), "visitor/"+link)
	}
	return nil
}

func main() {
	e, err := actor.NewEngine(actor.NewEngineConfig())
	if err != nil {
		log.Fatal(err)
	}
	pid := e.Spawn(NewManager(), "manager")

	time.Sleep(time.Microsecond * 500)

	e.Send(pid, VisitorRequest{links: []string{"https://linkedin.com"}})
	time.Sleep(time.Second * 10)
}
