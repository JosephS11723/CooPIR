package dbtypes

import (
	"bytes"
	"encoding/json"
)

//log level number
type Level int32

const (
	Critical Level = iota
	Error
	Warning
	Info
	Debug
)

func (s Level) String() string {

	var return_val string

	switch s {
	case Critical:
		return_val = "Critical"
	case Error:
		return_val = "Error"
	case Warning:
		return_val = "Warning"
	case Info:
		return_val = "Info"
	case Debug:
		return_val = "Debug"
	}
	return return_val
}

func (s Level) ID(lvl string) Level {

	var return_val Level

	switch lvl {
	case "Critical":
		return_val = Critical
	case "Error":
		return_val = Error
	case "Warning":
		return_val = Warning
	case "Info":
		return_val = Info
	case "Debug":
		return_val = Debug
	}
	return return_val
}

// MarshalJSON marshals the enum as a quoted json string
func (s Level) MarshalJSON() ([]byte, error) {
	buffer := bytes.NewBufferString(`"`)
	buffer.WriteString(s.String())
	buffer.WriteString(`"`)
	return buffer.Bytes(), nil
}

// UnmarshalJSON unmashals a quoted json string to the enum value
func (s *Level) UnmarshalJSON(b []byte) error {
	var j string
	err := json.Unmarshal(b, &j)
	if err != nil {
		return err
	}
	// Note that if the string cannot be found then it will be set to the zero value, 'Created' in this case.
	*s = s.ID(j)

	return nil
}

//-------------------------------------------------------------------------------------------------//

type CaseLog struct {
	Case    string      `json:"case"`
	Level   Level       `json:"level"`
	Time    string      `json:"time"`
	Content interface{} `json:"content"`
}

type UserLog struct {
	User    string      `json:"case"`
	Level   Level       `json:"level"`
	Time    string      `json:"time"`
	Content interface{} `json:"content"`
}

type FileLog struct {
	File    string      `json:"file"`
	Case    string      `json:"case"`
	Level   Level       `json:"level"`
	Time    string      `json:"time"`
	Content interface{} `json:"content"`
}

type AuthLog struct {
	User    string      `json:"user"`
	IP      string      `json:"ip"`
	Level   Level       `json:"level"`
	Time    string      `json:"time"`
	Content interface{} `json:"content"`
}
