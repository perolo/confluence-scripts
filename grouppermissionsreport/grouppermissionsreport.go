package grouppermissionsreport

import (
	"flag"
	"fmt"
	"github.com/magiconair/properties"
	"github.com/perolo/confluence-client/client"
	"github.com/perolo/confluence-scripts/schedulerutil"
	"github.com/perolo/confluence-scripts/utilities"
	"github.com/perolo/excellogger"
	"log"
	"path/filepath"
	"time"
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
	ConfHost  string `properties:"confhost"`
	ConfUser  string `properties:"confuser"`
	ConfPass  string `properties:"confpass"`
	ConfToken string `properties:"conftoken"`
	UseToken  bool   `properties:"usetoken"`
	Groups    bool   `properties:"groups"`
	Users     bool   `properties:"users"`
	Group     string `properties:"group"`
	File      string `properties:"file"`
	Simple    bool   `properties:"simple"`
	Report    bool   `properties:"report"`
	Reset     bool   `properties:"reset"`
}

func SpacePermissionsReport(propPtr string) {

	flag.Parse()

	p := properties.MustLoadFile(propPtr, properties.ISO_8859_1)

	// or through Decode
	var cfg ReportConfig
	if err := p.Decode(&cfg); err != nil {
		log.Fatal(err)
	}
	if cfg.UseToken {
		cfg.ConfPass = cfg.ConfToken
	}

	if cfg.Simple {
		cfg.File = fmt.Sprintf(cfg.File, "-groupreport-"+cfg.Group)
		CreateGroupPermissionsReport(cfg)
	} else {
		reportBase := cfg.File
		for _, group := range Groups {
			if schedulerutil.CheckScheduleDetail("GrouPermissionsReport-"+group, 7*time.Hour*24, cfg.Reset, schedulerutil.DummyFunc, "jiracategory.properties") {
				cfg.Group = group
				cfg.File = fmt.Sprintf(reportBase, "-groupreport-"+group)
				fmt.Printf("Category: %s \n", group)
				CreateGroupPermissionsReport(cfg)
			}
		}
	}
}

func CreateGroupPermissionsReport(cfg ReportConfig) {

	excellogger.NewFile(nil)

	excellogger.SetCellFontHeader()
	excellogger.WiteCellln("Introduction")
	excellogger.WiteCellln("Please Do not edit this page!")
	excellogger.WiteCellln("This page is created by the User Report script: " + "https://github.com/perolo/perolo/confluence-scripts" + "/" + "SpacePermissionsReport")
	t := time.Now()
	excellogger.WiteCellln("Created by: " + cfg.ConfUser + " : " + t.Format(time.RFC3339))
	excellogger.WiteCellln("")

	var config = client.ConfluenceConfig{}
	config.Username = cfg.ConfUser
	config.Password = cfg.ConfPass
	config.UseToken = cfg.UseToken
	config.URL = cfg.ConfHost
	config.Debug = false

	theClient := client.Client(&config)
	types := theClient.GetPermissionTypes()
	excellogger.SetCellFontHeader2()
	excellogger.WiteCellln("Users and Permissions for group: " + cfg.Group)
	excellogger.NextLine()
	excellogger.AutoFilterStart()
	excellogger.SetTableHeader()
	excellogger.WiteCell("Space Name")
	excellogger.SetTableHeader()
	excellogger.NextCol()
	excellogger.SetTableHeader()
	excellogger.WiteCell("Space Key")
	//excellogger.SetCellStyleRotate()
	excellogger.NextCol()
	excellogger.SetTableHeader()
	excellogger.WiteCell("Type")
	//excellogger.SetCellStyleRotate()
	excellogger.NextCol()
	excellogger.SetTableHeader()
	excellogger.WiteCell("Name")
	//excellogger.SetCellStyleRotate()
	excellogger.NextCol()
	excellogger.SetCellStyleRotateN(len(*types))
	excellogger.WriteColumnsln(*types)

	cont := true
	start := 0
	max := 20
	for cont {
		groupMembers, err, _ := theClient.GetGroupMembers(cfg.Group, &client.GetGroupMembersOptions{StartAt: start, MaxResults: max, ShowExtendedDetails: true})
		if err != nil {
			panic(err)
		}
		for _, member := range groupMembers.Users {
			spcont := true
			spstart := 0
			spmax := 20
			for spcont {
				spaces, _ := theClient.GetAllSpacesWithPermissions(member.Name, &client.GetGroupMembersOptions{StartAt: start, MaxResults: max, ShowExtendedDetails: true})
				for _, space := range spaces.Spaces {
					fmt.Printf("User: %s Space: %s Permissions: %s \n", member.Name, space.Name, space.Permissions)
					excellogger.ResetCol()
					excellogger.WiteCellnc(space.Name)
					//excellogger.WiteCellnc(space.Key)
					excellogger.WiteCellHyperLinknc(space.Key, cfg.ConfHost+"/spaces/spacepermissions.action?key="+space.Key)
					excellogger.WiteCellnc("User")
					excellogger.WiteCellnc(member.Name)
					for _, atype := range *types {
						if Contains(space.Permissions, atype) {
							excellogger.SetCellStyleCenter()
							excellogger.WiteCellnc("X")
						} else {
							excellogger.SetCellStyleCenter()
							excellogger.WiteCellnc("-")
						}
					}
					excellogger.NextLine()
				}
				if len(spaces.Spaces) != spmax {
					spcont = false
				} else {
					//					spcont = false // testing
					spstart = spstart + spmax
				}
			}
		}
		if len(groupMembers.Users) != max {
			cont = false
		} else {
			//			cont = false // testing
			start = start + max
		}
	}

	excellogger.SetAutoColWidth()
	excellogger.AutoFilterEnd()

	excellogger.SetColWidth("A", "A", 40)
	// Save xlsx file by the given path.
	excellogger.SaveAs(cfg.File)
	if cfg.Report {
		var config = client.ConfluenceConfig{}
		var copt client.OperationOptions
		config.Username = cfg.ConfUser
		config.Password = cfg.ConfPass
		config.UseToken = cfg.UseToken
		config.URL = cfg.ConfHost
		config.Debug = false
		confluenceClient := client.Client(&config)

		copt.Title = "Space Group Permissions Reports"
		copt.SpaceKey = "AAAD"
		_, name := filepath.Split(cfg.File)
		utilities.CheckPageExists(copt, confluenceClient)
		err := utilities.AddAttachmentAndUpload(confluenceClient, copt, name, cfg.File, "Created by Space Group Permissions Report")
		if err != nil {
			panic(err)
		}

	}
}
