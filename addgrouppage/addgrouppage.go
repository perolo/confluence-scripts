package addgrouppage

import (
	"fmt"
	"github.com/magiconair/properties"
	"github.com/perolo/confluence-prop/client"
	"github.com/perolo/confluence-scripts/utilities"
	"github.com/perolo/confluence-scripts/utilities/htmlutils"
	"io/ioutil"
	"log"
	"os"
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
	var copt client.OperationOptions

	fmt.Printf("%%%%%%%%%%  Add Group Page %%%%%%%%%%%%%%\n")

	p := properties.MustLoadFile(propPtr, properties.ISO_8859_1)

	if err := p.Decode(&cfg); err != nil {
		log.Fatal(err)
	}

	var config = client.ConfluenceConfig{}

	if cfg.UseToken {
		cfg.ConfPass = cfg.ConfToken
	}

	config.Username = cfg.ConfUser
	config.Password = cfg.ConfPass
	config.UseToken = cfg.UseToken
	config.URL = cfg.ConfHost
	config.Debug = false
	config.UseToken = cfg.UseToken

	fmt.Printf("->Connecting to %s\n", cfg.ConfHost)
	fmt.Printf("		->as %s\n", cfg.ConfUser)
	fmt.Printf("		->pass %s\n", cfg.ConfPass)
	fmt.Printf("		->using token:  %t\n", cfg.UseToken)

	confluence := client.Client(&config)

	copt.SpaceKey = cfg.Space
	copt.AncestorTitle = cfg.AncestorTitle

	var gropt client.GetGroupMembersOptions
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
				//	group := "gtt-all"
				copt.Title = "Group: " + group

				copt.Filepath = f.Name()
				copt.BodyOnly = true

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
				//	t := time.Now()

				if !utilities.CheckPageExists2(copt, confluence) {
					if confluence.AddOrUpdatePage(copt) {
						fmt.Printf("%s uploaded ok", copt.Title)
					}
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
