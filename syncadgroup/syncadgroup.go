package main

import (
	"flag"
	"fmt"
	"github.com/magiconair/properties"
	"github.com/perolo/ad-utils"
	"github.com/perolo/confluence-prop/client"
	excelutils "github.com/perolo/excel-utils"
	"log"
	"time"
)

// or through Decode
type Config struct {
	ConfHost     string `properties:"confhost"`
	User         string `properties:"user"`
	Pass         string `properties:"password"`
	Simple       bool   `properties:"simple"`
	AddOperation bool   `properties:"add"`
	Report       bool   `properties:"report"`
	Limited      bool   `properties:"limited"`
	AdGroup      string `properties:"adgroup"`
	Localgroup   string `properties:"localgroup"`
	File         string `properties:"file"`
	Bindusername string `properties:"bindusername"`
	Bindpassword string `properties:"bindpassword"`
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
		excelutils.AutoFilterStart()
		var headers = []string{"Report Function", "AD group", "Local Group", "Name", "Uname", "Error"}
		excelutils.WriteColumnsHeaderln(headers)

	}
}

func endReport(cfg Config) {
	if cfg.Report {

		file := fmt.Sprintf(cfg.File, "-Confluence")
		excelutils.AutoFilterEnd()
		excelutils.SaveAs(file)
	}
}

func main() {

	propPtr := flag.String("prop", "../confluence.properties", "a string")

	flag.Parse()

	p := properties.MustLoadFile(*propPtr, properties.ISO_8859_1)

	var cfg Config
	if err := p.Decode(&cfg); err != nil {
		log.Fatal(err)
	}

	initReport(cfg)

	var config = client.ConfluenceConfig{}
	config.Username = cfg.User
	config.Password = cfg.Pass
	config.URL = cfg.ConfHost
	config.Debug = false

	confClient := client.Client(&config)

	ad_utils.InitAD(cfg.Bindusername, cfg.Bindpassword)

	if cfg.Simple {
		SyncGroupInConfluence(cfg, confClient)
	} else {
		for _, syn := range GroupSyncs {
			//var adUnames []ad_utils.ADUser
			cfg.AdGroup = syn.AdGroup
			cfg.Localgroup = syn.LocalGroup
			SyncGroupInConfluence(cfg, confClient)
		}
	}
	endReport(cfg)
	ad_utils.CloseAD()
}

func SyncGroupInConfluence(cfg Config, confClient *client.ConfluenceClient) {
	var toolGroupMemberNames map[string]ad_utils.ADUser
	toolGroupMemberNames = make(map[string]ad_utils.ADUser)
	fmt.Printf("\n")
	fmt.Printf("SyncGroup AdGroup: %s LocalGroup: %s \n", cfg.AdGroup, cfg.Localgroup)
	fmt.Printf("\n")
	var adUnames, aderrs []ad_utils.ADUser
	if cfg.AdGroup != "" {
		adUnames, _, aderrs = ad_utils.GetUnamesInGroup(cfg.AdGroup)
		fmt.Printf("adUnames(%v): %s \n", len(adUnames), adUnames)
	}

	if cfg.Report {
		if !cfg.Limited {
			for _, adu := range adUnames {
				//			var row = []string{"AD group", "group", "fun", "Name", "Uname"}
				var row = []string{"AD Names", cfg.AdGroup, cfg.Localgroup, adu.Name, adu.Uname}
				excelutils.WriteColumnsln(row)
			}
		}
		for _, aderr := range aderrs {
			//			var row = []string{"AD group", "group", "fun", "Name", "Uname"}
			var row = []string{"AD Errors", cfg.AdGroup, cfg.Localgroup, aderr.Name, aderr.Uname, aderr.Err}
			excelutils.WriteColumnsln(row)
		}

	}
	if cfg.Localgroup != "" {
		getUnamesInConfluenceGroup(confClient, cfg.Localgroup, toolGroupMemberNames)
		if cfg.Report {
			if !cfg.Limited {
				for _, tgm := range toolGroupMemberNames {
					//			var row = []string{"AD group", "group", "fun", "Name", "Uname"}
					var row = []string{"JIRA Users", cfg.AdGroup, cfg.Localgroup, tgm.Name, tgm.Uname}
					excelutils.WriteColumnsln(row)
				}
			}
		}
	}

	if cfg.Localgroup != "" && cfg.AdGroup != "" {
		notInConfluence := ad_utils.Difference(adUnames, toolGroupMemberNames)
		fmt.Printf("notInConfluence(%v): %s \n", len(notInConfluence), notInConfluence)
		if cfg.Report {
			for _, nji := range notInConfluence {
				//			var row = []string{"AD group", "group", "fun", "Name", "Uname"}
				var row = []string{"Not in JIRA", cfg.AdGroup, cfg.Localgroup, nji.Name, nji.Uname}
				excelutils.WriteColumnsln(row)
			}
		}

		notInAD := ad_utils.Difference2(toolGroupMemberNames, adUnames)
		fmt.Printf("notInAD: %s \n", notInAD)
		if cfg.Report {
			for _, nad := range notInAD {
				//			var row = []string{"AD group", "group", "fun", "Name", "Uname"}
				var row = []string{"Not in AD", cfg.AdGroup, cfg.Localgroup, nad.Name, nad.Uname}
				excelutils.WriteColumnsln(row)
			}
		}

		if cfg.AddOperation {

			for _, notin := range notInConfluence {
				addUser := confClient.AddGroupMembers(cfg.Localgroup, []string{notin.Uname})

				fmt.Printf("Group: %s status: %s \n", cfg.Localgroup, addUser.Status)
				fmt.Printf("Message: %s \n", addUser.Message)
				fmt.Printf("Users Added: %s \n", addUser.UsersAdded)
				fmt.Printf("Users Skipped: %s \n", addUser.UsersSkipped)

			}
		}
	}
}
func getUnamesInConfluenceGroup(confClient *client.ConfluenceClient, localgroup string, confGroupMemberNames map[string]ad_utils.ADUser) {
	cont := true
	start := 0
	max := 50
	for cont {
		confGroupMembers, err := confClient.GetGroupMembers(localgroup, &client.GetGroupMembersOptions{StartAt: start, MaxResults: max, ShowBasicDetails: true})
		if err != nil {
			panic(err)
		}

		for _, confmember := range confGroupMembers.Users {
			if _, ok := confGroupMemberNames[confmember.Name]; !ok {
				var newUser ad_utils.ADUser
				newUser.Uname = confmember.Name
				newUser.Name = confmember.FullName
				confGroupMemberNames[confmember.Name] = newUser
			}
		}
		if len(confGroupMembers.Users) != max {
			cont = false
		} else {
			start = start + max
		}
	}
}
