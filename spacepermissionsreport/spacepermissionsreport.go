package spacepermissionsreport

import (
	"flag"
	"fmt"
	"log"
	"time"

	goconfluence "github.com/perolo/confluence-go-api"

	"github.com/magiconair/properties"
	"github.com/perolo/excellogger"
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
	ConfToken     string `properties:"conftoken"`
	UseToken      bool   `properties:"usetoken"`
	Groups        bool   `properties:"groups"`
	Users         bool   `properties:"users"`
	SpaceCategory string `properties:"spacecategory"`
	File          string `properties:"file"`
	Simple        bool   `properties:"simple"`
	Report        bool   `properties:"report"`
	//	Reset         bool   `properties:"reset"`
	Space         string `properties:"space"`
	AncestorTitle string `properties:"ancestortitle"`
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
			//			if schedulerutil.CheckScheduleDetail("SpacePermissionsReport-"+category, 7*time.Hour*24, cfg.Reset, schedulerutil.DummyFunc, "jiracategory.properties") {
			cfg.SpaceCategory = category
			cfg.File = fmt.Sprintf(reportBase, "-"+category)
			fmt.Printf("Category: %s \n", category)
			CreateSpacePermissionsReport(cfg)
			//			}
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

	/*
		var config = client.ConfluenceConfig{}
		config.Username = cfg.ConfUser
		config.Password = cfg.ConfPass
		config.UseToken = cfg.UseToken
		config.URL = cfg.ConfHost
		config.Debug = false

		theClient := client.Client(&config)
	*/

	var confluence *goconfluence.API
	var err error
	if cfg.UseToken {
		confluence, err = goconfluence.NewAPI(cfg.ConfHost, "", cfg.ConfToken)
	} else {
		confluence, err = goconfluence.NewAPI(cfg.ConfHost, cfg.ConfUser, cfg.ConfPass)
	}
	if err != nil {
		log.Fatal(err)
	}

	types, err2 := confluence.GetPermissionTypes()
	if err2 != nil {
		log.Fatal(err2)
	}

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
	//var spaces *confluence.Search
	for spcont {
		spopt := goconfluence.AllSpacesQuery{Start: spstart, Limit: spincrease, Label: cfg.SpaceCategory, Type: "global", Status: "current"}
		spaces, _ := confluence.GetAllSpaces(spopt)
		opt := goconfluence.PaginationOptions{}
		for _, space := range spaces.Results {
			if space.Type == "global" {
				noSpaces++
				fmt.Printf("Space: %s \n", space.Name)
				SpaceOwner := ""
				/*
					if cfg.SpaceCategory == "demo" {
						found, page := searchutils.SearchSpacePage(confluence, space.Key)
						if found {
							ownerFound, ownerName := searchutils.GetOwner(confluence, page)
							if ownerFound {
								SpaceOwner = ownerName
							}
						}
					}
				*/
				//htmlutils.WriteWrapLink(f, cfg.ConfHost+"/display/"+spaceKey+"/?pageId="+page.ID, "Space Description")

				if cfg.Groups {
					start := 0
					cont := true
					increase := 50
					for cont {
						opt.StartAt = start
						opt.MaxResults = increase
						groups, err := confluence.GetAllGroupsWithAnyPermission(space.Key, &opt)
						if err != nil {
							log.Fatal(err)
						}

						excellogger.NextCol()
						for _, group := range groups.Groups {
							excellogger.ResetCol()
							excellogger.WiteCellnc(space.Name)
							excellogger.WiteCellnc(SpaceOwner)
							//excellogger.WiteCellnc(space.Key)
							excellogger.WiteCellHyperLinknc(space.Key, cfg.ConfHost+"/spaces/spacepermissions.action?key="+space.Key)
							excellogger.WiteCellnc("Group")
							permissions, err := confluence.GetGroupPermissionsForSpace(space.Key, group)
							if err != nil {
								log.Fatal(err)
							}
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
					var start, increase int64
					start = 0
					cont := true
					increase = 50
					for cont {
						opt.StartAt = int(start)
						opt.MaxResults = int(increase)

						users, err := confluence.GetAllUsersWithAnyPermission(space.Key, &opt)
						if err != nil {
							log.Fatal(err)
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
							permissions, err := confluence.GetUserPermissionsForSpace(space.Key, user)
							if err != nil {
								log.Fatal(err)
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
		if int(spaces.Size) < spincrease {
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
		/*
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
		*/
		//		_, name := filepath.Split(cfg.File)
		res, err := confluence.GetPageId(cfg.Space, cfg.AncestorTitle)
		if err == nil {
			if res.Size == 1 {
				err := confluence.UppdateAttachment(cfg.Space, cfg.AncestorTitle, cfg.File)
				if err != nil {
					panic(err)
				}

			}

		}

		//) .AddAttachmentAndUpload(cfg.Space, name, cfg.File)
		//AddAttachmentAndUpload(confluenceClient, copt, name, cfg.File, "Created by Space Permissions Report")

	}
}
