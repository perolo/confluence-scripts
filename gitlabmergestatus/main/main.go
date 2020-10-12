package main

import (
	"flag"
	"github.com/perolo/confluence-scripts/gitlabmergestatus"
)

func main() {
	propPtr := flag.String("prop", "gitlabmergestatus.properties", "a properties file")

	flag.Parse()
	gitlabmergestatus.GitLabMergeReport(*propPtr)
}
