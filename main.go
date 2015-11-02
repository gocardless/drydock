package main

import (
	"flag"
	"log"
	"os"
	"regexp"
	"time"
)

// list images older than age matching match
// remove images older than age that aren't running

var (
	DryRun       = flag.Bool("dry-run", false, "don't run delete functions")
	ImageAge     = flag.String("age", "48h", "delete images older than age")
	ImagesToKeep = flag.Int("keep", 10, "keep at least this many images")
	ImagePattern = flag.String("pattern", "^.*$", "match image names")
	DockerHost   = flag.String("docker", "tcp://127.0.0.1:2375", "docker endpoint")
)

func main() {
	flag.Usage = usage
	flag.Parse()

	logOut := log.New(os.Stdout, "", log.LstdFlags)

	if *DryRun {
		logOut.Println("[INFO] dry run enabled")
	}

	drydock, err := NewDryDock(*DockerHost)
	if err != nil {
		log.Fatalf("[FATAL] drydock: %s\n", err)
	}

	drydock.Age, err = time.ParseDuration(*ImageAge)
	if err != nil {
		log.Fatalf("[FATAL] age: %s\n", err)
	}

	drydock.Keep = *ImagesToKeep

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

	logOut.Printf("[INFO] %d images scheduled for deletion\n", len(images))
	for _, image := range images {
		if imagesInUse.Exist(image) {
			log.Printf("[WARN] skipping %s, in use\n", image)
			continue
		}

		if !*DryRun {
			err := drydock.RemoveImage(image)
			if err != nil {
				log.Fatalf("[FATAL] remove image: %s\n", err)
			}
		}

		logOut.Printf("[INFO] deleted image %s\n", image)
	}
}
