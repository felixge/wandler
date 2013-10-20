package queue

import (
	"strings"
	"testing"
)

var newQueueTests = []struct {
	Dsn       string
	ExpectErr string
	Verify    func(*testing.T, Interface)
}{
	{"foo://bar", "unknown schema: foo", nil},
	{"redis://myhost:12345/foo", "", func(t *testing.T, queue Interface) {
		q, ok := queue.(*redisQueue)
		if !ok {
			t.Errorf("wrong type: %#v", q)
			return
		}

		if q.host != "myhost:12345" {
			t.Errorf("wrong host: %s", q.host)
			return
		}

		if q.prefix != "foo" {
			t.Errorf("wrong prefix: %s", q.prefix)
			return
		}
	}},
}

func TestNewQueue(t *testing.T) {
	for _, test := range newQueueTests {
		queue, err := NewQueue(test.Dsn, nil)
		if test.ExpectErr != "" {
			if err == nil {
				t.Errorf("Expected err, got: %#v", queue)
			} else if !strings.Contains(err.Error(), test.ExpectErr) {
				t.Errorf("Expected err: %s, got: %#v", test.ExpectErr, err)
			}
			continue
		}

		test.Verify(t, queue)
	}
}
