package scrape

import (
	"go-MyVIT/api/Godeps/_workspace/src/github.com/headzoo/surf/browser"
	"go-MyVIT/api/Godeps/_workspace/src/github.com/PuerkitoBio/goquery"
	"go-MyVIT/api/login"
	"strings"
	"sync"
)

type AcademicStruct struct {
	GradeSummary Grades `json:"grade summary"`
	History1 map[string]CourseDets `json:"history 1"`
	History2 StudentDets `json:"history 2"`
	Status string `json:"status"`
}

type CourseDets struct{
	CourseTitle string `json:"course_title"`
	CourtType string `json:"course_type"`
	Credit string `json:"credit"`
	Grade string `json:"grade"`
}

type StudentDets struct{
	CGPA string `json:"cgpa"`
	CEarned string `json:"credits earned"`
	CRegistered string `json:"credits registered"`
	Rank string `json:"rank"`
}

type Grades struct{
	A string `json:"A grades"`
	B string `json:"B grades"`
	C string `json:"C grades"`
	D string `json:"D grades"`
	E string `json:"E grades"`
	F string `json:"F grades"`
	N string `json:"N grades"`
	S string `json:"S grades"`
}

func Academics(bow *browser.Browser,regno, password, baseuri string) *AcademicStruct{
	response := login.NewLogin(bow,regno,password,baseuri)
	status := "Success"
	var wg sync.WaitGroup
	history1 := make(map[string]CourseDets)
	var history2 StudentDets
	var grade Grades
	if response.Status==0 {
		status = "Failure"
	} else {
		bow.Open(baseuri+"/student/student_history.asp")
		bow.Open(baseuri+"/student/student_history.asp")
		table := bow.Find("table").Eq(2)
		tr := table.Find("tr")
		wg.Add(3)
		go func(){
			defer wg.Done()
			tr.Each(func(i int, s *goquery.Selection){
				if i>0 {
					td := s.Find("td")
					history1[td.Eq(1).Text()] = CourseDets{
						CourseTitle: td.Eq(2).Text(),
						CourtType: td.Eq(3).Text(),
						Credit: td.Eq(4).Text(),
						Grade: td.Eq(5).Text(),
					}
				}
			})
		}()
		go func(){
			defer wg.Done()
			table = bow.Find("table").Eq(3)
			td := table.Find("tr").Eq(1).Find("td")
			history2 = StudentDets{
				CGPA : strings.TrimSpace(td.Eq(2).Text()),
				CEarned: td.Eq(1).Text(),
				CRegistered: td.Eq(0).Text(),
				Rank: td.Eq(3).Text(),
			}
		}()
		go func() {
			defer wg.Done()
			table = bow.Find("table").Eq(4)
			td := table.Find("tr").Eq(1).Find("td")
			grade = Grades{
				A: td.Eq(1).Text(),
				B: td.Eq(2).Text(),
				C: td.Eq(3).Text(),
				D: td.Eq(4).Text(),
				E: td.Eq(5).Text(),
				F: td.Eq(6).Text(),
				N: td.Eq(7).Text(),
				S: td.Eq(0).Text(),
			}
		}()
	}
	wg.Wait()
	return &AcademicStruct{
		GradeSummary: grade,
		History1: history1,
		History2: history2,
		Status: status,
	}
}