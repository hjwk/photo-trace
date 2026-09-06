package gpx

import (
	"encoding/xml"
	"fmt"
	"os"
	"time"

	"github.com/hjwk/photo-trace/internal/photo"
)

type GPX struct {
	XMLName xml.Name `xml:"gpx"`

	Version string `xml:"version,attr"`
	Creator string `xml:"creator,attr"`
	XMLNS   string `xml:"xmlns,attr"`

	Tracks []Track `xml:"trk"`
}
type Track struct {
	Name     string    `xml:"name"`
	Segments []Segment `xml:"trkseg"`
}

type Segment struct {
	Points []TrackPoint `xml:"trkpt"`
}

type TrackPoint struct {
	Latitude  float64    `xml:"lat,attr"`
	Longitude float64    `xml:"lon,attr"`
	Elevation *float64   `xml:"ele,omitempty"`
	Time      *time.Time `xml:"time,omitempty"`
}

func Write(path string, photos []photo.Photo) error {
	if len(photos) == 0 {
		return fmt.Errorf("cannot generate GPX with no photos")
	}

	points := make([]TrackPoint, 0, len(photos))

	for _, p := range photos {
		points = append(points, TrackPoint{
			Latitude:  p.Latitude,
			Longitude: p.Longitude,
			Elevation: p.Altitude,
			Time:      new(p.Timestamp),
		})
	}

	document := GPX{
		Version: "1.1",
		Creator: "photo-trace",
		XMLNS:   "http://www.topografix.com/GPX/1/1",
		Tracks: []Track{
			{
				Name: "Photo track",
				Segments: []Segment{
					{
						Points: points,
					},
				},
			},
		},
	}

	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := xml.NewEncoder(file)
	encoder.Indent("", "  ")

	if _, err := file.WriteString(xml.Header); err != nil {
		return err
	}

	if err := encoder.Encode(document); err != nil {
		return err
	}

	return encoder.Flush()
}
