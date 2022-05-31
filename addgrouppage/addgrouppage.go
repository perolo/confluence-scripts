package addgrouppage

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/magiconair/properties"
	goconfluence "github.com/perolo/confluence-go-api"
	"github.com/perolo/confluence-scripts/utilities/htmlutils"
)

type Config struct {
	ConfHost      string `properties:"confhost"`
	ConfUser      string `properties:"confuser"`
	ConfPass      string `properties:"confpass"`
	ConfToken     string `properties:"conftoken"`
	UseToken      bool   `properties:"usetoken"`
	Space         string `properties:"space"`
	AncestorTitle string `properties:"ancestortitle"`
}

func AddGroupPage(propPtr string) {
	var cfg Config

	fmt.Printf("%%%%%%%%%%  Add Group Page %%%%%%%%%%%%%%\n")

	p := properties.MustLoadFile(propPtr, properties.ISO_8859_1)

	if err := p.Decode(&cfg); err != nil {
		log.Fatal(err)
	}

	fmt.Printf("->Connecting to %s\n", cfg.ConfHost)
	fmt.Printf("		->as %s\n", cfg.ConfUser)
	fmt.Printf("		->pass %s\n", cfg.ConfPass)
	fmt.Printf("		->using token:  %t\n", cfg.UseToken)

	var confluence *goconfluence.API
	var err error
	if cfg.UseToken {
		confluence, err = goconfluence.NewAPI(cfg.ConfHost, "", cfg.ConfToken)
	} else {
		confluence, err = goconfluence.NewAPI(cfg.ConfHost, cfg.ConfUser, cfg.ConfPass)
	}
	if err != nil {
		log.Fatal(err)
	}

	var gropt goconfluence.GetGroupMembersOptions
	grcont := true
	grstart := 0
	grmax := 50
	for grcont {
		gropt.StartAt = grstart
		gropt.MaxResults = grmax

		groups, _ := confluence.GetGroups(&gropt)
		if groups.Status == "success" {
			for _, group := range groups.Groups {

				f, err := ioutil.TempFile(os.TempDir(), "page*.html")
				if err != nil {
					return
				}
				if f == nil {
					return
				}

				defer f.Close()

				htmlutils.WriteHeader2(f, "Introduction")
				htmlutils.WriteParagraf(f, "")
				htmlutils.WriteParagraf(f, "On this page the members of the confluence group is displayed using the User List macro")
				htmlutils.WriteParagraf(f, "The page displays the current membership - automatically updated")
				htmlutils.WriteParagraf(f, "")

				block := fmt.Sprintf("<p> <ac:structured-macro ac:macro-id=\"f8e08252-b52a-4740-a564-2ad87efa70f6\" "+
					"ac:name=\"userlister\" ac:schema-version=\"1\"> <ac:parameter ac:name=\"groups\">%s</ac:parameter>	"+
					"</ac:structured-macro> 	</p>", group)
				htmlutils.WriteParagraf(f, block)
				htmlutils.WriteParagraf(f, "")
				pagename := "Group: " + group
				ancestid, err2 := confluence.GetPageId(cfg.Space, cfg.AncestorTitle)
				if err2 == nil && ancestid.Size == 1 {
					s, _ := confluence.GetPageId(cfg.Space, pagename)
					if s.Size == 0 {
						fmt.Printf("Creating page: %s \n", pagename)
						confluence.AddPage(pagename, cfg.Space, f.Name(), true, true, ancestid.Results[0].ID)
					} else {
						fmt.Printf("Skip creating page: %s \n", pagename)
					}

				} else {
					return
				}
			}
		}
		if len(groups.Groups) != grmax {
			grcont = false
		} else {
			grstart = grstart + grmax
		}
	}
}
