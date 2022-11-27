package main

import (
	"log"
	"mime"
	"net"
	"net/http"

	"github.com/LQR471814/music-player/server/api"
	"github.com/LQR471814/music-player/server/env"
	"github.com/LQR471814/music-player/server/logging"

	"github.com/improbable-eng/grpc-web/go/grpcweb"
	"github.com/rs/cors"
	"google.golang.org/grpc"
)

func main() {
	err := mime.AddExtensionType(".js", "application/javascript")
	if err != nil {
		logging.Error.Fatal(err)
	}

	listener, err := net.Listen("tcp", env.Options.Address)
	if err != nil {
		logging.Error.Fatal(err)
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
	// 	logging.Error.Fatal(err)
	// }
	// for i, a := range albums {
	// 	log.Printf("%d. %s - %s ===================\n", i, a.AlbumArtist, a.Title)
	// 	for _, t := range a.Tracks {
	// 		log.Printf("\t%d. %s - %s\n", t.Disc, t.Artist, t.Title)
	// 	}
	// }

	log.Printf("Serving on %s...", listener.Addr().String())

	fileMux := &http.ServeMux{}
	HandleSubdirectory(fileMux, "/audio")
	HandleSubdirectory(fileMux, "/covers")
	fileMux.Handle("/", http.FileServer(http.Dir("static")))

	http.ServeTLS(listener, cors.AllowAll().Handler(SplitGRPCTraffic(
		fileMux,
		wrappedServer,
	)), "host-crt.pem", "host-key.pem")
}
