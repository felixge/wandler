package queue

// @TODO: config options for dial/read/write/dequeue timeouts (maybe one
// timeout value suffices)
// @TODO: support for retries
// @TODO: more reliable queueing, acks, rejects, etc.
// @TODO: support for HA
// @TODO: use abstraction around redis to allow for pooling in future
// @TODO: support for concurrent enqueue/dequeue

import (
	"fmt"
	"github.com/garyburd/redigo/redis"
	"io"
	"net/url"
	"strings"
	"sync"
)

func newRedisQueue(connURL *url.URL) (*redisQueue, error) {
	host := connURL.Host
	if !strings.Contains(host, ":") {
		host += ":6379"
	}
	prefix := strings.TrimPrefix(connURL.Path, "/")

	q := &redisQueue{
		host:   host,
		prefix: prefix,
	}
	return q, nil
}

type redisQueue struct {
	lock      sync.Mutex
	host      string
	prefix    string
	redisConn redis.Conn
}

func (q *redisQueue) Enqueue(msg string) error {
	q.lock.Lock()
	defer q.lock.Unlock()

	c, err := q.conn()
	if err != nil {
		return err
	}
	_, err = c.Do("RPUSH", q.key(), msg)
	if err != nil {
		return err
	}
	return nil
}

func (q *redisQueue) Dequeue() (string, error) {
	q.lock.Lock()
	defer q.lock.Unlock()

	c, err := q.conn()
	if err != nil {
		return "", err
	}

	r, err := redis.Strings(c.Do("BLPOP", q.key(), "1"))
	if err != nil {
		return "", err
	}

	if r == nil {
		return "", io.EOF
	} else if len(r) != 2 {
		return "", fmt.Errorf("redis queue: Expected 2 elements, got: %#v", r)
	}
	return r[1], nil
}

func (q *redisQueue) key() string {
	key := "queue"
	if q.prefix == "" {
		return key
	}
	return q.prefix + "." + key
}

func (q *redisQueue) conn() (redis.Conn, error) {
	if q.redisConn != nil {
		return q.redisConn, nil
	}
	return redis.Dial("tcp", q.host)
}
