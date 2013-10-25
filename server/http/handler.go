package http

import (
	"encoding/json"
	"github.com/felixge/wandler/job"
	"github.com/felixge/wandler/log"
	"github.com/felixge/wandler/queue"
	"github.com/nu7hatch/gouuid"
	"io"
	"net/http"
	"regexp"
	"strconv"
	"sync"
)

const multipartMaxMemory = 64 * 1024

func NewHandler(c HandlerConfig) (*Handler, error) {
	return &Handler{
		log:     c.Log,
		queue:   c.Queue,
		httpURL: c.HttpURL,
		pending: map[string]chan string{},
	}, nil
}

type HandlerConfig struct {
	Log     log.Interface
	Queue   queue.Interface
	HttpURL string
}

type Handler struct {
	lock    sync.Mutex
	log     log.Interface
	queue   queue.Interface
	httpURL string
	pending map[string]chan string
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.log.Info(
		"Incoming http request method=%s url=%s length=%d addr=%s",
		r.Method,
		r.URL,
		r.ContentLength,
		r.RemoteAddr,
	)

	if err := r.ParseMultipartForm(multipartMaxMemory); err != nil {
		h.log.Warn("Bad http request: %s", err)
		w.WriteHeader(http.StatusBadRequest)
		io.WriteString(w, err.Error())
		return
	}

	if m := regexp.MustCompile("/notify/([a-z0-9-]+)").FindStringSubmatch(r.URL.Path); len(m) == 2 {
		h.serveNotification(m[1], w, r)
		return
	}

	width, err := strconv.ParseInt(r.Form.Get("width"), 10, 32)
	if err != nil {
		h.log.Warn("Bad http request: %s", err)
		w.WriteHeader(http.StatusBadRequest)
		io.WriteString(w, err.Error())
		return
	}
	height, err := strconv.ParseInt(r.Form.Get("width"), 10, 32)
	if err != nil {
		h.log.Warn("Bad http request: %s", err)
		w.WriteHeader(http.StatusBadRequest)
		io.WriteString(w, err.Error())
		return
	}

	id, err := uuid.NewV4()
	if err != nil {
		h.log.Err("Could not generate uuid: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		io.WriteString(w, err.Error())
		return
	}

	ch := make(chan string, 1)
	h.lock.Lock()
	h.pending[id.String()] = ch
	h.lock.Unlock()

	j := &job.Image{
		Common: job.Common{
			Src:       r.Form.Get("src"),
			Dst:       r.Form.Get("dst"),
			NotifyURL: h.httpURL + "/notify/" + id.String(),
		},
		Width:  int(width),
		Height: int(height),
	}

	d, err := json.Marshal(j)
	if err != nil {
		h.log.Warn("Bad http request: %s", err)
		w.WriteHeader(http.StatusBadRequest)
		io.WriteString(w, err.Error())
		return
	}
	h.log.Debug("Enqueuing job: %s", d)
	h.queue.Enqueue(string(d))

	url := <-ch
	h.log.Debug("Received url callback: %s", url)
	w.Write([]byte(url))
}

func (h *Handler) serveNotification(id string, w http.ResponseWriter, r *http.Request) {
	h.lock.Lock()
	ch := h.pending[id]
	h.lock.Unlock()

	h.log.Debug("serving notifcation for: %s", id)
	ch <- r.Form.Get("url")
}
