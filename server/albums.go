package main

import (
	"bytes"
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io/ioutil"
	"mime"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/LQR471814/music-player/server/api"
	"github.com/LQR471814/music-player/server/env"
	"github.com/LQR471814/music-player/server/index"
	"github.com/LQR471814/music-player/server/logging"
	"github.com/LQR471814/music-player/server/utils"
	"github.com/disintegration/imaging"

	"github.com/dhowden/tag"
	"github.com/google/uuid"
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
	if len(split) == 1 {
		return split[0]
	}
	return strings.Join(split[:len(split)-1], ".")
}

func IconLocation(album *api.Album) string {
	return fmt.Sprintf("%s/%s.jpg", env.Options.IconDirectory, album.Id)
}

func SaveIcon(byteSlice []byte, album *api.Album) (string, error) {
	buffer := bytes.NewBuffer(byteSlice)

	img, _, err := image.Decode(buffer)
	if err != nil {
		return "", err
	}
	resized := imaging.Fill(img, 384, 384, imaging.Center, imaging.Linear)

	location := IconLocation(album)
	err = imaging.Save(resized, location)
	return location, err
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
	filename := RemoveExt(fileInfo.Name())
	if strings.Contains(filename, "-") {
		segments := strings.Split(filename, "-")
		a.Title = strings.Trim(strings.Join(segments[1:], "-"), " ")
		a.AlbumArtist = strings.Trim(segments[0], " ")
		return
	}
	a.Title = filename
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
			data, err := ioutil.ReadAll(f)
			if err != nil {
				logging.Warn.Println(err)
				return
			}

			mimetype := mime.TypeByExtension(filepath.Ext(e.Name()))
			location, err := SaveIcon(data, a)
			if err != nil {
				logging.Warn.Println(err)
				return
			}
			a.Cover = &api.Picture{
				Url:  location,
				Mime: mimetype,
			}
			return
		}
	}
}

func IsAudio(path string) bool {
	ext := filepath.Ext(path)
	valid := false
	for k := range extToFormat {
		if "."+k == ext {
			valid = true
		}
	}
	return valid
}

func HandleTrack(album *api.Album, path string) {
	err := syscall.Access(path, syscall.O_RDWR)
	if err != nil {
		logging.Warn.Println(path, err)
		return
	}

	if !IsAudio(path) {
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

	withoutHead := strings.Split(path, "/")[1:]
	track := &api.Track{
		Id:   uuid.NewString(),
		Path: filepath.Join(withoutHead...),
	}
	InferTrack(fileInfo, track)

	m, err := tag.ReadFrom(rawFile)
	if err != nil && err != tag.ErrNoTagsFound {
		logging.Warn.Println(path, err)
		return
	}

	if m != nil {
		album.Title = utils.Fallback([]string{m.Album(), album.Title})
		album.AlbumArtist = utils.Fallback([]string{m.AlbumArtist(), album.AlbumArtist})

		_, err := os.Stat(IconLocation(album))
		if m.Picture() != nil && len(m.Picture().Data) > 0 && os.IsNotExist(err) && album.Cover == nil {
			location, err := SaveIcon(m.Picture().Data, album)
			if err == nil {
				album.Cover = &api.Picture{
					Url:         location,
					Mime:        m.Picture().MIMEType,
					Description: m.Picture().Description,
				}
			} else {
				logging.Warn.Println(err)
			}
		}
		disc, _ := m.Disc()

		track.Artist = utils.Fallback([]string{m.Artist(), track.Artist})
		track.Composer = utils.Fallback([]string{m.Composer(), track.Composer})
		track.Genre = utils.Fallback([]string{m.Genre(), track.Genre})
		track.Year = utils.Fallback([]int32{int32(m.Year()), track.Year})
		track.Disc = utils.Fallback([]int32{int32(disc), track.Disc})
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

		if !fileInfo.IsDir() && !IsAudio(fileInfo.Name()) {
			continue
		}

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

type AlbumIndex struct {
	Index    *index.Index[*api.Album]
	Channels []chan *api.BatchedUpdate

	modifications []*api.Update
	cancelUpdate  chan bool
	updateLock    sync.Mutex
}

func NewIndex() *AlbumIndex {
	return &AlbumIndex{
		Index: index.NewProtoIndex(
			filepath.Join(env.Options.AudioDirectory, "index.pb"),
			func() (map[string]*api.Album, error) {
				albums := map[string]*api.Album{}
				pulled, err := PullAlbums()
				if err != nil {
					return nil, err
				}
				for _, a := range pulled {
					albums[a.Id] = a
				}
				return albums, nil
			},
		),
		Channels:      make([]chan *api.BatchedUpdate, 0),
		modifications: make([]*api.Update, 0),
	}
}

func (i *AlbumIndex) Update(update *api.Update) {
	i.updateLock.Lock()

	i.modifications = append(i.modifications, update)
	switch update.Payload.(type) {
	case *api.Update_Album:
		album := update.Payload.(*api.Update_Album).Album
		switch update.Action {
		case api.Action_ADD:
			album.Id = uuid.NewString()
			i.Index.Values[album.Id] = album
		case api.Action_REMOVE:
			delete(i.Index.Values, album.Id)
		case api.Action_OVERRIDE:
			i.Index.Values[album.Id] = album
		}
	case *api.Update_Track:
		updateTrack := update.Payload.(*api.Update_Track).Track
		track := updateTrack.Track
		albumId := updateTrack.AlbumId
		switch update.Action {
		case api.Action_ADD:
			track.Id = uuid.NewString()
			i.Index.Values[albumId].Tracks[track.Id] = track
		case api.Action_REMOVE:
			delete(i.Index.Values[albumId].Tracks, track.Id)
		case api.Action_OVERRIDE:
			i.Index.Values[albumId].Tracks[track.Id] = track
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
			case <-utils.Sleep(time.Second * 10):
				for _, c := range i.Channels {
					c <- &api.BatchedUpdate{
						Updates: i.modifications,
						Status:  &api.Status{Ok: true},
					}
				}
				i.Index.Store()
				return
			}
		}
	}()

	i.updateLock.Unlock()
}
