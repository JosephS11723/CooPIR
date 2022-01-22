package dbtypes

type User struct {
	UUID  string
	Name  string
	Email string
	Role  string
	Cases []string
}

type Case struct {
	UUID         string
	Name         string
	Date_created string
	View_access  string
	Edit_access  string
	Colaborators []string
	Files        []File
}

type File struct {
	Hash        string
	Filename    string
	File_dir    string
	Upload_date string
	Last_access []Access
}

type Access struct {
	Date string
	User string
}

type Authentication struct {
	UUID string
	Salt string
	Pass string
}
