package dbtypes

type Client struct {
	UUID  	string
	Name  	string
	Email 	string
	Role  	string
	Cases 	[]Case
}

type Case struct {
}

type Authentication struct {
	UUID 	string
	Salt 	string
	Pass 	string
}