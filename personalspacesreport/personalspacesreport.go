package personalspacesreport

import (
	"flag"
	"fmt"
	"github.com/magiconair/properties"
	"github.com/perolo/ad-utils"
	"github.com/perolo/confluence-prop/client"
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
	ConfHost     string `properties:"confhost"`
	ConfUser     string `properties:"confuser"`
	ConfPass     string `properties:"confpass"`
	ConfToken    string `properties:"conftoken"`
	UseToken     bool   `properties:"usetoken"`
	File         string `properties:"file"`
	Simple       bool   `properties:"simple"`
	Report       bool   `properties:"report"`
	Bindusername string `properties:"bindusername"`
	Bindpassword string `properties:"bindpassword"`
	BaseDN       string `properties:"basedn"`
}

func PersonalSpaceReport(propPtr string) {

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

	adutils.InitAD(cfg.Bindusername, cfg.Bindpassword)
	if cfg.Simple {
		cfg.File = fmt.Sprintf(cfg.File, "-"+"PersonalSpaces")
		CreatePersonalSpacesReport(cfg)
	} else {
		reportBase := cfg.File
		cfg.File = fmt.Sprintf(reportBase, "-"+"PersonalSpaces")
		CreatePersonalSpacesReport(cfg)
	}
	adutils.CloseAD()

}

func CreatePersonalSpacesReport(cfg ReportConfig) {

	excelutils.NewFile()

	excelutils.SetCellFontHeader()
	excelutils.WiteCellln("Introduction")
	excelutils.WiteCellln("Please Do not edit this page!")
	excelutils.WiteCellln("This page is created by the User Report script: " + "https://git.aa.st/perolo/confluence-scripts" + "/" + "PersonalSpacesReport")
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
	excelutils.SetCellFontHeader2()
	excelutils.WiteCellln("Users and Permissions")
	excelutils.NextLine()
	excelutils.AutoFilterStart()
	excelutils.SetTableHeader()
	excelutils.WiteCell("Space Name")
	excelutils.SetTableHeader()
	excelutils.NextCol()
	excelutils.SetTableHeader()
	excelutils.WiteCell("Space Key")
	excelutils.NextCol()
	excelutils.SetTableHeader()
	excelutils.WiteCell("Type")
	excelutils.NextCol()
	excelutils.SetTableHeader()
	excelutils.WiteCell("Name")
	excelutils.NextCol()

	excelutils.SetTableHeader()
	excelutils.WiteCell("DN")
	excelutils.NextCol()
	excelutils.SetTableHeader()
	excelutils.WiteCell("Mail")
	excelutils.NextCol()
	excelutils.SetTableHeader()
	excelutils.WiteCell("Comment")
	excelutils.NextCol()
	excelutils.NextLine()

	noSpaces := 0
	spstart := 0
	spincrease := 50
	spcont := true
	var spaces *client.ConfluenceSpaceResult
	for spcont {
		spopt := client.SpaceOptions{Start: spstart, Limit: spincrease, Type: "personal", Status: "current"}
		spaces, _ = theClient.GetSpaces(&spopt)
		opt := client.PaginationOptions{}
		for _, space := range spaces.Results {
			noSpaces++
			fmt.Printf("Space: %s \n", space.Name)
			start := 0
			cont := true
			increase := 50
			for cont {
				opt.StartAt = start
				opt.MaxResults = increase
				users, _ := theClient.GetAllUsersWithAnyPermission(space.Key, &opt)
				excelutils.NextCol()
				for _, user := range users.Users {
					permissions, _ := theClient.GetUserPermissionsForSpace(space.Key, user)
					if Contains(permissions.Permissions, "SETPAGEPERMISSIONS") {
						_, err := adutils.GetActiveUserDN(user, cfg.BaseDN)
						if err == nil {
							//excelutils.WiteCellnc(dn.DN)
							//excelutils.WiteCellnc(dn.Mail)
							//excelutils.WiteCellnc("")
						} else {
							excelutils.ResetCol()
							excelutils.WiteCellnc(space.Name)
							excelutils.WiteCellnc(space.Key)
							excelutils.WiteCellnc("User")
							excelutils.WiteCellnc(user)
							udn, err := adutils.GetAllUserDN(user, cfg.BaseDN)
							if err == nil {
								excelutils.WiteCellnc(udn.DN)
								excelutils.WiteCellnc(udn.Mail)
								excelutils.WiteCellnc("Deactivated!")
							} else {
								excelutils.WiteCellnc("")
								excelutils.WiteCellnc("")
								excelutils.WiteCellnc("Not Found!")
							}
							excelutils.NextLine()
						}
					}
				}
				if users.Total < start+increase {
					cont = false
				} else {
					start = start + increase
				}
			}
		}
		if spaces.Size < spincrease {
			spcont = false
		} else {
			spstart = spstart + spincrease
		}
	}
	excelutils.SetAutoColWidth()
	excelutils.AutoFilterEnd()

	excelutils.SetColWidth("A", "A", 60)
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

		copt.Title = "Personal Spaces Reports"
		copt.SpaceKey = "AAAD"
		_, name := filepath.Split(cfg.File)
		err := utilities.AddAttachmentAndUpload(confluenceClient, copt, name, cfg.File, "Created by PersonalSpacesReport")
		if err != nil {
			panic(err)
		}
	}
}
