package utilities

import (
	"encoding/json"
	"fmt"
	"github.com/kennygrant/sanitize"
	"github.com/perolo/confluence-prop/client"
	"github.com/perolo/confluence-scripts/utilities/htmlutils"
	"github.com/perolo/jira-scripts/jirautils"
	"io/ioutil"
	"log"
	"net/url"
	"os"
)

func CheckPageExists(copt client.OperationOptions, confluence *client.ConfluenceClient) {

	opt2 := client.PageOptions{Start: 0, Limit: 10}
	/*
		https://confluence.assaabloy.net/rest/api/content?spaceKey=STPIM&expand=metadata.labels
		https://confluence.assaabloy.net/rest/api/content?
	*/
	content := fmt.Sprintf("spaceKey=%s&title=%s", copt.SpaceKey, url.QueryEscape(copt.Title))

	pages := confluence.GetContent(content, &opt2)

	for _, page := range pages.Results {
		fmt.Printf("Pages name: %s type: %s\n", page.Title, page.Type)
	}

	if len(pages.Results) == 0 {

		f, err := ioutil.TempFile(os.TempDir(), "page*.html")
		if err != nil {
			return
		}
		if f == nil {
			return
		}

		defer f.Close()

		//	defer os.Remove(f.Name())
		//		var copt client.OperationOptions

		//		copt.Title = "Group User Report: " + group
		//		copt.SpaceKey = cfg.Space
		//		Utilities.Check(err)
		copt.Filepath = f.Name()
		copt.BodyOnly = true

		htmlutils.WriteHeader2(f, "Introduction")
		htmlutils.WriteParagraf(f, "Please Do not edit this page!")
		htmlutils.WriteParagraf(f, "This page is created by the Ad User Report script: "+htmlutils.WrapLink("https://github.com/perolo/confluence-scripts", "confluence-utils"))
		htmlutils.WriteParagraf(f, "The report is uploaded as attachment to this page")

		confluence.AddPage(copt.Title, copt.SpaceKey, copt.Filepath, true, false, 0)

		/*
			if confluence.AddOrUpdatePage(copt) {
				fmt.Printf("%s uploaded ok", copt.Title)
			}
		*/
	}
}

func CheckPageExists2(copt client.OperationOptions, confluence *client.ConfluenceClient) bool {

	opt2 := client.PageOptions{Start: 0, Limit: 10}
	/*
		https://confluence.assaabloy.net/rest/api/content?spaceKey=STPIM&expand=metadata.labels
		https://confluence.assaabloy.net/rest/api/content?
	*/
	content := fmt.Sprintf("spaceKey=%s&title=%s", copt.SpaceKey, url.QueryEscape(copt.Title))

	pages := confluence.GetContent(content, &opt2)

	for _, page := range pages.Results {
		fmt.Printf("Pages name: %s type: %s\n", page.Title, page.Type)
	}

	return len(pages.Results) != 0
}

func CreateAttachmentAndUpload(data interface{}, copt client.OperationOptions, confluence *client.ConfluenceClient, comment string) error {
	buf, err := json.Marshal(data)
	if err != nil {
		log.Fatal(err)
	}

	attname := sanitize.BaseName(copt.Title) + ".json"
	ff, err := ioutil.TempFile(os.TempDir(), attname)
	//ff, err := os.Create("C://temp/" + attname)
	jirautils.Check(err)
	_, err = ff.Write(buf)
	jirautils.Check(err)
	err = ff.Close()
	jirautils.Check(err)

	return AddAttachmentAndUpload(confluence, copt, attname, ff.Name(), comment)
}

func AddAttachmentAndUpload(confluence *client.ConfluenceClient, copt client.OperationOptions, attname string, fname string, comment string) error {
	//TODO Refactor to simplify, why copt?
	results := confluence.SearchPages(copt.Title, copt.SpaceKey)
	if results.Size == 1 {
		attId, _, err := confluence.GetPageAttachmentById(results.Results[0].ID, attname)
		if err != nil {
			if attId != nil && attId.Size == 0 {
				_, _, err = confluence.AddAttachment(results.Results[0].ID, attname, fname, comment)
				if err != nil {
					return fmt.Errorf("Failed to add attachemt to Page: %s err: %s \n", copt.Title, err)
				} else {
					fmt.Printf("Added attachment to page: %s \n", copt.Title)
				}
			} else {
				return fmt.Errorf("Failed to add attachment to Page: %s err: %s \n", copt.Title, err)
			}
		} else {
			_, _, err = confluence.UpdateAttachment(results.Results[0].ID, attId.Results[0].ID, attname, fname, comment)
			if err != nil {
				return fmt.Errorf("Failed to update attachment to Page: %s err: %s \n", copt.Title, err)
			}
		}
	} else {
		return fmt.Errorf("Failed to find Confluence Page: \"%s\" in Space: \"%s\" \n", copt.Title, copt.SpaceKey)
	}
	return nil
}
