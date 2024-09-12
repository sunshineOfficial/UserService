package handlers

import (
	"net/http"

	"github.com/go-chi/render"
)

func PingHandler(w http.ResponseWriter, r *http.Request) {
	render.PlainText(w, r, "pong")
}
