
</head>
<body>
    <h1>Go Job Scraper</h1>
    
  <p>A concurrent web scraper built in Go that extracts job listings from remote job websites.</p>

   <h2>Features</h2>
    <ul>
        <li>Concurrent scraping of multiple job websites</li>
        <li>Robust error handling and timeouts</li>
        <li>HTML entity cleaning and text normalization</li>
        <li>Multiple selector patterns for different site layouts</li>
        <li>Rate limiting to prevent IP blocking</li>
        <li>Browser-like headers to avoid request filtering</li>
    </ul>

  <h2>Prerequisites</h2>
    <ul>
        <li>Go 1.16 or higher</li>
        <li>Required packages:
            <ul>
                <li>github.com/PuerkitoBio/goquery</li>
                <li>golang.org/x/net/html</li>
                <li>golang.org/x/sync/errgroup</li>
            </ul>
        </li>
    </ul>

   <h2>Installation</h2>
    <pre><code>go get github.com/PuerkitoBio/goquery
go get golang.org/x/net/html
go get golang.org/x/sync/errgroup</code></pre>

  <h2>Usage</h2>
    <pre><code>make run</code></pre>

   <h2>Code Structure</h2>
    
  <h3>Job Structure</h3>
    <p>The scraper stores job information in the following structure:</p>
    <pre><code>type Job struct {
    Title       string
    Company     string
    Location    string
    Description string
}</code></pre>

   <h3>Main Components</h3>
    <ol>
        <li><strong>cleanText(s string) string</strong>
            <ul>
                <li>Removes HTML tags</li>
                <li>Normalizes whitespace</li>
                <li>Converts HTML entities to characters</li>
            </ul>
        </li>
        <li><strong>scrapeJobListings(ctx context.Context, url string) ([]Job, error)</strong>
            <ul>
                <li>Makes HTTP requests with browser-like headers</li>
                <li>Parses HTML content</li>
                <li>Extracts job information using CSS selectors</li>
                <li>Implements rate limiting</li>
            </ul>
        </li>
    </ol>

  <h2>Configuration</h2>
  <p>The following settings can be modified in the code:</p>
    <ul>
        <li>HTTP client timeout (default: 30 seconds)</li>
        <li>Rate limiting delay (default: 2 seconds)</li>
        <li>Target URLs (in the urls slice)</li>
    </ul>

   <h2>Example Output</h2>
    <pre><code>Scraping URL: https://remoteok.com/remote-software-jobs

Found job #1:
Title: Senior Frontend Developer
Company: TechCorp
Location: Remote
Description: We are looking for an experienced frontend developer...</code></pre>

   <div class="warning">
        <h3>⚠️ Important Notes</h3>
        <ul>
            <li>Respect the websites' robots.txt and terms of service</li>
            <li>Adjust rate limiting as needed to avoid IP blocking</li>
            <li>The scraper assumes specific HTML structure; site changes may require code updates</li>
            <li>Some websites may block automated scraping attempts</li>
        </ul>
    </div>

  <h2>Error Handling</h2>
    <p>The scraper handles several types of errors:</p>
    <ul>
        <li>Network timeouts</li>
        <li>Invalid HTML responses</li>
        <li>Non-200 HTTP status codes</li>
        <li>Context cancellation</li>
    </ul>

   <h2>Contributing</h2>
    <p>Contributions are welcome! Please feel free to submit a Pull Request. For major changes, please open an issue first to discuss what you would like to change.</p>

 
</body>
</html>
