package gitlabmergestatus

import (
	"encoding/json"
	"fmt"
	"git.aa.st/perolo/confluence-utils/Utilities"
	"github.com/kennygrant/sanitize"
	"github.com/magiconair/properties"
	"github.com/perolo/confluence-prop/client"
	"github.com/xanzy/go-gitlab"
	"io/ioutil"
	"log"
	"os"
	"sort"
	"time"
)

// or through Decode
type Config struct {
	GitLabHost       string `properties:"gitlabhost"`
	GitLabtoken      string `properties:"gitlabtoken"`
	GitProjId        int    `properties:"gitlabprojid"`
//	GitGroupId       int    `properties:"gitlabgroupid"`
	PageName string `properties:"pagename"`
	CreateAttachment bool   `properties:"createattachment"`
	User             string `properties:"user"`
	Pass             string `properties:"password"`
	ConfHost         string `properties:"confhost"`
}

var cfg Config

type Data struct {
	Project string            `json:"project"`
	Merges  []MergeData       `json:"mergerequests"`
	TopList []ContributorData `json:"toplist"`
}

type MergeData struct {
	MergeRequest string `json:"mergerequest"`
	Author       string `json:"author"`
	Start        string `json:"start"`
	End          string `json:"end"`
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
	var data Data
	var copt client.OperationOptions
	var confluence *client.ConfluenceClient
	fmt.Printf("%%%%%%%%%%  GitLabMergeReport %%%%%%%%%%%%%%\n")

	p := properties.MustLoadFile(propPtr, properties.ISO_8859_1)

	if err := p.Decode(&cfg); err != nil {
		log.Fatal(err)
	}
	if cfg.CreateAttachment {
		// Access Confluence
		var config = client.ConfluenceConfig{}
		config.Username = cfg.User
		config.Password = cfg.Pass
		config.URL = cfg.ConfHost
		//config.Debug = true

		confluence = client.Client(&config)

		data.Project = fmt.Sprintf("Project Id: %v", cfg.GitProjId)
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
	openmerges, _, err := gitlabclient.MergeRequests.ListProjectMergeRequests(cfg.GitProjId, &opt2, nil)
	Utilities.Check(err)

	for _, merge := range openmerges {

		fmt.Printf("Merge: %s Author: %s Upvotes: %d Downvotes: %d\n", merge.Title, merge.Author.Name, merge.Upvotes, merge.Downvotes)
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
		if cfg.CreateAttachment {
			var i MergeData
			i.Title = merge.Title
			i.Author = merge.Author.Name
			i.Status = merge.MergeStatus
			i.Link = merge.WebURL
			i.Start = merge.CreatedAt.Format("2006, 01, 02")
			i.End = time.Now().Format("2006, 01, 02")
			i.UpVotes = merge.Upvotes
			i.DownVotes = merge.Downvotes
			data.Merges = append(data.Merges, i)
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
		CreateAttachmentAndUpload(data, copt, confluence)
	}

}

func CreateAttachmentAndUpload(data Data, copt client.OperationOptions, confluence *client.ConfluenceClient) {
	buf, err := json.Marshal(data)
	if err != nil {
		log.Fatal(err)
	}

	attname := sanitize.BaseName(copt.Title) + ".json"
	ff, err := ioutil.TempFile(os.TempDir(), attname)
	//ff, err := os.Create("C://temp/" + attname)
	Utilities.Check(err)
	_, err = ff.Write(buf)
	Utilities.Check(err)
	err = ff.Close()
	Utilities.Check(err)

	results := confluence.SearchPages(copt.Title, copt.SpaceKey)
	if results.Size == 1 {
		attId, _, err := confluence.GetPageAttachmentById(results.Results[0].ID, attname)
		if err != nil {
			if attId != nil && attId.Size == 0 {
				_, _, err = confluence.AddAttachment(results.Results[0].ID, attname, ff.Name(), "Added with theReport Report")
				if err != nil {
					log.Fatal(err)
				} else {
					fmt.Printf("Added attachment to page: %s \n", copt.Title)
				}
			} else {
				log.Fatal(err)
			}
		} else {
			_, _, err = confluence.UpdateAttachment(results.Results[0].ID, attId.Results[0].ID, attname, ff.Name(), "Updated with theReport Report")
			if err != nil {
				log.Fatal(err)
			}
		}
	} else {
		fmt.Printf("Failed to find Page: %s \n", copt.Title)
	}
}
