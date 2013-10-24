package integration

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"testing"
	"time"
)

func TestImageResize(t *testing.T) {
	worker, err := NewWorker()
	if err != nil {
		t.Fatal(err)
	}
	defer worker.Kill()

	server, err := NewServer()
	if err != nil {
		t.Fatal(err)
	}
	defer server.Kill()

	data := url.Values{
		"src":    {"http://felixge.de/img/get-on-the-squirrel.png"},
		"width":  {"200"},
		"height": {"100"},
	}
	resp, err := http.PostForm("http://"+server.HttpAddr()+"/image", data)
	if err != nil {
		t.Fatal(err)
	}

	fmt.Printf("headers: %#v\n", resp.Header)
	defer resp.Body.Close()
	_, err = io.Copy(os.Stdout, resp.Body)
	if err != nil {
		t.Fatal(err)
	}

	time.Sleep(10*time.Second)
}
