package gitlabmergestatus

import (
	"fmt"
	"git.aa.st/perolo/confluence-utils/Utilities"
	"github.com/magiconair/properties"
	"github.com/perolo/confluence-prop/client"
	"github.com/perolo/confluence-scripts/utilities"
	"github.com/xanzy/go-gitlab"
	"log"
	"sort"
	"time"
)

// or through Decode
type Config struct {
	User     string `properties:"user"`
	Pass     string `properties:"password"`
	ConfHost string `properties:"confhost"`
}
type MergConfig struct {
	GitLabHost       string `properties:"gitlabhost"`
	GitLabtoken      string `properties:"gitlabtoken"`
	GitProjId        int    `properties:"gitlabprojid"`
	PageName         string `properties:"pagename"`
	CreateAttachment bool   `properties:"createattachment"`
}

type Data struct {
	Project string            `json:"project"`
	Description string        `json:"description"`
	Merges  []MergeData       `json:"mergerequests"`
	TopList []ContributorData `json:"toplist"`
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

func GitLabMergeReport(propPtr string) {
	var cfg Config
	var mergecfg MergConfig
	var data Data
	var copt client.OperationOptions
	var confluence *client.ConfluenceClient
	fmt.Printf("%%%%%%%%%%  GitLab Merge Report %%%%%%%%%%%%%%\n")

	p := properties.MustLoadFile(propPtr, properties.ISO_8859_1)

	if err := p.Decode(&cfg); err != nil {
		log.Fatal(err)
	}
	mergecfg.CreateAttachment = true
	if mergecfg.CreateAttachment {
		// Access Confluence
		var config = client.ConfluenceConfig{}
		config.Username = cfg.User
		config.Password = cfg.Pass
		config.URL = cfg.ConfHost
		//config.Debug = true

		confluence = client.Client(&config)

	}

	for _, report := range Reports {
		mergecfg.PageName = report.PageName
		mergecfg.GitLabHost = report.Host
		mergecfg.GitProjId = report.ProjId
		mergecfg.GitLabtoken = report.Token
		createProjectReport(confluence, data, copt, mergecfg)
	}

}

func createProjectReport(confluence *client.ConfluenceClient, data Data, copt client.OperationOptions, cfg MergConfig) {
	if cfg.CreateAttachment {
		data.Project = fmt.Sprintf("Project Id: %v", cfg.GitProjId)
		data.Description = "Open Merge Requests"
	}
	gitlabclient, err := gitlab.NewClient(cfg.GitLabtoken, gitlab.WithBaseURL(cfg.GitLabHost))
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	count := make(map[string]int)

	state := "opened"
	var opt2 gitlab.ListProjectMergeRequestsOptions
	cont := true
	page := 0
	opt2.PerPage = 100
	opt2.State = &state
	for cont {
		opt2.Page = page
		openmerges, _, err := gitlabclient.MergeRequests.ListProjectMergeRequests(cfg.GitProjId, &opt2, nil)
		Utilities.Check(err)

		for _, merge := range openmerges {

			fmt.Printf("Merge: %s Author: %s Upvotes: %d Downvotes: %d\n", merge.Title, merge.Author.Name, merge.Upvotes, merge.Downvotes)
			participants, _, err := gitlabclient.MergeRequests.GetMergeRequestParticipants(cfg.GitProjId, merge.IID, nil)
			Utilities.Check(err)
			for _, participant := range participants {
				//fmt.Printf("  participant: %s\n", participant.Name)
				if _, ok := count[participant.Name]; !ok {
					count[participant.Name] = 1
				} else {
					count[participant.Name] = count[participant.Name] + 1
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
				i.End = time.Now().Format("2006, 01, 02")
				i.UpVotes = merge.Upvotes
				i.DownVotes = merge.Downvotes
				data.Merges = append(data.Merges, i)
			}

		}
		if len(openmerges) != opt2.PerPage {
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
		utilities.CreateAttachmentAndUpload(data, copt, confluence)
	}
}
