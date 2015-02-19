package main

// array of image ids
type Images []string

// check if an image id exists in images list
func (imgs Images) Exist(id string) bool {
	for _, img := range imgs {
		if img == id {
			return true
		}
	}

	return false
}
