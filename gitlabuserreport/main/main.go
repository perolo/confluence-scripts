package main

import (
	"flag"
	"github.com/perolo/confluence-scripts/gitlabuserreport"
)

func main() {
	propPtr := flag.String("prop", "gitlab.properties", "a properties file")

	flag.Parse()
	gitlabuserreport.GitLabUserReport(*propPtr)
}
