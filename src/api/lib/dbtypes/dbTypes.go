package dbtypes

type User struct {
	Name       string   `json:"name"`
	Email      string   `json:"email"`
	Role       string   `json:"role"`
	Cases      []string `json:"cases"`
	SaltedHash string   `json:"saltedHash"`
}

type Case struct {
	Name          string   `json:"name"`
	Date_created  string   `json:"dateCreated"`
	View_access   string   `json:"viewAccess"`
	Edit_access   string   `json:"editAccess"`
	Collaborators []string `json:"collabs"`
}

type File struct {
	Hash        string `json:"hash"`
	Filename    string `json:"filename"`
	Case        string `json:"case"`
	File_dir    string `json:"fileDir"`
	Upload_date string `json:"uploadDate"`
	View_access string `json:"viewAccess"`
	Edit_access string `json:"editAccess"`
}

type Access struct {
	Filename string `json:"filename"`
	User     string `json:"user"`
	Date     string `json:"date"`
}
