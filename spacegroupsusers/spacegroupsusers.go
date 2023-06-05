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

type ReportConfig struct {
	ConfHost      string `properties:"confhost"`
	User          string `properties:"user"`
	Pass          string `properties:"pass"`
	File          string `properties:"file"`
	Space         string `properties:"space"`
	Report        bool   `properties:"report"`
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

const MAX_GROUPS = 40

func CreateGroupUserReport(cfg ReportConfig) { //nolint:funlen

	var confluence *goconfluence.API
	var err error
	confluence, err = goconfluence.NewAPI(cfg.ConfHost, cfg.User, cfg.Pass)
	if err != nil {
		log.Fatal(err)
	}

	var gopt goconfluence.GetGroupMembersOptions
	gopt.Start = 0
	gopt.Limit = MAX_GROUPS

	groups, err2 := confluence.GetGroups(&gopt)
	if err2 != nil {
		log.Fatal(err2)
	}
	if groups.Size > MAX_GROUPS {
		log.Fatal("Too many groups")
	}

	// TODO Assuming max MAX_GROUPS groups
	type bvec [MAX_GROUPS]bool
	usermap := make(map[string]bvec)
	indexgroup := 0
	for _, group := range groups.Groups {
		gopt.Start = 0
		gopt.Limit = MAX_GROUPS * 10
		users, err3 := confluence.GetGroupMembers(group.Name, &gopt)
		if err3 != nil {
			log.Fatal(err3)
		}
		if users.Size > MAX_GROUPS*10 {
			log.Fatal("Too many users")
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

	excellogger.NewFile(nil)
	excellogger.SetCellFontHeader()
	excellogger.WiteCellln("Introduction")
	excellogger.WiteCellln("Please Do not edit this page!")
	excellogger.WiteCellln("This page is created by the User Report script: " + "https://github.com/perolo/confluence-scripts" + "/" + "SpaceGroupsUsersReport")
	t := time.Now()
	excellogger.WiteCellln("Created by: " + cfg.User + " : " + t.Format(time.RFC3339))
	excellogger.WiteCellln("")

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

	excellogger.SetColWidth("A", "A", 40)
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
