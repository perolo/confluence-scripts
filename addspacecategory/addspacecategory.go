package addspacecategory

import (
	"fmt"
	"github.com/magiconair/properties"
	"github.com/perolo/confluence-prop/client"
	"log"
)

// or through Decode
type Config struct {
	ConfHost    string `properties:"confhost"`
	User        string `properties:"user"`
	Pass        string `properties:"password"`
	Space       string `properties:"space"`
	SearchLabel string `properties:"searchlabel"`
	AddLabel    string `properties:"addlabel"`
}

func AddSpaceCategory(propPtr string) {
	var cfg Config
	var config client.ConfluenceConfig
	var conf *client.ConfluenceClient

	fmt.Printf("%%%%%%%%%%  Add Space Category %%%%%%%%%%%%%%\n")

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
		spaces := conf.GetSpaces(&opt)
		for _, space := range spaces.Results {
			fmt.Printf("Add Label: %s to Space: %s \n", cfg.AddLabel, space.Name)
			conf.AddSpaceCategory(space.Key, cfg.AddLabel)
		}
		start = start + increase
		if spaces.Size < increase {
			cont = false
		}
	}
}
