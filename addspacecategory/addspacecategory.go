package addspacecategory

import (
	"fmt"
	"github.com/magiconair/properties"
	"github.com/perolo/confluence-go-api"
	"log"
)

type Config struct {
	ConfHost    string `properties:"confhost"`
	ConfPass    string `properties:"confpass"`
	ConfUser    string `properties:"confuser"`
	UseToken    bool   `properties:"usetoken"`
	ConfToken   string `properties:"conftoken"`
	Space       string `properties:"space"`
	SearchLabel string `properties:"searchlabel"`
	AddLabel    string `properties:"addlabel"`
}

func AddSpaceCategory(propPtr string) {
	var cfg Config
	var confClient *goconfluence.API
	var err error

	fmt.Printf("%%%%%%%%%%  Add Space Category %%%%%%%%%%%%%%\n")

	p := properties.MustLoadFile(propPtr, properties.ISO_8859_1)

	if err = p.Decode(&cfg); err != nil {
		log.Fatal(err)
	}
	if cfg.UseToken {
		confClient, err = goconfluence.NewAPI(cfg.ConfHost, "", cfg.ConfToken)
	} else {
		confClient, err = goconfluence.NewAPI(cfg.ConfHost, cfg.ConfUser, cfg.ConfPass)
	}
	if err != nil {
		log.Fatal(err)
	}
	//confClient.Debug = true

	start := 0
	increase := 50

	cont := true
	for cont {
		//opt := client.SpaceOptions{Start: start, Limit: increase, Label: cfg.SearchLabel}
		//		opt := goconfluence.AllSpacesQuery{Start: start, Limit: increase, Label: cfg.SearchLabel}
		opt := goconfluence.AllSpacesQuery{Start: start, Limit: increase, Type: "global"}
		//spaces, _ := conf.GetSpaces(&opt)
		spaces, err2 := confClient.GetAllSpaces(opt)
		if err2 != nil {
			log.Fatal(err2)
		}
		for _, space := range spaces.Results {
			fmt.Printf("Add Label: %s to Space: %s \n", cfg.AddLabel, space.Name)
			//conf.AddSpaceCategory(space.Key, cfg.AddLabel)
			_, err3 := confClient.AddSpaceCategory(space.Key, cfg.AddLabel)
			if err3 != nil {
				log.Fatal(err3)
			}
		}
		start = start + increase
		if spaces.Size < int64(increase) {
			cont = false
		}
	}
}
