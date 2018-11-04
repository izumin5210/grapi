package usecase

import "bytes"

type InitConfig struct {
	Source   string
	Revision string
	Branch   string
	Version  string
	HEAD     bool
	Package  string
}

func (c *InitConfig) BuildSpec() string {
	buf := bytes.NewBufferString("")
	if c.Source != "" {
		buf.WriteString(":" + c.Source)
	}
	var constraint string
	switch {
	case c.Revision != "":
		constraint = c.Revision
	case c.Branch != "":
		constraint = c.Branch
	case c.Version != "":
		constraint = c.Version
	case c.HEAD:
		constraint = "master"
	}
	if constraint != "" {
		buf.WriteString("@" + constraint)
	}
	return buf.String()
}
