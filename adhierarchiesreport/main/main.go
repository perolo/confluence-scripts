package main

import (
	"flag"
	"github.com/perolo/confluence-scripts/adhierarchiesreport"
)

func main() {
	propPtr := flag.String("prop", "confluence.properties", "a properties file")

	adhierarchiesreport.CreateAdHierarchiesReport(*propPtr, "#AAAB - Group Technology Team")
}
