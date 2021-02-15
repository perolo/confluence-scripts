package gitlabmergestatus

type ReportsType struct {
	Host string
	Token string
	ProjId int
	PageName string
}

var Reports = []ReportsType{
	{
		Host: "https://gitlab/api/v4/",
		Token: "sometoken",
		ProjId: 42,
		PageName: "Thanks for the fish",
	},
}

