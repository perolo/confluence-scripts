package main

import (
	"flag"

	"github.com/perolo/confluence-scripts/spacegroupsusers"
)

func main() {
	propPtr := flag.String("prop", "spacegroupsusers.properties", "a string")

	spacegroupsusers.SpaceGroupsUsersReport(*propPtr)
}
