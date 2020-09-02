package spacepermissionsreport

import (
	"github.com/perolo/confluence-prop/client"
	excelutilities "github.com/perolo/confluence-scripts/utilities"
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
	ConfHost string `properties:"confhost"`
	User     string `properties:"user"`
	Pass     string `properties:"password"`
	Groups   bool   `properties:"groups"`
	Users    bool   `properties:"users"`
	Space    string `properties:"space"`
	Label    string `properties:"label"`
	File     string `properties:"file"`
}

func SpacePermissionsReport(cfg SpacePermissionsReportConfig) {

	excelutilities.NewFile()

	excelutilities.SetCellFontHeader()
	excelutilities.WiteCellln("Introduction")
	excelutilities.WiteCellln("Please Do not edit this page!")
	excelutilities.WiteCellln("This page is created by the User Report script: " + "https://git.aa.st/perolo/confluence-utils")
	t := time.Now()
	excelutilities.WiteCellln("Created by: " + cfg.User + " : " + t.Format(time.RFC3339))
	excelutilities.WiteCellln("")

	var config = client.ConfluenceConfig{}
	config.Username = cfg.User
	config.Password = cfg.Pass
	config.URL = cfg.ConfHost
	config.Debug = false

	theClient := client.Client(&config)
	types := theClient.GetPermissionTypes()
	excelutilities.WiteCellln("Users and Permissions")
	excelutilities.NextLine()
	excelutilities.WiteCell("Space")
	excelutilities.SetCellStyleRotate()
	excelutilities.NextCol()
	excelutilities.WiteCell("Type")
	excelutilities.SetCellStyleRotate()
	excelutilities.NextCol()
	excelutilities.WiteCell("Group")
	excelutilities.SetCellStyleRotate()
	excelutilities.NextCol()
	excelutilities.SetCellStyleRotateN(len(*types))
	excelutilities.WriteColumnsln([]string (*types))
	noSpaces := 0
	spstart := 0
	spincrease := 10
	spcont := true
	var spaces *client.ConfluenceSpaceResult
	if cfg.Space != "" && cfg.Label == "" {
		//TBD
	} else {
		spopt := client.SpaceOptions{Start: spstart, Limit: spincrease, Label: cfg.Label}
		spaces = theClient.GetSpaces(&spopt)
	}
	opt := client.PaginationOptions{}
	for spcont {
		for _, space := range spaces.Results {
			if space.Type == "global" {
				noSpaces++
				if cfg.Groups {
					start := 0
					cont := true
					increase := 10
					for cont {
						opt.StartAt = start
						opt.MaxResults = increase
						groups := theClient.GetAllGroupsWithAnyPermission(space.Key, &opt)
						excelutilities.NextCol()
						for _, group := range groups.Groups {
							excelutilities.ResetCol()
							excelutilities.WiteCellnc(space.Key)
							excelutilities.WiteCellnc("Group")
							permissions := theClient.GetGroupPermissionsForSpace(space.Key, group)
							excelutilities.WiteCellnc(group)
							for _, atype := range *types {
								if Contains(permissions.Permissions, atype) {
									excelutilities.WiteCellnc("X")
								} else {
									excelutilities.WiteCellnc("-")
								}
							}
							excelutilities.NextLine()
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
						excelutilities.NextCol()
						for _, user := range users.Users {
							excelutilities.ResetCol()
							excelutilities.WiteCellnc(space.Key)
							excelutilities.WiteCellnc("User")
							permissions := theClient.GetUserPermissionsForSpace(space.Key, user)
							excelutilities.WiteCellnc(user)
							for _, atype := range *types {
								if Contains(permissions.Permissions, atype) {
									excelutilities.WiteCellnc("X")
								} else {
									excelutilities.WiteCellnc("-")
								}
							}
							excelutilities.NextLine()
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
	excelutilities.AutoFilter("A8")

	// Save xlsx file by the given path.
	excelutilities.SaveAs(cfg.File)
}