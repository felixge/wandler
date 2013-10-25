package integration

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"testing"
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
	u, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}

	fmt.Printf("url: %s\n", u)
}
