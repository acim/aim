package main

import (
	"fmt"
	"net/http"

	"go.acim.net/aim/rest"
)

type page struct {
	CurrentPage  int    `json:"current_page"`
	Data         []fact `json:"data"`
	FirstPageURL string `json:"first_page_url"`
	From         int    `json:"from"`
	LastPage     int    `json:"last_page"`
	LastPageURL  string `json:"last_page_url"`
	Links        []struct {
		URL    interface{} `json:"url"`
		Label  string      `json:"label"`
		Active bool        `json:"active"`
	} `json:"links"`
	NextPageURL string      `json:"next_page_url"`
	Path        string      `json:"path"`
	PerPage     uint        `json:"per_page"`
	PrevPageURL interface{} `json:"prev_page_url"`
	To          int         `json:"to"`
	Total       uint        `json:"total"`
}

type fact struct {
	Fact   string `json:"fact"`
	Length int    `json:"length"`
}

var _ rest.Paginator[page, string] = (*paginator)(nil)

type Client struct {
	*rest.Client[page, string]
}

func NewClient(c *http.Client) *Client {
	return &Client{Client: rest.NewClient[page, string](c, &paginator{page: 1, hasMore: true})}
}

type paginator struct {
	page    uint
	hasMore bool
}

func (c *paginator) HasNextPage() bool {
	return c.hasMore
}
func (c *paginator) NextPageURL() string {
	return fmt.Sprintf("https://catfact.ninja/facts?page=%d", c.page)
}

func (c *paginator) SetNextPage(page *page) {
	if page.NextPageURL == "" {
		c.hasMore = false

		return
	}

	c.page++
}

func (c *paginator) ExtractData(page *page) <-chan *string {
	ch := make(chan *string)

	go func() {
		for i := range page.Data {
			ch <- &page.Data[i].Fact
		}

		close(ch)
	}()

	return ch
}
