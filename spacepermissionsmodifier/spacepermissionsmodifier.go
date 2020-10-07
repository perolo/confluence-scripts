package spacepermissionsmodifier

import (
	"fmt"
	"github.com/perolo/confluence-prop/client"
	"github.com/manifoldco/promptui"
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
	ConfHost      string `properties:"confhost"`
	User          string `properties:"user"`
	Pass          string `properties:"password"`
	Groups        bool   `properties:"groups"`
	Users         bool   `properties:"users"`
	Space         string `properties:"space"`
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
		if cfg.Space != "" && cfg.SpaceCategory == "" {
			//TBD
		} else {
			spopt := client.SpaceOptions{Start: spstart, Limit: spincrease, Label: cfg.SpaceCategory}
			spaces = theClient.GetSpaces(&spopt)
		}
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
							if Contains(permissions.Permissions, selectedperm) {
								if mode == "Display" {
									fmt.Printf("Space: %s Permission\n", space.Name, selectedperm)
								} else if mode == "Remove" {
									fmt.Printf("Remove: %s \n", space.Name)
									confirm := promptui.Prompt{
										Label:     "Delete Permission",
										IsConfirm: true,
									}
									ok, err := confirm.Run()

									if err != nil {
										fmt.Printf("Prompt failed %v\n", err)
										return
									}

									fmt.Printf("You choose %q\n", ok)
								}
							} else {
								//								fmt.Printf("X: %s \n", space.Name)
								if mode == "Display" {
									fmt.Printf("Space: %s Permission\n", space.Name, selectedperm)
								} else if mode == "Add" {
									fmt.Printf("Add: %s %s\n", space.Name, selectedperm)
									confirm := promptui.Prompt{
										Label:     "Add Permission",
										IsConfirm: true,
									}
									ok, err := confirm.Run()

									if err != nil {
										fmt.Printf("Prompt failed %v\n", err)
										return
									}

									fmt.Printf("You choose %q\n", ok)
								}
							}

						}
						start = start + increase
						if groups.Total < increase {
							cont = false
						}
					}
				}
				/*
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
								for _, atype := range *types {
									if Contains(permissions.Permissions, atype) {
										fmt.Printf("X: %s \n", space.Name)
									} else {
										fmt.Printf("X: %s \n", space.Name)
									}
								}
							}
							start = start + increase
							if users.Total < increase {
								cont = false
							}
						}
					}*/
			}
		}
		spstart = spstart + spincrease
		if spaces.Size < spincrease {
			spcont = false
		}
	}
}
