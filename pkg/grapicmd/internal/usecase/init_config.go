package usecase

import "bytes"

type InitConfig struct {
	Revision string
	Branch   string
	Version  string
	HEAD     bool
	Package  string
}

func (c *InitConfig) BuildSpec() string {
	buf := bytes.NewBufferString("")
	var constraint string
	switch {
	case c.Revision != "":
		constraint = c.Revision
	case c.Branch != "":
		constraint = c.Branch
	case c.HEAD:
		constraint = "master"
	case c.Version != "":
		constraint = c.Version
	}
	if constraint != "" {
		buf.WriteString("@" + constraint)
	}
	return buf.String()
}
