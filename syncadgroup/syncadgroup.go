package syncadgroup

import (
	"flag"
	"fmt"
	"github.com/magiconair/properties"
	"github.com/perolo/ad-utils"
	"github.com/perolo/confluence-prop/client"
	"github.com/perolo/confluence-scripts/utilities"
	excelutils "github.com/perolo/excel-utils"
	"log"
	"path/filepath"
	"time"
)

// or through Decode
type Config struct {
	ConfHost        string `properties:"confhost"`
	User            string `properties:"user"`
	Pass            string `properties:"password"`
	Simple          bool   `properties:"simple"`
	AddOperation    bool   `properties:"add"`
	RemoveOperation bool   `properties:"remove"`
	Report          bool   `properties:"report"`
	Limited         bool   `properties:"limited"`
	AdGroup         string `properties:"adgroup"`
	Localgroup      string `properties:"localgroup"`
	File            string `properties:"file"`
	ConfUpload      bool   `properties:"confupload"`
	ConfPage        string `properties:"confluencepage"`
	ConfSpace       string `properties:"confluencespace"`
	ConfAttName     string `properties:"conlfuenceattachment"`
	Bindusername    string `properties:"bindusername"`
	Bindpassword    string `properties:"bindpassword"`
}

func initReport(cfg Config) {
	if cfg.Report {
		excelutils.NewFile()
		excelutils.SetCellFontHeader()
		excelutils.WiteCellln("Introduction")
		excelutils.WiteCellln("Please Do not edit this page!")
		excelutils.WiteCellln("This page is created by the projectreport script: github.com\\perolo\\confluence-scripts\\SyncADGroup")
		t := time.Now()
		excelutils.WiteCellln("Created by: " + cfg.User + " : " + t.Format(time.RFC3339))
		excelutils.WiteCellln("")
		excelutils.WiteCellln("The Report Function shows:")
		excelutils.WiteCellln("   AdNames - Name and user found in AD Group")
		excelutils.WiteCellln("   JIRA Users - Name and user found in JIRA Group")
		excelutils.WiteCellln("   Not in AD - Users in the Local Group not found in the AD")
		excelutils.WiteCellln("   Not in JIRA - Users in the AD not found in the JIRA Group")
		excelutils.WiteCellln("   AD Errors - Internal error when searching for user in AD")
		excelutils.WiteCellln("")
		excelutils.SetCellFontHeader2()
		excelutils.WiteCellln("Group Mapping")
		if cfg.Simple {
			excelutils.WriteColumnsHeaderln([]string{"AD Group", "Local group"})
			excelutils.WriteColumnsln([]string{cfg.AdGroup, cfg.Localgroup})
		} else {
			excelutils.WriteColumnsHeaderln([]string{"AD Group", "Local group"})
			for _, syn := range GroupSyncs {
				excelutils.WriteColumnsln([]string{syn.AdGroup, syn.LocalGroup})
			}
		}
		excelutils.WiteCellln("")
		excelutils.SetCellFontHeader2()
		excelutils.WiteCellln("Report")
		excelutils.AutoFilterStart()
		var headers = []string{"Report Function", "AD group", "Local Group", "Name", "Uname", "Mail", "Error", "DN"}
		excelutils.WriteColumnsHeaderln(headers)
	}
}

func endReport(cfg Config) {
	if cfg.Report {
		file := fmt.Sprintf(cfg.File, "-Confluence")
		excelutils.SetColWidth("A", "A", 60)
		excelutils.AutoFilterEnd()
		excelutils.SaveAs(file)
		if cfg.ConfUpload {
			var config = client.ConfluenceConfig{}
			var copt client.OperationOptions
			config.Username = cfg.User
			config.Password = cfg.Pass
			config.URL = cfg.ConfHost
			config.Debug = false
			confluenceClient := client.Client(&config)
			// Intentional override
			copt.Title = "Using AD groups for JIRA/Confluence"
			copt.SpaceKey = "STPIM"
			_, name := filepath.Split(file)
			cfg.ConfAttName = name
			utilities.AddAttachmentAndUpload(confluenceClient, copt, name, file, "Created by Sync AD group")
		}
	}
}

func ConfluenceSyncAdGroup(propPtr string) {

	//	propPtr := flag.String("prop", "confluence.properties", "a string")
	flag.Parse()
	p := properties.MustLoadFile(propPtr, properties.ISO_8859_1)
	var cfg Config
	if err := p.Decode(&cfg); err != nil {
		log.Fatal(err)
	}
	toolClient := toollogin(cfg)
	initReport(cfg)
	ad_utils.InitAD(cfg.Bindusername, cfg.Bindpassword)
	if cfg.Simple {
		SyncGroupInTool(cfg, toolClient)
	} else {
		for _, syn := range GroupSyncs {
			cfg.AdGroup = syn.AdGroup
			cfg.Localgroup = syn.LocalGroup
			SyncGroupInTool(cfg, toolClient)
		}
	}
	endReport(cfg)
	ad_utils.CloseAD()
}

