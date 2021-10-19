package sonarqubeuserreport

import (
	"fmt"
	"github.com/magiconair/properties"
	"github.com/perolo/sonarqube-client/sonarclient"
	"log"
)

type Config struct {
	Host     string `properties:"host"`
	User     string `properties:"user"`
	Pass     string `properties:"password"`
}

var sonarClient *sonarclient.SonarQubeClient

var cfg Config

func Sonarqubeuserreport(propPtr string) {

	fmt.Printf("%%%%%%%%%%  SonarQube User Report  %%%%%%%%%%%%%%\n")

	p := properties.MustLoadFile(propPtr, properties.ISO_8859_1)

	if err := p.Decode(&cfg); err != nil {
		log.Fatal(err)
	}

	var sonaConfig = sonarclient.SonarQubeConfig{}
	sonaConfig.Username = cfg.User
	sonaConfig.Password = cfg.Pass
	sonaConfig.URL = cfg.Host
	//config.Debug = true

	sonarClient = sonarclient.Client(&sonaConfig)

	groups := sonarClient.GetGroups()
	fmt.Printf("Group Count: %v\n", len(groups.Groups))
	for _, agroup := range groups.Groups {

		fmt.Printf("Group: %s\n", agroup.Name)
		members := sonarClient.GetGroupMembers(agroup.ID)
		for _, member := range members.Users {
			fmt.Printf("	Member: %s Login: %s \n", member.Name, member.Login)
		}
	}
}
