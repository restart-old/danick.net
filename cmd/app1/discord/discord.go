package discord

import (
	"fmt"
	"net/http"

	discordDM "code.danick.net/lib/discord/dm"
)

var Sessions = discordDM.NewRunning()

func StopFunc(w http.ResponseWriter, r *http.Request) {
	token := r.URL.Query().Get("token")
	if Sessions.TokenRunning(token) {
		Sessions.StopCleaning(token)
		fmt.Fprintf(w, "stopped session with this token")
		return
	}
	http.Error(w, "no session is started with this token", 202)
}
func StatusFunc(w http.ResponseWriter, r *http.Request) {
	token := r.URL.Query().Get("token")
	if Sessions.TokenRunning(token) {
		w.Write(Sessions.StatusCleaning(token))
		return
	}
	http.Error(w, "no session is started with this token", 202)
}
func StartFunc(w http.ResponseWriter, r *http.Request) {
	token := r.URL.Query().Get("token")
	if !Sessions.TokenRunning(token) {
		err := Sessions.StartCleaning(token)
		if err != nil {
			http.Error(w, "could not connect to discord.com using this token", 202)
			return
		}
		fmt.Fprintf(w, "session started with this token")
		return
	}
	http.Error(w, "a session is already started with this token", 201)
}
