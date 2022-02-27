package middlewares

import (
	"fmt"
	"log"
	"net/http"
)

func LogResponse(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println(fmt.Sprintf("request from %v - STATUS %v", r.RemoteAddr, 200))
		next.ServeHTTP(w, r)
	})
}
