package photo

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"time"

	"github.com/rwcarlsen/goexif/exif"
)

type Photo struct {
	Path      string
	Latitude  float64
	Longitude float64
	Altitude  *float64
	Timestamp time.Time
}

var ErrNoGPS = errors.New("photo has no GPS coordinates")

func ScanDirectory(root string) ([]Photo, error) {
	var photos []Photo

	err := filepath.WalkDir(root, func(path string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if entry.IsDir() {
			return nil
		}

		if !isJPEG(path) {
			return nil
		}

		p, err := readPhoto(path)
		if err != nil {
			if errors.Is(err, ErrNoGPS) {
				fmt.Printf("SKIP  %s: no GPS data\n", path)
				return nil
			}

			fmt.Printf("SKIP  %s: %v\n", path, err)
			return nil
		}

		photos = append(photos, *p)

		fmt.Printf(
			"FOUND %s: %.6f, %.6f\n",
			path,
			p.Latitude,
			p.Longitude,
		)

		return nil
	})

	if err != nil {
		return nil, err
	}

	sort.Slice(photos, func(i, j int) bool {
		return photos[i].Timestamp.Before(photos[j].Timestamp)
	})

	return photos, nil
}

func readPhoto(path string) (*Photo, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	x, err := exif.Decode(file)
	if err != nil {
		return nil, fmt.Errorf("read EXIF: %w", err)
	}

	lat, lon, err := x.LatLong()
	if err != nil {
		return nil, ErrNoGPS
	}

	timestamp, err := x.DateTime()
	if err != nil {
		// If EXIF doesn't contain a timestamp, use the file modification time.
		info, statErr := file.Stat()
		if statErr != nil {
			return nil, fmt.Errorf("read timestamp: %w", err)
		}

		timestamp = info.ModTime()
	}

	altitude := readAltitude(x)

	relativePath := path

	return &Photo{
		Path:      relativePath,
		Latitude:  lat,
		Longitude: lon,
		Altitude:  altitude,
		Timestamp: timestamp,
	}, nil
}

func readAltitude(x *exif.Exif) *float64 {
	tag, err := x.Get(exif.GPSAltitude)
	if err != nil {
		return nil
	}

	numerator, denominator, err := tag.Rat2(0)
	if err != nil || denominator == 0 {
		return nil
	}

	altitude := float64(numerator) / float64(denominator)

	ref, err := x.Get(exif.GPSAltitudeRef)
	if err == nil {
		value, err := ref.Int(0)
		if err == nil && value == 1 {
			altitude = -altitude
		}
	}

	return &altitude
}

func Filename(path string) string {
	return filepath.Base(path)
}
