package dbtypes

//log level number
type AccessLevel int

const (
	Users AccessLevel = iota
	Responder
	Supervisor
	Admin
)

func (s AccessLevel) Int() int {

	var return_val int

	switch s {
	case Users:
		return_val = 1
	case Responder:
		return_val = 2
	case Supervisor:
		return_val = 3
	case Admin:
		return_val = 4
	}
	return return_val
}

func (s AccessLevel) ToInt(lvl string) AccessLevel {
	var return_val AccessLevel

	switch lvl {
	case "users":
		return_val = Users
	case "responder":
		return_val = Responder
	case "supervisor":
		return_val = Supervisor
	case "admin":
		return_val = Admin
	}
	return return_val
}

func (s AccessLevel) String() string {
	var return_val string

	switch s {
	case Users:
		return_val = "users"
	case Responder:
		return_val = "responder"
	case Supervisor:
		return_val = "supervisor"
	case Admin:
		return_val = "admin"
	}
	return return_val
}