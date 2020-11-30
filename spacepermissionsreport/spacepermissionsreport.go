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

	excel_utils.NewFile()

	excel_utils.SetCellFontHeader()
	excel_utils.WiteCellln("Introduction")
	excel_utils.WiteCellln("Please Do not edit this page!")
	excel_utils.WiteCellln("This page is created by the User Report script: " + "https://git.aa.st/perolo/confluence-scripts" + "/" + "SpacePermissionsReport")
	t := time.Now()
	excel_utils.WiteCellln("Created by: " + cfg.User + " : " + t.Format(time.RFC3339))
	excel_utils.WiteCellln("")

	var config = client.ConfluenceConfig{}
	config.Username = cfg.User
	config.Password = cfg.Pass
	config.URL = cfg.ConfHost
	config.Debug = false

	theClient := client.Client(&config)
	types := theClient.GetPermissionTypes()
	excel_utils.SetCellFontHeader2()
	excel_utils.WiteCellln("Users and Permissions")
	excel_utils.NextLine()
	excel_utils.AutoFilterStart()
	excel_utils.SetTableHeader()
	excel_utils.WiteCell("Space")
	excel_utils.SetTableHeader()
	excel_utils.WiteCell("Key")
	excel_utils.SetCellStyleRotate()
	excel_utils.NextCol()
	excel_utils.SetTableHeader()
	excel_utils.WiteCell("Type")
	excel_utils.SetCellStyleRotate()
	excel_utils.NextCol()
	excel_utils.SetTableHeader()
	excel_utils.WiteCell("Group")
	excel_utils.SetCellStyleRotate()
	excel_utils.NextCol()
	excel_utils.SetTableHeader()
	excel_utils.WiteCell("Name")
	excel_utils.SetCellStyleRotate()
	excel_utils.NextCol()
	excel_utils.SetCellStyleRotateN(len(*types))
	excel_utils.WriteColumnsln([]string (*types))
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
						excel_utils.NextCol()
						for _, group := range groups.Groups {
							excel_utils.ResetCol()
							excel_utils.WiteCellnc(space.Name)
							excel_utils.WiteCellnc(space.Key)
							excel_utils.WiteCellnc("Group")
							permissions := theClient.GetGroupPermissionsForSpace(space.Key, group)
							excel_utils.WiteCellnc(group)
							for _, atype := range *types {
								if Contains(permissions.Permissions, atype) {
									excel_utils.SetCellStyleCenter()
									excel_utils.WiteCellnc("X")
								} else {
									excel_utils.SetCellStyleCenter()
									excel_utils.WiteCellnc("-")
								}
							}
							excel_utils.NextLine()
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
						excel_utils.NextCol()
						for _, user := range users.Users {
							excel_utils.ResetCol()
							excel_utils.WiteCellnc(space.Name)
							excel_utils.WiteCellnc(space.Key)
							excel_utils.WiteCellnc("User")
							permissions := theClient.GetUserPermissionsForSpace(space.Key, user)
							excel_utils.WiteCellnc(user)
							for _, atype := range *types {
								if Contains(permissions.Permissions, atype) {
									excel_utils.SetCellStyleCenter()
									excel_utils.WiteCellnc("X")
								} else {
									excel_utils.SetCellStyleCenter()
									excel_utils.WiteCellnc("-")
								}
							}
							excel_utils.NextLine()
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
	excel_utils.AutoFilterEnd()

	excel_utils.SetColWidth("A","A",40)
	excel_utils.SetColWidth("B","D",30)
	excel_utils.SetColWidth("E","R",5)
	// Save xlsx file by the given path.
	excel_utils.SaveAs(cfg.File)
}
