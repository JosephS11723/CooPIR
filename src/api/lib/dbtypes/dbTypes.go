package dbtypes

type User struct {
	UUID       string   `json:"uuid"`
	Name       string   `json:"name"`
	Email      string   `json:"email"`
	Role       string   `json:"role"`
	Cases      []string `json:"cases"`
	SaltedHash string   `json:"saltedHash"`
}

type Case struct {
	UUID         string   `json:"uuid"`
	Name         string   `json:"name"`
	Date_created string   `json:"dateCreated"`
	View_access  string   `json:"viewAccess"`
	Edit_access  string   `json:"editAccess"`
	Colaborators []string `json:"collabs"`
}

type File struct {
	Hash        string `json:"hash"`
	Filename    string `json:"filename"`
	File_dir    string `json:"fileDir"`
	Upload_date string `json:"uploadDate"`
}

type Access struct {
	Filesname string `json:"filename"`
	User      string `json:"user"`
	Date      string `json:"date"`
}
