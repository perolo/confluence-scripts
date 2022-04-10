package main

import (
	"flag"
	"fmt"
	"github.com/magiconair/properties"
	"github.com/perolo/confluence-client/client"
	"log"
)

func main() {

	propPtr := flag.String("prop", "../confluence.properties", "a string")

	flag.Parse()

	p := properties.MustLoadFile(*propPtr, properties.ISO_8859_1)

	// or through Decode
	type Config struct {
		ConfHost    string `properties:"confhost"`
		ConfPass    string `properties:"confpass"`
		ConfUser    string `properties:"confuser"`
		UseToken    bool   `properties:"usetoken"`
		ConfToken   string `properties:"conftoken"`
		Source      string `properties:"source"`
		Destination string `properties:"destination"`
	}
	var cfg Config
	if err := p.Decode(&cfg); err != nil {
		log.Fatal(err)
	}

	var config = client.ConfluenceConfig{}
	config.Username = cfg.ConfUser
	config.Password = cfg.ConfPass
	config.UseToken = cfg.UseToken
	config.URL = cfg.ConfHost
	//config.Debug = true

	if cfg.UseToken {
		config.Password = cfg.ConfToken
	}

	theClient := client.Client(&config)

	start := 0
	cont := true
	increase := 50

	noSpaces := 0
	for cont {
		opt := client.GroupOptions{Start: start, Limit: increase}
		spaces := theClient.GetAllSpacesForGroupPermissions(cfg.Source, &opt)

		for _, space := range spaces.Spaces {
			fmt.Printf("Space name: %s\n", space.Name)
			//			opt2 := client.SpaceOptions{Start: 0, Limit: 10,  Status: "current", SpaceKey: space.Key}
			//			opt2 := client.SpaceOptions{Start: 0, Limit: 10,  Status: "current", SpaceKey: space.Key}
			//			spaces := theClient.GetSpaces(&opt2)
			//			if (spaces.Size ==1 ) {
			p := space.Permissions
			added := theClient.AddSpacePermissionsForGroup(space.Key, cfg.Destination, p)
			fmt.Printf("Permissions added : %s\n", added.Added)
			noSpaces++
			//			} else {
			//				fmt.Printf("Archived Space : %s %v\n", space.Name, spaces.Size )

			//			}
		}

		start = start + increase
		if spaces.Total < start {
			cont = false
		}
	}
	fmt.Printf("Spaces : %d \n", noSpaces)

}
