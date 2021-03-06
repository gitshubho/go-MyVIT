/*
@Author Shubhodeep Mukherjee
@Organization Google Developers Group VIT Vellore
	Isn't Go sexy?
	I know right!!
	Just like Shubhodeep
	I mean, have you seen the guy? xP
	#GDGSwag
*/

package scrape

import (
	"go-MyVIT/api/Godeps/_workspace/src/github.com/headzoo/surf/browser"
	"strings"
)

type Personal struct {
	Name   string
	School string
}

func ShowPersonal(bow *browser.Browser, baseuri string) *Personal {
	name := ""
	school := ""
	if bow.Open(baseuri+"/student/home.asp") == nil {
		table := bow.Find("table").Eq(1)
		tr := table.Find("tr").Eq(0)
		font := tr.Find("font").Eq(0)
		s := strings.Split(strings.TrimSpace(font.Text())[10:], "-")
		name = strings.Title(strings.ToLower(strings.TrimSpace(s[0])))
		school = strings.TrimSpace(s[2])
	}
	return &Personal{
		Name:   name,
		School: school,
	}
}
