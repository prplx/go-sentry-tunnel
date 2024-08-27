package handlers

import (
	"fmt"
	"net/http"
)

func HandleTunnel(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, tunnel!")
}
