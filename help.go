package main

import (
	"fmt"
)

const VERSION = "0.0.3"

func usage() {
	fmt.Printf("DryDock %s\n", VERSION)
	fmt.Printf("usage: drydock [options]\n\n")

	fmt.Printf("Options:\n")
	fmt.Printf("  --dry-run                          don't delete images\n")
	fmt.Printf("  --age      <48h>                   delete images older than age\n")
	fmt.Printf("  --keep     <10>                    keep at least this many images\n")
	fmt.Printf("  --pattern  <^.*$>                  pattern for images to be deleted\n")
	fmt.Printf("  --docker   <tcp://127.0.0.1:2375>  docker host endpoint\n")
}
