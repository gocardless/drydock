package main

import (
	"log"
	"os"
	"regexp"
	"sort"
	"time"

	"github.com/fsouza/go-dockerclient"
)

type DryDock struct {
	Age     time.Duration
	Keep    int
	Pattern *regexp.Regexp
	docker  *docker.Client
	logOut  *log.Logger
}

// create a new DryDock assignment connected to Docker server
func NewDryDock(endpoint string) (*DryDock, error) {
	client, err := docker.NewClient(endpoint)
	if err != nil {
		return nil, err
	}

	return &DryDock{
		Age:     48 * time.Hour,
		Pattern: regexp.MustCompile(`^.*$`),

		docker: client,
		logOut: log.New(os.Stdout, "", log.LstdFlags),
	}, nil
}

// Sorts docker.APIImages, newest first
type byCreatedNewestFirst []docker.APIImages

func (a byCreatedNewestFirst) Len() int           { return len(a) }
func (a byCreatedNewestFirst) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a byCreatedNewestFirst) Less(i, j int) bool { return a[i].Created > a[j].Created }

// list images older than Age, excluding the newest Keep images
func (dd DryDock) ListImages() (Images, error) {
	images, err := dd.docker.ListImages(docker.ListImagesOptions{All: false})
	if err != nil {
		return nil, err
	}

	var imagesMatchingFilter []docker.APIImages
	for _, image := range images {
		if dd.matchRepoTags(image) {
			imagesMatchingFilter = append(imagesMatchingFilter, image)
		}
	}
	dd.logOut.Printf("[INFO] %d images matched pattern", len(imagesMatchingFilter))

	sort.Sort(byCreatedNewestFirst(imagesMatchingFilter))

	cutoff := time.Now().Add(-(dd.Age))
	var deleteStartingAtIndex int
	for i, image := range imagesMatchingFilter {
		created := time.Unix(image.Created, 0)

		if i >= dd.Keep && created.Before(cutoff) {
			deleteStartingAtIndex = i
			break
		}

		if i == len(imagesMatchingFilter)-1 {
			dd.logOut.Printf("[INFO] No images met keep/age criteria")
			return Images{}, nil
		}
	}

	imagesForDeletion := imagesMatchingFilter[deleteStartingAtIndex:len(imagesMatchingFilter)]
	dd.logOut.Printf("[INFO] %d images met keep/age criteria", len(imagesForDeletion))

	var imageIds Images
	for _, image := range imagesForDeletion {
		imageIds = append(imageIds, image.ID)
	}

	return imageIds, nil
}

// list images used by containers
func (dd DryDock) ListInUseImages() (Images, error) {
	containers, err := dd.docker.ListContainers(docker.ListContainersOptions{All: true})
	if err != nil {
		return nil, err
	}

	var imgs Images
	for _, container := range containers {
		image, err := dd.docker.InspectImage(container.Image)
		if err != nil {
			return nil, err
		}

		imgs = append(imgs, image.ID)
	}

	return imgs, nil
}

// remote an image by id
func (dd DryDock) RemoveImage(id string) error {
	return dd.docker.RemoveImage(id)
}

// match an images repo tags
func (dd DryDock) matchRepoTags(image docker.APIImages) bool {
	for _, tag := range image.RepoTags {
		if dd.Pattern.MatchString(tag) {
			return true
		}
	}

	return false
}
