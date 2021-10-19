package removespacecategory

import (
	"fmt"
	"github.com/magiconair/properties"
	"github.com/perolo/confluence-prop/client"
	"log"
)

type Config struct {
	ConfHost    string `properties:"confhost"`
	User        string `properties:"user"`
	Pass        string `properties:"password"`
	Space       string `properties:"space"`
	SearchLabel string `properties:"searchlabel"`
}

func RemoveSpaceCategory(propPtr string) {
	var cfg Config
	var config client.ConfluenceConfig
	var conf *client.ConfluenceClient

	fmt.Printf("%%%%%%%%%%  Re Space Category %%%%%%%%%%%%%%\n")

	p := properties.MustLoadFile(propPtr, properties.ISO_8859_1)

	if err := p.Decode(&cfg); err != nil {
		log.Fatal(err)
	}
	config = client.ConfluenceConfig{}
	config.Username = cfg.User
	config.Password = cfg.Pass
	config.URL = cfg.ConfHost
	//config.Debug = true

	conf = client.Client(&config)

	start := 0
	increase := 50

	cont := true
	for cont {
		opt := client.SpaceOptions{Start: start, Limit: increase, Label: cfg.SearchLabel}
		spaces, _ := conf.GetSpaces(&opt)
		for _, space := range spaces.Results {
			found := false
			id := 0

			cats := conf.GetSpaceCategories(space.Key)
			for _, cat := range cats.Categories {
				if cat.NiceName == cfg.SearchLabel {
					found = true
					id = cat.ID
				}
			}
			if found {
				fmt.Printf("Remove Label: %s from Space: %s \n", cfg.SearchLabel, space.Name)
				conf.RemoveSpaceCategory(space.Key, id)
			} else {
				fmt.Printf("Failed to Remove Label!: %s from Space: %s \n", cfg.SearchLabel, space.Name)
			}
		}
		start = start + increase
		if spaces.Size < increase {
			cont = false
		}
	}
}
