package main

import (
	"fmt"
	"io"
	"net/http"
)

func main() {
	url := "https://go.dev/learn/"

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

	fmt.Printf("Successfully fetching %s\n", url)
	fmt.Printf("Response status: %s\n", res.Status)
	fmt.Printf("Response body length: %d\n", len(body))

	if(len(body) > 200) {
		fmt.Printf("First 200 characters of the body: %s\n", string(body[:200]))
	} else {
		fmt.Printf("Response body: %s\n", string(body))
	}
}
