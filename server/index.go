package main

import (
	"encoding/gob"
	"io/ioutil"
	"mime"
	"music-player/server/api"
	"music-player/server/env"
	"music-player/server/logging"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/dhowden/tag"
	"github.com/google/uuid"
	"google.golang.org/protobuf/proto"
)

var extToFormat = map[string]api.Format{
	"flac": api.Format_FLAC,
	"ogg":  api.Format_OGG,
	"m4a":  api.Format_M4A,
	"mp3":  api.Format_MP3,
	"wav":  api.Format_WAV,
}

var discRegex *regexp.Regexp
var bandcampRegex *regexp.Regexp

func init() {
	bandcampRe, err := regexp.Compile(`(.+) - (.+) - ([0-9]+) (.+)`)
	if err != nil {
		logging.Error.Fatal(err)
		return
	}
	bandcampRegex = bandcampRe

	discRe, err := regexp.Compile(`([0-9]*)(\.| |-|\/|\||:|\)|\]|\})*(.+)`)
	if err != nil {
		logging.Error.Fatal(err)
		return
	}
	discRegex = discRe
}

func RemoveExt(filename string) string {
	split := strings.Split(filename, ".")
	return strings.Join(split[:len(split)-1], ".")
}

func InferTrack(fileInfo os.FileInfo, t *api.Track) {
	name := RemoveExt(fileInfo.Name())
	if bandcampRegex.MatchString(name) {
		groups := bandcampRegex.FindStringSubmatch(name)

		t.Artist = groups[1]
		t.Title = groups[4]

		disc, err := strconv.Atoi(groups[3])
		if err != nil {
			logging.Warn.Println(err)
		} else {
			t.Disc = int32(disc)
		}
		return
	}

	groups := discRegex.FindStringSubmatch(name)

	if groups[1] != "" {
		disc, err := strconv.Atoi(groups[1])
		if err != nil {
			logging.Warn.Println(err)
		} else {
			t.Disc = int32(disc)
		}
	}

	if strings.Contains(groups[3], "-") {
		segments := strings.Split(groups[3], "-")
		t.Title = strings.Trim(strings.Join(segments[1:], "-"), " ")
		t.Artist = strings.Trim(segments[0], " ")
	} else {
		t.Title = groups[3]
	}
	t.Format = extToFormat[filepath.Ext(fileInfo.Name())]
}

func InferAlbum(fileInfo os.FileInfo, a *api.Album) {
	if strings.Contains(fileInfo.Name(), "-") {
		segments := strings.Split(fileInfo.Name(), "-")
		a.Title = strings.Trim(strings.Join(segments[1:], "-"), " ")
		a.AlbumArtist = strings.Trim(segments[0], " ")
		return
	}
	a.Title = fileInfo.Name()
}

func InferCover(path string, a *api.Album) {
	entries, err := ioutil.ReadDir(path)
	if err != nil {
		logging.Warn.Println(err)
		return
	}
	for _, e := range entries {
		if RemoveExt(e.Name()) == "cover" {
			f, err := os.Open(filepath.Join(path, e.Name()))
			if err != nil {
				logging.Warn.Println(err)
				return
			}

			data := make([]byte, 0)
			_, err = f.Read(data)
			if err != nil {
				logging.Warn.Println(err)
				return
			}

			a.Cover = &api.Picture{
				Mime: mime.TypeByExtension(filepath.Ext(e.Name())),
				Data: data,
			}
			return
		}
	}
}

func HandleTrack(album *api.Album, path string) {
	err := syscall.Access(path, syscall.O_RDWR)
	if err != nil {
		logging.Warn.Println(path, err)
		return
	}

	ext := filepath.Ext(path)
	valid := false
	for k := range extToFormat {
		if "."+k == ext {
			valid = true
		}
	}
	if !valid {
		return
	}

	rawFile, err := os.Open(path)
	if err != nil {
		logging.Warn.Println(path, err)
		return
	}
	fileInfo, err := rawFile.Stat()
	if err != nil {
		logging.Warn.Println(path, err)
		return
	}

	track := &api.Track{
		Id:   uuid.NewString(),
		Path: path,
	}
	InferTrack(fileInfo, track)

	m, err := tag.ReadFrom(rawFile)
	if err != nil && err != tag.ErrNoTagsFound {
		logging.Warn.Println(path, err)
		return
	}

	if m != nil {
		album.Title = Fallback([]string{m.Album(), album.Title})
		album.AlbumArtist = Fallback([]string{m.AlbumArtist(), album.AlbumArtist})
		if m.Picture() != nil && len(m.Picture().Data) > 0 {
			album.Cover = &api.Picture{
				Data:        m.Picture().Data,
				Mime:        m.Picture().MIMEType,
				Description: m.Picture().Description,
			}
		}
		disc, _ := m.Disc()

		track.Artist = Fallback([]string{m.Artist(), track.Artist})
		track.Composer = Fallback([]string{m.Composer(), track.Composer})
		track.Genre = Fallback([]string{m.Genre(), track.Genre})
		track.Year = Fallback([]int32{int32(m.Year()), track.Year})
		track.Disc = Fallback([]int32{int32(disc), track.Disc})
	}
	album.Tracks[track.Id] = track
}

