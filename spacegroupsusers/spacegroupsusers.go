package spacegroupsusers

import (
	"flag"
	"fmt"
	"log"
	"time"

	goconfluence "github.com/perolo/confluence-go-api"

	"github.com/magiconair/properties"
	"github.com/perolo/excellogger"
)

// Contains tells whether a contains x.
func Contains(a []string, x string) bool {
	for _, n := range a {
		if x == n {
			return true
		}
	}
	return false
}

type ReportConfig struct {
	ConfHost      string `properties:"confhost"`
	ConfUser      string `properties:"confuser"`
	ConfPass      string `properties:"confpass"`
	ConfToken     string `properties:"conftoken"`
	UseToken      bool   `properties:"usetoken"`
	Groups        bool   `properties:"groups"`
	Users         bool   `properties:"users"`
	File          string `properties:"file"`
	Report        bool   `properties:"report"`
	Space         string `properties:"space"`
	AncestorTitle string `properties:"ancestortitle"`
}

func SpaceGroupsUsersReport(propPtr string) {

	flag.Parse()

	p := properties.MustLoadFile(propPtr, properties.ISO_8859_1)

	// or through Decode
	var cfg ReportConfig
	if err := p.Decode(&cfg); err != nil {
		log.Fatal(err)
	}
	cfg.File = fmt.Sprintf(cfg.File, "-GroupsUsers")
	CreateGroupUserReport(cfg)
}

func addUserToGroup(user string, group string) {

}

func CreateGroupUserReport(cfg ReportConfig) { //nolint:funlen

	excellogger.NewFile(nil)

	excellogger.SetCellFontHeader()
	excellogger.WiteCellln("Introduction")
	excellogger.WiteCellln("Please Do not edit this page!")
	excellogger.WiteCellln("This page is created by the User Report script: " + "https://github.com/perolo/confluence-scripts" + "/" + "SpaceGroupsUsersReport")
	t := time.Now()
	excellogger.WiteCellln("Created by: " + cfg.ConfUser + " : " + t.Format(time.RFC3339))
	excellogger.WiteCellln("")

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

	groups, err2 := confluence.GetGroups(nil)
	if err2 != nil {
		log.Fatal(err2)
	}

	type bvec [40]bool
	usermap := make(map[string]bvec)
	indexgroup := 0
	for _, group := range groups.Groups {

		users, err3 := confluence.GetGroupMembers(group.Name)
		if err3 != nil {
			log.Fatal(err3)
		}
		for _, member := range users.Members {
			if _, ok := usermap[member.DisplayName]; !ok {
				var newvec bvec
				newvec[indexgroup] = true
				usermap[member.DisplayName] = newvec
			} else {
				oldvec := usermap[member.DisplayName]
				oldvec[indexgroup] = true
				usermap[member.DisplayName] = oldvec
			}
		}
		indexgroup++
	}

	excellogger.SetCellFontHeader2()
	excellogger.WiteCellln("Users and Groups")
	excellogger.NextLine()
	excellogger.AutoFilterStart()
	excellogger.SetTableHeader()
	excellogger.WiteCellnc("Display Name")
	for _, group := range groups.Groups {
		excellogger.SetTableHeader()
		excellogger.WiteCellnc(group.Name)
	}
	excellogger.NextLine()
	for user, binmap := range usermap {
		excellogger.WiteCellnc(user)
		for b, _ := range groups.Groups {
			excellogger.WiteCellnc(excellogger.BoolToEmoji(binmap[b]))
		}
		excellogger.NextLine()
	}

	excellogger.SetAutoColWidth()
	excellogger.AutoFilterEnd()

	//excellogger.SetColWidth("A", "A", 40)
	// Save xlsx file by the given path.
	excellogger.SaveAs(cfg.File)
	if cfg.Report {
		res, err := confluence.GetPageId(cfg.Space, cfg.AncestorTitle)
		if err == nil {
			if res.Size == 1 {
				err := confluence.UppdateAttachment(cfg.Space, cfg.AncestorTitle, cfg.File)
				if err != nil {
					panic(err)
				}
			}
		}
	}
}
