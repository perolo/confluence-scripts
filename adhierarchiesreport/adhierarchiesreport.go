package adhierarchiesreport

import (
	"flag"
	"fmt"
	"github.com/magiconair/properties"
	"github.com/perolo/ad-utils"
	"github.com/perolo/confluence-prop/client"
	"github.com/perolo/confluence-scripts/utilities"
	"log"
)

// or through Decode
type Config struct {
	ConfHost        string `properties:"confhost"`
	User            string `properties:"user"`
	Pass            string `properties:"password"`
	Bindusername    string `properties:"bindusername"`
	Bindpassword    string `properties:"bindpassword"`
}


//func main() {
func CreateAdHierarchiesReport(propPtr, adgroup string) {
	var copt client.OperationOptions
	var confluence *client.ConfluenceClient

//	propPtr := flag.String("prop", "confluence.properties", "a string")
	flag.Parse()
	p := properties.MustLoadFile(propPtr, properties.ISO_8859_1)
	var cfg Config
	if err := p.Decode(&cfg); err != nil {
		log.Fatal(err)
	}
	var config = client.ConfluenceConfig{}
	config.Username = cfg.User
	config.Password = cfg.Pass
	config.URL = cfg.ConfHost
	//config.Debug = true

	confluence = client.Client(&config)

	ad_utils.InitAD(cfg.Bindusername, cfg.Bindpassword)
	//cfg.AdGroup = "#AAAB - Group Technology Team"
	var roothier [] ad_utils.ADHierarchy
	var newhierarchy ad_utils.ADHierarchy
	newhierarchy.Name = adgroup
	newhierarchy.Parent = ""
	roothier = append(roothier, newhierarchy)

	groups, hier, err := ad_utils.ExpandHierarchy(adgroup, roothier)
	if err != nil {
		fmt.Printf("Failed to parse AD hierarchy : %s \n", err)
	} else {
		fmt.Printf("adUnames(%v): %s \n", len(groups), groups)
		fmt.Printf("adUnames(%v): %s \n", len(hier), hier)
		copt.Title = "GTT Hierarchies - " + adgroup
		copt.SpaceKey = "~per.olofsson@assaabloy.com"
		utilities.CreateAttachmentAndUpload(hier, copt, confluence, "Created by AD Hierarchies Report")
	}
	ad_utils.CloseAD()
}
