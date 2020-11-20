package main

import (
	"flag"
	"fmt"
	"github.com/magiconair/properties"
	"github.com/perolo/ad-utils"
	"github.com/perolo/confluence-prop/client"
	"log"
)

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
	toolGroupMemberNames := make(map[string]string)

	SyncGroupInConfluence(adUnames, cfg, confClient, toolGroupMemberNames)
	ad_utils.CloseAD()
}

func SyncGroupInConfluence(adUnames []string, cfg Config, confClient *client.ConfluenceClient, toolGroupMemberNames map[string]string) {
	fmt.Printf("\n")
	fmt.Printf("SyncGroupInConfluence AdGroup: %s LocalGroup: %s \n", cfg.AdGroup, cfg.Localgroup)
	fmt.Printf("\n")
	adUnames, _ = ad_utils.GetUnamesInGroup(cfg.AdGroup)
	fmt.Printf("adUnames(%v): %s \n", len(adUnames), adUnames)

	getUnamesInConfluenceGroup(confClient, cfg, toolGroupMemberNames)

	notInConfluence := difference(adUnames, toolGroupMemberNames)
	fmt.Printf("notInConfluence(%v): %s \n", len(notInConfluence), notInConfluence)

	notInAD := difference2(toolGroupMemberNames, adUnames)
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
	cont := true
	start := 0
	max := 50
	for cont {
		confGroupMembers, err := confClient.GetGroupMembers(cfg.Localgroup, &client.GetGroupMembersOptions{StartAt: start, MaxResults: max, ShowBasicDetails: true})
		if err != nil {
			panic(err)
		}

		for _, confmember := range confGroupMembers.Users {
			if _, ok := confGroupMemberNames[confmember.Name]; !ok {
				confGroupMemberNames[confmember.Name] = confmember.Name
			}
		}
		if len(confGroupMembers.Users) != max {
			cont = false
		} else {
			start = start + max
		}
	}
}
