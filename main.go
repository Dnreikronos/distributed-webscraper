package main

import (
	"bytes"
	"context"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"golang.org/x/net/html"
	"golang.org/x/sync/errgroup"
	"log"
	"net/http"
	"regexp"
	"strings"
	"time"

)

type Job struct {
	Title       string
	Company     string
	Location    string
	Description string
}

func cleanText(s string) string {
	doc, err := html.Parse(strings.NewReader(s))
	if err != nil {
		return strings.TrimSpace(s)
	}

	var buf bytes.Buffer
	var extractText func(*html.Node)
	extractText = func(n *html.Node) {
		if n.Type == html.TextNode {
			buf.WriteString(n.Data)
			buf.WriteString(" ")
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			extractText(c)
		}
	}
	extractText(doc)

	text := buf.String()
	space := regexp.MustCompile(`\s+`)
	text = space.ReplaceAllString(text, " ")

	text = strings.ReplaceAll(text, "&nbsp;", " ")
	text = strings.ReplaceAll(text, "&amp;", "&")
	text = strings.ReplaceAll(text, "&lt;", "<")
	text = strings.ReplaceAll(text, "&gt;", ">")
	text = strings.ReplaceAll(text, "&quot;", "\"")

	return strings.TrimSpace(text)
}

func scrapeJobListings(ctx context.Context, url string) ([]Job, error) {
	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9")
	req.Header.Set("Cache-Control", "no-cache")
	req.Header.Set("Connection", "keep-alive")

	time.Sleep(2 * time.Second)

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch URL: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("received non-200 response code: %d", resp.StatusCode)
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to parse HTML: %w", err)
	}

	var jobs []Job

	doc.Find(".job").Each(func(i int, s *goquery.Selection) {
		title := cleanText(s.Find(".company_and_position h2").Text())
		company := cleanText(s.Find(".company_and_position h3").Text())
		location := cleanText(s.Find(".location").Text())
		description := cleanText(s.Find(".description").Text())

		fmt.Printf("\nFound job #%d:\nTitle: %s\nCompany: %s\nLocation: %s\nDescription preview: %.100s...\n",
			i+1, title, company, location, description)

		if title != "" || company != "" || location != "" || description != "" {
			job := Job{
				Title:       title,
				Company:     company,
				Location:    location,
				Description: description,
			}
			jobs = append(jobs, job)
		}
	})

	if len(jobs) == 0 {
		doc.Find(".job-listing").Each(func(i int, s *goquery.Selection) {
			title := cleanText(s.Find(".job-title").Text())
			company := cleanText(s.Find(".company-name").Text())
			location := cleanText(s.Find(".location").Text())
			description := cleanText(s.Find(".job-description").Text())

			if title != "" || company != "" || location != "" || description != "" {
				job := Job{
					Title:       title,
					Company:     company,
					Location:    location,
					Description: description,
				}
				jobs = append(jobs, job)
			}
		})
	}

	return jobs, nil
}

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
