package http

import (
	"encoding/json"
	"github.com/felixge/wandler/job"
	"github.com/felixge/wandler/log"
	"github.com/felixge/wandler/queue"
	"io"
	"net/http"
	"strconv"
)

const multipartMaxMemory = 64 * 1024

func NewHandler(c HandlerConfig) (*Handler, error) {
	return &Handler{
		log:   c.Log,
		queue: c.Queue,
	}, nil
}

type HandlerConfig struct {
	Log   log.Interface
	Queue queue.Interface
}

type Handler struct {
	log   log.Interface
	queue queue.Interface
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

	j := &job.Image{
		Common: job.Common{
			Src: r.Form.Get("src"),
			Dst: r.Form.Get("dst"),
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
}
