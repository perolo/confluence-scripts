package main

import (
	"flag"

	"github.com/perolo/confluence-scripts/addgrouppage"
)

func main() {
	propPtr := flag.String("prop", "addgrouppage.properties", "a string")
	flag.Parse()
	addgrouppage.AddGroupPage(*propPtr)

}
