package main

import (
	"context"
	"music-player/server/api"
)

type IndexServer struct {
	api.UnimplementedIndexServer
	Store *Index
}

func (s *IndexServer) Index(_ *api.Empty, server api.Index_IndexServer) error {
	updates := []*api.Update{}
	i := 0
	for _, a := range s.Store.Albums {
		updates = append(updates, &api.Update{
			Action: api.Action_ADD,
			Payload: &api.Update_Album{
				Album: a,
			},
		})
		i++
		if i == 10 {
			server.Send(&api.BatchedUpdate{
				Updates: updates,
			})
			updates = []*api.Update{}
			i = 0
		}
	}

	channel := make(chan *api.BatchedUpdate)
	s.Store.Channels = append(s.Store.Channels, channel)
	for {
		value := <-channel
		err := server.Send(value)
		if err != nil {
			break
		}
	}
	return nil
}

func (s *IndexServer) Modify(ctx context.Context, update *api.Update) (*api.Status, error) {
	s.Store.Update(update)
	return &api.Status{Ok: true}, nil
}

func (s *IndexServer) AddFrom(ctx context.Context, source *api.Source) (*api.Status, error) {
	return &api.Status{Ok: true}, nil
}
