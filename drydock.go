package main

import (
	"regexp"
	"time"

	"github.com/fsouza/go-dockerclient"
)

type DryDock struct {
	Age     time.Duration  // image age
	Pattern *regexp.Regexp // repo tag pattern

	docker *docker.Client // docker client
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
	}, nil
}

// list images older than age
func (dd DryDock) ListImages() (Images, error) {
	images, err := dd.docker.ListImages(docker.ListImagesOptions{All: false})
	if err != nil {
		return nil, err
	}

	before := time.Now().Add(-(dd.Age))

	var imgs Images
	for _, image := range images {
		created := time.Unix(image.Created, 0)

		if before.Before(created) || !dd.matchRepoTags(image) {
			continue
		}

		imgs = append(imgs, image.ID)
	}

	return imgs, nil
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
