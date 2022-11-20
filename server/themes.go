//* this file shall remain unused until further notice

package main

import (
	"bytes"
	"image"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"

	_ "image/gif"
	"image/jpeg"
	_ "image/png"

	"github.com/LQR471814/music-player/server/api"
	"github.com/LQR471814/music-player/server/env"
	"github.com/LQR471814/music-player/server/index"
	"github.com/LQR471814/music-player/server/logging"
	"github.com/LQR471814/music-player/server/utils"
	"github.com/disintegration/imaging"
	"github.com/mccutchen/palettor"
)

var backgroundDir = "env.Config.BackgroundDirectory"

type ThemeEntry struct {
	Path string
}

type ThemeIndex struct {
	Index *index.Index[ThemeEntry]
}

func PrepareWallpaper(name string, buffer io.Reader) (ThemeEntry, error) {
	path := filepath.Join(backgroundDir, name)
	lookup := path + ".lookup"

	_, err := os.Stat(path)
	if err == os.ErrNotExist && !env.Options.Reset {
		f, err := os.Create(path)
		if err != nil {
			return ThemeEntry{}, err
		}
		defer f.Close()

		_, err = io.Copy(f, buffer)
		if err != nil {
			return ThemeEntry{}, err
		}
	}

	_, err = os.Stat(lookup)
	if err == os.ErrNotExist && !env.Options.Reset {
		f, err := os.Create(lookup)
		if err != nil {
			return ThemeEntry{}, err
		}
		defer f.Close()

		img, _, err := image.Decode(buffer)
		if err != nil {
			return ThemeEntry{}, err
		}

		bounds := img.Bounds()
		w := bounds.Dx()
		h := bounds.Dy()

		aspectRatio := float64(h) / float64(w)

		resizeWidth := (2 * float64(env.Options.PaletteResolution)) / (aspectRatio + 1)
		width := int(resizeWidth)
		height := int(resizeWidth * aspectRatio)

		resized := imaging.Resize(img, width, height, imaging.Box)

		err = jpeg.Encode(f, resized, &jpeg.Options{
			Quality: 90,
		})
		if err != nil {
			return ThemeEntry{}, err
		}
	}

	return ThemeEntry{Path: path}, nil
}

// func ThemeFromBuffer(buffer io.Reader) (*api.Theme, error) {
// 	img, _, err := image.Decode(buffer)
// 	if err != nil {
// 		return nil, err
// 	}
// 	p, err := palettor.Extract(3, 1, img)
// 	if err != nil {
// 		return nil, err
// 	}
// 	entries := p.Entries()
// 	return &api.Theme{
// 		Colors: &api.Colors{
// 			Primary:    utils.ColorToString(entries[0].Color),
// 			Secondary:  utils.ColorToString(entries[1].Color),
// 			Background: utils.ColorToString(entries[2].Color),
// 		},
// 		Wallpaper: &api.Wallpaper{
// 			Type: api.WallpaperType_IMAGE,
// 		},
// 	}, nil
// }

func NewThemeIndex() *ThemeIndex {
	index := index.NewPrimitiveIndex(
		filepath.Join(backgroundDir, "index.pb"),
		func() (map[string]ThemeEntry, error) {
			backgrounds, err := ioutil.ReadDir(backgroundDir)
			if err != nil {
				return nil, err
			}
			themes := map[string]ThemeEntry{}
			for _, f := range backgrounds {
				if f.Name() == "index.pb" {
					continue
				}

				path := filepath.Join(backgroundDir, f.Name())
				file, err := os.Open(path)
				if err != nil {
					logging.Warn.Println(err)
					continue
				}
				defer file.Close()

				entry, err := PrepareWallpaper(f.Name(), file)
				if err != nil {
					logging.Warn.Println(err)
					continue
				}

				// theme, err := ThemeFromBuffer(file)
				// if err != nil {
				// 	logging.Warn.Println(err)
				// 	continue
				// }
				// theme.Wallpaper.Content = path
				themes[f.Name()] = entry
			}
			return themes, nil
		},
	)
	return &ThemeIndex{
		Index: index,
	}
}

func (s *ThemeIndex) Read(name string, bounds *api.Bounds) (*api.Theme, error) {
	lookup := s.Index.Values[name].Path + ".lookup"

	f, err := os.Open(lookup)
	if err != nil {
		return nil, err
	}

	img, _, err := image.Decode(f)
	if err != nil {
		return nil, err
	}

	cropped := imaging.Crop(img, image.Rectangle{
		Min: image.Point{
			X: int(bounds.X1),
			Y: int(bounds.Y1),
		},
		Max: image.Point{
			X: int(bounds.X2),
			Y: int(bounds.Y2),
		},
	})

	p, err := palettor.Extract(3, 1, cropped)
	if err != nil {
		logging.Error.Fatal(err)
	}
	entries := p.Entries()

	return &api.Theme{
		Wallpaper: &api.Wallpaper{
			Type:    api.WallpaperType_IMAGE,
			Content: s.Index.Values[name].Path,
		},
		Colors: &api.Colors{
			Primary:    utils.ColorToString(entries[0].Color),
			Secondary:  utils.ColorToString(entries[1].Color),
			Background: utils.ColorToString(entries[1].Color),
		},
		Effects: &api.Effects{
			BackgroundOpacity: 0.2,
			BorderOpacity:     0.3,
		},
	}, nil
}

func (s *ThemeIndex) Add(name string, imageData []byte) error {
	buff := bytes.NewBuffer(imageData)
	entry, err := PrepareWallpaper(name, buff)
	if err != nil {
		return err
	}
	s.Index.Values[name] = entry

	// theme, err := ThemeFromBuffer(buff)
	// if err != nil {
	// 	return err
	// }
	// s.Index.Values[name] = theme
	return s.Index.Store()
}

func (s *ThemeIndex) Remove(name string) error {
	delete(s.Index.Values, name)
	err := s.Index.Store()
	if err != nil {
		return err
	}
	err = os.Remove(filepath.Join(backgroundDir, name))
	if err != nil {
		return err
	}
	return os.Remove(filepath.Join(backgroundDir, name+".lookup"))
}
