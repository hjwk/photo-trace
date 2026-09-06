package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/hjwk/photo-trace/internal/gpx"
	"github.com/hjwk/photo-trace/internal/photo"
	"github.com/hjwk/photo-trace/internal/web"
)

func main() {
	input := flag.String(
		"input",
		"",
		"directory containing photos",
	)

	output := flag.String(
		"output",
		"./output",
		"output directory",
	)

	flag.Parse()

	if *input == "" {
		fmt.Println("Usage:")
		fmt.Println("  photo-trace -input ./photos -output ./output")
		os.Exit(1)
	}

	if err := run(*input, *output); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func run(input, output string) error {
	fmt.Printf("Scanning %s...\n\n", input)

	photos, err := photo.ScanDirectory(input)
	if err != nil {
		return fmt.Errorf("scan photos: %w", err)
	}

	if len(photos) == 0 {
		return fmt.Errorf("no geotagged JPEG photos found")
	}

	fmt.Printf("\nFound %d geotagged photos\n", len(photos))

	if err := os.MkdirAll(output, 0755); err != nil {
		return fmt.Errorf("create output directory: %w", err)
	}

	gpxPath := filepath.Join(output, "route.gpx")

	fmt.Printf("Writing %s...\n", gpxPath)

	if err := gpx.Write(gpxPath, photos); err != nil {
		return fmt.Errorf("write GPX: %w", err)
	}

	htmlPath := filepath.Join(output, "index.html")

	templatePath := filepath.Join(
		"templates",
		"index.html",
	)

	fmt.Printf("Writing %s...\n", htmlPath)

	if err := web.Write(
		htmlPath,
		photos,
		templatePath,
	); err != nil {
		return fmt.Errorf("write HTML: %w", err)
	}

	fmt.Println("Copying photos...")

	if err := web.CopyPhotos(
		output,
		photos,
	); err != nil {
		return fmt.Errorf("copy photos: %w", err)
	}

	fmt.Println()
	fmt.Println("Done!")
	fmt.Printf("  Photos: %d\n", len(photos))
	fmt.Printf("  GPX:    %s\n", gpxPath)
	fmt.Printf("  HTML:   %s\n", htmlPath)

	return nil
}
