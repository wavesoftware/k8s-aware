package server

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/wavesoftware/k8s-aware/pkg/utils/retcode"
)

// Serve an HTTP service.
func Serve() int {
	err := http.ListenAndServe(bind(), handler{})
	if err != nil {
		log.Println(err)
		return retcode.Calc(err)
	}
	return 0
}

func bind() string {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	return ":" + port
}

type handler struct{}

func (h handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hi there, I love %s!", r.URL.Path[1:])
}
