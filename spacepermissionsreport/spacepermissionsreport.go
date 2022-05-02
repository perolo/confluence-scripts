package spacepermissionsreport

import (
	"flag"
	"fmt"
	"github.com/magiconair/properties"
	"github.com/perolo/confluence-client/client"
	"github.com/perolo/confluence-scripts/schedulerutil"
	"github.com/perolo/confluence-scripts/utilities"
	"github.com/perolo/confluence-scripts/utilities/searchutils"
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
	ConfHost      string `properties:"confhost"`
	ConfUser      string `properties:"confuser"`
	ConfPass      string `properties:"confpass"`
	UseToken      bool   `properties:"usetoken"`
	Groups        bool   `properties:"groups"`
	Users         bool   `properties:"users"`
	SpaceCategory string `properties:"spacecategory"`
	File          string `properties:"file"`
	Simple        bool   `properties:"simple"`
	Report        bool   `properties:"report"`
	Reset         bool   `properties:"reset"`
}

func SpacePermissionsReport(propPtr string) {

	flag.Parse()

	p := properties.MustLoadFile(propPtr, properties.ISO_8859_1)

	// or through Decode
	var cfg ReportConfig
	if err := p.Decode(&cfg); err != nil {
		log.Fatal(err)
	}
	if cfg.Simple {
		cfg.File = fmt.Sprintf(cfg.File, "-"+cfg.SpaceCategory)
		CreateSpacePermissionsReport(cfg)
	} else {
		reportBase := cfg.File
		for _, category := range Categories {
			if schedulerutil.CheckScheduleDetail("SpacePermissionsReport-"+category, 7*time.Hour*24, cfg.Reset, schedulerutil.DummyFunc, "jiracategory.properties") {
				cfg.SpaceCategory = category
				cfg.File = fmt.Sprintf(reportBase, "-"+category)
				fmt.Printf("Category: %s \n", category)
				CreateSpacePermissionsReport(cfg)
			}
		}
	}
}

