package dbtypes

type User struct {
	UUID  string           `json:"uuid"`
	Name  string           `json:"name"`
	Email string           `json:"email"`
	Role  string           `json:"role"`
	Cases []string         `json:"cases"`
	Auth  []Authentication `json:"auth"`
}

type Case struct {
	UUID         string		`json:"uuid"`
	Name         string		`json:"name"`
	Date_created string		`json:"dateCreated"`
	View_access  string		`json:"viewAccess"`
	Edit_access  string		`json:"editAccess"`
	Colaborators []string	`json:"collabs"`
	Files        []File		`json:"files"`
}

type File struct {
	Hash        string		`json:"hash"`
	Filename    string		`json:"filename"`
	File_dir    string		`json:"fileDir"`
	Upload_date string		`json:"uploadDate"`
	Last_access []Access	`json:"lastAccess"`
}

type Access struct {
	User string		`json:"user"`
	Date string		`json:"date"`
}

type Authentication struct {
	Salt string 	`json:"salt"`
	Pass string 	`json:"pass"`
}