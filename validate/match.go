package validate

import (
	"fmt"
	"regexp"
)

type Match struct {
	Regexp *regexp.Regexp
}

func (m Match) IsSatisfied(obj interface{}) bool {
	str := obj.(string)
	return m.Regexp.MatchString(str)
}

func (m Match) DefaultMessage() string {
	return fmt.Sprint("Must match", m.Regexp)
}

var emailPattern = regexp.MustCompile("[\\w!#$%&'*+/=?^_`{|}~-]+(?:\\.[\\w!#$%&'*+/=?^_`{|}~-]+)*@(?:[\\w](?:[\\w-]*[\\w])?\\.)+[a-zA-Z0-9](?:[\\w-]*[\\w])?")

type Email struct {
	Match
}

func (e Email) DefaultMessage() string {
	return fmt.Sprint("Must be a valid email address")
}
var phonePattern= regexp.MustCompile("^[1][3-8]\\d{9}$")
type Phone struct{
	Match
}
func (e Phone) DefaultMessage() string {
	return fmt.Sprint("Must be a valid phone")
}

