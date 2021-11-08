package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/wavesoftware/k8s-aware/pkg/k8s"
	"github.com/wavesoftware/k8s-aware/pkg/utils/retcode"
)

// Serve an HTTP service.
func Serve() int {
	client, err := k8s.NewClient()
	if err != nil {
		return handleError(err)
	}
	err = http.ListenAndServe(bind(), handler{client})
	if err != nil {
		return handleError(err)
	}
	return 0
}

func handleError(err error) int {
	log.Println(err)
	return retcode.Calc(err)
}

func bind() string {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	return ":" + port
}

type handler struct {
	k8s.Info
}

func (h handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	pods, err := h.Pods()
	if err != nil {
		handleHTTPError(err, w)
		return
	}
	bytes, err := json.Marshal(pods)
	if err != nil {
		handleHTTPError(err, w)
		return
	}
	w.Header().Set("content-type", "application/json")
	_, _ = w.Write(bytes)
}

func handleHTTPError(err error, w http.ResponseWriter) {
	w.WriteHeader(http.StatusInternalServerError)
	_, _ = fmt.Fprintf(w, "Error: %v", err)
}
