package main

import (
	"flag"
	"github.com/magiconair/properties"
	"github.com/perolo/confluence-prop/client"
	"github.com/perolo/confluence-scripts/utilities"
	"log"
)

// or through Decode
type Config struct {
	ConfHost    string `properties:"confhost"`
	User        string `properties:"user"`
	Pass        string `properties:"password"`
	File        string `properties:"file"`
	ConfPage    string `properties:"confluencepage"`
	ConfSpace   string `properties:"confluencespace"`
	ConfAttName string `properties:"conlfuenceattachment"`
}

func main() {
	propPtr := flag.String("prop", "uploadattachment.properties", "a properties file")

	//	propPtr := flag.String("prop", "confluence.properties", "a string")
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

	var copt client.OperationOptions
	confluenceClient := client.Client(&config)
	// Intentional override
	//copt.Title = "Using AD groups for JIRA/Confluence"
	//copt.SpaceKey = "STPIM"
	copt.Filepath = cfg.ConfAttName

	//_, name := filepath.Split(file)
	//cfg.ConfAttName = name
	copt.SpaceKey = cfg.ConfSpace
	copt.Title = cfg.ConfPage
	utilities.AddAttachmentAndUpload(confluenceClient, copt, cfg.ConfAttName, cfg.File, "Uploaded by uploadattachment")
}
