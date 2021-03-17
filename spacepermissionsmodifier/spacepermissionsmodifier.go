package spacepermissionsmodifier

import (
	"fmt"
	"github.com/manifoldco/promptui"
	"github.com/perolo/confluence-prop/client"
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

type SpacePermissionsModConfig struct {
	ConfHost string `properties:"confhost"`
	User     string `properties:"user"`
	Pass     string `properties:"password"`
	Groups   bool   `properties:"groups"`
	Users    bool   `properties:"users"`
	//Space         string `properties:"space"`
	SpaceCategory string `properties:"spacecategory"`
	//File          string `properties:"file"`
}

func Modify(cfg SpacePermissionsModConfig) {

	fmt.Printf("Introduction\n")
	fmt.Printf("This script is maintained in: " + "https://git.aa.st/perolo/confluence-scripts" + "/" + "spacepermissionsmodifier \n")
	fmt.Printf("\n")
	fmt.Printf("Modify Permissions in SpaceCategory: %s\n", cfg.SpaceCategory)

	var config = client.ConfluenceConfig{}
	config.Username = cfg.User
	config.Password = cfg.Pass
	config.URL = cfg.ConfHost
	config.Debug = false

	promptMode := promptui.Select{
		Label: "Add or Remove Permission",
		Items: []string{"Display", "Add", "Remove"},
	}

	_, mode, err := promptMode.Run()

	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return
	}

	fmt.Printf("Selected mode: %q\n", mode)

	theClient := client.Client(&config)
	types := theClient.GetPermissionTypes()

	prompt := promptui.Select{
		Label: "Select Permission",
		Items: *types,
	}

	_, selectedperm, err := prompt.Run()

	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return
	}

	fmt.Printf("Selected permission: %q\n", selectedperm)

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
						for _, group := range groups.Groups {
							permissions := theClient.GetGroupPermissionsForSpace(space.Key, group)
							ModifyPerm(permissions, selectedperm, mode, space, group)
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
						for _, user := range users.Users {
							permissions := theClient.GetUserPermissionsForSpace(space.Key, user)
							ModifyPerm(permissions, selectedperm, mode, space, user)
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
}

func ModifyPerm(permissions *client.GetPermissionsForSpaceType, selectedperm string, mode string, space client.SpaceType, user string) {
	if Contains(permissions.Permissions, selectedperm) {
		if mode == "Display" {
			fmt.Printf("Display OK Space: %s Permission: %s for user: %s\n", space.Name, selectedperm, user)
		} else if mode == "Remove" {
			fmt.Printf("Remove Permission: %s in Space: %s for user: %s\n", selectedperm, space.Name, user)
			confirm := promptui.Prompt{
				Label:     "Delete Permission",
				IsConfirm: true,
			}
			ok, _ := confirm.Run()
			/*
				if err != nil {
					fmt.Printf("Prompt failed %v\n", err)
					panic("Que")
				}*/
			if ok == "y" {
				fmt.Printf("Removed permission: %s from Space: %s for user: %s\n", selectedperm, space.Name, user)
			} else if ok == "N" {
				fmt.Printf("Skipped Removing permission: %s from Space: %s for user: %s \n", selectedperm, space.Name, user)
			} else {
				panic("Que")
			}
		} else if mode == "Add" {
			fmt.Printf("Skipped Add permission: %s to Space: %s for user: %s\n", selectedperm, space.Name, user)
		} else {
			panic("Que")
		}
	} else {
		//								fmt.Printf("X: %s \n", space.Name)
		if mode == "Display" {
			fmt.Printf("Display NOK Space: %s Permission: %s for user: %s\n", space.Name, selectedperm, user)
		} else if mode == "Add" {
			fmt.Printf("Add permission: %s to Space: %s for user: %s\n", selectedperm, space.Name, user)
			confirm := promptui.Prompt{
				Label:     "Add Permission",
				IsConfirm: true,
			}
			ok, _ := confirm.Run()
			/*
				if err != nil {
					fmt.Printf("Prompt failed %v\n", err)
					panic("Que")
				}*/
			if ok == "y" {
				fmt.Printf("Added permission: %s from Space: %s for user: %s\n", selectedperm, space.Name, user)
			} else if ok == "N" {
				fmt.Printf("Skipped Add permission: %s from Space: %s for user: %s \n", selectedperm, space.Name, user)
			} else {
				panic("Que")
			}

		} else if mode == "Remove" {
			fmt.Printf("Skipped Remove permission: %s to Space: %s for user: %s\n", selectedperm, space.Name, user)
		} else {
			panic("Que")
		}
	}
}
