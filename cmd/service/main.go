package main

import (
	"os"

	"github.com/wavesoftware/k8s-aware/pkg/server"
)

func main() {
	os.Exit(server.Serve())
}
