package gitlabuserreport

import (
	"fmt"
	"github.com/magiconair/properties"
	"github.com/perolo/go-gitlab"
	"log"
)

type Config struct {
	GitLabHost  string `properties:"gitlabhost"`
	GitLabtoken string `properties:"gitlabtoken"`
	GitProjId   int    `properties:"gitprojid"`
}

var cfg Config

func GitLabUserReport(propPtr string) {

	fmt.Printf("%%%%%%%%%%  GitLabUserReport %%%%%%%%%%%%%%\n")

	p := properties.MustLoadFile(propPtr, properties.ISO_8859_1)

	if err := p.Decode(&cfg); err != nil {
		log.Fatal(err)
	}

	gitlabclient, err := gitlab.NewClient(cfg.GitLabtoken, gitlab.WithBaseURL(cfg.GitLabHost))
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	users, _, err := gitlabclient.Users.ListUsers(&gitlab.ListUsersOptions{})

	if err != nil {
		fmt.Printf("User : %v\n", err.Error())
	}
	for _, auser := range users {

		fmt.Printf("User: %s\n", auser.Username)
	}
	groups, _, err := gitlabclient.Groups.ListGroups(&gitlab.ListGroupsOptions{})
	if err != nil {
		fmt.Printf("Group: %v\n", err.Error())
	}
	var listGroupMembersOptions gitlab.ListGroupMembersOptions
	listGroupMembersOptions.Page = 0
	listGroupMembersOptions.PerPage = 30
	for _, agroup := range groups {

		fmt.Printf("Group: %s\n", agroup.Name)
		members, _, err := gitlabclient.Groups.ListGroupMembers(agroup.ID, &listGroupMembersOptions)
		if err != nil {
			fmt.Printf("Member: %v\n", err.Error())
		}
		for _, member := range members {
			fmt.Printf("	Member: %s Login: %s\n", member.Name, member.Username)
		}

	}
}
