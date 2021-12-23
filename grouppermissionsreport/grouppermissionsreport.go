package grouppermissionsreport

import (
	"flag"
	"fmt"
	"github.com/magiconair/properties"
	"github.com/perolo/confluence-prop/client"
	"github.com/perolo/confluence-scripts/schedulerutil"
	"github.com/perolo/confluence-scripts/utilities"
	"github.com/perolo/excel-utils"
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
	} else {
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

	excelutils.NewFile()

	excelutils.SetCellFontHeader()
	excelutils.WiteCellln("Introduction")
	excelutils.WiteCellln("Please Do not edit this page!")
	excelutils.WiteCellln("This page is created by the User Report script: " + "https://git.aa.st/perolo/confluence-scripts" + "/" + "SpacePermissionsReport")
	t := time.Now()
	excelutils.WiteCellln("Created by: " + cfg.ConfUser + " : " + t.Format(time.RFC3339))
	excelutils.WiteCellln("")

	var config = client.ConfluenceConfig{}
	config.Username = cfg.ConfUser
	config.Password = cfg.ConfPass
	config.UseToken = cfg.UseToken
	config.URL = cfg.ConfHost
	config.Debug = false

	theClient := client.Client(&config)
	types := theClient.GetPermissionTypes()
	excelutils.SetCellFontHeader2()
	excelutils.WiteCellln("Users and Permissions for group: " + cfg.Group)
	excelutils.NextLine()
	excelutils.AutoFilterStart()
	excelutils.SetTableHeader()
	excelutils.WiteCell("Space Name")
	excelutils.SetTableHeader()
	excelutils.NextCol()
	excelutils.SetTableHeader()
	excelutils.WiteCell("Space Key")
	//excelutils.SetCellStyleRotate()
	excelutils.NextCol()
	excelutils.SetTableHeader()
	excelutils.WiteCell("Type")
	//excelutils.SetCellStyleRotate()
	excelutils.NextCol()
	excelutils.SetTableHeader()
	excelutils.WiteCell("Name")
	//excelutils.SetCellStyleRotate()
	excelutils.NextCol()
	excelutils.SetCellStyleRotateN(len(*types))
	excelutils.WriteColumnsln(*types)

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
					excelutils.ResetCol()
					excelutils.WiteCellnc(space.Name)
					//excelutils.WiteCellnc(space.Key)
					excelutils.WiteCellHyperLinknc(space.Key, cfg.ConfHost+"/spaces/spacepermissions.action?key="+space.Key)
					excelutils.WiteCellnc("User")
					excelutils.WiteCellnc(member.Name)
					for _, atype := range *types {
						if Contains(space.Permissions, atype) {
							excelutils.SetCellStyleCenter()
							excelutils.WiteCellnc("X")
						} else {
							excelutils.SetCellStyleCenter()
							excelutils.WiteCellnc("-")
						}
					}
					excelutils.NextLine()
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

	excelutils.SetAutoColWidth()
	excelutils.AutoFilterEnd()

	excelutils.SetColWidth("A", "A", 40)
	// Save xlsx file by the given path.
	excelutils.SaveAs(cfg.File)
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
