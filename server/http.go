package main

import (
	"net/http"
	"strings"

	"github.com/LQR471814/music-player/server/logging"

	"github.com/gorilla/websocket"
)

func CORS(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodOptions && r.Header.Get("Access-Control-Request-Method") != "" {
			logging.Info.Println("got options request")
		}
		handler.ServeHTTP(w, r)
	})
}

func SplitGRPCTraffic(fallback http.Handler, grpcHandler http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			if strings.Contains(r.Header.Get("Content-Type"), "application/grpc") ||
				websocket.IsWebSocketUpgrade(r) {
				grpcHandler.ServeHTTP(w, r)
				return
			}
			fallback.ServeHTTP(w, r)
		},
	)
}

func HandleSubdirectory(mux *http.ServeMux, dir string) {
	mux.Handle(dir+"/", http.StripPrefix(dir, http.FileServer(http.Dir("."+dir))))
}
