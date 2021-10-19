package adhierarchiesreport

import (
	"flag"
	"fmt"
	"github.com/magiconair/properties"
	"github.com/perolo/ad-utils"
	"github.com/perolo/confluence-prop/client"
	"github.com/perolo/confluence-scripts/utilities"
	"log"
)

type Config struct {
	ConfHost        string `properties:"confhost"`
	User            string `properties:"user"`
	Pass            string `properties:"password"`
	Bindusername    string `properties:"bindusername"`
	Bindpassword    string `properties:"bindpassword"`
	BaseDN           string `properties:"basedn"`
}


func CreateAdHierarchiesReport(propPtr, adgroup string, expandUsers bool) {
	var copt client.OperationOptions
	var confluence *client.ConfluenceClient

	flag.Parse()
	p := properties.MustLoadFile(propPtr, properties.ISO_8859_1)
	var cfg Config
	if err := p.Decode(&cfg); err != nil {
		log.Fatal(err)
	}
	var config = client.ConfluenceConfig{}
	config.Username = cfg.User
	config.Password = cfg.Pass
	config.URL = cfg.ConfHost

	confluence = client.Client(&config)

	adutils.InitAD(cfg.Bindusername, cfg.Bindpassword)
	var roothier [] adutils.ADHierarchy
	var newhierarchy adutils.ADHierarchy
	newhierarchy.Name = adgroup
	newhierarchy.Parent = ""
	roothier = append(roothier, newhierarchy)

	groups, hier, err := adutils.ExpandHierarchy(adgroup, roothier, cfg.BaseDN)
	if err != nil {
		fmt.Printf("Failed to parse AD hierarchy : %s \n", err)
	} else {
		hier = append(hier, roothier...)
		fmt.Printf("adUnames(%v): %s \n", len(groups), groups)
		fmt.Printf("adUnames(%v): %s \n", len(hier), hier)
		copt.Title = "GTT Hierarchies - " + adgroup
		copt.SpaceKey = "~per.olofsson@assaabloy.com"
		if expandUsers {
			for _, h := range hier {
				users, _ := adutils.GetUnamesInGroup(h.Name, cfg.BaseDN)
				for _, u := range users {
					hier = append(hier,adutils.ADHierarchy{Name: u.Name, Parent: h.Name})
				}
			}

		}

		utilities.CheckPageExists(copt, confluence)
		err = utilities.CreateAttachmentAndUpload(hier, copt, confluence, "Created by AD Hierarchies Report")
		if err!= nil {
			panic(err)
		}

	}
	adutils.CloseAD()
}
