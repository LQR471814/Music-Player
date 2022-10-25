package main

import (
	"log"
	"music-player/server/api"
	"music-player/server/env"
	"net"
	"net/http"

	"github.com/improbable-eng/grpc-web/go/grpcweb"
	"google.golang.org/grpc"
)

func main() {
	listener, err := net.Listen("tcp", ":8000")
	if err != nil {
		log.Fatal(err)
	}

	grpcServer := grpc.NewServer()
	indexServer := &IndexServer{
		Store: NewIndex(),
	}
	api.RegisterIndexServer(grpcServer, indexServer)

	wrappedServer := grpcweb.WrapServer(
		grpcServer,
		grpcweb.WithOriginFunc(func(origin string) bool {
			return true
		}),
	)

	log.Printf("Serving on %s...", env.Options.Address)
	http.Serve(listener, SplitGRPCTraffic(
		http.FileServer(http.Dir(env.Options.AudioDirectory)),
		wrappedServer,
	))
	// albums, err := PullAlbums()
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// for _, a := range albums {
	// 	log.Printf("%s - %s ===================\n", a.AlbumArtist, a.Title)
	// 	for _, t := range a.Tracks {
	// 		log.Printf("\t%d. %s - %s\n", t.Disc, t.Artist, t.Title)
	// 	}
	// }
}
