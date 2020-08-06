package main

import (
	"flag"
	gitlabuserreport "git.aa.st/perolo/confluence-utils/GitLabUserReport"
)

func main() {
	propPtr := flag.String("prop", "gitlab.properties", "a properties file")

	flag.Parse()
	gitlabuserreport.GitLabUserReport(*propPtr)
}
