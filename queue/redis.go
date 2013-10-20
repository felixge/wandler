package queue

import (
	"github.com/adeven/redismq"
	"net/url"
	"strings"
)

func newRedisQueue(connURL *url.URL) (*redisQueue, error) {
	host := connURL.Host
	if !strings.Contains(host, ":") {
		host += ":6379"
	}
	prefix := strings.TrimPrefix(connURL.Path, "/")

	rmq := redismq.NewQueue(host, "", 0, prefix)
	consumer, err := rmq.AddConsumer("bob")
	if err != nil {
		return nil, err
	}

	q := &redisQueue{
		host:     host,
		prefix:   prefix,
		rmq:      rmq,
		consumer: consumer,
	}
	return q, nil
}

type redisQueue struct {
	host     string
	prefix   string
	rmq      *redismq.Queue
	consumer *redismq.Consumer
}

func (q *redisQueue) Enqueue(msg string) error {
	return q.rmq.Put(msg)
}

func (q *redisQueue) Dequeue() (string, error) {
	pkg, err := q.consumer.Get()
	if err != nil {
		return "", err
	}

	return pkg.Payload, pkg.Ack()
}
