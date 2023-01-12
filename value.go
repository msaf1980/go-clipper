package clipper

import (
	"bytes"
	"encoding/csv"
	"strings"
)

// Value is the interface to the dynamic value stored in a flag.
// (The default value is represented as a string.)
type Value interface {
	String() string
	Set(s string, doAppend bool) error
	Reset(interface{})
	Type() string
	Get() interface{}
}

type Arg interface {
	Value

	SetMaxLen(max int) Arg
	MaxLen() int
}

// None is disabled value (can't set).
// (The default value is represented as a string.)
type None struct{}

func (None) String() string {
	return ""
}

func (None) Set(v string, _ bool) error {
	return ErrorUnsupportedFlag{v}
}

func (None) Type() string {
	return "none"
}

func (None) Get() interface{} {
	return nil
}

func (None) Reset(interface{}) {}

func (n None) SetMaxLen(max int) Arg {
	return n
}

func (None) MaxLen() int {
	return -1
}

func readAsCSV(val string) ([]string, error) {
	if val == "" {
		return []string{}, nil
	}
	stringReader := strings.NewReader(val)
	csvReader := csv.NewReader(stringReader)
	return csvReader.Read()
}

func writeAsCSV(vals []string) (string, error) {
	b := bytes.Buffer{}
	w := csv.NewWriter(&b)
	err := w.Write(vals)
	if err != nil {
		return "", err
	}
	w.Flush()
	return strings.TrimSuffix(b.String(), "\n"), nil
}

func validateByValues(v string, validValues map[string]bool) (exist bool) {
	_, exist = validValues[v]
	return exist
}
