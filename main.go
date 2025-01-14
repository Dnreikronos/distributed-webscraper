package main

import (
	"bytes"
	"context"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"golang.org/x/net/html"
	"golang.org/x/sync/errgroup"
	"log"
	"log/slog"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/anthdm/hollywood/actor"
	"golang.org/x/net/html"
)

type Job struct {
	Title       string
	Company     string
	Location    string
	Description string
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

func main() {
	ctx := context.Background()
	g, ctx := errgroup.WithContext(ctx)

	urls := []string{
		"https://remoteok.com/remote-software-jobs",
	}

	for _, url := range urls {
		url := url
		g.Go(func() error {
			fmt.Printf("Scraping URL: %s\n", url)

			jobs, err := scrapeJobListings(ctx, url)
			if err != nil {
				return fmt.Errorf("failed to scrape %s: %w", url, err)
			}

			if len(jobs) == 0 {
				fmt.Println("No jobs found. The site structure might have changed or blocking requests.")
				return nil
			}

			fmt.Printf("\nFound %d jobs:\n\n", len(jobs))
			for i, job := range jobs {
				fmt.Printf("Job #%d:\n", i+1)
				fmt.Printf("Title: %s\nCompany: %s\nLocation: %s\nDescription: %.200s...\n\n",
					job.Title,
					job.Company,
					job.Location,
					job.Description)
			}
			return nil
		})
	}

	if err := g.Wait(); err != nil {
		log.Fatal("Error during scraping:", err)
	}
}
