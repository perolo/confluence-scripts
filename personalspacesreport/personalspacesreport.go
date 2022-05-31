package personalspacesreport

import (
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/magiconair/properties"
	goconfluence "github.com/perolo/confluence-go-api"
	"github.com/perolo/excellogger"
	"golang.org/x/exp/slices"
)

type ReportConfig struct {
	ConfHost      string `properties:"confhost"`
	ConfUser      string `properties:"confuser"`
	ConfPass      string `properties:"confpass"`
	ConfToken     string `properties:"conftoken"`
	UseToken      bool   `properties:"usetoken"`
	File          string `properties:"file"`
	Simple        bool   `properties:"simple"`
	Report        bool   `properties:"report"`
	Space         string `properties:"space"`
	AncestorTitle string `properties:"ancestortitle"`
	//	Bindusername string `properties:"bindusername"`
	//	Bindpassword string `properties:"bindpassword"`
	//	BaseDN       string `properties:"basedn"`
}

func PersonalSpaceReport(propPtr string) {

	flag.Parse()

	p := properties.MustLoadFile(propPtr, properties.ISO_8859_1)

	// or through Decode
	var cfg ReportConfig
	if err := p.Decode(&cfg); err != nil {
		log.Fatal(err)
	}
	/*
		adutils.InitAD(cfg.Bindusername, cfg.Bindpassword)

	*/
	if cfg.Simple {
		cfg.File = fmt.Sprintf(cfg.File, "-"+"PersonalSpaces")
		CreatePersonalSpacesReport(cfg)
	} else {
		reportBase := cfg.File
		cfg.File = fmt.Sprintf(reportBase, "-"+"PersonalSpaces")
		CreatePersonalSpacesReport(cfg)
	}
	//	adutils.CloseAD()

}

func CreatePersonalSpacesReport(cfg ReportConfig) {
	var confClient *goconfluence.API
	var err error
	if cfg.UseToken {
		confClient, err = goconfluence.NewAPI(cfg.ConfHost, "", cfg.ConfToken)
	} else {
		confClient, err = goconfluence.NewAPI(cfg.ConfHost, cfg.ConfUser, cfg.ConfPass)
	}
	if err != nil {
		log.Fatal(err)
	}
	//	confClient.Debug = true

	excellogger.NewFile(nil)

	excellogger.SetCellFontHeader()
	excellogger.WiteCellln("Introduction")
	excellogger.WiteCellln("Please Do not edit this page!")
	excellogger.WiteCellln("This page is created by the User Report script: " + "https://github.com/perolo/confluence-scripts" + "/" + "PersonalSpacesReport")
	t := time.Now()
	excellogger.WiteCellln("Created by: " + cfg.ConfUser + " : " + t.Format(time.RFC3339))
	excellogger.WiteCellln("")

	excellogger.SetCellFontHeader2()
	excellogger.WiteCellln("Users and Permissions")
	excellogger.NextLine()
	excellogger.AutoFilterStart()
	excellogger.SetTableHeader()
	excellogger.WiteCell("Space Name")
	excellogger.SetTableHeader()
	excellogger.NextCol()
	excellogger.SetTableHeader()
	excellogger.WiteCell("Space Key")
	excellogger.NextCol()
	excellogger.SetTableHeader()
	excellogger.WiteCell("Type")
	excellogger.NextCol()
	excellogger.SetTableHeader()
	excellogger.WiteCell("Name")
	excellogger.NextCol()
	/*
		excellogger.SetTableHeader()
		excellogger.WiteCell("DN")
		excellogger.NextCol()
		excellogger.SetTableHeader()
		excellogger.WiteCell("Mail")
		excellogger.NextCol()
		excellogger.SetTableHeader()
		excellogger.WiteCell("Comment")
		excellogger.NextCol()
		excellogger.NextLine()
	*/
	noSpaces := 0
	spstart := 0
	spincrease := 50
	spcont := true
	var spaces *goconfluence.AllSpaces
	types, err := confClient.GetPermissionTypes()
	for _, t := range *types {
		excellogger.SetTableHeader()
		excellogger.WiteCell(t)
		excellogger.NextCol()
	}
	if err != nil {
		log.Fatal(err)
	}
	excellogger.NextLine()

	for spcont {
		spopt := goconfluence.AllSpacesQuery{Start: spstart, Limit: spincrease, Type: "personal", Status: "current"}
		spaces, _ = confClient.GetAllSpaces(spopt)
		opt := goconfluence.PaginationOptions{}
		for _, space := range spaces.Results {
			noSpaces++
			fmt.Printf("Space: %s \n", space.Name)
			start := 0
			cont := true
			increase := 50
			for cont {
				opt.StartAt = start
				opt.MaxResults = increase
				users, err2 := confClient.GetAllUsersWithAnyPermission(space.Key, &opt)
				if err2 != nil {
					log.Fatal(err2)
				}
				excellogger.NextCol()
				for _, user := range users.Users {
					permissions, _ := confClient.GetUserPermissionsForSpace(space.Key, user)
					excellogger.ResetCol()
					excellogger.WiteCellnc(space.Name)
					excellogger.WiteCellnc(space.Key)
					excellogger.WiteCellnc("User")
					excellogger.WiteCellnc(user)
					for _, t := range *types {
						if slices.Contains(permissions.Permissions, t) {
							excellogger.WiteCellnc("x")
						} else {
							excellogger.WiteCellnc("-")
						}
					}

					excellogger.NextLine()
					/*
						if Contains(permissions.Permissions, "SETPAGEPERMISSIONS") {
							_, err := adutils.GetActiveUserDN(user, cfg.BaseDN)
							if err == nil {
								//excellogger.WiteCellnc(dn.DN)
								//excellogger.WiteCellnc(dn.Mail)
								//excellogger.WiteCellnc("")
							} else {
								excellogger.ResetCol()
								excellogger.WiteCellnc(space.Name)
								excellogger.WiteCellnc(space.Key)
								excellogger.WiteCellnc("User")
								excellogger.WiteCellnc(user)
								udn, err := adutils.GetAllUserDN(user, cfg.BaseDN)
								if err == nil {
									excellogger.WiteCellnc(udn.DN)
									excellogger.WiteCellnc(udn.Mail)
									excellogger.WiteCellnc("Deactivated!")
								} else {
									excellogger.WiteCellnc("")
									excellogger.WiteCellnc("")
									excellogger.WiteCellnc("Not Found!")
								}
								excellogger.NextLine()
							}
						}*/
				}
				if users.Total < int64(start+increase) {
					cont = false
				} else {
					start = start + increase
				}
			}
		}
		if spaces.Size < int64(spincrease) {
			spcont = false
		} else {
			spstart = spstart + spincrease
		}
	}
	excellogger.SetAutoColWidth()
	excellogger.AutoFilterEnd()

	excellogger.SetColWidth("A", "A", 60)
	// Save xlsx file by the given path.
	excellogger.SaveAs(cfg.File)
	if cfg.Report {
		if cfg.Report {
			res, err := confClient.GetPageId(cfg.Space, cfg.AncestorTitle)
			if err == nil {
				if res.Size == 1 {
					err := confClient.UppdateAttachment(cfg.Space, cfg.AncestorTitle, cfg.File)
					if err != nil {
						panic(err)
					}
				}
			}
		}
	}
}
