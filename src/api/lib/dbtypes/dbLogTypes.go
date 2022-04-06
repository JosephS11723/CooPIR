package dbtypes

import (
	"bytes"
	"encoding/json"
)

//log level number
type ErrorLevel int32

const (
	Debug ErrorLevel = iota
	Info
	Warning
	Error
	Critical
)

func (s ErrorLevel) String() string {

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

func (s ErrorLevel) toInt32(lvl string) ErrorLevel {

	var return_val ErrorLevel

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
func (s ErrorLevel) MarshalJSON() ([]byte, error) {
	buffer := bytes.NewBufferString(`"`)
	buffer.WriteString(s.String())
	buffer.WriteString(`"`)
	return buffer.Bytes(), nil
}

// UnmarshalJSON unmashals a quoted json string to the enum value
func (s *ErrorLevel) UnmarshalJSON(b []byte) error {
	var j string
	err := json.Unmarshal(b, &j)
	if err != nil {
		return err
	}
	// Note that if the string cannot be found then it will be set to the zero value, 'Created' in this case.
	*s = s.toInt32(j)

	return nil
}

//-------------------------------------------------------------------------------------------------//

type CaseLog struct {
	Case    string      `json:"case"`
	Level   ErrorLevel  `json:"level"`
	Time    string      `json:"time"`
	Content interface{} `json:"content"`
}

type UserLog struct {
	User    string      `json:"case"`
	Level   ErrorLevel  `json:"level"`
	Time    string      `json:"time"`
	Content interface{} `json:"content"`
}

type FileLog struct {
	File    string      `json:"file"`
	Case    string      `json:"case"`
	Level   ErrorLevel  `json:"level"`
	Time    string      `json:"time"`
	Content interface{} `json:"content"`
}

type AuthLog struct {
	User    string      `json:"user"`
	IP      string      `json:"ip"`
	Level   ErrorLevel  `json:"level"`
	Time    string      `json:"time"`
	Content interface{} `json:"content"`
}
