package main

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"

	"golang.org/x/net/html"
)

type Queue struct {
	elements []string
	mu sync.Mutex
}

func (q *Queue) enqueue(url string) {
	q.mu.Lock()
	defer q.mu.Unlock()
	q.elements = append(q.elements, url)
}

func (q *Queue) dequeue() string {
	q.mu.Lock()
	defer q.mu.Unlock()
	if len(q.elements) == 0 {
		return ""
	}
	firstUrl := q.elements[0]
	q.elements = q.elements[1:]
	return firstUrl
}

func (q *Queue) size() int {
	q.mu.Lock()
	defer q.mu.Unlock()
	return len(q.elements)
}

type CrawledSet struct {
	data map[string]bool
	mu sync.Mutex
}

type CrawledResult struct {
	data map[string]int
	mu sync.Mutex
}

func (c *CrawledResult) add(title string, numberOfLinks int) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.data[title] = numberOfLinks
}

func (c *CrawledSet) add(url string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.data[url] = true
}

func (c *CrawledSet) contains(url string) bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.data[url]
}

func (c *CrawledSet) size() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return len(c.data)
}

func getHref(token html.Token) (bool, string) {
	for _, a := range token.Attr {
		if a.Key == "href" {
			if len(a.Val) == 0 || !strings.HasPrefix(a.Val, "http") {
				return false, ""
			}
			return true, a.Val
		}
	}
	return false, ""
}

func parseHTML (content []byte) (title string, links []string) {
	tokenizer := html.NewTokenizer(bytes.NewReader(content))

	for {
		tokenType := tokenizer.Next()
		if tokenType == html.ErrorToken {
			break
		}

		token := tokenizer.Token()

		if(tokenType == html.StartTagToken) {
			if(token.Data == "title") {
				tokenizer.Next()
				titleToken := tokenizer.Token()
				title = titleToken.Data
			}

			if token.Data == "a" {
				if ok, href := getHref(token); ok {
					links = append(links, href)
				}
			}
		}
	}	
	return title, links
}

func fetchBody(url string) ([]byte, error) {
	res, err := http.Get(url)
	if(err != nil) {
		fmt.Printf("Error fetching %s: %v\n", url, err)
		return nil, err
	}

	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	return body, nil
}

func main() {
	url := "https://pkg.go.dev/golang.org/x/net/html"
	queue := &Queue{elements: make([]string, 0)}
	crawled := &CrawledSet{data: make(map[string]bool)}
	result := &CrawledResult{data: make(map[string]int)}
	queue.enqueue(url)
	maxPage := 5


	for(crawled.size() <= maxPage && queue.size() >0) {
		urlToCrawl := queue.dequeue()
		if(crawled.contains(urlToCrawl)) {
			continue
		}
		crawled.add(urlToCrawl)
		content, err := fetchBody(urlToCrawl)
		if(err != nil) {
			fmt.Printf("Error reading the response body: %v\n", err)
			return
		}
		title , links := parseHTML(content)
		result.add(title, len(links))

		for _, link := range links {
			if !crawled.contains(link) {
				queue.enqueue(link)
			}
		}
	}
	fmt.Printf("Crawled %d pages\n", maxPage)
	fmt.Print("Crawled result:\n")
	for key, value := range result.data {
    	fmt.Println(key, ":", value, "urls")
	}
}
