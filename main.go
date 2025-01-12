package main

import (
	"fmt"
	"io"
	"log"
	"log/slog"
	"net/http"
	"time"

	"github.com/anthdm/hollywood/actor"
	"golang.org/x/net/html"
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
	url := "https://linkedin.com.br"
	resp, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(exctractLinks(resp.Body))

	e, err := actor.NewEngine(actor.NewEngineConfig())
	if err != nil {
		log.Fatal(err)
	}
	pid := e.Spawn(NewManager(), "manager")

	time.Sleep(time.Microsecond * 500)

	e.Send(pid, VisitorRequest{links: []string{"https://linkedin.com"}})
	time.Sleep(time.Second * 10)
}

func exctractLinks(body io.Reader) []string {
	links := make([]string, 0)
	tokenizer := html.NewTokenizer(body)

	for {
		tokenType := tokenizer.Next()
		if tokenType == html.ErrorToken {
			return links
		}
		if tokenType == html.StartTagToken {
			token := tokenizer.Token()
			if token.Data == "a" {
				for _, attr := range token.Attr {
					if attr.Key == "href" {
						links = append(links, attr.Val)
					}
				}
			}
		}
	}
}
