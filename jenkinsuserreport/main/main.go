package main

import (
	"flag"
	"github.com/perolo/confluence-scripts/jenkinsuserreport"
)

func main() {
	propPtr := flag.String("prop", "jenkins.properties", "a properties file")

	flag.Parse()
	jenkinsuserreport.JenkinsUserReport(*propPtr)
}
