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
	UUID          string   `json:"uuid"`
	Name          string   `json:"name"`
	Description   string   `json:"description" default:""`
	Date_created  string   `json:"dateCreated"`
	View_access   string   `json:"viewAccess"`
	Edit_access   string   `json:"editAccess"`
	Collaborators []string `json:"collabs"`
}

type File struct {
	UUID        string   `json:"uuid"`
	MD5         string   `json:"md5"`
	SHA1        string   `json:"sha1"`
	SHA256      string   `json:"sha256"`
	SHA512      string   `json:"sha512"`
	Tags        []string `json:"tags"`
	Filename    string   `json:"filename"`
	Case        string   `json:"case"`
	Upload_date string   `json:"uploadDate"`
	View_access string   `json:"viewAccess"`
	Edit_access string   `json:"editAccess"`
}

type Access struct {
	UUID     string `json:"uuid"`
	Caseuuid string `json:"caseuuid"`
	Target   string `json:"target"`
	User     string `json:"user"`
	Time     string `json:"time"`
}

type UpdateDoc struct {
	Filter map[string]interface{}
	Update map[string]interface{}
}

//just makes it easy to unmarshal a UUID-less struct as well
//as a struct with a password; User has saltedHash
type NewUser struct {
	Name     string   `json:"name"`
	Email    string   `json:"email"`
	Role     string   `json:"role"`
	Cases    []string `json:"cases"`
	Password string   `json:"password"`
}
