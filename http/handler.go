package http

import (
	"github.com/felixge/wandler/job"
	"github.com/felixge/wandler/queue"
	"github.com/felixge/wandler/log"
	"io"
	"net/http"
)

const multipartMaxMemory = 64 * 1024

func NewHandler(c HandlerConfig) (*Handler, error) {
	return &Handler{
		log: c.Log,
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

	j := &job.ImageJob{
		Common: job.Common{
			Src: r.Form.Get("src"),
			Dst: r.Form.Get("dst"),
		},
		Width:  r.Form.Get("width"),
		Height: r.Form.Get("height"),
	}

	h.log.Debug("Enqueuing job: %+v", j)
	h.queue.Enqueue(j)
}