func toollogin(cfg Config) *client.ConfluenceClient {
	var config = client.ConfluenceConfig{}
	config.Username = cfg.User
	config.Password = cfg.Pass
	config.URL = cfg.ConfHost
	config.Debug = false
	return client.Client(&config)
}

func SyncGroupInTool(cfg Config, client *client.ConfluenceClient) {
	var toolGroupMemberNames map[string]ad_utils.ADUser
	fmt.Printf("\n")
	fmt.Printf("SyncGroup AdGroup: %s LocalGroup: %s \n", cfg.AdGroup, cfg.Localgroup)
	fmt.Printf("\n")
	var adUnames []ad_utils.ADUser
	if cfg.AdGroup != "" {
		adUnames, _ = ad_utils.GetUnamesInGroup(cfg.AdGroup)
		fmt.Printf("adUnames(%v): %s \n", len(adUnames), adUnames)
	}
	if cfg.Report {
		if !cfg.Limited {
			for _, adu := range adUnames {
				var row = []string{"AD Names", cfg.AdGroup, cfg.Localgroup, adu.Name, adu.Uname, adu.Mail, adu.Err, adu.DN}
				excelutils.WriteColumnsln(row)
			}
		}
	}
	if cfg.Localgroup != "" {
		toolGroupMemberNames = getUnamesInToolGroup(client, cfg.Localgroup)
		if cfg.Report {
			if !cfg.Limited {
				for _, tgm := range toolGroupMemberNames {
					var row = []string{"Confluence Users", cfg.AdGroup, cfg.Localgroup, tgm.Name, tgm.Uname, tgm.Mail, tgm.Err, tgm.DN}
					excelutils.WriteColumnsln(row)
				}
			}
		}
	}
	if cfg.Localgroup != "" && cfg.AdGroup != "" {
		notInTool := ad_utils.Difference(adUnames, toolGroupMemberNames)
		fmt.Printf("Not In Tool(%v): %s \n", len(notInTool), notInTool)
		if cfg.Report {
			for _, nji := range notInTool {
				var row = []string{"AD group users not found in Tool user group", cfg.AdGroup, cfg.Localgroup, nji.Name, nji.Uname, nji.Mail, nji.Err, nji.DN}
				excelutils.WriteColumnsln(row)
			}
		}
		notInAD := ad_utils.Difference2(toolGroupMemberNames, adUnames)
		fmt.Printf("notInAD(%v): %s \n", len(notInAD), notInAD)
		if cfg.Report {
			for _, nad := range notInAD {
				var row = []string{"Tool user group member not found in AD", cfg.AdGroup, cfg.Localgroup, nad.Name, nad.Uname, nad.Mail, nad.Err, nad.DN}
				excelutils.WriteColumnsln(row)
			}
		}
		if cfg.AddOperation {
			for _, notin := range notInTool {
				if notin.Err == "" {
					addUser := client.AddGroupMembers(cfg.Localgroup, []string{notin.Uname})
					fmt.Printf("Group: %s status: %s \n", cfg.Localgroup, addUser.Status)
					fmt.Printf("Message: %s \n", addUser.Message)
					fmt.Printf("Users Added: %s \n", addUser.UsersAdded)
					fmt.Printf("Users Skipped: %s \n", addUser.UsersSkipped)
				} else {
					fmt.Printf("Ad Problems skipping add: %s \n", notin.Uname)
				}
			}
		}

		if cfg.RemoveOperation {
			for _, notin := range notInAD {
				if notin.Err == "" {

					removeUser := client.RemoveGroupMembers(cfg.Localgroup, []string{notin.Uname})
					fmt.Printf("Remove user. Group: %s status: %s \n", cfg.Localgroup, removeUser.Status)
					fmt.Printf("Message: %s \n", removeUser.Message)
					fmt.Printf("Users Removed: %s \n", removeUser.UsersRemoved)
					fmt.Printf("Users Skipped: %s \n", removeUser.UsersSkipped)
				} else {
					fmt.Printf("Ad Problems skipping remove: %s \n", notin.Uname)
				}
				//fmt.Printf("Not Yet Implemented\n")
			}
		}

	}
}
func getUnamesInToolGroup(theClient *client.ConfluenceClient, localgroup string) map[string]ad_utils.ADUser {
	groupMemberNames := make(map[string]ad_utils.ADUser)
	cont := true
	start := 0
	max := 50
	for cont {
		groupMembers, err := theClient.GetGroupMembers(localgroup, &client.GetGroupMembersOptions{StartAt: start, MaxResults: max, ShowBasicDetails: true})
		if err != nil {
			panic(err)
		}
		for _, member := range groupMembers.Users {
			if _, ok := groupMemberNames[member.Name]; !ok {
				var newUser ad_utils.ADUser
				newUser.Uname = member.Name
				newUser.Name = member.FullName
				newUser.Mail = member.Email
				groupMemberNames[member.Name] = newUser
			}
		}
		if len(groupMembers.Users) != max {
			cont = false
		} else {
			start = start + max
		}
	}
	return groupMemberNames
}
