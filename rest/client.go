package rest

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type Paginator[T any, D any] interface {
	HasNextPage() bool
	NextPageURL() string
	SetNextPage(*T)
	ExtractData(*T) <-chan *D
}

type Client[T any, D any] struct {
	c     *http.Client
	p     Paginator[T, D]
	items chan *D
}

func NewClient[T any, D any](c *http.Client, p Paginator[T, D]) *Client[T, D] {
	return &Client[T, D]{c: c, p: p}
}

func (c *Client[T, D]) Items() <-chan *D {
	c.items = make(chan *D)

	go func() {
		for c.p.HasNextPage() {
			page, err := c.getNextPage()
			if err != nil {
				log.Printf("failed fetching page: %v", err)
				close(c.items)
				break
			}

			for d := range c.p.ExtractData(page) {
				c.items <- d
			}
		}

		close(c.items)
	}()

	return c.items
}

func (c *Client[T, D]) getNextPage() (*T, error) {
	res, err := c.c.Get(c.p.NextPageURL())
	if err != nil {
		return nil, fmt.Errorf("get: %w", err)
	}

	defer res.Body.Close()

	var page T

	if err = json.NewDecoder(res.Body).Decode(&page); err != nil {
		return nil, fmt.Errorf("decode: %w", err)
	}

	c.p.SetNextPage(&page)

	return &page, nil
}
