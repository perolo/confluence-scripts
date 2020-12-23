package main

import (
	"flag"
	"github.com/perolo/confluence-scripts/gitlabmergeflow"
)

func main() {
	propPtr := flag.String("prop", "gitlabmergestatus.properties", "a properties file")

	flag.Parse()
	gitlabmergeflow.GitLabMergeFlowReport(*propPtr)
}