func PullAlbums() ([]*api.Album, error) {
	albums := []*api.Album{}

	dirs, err := ioutil.ReadDir(env.Options.AudioDirectory)
	if err != nil {
		return nil, err
	}
	for _, fileInfo := range dirs {
		path := filepath.Join(env.Options.AudioDirectory, fileInfo.Name())

		album := &api.Album{
			Id:     uuid.NewString(),
			Tracks: make(map[string]*api.Track),
		}
		InferAlbum(fileInfo, album)

		if fileInfo.IsDir() {
			InferCover(path, album)
			subDirs, err := ioutil.ReadDir(path)
			if err != nil {
				logging.Warn.Println(err)
				continue
			}
			for _, trackFile := range subDirs {
				if !trackFile.IsDir() {
					HandleTrack(album, filepath.Join(path, trackFile.Name()))
				}
			}
		} else {
			HandleTrack(album, path)
		}
		albums = append(albums, album)
	}

	return albums, nil
}

func IndexLocation() string {
	return filepath.Join(env.Options.AudioDirectory, env.Options.IndexName)
}

type Index struct {
	Albums   map[string]*api.Album
	Channels []chan *api.BatchedUpdate

	modifications []*api.Update
	cancelUpdate  chan bool
	updateLock    sync.Mutex
}

func NewIndex() *Index {
	index := &Index{
		Albums:        make(map[string]*api.Album),
		Channels:      make([]chan *api.BatchedUpdate, 0),
		modifications: make([]*api.Update, 0),
	}

	_, err := os.Stat(IndexLocation())
	if err == nil {
		index.Load()
	} else {
		albums, err := PullAlbums()
		if err != nil {
			logging.Error.Fatal("failed to read albums:", err)
		}
		for _, a := range albums {
			index.Albums[a.Id] = a
		}
		index.Store()
	}

	return index
}

func (i *Index) Update(update *api.Update) {
	i.updateLock.Lock()

	i.modifications = append(i.modifications, update)
	switch update.Payload.(type) {
	case *api.Update_Album:
		album := update.Payload.(*api.Update_Album).Album
		switch update.Action {
		case api.Action_ADD:
			album.Id = uuid.NewString()
			i.Albums[album.Id] = album
		case api.Action_REMOVE:
			delete(i.Albums, album.Id)
		case api.Action_OVERRIDE:
			i.Albums[album.Id] = album
		}
	case *api.Update_Track:
		updateTrack := update.Payload.(*api.Update_Track).Track
		track := updateTrack.Track
		albumId := updateTrack.AlbumId
		switch update.Action {
		case api.Action_ADD:
			track.Id = uuid.NewString()
			i.Albums[albumId].Tracks[track.Id] = track
		case api.Action_REMOVE:
			delete(i.Albums[albumId].Tracks, track.Id)
		case api.Action_OVERRIDE:
			i.Albums[albumId].Tracks[track.Id] = track
		}
	}

	if i.cancelUpdate != nil {
		i.cancelUpdate <- true
	}
	i.cancelUpdate = make(chan bool)
	go func() {
		for {
			select {
			case <-i.cancelUpdate:
				return
			case <-Sleep(time.Second * 10):
				for _, c := range i.Channels {
					c <- &api.BatchedUpdate{Updates: i.modifications}
				}
				i.Store()
				return
			}
		}
	}()

	i.updateLock.Unlock()
}

func (i *Index) Load() {
	f, err := os.Open(IndexLocation())
	if err != nil {
		logging.Error.Println("could not open index:", err)
		return
	}

	decoded := map[string][]byte{}

	decoder := gob.NewDecoder(f)
	err = decoder.Decode(&decoded)
	if err != nil {
		logging.Error.Println("could not serialize index:", err)
	}

	for id, serialized := range decoded {
		album := &api.Album{}
		err := proto.Unmarshal(serialized, album)
		if err != nil {
			logging.Error.Println("error while deserializing album:", id, err)
		}
		i.Albums[id] = album
	}
}

func (i *Index) Store() {
	i.updateLock.Lock()
	serialize := map[string][]byte{}
	for _, a := range i.Albums {
		bytes, err := proto.Marshal(a)
		if err != nil {
			logging.Error.Println("error while serializing album:", a.Id, err)
			continue
		}
		serialize[a.Id] = bytes
	}
	i.updateLock.Unlock()

	f, err := os.Create(IndexLocation())
	if err != nil {
		logging.Error.Println("could not create index:", err)
		return
	}

	encoder := gob.NewEncoder(f)
	err = encoder.Encode(serialize)
	if err != nil {
		logging.Error.Println("could not serialize index:", err)
	}
}
