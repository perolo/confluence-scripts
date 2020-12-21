package spacepermissionsreport

import (
	"fmt"
	"github.com/perolo/confluence-prop/client"
	"github.com/perolo/excel-utils"
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

type SpacePermissionsReportConfig struct {
	ConfHost      string `properties:"confhost"`
	User          string `properties:"user"`
	Pass          string `properties:"password"`
	Groups        bool   `properties:"groups"`
	Users         bool   `properties:"users"`
	SpaceCategory string `properties:"spacecategory"`
	File          string `properties:"file"`
}

func SpacePermissionsReport(cfg SpacePermissionsReportConfig) {

	excelutils.NewFile()

	excelutils.SetCellFontHeader()
	excelutils.WiteCellln("Introduction")
	excelutils.WiteCellln("Please Do not edit this page!")
	excelutils.WiteCellln("This page is created by the User Report script: " + "https://git.aa.st/perolo/confluence-scripts" + "/" + "SpacePermissionsReport")
	t := time.Now()
	excelutils.WiteCellln("Created by: " + cfg.User + " : " + t.Format(time.RFC3339))
	excelutils.WiteCellln("")

	var config = client.ConfluenceConfig{}
	config.Username = cfg.User
	config.Password = cfg.Pass
	config.URL = cfg.ConfHost
	config.Debug = false

	theClient := client.Client(&config)
	types := theClient.GetPermissionTypes()
	excelutils.SetCellFontHeader2()
	excelutils.WiteCellln("Users and Permissions")
	excelutils.NextLine()
	excelutils.AutoFilterStart()
	excelutils.SetTableHeader()
	excelutils.WiteCell("Space")
	excelutils.SetTableHeader()
	excelutils.WiteCell("Key")
	excelutils.SetCellStyleRotate()
	excelutils.NextCol()
	excelutils.SetTableHeader()
	excelutils.WiteCell("Type")
	excelutils.SetCellStyleRotate()
	excelutils.NextCol()
	excelutils.SetTableHeader()
	excelutils.WiteCell("Group")
	excelutils.SetCellStyleRotate()
	excelutils.NextCol()
	excelutils.SetTableHeader()
	excelutils.WiteCell("Name")
	excelutils.SetCellStyleRotate()
	excelutils.NextCol()
	excelutils.SetCellStyleRotateN(len(*types))
	excelutils.WriteColumnsln([]string (*types))
	noSpaces := 0
	spstart := 0
	spincrease := 10
	spcont := true
	var spaces *client.ConfluenceSpaceResult
	for spcont {
		spopt := client.SpaceOptions{Start: spstart, Limit: spincrease, Label: cfg.SpaceCategory}
		spaces = theClient.GetSpaces(&spopt)
		opt := client.PaginationOptions{}
		for _, space := range spaces.Results {
			if space.Type == "global" {
				noSpaces++
				fmt.Printf("Space: %s \n", space.Name)
				if cfg.Groups {
					start := 0
					cont := true
					increase := 10
					for cont {
						opt.StartAt = start
						opt.MaxResults = increase
						groups := theClient.GetAllGroupsWithAnyPermission(space.Key, &opt)
						excelutils.NextCol()
						for _, group := range groups.Groups {
							excelutils.ResetCol()
							excelutils.WiteCellnc(space.Name)
							excelutils.WiteCellnc(space.Key)
							excelutils.WiteCellnc("Group")
							permissions := theClient.GetGroupPermissionsForSpace(space.Key, group)
							excelutils.WiteCellnc(group)
							for _, atype := range *types {
								if Contains(permissions.Permissions, atype) {
									excelutils.SetCellStyleCenter()
									excelutils.WiteCellnc("X")
								} else {
									excelutils.SetCellStyleCenter()
									excelutils.WiteCellnc("-")
								}
							}
							excelutils.NextLine()
						}
						start = start + increase
						if groups.Total < increase {
							cont = false
						}
					}
				}
				if cfg.Users {
					start := 0
					cont := true
					increase := 10
					for cont {
						opt.StartAt = start
						opt.MaxResults = increase
						users := theClient.GetAllUsersWithAnyPermission(space.Key, &opt)
						excelutils.NextCol()
						for _, user := range users.Users {
							excelutils.ResetCol()
							excelutils.WiteCellnc(space.Name)
							excelutils.WiteCellnc(space.Key)
							excelutils.WiteCellnc("User")
							permissions := theClient.GetUserPermissionsForSpace(space.Key, user)
							excelutils.WiteCellnc(user)
							for _, atype := range *types {
								if Contains(permissions.Permissions, atype) {
									excelutils.SetCellStyleCenter()
									excelutils.WiteCellnc("X")
								} else {
									excelutils.SetCellStyleCenter()
									excelutils.WiteCellnc("-")
								}
							}
							excelutils.NextLine()
						}
						start = start + increase
						if users.Total < increase {
							cont = false
						}
					}
				}
			}
		}
		spstart = spstart + spincrease
		if spaces.Size < spincrease {
			spcont = false
		}

	}
	excelutils.AutoFilterEnd()

	excelutils.SetColWidth("A","A",40)
	excelutils.SetColWidth("B","D",30)
	excelutils.SetColWidth("E","R",5)
	// Save xlsx file by the given path.
	excelutils.SaveAs(cfg.File)
}
