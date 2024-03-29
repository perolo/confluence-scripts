package gitlabmergeflow

import (
	"fmt"
	"github.com/magiconair/properties"
	"github.com/perolo/confluence-client/client"
	"github.com/perolo/confluence-scripts/gitlabmergestatus"
	"github.com/perolo/confluence-scripts/utilities"
	"github.com/perolo/go-gitlab"
	"github.com/perolo/jira-scripts/jirautils"
	"log"
	"sort"
	"time"
)

type Config struct {
	ConfUser  string `properties:"confuser"`
	ConfPass  string `properties:"confpass"`
	ConfHost  string `properties:"confhost"`
	ConfToken string `properties:"conftoken"`
	UseToken  bool   `properties:"usetoken"`
}
type MergConfig struct {
	GitLabHost       string `properties:"gitlabhost"`
	GitLabtoken      string `properties:"gitlabtoken"`
	GitProjId        int    `properties:"gitlabprojid"`
	PageName         string `properties:"pagename"`
	CreateAttachment bool   `properties:"createattachment"`
}

type Data struct {
	Project     string            `json:"project"`
	Description string            `json:"description"`
	Merges      []MergeData       `json:"mergerequests"`
	TopList     []ContributorData `json:"toplist"`
}

type MergeData struct {
	MergeRequest string `json:"mergerequest"`
	Author       string `json:"author"`
	Start        string `json:"start"`
	End          string `json:"end"`
	State        string `json:"state"`
	Status       string `json:"status"`
	Title        string `json:"title"`
	Link         string `json:"link"`
	UpVotes      int    `json:"upvotes"`
	DownVotes    int    `json:"downvotes"`
}
type ContributorData struct {
	USer          string `json:"user"`
	Contributions int    `json:"contributions"`
}

func GitLabMergeFlowReport(propPtr string) {
	var cfg Config
	var mergecfg MergConfig
	var data Data
	var copt client.OperationOptions
	var confluence *client.ConfluenceClient
	fmt.Printf("%%%%%%%%%%  GitLab Merge Flow %%%%%%%%%%%%%%\n")

	p := properties.MustLoadFile(propPtr, properties.ISO_8859_1)

	if err := p.Decode(&cfg); err != nil {
		log.Fatal(err)
	}
	if cfg.UseToken {
		cfg.ConfPass = cfg.ConfToken
	}

	mergecfg.CreateAttachment = true
	if mergecfg.CreateAttachment {
		// Access Confluence
		var config = client.ConfluenceConfig{}
		config.Username = cfg.ConfUser
		config.Password = cfg.ConfPass
		config.UseToken = cfg.UseToken
		config.URL = cfg.ConfHost
		//config.Debug = true

		confluence = client.Client(&config)
	}

	for _, report := range gitlabmergestatus.Reports {
		mergecfg.PageName = report.PageName + "-flow"
		mergecfg.GitLabHost = report.Host
		mergecfg.GitProjId = report.ProjId
		mergecfg.GitLabtoken = report.Token
		createProjectReport(confluence, data, copt, mergecfg)
	}

}

func createProjectReport(confluence *client.ConfluenceClient, data Data, copt client.OperationOptions, cfg MergConfig) {
	if cfg.CreateAttachment {
		data.Project = fmt.Sprintf("Project Id: %v", cfg.GitProjId)
		data.Description = "Master branch Merge requests updated within last 28 days"
	}
	gitlabclient, err := gitlab.NewClient(cfg.GitLabtoken, gitlab.WithBaseURL(cfg.GitLabHost))
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	count := make(map[string]int)

	//state := "opened"
	var opt2 gitlab.ListProjectMergeRequestsOptions
	window := time.Now().AddDate(0, 0, -28)
	opt2.UpdatedAfter = &window
	master := "master"
	opt2.TargetBranch = &master

	cont := true
	page := 0
	opt2.PerPage = 100

	for cont {
		opt2.Page = page

		flowmerges, _, err2 := gitlabclient.MergeRequests.ListProjectMergeRequests(cfg.GitProjId, &opt2, nil)
		jirautils.Check(err2)

		for _, merge := range flowmerges {

			fmt.Printf("Merge: %s Author: %s Upvotes: %d Downvotes: %d\n", merge.Title, merge.Author.Name, merge.Upvotes, merge.Downvotes)
			//ListMergeRequestNotes
			//GET /projects/:id/merge_requests/:merge_request_iid/notes
			notes, _, err3 := gitlabclient.Notes.ListMergeRequestNotes(cfg.GitProjId, merge.IID, nil)
			jirautils.Check(err3)
			for _, note := range notes {
				if _, ok := count[note.Author.Name]; !ok {
					count[note.Author.Name] = 1
				} else {
					count[note.Author.Name] = count[note.Author.Name] + 1
				}
			}

			if cfg.CreateAttachment {
				var i MergeData
				i.Title = merge.Title
				i.Author = merge.Author.Name
				i.State = merge.State
				i.Status = merge.MergeStatus
				i.Link = merge.WebURL
				i.Start = merge.CreatedAt.Format("2006, 01, 02")
				if merge.State == "merged" {
					fmt.Printf("MergeStatus: %s State: %s Target: %s\n", merge.MergeStatus, merge.State, merge.TargetBranch)
					if merge.MergedAt == nil {
						// Why is a Merged branch without mergedate? This is a Workaround
						i.End = time.Now().Format("2006, 01, 02")
					} else {
						i.End = merge.MergedAt.Format("2006, 01, 02")
					}
				} else {
					fmt.Printf("MergeStatus: %s State: %s Target: %s\n", merge.MergeStatus, merge.State, merge.TargetBranch)
					i.End = time.Now().Format("2006, 01, 02")
				}
				i.UpVotes = merge.Upvotes
				i.DownVotes = merge.Downvotes
				data.Merges = append(data.Merges, i)
			}
		}
		if len(flowmerges) != opt2.PerPage {
			cont = false
		} else {
			page++
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
		if cfg.CreateAttachment {
			var i ContributorData
			i.USer = kv.Key
			i.Contributions = kv.Value
			data.TopList = append(data.TopList, i)
		}

	}
	if cfg.CreateAttachment {
		copt.Title = cfg.PageName
		copt.SpaceKey = "~per.olofsson@assaabloy.com"
		err := utilities.CreateAttachmentAndUpload(data, copt, confluence, "Created by GitLab Merge Flow Report")
		if err != nil {
			panic(err)
		}
	}
}
