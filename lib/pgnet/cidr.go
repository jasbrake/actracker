package pgnet

import (
	"database/sql/driver"
	"errors"
	"fmt"
	"net"
)

// CIDR is an net.IPNet that supports marshalling and scanning with a Postgresql CIDR column.
type CIDR struct {
	*net.IPNet
}

func (c CIDR) Value() (driver.Value, error) {
	return c.String(), nil
}

func (c *CIDR) Scan(src interface{}) error {
	var t string
	switch src.(type) {
	case string:
		t = src.(string)
	case []byte:
		t = string(src.([]byte))
	default:
		return errors.New("incompatible type for CIDR")
	}
	_, cidr, err := net.ParseCIDR(t)
	if err != nil {
		return err
	}
	c.IPNet = cidr
	return nil
}

func (c CIDR) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf("\"%s\"", c.String())), nil
}
