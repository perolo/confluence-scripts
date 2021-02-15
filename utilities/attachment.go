package utilities

import (
	"encoding/json"
	"git.aa.st/perolo/confluence-utils/Utilities"
	"github.com/kennygrant/sanitize"
	"github.com/perolo/confluence-prop/client"
	"io/ioutil"
	"log"
	"os"
	"fmt"
)

func CreateAttachmentAndUpload(data interface{}, copt client.OperationOptions, confluence *client.ConfluenceClient, comment string) {
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

	AddAttachmentAndUpload(confluence, copt, attname, ff.Name(), comment)
}

func AddAttachmentAndUpload(confluence *client.ConfluenceClient, copt client.OperationOptions, attname string, fname string, comment string) {
	results := confluence.SearchPages(copt.Title, copt.SpaceKey)
	if results.Size == 1 {
		attId, _, err := confluence.GetPageAttachmentById(results.Results[0].ID, attname)
		if err != nil {
			if attId != nil && attId.Size == 0 {
				_, _, err = confluence.AddAttachment(results.Results[0].ID, attname, fname, comment)
				if err != nil {
					log.Fatal(err)
				} else {
					fmt.Printf("Added attachment to page: %s \n", copt.Title)
				}
			} else {
				log.Fatal(err)
			}
		} else {
			_, _, err = confluence.UpdateAttachment(results.Results[0].ID, attId.Results[0].ID, attname, fname, comment)
			if err != nil {
				log.Fatal(err)
			}
		}
	} else {
		fmt.Printf("Failed to find Page: %s \n", copt.Title)
	}
}

