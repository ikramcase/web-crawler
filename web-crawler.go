package main

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"golang.org/x/net/html"
)

func main() {
	// startURL := "https://www.facebook.com" // giving all the domains on login page
	// startURL := "https://www.kdd-technologies.com" //error host
	// startURL := "https://parserdigital.com/" //giving all the subdomains
	startURL := "https://www.linkedin.com/" //giving all the subdomains

	domain := getDomain(startURL)
	visited := make(map[string]bool)

	crawl(startURL, domain, visited)
}

func crawl(url string, domain string, visited map[string]bool) {
	if visited[url] {
		return
	}

	visited[url] = true
	fmt.Println("Visited:", url)

	time.Sleep(2 * time.Second) // Wait for 2 seconds before making the HTTP request

	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer resp.Body.Close()

	doc, err := html.Parse(resp.Body)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	links := extractLinks(doc, domain)
	for _, link := range links {
		crawl(link, domain, visited)
	}
}

func extractLinks(body *html.Node, domain string) []string {
	var links []string

	var traverse func(*html.Node)
	traverse = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "a" {
			for _, attr := range n.Attr {
				if attr.Key == "href" {
					link, err := url.Parse(attr.Val)
					if err == nil && link.IsAbs() && link.Hostname() == domain && strings.HasPrefix(link.Path, "/") {
						links = append(links, link.String())
					}
				}
			}
		}

		for c := n.FirstChild; c != nil; c = c.NextSibling {
			traverse(c)
		}
	}

	traverse(body)
	return links
}

func getDomain(urlStr string) string {
	u, err := url.Parse(urlStr)
	if err != nil {
		panic(err)
	}
	return u.Hostname()
}
