package index

import (
	"bytes"
	"encoding/gob"
	"os"
	"sync"
	"time"

	"github.com/LQR471814/music-player/server/env"
	"github.com/LQR471814/music-player/server/logging"
	"github.com/LQR471814/music-player/server/utils"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
)

func Sleep(d time.Duration) chan bool {
	channel := make(chan bool)
	go func() {
		time.Sleep(d)
		channel <- true
	}()
	return channel
}

type Index[T any] struct {
	Values   map[string]T
	Location string

	Channels   []chan T
	UpdateLock sync.Mutex

	Marshal   func(T) ([]byte, error)
	Unmarshal func([]byte) (T, error)
}

func NewIndex[T any](
	location string,
	fetch func() (map[string]T, error),
	marshal func(T) ([]byte, error),
	unmarshal func([]byte) (T, error),
) *Index[T] {
	index := &Index[T]{
		Values:     make(map[string]T),
		Location:   location,
		Channels:   make([]chan T, 0),
		UpdateLock: sync.Mutex{},
		Marshal:    marshal,
		Unmarshal:  unmarshal,
	}
	_, err := os.Stat(location)
	if err == nil && !env.Options.Reset {
		index.Load()
	} else {
		values, err := fetch()
		if err != nil {
			logging.Error.Fatal("failed to execute fetch():", err)
		}
		index.Values = values
		index.Store()
	}
	return index
}

func NewProtoIndex[T protoreflect.ProtoMessage](
	location string,
	fetch func() (map[string]T, error),
) *Index[T] {
	return NewIndex(
		location, fetch,
		func(p T) ([]byte, error) {
			return proto.Marshal(p)
		},
		func(b []byte) (T, error) {
			value := utils.InitializePointerType[T]()
			return value, proto.Unmarshal(b, value)
		},
	)
}

func NewPrimitiveIndex[T any](
	location string,
	fetch func() (map[string]T, error),
) *Index[T] {
	return NewIndex(
		location, fetch,
		func(p T) ([]byte, error) {
			buffer := bytes.NewBuffer(nil)
			encoder := gob.NewEncoder(buffer)
			err := encoder.Encode(p)
			return buffer.Bytes(), err
		},
		func(b []byte) (T, error) {
			var value T
			buffer := bytes.NewBuffer(b)
			decoder := gob.NewDecoder(buffer)
			err := decoder.Decode(&value)
			return value, err
		},
	)
}

func (i *Index[T]) Load() error {
	f, err := os.Open(i.Location)
	if err != nil {
		logging.Error.Println("could not open index:", err)
		return err
	}

	decoded := map[string][]byte{}

	decoder := gob.NewDecoder(f)
	err = decoder.Decode(&decoded)
	if err != nil {
		logging.Error.Println("could not serialize index:", err)
		return err
	}

	i.UpdateLock.Lock()
	for id, serialized := range decoded {
		value, err := i.Unmarshal(serialized)
		if err != nil {
			logging.Error.Println("error while deserializing proto:", id, err)
		}
		i.Values[id] = value
	}
	i.UpdateLock.Unlock()
	return nil
}

func (i *Index[T]) Store() error {
	i.UpdateLock.Lock()
	serialize := map[string][]byte{}
	for k, a := range i.Values {
		bytes, err := i.Marshal(a)
		if err != nil {
			logging.Error.Println("error while serializing proto:", k, err)
			continue
		}
		serialize[k] = bytes
	}
	i.UpdateLock.Unlock()

	f, err := os.Create(i.Location)
	if err != nil {
		logging.Error.Println("could not create index:", err)
		return err
	}

	encoder := gob.NewEncoder(f)
	err = encoder.Encode(serialize)
	if err != nil {
		logging.Error.Println("could not serialize index:", err)
	}
	return err
}
