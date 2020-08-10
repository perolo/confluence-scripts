package main

import (
	"flag"
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/magiconair/properties"
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

	propPtr := flag.String("prop", "jira.properties", "a string")

	flag.Parse()

	p := properties.MustLoadFile(*propPtr, properties.ISO_8859_1)

	// or through Decode
	type Config struct {
		//ConfHost     string `properties:"confhost"`
		User      string `properties:"user"`
		Pass      string `properties:"password"`
		JiraHost  string `properties:"jirahost"`
		ExcelFile string `properties:"excelfile"`
		/*
			AddOperation bool   `properties:"add"`
			ADgroup      string `properties:"adgroup"`
			Confgroup    string `properties:"confgroup"`
			Bindusername string `properties:"bindusername"`
			Bindpassword string `properties:"bindpassword"`

		*/
	}
	var cfg Config
	if err := p.Decode(&cfg); err != nil {
		log.Fatal(err)
	}
	fexcel, err := excelize.OpenFile(cfg.ExcelFile)
	if err != nil {
		fmt.Println(err)
		return
	}
	// Get value from cell by given worksheet name and axis.
	cell, err := fexcel.GetCellValue("Platform Summary", "B2")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(cell)
	// Get all the rows in the Sheet1.
	rows, err := fexcel.GetRows("Platform Summary")
	for _, row := range rows {
		for _, colCell := range row {
			fmt.Print(colCell, "\t")
		}
		fmt.Println()
	}

	// Save xlsx file by the given path.
	if err := fexcel.SaveAs(cfg.ExcelFile); err != nil {
		fmt.Println(err)
	}
/*
	var err error
	xlFile, err := xlsx.OpenFile(cfg.ExcelFile)
	if err != nil {
		fmt.Printf("Result: %v\n", err.Error())
		panic(err)
	}
	for _, sheet := range xlFile.Sheets {


		fmt.Printf("Cell: %s\n", sheet.Rows[3].Cells[3].Value)

		somerows := sheet.Rows[74:1000]
		for _, row := range somerows {
		}

	}
*/
/*
	tp := jira.BasicAuthTransport{
		Username: strings.TrimSpace(cfg.User),
		Password: strings.TrimSpace(cfg.Pass),
	}

	jiraClient, err := jira.NewClient(tp.Client(), strings.TrimSpace(cfg.JiraHost))
	if err != nil {
		fmt.Printf("\nerror: %v\n", err)
		return
	}

	components, _ , err := jiraClient.Project.GetComponents("RFA")
	if err != nil {
		fmt.Printf("Result: %v\n", err.Error())
		panic(err)
	}
	for _, component := range *components {
		if (component.Archived == false) {
			fmt.Printf("Component: %s , %s, %s, %s \n", component.Name, component.AssigneeType, component.RealAssignee.Name, component.Description)
		}
	}

*/
/*

	var err error
	xlFile, err = xlsx.OpenFile(excelFileName)
	if err != nil {
		fmt.Printf("Result: %v\n", err.Error())
		panic(err)
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
		for _, v := range confGroupMembers.Users {
			for kk, _ := range v {
				confGroupMemberNames = append(confGroupMemberNames, kk)
			}
		}
		fmt.Printf("confGroupMemberNames: %s \n", confGroupMemberNames)
	}

	notInConfluence := difference(adUnames, confGroupMemberNames)
	fmt.Printf("notInConfluence: %s \n", notInConfluence)

	notInAD := difference(confGroupMemberNames, adUnames)
	fmt.Printf("notInAD: %s \n", notInAD)

	if cfg.AddOperation {
		addUser := confClient.AddGroupMembers(cfg.Confgroup, notInConfluence)

		fmt.Printf("Group: %s status: %s \n", cfg.Confgroup, addUser.Status)

		fmt.Printf("Message: %s \n", addUser.Message)
		fmt.Printf("Users Added: %s \n", addUser.UsersAdded)
		fmt.Printf("Users Skipped: %s \n", addUser.UsersSkipped)
	}
	Utilities.CloseAD()

*/
}
