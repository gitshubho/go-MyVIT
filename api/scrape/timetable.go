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
	"go-MyVIT/api/Godeps/_workspace/src/github.com/PuerkitoBio/goquery"
	"go-MyVIT/api/Godeps/_workspace/src/github.com/headzoo/surf/browser"
	"os"
	"strconv"
	"strings"
	"sync"
)

type Timetable struct {
	Status     string              `json:"status"`
	Time_table map[string]Contents `json:"time_table"`
}

type Contents struct {
	Class_number        int          `json:"class_number"`
	Course_code         string       `json:"course_code"`
	Course_mode         string       `json:"course_mode"`
	Course_option       string       `json:"course_option"`
	Course_title        string       `json:"course_title"`
	Course_type         string       `json:"subject_type"`
	Faculty             string       `json:"faculty"`
	Ltpjc               string       `json:"ltpc"`
	Registration_status string       `json:"registration_status"`
	Slot                string       `json:"slot"`
	Venue               string       `json:"venue"`
	BillDate            string       `json:"bill_date"`
	BillNumber          string       `json:"bill_number"`
	ProjectTitle        string       `json:"project_title"`
	Timings             []TimeStruct `json:"timings"`
	Attendance          Subject      `json:"attendance"`
	Marks               Mrks         `json:"marks"`
}

type TimeStruct struct {
	Day       int    `json:"day"`
	StartTime string `json:"start_time"`
	EndTime   string `json:"end_time"`
}

/*
Function to show Timetable,
Calls NewLogin to login to academics,
@param bow (surf Browser) registration_no password
@return Timtable struct
*/
func ShowTimetable4(bow *browser.Browser, baseuri string) *Timetable {
	sem := os.Getenv("SEM")
	conts := make(map[string]Contents)
	status := "Success"
	if 1 != 1 {
		status = "Failure"
	} else {
		bow.Open(baseuri + "/student/course_regular.asp?sem=" + sem)
		//Twice loading due to Redirect policy defined by academics.vit.ac.in
		if bow.Open(baseuri+"/student/course_regular.asp?sem="+sem) == nil {
			tables := bow.Find("table")
			reg_table := tables.Eq(1)

			tr := reg_table.Find("tr")
			tr_len := tr.Length()
			var wg sync.WaitGroup
			tr.Each(func(i int, s *goquery.Selection) {
				if i > 0 && i < tr_len-2 {
					wg.Add(1)
					go func(conts map[string]Contents, s *goquery.Selection) {
						defer wg.Done()
						td := s.Find("td")
						code := td.Eq(3).Text()
						if code == "Embedded Lab" {
							code = td.Eq(1).Text() + "_L"
							cn, _ := strconv.Atoi(td.Eq(0).Text())
							conts[code] = Contents{
								Class_number:  cn,
								Course_code:   td.Eq(1).Text(),
								Course_mode:   td.Eq(5).Text(),
								Course_option: td.Eq(6).Text(),
								Course_title:  td.Eq(2).Text(),
								Course_type:   td.Eq(3).Text(),
								Faculty:       td.Eq(9).Text(),
								Ltpjc:         strings.TrimSpace(td.Eq(4).Text()),
								Slot:          td.Eq(7).Text(),
								Venue:         td.Eq(8).Text(),
							}
						} else {
							if td.Eq(5).Text() == "Lab Only" {
								code = code + "_L"
							} else if code == "Embedded Project" {
								code = code + "_P"
							}
							cn, _ := strconv.Atoi(td.Eq(2).Text())
							conts[code] = Contents{
								Class_number:        cn,
								Course_code:         td.Eq(3).Text(),
								Course_mode:         td.Eq(7).Text(),
								Course_option:       td.Eq(8).Text(),
								Course_title:        td.Eq(4).Text(),
								Course_type:         td.Eq(5).Text(),
								Faculty:             td.Eq(11).Text(),
								Ltpjc:               strings.TrimSpace(td.Eq(6).Text()),
								Registration_status: td.Eq(12).Text(),
								Slot:                td.Eq(9).Text(),
								Venue:               td.Eq(10).Text(),
							}

						}
					}(conts, s)
				}
				wg.Wait()
				if len(conts) == 0 {
					status = "Failure"
				}
			})
		}

	}
	return &Timetable{
		Status:     status,
		Time_table: conts,
	}
}
