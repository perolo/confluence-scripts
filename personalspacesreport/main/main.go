package main

import (
	"flag"

	"github.com/perolo/confluence-scripts/personalspacesreport"
)

func main() {
	propPtr := flag.String("prop", "personalspacesreport.properties", "a string")

	personalspacesreport.PersonalSpaceReport(*propPtr)
}
