package integration

import (
	"fmt"
	"github.com/felixge/wandler/queue"
	"testing"
)

func TestEnqueueDequeue(t *testing.T) {
	q, err := queue.NewQueue("redis://localhost/wandler.test", nil)
	if err != nil {
		t.Fatal(err)
	}

	for i := 0; i < 100; i++ {
		input := fmt.Sprintf("foo:%d", i)
		if err := q.Enqueue(input); err != nil {
			t.Fatal(err)
		}

		if output, err := q.Dequeue(); err != nil {
			t.Fatal(err)
		} else if output != input {
			t.Fatalf("wrong message: %s != %s", output, input)
		}
	}
}
