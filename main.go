package main

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"strings"

	"golang.org/x/net/html"
)

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

func main() {
	url := "https://pkg.go.dev/golang.org/x/net/html"

	res, err := http.Get(url)
	if(err != nil) {
		fmt.Printf("Error fetching %s: %v\n", url, err)
		return
	}

	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)

	if(err != nil) {
		fmt.Printf("Error reading the response body: %v\n", err)
		return
	}

	title, links := parseHTML(body)

	fmt.Printf("Title: %s\n", title)
	fmt.Printf("Found %d links\n", len(links))

	for i, link := range links {
		if (i >= 10) {
			fmt.Printf("...and %d more links\n", len(links) - 10)
			break
		}
		fmt.Printf(" %d: %s\n", i + 1, link)
	}
}
