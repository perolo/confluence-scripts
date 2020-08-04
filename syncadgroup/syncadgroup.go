package main

import (
	"flag"
	"fmt"
	"git.aa.st/perolo/confluence-utils/Utilities"
	"github.com/magiconair/properties"
	"github.com/perolo/confluence-prop/client"
	"log"
)

func difference(a, b []string) []string {
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
func main() {

	propPtr := flag.String("prop", "../confluence.properties", "a string")

	flag.Parse()

	p := properties.MustLoadFile(*propPtr, properties.ISO_8859_1)

	// or through Decode
	type Config struct {
		ConfHost    string `properties:"confhost"`
		User        string `properties:"user"`
		Pass        string `properties:"password"`
		AddOperation bool `properties:"add"`
		ADgroup       string `properties:"adgroup"`
		Confgroup       string `properties:"confgroup"`
		Bindusername string `properties:"bindusername"`
		Bindpassword string `properties:"bindpassword"`
	}

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

	Utilities.InitAD(cfg.Bindusername, cfg.Bindpassword)

	var adUnames []string

	adUnames, _ = Utilities.GetUnamesInGroup(cfg.ADgroup)
	fmt.Printf("adUnames: %s \n", adUnames)

	confGroupMembers := confClient.GetGroupMembers(cfg.Confgroup)
	var confGroupMemberNames []string

	if confGroupMembers.Status == "success" {
		for _,v := range confGroupMembers.Users {
			for kk,_ := range v {
				confGroupMemberNames = append(confGroupMemberNames, kk)
			}
		}
		fmt.Printf("confGroupMemberNames: %s \n", confGroupMemberNames)
	}

	notInConfluence := difference(adUnames,confGroupMemberNames)
	fmt.Printf("notInConfluence: %s \n", notInConfluence)

	notInAD := difference(confGroupMemberNames,adUnames)
	fmt.Printf("notInAD: %s \n", notInAD)

	if (cfg.AddOperation) {
		addUser := confClient.AddGroupMembers(cfg.Confgroup, notInConfluence)

		fmt.Printf("Group: %s status: %s \n", cfg.Confgroup, addUser.Status)

		fmt.Printf("Message: %s \n", addUser.Message)
		fmt.Printf("Users Added: %s \n", addUser.UsersAdded)
		fmt.Printf("Users Skipped: %s \n", addUser.UsersSkipped)
	}
	Utilities.CloseAD()
}
