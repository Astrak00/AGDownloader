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
	FileType string
}

type Course struct {
	Name string
	ID   string
}

type FileStore struct {
	FileName string
	FileURL  string
	FileType string
	Dir      string
}

type UserInfo struct {
	FullName string
	UserID   string
}
