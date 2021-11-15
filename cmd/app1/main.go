package main

import (
	"net/http"

	"code.danick.net/cmd/app1/discord"
	"code.danick.net/lib/ratelimit"
	"github.com/gorilla/mux"
)

func main() {
	router := mux.NewRouter()
	router.Use(ratelimit.Middleware)
	router.Handle("/", http.RedirectHandler("https://danick.net/docs/", 308))
	router.Handle("/{[A-Z]|[a-z]\\w++}", http.RedirectHandler("https://danick.net/docs/", 308))

	router.HandleFunc("/api/v1/clean/start", discord.StartFunc)
	router.HandleFunc("/api/v1/clean/stop", discord.StopFunc)
	router.HandleFunc("/api/v1/clean/status", discord.StatusFunc)
	http.ListenAndServe(":80", router)
}
