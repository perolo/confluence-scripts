package main

import (
	"flag"
	sonarqubeuserreport "github.com/perolo/confluence-scripts/sonarqubeuserreport"
)

func main() {
	propPtr := flag.String("prop", "sonar.properties", "a string")
	flag.Parse()

	sonarqubeuserreport.Sonarqubeuserreport(*propPtr)
}
