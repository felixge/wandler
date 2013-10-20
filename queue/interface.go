package queue

import (
	"fmt"
	"github.com/felixge/wandler/log"
	"net/url"
)

type Interface interface {
	Enqueue(msg string) error
	Dequeue() (string, error)
}

func NewQueue(connURL string, l log.Interface) (Interface, error) {
	u, err := url.Parse(connURL)
	if err != nil {
		return nil, err
	}

	switch u.Scheme {
	case "redis":
		return newRedisQueue(u)
	default:
		return nil, fmt.Errorf("unknown schema: %s", u.Scheme)
	}
}
