package main

import (
	"flag"
	"fmt"
	"github.com/magiconair/properties"
	"github.com/perolo/ad-utils"
	"github.com/perolo/confluence-prop/client"
	"log"
)
type GroupSyncType struct {
	AdGroup    string
	LocalGroup string
}
var GroupSyncs = []GroupSyncType{
	{AdGroup: "AD Group 1", LocalGroup: "Local 1"},
	{AdGroup: "AD Group 2", LocalGroup: "Local 2"},
}
func difference(a []string, b map[string]string) []string {
	mb := make(map[string]struct{}, len(b))
	for _, x := range b {
		mb[x] = struct{}{}
	}
	var diff []string
	for _, x := range a {
		if _, found := mb[x]; !found {
			diff = append(diff, x)
		}
	}
	return diff
}
func difference2(a map[string]string, b []string) []string {
	mb := make(map[string]struct{}, len(b))
	for _, x := range b {
		mb[x] = struct{}{}
	}
	var diff []string
	for _, x := range a {
		if _, found := mb[x]; !found {
			diff = append(diff, x)
		}
	}
	return diff
}

// or through Decode
type Config struct {
	ConfHost     string `properties:"confhost"`
	User         string `properties:"user"`
	Pass         string `properties:"password"`
	AddOperation bool   `properties:"add"`
	AdGroup      string `properties:"adgroup"`
	Localgroup   string `properties:"localgroup"`
	Bindusername string `properties:"bindusername"`
	Bindpassword string `properties:"bindpassword"`
}


func main() {

	propPtr := flag.String("prop", "../confluence.properties", "a string")

	flag.Parse()

	p := properties.MustLoadFile(*propPtr, properties.ISO_8859_1)

	var cfg Config
	if err := p.Decode(&cfg); err != nil {
		log.Fatal(err)
	}

	var config = client.ConfluenceConfig{}
	config.Username = cfg.User
	config.Password = cfg.Pass
	config.URL = cfg.ConfHost
	config.Debug = false

	confClient := client.Client(&config)

	ad_utils.InitAD(cfg.Bindusername, cfg.Bindpassword)

	for _, syn := range GroupSyncs {
		var adUnames []string
		confGroupMemberNames := make(map[string]string)
		cfg.AdGroup = syn.AdGroup
		cfg.Localgroup = syn.LocalGroup
		SyncGroupInConfluence(adUnames, cfg, confClient, confGroupMemberNames)
	}

	ad_utils.CloseAD()
}

func main2() {

	propPtr := flag.String("prop", "../confluence.properties", "a string")

	flag.Parse()

	p := properties.MustLoadFile(*propPtr, properties.ISO_8859_1)

	var cfg Config
	if err := p.Decode(&cfg); err != nil {
		log.Fatal(err)
	}

	var config = client.ConfluenceConfig{}
	config.Username = cfg.User
	config.Password = cfg.Pass
	config.URL = cfg.ConfHost
	config.Debug = false

	confClient := client.Client(&config)

	ad_utils.InitAD(cfg.Bindusername, cfg.Bindpassword)

	var adUnames []string
	confGroupMemberNames := make(map[string]string)

	SyncGroupInConfluence(adUnames, cfg, confClient, confGroupMemberNames)
	ad_utils.CloseAD()
}

func SyncGroupInConfluence(adUnames []string, cfg Config, confClient *client.ConfluenceClient, confGroupMemberNames map[string]string) {
	fmt.Printf("\n")
	fmt.Printf("SyncGroupInConfluence AdGroup: %s LocalGroup: %s \n", cfg.AdGroup, cfg.Localgroup)
	fmt.Printf("\n")
	adUnames, _ = ad_utils.GetUnamesInGroup(cfg.AdGroup)
	fmt.Printf("adUnames: %s \n", adUnames)

	getUnamesInConfluenceGroup(confClient, cfg, confGroupMemberNames)

	notInConfluence := difference(adUnames, confGroupMemberNames)
	fmt.Printf("notInConfluence: %s \n", notInConfluence)

	notInAD := difference2(confGroupMemberNames, adUnames)
	fmt.Printf("notInAD: %s \n", notInAD)

	if cfg.AddOperation {
		if notInConfluence != nil{

			addUser := confClient.AddGroupMembers(cfg.Localgroup, notInConfluence)

			fmt.Printf("Group: %s status: %s \n", cfg.Localgroup, addUser.Status)
			fmt.Printf("Message: %s \n", addUser.Message)
			fmt.Printf("Users Added: %s \n", addUser.UsersAdded)
			fmt.Printf("Users Skipped: %s \n", addUser.UsersSkipped)

		}
	}
}

func getUnamesInConfluenceGroup(confClient *client.ConfluenceClient, cfg Config, confGroupMemberNames map[string]string) {
	contconf := true
	startconf := 0
	maxconf := 20
	for contconf {
		confGroupMembers, _ := confClient.GetGroupMembers(cfg.Localgroup, &client.GetGroupMembersOptions{StartAt: startconf, MaxResults: maxconf, ShowBasicDetails: true})

		for _, confmember := range confGroupMembers.Users {
			if _, ok := confGroupMemberNames[confmember.Name]; !ok {
				confGroupMemberNames[confmember.Name] = confmember.Name
			}
		}
		if len(confGroupMembers.Users) != maxconf {
			contconf = false
		} else {
			startconf = startconf + maxconf
		}
	}
}
