package main

import (
	"flag"
	"log"
	"regexp"
	"time"
)

// list images older than age matching match
// remove images older than age that aren't running

var (
	DryRun       = flag.Bool("dry-run", false, "don't run delete functions")
	ImageAge     = flag.String("age", "48h", "delete images older than age")
	ImagePattern = flag.String("pattern", "^.*$", "match image names")
	DockerHost   = flag.String("docker", "tcp://127.0.0.1:2375", "docker endpoint")
)

func main() {
	flag.Usage = usage
	flag.Parse()

	if *DryRun {
		log.Println("[INFO] dry run enabled")
	}

	drydock, err := NewDryDock(*DockerHost)
	if err != nil {
		log.Fatalf("[FATAL] drydock: %s\n", err)
	}

	drydock.Age, err = time.ParseDuration(*ImageAge)
	if err != nil {
		log.Fatalf("[FATAL] age: %s\n", err)
	}

	drydock.Pattern, err = regexp.Compile(*ImagePattern)
	if err != nil {
		log.Fatalf("[FATAL] pattern: %s\n", err)
	}

	images, err := drydock.ListImages()
	if err != nil {
		log.Fatalf("[FATAL] images: %s\n", err)
	}

	imagesInUse, err := drydock.ListInUseImages()
	if err != nil {
		log.Fatalf("[FATAL] images in use: %s\n", err)
	}

	log.Printf("[INFO] %d images scheduled for deletion\n", len(images))
	for _, image := range images {
		if imagesInUse.Exist(image) {
			log.Printf("[INFO] skipping %s, in use\n", image)
			continue
		}

		if !*DryRun {
			err := drydock.RemoveImage(image)
			if err != nil {
				log.Fatalf("[FATAL] remove image: %s\n", err)
			}
		}

		log.Printf("[INFO] deleted image %s\n", image)
	}
}
