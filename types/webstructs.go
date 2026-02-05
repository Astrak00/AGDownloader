package types

import (
	"io"
	"log"
	"net/http"
)

type WebCourse []struct {
	Hiddenbynumsections int `json:"hiddenbynumsections"`
	ID                  int `json:"id"`
	Modules             []struct {
		Afterlink        *string `json:"afterlink"`
		Availabilityinfo string  `json:"availabilityinfo,omitempty"`
		Completion       int     `json:"completion,omitempty"`
		Contents         []struct {
			Author         *string `json:"author"`
			Filename       string  `json:"filename"`
			Filepath       *string `json:"filepath"`
			Filesize       int     `json:"filesize"`
			Fileurl        string  `json:"fileurl"`
			Isexternalfile bool    `json:"isexternalfile"`
			License        *string `json:"license"`
			Mimetype       string  `json:"mimetype,omitempty"`
			Sortorder      *int    `json:"sortorder"`
			Timecreated    *int    `json:"timecreated"`
			Timemodified   int     `json:"timemodified"`
			Type           string  `json:"type"`
			Userid         *int    `json:"userid"`
		} `json:"contents,omitempty"`
		Contentsinfo *struct {
			Filescount     int      `json:"filescount"`
			Filessize      int      `json:"filessize"`
			Lastmodified   int      `json:"lastmodified"`
			Mimetypes      []string `json:"mimetypes"`
			Repositorytype string   `json:"repositorytype"`
		} `json:"contentsinfo,omitempty"`
		Contextid  int    `json:"contextid"`
		Customdata string `json:"customdata"`
		Dates      []struct {
			Dataid    string `json:"dataid"`
			Label     string `json:"label"`
			Timestamp int    `json:"timestamp"`
		} `json:"dates"`
		Description         string `json:"description"`
		Downloadcontent     int    `json:"downloadcontent"`
		ID                  int    `json:"id"`
		Indent              int    `json:"indent"`
		Instance            int    `json:"instance"`
		Modicon             string `json:"modicon"`
		Modname             string `json:"modname"`
		Modplural           string `json:"modplural"`
		Name                string `json:"name"`
		Noviewlink          bool   `json:"noviewlink"`
		Onclick             string `json:"onclick"`
		URL                 string `json:"url,omitempty"`
		Uservisible         bool   `json:"uservisible"`
		Visible             int    `json:"visible"`
		Visibleoncoursepage int    `json:"visibleoncoursepage"`
	} `json:"modules"`
	Name          string `json:"name"`
	Section       int    `json:"section"`
	Summary       string `json:"summary"`
	Summaryformat int    `json:"summaryformat"`
	Uservisible   bool   `json:"uservisible"`
	Visible       int    `json:"visible"`
}

type WebUser []struct {
	Category                 int    `json:"category"`
	Completed                any    `json:"completed"`
	Completionhascriteria    bool   `json:"completionhascriteria"`
	Completionusertracked    bool   `json:"completionusertracked"`
	Displayname              string `json:"displayname"`
	Enablecompletion         any    `json:"enablecompletion"`
	Enddate                  any    `json:"enddate"`
	Format                   any    `json:"format"`
	Fullname                 string `json:"fullname"`
	Hidden                   bool   `json:"hidden"`
	ID                       int    `json:"id"`
	Idnumber                 string `json:"idnumber"`
	Isfavourite              bool   `json:"isfavourite"`
	Lang                     string `json:"lang"`
	Lastaccess               any    `json:"lastaccess"`
	Marker                   any    `json:"marker"`
	Overviewfiles            []any  `json:"overviewfiles"`
	Progress                 any    `json:"progress"`
	Shortname                string `json:"shortname"`
	Showactivitydates        bool   `json:"showactivitydates"`
	Showcompletionconditions any    `json:"showcompletionconditions"`
	Showgrades               any    `json:"showgrades"`
	Startdate                int    `json:"startdate"`
	Summary                  string `json:"summary"`
	Summaryformat            int    `json:"summaryformat"`
	Timemodified             any    `json:"timemodified"`
	Visible                  int    `json:"visible"`
}

type TimelineCourses struct {
	Courses []struct {
		ID                       int    `json:"id"`
		Fullname                 string `json:"fullname"`
		Shortname                string `json:"shortname"`
		Idnumber                 string `json:"idnumber"`
		Summary                  string `json:"summary"`
		Summaryformat            int    `json:"summaryformat"`
		Startdate                int    `json:"startdate"`
		Enddate                  int    `json:"enddate"`
		Visible                  bool   `json:"visible"`
		Showactivitydates        bool   `json:"showactivitydates"`
		Showcompletionconditions any    `json:"showcompletionconditions"`
		Fullnamedisplay          string `json:"fullnamedisplay"`
		Viewurl                  string `json:"viewurl"`
		Courseimage              string `json:"courseimage"`
		Progress                 int    `json:"progress"`
		Hasprogress              bool   `json:"hasprogress"`
		Isfavourite              bool   `json:"isfavourite"`
		Hidden                   bool   `json:"hidden"`
		Showshortname            bool   `json:"showshortname"`
		Coursecategory           string `json:"coursecategory"`
	} `json:"courses"`
}

func GetJson(URL string) []byte {
	// Get the json from the URL
	resp, err := http.Get(URL)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	// Read the json
	jsonData, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	return jsonData
}
