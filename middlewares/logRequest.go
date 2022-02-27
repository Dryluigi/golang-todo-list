package middlewares

import (
	"fmt"
	"log"
	"net/http"
)

func LogRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println(fmt.Sprintf("%v -> %v %v%v", r.RemoteAddr, r.Method, r.Host, r.URL.Path))

		next.ServeHTTP(w, r)
	})
}
