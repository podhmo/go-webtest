package httpbin

import (
	"net/http"
)

// Handler :
func Handler() http.Handler {
	mux := http.NewServeMux()
	BindHandlers(mux)
	return mux
}

// BindHandlers :
func BindHandlers(mux *http.ServeMux) {
	mux.HandleFunc("/status/", Status)
	mux.HandleFunc("/get", Get)
}

func Run(port string, mux http.Handler) error {
	return http.ListenAndServe(port, mux)
}
