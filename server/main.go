package main

import (
	"log"
	"net"
	"net/http"

	"github.com/LQR471814/music-player/server/api"
	"github.com/LQR471814/music-player/server/env"

	"github.com/improbable-eng/grpc-web/go/grpcweb"
	"github.com/rs/cors"
	"google.golang.org/grpc"
)

func main() {
	listener, err := net.Listen("tcp", env.Options.Address)
	if err != nil {
		log.Fatal(err)
	}

	grpcServer := grpc.NewServer()
	apiServer := &APIServer{
		AlbumStore: NewIndex(),
	}
	api.RegisterAPIServer(grpcServer, apiServer)

	wrappedServer := grpcweb.WrapServer(
		grpcServer,
		grpcweb.WithOriginFunc(func(origin string) bool {
			return true
		}),
	)

	// albums, err := PullAlbums()
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// for i, a := range albums {
	// 	log.Printf("%d. %s - %s ===================\n", i, a.AlbumArtist, a.Title)
	// 	for _, t := range a.Tracks {
	// 		log.Printf("\t%d. %s - %s\n", t.Disc, t.Artist, t.Title)
	// 	}
	// }

	log.Printf("Serving on %s...", listener.Addr().String())
	// http.Serve(listener, cors.AllowAll().Handler(SplitGRPCTraffic(
	// 	http.FileServer(http.Dir(env.Options.AudioDirectory)),
	// 	wrappedServer,
	// )))
	http.ServeTLS(listener, cors.AllowAll().Handler(SplitGRPCTraffic(
		http.FileServer(http.Dir(env.Options.AudioDirectory)),
		wrappedServer,
	)), "host-crt.pem", "host-key.pem")
}
