package types

const (
	Domain     = "aulaglobal.uc3m.es"
	Webservice = "/webservice/rest/server.php"
	Service    = "aulaglobal_mobile"
	TokenDir   = "token-file"
)

type Prog_args struct {
	Language      int
	UserToken     string
	DirPath       string
	MaxGoroutines int
	CoursesList   []string
}

type File struct {
	FileName string
	FileURL  string
}

type Course struct {
	Name string
	ID   string
}

type FileStore struct {
	FileName string
	FileURL  string
	Dir      string
}

type UserInfo struct {
	FullName string
	UserID   string
}

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
