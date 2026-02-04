package types

const (
	Domain     = "aulaglobal.uc3m.es"
	Webservice = "/webservice/rest/server.php"
	Service    = "aulaglobal_mobile"
	TokenDir   = "aulaglobal-token"
)

type ProgramArgs struct {
	Language      int
	UserToken     string
	DirPath       string
	MaxGoroutines int
	CoursesList   []string
	WebUI         bool
	Timeline      bool
}

// CHeck if all the arguments are assigned
func (p ProgramArgs) CheckAllAsigned() bool {
	if p.Language == 0 || p.UserToken == "" || p.DirPath == "" || p.MaxGoroutines == 0 {
		return false
	}
	return true
}

type File struct {
	FileName string
	FileURL  string
}

type Course struct {
	Name string
	ID   string
}

// Define a named type for a slice of Course
type Courses []Course

// Map the courses to obtain a []string with the names of the courses
func (c Courses) GetCoursesName() []string {
	coursesNames := make([]string, len(c))
	for i, course := range c {
		coursesNames[i] = course.Name
	}
	return coursesNames
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
