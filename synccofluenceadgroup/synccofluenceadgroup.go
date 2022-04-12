package synccofluenceadgroup

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/magiconair/properties"
	"github.com/perolo/ad-utils"
	"github.com/perolo/confluence-client/client"
	"github.com/perolo/confluence-scripts/utilities"
	"github.com/perolo/excel-utils"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

type Config struct {
	ConfHost        string `properties:"confhost"`
	ConfUser        string `properties:"confuser"`
	ConfPass        string `properties:"confpass"`
	ConfToken       string `properties:"conftoken"`
	UseToken        bool   `properties:"usetoken"`
	Simple          bool   `properties:"simple"`
	AddOperation    bool   `properties:"add"`
	RemoveOperation bool   `properties:"remove"`
	Report          bool   `properties:"report"`
	Limited         bool   `properties:"limited"`
	AutoDisable     bool   `properties:"autodisable"`
	AdGroup         string `properties:"adgroup"`
	Localgroup      string `properties:"localgroup"`
	File            string `properties:"file"`
	ConfUpload      bool   `properties:"confupload"`
	ConfPage        string `properties:"confluencepage"`
	ConfSpace       string `properties:"confluencespace"`
	ConfAttName     string `properties:"conlfuenceattachment"`
	Reset           bool   `properties:"reset"`
	Bindusername    string `properties:"bindusername"`
	Bindpassword    string `properties:"bindpassword"`
	BaseDN          string `properties:"basedn"`
}

