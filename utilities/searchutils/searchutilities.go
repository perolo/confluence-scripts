package searchutils

import (
	"github.com/perolo/confluence-client/client"
	"golang.org/x/net/html"
	"regexp"
	"strings"
)

func SearchSpacePage(clie *client.ConfluenceClient, key string) (bool, *client.ConfluencePage) {
	cql := "space=" + key + " AND label=space_description"
	retry := 3
	for retry > 3 {
		results := clie.SearchPage(cql)
		if results.Size == 1 {
			//	log.Println("Page found")
			item := results.Results[0]
			return true, &item
		} else {

			retry--
		}
	}
	return false, nil
}

func GetArray(txt string) []string {
	var body = strings.NewReader("<html>" + txt + "</html>")
	z := html.NewTokenizer(body)
	content := []string{}

	// While have not hit the </html> tag
	for z.Token().Data != "html" {
		tt := z.Next()
		if tt == html.StartTagToken {
			t := z.Token()
			//fmt.Println(t.Data)
			if t.Data == "td" || t.Data == "th" {
				qq := ""
				inner := z.Next()
				text2 := (string)(z.Raw())
				//fmt.Println(text2)
				if inner == html.StartTagToken {
					tagcounter := 1
					qq = text2
					for {
						inner = z.Next()
						text3 := (string)(z.Raw())
						//fmt.Println(text3)
						qq = qq + text3
						switch inner {
						case html.StartTagToken:
							tagcounter++
						case html.EndTagToken:
							tagcounter--
						}
						if tagcounter == 0 {
							content = append(content, qq)
							break
						}
					}
				}
				if inner == html.TextToken {
					text := (string)(z.Text())
					t := qq + strings.TrimSpace(text)
					content = append(content, t)
				}
			}
		}
	}
	// Print to check the slice's content
	//fmt.Println(content)
	return content
}
func GetOwner(clie *client.ConfluenceClient, page *client.ConfluencePage) (bool, string) {
	results := clie.GetPageByID(page.ID)

	if len(results.Body.View.Value) > 0 {
		//	log.Println("Page found")
		vec := GetArray(results.Body.View.Value)
		for cc, it := range vec {
			if it == "Owner" {
				item := vec[cc+1]
				return true, item
			}
		}

	}
	return false, ""

}

func GetSplitLine(txt string) (before, after string) {

	zp := regexp.MustCompile(`<hr/>`)
	vec := zp.Split(txt, -1)
	//fmt.Printf("%q\n", ) // ["pi" "a"]

	//	var body= strings.NewReader("<html>" + txt + "</html>")
	//	z := html.NewTokenizer(body)
	before = vec[0]
	after = vec[1]
	return before, after
}
