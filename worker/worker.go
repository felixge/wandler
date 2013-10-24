package worker

import (
	"encoding/json"
	"fmt"
	"github.com/felixge/wandler/job"
	"github.com/felixge/wandler/log"
	"github.com/felixge/wandler/queue"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
)

func NewWorker(l log.Interface, q queue.Interface) (*Worker, error) {
	return &Worker{
		queue: q,
		log:   l,
	}, nil
}

type Worker struct {
	queue queue.Interface
	log   log.Interface
}

func (w *Worker) Run() error {
	for {
		msg, err := w.queue.Dequeue()
		if err == io.EOF {
			continue
		} else if err != nil {
			panic(err)
		}

		j := &job.Image{}
		if err := json.Unmarshal([]byte(msg), j); err != nil {
			w.log.Err("could not unmarshal: %s", err)
			continue
		}

		w.log.Debug("new job: %#v", j)
		w.execute(j)
	}
	return nil
}

func (w *Worker) execute(j *job.Image) error {
	w.log.Debug("Downloading: %s", j.Src)
	resp, err := http.Get(j.Src)
	if err != nil {
		return w.log.Err("could not request: %s: %s", j.Src, err)
	}
	defer resp.Body.Close()

	name := "download"+filepath.Ext(j.Src)
	file, err := os.OpenFile(name, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0666)
	if err != nil {
		return w.log.Err("could not open file: %s", err)
	}
	defer file.Close()

	if _, err := io.Copy(file, resp.Body); err != nil {
		return w.log.Err("could not download file: %s", err)
	}

	w.log.Debug("downloaded: %s: %s", j.Src, name)

	outputName := "result"+filepath.Ext(name)
	resize := fmt.Sprintf("%dx%d!", j.Width, j.Height)
	cmd := exec.Command("convert", name, "-resize", resize, outputName)
	if err := cmd.Run(); err != nil {
		return w.log.Err("could not execute cmd: %#v: %s", cmd, err)
	}
	w.log.Debug("result: %s", outputName)
	return nil
}
