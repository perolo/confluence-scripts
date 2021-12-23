package main

import (
	"flag"
	"fmt"
	"github.com/magiconair/properties"
	"github.com/perolo/confluence-prop/client"
	"github.com/perolo/confluence-scripts/utilities"
	"log"
	"os"
)

type Config struct {
	ConfHost    string `properties:"confhost"`
	ConfPass    string `properties:"confpass"`
	ConfUser    string `properties:"confuser"`
	UseToken    bool   `properties:"usetoken"`
	ConfToken   string `properties:"conftoken"`
	File        string `properties:"file"`
	ConfPage    string `properties:"confluencepage"`
	ConfSpace   string `properties:"confluencespace"`
	ConfAttName string `properties:"conlfuenceattachment"`
	//TODO Add remove of previous versions to reduce disk space
	//TODO Set range of history
}

func main() {
	propPtr := flag.String("prop", "uploadattachment.properties", "a properties file")

	flag.Parse()
	p := properties.MustLoadFile(*propPtr, properties.ISO_8859_1)
	var cfg Config
	if err := p.Decode(&cfg); err != nil {
		log.Fatal(err)
	}
	var config = client.ConfluenceConfig{}
	config.Username = cfg.ConfUser
	config.Password = cfg.ConfPass
	config.UseToken = cfg.UseToken
	config.URL = cfg.ConfHost
	config.Debug = false
	if cfg.UseToken {
		config.Password = cfg.ConfToken
	}

	var copt client.OperationOptions
	confluenceClient := client.Client(&config)

	copt.Filepath = cfg.ConfAttName

	copt.SpaceKey = cfg.ConfSpace
	copt.Title = cfg.ConfPage
	// TODO Add verify Space OK
	// TODO Add verify Page OK
	// TODO Add verify Attachment OK
	err := utilities.AddAttachmentAndUpload(confluenceClient, copt, cfg.ConfAttName, cfg.File, "Uploaded by uploadattachment")
	if err != nil {
		fmt.Printf("%s\n", err.Error())
		os.Exit(1)
	}
}
