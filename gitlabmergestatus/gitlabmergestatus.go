package gitlabmergestatus

import (
	"fmt"
	"git.aa.st/perolo/confluence-utils/Utilities"
	"github.com/magiconair/properties"
	"github.com/xanzy/go-gitlab"
	"log"
	"sort"
)

// or through Decode
type Config struct {
	GitLabHost  string `properties:"gitlabhost"`
	GitLabtoken string `properties:"gitlabtoken"`
	GitProjId   int    `properties:"gitlabprojid"`
	GitGroupId  int    `properties:"gitlabgroupid"`
}

var cfg Config

func GitLabMergeReport(propPtr string) {

	fmt.Printf("%%%%%%%%%%  GitLabMergeReport %%%%%%%%%%%%%%\n")

	p := properties.MustLoadFile(propPtr, properties.ISO_8859_1)

	if err := p.Decode(&cfg); err != nil {
		log.Fatal(err)
	}

	gitlabclient, err := gitlab.NewClient(cfg.GitLabtoken, gitlab.WithBaseURL(cfg.GitLabHost))
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	count := make(map[string]int)

	state := "opened"
	var opt2 gitlab.ListProjectMergeRequestsOptions
	opt2.Page = 0
	opt2.PerPage = 100
	opt2.State = &state
	//opt2 = gitlab.ListProjectMergeRequestsOptions{State: &state, Page: 0}
	merges, _, err := gitlabclient.MergeRequests.ListProjectMergeRequests(cfg.GitProjId, &opt2, nil)
	Utilities.Check(err)
	//openmenmerges := len(merges)

	for _, merge := range merges {

		fmt.Printf("Merge: %s Author: %s Upvotes: %d Donwwotes: %d\n", merge.Title, merge.Author.Name, merge.Upvotes, merge.Downvotes)
		participants, _, err := gitlabclient.MergeRequests.GetMergeRequestParticipants(cfg.GitProjId, merge.Iid, nil)
		Utilities.Check(err)
		for _, participant := range participants {
			fmt.Printf("  participant: %s\n", participant.Name)
			if _, ok := count[participant.Name]; !ok {
				count[participant.Name] = 1
			} else {
				count[participant.Name] = count[participant.Name] + 1
			}
		}
	}
	fmt.Printf("\n")
	fmt.Printf("Top List: \n")

	type kv struct {
		Key   string
		Value int
	}

	var ss []kv
	for k, v := range count {
		ss = append(ss, kv{k, v})
	}
	sort.Slice(ss, func(i, j int) bool {
		return ss[i].Value > ss[j].Value
	})

	for _, kv := range ss {
		fmt.Printf("%s, %d\n", kv.Key, kv.Value)
	}


	//fmt.Printf("Hello %i \n", openmenmerges)
}