func initReport(cfg Config) {
	if cfg.Report {
		excelutils.NewFile()
		excelutils.SetCellFontHeader()
		excelutils.WiteCellln("Introduction")
		excelutils.WiteCellln("Please Do not edit this page!")
		excelutils.WiteCellln("This page is created by the projectreport script: github.com\\perolo\\confluence-scripts\\SyncADGroup")
		t := time.Now()
		excelutils.WiteCellln("Created by: " + cfg.ConfUser + " : " + t.Format(time.RFC3339))
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
		excelutils.WriteColumnsHeaderln([]string{"AD Group", "Local group", "Add", "Remove", "Ad Count", "Local Count"})
		if cfg.Simple {
			excelutils.WriteColumnsln([]string{cfg.AdGroup, cfg.Localgroup, strconv.FormatBool(cfg.AddOperation), strconv.FormatBool(cfg.RemoveOperation)})
		} else {
			for _, syn := range GroupSyncs {
				if syn.InConfluence {
					excelutils.WriteColumnsln([]string{syn.AdGroup, syn.LocalGroup, excelutils.BoolToEmoji(syn.DoAdd), excelutils.BoolToEmoji(syn.DoRemove)})
				}
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

func endReport(cfg Config) error {
	if cfg.Report {
		file := fmt.Sprintf(cfg.File, "-Confluence")
		excelutils.SetAutoColWidth()
		excelutils.SetColWidth("A", "A", 50)
		excelutils.AutoFilterEnd()
		excelutils.SaveAs(file)
		if cfg.ConfUpload {
			var config = client.ConfluenceConfig{}
			var copt client.OperationOptions
			config.Username = cfg.ConfUser
			config.Password = cfg.ConfPass
			config.UseToken = cfg.UseToken
			config.URL = cfg.ConfHost
			config.Debug = false
			confluenceClient := client.Client(&config)
			// Intentional override
			copt.Title = "Using AD groups for JIRA/Confluence"
			copt.SpaceKey = "AAAD"
			_, name := filepath.Split(file)
			cfg.ConfAttName = name
			return utilities.AddAttachmentAndUpload(confluenceClient, copt, name, file, "Created by Sync AD group")

		}
	}
	return nil
}

func ConfluenceSyncAdGroup(propPtr string) {
	//	propPtr := flag.String("prop", "confluence.properties", "a string")
	flag.Parse()
	p := properties.MustLoadFile(propPtr, properties.ISO_8859_1)
	var cfg Config
	if err := p.Decode(&cfg); err != nil {
		log.Fatal(err)
	}
	// Temporary workaround solution - need to find better?
	if cfg.UseToken {
		//cfg.ConfUser = cfg.ConfUser
		cfg.ConfPass = cfg.ConfToken
	}

	toolClient := toollogin(cfg)
	initReport(cfg)
	adutils.InitAD(cfg.Bindusername, cfg.Bindpassword)
	x := 15
	if cfg.Simple {
		SyncGroupInTool(cfg, toolClient)
	} else {
		for _, syn := range GroupSyncs {
			// If this is enabled the reports are partial
			//			if schedulerutil.CheckScheduleDetail(fmt.Sprintf("ConfluenceSyncAdGroup-%s", syn.LocalGroup), time.Hour*24, cfg.Reset, schedulerutil.DummyFunc, "jiracategory.properties") {

			adCount := 0
			groupCount := 0
			if !syn.InJira && !syn.InConfluence {
				log.Fatal("Error in setup")
			}
			if syn.InConfluence {
				cfg.AdGroup = syn.AdGroup
				cfg.Localgroup = syn.LocalGroup
				cfg.AddOperation = syn.DoAdd
				cfg.RemoveOperation = syn.DoRemove
				cfg.AutoDisable = syn.AutoDisable
				adCount, groupCount = SyncGroupInTool(cfg, toolClient)
				// Dirty Solution - find a better?
				excelutils.SetCell(fmt.Sprintf("%v", adCount), 5, x)
				excelutils.SetCell(fmt.Sprintf("%v", groupCount), 6, x)
				if adCount == groupCount {
					excelutils.SetCellStyleColor("green")
				}
			}
			x++
		}
	}
	err := endReport(cfg)
	if err != nil {
		panic(err)
	}
	adutils.CloseAD()
}

func toollogin(cfg Config) *client.ConfluenceClient {
	var config = client.ConfluenceConfig{}
	config.Username = cfg.ConfUser
	config.Password = cfg.ConfPass
	config.UseToken = cfg.UseToken
	config.URL = cfg.ConfHost
	config.Debug = false
	return client.Client(&config)
}

func SyncGroupInTool(cfg Config, client *client.ConfluenceClient) (adcount int, localcount int) {
	var toolGroupMemberNames map[string]adutils.ADUser
	deactCounter := 0
	fmt.Printf("\n")
	fmt.Printf("SyncGroup Confluence AdGroup: %s LocalGroup: %s \n", cfg.AdGroup, cfg.Localgroup)
	fmt.Printf("\n")
	var adUnames []adutils.ADUser
	if cfg.AdGroup != "" {
		adUnames, _ = adutils.GetUnamesInGroup(cfg.AdGroup, cfg.BaseDN)
		fmt.Printf("adUnames(%v)\n", len(adUnames))
		if len(adUnames) == 0 {
			fmt.Printf("Warning empty AD group! adUnames(%v)\n", len(adUnames))
			panic(nil)
		}
	}
	if cfg.Report {
		if !cfg.Limited {
			for _, adu := range adUnames {
				var row = []string{"AD Names", cfg.AdGroup, cfg.Localgroup, adu.Name, adu.Uname, adu.Mail, adu.Err, adu.DN}
				excelutils.WriteColumnsln(row)
			}
		}
		adcount = len(adUnames)
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
		localcount = len(toolGroupMemberNames)
	}
	if cfg.Localgroup != "" && cfg.AdGroup != "" {
		notInTool := adutils.Difference(adUnames, toolGroupMemberNames)
		if len(notInTool) == 0 {
			fmt.Printf("Not In Tool(%v)\n", len(notInTool))
		} else {
			fmt.Printf("Not In Tool(%v) ", len(notInTool))
			for _, nit := range notInTool {
				fmt.Printf("%s, ", nit.Uname)
			}
			fmt.Printf("\n")
		}
		if cfg.Report {
			for _, nji := range notInTool {
				var row = []string{"AD group users not found in Tool user group", cfg.AdGroup, cfg.Localgroup, nji.Name, nji.Uname, nji.Mail, nji.Err, nji.DN}
				excelutils.WriteColumnsln(row)
			}
		}
		notInAD := adutils.Difference2(toolGroupMemberNames, adUnames)
		if len(notInAD) == 0 {
			fmt.Printf("notInAD(%v)\n", len(notInAD))
		} else {
			fmt.Printf("notInAD(%v) ", len(notInAD))
			for _, nit := range notInAD {
				fmt.Printf("%s, ", nit.Uname)
			}
			fmt.Printf("\n")
		}
		if cfg.Report {
			for _, nad := range notInAD {
				if nad.DN == "" {

					dn, err := adutils.GetActiveUserDN(nad.Uname, cfg.BaseDN)
					if err == nil {
						nad.DN = dn.DN
						nad.Mail = dn.Mail
						nad.Name = dn.Name
					} else {
						udn, err := adutils.GetAllUserDN(nad.Uname, cfg.BaseDN)
						if err == nil {
							nad.DN = udn.DN
							nad.Mail = udn.Mail
							nad.Name = udn.Name
							nad.Err = "Deactivated"
							// Avoid being kicked out
							deactCounter++
							if deactCounter > 10 {
								_, errn := adutils.GetAllUserDN("perolo", cfg.BaseDN)
								if errn != nil {
									fmt.Printf("Error: finding %s \n", "perolo")
									panic(errn)
								}
								deactCounter = 0
							}
							if cfg.AutoDisable {
								TryDeactivateUserConfluence(client, nad.Uname)
							}
						} else {
							edn, err := adutils.GetAllEmailDN(nad.Mail, cfg.BaseDN)
							if err == nil {
								nad.DN = edn[0].DN
								nad.Mail = edn[0].Mail
								nad.Err = edn[0].Err
								for _, ldn := range edn {
									var row2 = []string{"Tool user group member not found in AD group (multiple?)", cfg.AdGroup, cfg.Localgroup, nad.Name, nad.Uname, ldn.Mail, ldn.Err, ldn.DN}
									excelutils.WriteColumnsln(row2)
								}
							} else {

								nad.Err = err.Error()
							}
						}
					}

				}
				var row = []string{"Tool user group member not found in AD group", cfg.AdGroup, cfg.Localgroup, nad.Name, nad.Uname, nad.Mail, nad.Err, nad.DN}
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
					fmt.Printf("About to remove user. Group: %s uname: %s Name: %s \n", cfg.Localgroup, notin.Uname, notin.Name)

					fmt.Printf("Remove [y/n]: ")

					reader := bufio.NewReader(os.Stdin)
					response, err := reader.ReadString('\n')
					if err != nil {
						log.Fatal(err)
					}

					response = strings.ToLower(strings.TrimSpace(response))

					if response == "y" || response == "yes" {
						removeUser := client.RemoveGroupMembers(cfg.Localgroup, []string{notin.Uname})
						fmt.Printf("Remove user. Group: %s status: %s \n", cfg.Localgroup, removeUser.Status)
						fmt.Printf("Message: %s \n", removeUser.Message)
						fmt.Printf("Users Removed: %s \n", removeUser.UsersRemoved)
						fmt.Printf("Users Skipped: %s \n", removeUser.UsersSkipped)

					} else {
						fmt.Printf("Respone No -  skipping remove: %s \n", notin.Uname)
					}
					//fmt.Printf("Not Yet Implemented\n")
				}
			}
		}

	}
	return adcount, localcount
}

func TryDeactivateUserConfluence(client *client.ConfluenceClient, deactuser string) {
	deactUser, resp := client.GetUserDetails(deactuser)
	if resp.StatusCode == 200 {
		if deactUser.HasAccessToUseConfluence {
			fmt.Printf("Deactivating User: %s  \n", deactuser)
			mess, resp2 := client.DeactivateUser(deactuser)
			fmt.Printf("User: %s Deactivated, message: %s response: %v \n", deactuser, mess.Message, resp2.StatusCode)
			if resp2.StatusCode != 200 {
				fmt.Printf("User: %s Deactivated, message: %s response: %v \n", deactuser, mess.Message, resp2.StatusCode)
			}
		} else {
			fmt.Printf("User: %s, Already Deactivated \n", deactuser)
		}
	} else {
		fmt.Printf("Error: \n")
		panic(nil)
	}
}
func getUnamesInToolGroup(theClient *client.ConfluenceClient, localgroup string) map[string]adutils.ADUser {
	groupMemberNames := make(map[string]adutils.ADUser)
	cont := true
	start := 0
	max := 50
	for cont {
		groupMembers, resp, err := theClient.GetGroupMembers(localgroup, &client.GetGroupMembersOptions{StartAt: start, MaxResults: max, ShowExtendedDetails: true})

		if err != nil {
			panic(err) // theClient.AddGroups(localgroup)
		} else {
			if resp.StatusCode == 500 {
				theClient.AddGroups([]string{localgroup})
			}
		}
		for _, member := range groupMembers.Users {
			if _, ok := groupMemberNames[member.Name]; !ok {
				if member.HasAccessToUseConfluence {
					var newUser adutils.ADUser
					newUser.Uname = member.Name
					newUser.Name = member.FullName
					newUser.Mail = member.Email
					groupMemberNames[member.Name] = newUser
				}
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
