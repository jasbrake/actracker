package pgnet

import (
	"database/sql/driver"
	"errors"
	"fmt"
	"net"
)

// INET is an net.IP that supports marshalling and scanning with a Postgresql INET column.
type INET struct {
	net.IP
}

func (i INET) Value() (driver.Value, error) {
	return i.String(), nil
}

func (i *INET) Scan(src interface{}) error {
	var t string
	switch src.(type) {
	case string:
		t = src.(string)
	case []byte:
		t = string(src.([]byte))
	default:
		return errors.New("incompatible type for IP")
	}
	ip := net.ParseIP(t)
	if ip == nil {
		return errors.New("unable to parse IP")
	}
	i.IP = ip
	return nil
}

func (i INET) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf("\"%s\"", i.String())), nil
}
