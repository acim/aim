package main

import (
	"fmt"
	"net/http"
)

func main() {
	counter := 1

	c := NewClient(http.DefaultClient)

	for item := range c.Items() {
		fmt.Printf("%d. %s\n", counter, *item)
		counter++
	}
}
