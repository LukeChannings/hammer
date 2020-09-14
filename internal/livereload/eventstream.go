package livereload

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"path"
	"strings"
	"time"

	"github.com/fsnotify/fsnotify"
)

type message struct {
	Path     string `json:"path"`
	Content  string `json:"content"`
	FileType string `json:"fileType"`
}

var LiveReloadPath string = "/__livereload__"

func Handler(w http.ResponseWriter, r *http.Request) {
	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming unsupported!", http.StatusInternalServerError)
		return
	}
	h := w.Header()
	h.Set("connection", "keep-alive")
	h.Set("cache-control", "no-cache")
	h.Set("content-type", "text/event-stream")
	h.Set("access-control-allow-origin", "*")

	done := make(chan bool)
	go func() {
		ticker := time.NewTicker(5 * time.Second)
		defer ticker.Stop()
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				if event.Op&fsnotify.Write == fsnotify.Write {
					fileType := strings.TrimLeft(path.Ext(event.Name), ".")

					data, err := ioutil.ReadFile(event.Name)

					if err != nil {
						fmt.Println(err)
						return
					}

					msg := message{
						Path:     event.Name,
						Content:  string(data),
						FileType: fileType,
					}

					msgWire, err := json.Marshal(msg)

					fmt.Printf("%s updated\n", event.Name)

					if err != nil {
						fmt.Println(err)
					} else {
						fmt.Fprintf(w, "data: %s\n\n", msgWire)
					}
					flusher.Flush()
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Println("error:", err)
			case <-ticker.C:
				fmt.Fprintf(w, "data: heartbeat\n\n")
				flusher.Flush()
			case <-r.Context().Done():
				ticker.Stop()
				done <- true
				return
			}
		}
	}()

	<-done
}
