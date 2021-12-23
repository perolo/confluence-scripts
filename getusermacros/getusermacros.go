package main

import (
	"flag"
	"fmt"
	"github.com/magiconair/properties"
	"github.com/perolo/confluence-prop/client"
	"github.com/perolo/confluence-scripts/getusermacros/macronames"
	"log"
	"os"
)

type Config struct {
	ConfHost  string `properties:"confhost"`
	ConfPass  string `properties:"confpass"`
	ConfUser  string `properties:"confuser"`
	UseToken  bool   `properties:"usetoken"`
	ConfToken string `properties:"conftoken"`
	MacroPath string `properties:"macropath"`
}

var cfg Config

func Check(e error) {
	if e != nil {
		panic(e)
	}
}

func getMacro(confluenceClient *client.ConfluenceClient, key string) (bool, string) {

	rresp, ccont := confluenceClient.GetPage("/admin/updateusermacro-start.action?macro=" + key)

	if ccont.StatusCode == 200 {
		theBody := string(rresp)

		return true, theBody
	} else {
		return false, ""
	}
}

func main() {

	propPtr := flag.String("prop", "jiracategory.properties", "a string")

	flag.Parse()
	fmt.Println(propPtr)

	p := properties.MustLoadFile(*propPtr, properties.ISO_8859_1)

	if err := p.Decode(&cfg); err != nil {
		log.Fatal(err)
	}

	// Access Confluence
	var config = client.ConfluenceConfig{}
	config.Username = cfg.ConfUser
	config.Password = cfg.ConfPass
	config.UseToken = cfg.UseToken
	config.URL = cfg.ConfHost
	//config.Debug = true
	if cfg.UseToken {
		config.Password = cfg.ConfToken
	}

	confluence := client.Client(&config)

	for _, macro := range macronames.MacroNames {
		ok, text := getMacro(confluence, macro)
		if ok {
			fmt.Printf("Saving : %s\n", macro)
			f, err := os.Create(cfg.MacroPath + macro + ".html")
			Check(err)
			_, err = f.Write([]byte(text))
			Check(err)
			err = f.Close()
			Check(err)
		} else {
			fmt.Printf("Failed to retrieve macro: %s\n", macro)
		}

	}

}