func CreateSpacePermissionsReport(cfg ReportConfig) { //nolint:funlen

	excellogger.NewFile(nil)

	excellogger.SetCellFontHeader()
	excellogger.WiteCellln("Introduction")
	excellogger.WiteCellln("Please Do not edit this page!")
	excellogger.WiteCellln("This page is created by the User Report script: " + "https://github.com/perolo/confluence-scripts" + "/" + "SpacePermissionsReport")
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
	excellogger.WiteCellln("Users and Permissions")
	excellogger.NextLine()
	excellogger.AutoFilterStart()
	excellogger.SetTableHeader()
	excellogger.WiteCellnc("Space Name")
	excellogger.SetTableHeader()
	excellogger.WiteCellnc("Space Owner")
	excellogger.SetTableHeader()
	//	excellogger.NextCol()
	//	excellogger.SetTableHeader()
	excellogger.WiteCellnc("Space Key")
	//excellogger.SetCellStyleRotate()
	//	excellogger.NextCol()
	excellogger.SetTableHeader()
	excellogger.WiteCellnc("Type")
	//excellogger.SetCellStyleRotate()
	//	excellogger.NextCol()
	excellogger.SetTableHeader()
	excellogger.WiteCellnc("Name")
	//excellogger.SetCellStyleRotate()
	//	excellogger.NextCol()
	excellogger.SetCellStyleRotateN(len(*types))
	excellogger.WriteColumnsln(*types)
	noSpaces := 0
	spstart := 0
	spincrease := 50
	spcont := true
	var spaces *client.ConfluenceSpaceResult
	for spcont {
		spopt := client.SpaceOptions{Start: spstart, Limit: spincrease, Label: cfg.SpaceCategory, Type: "global", Status: "current"}
		spaces, _ = theClient.GetSpaces(&spopt)
		opt := client.PaginationOptions{}
		for _, space := range spaces.Results {
			if space.Type == "global" {
				noSpaces++
				fmt.Printf("Space: %s \n", space.Name)
				SpaceOwner := ""
				if cfg.SpaceCategory == "demo" {
					found, page := searchutils.SearchSpacePage(theClient, space.Key)
					if found {
						ownerFound, ownerName := searchutils.GetOwner(theClient, page)
						if ownerFound {
							SpaceOwner = ownerName
						}
					}
				}

				//htmlutils.WriteWrapLink(f, cfg.ConfHost+"/display/"+spaceKey+"/?pageId="+page.ID, "Space Description")

				if cfg.Groups {
					start := 0
					cont := true
					increase := 50
					for cont {
						opt.StartAt = start
						opt.MaxResults = increase
						groups := theClient.GetAllGroupsWithAnyPermission(space.Key, &opt)
						excellogger.NextCol()
						for _, group := range groups.Groups {
							excellogger.ResetCol()
							excellogger.WiteCellnc(space.Name)
							excellogger.WiteCellnc(SpaceOwner)
							//excellogger.WiteCellnc(space.Key)
							excellogger.WiteCellHyperLinknc(space.Key, cfg.ConfHost+"/spaces/spacepermissions.action?key="+space.Key) //https://confluence.assaabloy.net/spaces/spacepermissions.action?key=REL
							excellogger.WiteCellnc("Group")
							permissions := theClient.GetGroupPermissionsForSpace(space.Key, group)
							excellogger.WiteCellnc(group)
							for _, atype := range *types {
								if Contains(permissions.Permissions, atype) {
									excellogger.SetCellStyleCenter()
									excellogger.WiteCellnc("X")
								} else {
									excellogger.SetCellStyleCenter()
									excellogger.WiteCellnc("-")
								}
							}
							excellogger.NextLine()
						}
						if groups.Total < start+increase {
							cont = false
						} else {
							start = start + increase
						}
					}
				}
				if cfg.Users {
					start := 0
					cont := true
					increase := 50
					for cont {
						opt.StartAt = start
						opt.MaxResults = increase

						users, resp := theClient.GetAllUsersWithAnyPermission(space.Key, &opt)
						if resp.StatusCode < 200 || resp.StatusCode > 300 {
							// one restry...
							users, _ = theClient.GetAllUsersWithAnyPermission(space.Key, &opt)
						}
						//users, err := retry(3,200, theClient.GetAllUsersWithAnyPermission(space.Key, &opt))
						excellogger.NextCol()
						for _, user := range users.Users {
							excellogger.ResetCol()
							excellogger.WiteCellnc(space.Name)
							excellogger.WiteCellnc(SpaceOwner)
							//excellogger.WiteCellnc(space.Key)
							excellogger.WiteCellHyperLinknc(space.Key, cfg.ConfHost+"/spaces/spacepermissions.action?key="+space.Key)
							excellogger.WiteCellnc("User")
							permissions, resp := theClient.GetUserPermissionsForSpace(space.Key, user)
							if resp.StatusCode < 200 || resp.StatusCode > 300 {
								// one restry...
								permissions, _ = theClient.GetUserPermissionsForSpace(space.Key, user)
							}
							excellogger.WiteCellnc(user)
							for _, atype := range *types {
								if Contains(permissions.Permissions, atype) {
									excellogger.SetCellStyleCenter()
									excellogger.WiteCellnc("X")
								} else {
									excellogger.SetCellStyleCenter()
									excellogger.WiteCellnc("-")
								}
							}
							excellogger.NextLine()
						}
						if users.Total < start+increase {
							cont = false
						} else {
							start = start + increase
						}
					}
				}
			}
		}
		if spaces.Size < spincrease {
			spcont = false
		} else {
			spstart = spstart + spincrease
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

		copt.Title = "Space Permissions Reports"
		copt.SpaceKey = "AAAD"
		_, name := filepath.Split(cfg.File)
		err := utilities.AddAttachmentAndUpload(confluenceClient, copt, name, cfg.File, "Created by Space Permissions Report")
		if err != nil {
			panic(err)
		}

	}
}
