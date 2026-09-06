package web

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/hjwk/photo-trace/internal/photo"
)

type PhotoData struct {
	Filename  string   `json:"filename"`
	Latitude  float64  `json:"latitude"`
	Longitude float64  `json:"longitude"`
	Altitude  *float64 `json:"altitude,omitempty"`
	Timestamp string   `json:"timestamp"`
}

type PageData struct {
	PhotosJSON template.JS
}

func Write(
	outputPath string,
	photos []photo.Photo,
	templatePath string,
) error {
	data := make([]PhotoData, 0, len(photos))

	for _, p := range photos {
		data = append(data, PhotoData{
			Filename:  photo.Filename(p.Path),
			Latitude:  p.Latitude,
			Longitude: p.Longitude,
			Altitude:  p.Altitude,
			Timestamp: p.Timestamp.Format(time.RFC3339),
		})
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("encode photo data: %w", err)
	}

	tmpl, err := template.ParseFiles(templatePath)
	if err != nil {
		return fmt.Errorf("parse template: %w", err)
	}

	file, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer file.Close()

	pageData := PageData{
		PhotosJSON: template.JS(jsonData),
	}

	return tmpl.Execute(file, pageData)
}

func CopyPhotos(
	outputRoot string,
	photos []photo.Photo,
) error {
	destination := filepath.Join(outputRoot, "photos")

	if err := os.MkdirAll(destination, 0755); err != nil {
		return err
	}

	for _, p := range photos {
		source := p.Path
		target := filepath.Join(destination, photo.Filename(p.Path))

		if err := copyFile(source, target); err != nil {
			return fmt.Errorf("copy %s: %w", source, err)
		}
	}

	return nil
}

func copyFile(source, destination string) error {
	input, err := os.Open(source)
	if err != nil {
		return err
	}
	defer input.Close()

	output, err := os.Create(destination)
	if err != nil {
		return err
	}
	defer output.Close()

	_, err = io.Copy(output, input)

	return err
}
