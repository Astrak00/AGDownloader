package types

const (
	Domain            = "aulaglobal.uc3m.es"
	Webservice        = "/webservice/rest/server.php"
	Service           = "aulaglobal_mobile"
	prompt_courses_en = "Select the courses you want to download\n"
	prompt_courses_es = "Selecciona los cursos que quieres descargar\n"
)

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
