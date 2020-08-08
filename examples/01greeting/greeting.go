package greeting

import (
	"fmt"
	"net/http"
)

func Greeting(hour int) string {
	if hour < 12 {
		return "Good morning"
	} else if 12 < hour && hour < 18 {
		return "Good afternoon"
	} else {
		return "Good evening"
	}
}

type App struct {
	Clock Clock
}

func (app *App) Greeting(w http.ResponseWriter, r *http.Request) {
	now := app.Clock.Now()
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, `{"message": %q, "time": %q}`, Greeting(now.Hour()), now)
}

func (app *App) Mount(mux *http.ServeMux) {
	mux.HandleFunc("/greeting", app.Greeting)
}
