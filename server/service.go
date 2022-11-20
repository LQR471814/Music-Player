package main

import (
	"context"

	"github.com/LQR471814/music-player/server/api"
	"github.com/LQR471814/music-player/server/logging"
)

type APIServer struct {
	api.UnimplementedAPIServer
	AlbumStore *AlbumIndex
}

func (s *APIServer) Index(_ *api.Empty, server api.API_IndexServer) error {
	updates := []*api.Update{}
	i := 0
	for _, a := range s.AlbumStore.Index.Values {
		updates = append(updates, &api.Update{
			Action: api.Action_ADD,
			Payload: &api.Update_Album{
				Album: a,
			},
		})
		i++
		if i == 10 {
			err := server.Send(&api.BatchedUpdate{
				Updates: updates,
			})
			if err != nil {
				logging.Warn.Println(err)
				return err
			}
			updates = []*api.Update{}
			i = 0
		}
	}
	if i > 0 {
		err := server.Send(&api.BatchedUpdate{
			Updates: updates,
			Status:  &api.Status{Ok: true},
		})
		if err != nil {
			logging.Warn.Println(err)
			return err
		}
	}

	channel := make(chan *api.BatchedUpdate)
	s.AlbumStore.Channels = append(s.AlbumStore.Channels, channel)
	for {
		value := <-channel
		err := server.Send(value)
		if err != nil {
			break
		}
	}
	return nil
}

func (s *APIServer) Modify(ctx context.Context, update *api.Update) (*api.Status, error) {
	s.AlbumStore.Update(update)
	return &api.Status{Ok: true}, nil
}

func (s *APIServer) AddFrom(ctx context.Context, source *api.Source) (*api.Status, error) {
	return &api.Status{Ok: true}, nil
}
